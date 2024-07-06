package services

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

const (
	// OneDay is the duration of one day.
	OneDay = 24 * time.Hour
	// Set to no-cache to disable caching.
	CacheControl = "no-cache"
	// Unknown is the default value for unknown fields.
	Unknown = "Unknown"
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
		log := logger.Get()
		log.Error().Err(err).Msg("failed to parse if-modified-since header")
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
	log := logger.Get()

	switch req.Type {
	case api.EventLoadEventHit:
		hostname := req.EventLoad.U.Hostname()
		log = log.With().Str("hostname", hostname).Logger()

		// Verify hostname exists
		if !h.hostnames.Has(hostname) {
			log.Warn().Msg("hit: website not found")
			return ErrNotFound(model.ErrWebsiteNotFound), nil
		}

		pathname := req.EventLoad.U.Path
		// Remove trailing slash if it exists
		if pathname != "/" {
			pathname = strings.TrimSuffix(pathname, "/")
		}

		// Get request from context
		reqBody, ok := ctx.Value(model.RequestKeyBody).(*http.Request)
		if !ok {
			log.Error().Msg("hit: failed to get request key from context")
			return ErrInternalServerError(model.ErrRequestContext), nil
		}

		// Parse user agent first to catch early if it is a bot.
		rawUserAgent := reqBody.Header.Get("User-Agent")
		ua := h.useragent.Parse(rawUserAgent)

		// If the user agent is a bot, we want to ignore it.
		if ua.Bot {
			log.Debug().Str("user_agent", rawUserAgent).Msg("hit: user agent is a bot")
			return &api.PostEventHitNoContent{}, nil
		}

		uaBrowser := ua.Browser
		if uaBrowser == "" {
			uaBrowser = Unknown
		}

		uaOS := ua.OS
		if uaOS == "" {
			uaOS = Unknown
		}

		uaDevice := Unknown
		switch {
		case ua.Desktop:
			uaDevice = "Desktop"
		case ua.Mobile:
			uaDevice = "Mobile"
		case ua.Tablet:
			uaDevice = "Tablet"
		case ua.TV:
			uaDevice = "TV"
		}

		if ua.Browser == "" || ua.OS == "" || uaDevice == Unknown {
			log.Debug().Str("user_agent", rawUserAgent).Msg("hit: unknown user agent")
		}

		if ua.Browser == "" && ua.OS == "" && uaDevice == Unknown {
			// Do not log the event if every element of the user agent is unknown.
			return &api.PostEventHitNoContent{}, nil
		}

		// Parse referrer URL and remove any query parameters or self-referencing
		// hostnames.
		referrerHost := ""
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

		referrerGroup := ""
		if referrerHost != "" {
			// Get the referrer group from the referrer URL.
			referrerGroup = h.referrer.Parse(referrerHost)
		}

		// Get country code from user's timezone. This is used as a best effort
		// to determine the country of the user's location without compromising
		// their privacy using IP addresses.
		var countryName string
		countryCode, err := h.timezoneMap.GetCode(req.EventLoad.T.Value)
		if err != nil {
			log.Debug().Err(err).Msg("hit: failed to get country code from timezone")
			countryCode = ""
			countryName = Unknown
		}

		if countryCode != "" {
			countryName, err = h.codeCountryMap.GetCountry(countryCode)
			if err != nil {
				log.Debug().Err(err).Msg("hit: failed to get country name from country code")
				countryName = Unknown
			}
		}

		// Get users language from Accept-Language header
		languages, _, err := language.ParseAcceptLanguage(reqBody.Header.Get("Accept-Language"))
		if err != nil {
			log.Debug().Err(err).Msg("hit: failed to parse accept language header")
		}

		// Get the first language from the list which is the most preferred and convert it to a language name
		languageBase := Unknown
		languageDialect := Unknown
		if len(languages) > 0 {
			// Narrow down the language to the base language (e.g. en-US -> en)
			base, _ := languages[0].Base()
			languageBase = display.English.Tags().Name(language.Make(base.String()))
			languageDialect = display.English.Tags().Name(languages[0])
		}

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
			ReferrerHost:    referrerHost,
			ReferrerGroup:   referrerGroup,
			Country:         countryName,
			LanguageBase:    languageBase,
			LanguageDialect: languageDialect,

			BrowserName: uaBrowser,
			OS:          uaOS,
			DeviceType:  uaDevice,

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
			Str("referrer_host", event.ReferrerHost).
			Str("referrer_group", event.ReferrerGroup).
			Str("country", countryName).
			Str("language_base", event.LanguageBase).
			Str("language_dialect", event.LanguageDialect).
			Str("browser_name", event.BrowserName).
			Str("os", event.OS).
			Str("device_type", event.DeviceType).
			Logger()

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

		log = log.With().
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
		log.Error().Str("type", string(req.Type)).Msg("hit: invalid event hit type")
		return ErrBadRequest(model.ErrInvalidTrackerEvent), nil
	}

	return &api.PostEventHitNoContent{}, nil
}
