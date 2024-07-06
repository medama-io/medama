package services

import (
	"context"
	"errors"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
)

func (h *Handler) GetWebsiteIDSummary(ctx context.Context, params api.GetWebsiteIDSummaryParams) (api.GetWebsiteIDSummaryRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:    db.NewFilter(db.FilterPathname, params.Path),
		Referrer:    db.NewFilter(db.FilterReferrer, params.Referrer),
		UTMSource:   db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:     db.NewFilter(db.FilterBrowser, params.Browser),
		OS:          db.NewFilter(db.FilterOS, params.Os),
		Device:      db.NewFilter(db.FilterDevice, params.Device),
		Country:     db.NewFilter(db.FilterCountry, params.Country),
		Language:    db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),
	}

	// Get summary
	currentSummary, err := h.analyticsDB.GetWebsiteSummary(ctx, filters)
	if err != nil {
		log.Error().Err(err).Msg("failed to get website summary")
		return ErrInternalServerError(err), nil
	}

	resp := &api.StatsSummary{
		Current: api.StatsSummaryCurrent{
			Visitors:  currentSummary.Visitors,
			Pageviews: currentSummary.Pageviews,
			Bounces:   currentSummary.Bounces,
			Duration:  currentSummary.Duration,
		},
	}

	// Include previous summary if requested.
	if params.Previous.Value && params.Start.IsSet() && params.End.IsSet() {
		// Make a copy of filters to avoid modifying the original.
		filters := *filters
		// Update filter periods to get previous summary.
		// Calculate the difference between the start and end dates and
		// subtract that from the start date to get the previous period.
		difference := params.End.Value.Sub(params.Start.Value)
		filters.PeriodStart = params.Start.Value.Add(-difference).Format(model.DateFormat)
		filters.PeriodEnd = params.Start.Value.Format(model.DateFormat)

		previousSummary, err := h.analyticsDB.GetWebsiteSummary(ctx, &filters)
		if err != nil {
			log.Error().Err(err).Msg("failed to get previous website summary")
			return ErrInternalServerError(err), nil
		}

		resp.Previous = api.NewOptStatsSummaryPrevious(
			api.StatsSummaryPrevious{
				Visitors:  previousSummary.Visitors,
				Pageviews: previousSummary.Pageviews,
				Bounces:   previousSummary.Bounces,
				Duration:  previousSummary.Duration,
			},
		)
	}

	// Return bucketed interval values if requested.
	if params.Interval.Value != "" {
		interval, err := h.analyticsDB.GetWebsiteIntervals(ctx, filters, params.Interval.Value)
		if err != nil {
			log.Error().Err(err).Msg("failed to get website intervals")
			if errors.Is(err, model.ErrInvalidParameter) {
				return ErrBadRequest(err), nil
			}

			return ErrInternalServerError(err), nil
		}

		resp.Interval = []api.StatsSummaryIntervalItem{}
		for _, i := range interval {
			resp.Interval = append(resp.Interval, api.StatsSummaryIntervalItem{
				Date:      i.Interval,
				Visitors:  api.NewOptInt(i.Visitors),
				Pageviews: api.NewOptInt(i.Pageviews),
				Bounces:   api.NewOptInt(i.Bounces),
				Duration:  api.NewOptInt(i.Duration),
			})
		}
	}

	return resp, nil
}

