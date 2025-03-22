package sqlite

import (
	"context"
	"database/sql"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
	"github.com/ncruces/go-sqlite3"
)

func (c *Client) CreateWebsite(ctx context.Context, website *model.Website) error {
	exec := `--sql
	INSERT INTO websites (
		user_id,
		hostname,
		date_created,
		date_updated
	) VALUES (
		:user_id,
		:hostname,
		:date_created,
		:date_updated
	)`

	paramMap := map[string]any{
		"user_id":      website.UserID,
		"hostname":     website.Hostname,
		"date_created": website.DateCreated,
		"date_updated": website.DateUpdated,
	}

	_, err := c.DB.NamedExecContext(ctx, exec, paramMap)
	if err != nil {
		switch {
		case errors.Is(err, sqlite3.CONSTRAINT_PRIMARYKEY):
			return model.ErrWebsiteExists
		case errors.Is(err, sqlite3.CONSTRAINT_FOREIGNKEY):
			return model.ErrUserNotFound
		}

		log := logger.Get()
		log.Error().
			Str("user_id", website.UserID).
			Str("hostname", website.Hostname).
			Int64("date_created", website.DateCreated).
			Int64("date_updated", website.DateUpdated).
			Err(err).
			Msg("failed to create website")

		return errors.Wrap(err, "db")
	}

	return nil
}

func (c *Client) ListWebsites(ctx context.Context, userID string) ([]*model.Website, error) {
	var websites []*model.Website
	log := logger.Get()

	query := `--sql
	SELECT user_id, hostname, date_created, date_updated FROM websites WHERE user_id = ?`

	err := c.SelectContext(ctx, &websites, query, userID)
	if err != nil {
		log.Error().
			Str("user_id", userID).
			Err(err).
			Msg("failed to list websites")

		return nil, errors.Wrap(err, "db")
	}

	if len(websites) == 0 {
		log.Debug().Str("user_id", userID).Msg("no websites found")
		// Return empty slice instead of nil
		return []*model.Website{}, nil
	}

	return websites, nil
}

func (c *Client) ListAllHostnames(ctx context.Context) ([]string, error) {
	query := `--sql
	SELECT hostname FROM websites`

	var hostnames []string
	err := c.SelectContext(ctx, &hostnames, query)
	if err != nil {
		log := logger.Get()
		log.Error().Err(err).Msg("failed to list all hostnames")
		return nil, errors.Wrap(err, "db")
	}

	return hostnames, nil
}

func (c *Client) UpdateWebsite(ctx context.Context, website *model.Website) error {
	log := logger.Get()
	// Update all values except user_id
	exec := `--sql
	UPDATE websites SET hostname = :hostname, date_updated = :date_updated WHERE hostname = :hostname`

	paramMap := map[string]any{
		"hostname":     website.Hostname,
		"date_updated": website.DateUpdated,
	}

	res, err := c.DB.NamedExecContext(ctx, exec, paramMap)
	if err != nil {
		if errors.Is(err, sqlite3.CONSTRAINT_PRIMARYKEY) {
			return model.ErrWebsiteExists
		}

		log.Error().
			Str("hostname", website.Hostname).
			Int64("date_updated", website.DateUpdated).
			Err(err).
			Msg("failed to update website")

		return errors.Wrap(err, "db")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error().
			Str("hostname", website.Hostname).
			Int64("date_updated", website.DateUpdated).
			Err(err).
			Msg("failed to get rows affected")

		return errors.Wrap(err, "db")
	}

	if rowsAffected == 0 {
		log.Debug().Str("hostname", website.Hostname).Msg("website not found")
		return model.ErrWebsiteNotFound
	}

	return nil
}

func (c *Client) GetWebsite(ctx context.Context, hostname string) (*model.Website, error) {
	var website model.Website
	log := logger.Get()

	query := `--sql
	SELECT user_id, hostname, date_created, date_updated FROM websites WHERE hostname = ?`

	err := c.QueryRowxContext(ctx, query, hostname).StructScan(&website)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().Str("hostname", hostname).Msg("website not found")
			return nil, model.ErrWebsiteNotFound
		}

		log.Error().Str("hostname", hostname).Err(err).Msg("failed to get website")

		return nil, errors.Wrap(err, "db")
	}

	return &website, nil
}

func (c *Client) DeleteWebsite(ctx context.Context, hostname string) error {
	log := logger.Get()
	exec := `--sql
	DELETE FROM websites WHERE hostname = ?`

	res, err := c.DB.ExecContext(ctx, exec, hostname)
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Err(err).
			Msg("failed to delete website")

		return errors.Wrap(err, "db")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Err(err).
			Msg("failed to get rows affected")

		return errors.Wrap(err, "db")
	}

	if rowsAffected == 0 {
		log.Debug().Str("hostname", hostname).Msg("website not found")
		return model.ErrWebsiteNotFound
	}

	return nil
}
