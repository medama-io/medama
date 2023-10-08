package services

import (
	"context"
	"errors"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

func (h *Handler) DeleteWebsitesID(ctx context.Context, params api.DeleteWebsitesIDParams) (api.DeleteWebsitesIDRes, error) {
	err := h.db.DeleteWebsite(ctx, params.ID)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteNotFound) {
			return ErrNotFound(err), nil
		}

		return nil, err
	}

	return nil, nil
}

func (h *Handler) GetWebsites(ctx context.Context, params api.GetWebsitesParams) (api.GetWebsitesRes, error) {
	return nil, nil
}

func (h *Handler) GetWebsitesID(ctx context.Context, params api.GetWebsitesIDParams) (api.GetWebsitesIDRes, error) {
	return nil, nil
}

func (h *Handler) PatchWebsitesID(ctx context.Context, req api.OptWebsitePatch, params api.PatchWebsitesIDParams) (api.PatchWebsitesIDRes, error) {
	return nil, nil
}

func (h *Handler) PostWebsites(ctx context.Context, req api.OptWebsiteCreate) (api.PostWebsitesRes, error) {
	return nil, nil
}
