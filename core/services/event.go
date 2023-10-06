package services

import (
	"context"

	"github.com/medama-io/medama/api"
)

func (h *Handler) GetEventPing(ctx context.Context, params api.GetEventPingParams) (*api.GetEventPingOK, error) {
	return nil, nil
}

func (h *Handler) PostEventHit(ctx context.Context, req api.OptEventHit, params api.PostEventHitParams) (api.PostEventHitRes, error) {
	return nil, nil
}
