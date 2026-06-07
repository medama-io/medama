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

func TestGetTenantSettings(t *testing.T) {
	assert, ctx, handler, _ := NewTestHandler(t)

	resp, err := handler.GetTenantSettings(ctx, api.GetTenantSettingsParams{})
	require.NoError(t, err)

	settings, ok := resp.(*api.TenantSettingsHeaders)
	require.True(t, ok)

	assert.Equal(
		[]api.TenantSettingsScriptTypeItem{api.TenantSettingsScriptTypeItemDefault},
		settings.Response.ScriptType,
	)
	assert.Equal(api.NewOptBool(true), settings.Response.BlockAbusiveIPs)
	assert.Equal(api.NewOptBool(true), settings.Response.BlockTorExitNodes)
	assert.Empty(settings.Response.BlockedIPs)
}

func TestPatchTenantSettings(t *testing.T) {
	assert, ctx, handler, _ := NewTestHandler(t)

	req := &api.TenantSettings{
		ScriptType: []api.TenantSettingsScriptTypeItem{
			api.TenantSettingsScriptTypeItemClickEvents,
			api.TenantSettingsScriptTypeItemPageEvents,
		},
		BlockAbusiveIPs:   api.NewOptBool(false),
		BlockTorExitNodes: api.NewOptBool(false),
		BlockedIPs:        []netip.Addr{netip.MustParseAddr("10.0.0.1")},
	}

	resp, err := handler.PatchTenantSettings(ctx, req, api.PatchTenantSettingsParams{})
	require.NoError(t, err)

	settings, ok := resp.(*api.TenantSettingsHeaders)
	require.True(t, ok)

	assert.Equal(
		[]api.TenantSettingsScriptTypeItem{
			api.TenantSettingsScriptTypeItemClickEvents,
			api.TenantSettingsScriptTypeItemPageEvents,
		},
		settings.Response.ScriptType,
	)
	assert.Equal(api.NewOptBool(false), settings.Response.BlockAbusiveIPs)
	assert.Equal(api.NewOptBool(false), settings.Response.BlockTorExitNodes)
	assert.Equal([]netip.Addr{netip.MustParseAddr("10.0.0.1")}, settings.Response.BlockedIPs)
}

func TestPatchTenantSettingsPartial(t *testing.T) {
	assert, ctx, handler, sqliteClient := NewTestHandler(t)

	err := sqliteClient.UpdateTenantSettings(ctx, &db.UpdateTenantSettings{
		ScriptType:        ptr("click-events"),
		BlockAbusiveIPs:   ptr("false"),
		BlockTorExitNodes: ptr("false"),
		BlockedIPs:        ptr("10.0.0.1"),
	})
	require.NoError(t, err)

	req := &api.TenantSettings{
		ScriptType: []api.TenantSettingsScriptTypeItem{
			api.TenantSettingsScriptTypeItemPageEvents,
		},
	}

	resp, err := handler.PatchTenantSettings(ctx, req, api.PatchTenantSettingsParams{})
	require.NoError(t, err)

	settings, ok := resp.(*api.TenantSettingsHeaders)
	require.True(t, ok)

	assert.Equal(
		[]api.TenantSettingsScriptTypeItem{api.TenantSettingsScriptTypeItemPageEvents},
		settings.Response.ScriptType,
	)
	assert.Equal(api.NewOptBool(false), settings.Response.BlockAbusiveIPs)
	assert.Equal(api.NewOptBool(false), settings.Response.BlockTorExitNodes)
	assert.Equal([]netip.Addr{netip.MustParseAddr("10.0.0.1")}, settings.Response.BlockedIPs)
}

func TestPatchTenantSettingsDemoMode(t *testing.T) {
	assert, ctx, handler, _ := NewTestHandlerDemoMode(t)

	req := &api.TenantSettings{
		ScriptType: []api.TenantSettingsScriptTypeItem{
			api.TenantSettingsScriptTypeItemClickEvents,
		},
	}

	resp, err := handler.PatchTenantSettings(ctx, req, api.PatchTenantSettingsParams{})
	require.NoError(t, err)

	forbidden, ok := resp.(*api.ForbiddenErrorHeaders)
	require.True(t, ok)

	assert.Equal(int32(http.StatusForbidden), forbidden.Response.Error.Code)
	assert.Equal(model.ErrDemoMode.Error(), forbidden.Response.Error.Message)
}

func ptr(s string) *string {
	return &s
}
