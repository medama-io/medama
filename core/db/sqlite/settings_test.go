package sqlite_test

import (
	"testing"

	"github.com/medama-io/medama/model"
)

func TestGetSetting(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	setting, err := client.GetSetting(ctx, model.SettingsKeyLanguage)
	assert.NoError(err)
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
	assert.NoError(err)
	assert.NotNil(settings)
	assert.Equal("en", settings.Language)
	assert.Equal("default", settings.ScriptType)
	assert.Equal(0, settings.Threads)
	assert.Equal("", settings.MemoryLimit)
}

func TestUpdateSetting(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	err := client.UpdateSetting(ctx, model.SettingsKeyLanguage, "jp")
	assert.NoError(err)

	setting, err := client.GetSetting(ctx, model.SettingsKeyLanguage)
	assert.NoError(err)
	assert.Equal("jp", setting)
}

func TestUpdateSettings(t *testing.T) {
	assert, ctx, client := SetupDatabase(t)

	user, err := client.GetUserByUsername(ctx, "admin")
	assert.NoError(err)

	err = client.UpdateSettings(ctx, user.ID, &model.GlobalSettings{
		DuckDBSettings: model.DuckDBSettings{
			Threads:     4,
			MemoryLimit: "1GB",
		},
		UserSettings: model.UserSettings{
			Language:   "jp",
			ScriptType: "tagged-events",
		},
	})
	assert.NoError(err)

	settings, err := client.GetSettings(ctx)
	assert.NoError(err)
	assert.Equal("jp", settings.Language)
	assert.Equal("tagged-events", settings.ScriptType)
	assert.Equal(4, settings.Threads)
	assert.Equal("1GB", settings.MemoryLimit)
}
