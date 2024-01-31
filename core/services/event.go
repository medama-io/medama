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
	// CorsMethods is the methods to allow for CORS.
	CorsMethods = "GET"
	// CorsHeaders is the headers to allow for CORS.
	CorsHeaders = "If-Modified-Since, Content-Type"
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
	currentDay := time.Now().Truncate(OneDay)

	// If it is not set, it is a unique user.
	if ifModified == "" {
		// Return body to activate caching.
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

		// Return body to activate caching.
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
		AccessControlAllowOrigin:  CorsOrigin,
		AccessControlAllowMethods: CorsMethods,
		AccessControlAllowHeaders: CorsHeaders,
		LastModified:              ifModified,
		CacheControl:              CacheControl,
		Response: api.GetEventPingOK{
			Data: body,
		},
	}, nil
}

func (h *Handler) PostEventHit(ctx context.Context, req api.EventHit, params api.PostEventHitParams) (api.PostEventHitRes, error) {
	attributes := []slog.Attr{}

	switch req.Type {
	case api.EventLoadEventHit:
		hostname := req.EventLoad.U.Hostname()
		pathname := req.EventLoad.U.Path

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

		// Parse referrer URL and remove any query parameters or self-referencing
		// hostnames.
		var referrerHost string
		if req.EventLoad.R.Value != "" {
			referrer, err := url.Parse(req.EventLoad.R.Value)
			if err != nil {
				attributes = append(attributes, slog.String("error", err.Error()))
				slog.LogAttrs(ctx, slog.LevelError, "failed to parse referrer URL", attributes...)
				return ErrBadRequest(err), nil
			}

			// If the referrer hostname is the same as the current hostname, we
			// want to remove it.
			referrerHost = referrer.Hostname()
			if referrerHost == hostname {
				referrerHost = ""
			}
		}

		// Get country code from user's timezone. This is used as a best effort
		// to determine the country of the user's location without compromising
		// their privacy using IP addresses.
		countryCode, err := h.timezoneMap.GetCode(req.EventLoad.T.Value)
		if err != nil {
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelDebug, "failed to get country code from timezone", attributes...)
			countryCode = ""
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
		uaOS := model.NewOSName(ua.OS)
		uaDeviceType := model.NewDeviceType(ua.Desktop, ua.Mobile, ua.Tablet, ua.TV)

		// Get utm source, medium, and campaigm from URL query parameters.
		queries := req.EventLoad.U.Query()
		utmSource := queries.Get("utm_source")
		utmMedium := queries.Get("utm_medium")
		utmCampaign := queries.Get("utm_campaign")

		event := &model.PageViewHit{
			// Required
			BID:          req.EventLoad.B,
			Hostname:     hostname,
			Pathname:     pathname,
			IsUniqueUser: req.EventLoad.P,
			IsUniquePage: req.EventLoad.Q,
			// Optional
			Referrer:    referrerHost,
			CountryCode: countryCode,
			Language:    acceptLanguage,

			BrowserName: uaBrowser,
			OS:          uaOS,
			DeviceType:  uaDeviceType,

			UTMSource:   utmSource,
			UTMMedium:   utmMedium,
			UTMCampaign: utmCampaign,
		}

		attributes = append(attributes,
			slog.String("bid", event.BID),
			slog.String("event_type", string(req.Type)),
			slog.String("hostname", event.Hostname),
			slog.String("pathname", event.Pathname),
			slog.Bool("is_unique_user", event.IsUniqueUser),
			slog.Bool("is_unique_page", event.IsUniquePage),
			slog.String("referrer", event.Referrer),
			slog.String("country_code", countryCode),
			slog.String("language", event.Language),
			slog.String("browser_name", event.BrowserName.String()),
			slog.String("os", event.OS.String()),
			slog.String("device_type", event.DeviceType.String()),
		)

		// TODO: Remove temporary raw user agent logging for debugging
		if event.BrowserName == model.UnknownBrowser || event.OS == model.UnknownOS || event.DeviceType == model.UnknownDevice {
			attributes = append(attributes, slog.String("user_agent", rawUserAgent))
		}

		err = h.analyticsDB.AddPageView(ctx, event)
		if err != nil {
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to add page view", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Log success
		slog.LogAttrs(ctx, slog.LevelDebug, "added page view", attributes...)
	case api.EventUnloadEventHit:
		event := &model.PageViewDuration{
			BID:        req.EventUnload.B,
			DurationMs: req.EventUnload.M,
		}

		attributes = append(attributes,
			slog.String("bid", event.BID),
			slog.String("event_type", string(req.Type)),
			slog.Int("duration_ms", event.DurationMs),
		)

		err := h.analyticsDB.UpdatePageView(ctx, event)
		if err != nil {
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to update page view", attributes...)
			return ErrInternalServerError(err), nil
		}

		// Log success
		slog.LogAttrs(ctx, slog.LevelDebug, "updated page view", attributes...)

	default:
		attributes = append(attributes, slog.String("type", string(req.Type)))
		slog.LogAttrs(ctx, slog.LevelError, "invalid event hit type", attributes...)
		return ErrBadRequest(model.ErrInvalidTrackerEvent), nil
	}

	return &api.PostEventHitNoContent{}, nil
}
