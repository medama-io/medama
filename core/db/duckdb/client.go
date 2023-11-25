package duckdb

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/medama-io/medama/model"
)

type Handler interface {
	// Events
	AddPageView(ctx context.Context, event *model.PageView) error
	UpdatePageView(ctx context.Context, event *model.PageViewUpdate) error
	// Pages
	GetWebsitePages(ctx context.Context, hostname string) ([]*model.StatsPages, error)
	GetWebsitePagesSummary(ctx context.Context, hostname string) ([]*model.StatsPagesSummary, error)
	// Locales
	// GetWebsiteCountries(ctx context.Context, hostname string) ([]*model.StatsCountries, error)
	// GetWebsiteLanguages(ctx context.Context, hostname string) ([]*model.StatsLanguages, error)
	// Referrers
	// GetWebsiteReferrers(ctx context.Context, hostname string) ([]*model.StatsReferrers, error)
	// Summary
	GetWebsiteSummary(ctx context.Context, hostname string) (*model.StatsSummary, error)
	// Time
	GetWebsiteTime(ctx context.Context, hostname string) ([]*model.StatsTime, error)
	GetWebsiteTimeSummary(ctx context.Context, hostname string) ([]*model.StatsTimeSummary, error)
	// Types
	// GetWebsiteBrowsers(ctx context.Context, hostname string) ([]*model.StatsBrowsers, error)
	// GetWebsiteBrowsersSummary(ctx context.Context, hostname string) ([]*model.StatsBrowsersSummary, error)
	// GetWebsiteOS(ctx context.Context, hostname string) ([]*model.StatsOS, error)
	// GetWebsiteDevices(ctx context.Context, hostname string) ([]*model.StatsDevices, error)
	// UTM
	// GetWebsiteUTMSources(ctx context.Context, hostname string) ([]*model.StatsUTMSources, error)
	// GetWebsiteUTMMediums(ctx context.Context, hostname string) ([]*model.StatsUTMMediums, error)
	// GetWebsiteUTMCampaigns(ctx context.Context, hostname string) ([]*model.StatsUTMCampaigns, error)
}

type Client struct {
	*sqlx.DB
}

// Compile time check for Handler.
var _ Handler = (*Client)(nil)

// NewClient returns a new instance of Client with the given configuration.
func NewClient(host string) (*Client, error) {
	db, err := sqlx.Connect("duckdb", host)
	if err != nil {
		return nil, err
	}

	// Enable ICU extension
	_, err = db.Exec(`--sql
		INSTALL icu;`)
	if err != nil {
		return nil, err
	}

	// Load ICU extension
	_, err = db.Exec(`--sql
		LOAD icu;`)
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}
