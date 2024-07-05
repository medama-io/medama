package db

import (
	"context"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

// AppClient is the interface that groups all database operations related to
// user or website management.
type AppClient interface {
	// Users
	// CreateUser adds a new user to the database.
	CreateUser(ctx context.Context, user *model.User) error
	// GetUser retrieves a user from the database by id.
	GetUser(ctx context.Context, id string) (*model.User, error)
	// GetUserByUsername retrieves a user from the database by username.
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	// UpdateUserUsername updates a user's username in the database.
	UpdateUserUsername(ctx context.Context, id string, username string, dateUpdated int64) error
	// UpdateUserPassword updates a user's password in the database.
	UpdateUserPassword(ctx context.Context, id string, password string, dateUpdated int64) error
	// DeleteUser deletes a user from the database.
	DeleteUser(ctx context.Context, id string) error

	// Websites
	// CreateWebsite adds a new website to the database.
	CreateWebsite(ctx context.Context, website *model.Website) error
	// ListWebsites retrieves a list of websites from the database.
	ListWebsites(ctx context.Context, userID string) ([]*model.Website, error)
	// ListAllHostnames returns all hostnames from the database.
	ListAllHostnames(ctx context.Context) ([]string, error)
	// UpdateWebsite updates a website in the database.
	UpdateWebsite(ctx context.Context, website *model.Website) error
	// GetWebsite retrieves a website from the database by id.
	GetWebsite(ctx context.Context, id string) (*model.Website, error)
	// DeleteWebsite deletes a website from the database.
	DeleteWebsite(ctx context.Context, id string) error
}

// AnalyticsClient is the interface that groups all database operations related
// to analytics and events.
type AnalyticsClient interface {
	// Settings
	GetSettingsUsage(ctx context.Context) (*model.GetSettingsUsage, error)
	PatchSettingsUsage(ctx context.Context, settings *model.GetSettingsUsage) error
	// Events
	AddPageView(ctx context.Context, event *model.PageViewHit) error
	UpdatePageView(ctx context.Context, event *model.PageViewDuration) error
	// Pages
	GetWebsitePages(ctx context.Context, filter *Filters) ([]*model.StatsPages, error)
	GetWebsitePagesSummary(ctx context.Context, filter *Filters) ([]*model.StatsPagesSummary, error)
	// Locales
	GetWebsiteCountries(ctx context.Context, filter *Filters) ([]*model.StatsCountries, error)
	GetWebsiteCountriesSummary(ctx context.Context, filter *Filters) ([]*model.StatsCountriesSummary, error)
	GetWebsiteLanguages(ctx context.Context, isLocale bool, filter *Filters) ([]*model.StatsLanguages, error)
	GetWebsiteLanguagesSummary(ctx context.Context, isLocale bool, filter *Filters) ([]*model.StatsLanguagesSummary, error)
	// Referrers
	GetWebsiteReferrers(ctx context.Context, filter *Filters) ([]*model.StatsReferrers, error)
	GetWebsiteReferrersSummary(ctx context.Context, filter *Filters) ([]*model.StatsReferrerSummary, error)
	// Summary
	GetWebsiteSummary(ctx context.Context, filter *Filters) (*model.StatsSummarySingle, error)
	GetWebsiteIntervals(ctx context.Context, filter *Filters, interval api.GetWebsiteIDSummaryInterval) ([]*model.StatsIntervals, error)
	GetWebsiteSummaryLast24Hours(ctx context.Context, hostname string) (*model.StatsSummaryLast24Hours, error)
	// Time
	GetWebsiteTime(ctx context.Context, filter *Filters) ([]*model.StatsTime, error)
	GetWebsiteTimeSummary(ctx context.Context, filter *Filters) ([]*model.StatsTimeSummary, error)
	// Types
	GetWebsiteBrowsers(ctx context.Context, filter *Filters) ([]*model.StatsBrowsers, error)
	GetWebsiteBrowsersSummary(ctx context.Context, filter *Filters) ([]*model.StatsBrowsersSummary, error)
	GetWebsiteOS(ctx context.Context, filter *Filters) ([]*model.StatsOS, error)
	GetWebsiteOSSummary(ctx context.Context, filter *Filters) ([]*model.StatsOSSummary, error)
	GetWebsiteDevices(ctx context.Context, filter *Filters) ([]*model.StatsDevices, error)
	GetWebsiteDevicesSummary(ctx context.Context, filter *Filters) ([]*model.StatsDevicesSummary, error)
	// UTM
	GetWebsiteUTMSources(ctx context.Context, filter *Filters) ([]*model.StatsUTMSources, error)
	GetWebsiteUTMSourcesSummary(ctx context.Context, filter *Filters) ([]*model.StatsUTMSourcesSummary, error)
	GetWebsiteUTMMediums(ctx context.Context, filter *Filters) ([]*model.StatsUTMMediums, error)
	GetWebsiteUTMMediumsSummary(ctx context.Context, filter *Filters) ([]*model.StatsUTMMediumsSummary, error)
	GetWebsiteUTMCampaigns(ctx context.Context, filter *Filters) ([]*model.StatsUTMCampaigns, error)
	GetWebsiteUTMCampaignsSummary(ctx context.Context, filter *Filters) ([]*model.StatsUTMCampaignsSummary, error)
}
