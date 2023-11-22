package services

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

const OneDay = 24 * time.Hour

func (h *Handler) GetEventPing(ctx context.Context, params api.GetEventPingParams) (api.GetEventPingRes, error) {
	attributes := []slog.Attr{}
	// Check if if-modified-since header is set
	ifModified := params.IfModifiedSince.Value
	attributes = append(attributes, slog.String("if_modified", ifModified))

	// If it is not set, it is a unique user.
	if ifModified == "" {
		// Get current day but reset the time to 00:00:00
		currentDay := time.Now().UTC().Truncate(OneDay)

		// Return body to activate caching which is the number of seconds
		body := strings.NewReader(strconv.Itoa(currentDay.Second()))

		lastModified := currentDay.Format(http.TimeFormat)

		attributes = append(attributes, slog.String("last_modified", lastModified))
		slog.LogAttrs(ctx, slog.LevelDebug, "last modified header not set", attributes...)

		return &api.GetEventPingOKHeaders{
			LastModified: lastModified,
			Response:     api.GetEventPingOK{Data: body},
		}, nil
	}

	// Otherwise, this is not a unique user.
	// Increment the last modified date by one second.
	lastModifiedTime, err := time.Parse(http.TimeFormat, ifModified)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to parse if modified since header", attributes...)
		return ErrBadRequest(err), nil
	}
	lastModifiedTime = lastModifiedTime.Add(time.Second)

	// Return body to activate caching
	body := strings.NewReader(strconv.Itoa(lastModifiedTime.Second()))

	lastModified := lastModifiedTime.Format(http.TimeFormat)
	attributes = append(attributes, slog.String("last_modified", lastModified))
	slog.LogAttrs(ctx, slog.LevelDebug, "last modified header set", attributes...)
	return &api.GetEventPingOKHeaders{
		LastModified: lastModified,
		Response:     api.GetEventPingOK{Data: body},
	}, nil
}

func (h *Handler) PostEventHit(ctx context.Context, req *api.EventHit, params api.PostEventHitParams) (api.PostEventHitRes, error) {
	attributes := []slog.Attr{}

	// Split url into hostname and pathname
	u, err := url.Parse(req.U)
	if err != nil {
		attributes = append(attributes, slog.String("url", req.U), slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to parse url", attributes...)
		return ErrBadRequest(err), nil
	}
	hostname := u.Hostname()
	pathname := u.Path

	// Verify hostname exists
	exists, err := h.db.WebsiteExists(ctx, hostname)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to check if website exists", attributes...)
		return ErrInternalServerError(err), nil
	}
	if !exists {
		attributes = append(attributes, slog.String("hostname", hostname))
		slog.LogAttrs(ctx, slog.LevelWarn, "website not found", attributes...)
		// return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Get request from context
	reqBody, ok := ctx.Value(model.RequestKeyBody).(*http.Request)
	if !ok {
		slog.LogAttrs(ctx, slog.LevelError, "failed to get request from context", attributes...)
		return ErrInternalServerError(errors.New("failed to get request from context")), nil
	}

	// Get users language from Accept-Language header
	acceptLanguage := reqBody.Header.Get("Accept-Language")

	// Parse user agent
	rawUserAgent := reqBody.Header.Get("User-Agent")
	ua := h.useragent.Parse(rawUserAgent)

	uaBrowser := model.NewBrowserName(ua.Browser)
	uaVersion := ua.GetMajorVersion()
	uaOS := model.NewOSName(ua.OS)
	uaDeviceType := model.NewDeviceType(ua.Desktop, ua.Mobile, ua.Tablet, ua.TV)
	isUnknownUA := false
	// If there are unfilled fields, we want to mark this as an unknown user agent
	// and store the raw user agent string.
	if uaBrowser == 0 || uaOS == 0 || uaDeviceType == 0 || uaVersion == "" {
		isUnknownUA = true
	}

	// Verify screen height and width for overflow as we store them as uint16
	// in the database.
	screenHeight := uint16(req.H.Value)
	screenWidth := uint16(req.W.Value)
	if screenHeight > 65535 || screenWidth > 65535 {
		attributes = append(attributes, slog.Int("screen_height", req.H.Value), slog.Int("screen_width", req.W.Value))
		slog.LogAttrs(ctx, slog.LevelDebug, "screen height or width is too large", attributes...)
		return ErrBadRequest(errors.New("screen height or width is too large")), nil
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
			Referrer:       req.R.Value,
			Title:          req.T.Value,
			Timezone:       req.D.Value,
			Language:       acceptLanguage,
			BrowserName:    uaBrowser,
			BrowserVersion: uaVersion,
			OS:             uaOS,
			DeviceType:     uaDeviceType,

			ScreenWidth:  uint16(req.W.Value),
			ScreenHeight: uint16(req.H.Value),
			DateCreated:  dateCreated,
		}

		// If the user agent was unable to be parsed, store the raw user agent
		// string.
		if isUnknownUA {
			event.RawUserAgent = rawUserAgent
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
			slog.String("browser_name", event.BrowserName.String()),
			slog.String("browser_version", event.BrowserVersion),
			slog.String("os", event.OS.String()),
			slog.String("device_type", event.DeviceType.String()),
			slog.String("raw_user_agent", event.RawUserAgent),
			slog.Int("screen_width", int(event.ScreenWidth)),
			slog.Int("screen_height", int(event.ScreenHeight)),
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
