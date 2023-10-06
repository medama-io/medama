package migrations

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/medama-io/medama/db/sqlite"
)

type Migration struct {
	ID   uint
	Name string
	Up   func(client *sqlite.Client) error
	Down func(client *sqlite.Client) error
}

type Service struct {
	client     *sqlite.Client
	migrations []*Migration
}

// CreateMigrationsTable creates the migrations table.
func CreateMigrationsTable(c *sqlite.Client) error {
	_, err := c.Exec("CREATE TABLE IF NOT EXISTS migrations (id INTEGER PRIMARY KEY, name VARCHAR(255) NOT NULL, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return err
	}

	return nil
}

// NewMigrationsService creates a new migrations service.
func NewMigrationsService(c *sqlite.Client) *Service {
	// Setup migration functions
	migrations := []*Migration{
		{ID: 1, Name: "0001_schema.go", Up: Up0001, Down: Down0001},
	}

	err := CreateMigrationsTable(c)
	if err != nil {
		slog.Error("failed to create migrations table", "error", err)
		return nil
	}
	slog.Debug("migrations table found")

	return &Service{
		client:     c,
		migrations: migrations,
	}
}

// AutoMigrate automatically migrates the schema, to keep your schema update to date.
func (s *Service) AutoMigrate() error {
	// Iterate over all migrations and check if they exist, else run
	for _, migration := range s.migrations {
		var id uint
		err := s.client.Get(&id, "SELECT id FROM migrations WHERE id = ?", migration.ID)

		switch {
		case errors.Is(err, sql.ErrNoRows): // Run migration
			slog.Warn("running migration, do not close the application", "id", migration.ID, "name", migration.Name)
			err = migration.Up(s.client)
			if err != nil {
				slog.Error("failed to run migration", "id", migration.ID, "name", migration.Name, "error", err)
				return err
			}

			// Insert migration into migrations table
			_, err = s.client.Exec("INSERT INTO migrations (id, name) VALUES (?, ?)", migration.ID, migration.Name)
			if err != nil {
				return err
			}

			slog.Info("migrated", "id", migration.ID, "name", migration.Name)

		case err == nil: // Migration already exists, skip
			slog.Debug("migration already exists", "id", migration.ID, "name", migration.Name)
			continue

		default:
			return err
		}
	}
	return nil
}
