package services

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

func (h *Handler) DeleteWebsitesID(ctx context.Context, params api.DeleteWebsitesIDParams) (api.DeleteWebsitesIDRes, error) {
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

	// Delete website
	err = h.db.DeleteWebsite(ctx, params.Hostname)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteNotFound) {
			return ErrNotFound(err), nil
		}

		return nil, errors.Wrap(err, "services")
	}

	// Delete all views associated with website
	err = h.analyticsDB.DeleteWebsite(ctx, params.Hostname)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteNotFound) {
			return ErrNotFound(err), nil
		}

		return nil, errors.Wrap(err, "services")
	}

	return &api.DeleteWebsitesIDOK{}, nil
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
	for _, w := range websites {
		*websitesGet = append(*websitesGet, api.WebsiteGet{
			Name:     w.Hostname,
			Hostname: w.Hostname,
		})
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
		Name:     website.Hostname,
		Hostname: website.Hostname,
	}, nil
}

func (h *Handler) PatchWebsitesID(ctx context.Context, req *api.WebsitePatch, params api.PatchWebsitesIDParams) (api.PatchWebsitesIDRes, error) {
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

	if req.Name.Value != "" {
		website.Name = req.Name.Value
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

	return &api.WebsiteGet{
		Name:     website.Hostname,
		Hostname: website.Hostname,
	}, nil
}

func (h *Handler) PostWebsites(ctx context.Context, req *api.WebsiteCreate) (api.PostWebsitesRes, error) {
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
		req.Name,
		dateCreated,
		dateCreated,
	)

	err := h.db.CreateWebsite(ctx, websiteCreate)
	if err != nil {
		return nil, errors.Wrap(err, "services")
	}

	return &api.WebsiteGet{
		Name:     req.Hostname,
		Hostname: req.Hostname,
	}, nil
}
