package services

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/iputils"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
	"go.jetify.com/typeid"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

const (
	// OneDay is the duration of one day.
	OneDay = 24 * time.Hour
	// Set to no-cache to disable caching.
	NoCache = "no-cache"
	// Unknown is the default value for unknown fields.
	Unknown = "Unknown"

	// Constants for unique user tracking.
	Zero = "0"
	One  = "1"

	// Cache control max-age prefix.
	maxAgePrefix       = "max-age="
	int64MaxBufferSize = 20
	base10             = 10
)

// readerPool is a pool of strings.Reader to reduce memory allocations.
var readerPool = sync.Pool{
	New: func() any {
		return new(strings.Reader)
	},
}

func getReader(s string) *strings.Reader {
	r := readerPool.Get().(*strings.Reader)
	r.Reset(s)
	return r
}

func putReader(r *strings.Reader) {
	readerPool.Put(r)
}

func (h *Handler) GetEventPing(
	_ctx context.Context,
	params api.GetEventPingParams,
) (api.GetEventPingRes, error) {
	// Check if if-modified-since header is set
	ifModified := params.IfModifiedSince.Value

	// Get current day but reset the time to 00:00:00
	currentDay := time.Now().Truncate(OneDay)

	// If it is not set, it is a unique user.
	if ifModified == "" {
		// Return "0" response body to activate caching.
		body := getReader(Zero)
		defer putReader(body)

		lastModified := currentDay.Format(http.TimeFormat)

		return &api.GetEventPingOKHeaders{
			LastModified: lastModified,
			CacheControl: NoCache,
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
		// Return "0" body to activate caching.
		body := getReader(Zero)
		defer putReader(body)

		return &api.GetEventPingOKHeaders{
			LastModified: currentDay.Format(http.TimeFormat),
			CacheControl: NoCache, // Keep no-cache for unique users
			Response:     api.GetEventPingOK{Data: body},
		}, nil
	}

	// Otherwise, this is not a unique user.
	body := getReader(One)
	defer putReader(body)

	// Calculate time until lastModifiedTime + OneDay
	nextResetTime := lastModifiedTime.Add(OneDay)
	secondsUntilReset := int(time.Until(nextResetTime).Seconds())

	// Preallocate cache control header with max length and then append the max age.
	maxLen := len(maxAgePrefix) + int64MaxBufferSize // 20 digits is enough for int64.MaxValue
	cacheControl := make([]byte, 0, maxLen)
	cacheControl = append(cacheControl, maxAgePrefix...)
	cacheControl = strconv.AppendInt(cacheControl, int64(secondsUntilReset), base10)

	// Return not modified if the last modified time is today (not unique user).
	return &api.GetEventPingOKHeaders{
		LastModified: ifModified,
		CacheControl: string(cacheControl),
		Response: api.GetEventPingOK{
			Data: body,
		},
	}, nil
}

const (
	// IsBotThreshold is the threshold of unknown metrics for determining if a
	// user agent is a bot.
	IsBotThreshold = 2
)

func (h *Handler) PostEventHit(
	ctx context.Context,
	req api.EventHit,
	_params api.PostEventHitParams,
) (api.PostEventHitRes, error) {
	log := logger.Get()

	reqBody, ok := ctx.Value(model.RequestKeyBody).(*http.Request)
	if !ok {
		log.Error().Msg("hit: failed to get request key from context")
		return ErrInternalServerError(model.ErrRequestContext), nil
	}

	clientIP, err := iputils.GetIP(reqBody)
	if err != nil {
		log.Debug().Err(err).Msg("hit: failed to extract client IP")
		return ErrBadRequest(err), nil
	}

	// Check if IP is blocked
	if h.RuntimeConfig.IPFilter.HasIP(clientIP) {
		log.Debug().Msg("hit: client IP is blocked")
		return &api.PostEventHitNoContent{}, nil
	}

	// If this counter exceeds 2, we want to return early as the event is likely
	// a bot.
	//
	// Ensure all functions that increment this counter occur at the beginning
	// rather than the end of the function to bail out early.
	unknownCounter := 0

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

		// Parse user agent first to catch early if it is a bot.
		rawUserAgent := reqBody.Header.Get("User-Agent")
		ua := h.useragent.Parse(rawUserAgent)

		// If the user agent is a bot, we want to ignore it.
		if ua.IsBot() {
			log.Debug().Str("user_agent", rawUserAgent).Msg("hit: user agent is a bot")
			return &api.PostEventHitNoContent{}, nil
		}

		uaBrowser := ua.Browser()
		if uaBrowser == "" {
			uaBrowser = Unknown
			unknownCounter++
		}

		uaOS := ua.OS()
		if uaOS == "" {
			uaOS = Unknown
			unknownCounter++
		}

		uaDevice := ua.Device()
		if uaDevice == "" {
			uaDevice = Unknown
			unknownCounter++
		}

		if uaBrowser == Unknown || uaOS == Unknown || uaDevice == Unknown {
			log.Debug().Str("user_agent", rawUserAgent).Msg("hit: unknown user agent")
			if unknownCounter >= IsBotThreshold {
				return &api.PostEventHitNoContent{}, nil
			}
		}

		// Get country code from user's timezone. This is used as a best effort
		// to determine the country of the user's location without compromising
		// their privacy using IP addresses.
		countryName, err := h.timezoneCountryMap.GetCountry(req.EventLoad.T.Value)
		if err != nil {
			log.Debug().Err(err).Msg("hit: failed to get country name from timezone")

			unknownCounter++
			if unknownCounter >= IsBotThreshold {
				return &api.PostEventHitNoContent{}, nil
			}

			countryName = Unknown
		}

		// Get users language from Accept-Language header
		languages, _, err := language.ParseAcceptLanguage(reqBody.Header.Get("Accept-Language"))
		if err != nil {
			log.Debug().Err(err).Msg("hit: failed to parse accept language header")

			unknownCounter++
			if unknownCounter >= IsBotThreshold {
				return &api.PostEventHitNoContent{}, nil
			}
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

		// Parse referrer URL and extract the host and group name.
		referrer, err := h.referrer.Parse(req.EventLoad.R.Value, hostname)
		if err != nil {
			log.Warn().Err(err).Msg("hit: failed to parse referrer URL")
			return ErrBadRequest(err), nil
		}

		// If the referrer is spam, we want to ignore it.
		if referrer.IsSpam {
			log.Debug().Str("referrer_host", referrer.Host).Msg("hit: referrer is spam")
			return &api.PostEventHitNoContent{}, nil
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
			ReferrerHost:    referrer.Host,
			ReferrerGroup:   referrer.Group,
			Country:         countryName,
			LanguageBase:    languageBase,
			LanguageDialect: languageDialect,

			BrowserName: uaBrowser.String(),
			OS:          uaOS.String(),
			DeviceType:  uaDevice.String(),

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

		if req.EventLoad.D.IsSet() {
			// Generate batch ID to group all the properties of the same event.
			batchIDType, err := typeid.WithPrefix("event")
			if err != nil {
				return nil, errors.Wrap(err, "typeid custom event")
			}
			batchID := batchIDType.String()

			events := make([]model.EventHit, 0, len(req.EventLoad.D.Value))

			for name, item := range req.EventLoad.D.Value {
				var value string

				switch item.Type {
				case api.StringEventLoadDItem:
					value = item.String
				case api.IntEventLoadDItem:
					value = strconv.Itoa(item.Int)
				case api.BoolEventLoadDItem:
					value = strconv.FormatBool(item.Bool)
				default:
					return nil, errors.New(
						"invalid custom event property type: " + string(item.Type),
					)
				}

				events = append(events, model.EventHit{
					BID:     event.BID,
					BatchID: batchID,
					Group:   hostname,
					Name:    name,
					Value:   value,
				})
			}

			log = log.With().
				Str("event_type", string(req.Type)).
				Int("event_count", len(events)).
				Logger()

			err = h.analyticsDB.AddPageView(ctx, event, &events)
			if err != nil {
				log.Error().Err(err).Msg("hit: failed to add page view")
				return ErrInternalServerError(err), nil
			}
		} else {
			err = h.analyticsDB.AddPageView(ctx, event, nil)
			if err != nil {
				log.Error().Err(err).Msg("hit: failed to add page view")
				return ErrInternalServerError(err), nil
			}
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

	case api.EventCustomEventHit:
		if (len(req.EventCustom.D)) == 0 {
			return ErrBadRequest(model.ErrInvalidProperties), nil
		}

		group := req.EventCustom.G
		log = log.With().Str("group_name", group).Logger()

		// Verify hostname exists as hostname is used as the group name.
		if !h.hostnames.Has(group) {
			log.Warn().Msg("hit: website not found")
			return ErrNotFound(model.ErrWebsiteNotFound), nil
		}

		// Generate batch ID to group all the properties of the same event.
		batchIDType, err := typeid.WithPrefix("event")
		if err != nil {
			return nil, errors.Wrap(err, "typeid custom event")
		}
		batchID := batchIDType.String()

		events := make([]model.EventHit, 0, len(req.EventCustom.D))

		for name, item := range req.EventCustom.D {
			var value string

			switch item.Type {
			case api.StringEventCustomDItem:
				value = item.String
			case api.IntEventCustomDItem:
				value = strconv.Itoa(item.Int)
			case api.BoolEventCustomDItem:
				value = strconv.FormatBool(item.Bool)
			default:
				return nil, errors.New("invalid custom event property type: " + string(item.Type))
			}

			events = append(events, model.EventHit{
				BID:     req.EventCustom.B.Or(""),
				BatchID: batchID,
				Group:   group,
				Name:    name,
				Value:   value,
			})
		}

		log = log.With().
			Str("event_type", string(req.Type)).
			Str("group", group).
			Int("event_count", len(events)).
			Logger()

		err = h.analyticsDB.AddEvents(ctx, &events)
		if err != nil {
			log.Error().Err(err).Msg("hit: failed to add event")
			return ErrInternalServerError(err), nil
		}

		log.Debug().Msg("hit: added custom events")
	default:
		log.Error().Str("type", string(req.Type)).Msg("hit: invalid event hit type")
		return ErrBadRequest(model.ErrInvalidTrackerEvent), nil
	}

	return &api.PostEventHitNoContent{}, nil
}
