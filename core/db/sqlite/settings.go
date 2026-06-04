package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/model"
)

type SystemSetting struct {
	Key   model.SettingsKey
	Value string
}

type UpdateSystemSettings struct {
	ScriptType        *string
	BlockAbusiveIPs   *string
	BlockTorExitNodes *string
	BlockedIPs        *string
}

func (c *Client) GetSystemSettings(ctx context.Context) (*model.SystemSettings, error) {
	var selectSettings []*SystemSetting

	settingKeysToLoad := fmt.Sprintf(
		"%s,%s,%s,%s",
		model.SettingsKeyScriptType,
		model.SettingsKeyBlockAbusiveIPs,
		model.SettingsKeyBlockTorExitNodes,
		model.SettingsKeyBlockedIPs,
	)

	err := c.SelectContext(
		ctx,
		&selectSettings,
		"SELECT key, value FROM system_settings WHERE key IN (?)",
		settingKeysToLoad,
	)

	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	systemSettings := model.NewDefaultSystemSettings()

	for _, setting := range selectSettings {
		switch setting.Key {
		case model.SettingsKeyScriptType:
			systemSettings.ScriptType = setting.Value
		case model.SettingsKeyBlockAbusiveIPs:
			systemSettings.BlockAbusiveIPs = setting.Value
		case model.SettingsKeyBlockTorExitNodes:
			systemSettings.BlockTorExitNodes = setting.Value
		case model.SettingsKeyBlockedIPs:
			systemSettings.BlockedIPs = setting.Value
		}
	}

	return systemSettings, nil
}

func (c *Client) UpdateSystemSettings(
	ctx context.Context,
	settings *UpdateSystemSettings,
) error {
	tx := c.DB.MustBeginTx(ctx, nil)

	propertiesToUpdate := map[model.SettingsKey]*string{
		model.SettingsKeyScriptType:        settings.ScriptType,
		model.SettingsKeyBlockAbusiveIPs:   settings.BlockAbusiveIPs,
		model.SettingsKeyBlockTorExitNodes: settings.BlockTorExitNodes,
		model.SettingsKeyBlockedIPs:        settings.BlockedIPs,
	}

	for key, value := range propertiesToUpdate {
		if value == nil {
			continue
		}

		_, err := tx.NamedExecContext(
			ctx,
			`--sql
			INSERT INTO system_settings (key, value, date_updated)
			VALUES (:key, :value, :date_updated)
			ON CONFLICT(key) DO UPDATE SET value=excluded.value, date_updated=excluded.date_updated`,
			map[string]any{
				"key":          key,
				"value":        value,
				"date_updated": time.Now().Unix(),
			},
		)

		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, fmt.Sprintf("failed to persist %s setting", key))
		}
	}

	return tx.Commit()
}

// ===

func (c *Client) GetSetting(ctx context.Context, key model.SettingsKey) (string, error) {
	name := "$." + string(key)
	query := `
    SELECT
        JSON_EXTRACT(settings, ?) AS value
    FROM users
    WHERE JSON_EXTRACT(settings, ?) IS NOT NULL
    LIMIT 1`

	var value sql.NullString

	err := c.GetContext(ctx, &value, query, name, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", model.ErrSettingNotFound
		}

		return "", errors.Wrap(err, "db")
	}

	if !value.Valid {
		return "", model.ErrSettingNotFound
	}

	return value.String, nil
}

func (c *Client) GetSettings(ctx context.Context) (*model.UserSettings, error) {
	query := `--sql
	SELECT
		COALESCE(JSON_EXTRACT(settings, '$.language'), 'en') AS language,
		COALESCE(JSON_EXTRACT(settings, '$.script_type'), 'default') AS script_type,
		COALESCE(JSON_EXTRACT(settings, '$.block_abusive_ips'), 'true') AS block_abusive_ips,
		COALESCE(JSON_EXTRACT(settings, '$.block_tor_exit_nodes'), 'true') AS block_tor_exit_nodes,
		COALESCE(JSON_EXTRACT(settings, '$.blocked_ips'), '') AS blocked_ips
	FROM users LIMIT 1`

	settings := model.NewDefaultUserSettings()

	err := c.GetContext(ctx, settings, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrSettingNotFound
	}

	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	return settings, nil
}

func (c *Client) UpdateSetting(ctx context.Context, key model.SettingsKey, value string) error {
	query := `--sql
    UPDATE users
    SET settings = JSON_SET(settings, :key, :value),
        date_updated = :date_updated`

	params := map[string]any{
		"key":          "$." + string(key),
		"value":        value,
		dateUpdatedKey: time.Now().Unix(),
	}

	result, err := c.NamedExecContext(ctx, query, params)
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
func (c *Client) UpdateSettings(
	ctx context.Context,
	userID string,
	settings *model.UserSettings,
) error {
	query := `--sql
    UPDATE users
    SET settings = :settings,
        date_updated = :date_updated
	WHERE id = :user_id`

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return errors.Wrap(err, "marshaling settings")
	}

	params := map[string]any{
		dateUpdatedKey: time.Now().Unix(),
		"settings":     string(settingsJSON),
		"user_id":      userID,
	}

	result, err := c.NamedExecContext(ctx, query, params)
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
