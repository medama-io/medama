package services

import (
	"context"
	"errors"
	"time"

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
		return nil, err
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

		return nil, err
	}

	return nil, nil
}

func (h *Handler) GetWebsites(ctx context.Context, params api.GetWebsitesParams) (api.GetWebsitesRes, error) {
	// Get user ID from context
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}

	websites, err := h.db.ListWebsites(ctx, userId)
	if err != nil {
		return nil, err
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

		return nil, err
	}

	if website.UserID != userId {
		return ErrUnauthorised(model.ErrWebsiteNotFound), nil
	}

	return &api.WebsiteGet{
		Name:     website.Hostname,
		Hostname: website.Hostname,
	}, nil
}

func (h *Handler) PatchWebsitesID(ctx context.Context, req api.OptWebsitePatch, params api.PatchWebsitesIDParams) (api.PatchWebsitesIDRes, error) {
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

		return nil, err
	}

	if website.UserID != userId {
		return ErrUnauthorised(model.ErrWebsiteNotFound), nil
	}

	// Update values
	if req.Value.Hostname.Value != "" {
		website.Hostname = req.Value.Hostname.Value
	}

	if req.Value.Name.Value != "" {
		website.Name = req.Value.Name.Value
	}

	website.DateUpdated = time.Now().Unix()

	// Update website
	err = h.db.UpdateWebsite(ctx, website)
	if err != nil {
		if errors.Is(err, model.ErrWebsiteNotFound) {
			return ErrNotFound(err), nil
		}

		return nil, err
	}

	return nil, nil
}

func (h *Handler) PostWebsites(ctx context.Context, req api.OptWebsiteCreate) (api.PostWebsitesRes, error) {
	// Get user ID from context
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}

	// Create website
	dateCreated := time.Now().Unix()
	websiteCreate := model.NewWebsite(
		userId,
		req.Value.Hostname,
		req.Value.Name,
		dateCreated,
		dateCreated,
	)

	err := h.db.CreateWebsite(ctx, websiteCreate)
	if err != nil {
		return nil, err
	}

	return &api.WebsiteGet{
		Name:     req.Value.Hostname,
		Hostname: req.Value.Hostname,
	}, nil
}
