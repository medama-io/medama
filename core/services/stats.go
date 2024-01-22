package services

import (
	"context"
	"log/slog"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/db"
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

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:         db.NewFilter(db.FilterPathname, params.Path),
		ReferrerHostname: db.NewFilter(db.FilterReferrerHostname, params.ReferrerHost),
		ReferrerPathname: db.NewFilter(db.FilterReferrerPathname, params.ReferrerPath),
		UTMSource:        db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:        db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign:      db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:          db.NewFilter(db.FilterBrowser, params.Browser),
		BrowserVersion:   db.NewFilter(db.FilterBrowserVersion, params.BrowserVersion),
		OS:               db.NewFilter(db.FilterOS, params.Os),
		Device:           db.NewFilter(db.FilterDevice, params.Device),
		Country:          db.NewFilter(db.FilterCountry, params.Country),
		Language:         db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Get summary
	currentSummary, err := h.analyticsDB.GetWebsiteSummary(ctx, filters)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to get website summary", attributes...)
		return ErrInternalServerError(err), nil
	}
	current := api.StatsSummaryCurrent{
		Uniques:   currentSummary.Uniques,
		Pageviews: currentSummary.Pageviews,
		Bounces:   currentSummary.Bounces,
		Duration:  currentSummary.Duration,
		Active:    currentSummary.Active,
	}

	// Include previous summary if requested
	if params.Previous.Value && params.Start.IsSet() && params.End.IsSet() {
		// Update filter periods to get previous summary
		// Calculate the difference between the start and end dates
		// and subtract that from the start date to get the previous period.
		difference := params.End.Value.Sub(params.Start.Value)
		filters.PeriodStart = params.Start.Value.Add(-difference).Format(model.DateFormat)
		filters.PeriodEnd = params.Start.Value.Format(model.DateFormat)

		previousSummary, err := h.analyticsDB.GetWebsiteSummary(ctx, filters)
		if err != nil {
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get previous website summary", attributes...)
			return ErrInternalServerError(err), nil
		}

		return &api.StatsSummary{
			Current: current,
			Previous: api.NewOptStatsSummaryPrevious(
				api.StatsSummaryPrevious{
					Uniques:   previousSummary.Uniques,
					Pageviews: previousSummary.Pageviews,
					Bounces:   previousSummary.Bounces,
					Duration:  previousSummary.Duration,
				},
			),
		}, nil
	}

	return &api.StatsSummary{
		Current: current,
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

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:         db.NewFilter(db.FilterPathname, params.Path),
		ReferrerHostname: db.NewFilter(db.FilterReferrerHostname, params.ReferrerHost),
		ReferrerPathname: db.NewFilter(db.FilterReferrerPathname, params.ReferrerPath),
		UTMSource:        db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:        db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign:      db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:          db.NewFilter(db.FilterBrowser, params.Browser),
		BrowserVersion:   db.NewFilter(db.FilterBrowserVersion, params.BrowserVersion),
		OS:               db.NewFilter(db.FilterOS, params.Os),
		Device:           db.NewFilter(db.FilterDevice, params.Device),
		Country:          db.NewFilter(db.FilterCountry, params.Country),
		Language:         db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Check parameter if it is asking for summary
	switch params.Summary.Value {
	case true:
		// Get summary
		pages, err := h.analyticsDB.GetWebsitePagesSummary(ctx, filters)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website pages summary", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsPages{}
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
		pages, err := h.analyticsDB.GetWebsitePages(ctx, filters)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website pages", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsPages{}
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

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:         db.NewFilter(db.FilterPathname, params.Path),
		ReferrerHostname: db.NewFilter(db.FilterReferrerHostname, params.ReferrerHost),
		ReferrerPathname: db.NewFilter(db.FilterReferrerPathname, params.ReferrerPath),
		UTMSource:        db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:        db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign:      db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:          db.NewFilter(db.FilterBrowser, params.Browser),
		BrowserVersion:   db.NewFilter(db.FilterBrowserVersion, params.BrowserVersion),
		OS:               db.NewFilter(db.FilterOS, params.Os),
		Device:           db.NewFilter(db.FilterDevice, params.Device),
		Country:          db.NewFilter(db.FilterCountry, params.Country),
		Language:         db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Check parameter if it is asking for summary
	switch params.Summary.Value {
	case true:
		// Get summary
		times, err := h.analyticsDB.GetWebsiteTimeSummary(ctx, filters)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website time summary", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsTime{}
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
		times, err := h.analyticsDB.GetWebsiteTime(ctx, filters)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website time", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsTime{}
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

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:         db.NewFilter(db.FilterPathname, params.Path),
		ReferrerHostname: db.NewFilter(db.FilterReferrerHostname, params.ReferrerHost),
		ReferrerPathname: db.NewFilter(db.FilterReferrerPathname, params.ReferrerPath),
		UTMSource:        db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:        db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign:      db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:          db.NewFilter(db.FilterBrowser, params.Browser),
		BrowserVersion:   db.NewFilter(db.FilterBrowserVersion, params.BrowserVersion),
		OS:               db.NewFilter(db.FilterOS, params.Os),
		Device:           db.NewFilter(db.FilterDevice, params.Device),
		Country:          db.NewFilter(db.FilterCountry, params.Country),
		Language:         db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Check parameter if it is asking for summary
	switch params.Summary.Value {
	case true:
		// Get summary
		referrers, err := h.analyticsDB.GetWebsiteReferrersSummary(ctx, filters)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website referrers summary", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsReferrers{}
		for _, page := range referrers {
			res = append(res, api.StatsReferrersItem{
				ReferrerHost:     page.ReferrerHostname,
				Uniques:          page.Uniques,
				UniquePercentage: page.UniquePercentage,
			})
		}

		return &res, nil
	case false:
		// Get referrers
		referrers, err := h.analyticsDB.GetWebsiteReferrers(ctx, filters)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website referrers", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsReferrers{}
		for _, page := range referrers {
			res = append(res, api.StatsReferrersItem{
				ReferrerHost:     page.ReferrerHostname,
				ReferrerPath:     api.NewOptString(page.ReferrerPathname),
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

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:         db.NewFilter(db.FilterPathname, params.Path),
		ReferrerHostname: db.NewFilter(db.FilterReferrerHostname, params.ReferrerHost),
		ReferrerPathname: db.NewFilter(db.FilterReferrerPathname, params.ReferrerPath),
		UTMSource:        db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:        db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign:      db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:          db.NewFilter(db.FilterBrowser, params.Browser),
		BrowserVersion:   db.NewFilter(db.FilterBrowserVersion, params.BrowserVersion),
		OS:               db.NewFilter(db.FilterOS, params.Os),
		Device:           db.NewFilter(db.FilterDevice, params.Device),
		Country:          db.NewFilter(db.FilterCountry, params.Country),
		Language:         db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Get sources
	sources, err := h.analyticsDB.GetWebsiteUTMSources(ctx, filters)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to get website utm sources", attributes...)
		return ErrInternalServerError(err), nil
	}

	// Create API response
	res := api.StatsUTMSources{}
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

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:         db.NewFilter(db.FilterPathname, params.Path),
		ReferrerHostname: db.NewFilter(db.FilterReferrerHostname, params.ReferrerHost),
		ReferrerPathname: db.NewFilter(db.FilterReferrerPathname, params.ReferrerPath),
		UTMSource:        db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:        db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign:      db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:          db.NewFilter(db.FilterBrowser, params.Browser),
		BrowserVersion:   db.NewFilter(db.FilterBrowserVersion, params.BrowserVersion),
		OS:               db.NewFilter(db.FilterOS, params.Os),
		Device:           db.NewFilter(db.FilterDevice, params.Device),
		Country:          db.NewFilter(db.FilterCountry, params.Country),
		Language:         db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Get mediums
	mediums, err := h.analyticsDB.GetWebsiteUTMMediums(ctx, filters)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to get website utm mediums", attributes...)
		return ErrInternalServerError(err), nil
	}

	// Create API response
	res := api.StatsUTMMediums{}
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

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:         db.NewFilter(db.FilterPathname, params.Path),
		ReferrerHostname: db.NewFilter(db.FilterReferrerHostname, params.ReferrerHost),
		ReferrerPathname: db.NewFilter(db.FilterReferrerPathname, params.ReferrerPath),
		UTMSource:        db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:        db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign:      db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:          db.NewFilter(db.FilterBrowser, params.Browser),
		BrowserVersion:   db.NewFilter(db.FilterBrowserVersion, params.BrowserVersion),
		OS:               db.NewFilter(db.FilterOS, params.Os),
		Device:           db.NewFilter(db.FilterDevice, params.Device),
		Country:          db.NewFilter(db.FilterCountry, params.Country),
		Language:         db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Get campaigns
	campaigns, err := h.analyticsDB.GetWebsiteUTMCampaigns(ctx, filters)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to get website utm campaigns", attributes...)
		return ErrInternalServerError(err), nil
	}

	// Create API response
	res := api.StatsUTMCampaigns{}
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

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:         db.NewFilter(db.FilterPathname, params.Path),
		ReferrerHostname: db.NewFilter(db.FilterReferrerHostname, params.ReferrerHost),
		ReferrerPathname: db.NewFilter(db.FilterReferrerPathname, params.ReferrerPath),
		UTMSource:        db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:        db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign:      db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:          db.NewFilter(db.FilterBrowser, params.Browser),
		BrowserVersion:   db.NewFilter(db.FilterBrowserVersion, params.BrowserVersion),
		OS:               db.NewFilter(db.FilterOS, params.Os),
		Device:           db.NewFilter(db.FilterDevice, params.Device),
		Country:          db.NewFilter(db.FilterCountry, params.Country),
		Language:         db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Check parameter if it is asking for summary
	switch params.Summary.Value {
	case true:
		// Get summary
		browsers, err := h.analyticsDB.GetWebsiteBrowsersSummary(ctx, filters)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website browsers summary", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsBrowsers{}
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
		browsers, err := h.analyticsDB.GetWebsiteBrowsers(ctx, filters)
		if err != nil {
			attributes = append(attributes, slog.Bool("summary", params.Summary.Value), slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get website browsers", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsBrowsers{}
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

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:         db.NewFilter(db.FilterPathname, params.Path),
		ReferrerHostname: db.NewFilter(db.FilterReferrerHostname, params.ReferrerHost),
		ReferrerPathname: db.NewFilter(db.FilterReferrerPathname, params.ReferrerPath),
		UTMSource:        db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:        db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign:      db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:          db.NewFilter(db.FilterBrowser, params.Browser),
		BrowserVersion:   db.NewFilter(db.FilterBrowserVersion, params.BrowserVersion),
		OS:               db.NewFilter(db.FilterOS, params.Os),
		Device:           db.NewFilter(db.FilterDevice, params.Device),
		Country:          db.NewFilter(db.FilterCountry, params.Country),
		Language:         db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Get OS
	os, err := h.analyticsDB.GetWebsiteOS(ctx, filters)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to get website os", attributes...)
		return ErrInternalServerError(err), nil
	}

	// Create API response
	res := api.StatsOS{}
	for _, page := range os {
		res = append(res, api.StatsOSItem{
			Os:               page.OS.String(),
			Uniques:          page.Uniques,
			UniquePercentage: page.UniquePercentage,
		})
	}

	return &res, nil
}

func (h *Handler) GetWebsiteIDDevice(ctx context.Context, params api.GetWebsiteIDDeviceParams) (api.GetWebsiteIDDeviceRes, error) {
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

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:         db.NewFilter(db.FilterPathname, params.Path),
		ReferrerHostname: db.NewFilter(db.FilterReferrerHostname, params.ReferrerHost),
		ReferrerPathname: db.NewFilter(db.FilterReferrerPathname, params.ReferrerPath),
		UTMSource:        db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:        db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign:      db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:          db.NewFilter(db.FilterBrowser, params.Browser),
		BrowserVersion:   db.NewFilter(db.FilterBrowserVersion, params.BrowserVersion),
		OS:               db.NewFilter(db.FilterOS, params.Os),
		Device:           db.NewFilter(db.FilterDevice, params.Device),
		Country:          db.NewFilter(db.FilterCountry, params.Country),
		Language:         db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Get devices
	devices, err := h.analyticsDB.GetWebsiteDevices(ctx, filters)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to get website devices", attributes...)
		return ErrInternalServerError(err), nil
	}

	// Create API response
	res := api.StatsDevices{}
	for _, page := range devices {
		res = append(res, api.StatsDevicesItem{
			Device:           page.Device.String(),
			Uniques:          page.Uniques,
			UniquePercentage: page.UniquePercentage,
		})
	}

	return &res, nil
}

func (h *Handler) GetWebsiteIDScreen(ctx context.Context, params api.GetWebsiteIDScreenParams) (api.GetWebsiteIDScreenRes, error) {
	return nil, nil
}

func (h *Handler) GetWebsiteIDLanguage(ctx context.Context, params api.GetWebsiteIDLanguageParams) (api.GetWebsiteIDLanguageRes, error) {
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

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:         db.NewFilter(db.FilterPathname, params.Path),
		ReferrerHostname: db.NewFilter(db.FilterReferrerHostname, params.ReferrerHost),
		ReferrerPathname: db.NewFilter(db.FilterReferrerPathname, params.ReferrerPath),
		UTMSource:        db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:        db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign:      db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:          db.NewFilter(db.FilterBrowser, params.Browser),
		BrowserVersion:   db.NewFilter(db.FilterBrowserVersion, params.BrowserVersion),
		OS:               db.NewFilter(db.FilterOS, params.Os),
		Device:           db.NewFilter(db.FilterDevice, params.Device),
		Country:          db.NewFilter(db.FilterCountry, params.Country),
		Language:         db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Get languages
	languages, err := h.analyticsDB.GetWebsiteLanguages(ctx, filters)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to get website languages", attributes...)
		return ErrInternalServerError(err), nil
	}

	// Create API response
	res := api.StatsLanguages{}
	for _, page := range languages {
		res = append(res, api.StatsLanguagesItem{
			Language:         page.Language,
			Uniques:          page.Uniques,
			UniquePercentage: page.UniquePercentage,
		})
	}

	return &res, nil
}

func (h *Handler) GetWebsiteIDCountry(ctx context.Context, params api.GetWebsiteIDCountryParams) (api.GetWebsiteIDCountryRes, error) {
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

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:         db.NewFilter(db.FilterPathname, params.Path),
		ReferrerHostname: db.NewFilter(db.FilterReferrerHostname, params.ReferrerHost),
		ReferrerPathname: db.NewFilter(db.FilterReferrerPathname, params.ReferrerPath),
		UTMSource:        db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:        db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign:      db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:          db.NewFilter(db.FilterBrowser, params.Browser),
		BrowserVersion:   db.NewFilter(db.FilterBrowserVersion, params.BrowserVersion),
		OS:               db.NewFilter(db.FilterOS, params.Os),
		Device:           db.NewFilter(db.FilterDevice, params.Device),
		Country:          db.NewFilter(db.FilterCountry, params.Country),
		Language:         db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Get countries
	countries, err := h.analyticsDB.GetWebsiteCountries(ctx, filters)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to get website countries", attributes...)
		return ErrInternalServerError(err), nil
	}

	// Create API response
	res := api.StatsCountries{}
	for _, page := range countries {
		// Convert country code to country name
		country, err := h.codeCountryMap.GetCountry(page.Country)
		if err != nil {
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get country name", attributes...)
			return ErrInternalServerError(err), nil
		}

		res = append(res, api.StatsCountriesItem{
			Country:          country,
			Uniques:          page.Uniques,
			UniquePercentage: page.UniquePercentage,
		})
	}

	return &res, nil
}
