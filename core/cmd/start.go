package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
	"github.com/ogen-go/ogen/middleware"
)

type StartCommand struct {
	Debug       bool
	Server      ServerConfig
	AppDB       AppDBConfig
	AnalyticsDB AnalyticsDBConfig
}

// NewStartCommand creates a new start command.
func NewStartCommand() (*StartCommand, error) {
	serverConfig, err := NewServerConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create server config")
	}

	appConfig, err := NewAppDBConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create app db config")
	}

	analyticsConfig, err := NewAnalyticsDBConfig()
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
	fs := flag.NewFlagSet("start", flag.ContinueOnError)
	fs.BoolVar(&s.Debug, "debug", false, "Enable verbose debug logging")
	fs.Int64Var(&s.Server.Port, "port", DefaultPort, "Port to listen on")

	// Parse flags
	err := fs.Parse(args)
	if err != nil {
		return errors.Wrap(err, "failed to parse flags")
	}

	return nil
}

// Run executes the start command.
func (s *StartCommand) Run(ctx context.Context) error {
	util.SetupLogger(os.Stdout, s.Debug)
	slog.Info(GetVersion())

	// Setup database
	sqlite, err := sqlite.NewClient(s.AppDB.Host)
	if err != nil {
		return errors.Wrap(err, "failed to create sqlite client")
	}

	duckdb, err := duckdb.NewClient(s.AnalyticsDB.Host)
	if err != nil {
		return errors.Wrap(err, "failed to create duckdb client")
	}

	// Run migrations
	m := migrations.NewMigrationsService(ctx, sqlite, duckdb)
	if m == nil {
		return errors.New("could not create migrations service")
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
	service, err := services.NewService(auth, sqlite, duckdb)
	if err != nil {
		return errors.Wrap(err, "failed to create handlers")
	}

	authMiddleware := middlewares.NewAuthHandler(auth)
	mw := []middleware.Middleware{
		middlewares.RequestLogger(),
		middlewares.RequestContext(),
		middlewares.Recovery(),
	}
	h, err := api.NewServer(service,
		authMiddleware,
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
	mux.Handle("/", h)

	srv := &http.Server{
		Addr:         ":" + strconv.FormatInt(s.Server.Port, 10),
		Handler:      mux,
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

		slog.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(ctx, s.Server.TimeoutIdle)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Could not gracefully shutdown the server", "error", err)
		}

		close(closed)
	}()

	slog.Info(fmt.Sprintf("Starting server at http://localhost:%d", s.Server.Port))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("Could not listen on", "port", s.Server.Port, "error", err)
	}

	<-closed
	return nil
}
