package services

import (
	"context"
	"strings"

	"github.com/go-faster/errors"
	"github.com/medama-io/go-referrer-parser"
	tz "github.com/medama-io/go-timezone-country"
	"github.com/medama-io/go-useragent"
	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util"
	"github.com/medama-io/medama/util/logger"
)

type ScriptType struct {
	// Default script that collects page view data.
	Default bool
	// Script that collects page view data and custom event properties.
	TaggedEvent bool
}

// This is a runtime config that is read from user settings in the database.
type RuntimeConfig struct {
	// Tracker settings.
	// Choose what features of script to serve from /script.js.
	ScriptType ScriptType
}

type Handler struct {
	auth        *util.AuthService
	db          *sqlite.Client
	analyticsDB *duckdb.Client

	// Parsing libraries
	useragent      *useragent.Parser
	referrer       *referrer.Parser
	timezoneMap    *tz.TimezoneCodeMap
	codeCountryMap *tz.CodeCountryMap

	// Cache store for hostnames
	hostnames *util.CacheStore

	// Runtime config
	RuntimeConfig *RuntimeConfig
}

// NewService returns a new instance of the ogen service handler.
func NewService(ctx context.Context, auth *util.AuthService, sqlite *sqlite.Client, duckdb *duckdb.Client) (*Handler, error) {
	// Load timezone and country maps
	tzMap, err := tz.NewTimezoneCodeMap()
	if err != nil {
		return nil, errors.Wrap(err, "services init")
	}

	codeCountryMap, err := tz.NewCodeCountryMap()
	if err != nil {
		return nil, errors.Wrap(err, "services init")
	}

	// Load referrer parser
	referrerParser, err := referrer.NewParser()
	if err != nil {
		return nil, errors.Wrap(err, "services init")
	}

	// Load hostname cache
	hostnameCache := util.NewCacheStore()
	hostnames, err := sqlite.ListAllHostnames(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "services init")
	}
	hostnameCache.AddAll(hostnames)

	runtimeConfig, err := NewRuntimeConfig(ctx, sqlite)
	if err != nil {
		return nil, errors.Wrap(err, "services init")
	}

	return &Handler{
		auth:           auth,
		db:             sqlite,
		analyticsDB:    duckdb,
		useragent:      useragent.NewParser(),
		referrer:       referrerParser,
		timezoneMap:    &tzMap,
		codeCountryMap: &codeCountryMap,
		hostnames:      &hostnameCache,
		RuntimeConfig:  runtimeConfig,
	}, nil
}

// NewRuntimeConfig creates a new runtime config.
func NewRuntimeConfig(ctx context.Context, user *sqlite.Client) (*RuntimeConfig, error) {
	// Load the script type from the database.
	settings, err := user.GetSettings(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "runtime config")
	}

	return &RuntimeConfig{
		ScriptType: convertScriptType(settings.ScriptType),
	}, nil
}

func (r *RuntimeConfig) UpdateConfig(ctx context.Context, meta *sqlite.Client, settings *model.UserSettings) error {
	if settings.ScriptType != "" {
		err := meta.UpdateSetting(ctx, model.SettingsKeyScriptType, settings.ScriptType)
		if err != nil {
			return errors.Wrap(err, "script type update config")
		}
		r.ScriptType = convertScriptType(settings.ScriptType)

		log := logger.Get()
		log.Warn().Str("script_type", settings.ScriptType).Msg("updated script type")
	}

	return nil
}

// Convert array of script type features split by comma to a ScriptType struct.
func convertScriptType(scriptType string) ScriptType {
	features := strings.Split(scriptType, ",")

	types := ScriptType{}
	for _, feature := range features {
		switch feature {
		case "default":
			types.Default = true
		case "tagged-events":
			types.TaggedEvent = true
		}
	}

	return types
}
