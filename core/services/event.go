package services

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

const (
	// CorsOrigin is the origin to allow for CORS.
	CorsOrigin = "*"

	// OneDay is the duration of one day.
	OneDay = 24 * time.Hour
	// Set to no-cache to disable caching.
	CacheControl = "no-cache"
)

func (h *Handler) GetEventPing(ctx context.Context, params api.GetEventPingParams) (api.GetEventPingRes, error) {
	attributes := []slog.Attr{}
	// Check if if-modified-since header is set
	ifModified := params.IfModifiedSince.Value
	attributes = append(attributes, slog.String("if_modified", ifModified))

	// Get current day but reset the time to 00:00:00
	currentDay := time.Now().UTC().Truncate(OneDay)

	// If it is not set, it is a unique user.
	if ifModified == "" {
		// Return body to activate caching which is the number of seconds
		body := strings.NewReader("0")

		lastModified := currentDay.Format(http.TimeFormat)

		attributes = append(attributes, slog.String("last_modified", lastModified))
		slog.LogAttrs(ctx, slog.LevelDebug, "last modified header not set", attributes...)

		return &api.GetEventPingOKHeaders{
			AccessControlAllowOrigin: CorsOrigin,
			LastModified:             lastModified,
			CacheControl:             CacheControl,
			Response:                 api.GetEventPingOK{Data: body},
		}, nil
	}

	// Parse the if-modified-since header and check if it is older than a day.
	lastModifiedTime, err := time.Parse(http.TimeFormat, ifModified)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to parse if modified since header", attributes...)
		return ErrBadRequest(err), nil
	}

	// If the last modified time is one day ago, we want to reset the cache
	// and mark as a unique user.
	if lastModifiedTime.Before(currentDay) {
		lastModifiedTime = currentDay

		// Return body to activate caching which is the number of seconds
		body := strings.NewReader("0")

		lastModified := lastModifiedTime.Format(http.TimeFormat)

		attributes = append(attributes, slog.String("last_modified", lastModified))
		slog.LogAttrs(ctx, slog.LevelDebug, "last modified header set", attributes...)

		return &api.GetEventPingOKHeaders{
			AccessControlAllowOrigin: CorsOrigin,
			LastModified:             lastModified,
			CacheControl:             CacheControl,
			Response:                 api.GetEventPingOK{Data: body},
		}, nil
	}

	// Otherwise, this is not a unique user.
	body := strings.NewReader("1")

	// Return not modified if the last modified time is today (not unique user).
	slog.LogAttrs(ctx, slog.LevelDebug, "last modified header set", attributes...)
	return &api.GetEventPingOKHeaders{
		AccessControlAllowOrigin: CorsOrigin,
		LastModified:             ifModified,
		CacheControl:             CacheControl,
		Response: api.GetEventPingOK{
			Data: body,
		},
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
		return ErrNotFound(model.ErrWebsiteNotFound), nil
	}

	// Add to database
	switch req.E {
	case "load":
		// If is unique is not set, default to true
		isUnique, exists := req.P.Get()
		if !exists {
			isUnique = true
		}

		// Get country code from user's timezone. This is used as a best effort
		// to determine the country of the user's location without compromising
		// their privacy using IP addresses.
		countryCode, err := h.timezoneMap.GetCode(req.D.Value)
		if err != nil {
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to get country code from timezone", attributes...)
			return ErrInternalServerError(model.ErrInvalidTimezone), nil
		}

		// Get request from context
		reqBody, ok := ctx.Value(model.RequestKeyBody).(*http.Request)
		if !ok {
			slog.LogAttrs(ctx, slog.LevelError, "failed to get request from context", attributes...)
			return ErrInternalServerError(model.ErrRequestContext), nil
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
			return ErrBadRequest(model.ErrInvalidScreenSize), nil
		}

		// Get utm source, medium, and campaigm from URL query parameters.
		queries := u.Query()
		utmSource := queries.Get("utm_source")
		utmMedium := queries.Get("utm_medium")
		utmCampaign := queries.Get("utm_campaign")

		event := &model.PageView{
			// Required
			BID:      req.B,
			Hostname: hostname,
			Pathname: pathname,
			// Optional
			IsUnique:       isUnique,
			Referrer:       req.R.Value,
			Title:          req.T.Value,
			CountryCode:    countryCode,
			Language:       acceptLanguage,
			BrowserName:    uaBrowser,
			BrowserVersion: uaVersion,
			OS:             uaOS,
			DeviceType:     uaDeviceType,

			ScreenWidth:  screenWidth,
			ScreenHeight: screenHeight,

			UTMSource:   utmSource,
			UTMMedium:   utmMedium,
			UTMCampaign: utmCampaign,
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
			slog.Bool("is_unique", event.IsUnique),
			slog.String("referrer", event.Referrer),
			slog.String("title", event.Title),
			slog.String("country_code", countryCode),
			slog.String("language", event.Language),
			slog.String("browser_name", event.BrowserName.String()),
			slog.String("browser_version", event.BrowserVersion),
			slog.String("os", event.OS.String()),
			slog.String("device_type", event.DeviceType.String()),
			slog.String("raw_user_agent", event.RawUserAgent),
			slog.Int("screen_width", int(event.ScreenWidth)),
			slog.Int("screen_height", int(event.ScreenHeight)),
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
		return ErrBadRequest(model.ErrInvalidTrackerEvent), nil
	}

	return &api.PostEventHitOK{}, nil
}
