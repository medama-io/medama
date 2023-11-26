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
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website pages summary", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		var res api.StatsPages
		for _, page := range pages {
			res = append(res, api.StatsPagesItem{
				Path:             page.Pathname,
				Uniques:          page.Uniques,
				UniquePercentage: page.UniquePercentage,
			})
		}

		return &res, nil
	case false:
		// Get pages
		pages, err := h.analyticsDB.GetWebsitePages(ctx, params.Hostname)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website pages", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		var res api.StatsPages
		for _, page := range pages {
			res = append(res, api.StatsPagesItem{
				Path:             page.Pathname,
				Uniques:          page.Uniques,
				UniquePercentage: page.UniquePercentage,
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

func (h *Handler) GetWebsiteIDTime(ctx context.Context, params api.GetWebsiteIDTimeParams) (api.GetWebsiteIDTimeRes, error) {
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
		times, err := h.analyticsDB.GetWebsiteTimeSummary(ctx, params.Hostname)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website time summary", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		var res api.StatsTime
		for _, page := range times {
			res = append(res, api.StatsTimeItem{
				Path:               page.Pathname,
				Duration:           page.Duration,
				DurationPercentage: page.DurationPercentage,
			})
		}

		return &res, nil
	case false:
		// Get time
		times, err := h.analyticsDB.GetWebsiteTime(ctx, params.Hostname)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website time", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		var res api.StatsTime
		for _, page := range times {
			res = append(res, api.StatsTimeItem{
				Path:                  page.Pathname,
				Duration:              page.Duration,
				DurationPercentage:    page.DurationPercentage,
				DurationUpperQuartile: api.NewOptInt(page.DurationUpperQuartile),
				DurationLowerQuartile: api.NewOptInt(page.DurationLowerQuartile),
				Title:                 api.NewOptString(page.Title),
				Bounces:               api.NewOptInt(page.Bounces),
				Uniques:               api.NewOptInt(page.Uniques),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDReferrers(ctx context.Context, params api.GetWebsiteIDReferrersParams) (api.GetWebsiteIDReferrersRes, error) {
	attributes := []slog.Attr{
		slog.String("hostname", params.Hostname),
	}

	// Check if website exists
	exists, err := h.db.WebsiteExists(ctx, params.Hostname)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to check if website exists", attributes...)
		return ErrInternalServerError(err), nil
	}
	if !exists {
		slog.LogAttrs(ctx, slog.LevelDebug, "website not found", attributes...)
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Check parameter if it is asking for summary
	switch params.Summary.Value {
	case true:
		// Get summary
		referrers, err := h.analyticsDB.GetWebsiteReferrersSummary(ctx, params.Hostname)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website referrers summary", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		var res api.StatsReferrers
		for _, page := range referrers {
			res = append(res, api.StatsReferrersItem{
				Referrer:         page.Referrer,
				Uniques:          page.Uniques,
				UniquePercentage: page.UniquePercentage,
			})
		}

		return &res, nil
	case false:
		// Get referrers
		referrers, err := h.analyticsDB.GetWebsiteReferrers(ctx, params.Hostname)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website referrers", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		var res api.StatsReferrers
		for _, page := range referrers {
			res = append(res, api.StatsReferrersItem{
				Referrer:         page.Referrer,
				Uniques:          page.Uniques,
				UniquePercentage: page.UniquePercentage,
				Bounces:          api.NewOptInt(page.Bounces),
				Duration:         api.NewOptInt(page.Duration),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDSources(ctx context.Context, params api.GetWebsiteIDSourcesParams) (api.GetWebsiteIDSourcesRes, error) {
	attributes := []slog.Attr{
		slog.String("hostname", params.Hostname),
	}

	// Check if website exists
	exists, err := h.db.WebsiteExists(ctx, params.Hostname)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to check if website exists", attributes...)
		return ErrInternalServerError(err), nil
	}
	if !exists {
		slog.LogAttrs(ctx, slog.LevelDebug, "website not found", attributes...)
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Get sources
	sources, err := h.analyticsDB.GetWebsiteUTMSources(ctx, params.Hostname)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to get website utm sources", attributes...)
		return ErrInternalServerError(err), nil
	}

	// Create API response
	var res api.StatsUTMSources
	for _, page := range sources {
		res = append(res, api.StatsUTMSourcesItem{
			Source:           page.Source,
			Uniques:          page.Uniques,
			UniquePercentage: page.UniquePercentage,
		})
	}

	return &res, nil
}

func (h *Handler) GetWebsiteIDMediums(ctx context.Context, params api.GetWebsiteIDMediumsParams) (api.GetWebsiteIDMediumsRes, error) {
	attributes := []slog.Attr{
		slog.String("hostname", params.Hostname),
	}

	// Check if website exists
	exists, err := h.db.WebsiteExists(ctx, params.Hostname)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to check if website exists", attributes...)
		return ErrInternalServerError(err), nil
	}
	if !exists {
		slog.LogAttrs(ctx, slog.LevelDebug, "website not found", attributes...)
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Get mediums
	mediums, err := h.analyticsDB.GetWebsiteUTMMediums(ctx, params.Hostname)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to get website utm mediums", attributes...)
		return ErrInternalServerError(err), nil
	}

	// Create API response
	var res api.StatsUTMMediums
	for _, page := range mediums {
		res = append(res, api.StatsUTMMediumsItem{
			Medium:           page.Medium,
			Uniques:          page.Uniques,
			UniquePercentage: page.UniquePercentage,
		})
	}

	return &res, nil
}

func (h *Handler) GetWebsiteIDCampaigns(ctx context.Context, params api.GetWebsiteIDCampaignsParams) (api.GetWebsiteIDCampaignsRes, error) {
	attributes := []slog.Attr{
		slog.String("hostname", params.Hostname),
	}

	// Check if website exists
	exists, err := h.db.WebsiteExists(ctx, params.Hostname)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to check if website exists", attributes...)
		return ErrInternalServerError(err), nil
	}
	if !exists {
		slog.LogAttrs(ctx, slog.LevelDebug, "website not found", attributes...)
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Get campaigns
	campaigns, err := h.analyticsDB.GetWebsiteUTMCampaigns(ctx, params.Hostname)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to get website utm campaigns", attributes...)
		return ErrInternalServerError(err), nil
	}

	// Create API response
	var res api.StatsUTMCampaigns
	for _, page := range campaigns {
		res = append(res, api.StatsUTMCampaignsItem{
			Campaign:         page.Campaign,
			Uniques:          page.Uniques,
			UniquePercentage: page.UniquePercentage,
		})
	}

	return &res, nil
}

func (h *Handler) GetWebsiteIDBrowsers(ctx context.Context, params api.GetWebsiteIDBrowsersParams) (api.GetWebsiteIDBrowsersRes, error) {
	attributes := []slog.Attr{
		slog.String("hostname", params.Hostname),
	}

	// Check if website exists
	exists, err := h.db.WebsiteExists(ctx, params.Hostname)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to check if website exists", attributes...)
		return ErrInternalServerError(err), nil
	}
	if !exists {
		slog.LogAttrs(ctx, slog.LevelDebug, "website not found", attributes...)
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Check parameter if it is asking for summary
	switch params.Summary.Value {
	case true:
		// Get summary
		browsers, err := h.analyticsDB.GetWebsiteBrowsersSummary(ctx, params.Hostname)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website browsers summary", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		var res api.StatsBrowsers
		for _, page := range browsers {
			res = append(res, api.StatsBrowsersItem{
				Browser:          page.Browser.String(),
				Uniques:          page.Uniques,
				UniquePercentage: page.UniquePercentage,
			})
		}

		return &res, nil
	case false:
		// Get browsers
		browsers, err := h.analyticsDB.GetWebsiteBrowsers(ctx, params.Hostname)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website browsers", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		var res api.StatsBrowsers
		for _, page := range browsers {
			res = append(res, api.StatsBrowsersItem{
				Browser:          page.Browser.String(),
				Uniques:          page.Uniques,
				UniquePercentage: page.UniquePercentage,
				Version:          api.NewOptString(page.Version),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDOs(ctx context.Context, params api.GetWebsiteIDOsParams) (api.GetWebsiteIDOsRes, error) {
	return nil, nil
}

func (h *Handler) GetWebsiteIDScreen(ctx context.Context, params api.GetWebsiteIDScreenParams) (api.GetWebsiteIDScreenRes, error) {
	return nil, nil
}

func (h *Handler) GetWebsiteIDDevice(ctx context.Context, params api.GetWebsiteIDDeviceParams) (api.GetWebsiteIDDeviceRes, error) {
	return nil, nil
}

func (h *Handler) GetWebsiteIDLanguage(ctx context.Context, params api.GetWebsiteIDLanguageParams) (api.GetWebsiteIDLanguageRes, error) {
	return nil, nil
}

func (h *Handler) GetWebsiteIDCountry(ctx context.Context, params api.GetWebsiteIDCountryParams) (api.GetWebsiteIDCountryRes, error) {
	return nil, nil
}
