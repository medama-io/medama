package services

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
)

func (h *Handler) DeleteWebsitesID(ctx context.Context, params api.DeleteWebsitesIDParams) (api.DeleteWebsitesIDRes, error) {
	log := logger.Get()
	if h.auth.IsDemoMode {
		log.Debug().Msg("delete website rejected in demo mode")
		return ErrForbidden(model.ErrDemoMode), nil
	}

	// Check if user owns website
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}

	websites, err := h.db.ListWebsites(ctx, userId)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteNotFound) {
			return ErrNotFound(err), nil
		}

		return nil, errors.Wrap(err, "services")
	}

	var website *model.Website
	for _, w := range websites {
		if w.Hostname == params.Hostname {
			website = w
			break
		}
	}

	if website == nil {
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Delete all views associated with website
	err = h.analyticsDB.DeleteWebsite(ctx, params.Hostname)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteNotFound) {
			return ErrNotFound(err), nil
		}

		return nil, errors.Wrap(err, "services")
	}

	// Delete website
	err = h.db.DeleteWebsite(ctx, params.Hostname)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteNotFound) {
			return ErrNotFound(err), nil
		}

		return nil, errors.Wrap(err, "services")
	}

	// Remove website from hostname cache
	h.hostnames.Remove(params.Hostname)

	return &api.DeleteWebsitesIDNoContent{}, nil
}

func (h *Handler) GetWebsites(ctx context.Context, params api.GetWebsitesParams) (api.GetWebsitesRes, error) {
	// Get user ID from context
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}

	websites, err := h.db.ListWebsites(ctx, userId)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteNotFound) {
			return ErrNotFound(err), nil
		}
		return nil, errors.Wrap(err, "services")
	}

	// Map to API response
	websitesGet := &api.GetWebsitesOKApplicationJSON{}

	// If summary is requested, include visitors per website
	if ok := params.Summary.Or(false); ok {
		for _, w := range websites {
			views, err := h.analyticsDB.GetWebsiteSummaryLast24Hours(ctx, w.Hostname)
			if err != nil {
				if errors.Is(err, model.ErrWebsiteNotFound) {
					return ErrNotFound(err), nil
				}

				return nil, errors.Wrap(err, w.Hostname)
			}

			*websitesGet = append(*websitesGet, api.WebsiteGet{
				Hostname: w.Hostname,
				Summary: api.NewOptWebsiteGetSummary(api.WebsiteGetSummary{
					Visitors: views.Visitors,
				}),
			})
		}
		// Otherwise, return only hostnames
	} else {
		for _, w := range websites {
			*websitesGet = append(*websitesGet, api.WebsiteGet{
				Hostname: w.Hostname,
			})
		}
	}

	return websitesGet, nil
}

func (h *Handler) GetWebsitesID(ctx context.Context, params api.GetWebsitesIDParams) (api.GetWebsitesIDRes, error) {
	// Get user ID from context
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}

	website, err := h.db.GetWebsite(ctx, params.Hostname)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteNotFound) {
			return ErrNotFound(err), nil
		}

		return nil, errors.Wrap(err, "services")
	}

	if website.UserID != userId {
		return ErrUnauthorised(model.ErrWebsiteNotFound), nil
	}

	return &api.WebsiteGet{
		Hostname: website.Hostname,
	}, nil
}

func (h *Handler) PatchWebsitesID(ctx context.Context, req *api.WebsitePatch, params api.PatchWebsitesIDParams) (api.PatchWebsitesIDRes, error) {
	log := logger.Get()
	if h.auth.IsDemoMode {
		log.Debug().Msg("patch website rejected in demo mode")
		return ErrForbidden(model.ErrDemoMode), nil
	}

	// Get user ID from context
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}

	website, err := h.db.GetWebsite(ctx, params.Hostname)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteNotFound) {
			return ErrNotFound(err), nil
		}

		return nil, errors.Wrap(err, "services")
	}

	if website.UserID != userId {
		return ErrUnauthorised(model.ErrWebsiteNotFound), nil
	}

	// Update values
	if req.Hostname.Value != "" {
		website.Hostname = req.Hostname.Value
	}

	website.DateUpdated = time.Now().Unix()

	// Update website
	err = h.db.UpdateWebsite(ctx, website)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteNotFound) {
			return ErrNotFound(err), nil
		}

		return nil, errors.Wrap(err, "services")
	}

	// If hostname was updated, remove old hostname from cache
	// and add new hostname to cache
	if req.Hostname.Value != "" {
		h.hostnames.Remove(params.Hostname)
		h.hostnames.Add(req.Hostname.Value)
	}

	return &api.WebsiteGet{
		Hostname: website.Hostname,
	}, nil
}

func (h *Handler) PostWebsites(ctx context.Context, req *api.WebsiteCreate) (api.PostWebsitesRes, error) {
	log := logger.Get()
	if h.auth.IsDemoMode {
		log.Debug().Msg("post website rejected in demo mode")
		return ErrForbidden(model.ErrDemoMode), nil
	}

	// Get user ID from context
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}

	// Create website
	dateCreated := time.Now().Unix()
	websiteCreate := model.NewWebsite(
		userId,
		req.Hostname,
		dateCreated,
		dateCreated,
	)

	err := h.db.CreateWebsite(ctx, websiteCreate)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteExists) {
			return ErrConflict(err), nil
		}

		return nil, errors.Wrap(err, "services")
	}

	// Add hostname to cache
	h.hostnames.Add(req.Hostname)

	return &api.WebsiteGet{
		Hostname: req.Hostname,
	}, nil
}
