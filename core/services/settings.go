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

func (h *Handler) GetSystemSettings(
	ctx context.Context,
	_params api.GetSystemSettingsParams,
) (api.GetSystemSettingsRes, error) {
	settings, err := h.db.GetSystemSettings(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get system settings")
	}

	response, err := buildSystemSettingsResponse(settings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build system settings response")
	}

	return response, nil
}

func (h *Handler) PatchSystemSettings(
	ctx context.Context,
	req *api.SystemSettings,
	_params api.PatchSystemSettingsParams,
) (api.PatchSystemSettingsRes, error) {
	log := logger.Get()
	if h.auth.IsDemoMode {
		log.Debug().Msg("patch user rejected in demo mode")
		return ErrForbidden(model.ErrDemoMode), nil
	}

	// Convert system settings from request to model format
	modifiedSettings := &db.UpdateSystemSettings{}

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

	// Update system settings in database
	err := h.db.UpdateSystemSettings(ctx, modifiedSettings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update system settings")
	}

	// Retrieve fresh settings
	settings, err := h.db.GetSystemSettings(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update system settings")
	}

	// Update runtime settings
	err = h.RuntimeConfig.UpdateConfig(settings)
	if err != nil {
		return nil, errors.Wrap(err, "update system settings: update runtime config")
	}

	// Build system settings response
	response, err := buildSystemSettingsResponse(settings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update system settings")
	}

	return response, nil
}

func buildSystemSettingsResponse(
	settings *model.SystemSettings,
) (*api.SystemSettingsHeaders, error) {
	scriptFeatures := []api.SystemSettingsScriptTypeItem{}

	if settings.ScriptType != "" {
		for v := range strings.SplitSeq(settings.ScriptType, ",") {
			scriptFeatures = append(scriptFeatures, api.SystemSettingsScriptTypeItem(v))
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

	return &api.SystemSettingsHeaders{
		Response: api.SystemSettings{
			ScriptType:        scriptFeatures,
			BlockAbusiveIPs:   api.NewOptBool(blockAbusiveIPs),
			BlockTorExitNodes: api.NewOptBool(blockTorExitNodes),
			BlockedIPs:        blockedIPs,
		},
	}, nil
}
