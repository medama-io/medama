package services

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

func (h *Handler) GetEventPing(ctx context.Context, params api.GetEventPingParams) (api.GetEventPingRes, error) {
	return &api.GetEventPingOK{}, nil
}

func (h *Handler) PostEventHit(ctx context.Context, req *api.EventHit, params api.PostEventHitParams) (api.PostEventHitRes, error) {
	attributes := []slog.Attr{}

	// Get request from context
	reqBody, ok := ctx.Value(model.RequestKeyBody).(*http.Request)
	if !ok {
		slog.LogAttrs(ctx, slog.LevelError, "failed to get request from context", attributes...)
		return nil, model.ErrInternalServerError
	}

	// Get users language from Accept-Language header
	acceptLanguage := reqBody.Header.Get("Accept-Language")

	// Split url into hostname and pathname
	u, err := url.Parse(req.U)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to parse url", attributes...)
		return nil, model.ErrInternalServerError
	}
	hostname := u.Hostname()
	pathname := u.Path

	// Verify hostname exists
	exists, err := h.db.WebsiteExists(ctx, hostname)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to check if website exists", attributes...)
		return nil, model.ErrInternalServerError
	}
	if !exists {
		attributes = append(attributes, slog.String("hostname", hostname))
		slog.LogAttrs(ctx, slog.LevelError, "website not found", attributes...)
		return ErrNotFound(err), nil
	}

	// Add to database
	switch req.E {
	case "load":
		// Get date created
		dateCreated := time.Now().Unix()

		event := &model.PageView{
			// Required
			BID:      req.B,
			Hostname: hostname,
			Pathname: pathname,
			// Optional
			Referrer:     req.R.Value,
			Title:        req.T.Value,
			Timezone:     req.D.Value,
			Language:     acceptLanguage,
			ScreenWidth:  req.W.Value,
			ScreenHeight: req.H.Value,
			DateCreated:  dateCreated,
		}

		attributes = append(attributes,
			slog.String("bid", event.BID),
			slog.String("event_type", req.E),
			slog.String("hostname", event.Hostname),
			slog.String("pathname", event.Pathname),
			slog.String("referrer", event.Referrer),
			slog.String("title", event.Title),
			slog.String("timezone", event.Timezone),
			slog.String("language", event.Language),
			slog.Int("screen_width", event.ScreenWidth),
			slog.Int("screen_height", event.ScreenHeight),
			slog.Int64("date_created", event.DateCreated),
		)

		err = h.analyticsDB.AddPageView(ctx, event)
		if err != nil {
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to add page view", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Log success
		slog.LogAttrs(ctx, slog.LevelDebug, "added page view", attributes...)
	case "pagehide", "unload", "hidden":
		event := &model.PageViewUpdate{
			BID:        req.B,
			DurationMs: req.M.Value,
		}

		attributes = append(attributes,
			slog.String("bid", event.BID),
			slog.String("event_type", req.E),
			slog.Int("duration_ms", event.DurationMs),
		)

		err = h.analyticsDB.UpdatePageView(ctx, event)
		if err != nil {
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to update page view", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Log success
		slog.LogAttrs(ctx, slog.LevelDebug, "updated page view", attributes...)
	default:
		attributes = append(attributes, slog.String("event_type", req.E))
		slog.LogAttrs(ctx, slog.LevelError, "invalid event type", attributes...)
		return ErrBadRequest(err), nil
	}

	return &api.PostEventHitOK{}, nil
}
