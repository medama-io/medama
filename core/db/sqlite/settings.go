package sqlite

import (
	"context"
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

	err := c.SelectContext(
		ctx,
		&selectSettings,
		"SELECT key, value FROM system_settings",
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to load system settings")
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
