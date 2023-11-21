package services

import (
	"context"
	"log/slog"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

func (h *Handler) GetWebsiteIDSummary(ctx context.Context, params api.GetWebsiteIDSummaryParams) (api.GetWebsiteIDSummaryRes, error) {
	attributes := []slog.Attr{
		slog.String("hostname", params.Hostname),
	}

	// Check if website exists
	exists, err := h.db.WebsiteExists(ctx, params.Hostname)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to check if website exists", attributes...)
		return ErrInternalServerError(err), nil
	} else if !exists {
		slog.LogAttrs(ctx, slog.LevelDebug, "website not found", attributes...)
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Get summary
	summary, err := h.analyticsDB.GetWebsiteSummary(ctx, params.Hostname)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to get website summary", attributes...)
		return ErrInternalServerError(err), nil
	}

	return &api.StatsSummary{
		Uniques:   summary.Uniques,
		Pageviews: summary.Pageviews,
		Bounces:   summary.Bounces,
		Duration:  summary.Duration,
	}, nil
}

func (h *Handler) GetWebsitesIDActive(ctx context.Context, params api.GetWebsitesIDActiveParams) (api.GetWebsitesIDActiveRes, error) {
	return nil, nil
}
