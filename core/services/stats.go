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
		Active:    api.NewOptInt(summary.Active),
	}, nil
}

func (h *Handler) GetWebsiteIDPages(ctx context.Context, params api.GetWebsiteIDPagesParams) (api.GetWebsiteIDPagesRes, error) {
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

	// Check parameter if it is asking for summary
	switch params.Summary.Value {
	case true:
		// Get summary
		pages, err := h.analyticsDB.GetWebsitePagesSummary(ctx, params.Hostname)
		if err != nil {
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website pages summary", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		var res api.StatsPages
		for _, page := range pages {
			res = append(res, api.StatsPagesItem{
				Path:             api.NewOptString(page.Pathname),
				Uniques:          page.Uniques,
				Uniquepercentage: page.UniquesPercentage,
			})
		}

		return &res, nil
	case false:
		// Get pages
		pages, err := h.analyticsDB.GetWebsitePages(ctx, params.Hostname)
		if err != nil {
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website pages", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		var res api.StatsPages
		for _, page := range pages {
			res = append(res, api.StatsPagesItem{
				Path:             api.NewOptString(page.Pathname),
				Uniques:          page.Uniques,
				Uniquepercentage: page.UniquesPercentage,
				Title:            api.NewOptString(page.Title),
				Pageviews:        api.NewOptInt(page.Pageviews),
				Bounces:          api.NewOptInt(page.Bounces),
				Duration:         api.NewOptInt(page.Duration),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}
