package services

import (
	"context"

	"github.com/medama-io/medama/api"
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
	filters := createFilters(params, params.Hostname)

	// Get the properties for the website
	properties, err := h.analyticsDB.GetWebsiteCustomProperties(ctx, filters)
	if err != nil {
		log.Error().Err(err).Msg("failed to get website properties")
		return ErrInternalServerError(model.ErrInternalServerError), nil
	}

	// Return the properties
	resp := api.StatsProperties{}
	for _, p := range properties {
		resp = append(resp, api.StatsPropertiesItem{
			Name:     p.Name,
			Value:    p.Value,
			Events:   p.Events,
			Visitors: p.Visitors,
		})
	}

	return &resp, nil
}
