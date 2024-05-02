package services

import (
	"context"

	"github.com/go-faster/errors"
	tz "github.com/medama-io/go-timezone-country"
	"github.com/medama-io/go-useragent"
	"github.com/medama-io/medama/db/duckdb"
	"github.com/medama-io/medama/db/sqlite"
	"github.com/medama-io/medama/util"
)

type Handler struct {
	auth           *util.AuthService
	db             *sqlite.Client
	analyticsDB    *duckdb.Client
	useragent      *useragent.Parser
	timezoneMap    *tz.TimezoneCodeMap
	codeCountryMap *tz.CodeCountryMap
	hostnames      *util.CacheStore
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

	// Load hostname cache
	hostnameCache := util.NewCacheStore()
	hostnames, err := sqlite.ListAllHostnames(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "services init")
	}
	hostnameCache.AddAll(hostnames)

	return &Handler{
		auth:           auth,
		db:             sqlite,
		analyticsDB:    duckdb,
		useragent:      useragent.NewParser(),
		timezoneMap:    &tzMap,
		codeCountryMap: &codeCountryMap,
		hostnames:      &hostnameCache,
	}, nil
}
