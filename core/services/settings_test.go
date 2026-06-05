package services_test

import (
	"net/http"
	"net/netip"
	"testing"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
	"github.com/stretchr/testify/require"
)

func TestGetSystemSettings(t *testing.T) {
	assert, ctx, handler, _ := NewTestHandler(t)

	resp, err := handler.GetSystemSettings(ctx, api.GetSystemSettingsParams{})
	require.NoError(t, err)

	settings, ok := resp.(*api.SystemSettingsHeaders)
	require.True(t, ok)

	assert.Equal(
		[]api.SystemSettingsScriptTypeItem{api.SystemSettingsScriptTypeItemDefault},
		settings.Response.ScriptType,
	)
	assert.Equal(api.NewOptBool(true), settings.Response.BlockAbusiveIPs)
	assert.Equal(api.NewOptBool(true), settings.Response.BlockTorExitNodes)
	assert.Empty(settings.Response.BlockedIPs)
}

func TestPatchSystemSettings(t *testing.T) {
	assert, ctx, handler, _ := NewTestHandler(t)

	req := &api.SystemSettings{
		ScriptType: []api.SystemSettingsScriptTypeItem{
			api.SystemSettingsScriptTypeItemClickEvents,
			api.SystemSettingsScriptTypeItemPageEvents,
		},
		BlockAbusiveIPs:   api.NewOptBool(false),
		BlockTorExitNodes: api.NewOptBool(false),
		BlockedIPs:        []netip.Addr{netip.MustParseAddr("10.0.0.1")},
	}

	resp, err := handler.PatchSystemSettings(ctx, req, api.PatchSystemSettingsParams{})
	require.NoError(t, err)

	settings, ok := resp.(*api.SystemSettingsHeaders)
	require.True(t, ok)

	assert.Equal(
		[]api.SystemSettingsScriptTypeItem{
			api.SystemSettingsScriptTypeItemClickEvents,
			api.SystemSettingsScriptTypeItemPageEvents,
		},
		settings.Response.ScriptType,
	)
	assert.Equal(api.NewOptBool(false), settings.Response.BlockAbusiveIPs)
	assert.Equal(api.NewOptBool(false), settings.Response.BlockTorExitNodes)
	assert.Equal([]netip.Addr{netip.MustParseAddr("10.0.0.1")}, settings.Response.BlockedIPs)
}

func TestPatchSystemSettingsPartial(t *testing.T) {
	assert, ctx, handler, sqliteClient := NewTestHandler(t)

	err := sqliteClient.UpdateSystemSettings(ctx, &db.UpdateSystemSettings{
		ScriptType:        ptr("click-events"),
		BlockAbusiveIPs:   ptr("false"),
		BlockTorExitNodes: ptr("false"),
		BlockedIPs:        ptr("10.0.0.1"),
	})
	require.NoError(t, err)

	req := &api.SystemSettings{
		ScriptType: []api.SystemSettingsScriptTypeItem{
			api.SystemSettingsScriptTypeItemPageEvents,
		},
	}

	resp, err := handler.PatchSystemSettings(ctx, req, api.PatchSystemSettingsParams{})
	require.NoError(t, err)

	settings, ok := resp.(*api.SystemSettingsHeaders)
	require.True(t, ok)

	assert.Equal(
		[]api.SystemSettingsScriptTypeItem{api.SystemSettingsScriptTypeItemPageEvents},
		settings.Response.ScriptType,
	)
	assert.Equal(api.NewOptBool(false), settings.Response.BlockAbusiveIPs)
	assert.Equal(api.NewOptBool(false), settings.Response.BlockTorExitNodes)
	assert.Equal([]netip.Addr{netip.MustParseAddr("10.0.0.1")}, settings.Response.BlockedIPs)
}

func TestPatchSystemSettingsDemoMode(t *testing.T) {
	assert, ctx, handler, _ := NewTestHandlerDemoMode(t)

	req := &api.SystemSettings{
		ScriptType: []api.SystemSettingsScriptTypeItem{
			api.SystemSettingsScriptTypeItemClickEvents,
		},
	}

	resp, err := handler.PatchSystemSettings(ctx, req, api.PatchSystemSettingsParams{})
	require.NoError(t, err)

	forbidden, ok := resp.(*api.ForbiddenErrorHeaders)
	require.True(t, ok)

	assert.Equal(int32(http.StatusForbidden), forbidden.Response.Error.Code)
	assert.Equal(model.ErrDemoMode.Error(), forbidden.Response.Error.Message)
}

func ptr(s string) *string {
	return &s
}
