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
	filters := db.CreateFilters(params, params.Hostname)

	// Get summary
	currentSummary, err := h.analyticsDB.GetWebsiteSummary(ctx, filters)
	if err != nil {
		log.Error().Err(err).Msg("failed to get website summary")
		return ErrInternalServerError(err), nil
	}

	resp := api.StatsSummary{
		Current: api.StatsSummaryCurrent{
			Visitors:         currentSummary.Visitors,
			Pageviews:        currentSummary.Pageviews,
			BouncePercentage: currentSummary.BounceRate,
			Duration:         currentSummary.Duration,
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
				Visitors:         previousSummary.Visitors,
				Pageviews:        previousSummary.Pageviews,
				BouncePercentage: previousSummary.BounceRate,
				Duration:         previousSummary.Duration,
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

		resp.Interval = make([]api.StatsSummaryIntervalItem, 0, len(interval))
		for _, i := range interval {
			resp.Interval = append(resp.Interval, api.StatsSummaryIntervalItem{
				Date:             i.Interval,
				Visitors:         api.NewOptInt(i.Visitors),
				Pageviews:        api.NewOptInt(i.Pageviews),
				BouncePercentage: api.NewOptFloat32(i.BounceRate),
				Duration:         api.NewOptInt(i.Duration),
			})
		}
	}

	return &api.StatsSummaryHeaders{
		Response: resp,
	}, nil
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
	filters := db.CreateFilters(params, params.Hostname)

	// Check parameter if it is asking for summary.
	if params.Summary.Value {
		// Get summary.
		pages, err := h.analyticsDB.GetWebsitePagesSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website pages summary")
			return ErrInternalServerError(err), nil
		}

		resp := make(api.StatsPages, 0, len(pages))
		for _, page := range pages {
			resp = append(resp, api.StatsPagesItem{
				Path:               page.Pathname,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &api.StatsPagesHeaders{
			Response: resp,
		}, nil
	}

	// Get pages
	pages, err := h.analyticsDB.GetWebsitePages(ctx, filters)
	if err != nil {
		log.Error().
			Err(err).
			Bool("summary", params.Summary.Value).
			Msg("failed to get website pages")
		return ErrInternalServerError(err), nil
	}

	resp := make(api.StatsPages, 0, len(pages))
	for _, page := range pages {
		resp = append(resp, api.StatsPagesItem{
			Path:                page.Pathname,
			Visitors:            page.Visitors,
			VisitorsPercentage:  page.VisitorsPercentage,
			Pageviews:           api.NewOptInt(page.Pageviews),
			PageviewsPercentage: api.NewOptFloat32(page.PageviewsPercentage),
			BouncePercentage:    api.NewOptFloat32(page.BounceRate),
			Duration:            api.NewOptInt(page.Duration),
		})
	}

	return &api.StatsPagesHeaders{
		Response: resp,
	}, nil
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
	filters := db.CreateFilters(params, params.Hostname)

	// Check parameter if it is asking for summary
	if params.Summary.Value {
		// Get summary
		times, err := h.analyticsDB.GetWebsiteTimeSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website time summary")
			return ErrInternalServerError(err), nil
		}

		resp := make(api.StatsTime, 0, len(times))
		for _, page := range times {
			resp = append(resp, api.StatsTimeItem{
				Path:               page.Pathname,
				Duration:           page.Duration,
				DurationPercentage: page.DurationPercentage,
			})
		}

		return &api.StatsTimeHeaders{
			Response: resp,
		}, nil
	}

	// Get time
	times, err := h.analyticsDB.GetWebsiteTime(ctx, filters)
	if err != nil {
		log.Error().
			Err(err).
			Bool("summary", params.Summary.Value).
			Msg("failed to get website time")
		return ErrInternalServerError(err), nil
	}

	resp := make(api.StatsTime, 0, len(times))
	for _, page := range times {
		resp = append(resp, api.StatsTimeItem{
			Path:                  page.Pathname,
			Duration:              page.Duration,
			DurationPercentage:    page.DurationPercentage,
			DurationUpperQuartile: api.NewOptInt(page.DurationUpperQuartile),
			DurationLowerQuartile: api.NewOptInt(page.DurationLowerQuartile),
			Visitors:              api.NewOptInt(page.Visitors),
		})
	}

	return &api.StatsTimeHeaders{
		Response: resp,
	}, nil
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
	filters := db.CreateFilters(params, params.Hostname)

	// Check parameter if it is asking for summary
	if params.Summary.Value {
		// Get summary
		referrers, err := h.analyticsDB.GetWebsiteReferrersSummary(ctx, params.Grouped.Value, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website referrers summary")
			return ErrInternalServerError(err), nil
		}

		resp := make(api.StatsReferrers, 0, len(referrers))
		for _, page := range referrers {
			resp = append(resp, api.StatsReferrersItem{
				Referrer:           page.Referrer,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &api.StatsReferrersHeaders{
			Response: resp,
		}, nil
	}

	// Get referrers
	referrers, err := h.analyticsDB.GetWebsiteReferrers(ctx, params.Grouped.Value, filters)
	if err != nil {
		log.Error().
			Err(err).
			Bool("summary", params.Summary.Value).
			Msg("failed to get website referrers")
		return ErrInternalServerError(err), nil
	}

	resp := make(api.StatsReferrers, 0, len(referrers))
	for _, page := range referrers {
		resp = append(resp, api.StatsReferrersItem{
			Referrer:           page.Referrer,
			Visitors:           page.Visitors,
			VisitorsPercentage: page.VisitorsPercentage,
			BouncePercentage:   api.NewOptFloat32(page.BounceRate),
			Duration:           api.NewOptInt(page.Duration),
		})
	}

	return &api.StatsReferrersHeaders{
		Response: resp,
	}, nil
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
	filters := db.CreateFilters(params, params.Hostname)

	if params.Summary.Value {
		// Get summary
		sources, err := h.analyticsDB.GetWebsiteUTMSourcesSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website sources summary")
			return ErrInternalServerError(err), nil
		}

		resp := make(api.StatsUTMSources, 0, len(sources))
		for _, page := range sources {
			resp = append(resp, api.StatsUTMSourcesItem{
				Source:             page.Source,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &api.StatsUTMSourcesHeaders{
			Response: resp,
		}, nil
	}

	// Get sources
	sources, err := h.analyticsDB.GetWebsiteUTMSources(ctx, filters)
	if err != nil {
		log.Error().Err(err).Msg("failed to get website utm sources")
		return ErrInternalServerError(err), nil
	}

	resp := make(api.StatsUTMSources, 0, len(sources))
	for _, page := range sources {
		resp = append(resp, api.StatsUTMSourcesItem{
			Source:             page.Source,
			Visitors:           page.Visitors,
			VisitorsPercentage: page.VisitorsPercentage,
			BouncePercentage:   api.NewOptFloat32(page.BounceRate),
			Duration:           api.NewOptInt(page.Duration),
		})
	}

	return &api.StatsUTMSourcesHeaders{
		Response: resp,
	}, nil
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
	filters := db.CreateFilters(params, params.Hostname)

	if params.Summary.Value {
		// Get summary
		mediums, err := h.analyticsDB.GetWebsiteUTMMediumsSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website mediums summary")
			return ErrInternalServerError(err), nil
		}

		resp := make(api.StatsUTMMediums, 0, len(mediums))
		for _, page := range mediums {
			resp = append(resp, api.StatsUTMMediumsItem{
				Medium:             page.Medium,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &api.StatsUTMMediumsHeaders{
			Response: resp,
		}, nil
	}

	// Get mediums
	mediums, err := h.analyticsDB.GetWebsiteUTMMediums(ctx, filters)
	if err != nil {
		log.Error().
			Err(err).
			Bool("summary", params.Summary.Value).
			Msg("failed to get website utm mediums")
		return ErrInternalServerError(err), nil
	}

	resp := make(api.StatsUTMMediums, 0, len(mediums))
	for _, page := range mediums {
		resp = append(resp, api.StatsUTMMediumsItem{
			Medium:             page.Medium,
			Visitors:           page.Visitors,
			VisitorsPercentage: page.VisitorsPercentage,
			BouncePercentage:   api.NewOptFloat32(page.BounceRate),
			Duration:           api.NewOptInt(page.Duration),
		})
	}

	return &api.StatsUTMMediumsHeaders{
		Response: resp,
	}, nil
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
	filters := db.CreateFilters(params, params.Hostname)

	if params.Summary.Value {
		// Get summary
		campaigns, err := h.analyticsDB.GetWebsiteUTMCampaignsSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website utm campaigns summary")
			return ErrInternalServerError(err), nil
		}

		resp := make(api.StatsUTMCampaigns, 0, len(campaigns))
		for _, page := range campaigns {
			resp = append(resp, api.StatsUTMCampaignsItem{
				Campaign:           page.Campaign,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &api.StatsUTMCampaignsHeaders{
			Response: resp,
		}, nil
	}

	// Get campaigns
	campaigns, err := h.analyticsDB.GetWebsiteUTMCampaigns(ctx, filters)
	if err != nil {
		log.Error().
			Err(err).
			Bool("summary", params.Summary.Value).
			Msg("failed to get website utm campaigns")
		return ErrInternalServerError(err), nil
	}

	resp := make(api.StatsUTMCampaigns, 0, len(campaigns))
	for _, page := range campaigns {
		resp = append(resp, api.StatsUTMCampaignsItem{
			Campaign:           page.Campaign,
			Visitors:           page.Visitors,
			VisitorsPercentage: page.VisitorsPercentage,
			BouncePercentage:   api.NewOptFloat32(page.BounceRate),
			Duration:           api.NewOptInt(page.Duration),
		})
	}

	return &api.StatsUTMCampaignsHeaders{
		Response: resp,
	}, nil
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
	filters := db.CreateFilters(params, params.Hostname)

	if params.Summary.Value {
		// Get summary
		browsers, err := h.analyticsDB.GetWebsiteBrowsersSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website browsers summary")
			return ErrInternalServerError(err), nil
		}

		resp := make(api.StatsBrowsers, 0, len(browsers))
		for _, page := range browsers {
			resp = append(resp, api.StatsBrowsersItem{
				Browser:            page.Browser,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &api.StatsBrowsersHeaders{
			Response: resp,
		}, nil
	}

	// Get browsers
	browsers, err := h.analyticsDB.GetWebsiteBrowsers(ctx, filters)
	if err != nil {
		log.Error().
			Err(err).
			Bool("summary", params.Summary.Value).
			Msg("failed to get website browsers")
		return ErrInternalServerError(err), nil
	}

	resp := make(api.StatsBrowsers, 0, len(browsers))
	for _, page := range browsers {
		resp = append(resp, api.StatsBrowsersItem{
			Browser:            page.Browser,
			Visitors:           page.Visitors,
			VisitorsPercentage: page.VisitorsPercentage,
			BouncePercentage:   api.NewOptFloat32(page.BounceRate),
			Duration:           api.NewOptInt(page.Duration),
		})
	}

	return &api.StatsBrowsersHeaders{
		Response: resp,
	}, nil
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
	filters := db.CreateFilters(params, params.Hostname)

	if params.Summary.Value {
		// Get summary
		os, err := h.analyticsDB.GetWebsiteOSSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website os summary")
			return ErrInternalServerError(err), nil
		}

		resp := make(api.StatsOS, 0, len(os))
		for _, page := range os {
			resp = append(resp, api.StatsOSItem{
				Os:                 page.OS,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &api.StatsOSHeaders{
			Response: resp,
		}, nil
	}

	// Get OS
	os, err := h.analyticsDB.GetWebsiteOS(ctx, filters)
	if err != nil {
		log.Error().
			Err(err).
			Bool("summary", params.Summary.Value).
			Msg("failed to get website os")
		return ErrInternalServerError(err), nil
	}

	resp := make(api.StatsOS, 0, len(os))
	for _, page := range os {
		resp = append(resp, api.StatsOSItem{
			Os:                 page.OS,
			Visitors:           page.Visitors,
			VisitorsPercentage: page.VisitorsPercentage,
			BouncePercentage:   api.NewOptFloat32(page.BounceRate),
			Duration:           api.NewOptInt(page.Duration),
		})
	}

	return &api.StatsOSHeaders{
		Response: resp,
	}, nil
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
	filters := db.CreateFilters(params, params.Hostname)

	if params.Summary.Value {
		// Get summary
		devices, err := h.analyticsDB.GetWebsiteDevicesSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website devices summary")
			return ErrInternalServerError(err), nil
		}

		resp := make(api.StatsDevices, 0, len(devices))
		for _, page := range devices {
			resp = append(resp, api.StatsDevicesItem{
				Device:             page.Device,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &api.StatsDevicesHeaders{
			Response: resp,
		}, nil
	}

	// Get devices
	devices, err := h.analyticsDB.GetWebsiteDevices(ctx, filters)
	if err != nil {
		log.Error().
			Err(err).
			Bool("summary", params.Summary.Value).
			Msg("failed to get website devices")
		return ErrInternalServerError(err), nil
	}

	resp := make(api.StatsDevices, 0, len(devices))
	for _, page := range devices {
		resp = append(resp, api.StatsDevicesItem{
			Device:             page.Device,
			Visitors:           page.Visitors,
			VisitorsPercentage: page.VisitorsPercentage,
			BouncePercentage:   api.NewOptFloat32(page.BounceRate),
			Duration:           api.NewOptInt(page.Duration),
		})
	}

	return &api.StatsDevicesHeaders{
		Response: resp,
	}, nil
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
	filters := db.CreateFilters(params, params.Hostname)

	if params.Summary.Value {
		// Get summary
		languages, err := h.analyticsDB.GetWebsiteLanguagesSummary(ctx, params.Locale.Value, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website languages summary")
			return ErrInternalServerError(err), nil
		}

		resp := make(api.StatsLanguages, 0, len(languages))
		for _, page := range languages {
			resp = append(resp, api.StatsLanguagesItem{
				Language:           page.Language,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &api.StatsLanguagesHeaders{
			Response: resp,
		}, nil
	}

	// Get languages
	languages, err := h.analyticsDB.GetWebsiteLanguages(ctx, params.Locale.Value, filters)
	if err != nil {
		log.Error().
			Err(err).
			Bool("summary", params.Summary.Value).
			Msg("failed to get website languages")
		return ErrInternalServerError(err), nil
	}

	resp := make(api.StatsLanguages, 0, len(languages))
	for _, page := range languages {
		resp = append(resp, api.StatsLanguagesItem{
			Language:           page.Language,
			Visitors:           page.Visitors,
			VisitorsPercentage: page.VisitorsPercentage,
			BouncePercentage:   api.NewOptFloat32(page.BounceRate),
			Duration:           api.NewOptInt(page.Duration),
		})
	}

	return &api.StatsLanguagesHeaders{
		Response: resp,
	}, nil
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
	filters := db.CreateFilters(params, params.Hostname)

	if params.Summary.Value {
		// Get summary
		countries, err := h.analyticsDB.GetWebsiteCountriesSummary(ctx, filters)
		if err != nil {
			log.Error().
				Err(err).
				Bool("summary", params.Summary.Value).
				Msg("failed to get website countries summary")
			return ErrInternalServerError(err), nil
		}

		resp := make(api.StatsCountries, 0, len(countries))
		for _, page := range countries {
			resp = append(resp, api.StatsCountriesItem{
				Country:            page.Country,
				Visitors:           page.Visitors,
				VisitorsPercentage: page.VisitorsPercentage,
			})
		}

		return &api.StatsCountriesHeaders{
			Response: resp,
		}, nil
	}

	// Get countries
	countries, err := h.analyticsDB.GetWebsiteCountries(ctx, filters)
	if err != nil {
		log.Error().
			Err(err).
			Bool("summary", params.Summary.Value).
			Msg("failed to get website countries")
		return ErrInternalServerError(err), nil
	}

	resp := make(api.StatsCountries, 0, len(countries))
	for _, page := range countries {
		resp = append(resp, api.StatsCountriesItem{
			Country:            page.Country,
			Visitors:           page.Visitors,
			VisitorsPercentage: page.VisitorsPercentage,
			BouncePercentage:   api.NewOptFloat32(page.BounceRate),
			Duration:           api.NewOptInt(page.Duration),
		})
	}

	return &api.StatsCountriesHeaders{
		Response: resp,
	}, nil
}
