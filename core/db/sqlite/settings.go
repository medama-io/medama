package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/model"
)

func (c *Client) GetSetting(ctx context.Context, key model.SettingsKey) (string, error) {
	name := "$." + string(key)
	query := `--sql
	SELECT
		JSON_EXTRACT(settings, ?) AS value
	FROM users LIMIT 1`

	var value string
	err := c.DB.GetContext(ctx, &value, query, name)
	if errors.Is(err, sql.ErrNoRows) {
		return "", model.ErrSettingNotFound
	}

	if err != nil {
		return "", errors.Wrap(err, "db")
	}
	return value, nil
}

func (c *Client) GetSettings(ctx context.Context) (*model.Settings, error) {
	query := `--sql
	SELECT
		JSON_EXTRACT(settings, '$.language') AS language,
		JSON_EXTRACT(settings, '$.script_type') AS script_type,
		JSON_EXTRACT(settings, '$.threads') AS threads,
		JSON_EXTRACT(settings, '$.memory_limit') AS memory_limit
	FROM users LIMIT 1`

	settings := &model.Settings{}
	err := c.DB.GetContext(ctx, settings, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrSettingNotFound
	}

	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	return settings, nil
}

func (c *Client) UpdateSetting(ctx context.Context, key model.SettingsKey, value string, dateUpdated int64) error {
	query := `--sql
    UPDATE users
    SET settings = JSON_SET(settings, :key, :value),
        date_updated = :date_updated`

	params := map[string]interface{}{
		"key":          "$." + string(key),
		"value":        value,
		"date_updated": dateUpdated,
	}

	result, err := c.DB.NamedExecContext(ctx, query, params)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "getting rows affected")
	}

	if rowsAffected == 0 {
		return model.ErrUserNotFound
	}

	return nil
}

// UpdateSettings updates a user's settings in the database.
func (c *Client) UpdateSettings(ctx context.Context, userID string, settings *model.Settings, dateUpdated int64) error {
	query := `--sql
    UPDATE users
    SET settings = :settings,
        date_updated = :date_updated
	WHERE id = :user_id`

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return errors.Wrap(err, "marshaling settings")
	}

	params := map[string]interface{}{
		"settings":     string(settingsJSON),
		"date_updated": dateUpdated,
		"user_id":      userID,
	}

	result, err := c.DB.NamedExecContext(ctx, query, params)
	if err != nil {
		return errors.Wrap(err, "db")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "getting rows affected")
	}

	if rowsAffected == 0 {
		return model.ErrUserNotFound
	}

	return nil
}
