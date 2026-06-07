package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
)

type TenantSetting struct {
	Key   model.SettingsKey `db:"key"`
	Value string            `db:"value"`
}

func (c *Client) GetTenantSettings(ctx context.Context) (*model.TenantSettings, error) {
	var selectSettings []*TenantSetting

	err := c.SelectContext(ctx, &selectSettings, "SELECT key, value FROM tenant_settings")
	if err != nil {
		return nil, errors.Wrap(err, "failed to load tenant settings")
	}

	tenantSettings := model.NewDefaultTenantSettings()

	for _, setting := range selectSettings {
		switch setting.Key {
		case model.SettingsKeyScriptType:
			tenantSettings.ScriptType = setting.Value
		case model.SettingsKeyBlockAbusiveIPs:
			tenantSettings.BlockAbusiveIPs = setting.Value
		case model.SettingsKeyBlockTorExitNodes:
			tenantSettings.BlockTorExitNodes = setting.Value
		case model.SettingsKeyBlockedIPs:
			tenantSettings.BlockedIPs = setting.Value
		case model.SettingsKeyLanguage:
			// exhaustive:ignore
		}
	}

	return tenantSettings, nil
}

func (c *Client) UpdateTenantSettings(
	ctx context.Context,
	settings *db.UpdateTenantSettings,
) error {
	tx := c.MustBeginTx(ctx, nil)

	//nolint:exhaustive
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
			INSERT INTO tenant_settings (key, value, date_updated)
			VALUES (:key, :value, :date_updated)
			ON CONFLICT(key) DO UPDATE SET value=excluded.value, date_updated=excluded.date_updated`,
			map[string]any{
				"key":          key,
				"value":        value,
				"date_updated": time.Now().Unix(),
			},
		)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return errors.Wrap(
					errors.Join(err, rbErr),
					fmt.Sprintf("failed to persist %s setting", key),
				)
			}

			return errors.Wrap(err, fmt.Sprintf("failed to persist %s setting", key))
		}
	}

	return tx.Commit()
}
