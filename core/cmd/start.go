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

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/middlewares"
	"github.com/medama-io/medama/migrations"
	"github.com/medama-io/medama/services"
	"github.com/medama-io/medama/util"
	"github.com/ogen-go/ogen/middleware"
)

type StartCommand struct {
	Debug    bool
	Server   ServerConfig
	Database DatabaseConfig
}

// NewStartCommand creates a new start command.
func NewStartCommand() *StartCommand {
	return &StartCommand{
		Server: ServerConfig{
			AppEnv:               AppEnvDevelopment,
			CacheCleanupInterval: DefaultCacheCleanupInterval,
			TimeoutRead:          DefaultTimeoutRead,
			TimeoutWrite:         DefaultTimeoutWrite,
			TimeoutIdle:          DefaultTimeoutIdle,
		},
		Database: DatabaseConfig{
			Host: DefaultDatabaseHost,
		},
	}
}

// ParseFlags parses the command line flags for the start command.
func (s *StartCommand) ParseFlags(args []string) error {
	fs := flag.NewFlagSet("start", flag.ContinueOnError)
	fs.BoolVar(&s.Debug, "debug", false, "Enable verbose debug logging")
	fs.Int64Var(&s.Server.Port, "port", DefaultPort, "Port to listen on")

	// Parse flags
	return fs.Parse(args)
}

// Run executes the start command.
func (s *StartCommand) Run(ctx context.Context) error {
	util.SetupLogger(os.Stdout, s.Debug)
	slog.Info(GetVersion())

	// Setup database
	db, err := sqlite.NewClient(s.Database.Host)
	if err != nil {
		return err
	}

	// Run migrations
	m := migrations.NewMigrationsService(ctx, db)
	if m == nil {
		slog.Error("Could not create migrations service")
		return err
	}
	err = m.AutoMigrate(ctx)
	if err != nil {
		slog.Error("Could not run migrations", "error", err)
		return err
	}

	// Setup handlers
	service := services.NewService(db)
	mw := []middleware.Middleware{
		middlewares.Recovery(),
		middlewares.RequestLogger(),
	}
	h, err := api.NewServer(service,
		api.WithMiddleware(mw...),
		api.WithErrorHandler(middlewares.ErrorHandler),
		api.WithNotFound(middlewares.NotFound()),
	)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:         ":" + strconv.FormatInt(s.Server.Port, 10),
		Handler:      h,
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
