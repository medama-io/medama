package services

import (
	"context"
	"slices"
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

// This is a runtime config that is read from user settings in the database.
type RuntimeConfig struct {
	// Tracker settings.
	// Choose what features of script to serve from /script.js.
	ScriptFileName string
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
		RuntimeConfig:  &runtimeConfig,
	}, nil
}

// NewRuntimeConfig creates a new runtime config.
func NewRuntimeConfig(ctx context.Context, user *sqlite.Client) (RuntimeConfig, error) {
	// Load the script type from the database.
	settings, err := user.GetSettings(ctx)
	if err != nil {
		return RuntimeConfig{}, errors.Wrap(err, "runtime config")
	}

	return RuntimeConfig{
		ScriptFileName: convertScriptType(settings.ScriptType),
	}, nil
}

func (r *RuntimeConfig) UpdateConfig(ctx context.Context, meta *sqlite.Client, settings *model.UserSettings) error {
	if settings.ScriptType != "" {
		err := meta.UpdateSetting(ctx, model.SettingsKeyScriptType, settings.ScriptType)
		if err != nil {
			return errors.Wrap(err, "script type update config")
		}
		r.ScriptFileName = convertScriptType(settings.ScriptType)

		log := logger.Get()
		log.Warn().Str("script_type", settings.ScriptType).Msg("updated script type")
	}

	return nil
}

// Convert array of script type features split by comma to a script file name.
func convertScriptType(scriptType string) string {
	features := strings.Split(scriptType, ",")

	// Hot path for basic script.
	if scriptType == "default" || len(features) == 0 {
		return "/scripts/default.js"
	}

	filteredFeatures := make([]string, 0, len(features))

	// Filter out the default feature.
	for _, feature := range features {
		if feature != "default" {
			filteredFeatures = append(filteredFeatures, feature)
		}
	}

	// Alphabetically sort the features as script files are named alphabetically.
	slices.Sort(filteredFeatures)

	var sb strings.Builder
	sb.WriteString("/scripts/")
	sb.WriteString(strings.Join(filteredFeatures, "."))
	sb.WriteString(".js")

	return sb.String()
}
