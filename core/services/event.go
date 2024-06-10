package services

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

const (
	// OneDay is the duration of one day.
	OneDay = 24 * time.Hour
	// Set to no-cache to disable caching.
	CacheControl = "no-cache"
)

func (h *Handler) GetEventPing(ctx context.Context, params api.GetEventPingParams) (api.GetEventPingRes, error) {
	// Check if if-modified-since header is set
	ifModified := params.IfModifiedSince.Value

	// Get current day but reset the time to 00:00:00
	currentDay := time.Now().Truncate(OneDay)

	// If it is not set, it is a unique user.
	if ifModified == "" {
		// Return body to activate caching.
		body := strings.NewReader("0")

		lastModified := currentDay.Format(http.TimeFormat)

		return &api.GetEventPingOKHeaders{
			LastModified: lastModified,
			CacheControl: CacheControl,
			Response:     api.GetEventPingOK{Data: body},
		}, nil
	}

	// Parse the if-modified-since header and check if it is older than a day.
	lastModifiedTime, err := time.Parse(http.TimeFormat, ifModified)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to parse if-modified-since header")
		return ErrBadRequest(err), nil
	}

	// If the last modified time is one day ago, we want to reset the cache
	// and mark as a unique user.
	if lastModifiedTime.Before(currentDay) {
		lastModifiedTime = currentDay

		// Return body to activate caching.
		body := strings.NewReader("0")

		lastModified := lastModifiedTime.Format(http.TimeFormat)

		return &api.GetEventPingOKHeaders{
			LastModified: lastModified,
			CacheControl: CacheControl,
			Response:     api.GetEventPingOK{Data: body},
		}, nil
	}

	// Otherwise, this is not a unique user.
	body := strings.NewReader("1")

	// Return not modified if the last modified time is today (not unique user).
	return &api.GetEventPingOKHeaders{
		LastModified: ifModified,
		CacheControl: CacheControl,
		Response: api.GetEventPingOK{
			Data: body,
		},
	}, nil
}

func (h *Handler) PostEventHit(ctx context.Context, req api.EventHit, params api.PostEventHitParams) (api.PostEventHitRes, error) {
	switch req.Type {
	case api.EventLoadEventHit:
		hostname := req.EventLoad.U.Hostname()
		pathname := req.EventLoad.U.Path
		log := zerolog.Ctx(ctx).With().Str("hostname", hostname).Logger()

		// Verify hostname exists
		exists := h.hostnames.Has(hostname)
		if !exists {
			log.Warn().Msg("hit: website not found")
			return ErrNotFound(model.ErrWebsiteNotFound), nil
		}

		// Parse referrer URL and remove any query parameters or self-referencing
		// hostnames.
		var referrerHost string
		if req.EventLoad.R.Value != "" {
			referrer, err := url.Parse(req.EventLoad.R.Value)
			if err != nil {
				log.Warn().Err(err).Msg("hit: failed to parse referrer URL")
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
			log.Debug().Err(err).Msg("hit: failed to get country code from timezone")
			countryCode = ""
		}

		// Get request from context
		reqBody, ok := ctx.Value(model.RequestKeyBody).(*http.Request)
		if !ok {
			log.Error().Msg("hit: failed to get request key from context")
			return ErrInternalServerError(model.ErrRequestContext), nil
		}

		// Get users language from Accept-Language header
		languages, _, err := language.ParseAcceptLanguage(reqBody.Header.Get("Accept-Language"))
		if err != nil {
			log.Debug().Err(err).Msg("hit: failed to parse accept language header")
		}
		// Get the first language from the list which is the most preferred and convert it to a language name
		language := "Unknown"
		if len(languages) > 0 {
			language = display.English.Tags().Name(languages[0])
		}

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
			Language:    language,

			BrowserName: uaBrowser,
			OS:          uaOS,
			DeviceType:  uaDeviceType,

			UTMSource:   utmSource,
			UTMMedium:   utmMedium,
			UTMCampaign: utmCampaign,
		}

		log = log.With().
			Str("bid", event.BID).
			Str("event_type", string(req.Type)).
			Str("pathname", event.Pathname).
			Bool("is_unique_user", event.IsUniqueUser).
			Bool("is_unique_page", event.IsUniquePage).
			Str("referrer", event.Referrer).
			Str("country_code", countryCode).
			Str("language", event.Language).
			Str("browser_name", event.BrowserName.String()).
			Str("os", event.OS.String()).
			Str("device_type", event.DeviceType.String()).
			Logger()

		// TODO: Remove temporary raw user agent logging for debugging
		if event.BrowserName == model.UnknownBrowser || event.OS == model.UnknownOS || event.DeviceType == model.UnknownDevice {
			log.Debug().Str("user_agent", rawUserAgent).Msg("hit: unknown user agent")
		}

		err = h.analyticsDB.AddPageView(ctx, event)
		if err != nil {
			log.Error().Err(err).Msg("hit: failed to add page view")
			return ErrInternalServerError(err), nil
		}

		// Log success
		log.Debug().Msg("hit: added page view")
	case api.EventUnloadEventHit:
		event := &model.PageViewDuration{
			BID:        req.EventUnload.B,
			DurationMs: req.EventUnload.M,
		}

		log := zerolog.Ctx(ctx).With().
			Str("bid", event.BID).
			Str("event_type", string(req.Type)).
			Int("duration_ms", event.DurationMs).
			Logger()

		err := h.analyticsDB.UpdatePageView(ctx, event)
		if err != nil {
			log.Error().Err(err).Msg("hit: failed to update page view")
			return ErrInternalServerError(err), nil
		}

		// Log success
		log.Debug().Msg("hit: updated page view")

	default:
		zerolog.Ctx(ctx).Error().Str("type", string(req.Type)).Msg("hit: invalid event hit type")
		return ErrBadRequest(model.ErrInvalidTrackerEvent), nil
	}

	return &api.PostEventHitNoContent{}, nil
}
