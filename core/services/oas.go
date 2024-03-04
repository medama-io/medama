package services

import (
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
}

// NewService returns a new instance of the ogen service handler.
func NewService(auth *util.AuthService, sqlite *sqlite.Client, duckdb *duckdb.Client) (*Handler, error) {
	tzMap, err := tz.NewTimezoneCodeMap()
	if err != nil {
		return nil, errors.Wrap(err, "services")
	}

	codeCountryMap, err := tz.NewCodeCountryMap()
	if err != nil {
		return nil, errors.Wrap(err, "services")
	}

	return &Handler{
		auth:           auth,
		db:             sqlite,
		analyticsDB:    duckdb,
		useragent:      useragent.NewParser(),
		timezoneMap:    &tzMap,
		codeCountryMap: &codeCountryMap,
	}, nil
}
