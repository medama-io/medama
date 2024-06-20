package main

import (
	"context"
	"flag"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"

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
)

type StartCommand struct {
	Server      ServerConfig
	AppDB       AppDBConfig
	AnalyticsDB AnalyticsDBConfig
}

// NewStartCommand creates a new start command.
func NewStartCommand(useEnv bool) (*StartCommand, error) {
	serverConfig, err := NewServerConfig(useEnv)
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
	fs.StringVar(&s.Server.Level, "level", s.Server.Level, "Logger level (debug, info, warn, error)")

	// Database settings.
	fs.StringVar(&s.AppDB.Host, "appdb", s.AppDB.Host, "Path to app database.")
	fs.StringVar(&s.AnalyticsDB.Host, "analyticsdb", s.AnalyticsDB.Host, "Path to analytics database.")

	// Misc settings.
	fs.BoolVar(&s.Server.UseEnvironment, "env", false, "Opt-in to allow environment variables to be used for configuration. Flags will still override environment variables.")

	// Handle array type flags.
	corsAllowedOrigins := fs.String("corsorigins", strings.Join(s.Server.CORSAllowedOrigins, ","), "Comma separated list of allowed CORS origins on API routes. Useful for external dashboards that may host the frontend on a different domain.")

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
	auth, err := util.NewAuthService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create auth service")
	}

	// Setup handlers
	service, err := services.NewService(ctx, auth, sqlite, duckdb)
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

	// SPA client. We need to serve index.html to all routes that are not /api.
	client, err := generate.SPAClient()
	if err != nil {
		return errors.Wrap(err, "failed to create spa client")
	}
	clientServer := http.FileServer(http.FS(client))

	// Read index.html and tracker script once during initialisation.
	indexFile, err := readFile(client, "index.html")
	if err != nil {
		return errors.Wrap(err, "could not read index.html")
	}
	trackerFile, err := readFile(client, "medama.js")
	if err != nil {
		return errors.Wrap(err, "could not read medama.js")
	}

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uPath := path.Clean(r.URL.Path)
		if strings.HasPrefix(uPath, "/assets/") || strings.HasPrefix(uPath, "/favicon.ico") || strings.HasPrefix(uPath, "/manifest") {
			clientServer.ServeHTTP(w, r)
			return
		}

		// Serve tracker script
		if uPath == "/script.js" {
			w.Header().Set("Content-Type", "application/javascript")
			_, err := w.Write(trackerFile)
			if err != nil {
				http.Error(w, "could not serve tracker script", http.StatusInternalServerError)
			}
			return
		}

		// Serve index.html for any other path
		w.Header().Set("Content-Type", "text/html")
		_, err := w.Write(indexFile)
		if err != nil {
			http.Error(w, "could not serve index.html", http.StatusInternalServerError)
		}
	}))

	// Apply custom CORS middleware to the mux handler
	handler := middlewares.CORSAllowedOriginsMiddleware(s.Server.CORSAllowedOrigins)(mux)

	srv := &http.Server{
		Addr:         ":" + strconv.FormatInt(s.Server.Port, 10),
		Handler:      handler,
		ReadTimeout:  s.Server.TimeoutRead,
		WriteTimeout: s.Server.TimeoutWrite,
		IdleTimeout:  s.Server.TimeoutIdle,
	}

	// Graceful shutdown
	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info().Msg("Shutting down server...")

		ctx, cancel := context.WithTimeout(ctx, s.Server.TimeoutIdle)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Could not gracefully shutdown the server")
		}

		close(closed)
	}()

	log.Info().Msgf("Starting server at http://localhost:%d", s.Server.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Panic().Err(err).Msgf("Could not listen on port: %d", s.Server.Port)
	}

	<-closed
	return nil
}

func readFile(filesystem fs.FS, file string) ([]byte, error) {
	f, err := filesystem.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
