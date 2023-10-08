package services

import (
	"context"

	"github.com/medama-io/medama/api"
)

func (h *Handler) DeleteWebsitesID(ctx context.Context, params api.DeleteWebsitesIDParams) (api.DeleteWebsitesIDRes, error) {
	err := h.db.DeleteWebsite(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) GetWebsiteIDSummary(ctx context.Context, params api.GetWebsiteIDSummaryParams) (api.GetWebsiteIDSummaryRes, error) {
	_, err := h.db.GetWebsite(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.StatsSummary{
		Uniques:   api.NewOptInt(0),
		Pageviews: api.NewOptInt(0),
		Bounces:   api.NewOptFloat32(0),
		Duration:  api.NewOptInt(0),
	}, nil
}

func (h *Handler) GetWebsites(ctx context.Context) (api.GetWebsitesRes, error) {
	return nil, nil
}

func (h *Handler) GetWebsitesID(ctx context.Context, params api.GetWebsitesIDParams) (api.GetWebsitesIDRes, error) {
	return nil, nil
}

func (h *Handler) GetWebsitesIDActive(ctx context.Context, params api.GetWebsitesIDActiveParams) (api.GetWebsitesIDActiveRes, error) {
	return nil, nil
}

func (h *Handler) PatchWebsitesID(ctx context.Context, req api.OptWebsitePatch, params api.PatchWebsitesIDParams) (api.PatchWebsitesIDRes, error) {
	return nil, nil
}

func (h *Handler) PostWebsites(ctx context.Context, req api.OptWebsiteCreate) (api.PostWebsitesRes, error) {
	return nil, nil
}
