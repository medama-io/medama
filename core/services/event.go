package services

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

func (h *Handler) GetEventPing(ctx context.Context, params api.GetEventPingParams) (api.GetEventPingRes, error) {
	return nil, nil
}

func (h *Handler) PostEventHit(ctx context.Context, req *api.EventHit, params api.PostEventHitParams) (api.PostEventHitRes, error) {
	attributes := []slog.Attr{}

	// Get request from context
	reqBody, ok := ctx.Value(model.RequestKeyBody).(*http.Request)
	if !ok {
		slog.LogAttrs(ctx, slog.LevelError, "failed to get request from context", attributes...)
		return nil, model.ErrInternalServerError
	}

	// Get values
	dateCreated := time.Now().Unix()

	event := &model.PageView{
		BID:         req.B,
		EventType:   req.E,
		Hostname:    reqBody.Host,
		Pathname:    reqBody.URL.Path,
		Referrer:    req.R.Value,
		DateCreated: dateCreated,
	}

	attributes = append(attributes,
		slog.String("bid", event.BID),
		slog.String("event_type", event.EventType),
		slog.String("hostname", event.Hostname),
		slog.String("pathname", event.Pathname),
		slog.String("referrer", event.Referrer),
		slog.Int64("date_created", dateCreated),
	)

	// Add to database
	err := h.analyticsDB.AddPageView(ctx, event)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to add page view", attributes...)
		return nil, model.ErrInternalServerError
	}

	// Log success
	slog.LogAttrs(ctx, slog.LevelDebug, "added page view", attributes...)

	return &api.PostEventHitOK{}, nil
}
