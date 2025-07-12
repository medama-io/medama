package main

import (
	"context"
	"crypto/tls"
	"flag"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/go-faster/errors"
	generate "github.com/medama-io/medama"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/middlewares"
	"github.com/medama-io/medama/migrations"
	"github.com/medama-io/medama/services"
	"github.com/medama-io/medama/util"
	"github.com/medama-io/medama/util/logger"
	"github.com/ogen-go/ogen/middleware"
	"github.com/rs/zerolog"
)

const (
	HTTPPort        = 80
	HTTPSPort       = 443
	shutdownTimeout = 5 * time.Second
)

type StartCommand struct {
	Server      ServerConfig
	AppDB       AppDBConfig
	AnalyticsDB AnalyticsDBConfig
}

// NewStartCommand creates a new start command.
func NewStartCommand(useEnv bool, version string, commit string) (*StartCommand, error) {
	serverConfig, err := NewServerConfig(useEnv, version, commit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create server config")
	}

	appConfig, err := NewAppDBConfig(useEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create app db config")
	}

	analyticsConfig, err := NewAnalyticsDBConfig(useEnv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create analytics db config")
	}

	return &StartCommand{
		Server:      *serverConfig,
		AppDB:       *appConfig,
		AnalyticsDB: *analyticsConfig,
	}, nil
}

// ParseFlags parses the command line flags for the start command.
func (s *StartCommand) ParseFlags(args []string) error {
	fs := flag.NewFlagSet("start", flag.ExitOnError)

	// General settings.
	fs.Int64Var(&s.Server.Port, "port", s.Server.Port, "Port to listen on.")
	fs.StringVar(&s.Server.Logger, "logger", s.Server.Logger, "Logger format (json, pretty)")
	fs.StringVar(
		&s.Server.Level,
		"level",
		s.Server.Level,
		"Logger level (debug, info, warn, error)",
	)

	// AutoSSL settings.
	fs.StringVar(
		&s.Server.AutoSSLDomain,
		"autossl",
		s.Server.AutoSSLDomain,
		"Automatically provision SSL certificates for the specified domain and redirect HTTP requests to HTTPS.\n\nRequires the server to run on ports 80 and 443. The domain must be publicly accessible and resolve to the server.",
	)
	fs.StringVar(
		&s.Server.AutoSSLEmail,
		"autosslemail",
		s.Server.AutoSSLEmail,
		"Email address to optionally send SSL certificate notifications.",
	)

	// Database settings.
	fs.StringVar(&s.AppDB.Host, "appdb", s.AppDB.Host, "Path to app database.")
	fs.StringVar(
		&s.AnalyticsDB.Host,
		"analyticsdb",
		s.AnalyticsDB.Host,
		"Path to analytics database.",
	)

	// Misc settings.
	fs.BoolVar(&s.Server.Profiler, "profiler", s.Server.Profiler, "Enable debug profiling.")
	fs.BoolVar(
		&s.Server.UseEnvironment,
		"env",
		false,
		"Opt-in to allow environment variables to be used for configuration. Flags will still override environment variables.",
	)
	fs.BoolVar(
		&s.Server.DemoMode,
		"demo",
		s.Server.DemoMode,
		"Enable demo mode restricting all POST/PATCH/DELETE actions (except login).",
	)

	// Handle array type flags.
	corsAllowedOrigins := fs.String(
		"corsorigins",
		strings.Join(s.Server.CORSAllowedOrigins, ","),
		"Comma separated list of allowed CORS origins on API routes. Useful for external dashboards that may host the frontend on a different domain.",
	)

	// Parse flags.
	err := fs.Parse(args)
	if err != nil {
		return errors.Wrap(err, "failed to parse flags")
	}

	if *corsAllowedOrigins != "" {
		s.Server.CORSAllowedOrigins = strings.Split(*corsAllowedOrigins, ",")
	}

	return nil
}

