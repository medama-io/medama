package services

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/iputils"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
)

func (h *Handler) GetTenantSettings(
	ctx context.Context,
	_params api.GetTenantSettingsParams,
) (api.GetTenantSettingsRes, error) {
	settings, err := h.db.GetTenantSettings(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tenant settings")
	}

	response, err := buildTenantSettingsResponse(settings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build tenant settings response")
	}

	return response, nil
}

func (h *Handler) PatchTenantSettings(
	ctx context.Context,
	req *api.TenantSettings,
	_params api.PatchTenantSettingsParams,
) (api.PatchTenantSettingsRes, error) {
	log := logger.Get()
	if h.auth.IsDemoMode {
		log.Debug().Msg("patch user rejected in demo mode")
		return ErrForbidden(model.ErrDemoMode), nil
	}

	// Convert tenant settings from request to model format
	modifiedSettings := &db.UpdateTenantSettings{}

	if req.ScriptType != nil {
		features := make([]string, 0, len(req.ScriptType))
		for _, v := range req.ScriptType {
			features = append(features, string(v))
		}

		scriptType := strings.Join(features, ",")
		modifiedSettings.ScriptType = &scriptType
	}

	if v, ok := req.BlockAbusiveIPs.Get(); ok {
		blockAbusiveIPs := strconv.FormatBool(v)
		modifiedSettings.BlockAbusiveIPs = &blockAbusiveIPs
	}

	if v, ok := req.BlockTorExitNodes.Get(); ok {
		blockTorExitNodes := strconv.FormatBool(v)
		modifiedSettings.BlockTorExitNodes = &blockTorExitNodes
	}

	if req.BlockedIPs != nil {
		blockedIPs := iputils.GetAddrListString(req.BlockedIPs)
		modifiedSettings.BlockedIPs = &blockedIPs
	}

	// Update tenant settings in database
	err := h.db.UpdateTenantSettings(ctx, modifiedSettings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update tenant settings")
	}

	// Retrieve fresh settings
	settings, err := h.db.GetTenantSettings(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update tenant settings")
	}

	// Update runtime settings
	err = h.RuntimeConfig.UpdateConfig(settings)
	if err != nil {
		return nil, errors.Wrap(err, "update system settings: update runtime config")
	}

	// Build tenant settings response
	response, err := buildTenantSettingsResponse(settings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update tenant settings")
	}

	return response, nil
}

func buildTenantSettingsResponse(
	settings *model.TenantSettings,
) (*api.TenantSettingsHeaders, error) {
	scriptFeatures := []api.TenantSettingsScriptTypeItem{}

	if settings.ScriptType != "" {
		for v := range strings.SplitSeq(settings.ScriptType, ",") {
			scriptFeatures = append(scriptFeatures, api.TenantSettingsScriptTypeItem(v))
		}
	}

	blockedIPs, err := iputils.GetAddrList(settings.BlockedIPs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse blocked IPs")
	}

	blockAbusiveIPs, err := strconv.ParseBool(settings.BlockAbusiveIPs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse block abusive IPs setting")
	}

	blockTorExitNodes, err := strconv.ParseBool(settings.BlockTorExitNodes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse block Tor exit nodes setting")
	}

	return &api.TenantSettingsHeaders{
		Response: api.TenantSettings{
			ScriptType:        scriptFeatures,
			BlockAbusiveIPs:   api.NewOptBool(blockAbusiveIPs),
			BlockTorExitNodes: api.NewOptBool(blockTorExitNodes),
			BlockedIPs:        blockedIPs,
		},
	}, nil
}
