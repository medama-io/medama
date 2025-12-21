package sqlite_test

import (
	"testing"

	"github.com/medama-io/medama/model"
	"github.com/stretchr/testify/require"
)

func TestGetSetting(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	setting, err := client.GetSetting(ctx, model.SettingsKeyLanguage)
	require.NoError(t, err)
	assert.NotNil(setting)
	assert.Equal("en", setting)
}

func TestGetUnknownSetting(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	setting, err := client.GetSetting(ctx, "unknown")
	assert.ErrorIs(err, model.ErrSettingNotFound)
	assert.Empty(setting)
}

func TestGetSettings(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	settings, err := client.GetSettings(ctx)
	require.NoError(t, err)
	assert.NotNil(settings)
	assert.Equal("en", settings.Language)
	assert.Equal("default", settings.ScriptType)
}

func TestUpdateSetting(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	err := client.UpdateSetting(ctx, model.SettingsKeyLanguage, "jp")
	require.NoError(t, err)

	setting, err := client.GetSetting(ctx, model.SettingsKeyLanguage)
	require.NoError(t, err)
	assert.Equal("jp", setting)
}

func TestUpdateSettings(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	user, err := client.GetUserByUsername(ctx, "admin")
	require.NoError(t, err)

	err = client.UpdateSettings(ctx, user.ID, &model.UserSettings{
		Language:   "jp",
		ScriptType: "tagged-events",
	})
	require.NoError(t, err)

	settings, err := client.GetSettings(ctx)
	require.NoError(t, err)
	assert.Equal("jp", settings.Language)
	assert.Equal("tagged-events", settings.ScriptType)
}