func (h *Handler) GetWebsiteIDPages(ctx context.Context, params api.GetWebsiteIDPagesParams) (api.GetWebsiteIDPagesRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists.
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query.
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:    db.NewFilter(db.FilterPathname, params.Path),
		Referrer:    db.NewFilter(db.FilterReferrer, params.Referrer),
		UTMSource:   db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:     db.NewFilter(db.FilterBrowser, params.Browser),
		OS:          db.NewFilter(db.FilterOS, params.Os),
		Device:      db.NewFilter(db.FilterDevice, params.Device),
		Country:     db.NewFilter(db.FilterCountry, params.Country),
		Language:    db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),

		// Pagination
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	// Check parameter if it is asking for summary.
	switch params.Summary.Value {
	case true:
		// Get summary.
		pages, err := h.analyticsDB.GetWebsitePagesSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website pages summary")
			return ErrInternalServerError(err), nil
		}

		// Create API response.
		res := api.StatsPages{}
		for _, page := range pages {
			res = append(res, api.StatsPagesItem{
				Path:               page.Pathname,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &res, nil
	case false:
		// Get pages
		pages, err := h.analyticsDB.GetWebsitePages(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website pages")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsPages{}
		for _, page := range pages {
			res = append(res, api.StatsPagesItem{
				Path:                page.Pathname,
				Visitors:            page.Visitors,
				VisitorsPercentage:  page.VisitorsPercentage,
				Pageviews:           api.NewOptInt(page.Pageviews),
				PageviewsPercentage: api.NewOptFloat32(page.PageviewsPercentage),
				Bounces:             api.NewOptInt(page.Bounces),
				Duration:            api.NewOptInt(page.Duration),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDTime(ctx context.Context, params api.GetWebsiteIDTimeParams) (api.GetWebsiteIDTimeRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:    db.NewFilter(db.FilterPathname, params.Path),
		Referrer:    db.NewFilter(db.FilterReferrer, params.Referrer),
		UTMSource:   db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:     db.NewFilter(db.FilterBrowser, params.Browser),
		OS:          db.NewFilter(db.FilterOS, params.Os),
		Device:      db.NewFilter(db.FilterDevice, params.Device),
		Country:     db.NewFilter(db.FilterCountry, params.Country),
		Language:    db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),

		// Pagination
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	// Check parameter if it is asking for summary
	switch params.Summary.Value {
	case true:
		// Get summary
		times, err := h.analyticsDB.GetWebsiteTimeSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website time summary")
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
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website time")
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
				Visitors:              api.NewOptInt(page.Visitors),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDReferrers(ctx context.Context, params api.GetWebsiteIDReferrersParams) (api.GetWebsiteIDReferrersRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:    db.NewFilter(db.FilterPathname, params.Path),
		Referrer:    db.NewFilter(db.FilterReferrer, params.Referrer),
		UTMSource:   db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:     db.NewFilter(db.FilterBrowser, params.Browser),
		OS:          db.NewFilter(db.FilterOS, params.Os),
		Device:      db.NewFilter(db.FilterDevice, params.Device),
		Country:     db.NewFilter(db.FilterCountry, params.Country),
		Language:    db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),

		// Pagination
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	// Check parameter if it is asking for summary
	switch params.Summary.Value {
	case true:
		// Get summary
		referrers, err := h.analyticsDB.GetWebsiteReferrersSummary(ctx, params.Grouped.Value, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website referrers summary")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsReferrers{}
		for _, page := range referrers {
			res = append(res, api.StatsReferrersItem{
				Referrer:           page.Referrer,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &res, nil
	case false:
		// Get referrers
		referrers, err := h.analyticsDB.GetWebsiteReferrers(ctx, params.Grouped.Value, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website referrers")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsReferrers{}
		for _, page := range referrers {
			res = append(res, api.StatsReferrersItem{
				Referrer:           page.Referrer,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
				Bounces:            api.NewOptInt(page.Bounces),
				Duration:           api.NewOptInt(page.Duration),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDSources(ctx context.Context, params api.GetWebsiteIDSourcesParams) (api.GetWebsiteIDSourcesRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:    db.NewFilter(db.FilterPathname, params.Path),
		Referrer:    db.NewFilter(db.FilterReferrer, params.Referrer),
		UTMSource:   db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:     db.NewFilter(db.FilterBrowser, params.Browser),
		OS:          db.NewFilter(db.FilterOS, params.Os),
		Device:      db.NewFilter(db.FilterDevice, params.Device),
		Country:     db.NewFilter(db.FilterCountry, params.Country),
		Language:    db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),

		// Pagination
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	switch params.Summary.Value {
	case true:
		// Get summary
		sources, err := h.analyticsDB.GetWebsiteUTMSourcesSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website sources summary")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsUTMSources{}
		for _, page := range sources {
			res = append(res, api.StatsUTMSourcesItem{
				Source:             page.Source,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &res, nil
	case false:
		// Get sources
		sources, err := h.analyticsDB.GetWebsiteUTMSources(ctx, filters)
		if err != nil {
			log.Error().Err(err).Msg("failed to get website utm sources")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsUTMSources{}
		for _, page := range sources {
			res = append(res, api.StatsUTMSourcesItem{
				Source:             page.Source,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
				Bounces:            api.NewOptInt(page.Bounces),
				Duration:           api.NewOptInt(page.Duration),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDMediums(ctx context.Context, params api.GetWebsiteIDMediumsParams) (api.GetWebsiteIDMediumsRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:    db.NewFilter(db.FilterPathname, params.Path),
		Referrer:    db.NewFilter(db.FilterReferrer, params.Referrer),
		UTMSource:   db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:     db.NewFilter(db.FilterBrowser, params.Browser),
		OS:          db.NewFilter(db.FilterOS, params.Os),
		Device:      db.NewFilter(db.FilterDevice, params.Device),
		Country:     db.NewFilter(db.FilterCountry, params.Country),
		Language:    db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),

		// Pagination
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	switch params.Summary.Value {
	case true:
		// Get summary
		mediums, err := h.analyticsDB.GetWebsiteUTMMediumsSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website mediums summary")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsUTMMediums{}
		for _, page := range mediums {
			res = append(res, api.StatsUTMMediumsItem{
				Medium:             page.Medium,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &res, nil
	case false:
		// Get mediums
		mediums, err := h.analyticsDB.GetWebsiteUTMMediums(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website utm mediums")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsUTMMediums{}
		for _, page := range mediums {
			res = append(res, api.StatsUTMMediumsItem{
				Medium:             page.Medium,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
				Bounces:            api.NewOptInt(page.Bounces),
				Duration:           api.NewOptInt(page.Duration),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDCampaigns(ctx context.Context, params api.GetWebsiteIDCampaignsParams) (api.GetWebsiteIDCampaignsRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:    db.NewFilter(db.FilterPathname, params.Path),
		Referrer:    db.NewFilter(db.FilterReferrer, params.Referrer),
		UTMSource:   db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:     db.NewFilter(db.FilterBrowser, params.Browser),
		OS:          db.NewFilter(db.FilterOS, params.Os),
		Device:      db.NewFilter(db.FilterDevice, params.Device),
		Country:     db.NewFilter(db.FilterCountry, params.Country),
		Language:    db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),

		// Pagination
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	switch params.Summary.Value {
	case true:
		// Get summary
		campaigns, err := h.analyticsDB.GetWebsiteUTMCampaignsSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website utm campaigns summary")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsUTMCampaigns{}
		for _, page := range campaigns {
			res = append(res, api.StatsUTMCampaignsItem{
				Campaign:           page.Campaign,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &res, nil
	case false:
		// Get campaigns
		campaigns, err := h.analyticsDB.GetWebsiteUTMCampaigns(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website utm campaigns")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsUTMCampaigns{}
		for _, page := range campaigns {
			res = append(res, api.StatsUTMCampaignsItem{
				Campaign:           page.Campaign,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
				Bounces:            api.NewOptInt(page.Bounces),
				Duration:           api.NewOptInt(page.Duration),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDBrowsers(ctx context.Context, params api.GetWebsiteIDBrowsersParams) (api.GetWebsiteIDBrowsersRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:    db.NewFilter(db.FilterPathname, params.Path),
		Referrer:    db.NewFilter(db.FilterReferrer, params.Referrer),
		UTMSource:   db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:     db.NewFilter(db.FilterBrowser, params.Browser),
		OS:          db.NewFilter(db.FilterOS, params.Os),
		Device:      db.NewFilter(db.FilterDevice, params.Device),
		Country:     db.NewFilter(db.FilterCountry, params.Country),
		Language:    db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),

		// Pagination
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	switch params.Summary.Value {
	case true:
		// Get summary
		browsers, err := h.analyticsDB.GetWebsiteBrowsersSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website browsers summary")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsBrowsers{}
		for _, page := range browsers {
			res = append(res, api.StatsBrowsersItem{
				Browser:            page.Browser,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &res, nil
	case false:
		// Get browsers
		browsers, err := h.analyticsDB.GetWebsiteBrowsers(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website browsers")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsBrowsers{}
		for _, page := range browsers {
			res = append(res, api.StatsBrowsersItem{
				Browser:            page.Browser,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
				Bounces:            api.NewOptInt(page.Bounces),
				Duration:           api.NewOptInt(page.Duration),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDOs(ctx context.Context, params api.GetWebsiteIDOsParams) (api.GetWebsiteIDOsRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:    db.NewFilter(db.FilterPathname, params.Path),
		Referrer:    db.NewFilter(db.FilterReferrer, params.Referrer),
		UTMSource:   db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:     db.NewFilter(db.FilterBrowser, params.Browser),
		OS:          db.NewFilter(db.FilterOS, params.Os),
		Device:      db.NewFilter(db.FilterDevice, params.Device),
		Country:     db.NewFilter(db.FilterCountry, params.Country),
		Language:    db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),

		// Pagination
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	switch params.Summary.Value {
	case true:
		// Get summary
		os, err := h.analyticsDB.GetWebsiteOSSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website os summary")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsOS{}
		for _, page := range os {
			res = append(res, api.StatsOSItem{
				Os:                 page.OS,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &res, nil
	case false:
		// Get OS
		os, err := h.analyticsDB.GetWebsiteOS(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website os")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsOS{}
		for _, page := range os {
			res = append(res, api.StatsOSItem{
				Os:                 page.OS,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
				Bounces:            api.NewOptInt(page.Bounces),
				Duration:           api.NewOptInt(page.Duration),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDDevice(ctx context.Context, params api.GetWebsiteIDDeviceParams) (api.GetWebsiteIDDeviceRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:    db.NewFilter(db.FilterPathname, params.Path),
		Referrer:    db.NewFilter(db.FilterReferrer, params.Referrer),
		UTMSource:   db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:     db.NewFilter(db.FilterBrowser, params.Browser),
		OS:          db.NewFilter(db.FilterOS, params.Os),
		Device:      db.NewFilter(db.FilterDevice, params.Device),
		Country:     db.NewFilter(db.FilterCountry, params.Country),
		Language:    db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),

		// Pagination
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	switch params.Summary.Value {
	case true:
		// Get summary
		devices, err := h.analyticsDB.GetWebsiteDevicesSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website devices summary")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsDevices{}
		for _, page := range devices {
			res = append(res, api.StatsDevicesItem{
				Device:             page.Device,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &res, nil
	case false:
		// Get devices
		devices, err := h.analyticsDB.GetWebsiteDevices(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website devices")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsDevices{}
		for _, page := range devices {
			res = append(res, api.StatsDevicesItem{
				Device:             page.Device,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
				Bounces:            api.NewOptInt(page.Bounces),
				Duration:           api.NewOptInt(page.Duration),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDLanguage(ctx context.Context, params api.GetWebsiteIDLanguageParams) (api.GetWebsiteIDLanguageRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:    db.NewFilter(db.FilterPathname, params.Path),
		Referrer:    db.NewFilter(db.FilterReferrer, params.Referrer),
		UTMSource:   db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:     db.NewFilter(db.FilterBrowser, params.Browser),
		OS:          db.NewFilter(db.FilterOS, params.Os),
		Device:      db.NewFilter(db.FilterDevice, params.Device),
		Country:     db.NewFilter(db.FilterCountry, params.Country),
		Language:    db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),

		// Pagination
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	switch params.Summary.Value {
	case true:
		// Get summary
		languages, err := h.analyticsDB.GetWebsiteLanguagesSummary(ctx, params.Locale.Value, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website languages summary")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsLanguages{}
		for _, page := range languages {
			res = append(res, api.StatsLanguagesItem{
				Language:           page.Language,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &res, nil
	case false:
		// Get languages
		languages, err := h.analyticsDB.GetWebsiteLanguages(ctx, params.Locale.Value, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website languages")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsLanguages{}
		for _, page := range languages {
			res = append(res, api.StatsLanguagesItem{
				Language:           page.Language,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
				Bounces:            api.NewOptInt(page.Bounces),
				Duration:           api.NewOptInt(page.Duration),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}

func (h *Handler) GetWebsiteIDCountry(ctx context.Context, params api.GetWebsiteIDCountryParams) (api.GetWebsiteIDCountryRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query
	filters := &db.Filters{
		Hostname: params.Hostname,

		Pathname:    db.NewFilter(db.FilterPathname, params.Path),
		Referrer:    db.NewFilter(db.FilterReferrer, params.Referrer),
		UTMSource:   db.NewFilter(db.FilterUTMSource, params.UtmSource),
		UTMMedium:   db.NewFilter(db.FilterUTMMedium, params.UtmMedium),
		UTMCampaign: db.NewFilter(db.FilterUTMCampaign, params.UtmCampaign),
		Browser:     db.NewFilter(db.FilterBrowser, params.Browser),
		OS:          db.NewFilter(db.FilterOS, params.Os),
		Device:      db.NewFilter(db.FilterDevice, params.Device),
		Country:     db.NewFilter(db.FilterCountry, params.Country),
		Language:    db.NewFilter(db.FilterLanguage, params.Language),

		// YYYY-MM-DD
		PeriodStart: params.Start.Value.Format(model.DateFormat),
		PeriodEnd:   params.End.Value.Format(model.DateFormat),

		// Pagination
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	switch params.Summary.Value {
	case true:
		// Get summary
		countries, err := h.analyticsDB.GetWebsiteCountriesSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website countries summary")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsCountries{}
		for _, page := range countries {
			res = append(res, api.StatsCountriesItem{
				Country:            page.Country,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &res, nil
	case false:
		// Get countries
		countries, err := h.analyticsDB.GetWebsiteCountries(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website countries")
			return ErrInternalServerError(err), nil
		}

		// Create API response
		res := api.StatsCountries{}
		for _, page := range countries {
			res = append(res, api.StatsCountriesItem{
				Country:            page.Country,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
				Bounces:            api.NewOptInt(page.Bounces),
				Duration:           api.NewOptInt(page.Duration),
			})
		}

		return &res, nil
	default:
		return ErrBadRequest(model.ErrInvalidParameter), nil
	}
}
