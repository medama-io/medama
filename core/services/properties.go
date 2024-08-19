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
	filters.IsEvent = true

	// Get the properties for the website
	properties, err := h.analyticsDB.GetWebsiteCustomProperties(ctx, filters)
	if err != nil {
		log.Error().Err(err).Msg("failed to get website properties")
		return ErrInternalServerError(model.ErrInternalServerError), nil
	}

	// Return the properties in desired format
	respMap := make(map[string]*api.StatsPropertiesItem)

	for _, p := range properties {
		item, exists := respMap[p.Name]
		if !exists {
			item = &api.StatsPropertiesItem{
				Name:  p.Name,
				Items: []api.StatsPropertiesItemItemsItem{},
			}
			respMap[p.Name] = item
		}

		// Aggregate the values
		item.Events += p.Events
		item.Visitors += p.Visitors

		// Add the property
		item.Items = append(item.Items, api.StatsPropertiesItemItemsItem{
			Value:    p.Value,
			Events:   p.Events,
			Visitors: p.Visitors,
		})
	}

	resp := make(api.StatsProperties, 0, len(respMap))
	for _, item := range respMap {
		resp = append(resp, *item)
	}

	return &resp, nil
}
