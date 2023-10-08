package services

import (
	"context"
	"errors"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

func (h *Handler) GetWebsiteIDSummary(ctx context.Context, params api.GetWebsiteIDSummaryParams) (api.GetWebsiteIDSummaryRes, error) {
	_, err := h.db.GetWebsite(ctx, params.ID)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteNotFound) {
			return ErrNotFound(err), nil
		}

		return nil, err
	}

	return &api.StatsSummary{
		Uniques:   api.NewOptInt(0),
		Pageviews: api.NewOptInt(0),
		Bounces:   api.NewOptFloat32(0),
		Duration:  api.NewOptInt(0),
	}, nil
}

func (h *Handler) GetWebsitesIDActive(ctx context.Context, params api.GetWebsitesIDActiveParams) (api.GetWebsitesIDActiveRes, error) {
	return nil, nil
}
