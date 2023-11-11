package migrations

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
)

type MigrationType string

const (
	SQLite MigrationType = "sqlite"
	DuckDB MigrationType = "duckdb"
)

type Migration[Client sqlite.Client | duckdb.Client] struct {
	ID   int
	Name string
	Type MigrationType
	Up   func(client *Client) error
	Down func(client *Client) error
}

type Service struct {
	duckdb           *duckdb.Client
	duckdbMigrations []*Migration[duckdb.Client]
	sqlite           *sqlite.Client
	sqliteMigrations []*Migration[sqlite.Client]
}

// CreateMigrationsTable creates the migrations table.
func CreateMigrationsTable(c *sqlite.Client) error {
	exec := `--sql
	CREATE TABLE IF NOT EXISTS migrations (id INTEGER PRIMARY KEY, name VARCHAR(255) NOT NULL, type VARCHAR(255) NOT NULL, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`
	_, err := c.Exec(exec)
	if err != nil {
		return err
	}

	return nil
}

// NewMigrationsService creates a new migrations service.
func NewMigrationsService(ctx context.Context, sqliteC *sqlite.Client, duckdbC *duckdb.Client) *Service {
	// Setup migration functions
	sqliteMigrations := []*Migration[sqlite.Client]{
		{ID: 1, Name: "0001_sqlite_schema.go", Type: SQLite, Up: Up0001, Down: Down0001},
	}

	duckdbMigrations := []*Migration[duckdb.Client]{}

	err := CreateMigrationsTable(sqliteC)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create migrations table", "error", err)
		return nil
	}
	slog.DebugContext(ctx, "migrations table found")

	return &Service{
		duckdb:           duckdbC,
		duckdbMigrations: duckdbMigrations,
		sqlite:           sqliteC,
		sqliteMigrations: sqliteMigrations,
	}
}

// runMigrator uses the given client to run the migrations.
func runMigrator[Client sqlite.Client | duckdb.Client](ctx context.Context, sqlite *sqlite.Client, client *Client, migrations []*Migration[Client]) error {
	// Iterate over all migrations and check if they exist, else run
	for _, migration := range migrations {
		var id int
		err := sqlite.GetContext(ctx, &id, "SELECT id FROM migrations WHERE id = ?", migration.ID)

		attributes := []slog.Attr{
			slog.Int("id", migration.ID),
			slog.String("name", migration.Name),
			slog.String("type", string(migration.Type)),
		}
		switch {
		// Run migration if it does not exist in migrations table
		case errors.Is(err, sql.ErrNoRows):
			slog.LogAttrs(ctx, slog.LevelWarn, "running migration, do not close the application", attributes...)

			// Run migration
			err = migration.Up(client)

			if err != nil {
				attributes = append(attributes, slog.String("error", err.Error()))
				slog.LogAttrs(ctx, slog.LevelError, "failed to run migration", attributes...)
				return err
			}

			// Insert migration into migrations table
			exec := `--sql
			INSERT INTO migrations (id, name, type) VALUES (?, ?, ?)`
			_, err = sqlite.ExecContext(ctx, exec, migration.ID, migration.Name, migration.Type)
			if err != nil {
				attributes = append(attributes, slog.String("error", err.Error()))
				slog.LogAttrs(ctx, slog.LevelError, "failed to insert migration into migrations table", attributes...)
				return err
			}

			slog.LogAttrs(ctx, slog.LevelInfo, "migrated", attributes...)

		case err == nil: // Migration already exists, skip
			slog.LogAttrs(ctx, slog.LevelDebug, "migration already exists", attributes...)
			continue

		default:
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to check if migration exists", attributes...)
			return err
		}
	}

	return nil
}

// AutoMigrate automatically migrates the schema, to keep your schema update to date.
func (s *Service) AutoMigrate(ctx context.Context) error {
	// SQLite
	err := runMigrator(ctx, s.sqlite, s.sqlite, s.sqliteMigrations)
	if err != nil {
		return err
	}

	// DuckDB
	err = runMigrator(ctx, s.sqlite, s.duckdb, s.duckdbMigrations)
	if err != nil {
		return err
	}

	return nil
}
