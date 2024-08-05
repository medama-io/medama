package services

import (
	"context"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/medama-io/go-referrer-parser"
	tz "github.com/medama-io/go-timezone-country"
	"github.com/medama-io/go-useragent"
	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util"
)

// This is a runtime config that is read from user settings in the database.
type RuntimeConfig struct {
	// Tracker settings.
	// Choose what type of script to serve from /script.js.
	//
	// Options:
	//
	// - "default" - Default script that collects page view data.
	//
	// - "tagged-events" - Script that collects page view data and custom event properties.
	ScriptType string
	// Usage settings.
	// Number of threads to use for processing events.
	Threads int
	// Memory limit for processing events.
	MemoryLimit string
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
	runtimeConfig *RuntimeConfig
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

	runtimeConfig, err := NewRuntimeConfig(ctx, sqlite, duckdb)
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
		runtimeConfig:  runtimeConfig,
	}, nil
}

// NewRuntimeConfig creates a new runtime config.
func NewRuntimeConfig(ctx context.Context, user *sqlite.Client, analytics *duckdb.Client) (*RuntimeConfig, error) {
	// Load the script type from the database.
	settings, err := user.GetSettings(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "runtime config")
	}

	// Set the DuckDB settings.
	if settings.MemoryLimit != "" || settings.Threads != 0 {
		err := analytics.SetDuckDBSettings(ctx, &settings.DuckDBSettings)
		if err != nil {
			return nil, errors.Wrap(err, "runtime config")
		}
	}

	metadata, err := analytics.GetDuckDBSettings(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "runtime config")
	}

	return &RuntimeConfig{
		ScriptType:  settings.ScriptType,
		Threads:     metadata.Threads,
		MemoryLimit: metadata.MemoryLimit,
	}, nil
}

func (r *RuntimeConfig) UpdateConfig(ctx context.Context, meta *sqlite.Client, analytics *duckdb.Client, settings *model.GlobalSettings) error {
	if settings.ScriptType != "" {
		err := meta.UpdateSetting(ctx, model.SettingsKeyScriptType, settings.ScriptType)
		if err != nil {
			return errors.Wrap(err, "script type update config")
		}
		r.ScriptType = settings.ScriptType
	}

	if settings.Threads != 0 {
		err := meta.UpdateSetting(ctx, model.SettingsKeyThreads, strconv.Itoa(settings.Threads))
		if err != nil {
			return errors.Wrap(err, "threads update config")
		}
		r.Threads = settings.Threads
	}

	if settings.MemoryLimit != "" {
		err := meta.UpdateSetting(ctx, model.SettingsKeyMemoryLimit, settings.MemoryLimit)
		if err != nil {
			return errors.Wrap(err, "memory limit update config")
		}
		r.MemoryLimit = settings.MemoryLimit
	}

	// Update metadata settings to persist choices.
	return analytics.SetDuckDBSettings(ctx, &model.DuckDBSettings{
		Threads:     r.Threads,
		MemoryLimit: r.MemoryLimit,
	})
}
