package migrations

import (
	"context"
	"database/sql"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/util/logger"
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
	CREATE TABLE IF NOT EXISTS migrations (id INTEGER PRIMARY KEY, name TEXT NOT NULL, type TEXT NOT NULL, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`
	_, err := c.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "migration")
	}

	return nil
}

// NewMigrationsService creates a new migrations service.
func NewMigrationsService(ctx context.Context, sqliteC *sqlite.Client, duckdbC *duckdb.Client) (*Service, error) {
	// Setup migration functions
	sqliteMigrations := []*Migration[sqlite.Client]{
		{ID: 1, Name: "0001_sqlite_schema.go", Type: SQLite, Up: Up0001, Down: Down0001},
	}

	duckdbMigrations := []*Migration[duckdb.Client]{
		{ID: 2, Name: "0002_duckdb_schema.go", Type: DuckDB, Up: Up0002, Down: Down0002},
		{ID: 3, Name: "0003_duckdb_referrer.go", Type: DuckDB, Up: Up0003, Down: Down0003},
		{ID: 4, Name: "0004_duckdb_events.go", Type: DuckDB, Up: Up0004, Down: Down0004},
	}

	log := logger.Get()
	err := CreateMigrationsTable(sqliteC)
	if err != nil {
		log.Error().Err(err).Msg("failed to create migrations table")
		return nil, err
	}
	log.Debug().Msg("migrations table found")

	return &Service{
		duckdb:           duckdbC,
		duckdbMigrations: duckdbMigrations,
		sqlite:           sqliteC,
		sqliteMigrations: sqliteMigrations,
	}, nil
}

// runMigrator uses the given client to run the migrations.
func runMigrator[Client sqlite.Client | duckdb.Client](ctx context.Context, sqlite *sqlite.Client, client *Client, migrations []*Migration[Client]) error {
	// Iterate over all migrations and check if they exist, else run
	for _, migration := range migrations {
		var id int
		err := sqlite.GetContext(ctx, &id, "SELECT id FROM migrations WHERE id = ?", migration.ID)

		log := logger.Get().With().
			Int("id", migration.ID).
			Str("name", migration.Name).
			Str("type", string(migration.Type)).
			Logger()

		switch {
		// Run migration if it does not exist in migrations table
		case errors.Is(err, sql.ErrNoRows):
			log.Warn().
				Msg("running migration, do not close the application")

			// Run migration
			err = migration.Up(client)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to run migration")

				return err
			}

			// Insert migration into migrations table
			exec := `--sql
			INSERT INTO migrations (id, name, type) VALUES (?, ?, ?)`
			_, err = sqlite.ExecContext(ctx, exec, migration.ID, migration.Name, migration.Type)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to insert migration into migrations table")

				return errors.Wrap(err, "migration")
			}

			log.Info().Msg("migrated")

		case err == nil: // Migration already exists, skip
			log.Debug().Msg("migration already exists")
			continue

		default:
			log.Error().
				Err(err).
				Msg("failed to check if migration exists")
			return errors.Wrap(err, "migration")
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
