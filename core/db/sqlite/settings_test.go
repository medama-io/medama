package sqlite_test

import (
	"testing"

	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/model"
	"github.com/stretchr/testify/require"
)

func TestGetSystemSettingsDefaults(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	settings, err := client.GetSystemSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, settings)

	defaults := model.NewDefaultSystemSettings()
	assert.Equal(defaults.ScriptType, settings.ScriptType)
	assert.Equal(defaults.BlockAbusiveIPs, settings.BlockAbusiveIPs)
	assert.Equal(defaults.BlockTorExitNodes, settings.BlockTorExitNodes)
	assert.Equal(defaults.BlockedIPs, settings.BlockedIPs)
}

func TestGetSystemSettingsCustomValues(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	for key, value := range map[model.SettingsKey]string{
		model.SettingsKeyScriptType:       "tagged-events",
		model.SettingsKeyBlockAbusiveIPs:   "false",
		model.SettingsKeyBlockTorExitNodes: "false",
		model.SettingsKeyBlockedIPs:        "10.0.0.1,10.0.0.2",
	} {
		_, err := client.ExecContext(ctx,
			"INSERT INTO system_settings (key, value, date_updated) VALUES (?, ?, ?)",
			key, value, 1)
		require.NoError(t, err)
	}

	settings, err := client.GetSystemSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, settings)

	assert.Equal("tagged-events", settings.ScriptType)
	assert.Equal("false", settings.BlockAbusiveIPs)
	assert.Equal("false", settings.BlockTorExitNodes)
	assert.Equal("10.0.0.1,10.0.0.2", settings.BlockedIPs)
}

func TestGetSystemSettingsPartialValues(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	_, err := client.ExecContext(ctx,
		"INSERT INTO system_settings (key, value, date_updated) VALUES (?, ?, ?)",
		model.SettingsKeyScriptType, "tagged-events", 1)
	require.NoError(t, err)

	_, err = client.ExecContext(ctx,
		"INSERT INTO system_settings (key, value, date_updated) VALUES (?, ?, ?)",
		model.SettingsKeyBlockedIPs, "10.0.0.1", 1)
	require.NoError(t, err)

	settings, err := client.GetSystemSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, settings)

	assert.Equal("tagged-events", settings.ScriptType)
	assert.Equal("10.0.0.1", settings.BlockedIPs)
	assert.Equal("true", settings.BlockAbusiveIPs)
	assert.Equal("true", settings.BlockTorExitNodes)
}

func TestUpdateSystemSettingsAll(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	scriptType := "tagged-events"
	blockAbusive := "false"
	blockTor := "false"
	blockedIPs := "10.0.0.1,10.0.0.2"

	err := client.UpdateSystemSettings(ctx, &sqlite.UpdateSystemSettings{
		ScriptType:        &scriptType,
		BlockAbusiveIPs:   &blockAbusive,
		BlockTorExitNodes: &blockTor,
		BlockedIPs:        &blockedIPs,
	})
	require.NoError(t, err)

	settings, err := client.GetSystemSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, settings)

	assert.Equal("tagged-events", settings.ScriptType)
	assert.Equal("false", settings.BlockAbusiveIPs)
	assert.Equal("false", settings.BlockTorExitNodes)
	assert.Equal("10.0.0.1,10.0.0.2", settings.BlockedIPs)
}

func TestUpdateSystemSettingsPartial(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	blAbusive := "false"
	blTor := "true"
	blIPs := "1.1.1.1"
	scType := "spa"

	err := client.UpdateSystemSettings(ctx, &sqlite.UpdateSystemSettings{
		ScriptType:        &scType,
		BlockAbusiveIPs:   &blAbusive,
		BlockTorExitNodes: &blTor,
		BlockedIPs:        &blIPs,
	})
	require.NoError(t, err)

	newScript := "tagged-events"
	newIPs := "2.2.2.2"
	err = client.UpdateSystemSettings(ctx, &sqlite.UpdateSystemSettings{
		ScriptType: &newScript,
		BlockedIPs: &newIPs,
	})
	require.NoError(t, err)

	settings, err := client.GetSystemSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, settings)

	assert.Equal("tagged-events", settings.ScriptType)
	assert.Equal("false", settings.BlockAbusiveIPs)
	assert.Equal("true", settings.BlockTorExitNodes)
	assert.Equal("2.2.2.2", settings.BlockedIPs)
}

func TestUpdateSystemSettingsEmpty(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	err := client.UpdateSystemSettings(ctx, &sqlite.UpdateSystemSettings{})
	require.NoError(t, err)

	settings, err := client.GetSystemSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, settings)

	defaults := model.NewDefaultSystemSettings()
	assert.Equal(defaults.ScriptType, settings.ScriptType)
	assert.Equal(defaults.BlockAbusiveIPs, settings.BlockAbusiveIPs)
	assert.Equal(defaults.BlockTorExitNodes, settings.BlockTorExitNodes)
	assert.Equal(defaults.BlockedIPs, settings.BlockedIPs)
}
