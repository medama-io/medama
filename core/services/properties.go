package services

import (
	"context"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/db"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
)

func (h *Handler) GetWebsiteIDProperties(ctx context.Context, params api.GetWebsiteIDPropertiesParams) (api.GetWebsiteIDPropertiesRes, error) {
	log := logger.Get().With().Str("hostname", params.Hostname).Logger()

	// Check if website exists
	exists := h.hostnames.Has(params.Hostname)
	if !exists {
		log.Debug().Msg("website not found")
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Create filter for database query
	filters := db.CreateFilters(params, params.Hostname)
	filters.SortByEventDates = true

	// Get the properties for the website
	properties, err := h.analyticsDB.GetWebsiteCustomProperties(ctx, filters)
	if err != nil {
		log.Error().Err(err).Msg("failed to get website properties")
		return ErrInternalServerError(model.ErrInternalServerError), nil
	}

	resp := make(api.StatsProperties, 0, len(properties))
	for _, p := range properties {
		item := api.StatsPropertiesItem{
			Events:           p.Events,
			EventsPercentage: p.EventsPercentage,
		}

		if p.Name != "" {
			item.Name = api.NewOptString(p.Name)
		}

		if p.Value != "" {
			item.Value = api.NewOptString(p.Value)
		}

		resp = append(resp, item)
	}

	return &resp, nil
}