// Run executes the start command.
func (s *StartCommand) Run(ctx context.Context) error {
	log, err := logger.Init(s.Server.Logger, s.Server.Level)
	if err != nil {
		return errors.Wrap(err, "failed to setup logger")
	}
	log.Info().Msg(GetVersion())
	log.Debug().Interface("config", s).Msg("")

	// Setup database
	sqlite, err := sqlite.NewClient(s.AppDB.Host)
	if err != nil {
		return errors.Wrap(err, "failed to create sqlite client")
	}
	defer sqlite.Close()

	duckdb, err := duckdb.NewClient(s.AnalyticsDB.Host)
	if err != nil {
		return errors.Wrap(err, "failed to create duckdb client")
	}
	defer duckdb.Close()

	// Run migrations
	m, err := migrations.NewMigrationsService(ctx, sqlite, duckdb)
	if err != nil {
		return errors.Wrap(err, "failed to create migrations service")
	}
	err = m.AutoMigrate(ctx)
	if err != nil {
		return errors.Wrap(err, "could not run migrations")
	}

	// Setup auth service
	auth, err := util.NewAuthService(ctx, s.Server.DemoMode)
	if err != nil {
		return errors.Wrap(err, "failed to create auth service")
	}

	// Setup handlers
	service, err := services.NewService(ctx, auth, sqlite, duckdb, s.Server.Commit)
	if err != nil {
		return errors.Wrap(err, "failed to create handlers")
	}

	mw := []middleware.Middleware{
		middlewares.RequestLogger(),
		middlewares.RequestContext(),
		middlewares.Recovery(),
	}
	apiHandler, err := api.NewServer(service,
		middlewares.NewAuthHandler(auth),
		api.WithMiddleware(mw...),
		api.WithErrorHandler(middlewares.ErrorHandler),
		api.WithNotFound(middlewares.NotFound()),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create server")
	}

	// We need to add additional static routes for the web app.
	mux := http.NewServeMux()
	mux.Handle("/openapi.yaml", http.FileServer(http.FS(generate.OpenAPIDocument)))
	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	// Start CPU profiling if enabled.
	if s.Server.Profiler {
		log.Warn().Msg("Enabling debug profiler...")
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	// SPA client.
	err = services.SetupAssetHandler(mux, service.RuntimeConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup asset handler")
	}

	// Apply custom CORS middleware to the mux handler
	handler := middlewares.CORSAllowedOriginsMiddleware(s.Server.CORSAllowedOrigins)(mux)

	// Compression middleware.
	handler = middlewares.Compress()(handler)

	// X-API-Commit header for client-side cache busting.
	handler = middlewares.XAPICommitMiddleware(s.Server.Commit)(handler)

	// RateLimiter middleware to limit requests with coarse IP prefixes. Ensure this is applied last to the handler chain.
	handler = middlewares.NewRateLimiter(handler)

	return s.serve(ctx, log, handler)
}

// serve starts the HTTP server with the given handler. If AutoSSL is enabled, it will also provision certificates
// and redirect HTTP to HTTPS.
func (s *StartCommand) serve(ctx context.Context, log zerolog.Logger, mux http.Handler) error {
	var (
		httpListener  net.Listener
		httpsListener net.Listener
		httpWg        sync.WaitGroup

		cfg *certmagic.Config
	)

	isSSL := s.Server.AutoSSLDomain != ""

	// If AutoSSL is enabled, we need to provision SSL certificates for the domain and redirect HTTP to HTTPS.
	// Otherwise, we can start the server as is.
	if isSSL {
		certmagic.DefaultACME.Agreed = true

		if s.Server.AutoSSLEmail != "" {
			certmagic.DefaultACME.Email = s.Server.AutoSSLEmail
		}

		cfg = certmagic.NewDefault()

		err := cfg.ManageSync(ctx, []string{s.Server.AutoSSLDomain})
		if err != nil {
			return errors.Wrap(err, "failed to provision ssl certificate")
		}
	}

	// Create HTTP/S listeners.
	httpPort := ":" + strconv.FormatInt(s.Server.Port, 10)
	if isSSL {
		httpPort = ":" + strconv.FormatInt(HTTPPort, 10)
	}

	httpListener, err := net.Listen("tcp", httpPort)
	if err != nil {
		return errors.Wrap(err, "failed to create http listener")
	}

	if isSSL {
		httpsPort := ":" + strconv.FormatInt(HTTPSPort, 10)

		tlsConfig := cfg.TLSConfig()
		tlsConfig.NextProtos = append([]string{"h2", "http/1.1"}, tlsConfig.NextProtos...)

		httpsListener, err = tls.Listen("tcp", httpsPort, tlsConfig)
		if err != nil {
			httpListener.Close()
			return errors.Wrap(err, "failed to create https listener")
		}
	}

	httpWg.Add(1)
	defer httpWg.Done()

	// Cleanup listeners when all servers are done.
	go func() {
		httpWg.Wait()
		httpListener.Close()
		if isSSL {
			httpsListener.Close()
		}
	}()

	var (
		httpServer  *http.Server
		httpsServer *http.Server
	)

	if !isSSL {
		httpServer = &http.Server{
			Addr:              ":" + strconv.FormatInt(s.Server.Port, 10),
			Handler:           mux,
			ReadHeaderTimeout: s.Server.TimeoutReadHeader,
			ReadTimeout:       s.Server.TimeoutRead,
			WriteTimeout:      s.Server.TimeoutWrite,
			IdleTimeout:       s.Server.TimeoutIdle,
			BaseContext:       func(_ net.Listener) context.Context { return ctx },
		}
	} else {
		// The HTTP server solves the ACME challenges and redirects to HTTPS.
		httpServer = &http.Server{
			ReadHeaderTimeout: shutdownTimeout,
			ReadTimeout:       shutdownTimeout,
			WriteTimeout:      shutdownTimeout,
			IdleTimeout:       shutdownTimeout,
			BaseContext:       func(_ net.Listener) context.Context { return ctx },
		}

		if len(cfg.Issuers) > 0 {
			if am, ok := cfg.Issuers[0].(*certmagic.ACMEIssuer); ok {
				httpServer.Handler = am.HTTPChallengeHandler(middlewares.HTTPSRedirectFunc())
			}
		}

		// The HTTPS server simply serves the user's handler.
		httpsServer = &http.Server{
			Handler:           mux,
			ReadHeaderTimeout: s.Server.TimeoutReadHeader,
			ReadTimeout:       s.Server.TimeoutRead,
			WriteTimeout:      s.Server.TimeoutWrite,
			IdleTimeout:       s.Server.TimeoutIdle,
			BaseContext:       func(_ net.Listener) context.Context { return ctx },
		}
	}

	// Graceful shutdown with signal handling.
	closed := make(chan struct{})

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		<-stop
		log.Info().Msg("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		if isSSL {
			if err := httpsServer.Shutdown(shutdownCtx); err != nil {
				log.Error().Err(err).Msg("Could not gracefully shutdown the HTTPS server")
			}
		}

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Could not gracefully shutdown the HTTP server")
		}

		close(closed)
	}()

	if !isSSL {
		log.Info().Msgf("Starting server at http://localhost:%d", s.Server.Port)

		if err := httpServer.Serve(httpListener); err != http.ErrServerClosed {
			return errors.Wrap(err, "could not listen on port")
		}

		<-closed
		return nil
	}

	log.Printf(
		"Starting server with automatic SSL for %s configured.\n\nServing HTTP->HTTPS on %s and %s",
		s.Server.AutoSSLDomain,
		httpListener.Addr(),
		httpsListener.Addr(),
	)

	go func() {
		if err := httpServer.Serve(httpListener); err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Could not listen on http port")
		}
	}()

	if err := httpsServer.Serve(httpsListener); err != http.ErrServerClosed {
		return errors.Wrap(err, "could not listen on https port")
	}

	<-closed
	return nil
}
