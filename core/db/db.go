package db

import (
	"context"

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
	// UpdateWebsite updates a website in the database.
	UpdateWebsite(ctx context.Context, website *model.Website) error
	// GetWebsite retrieves a website from the database by id.
	GetWebsite(ctx context.Context, id string) (*model.Website, error)
	// WebsiteExists checks if a website exists in the database.
	WebsiteExists(ctx context.Context, id string) (bool, error)
	// DeleteWebsite deletes a website from the database.
	DeleteWebsite(ctx context.Context, id string) error
}

// AnalyticsClient is the interface that groups all database operations related
// to analytics and events.
type AnalyticsClient interface {
	// Events
	AddPageView(ctx context.Context, event *model.PageView) error
	UpdatePageView(ctx context.Context, event *model.PageViewUpdate) error
	// Pages
	GetWebsitePages(ctx context.Context, filter Filter) ([]*model.StatsPages, error)
	GetWebsitePagesSummary(ctx context.Context, filter Filter) ([]*model.StatsPagesSummary, error)
	// Locales
	GetWebsiteCountries(ctx context.Context, filter Filter) ([]*model.StatsCountries, error)
	GetWebsiteLanguages(ctx context.Context, filter Filter) ([]*model.StatsLanguages, error)
	// Referrers
	GetWebsiteReferrers(ctx context.Context, filter Filter) ([]*model.StatsReferrers, error)
	GetWebsiteReferrersSummary(ctx context.Context, filter Filter) ([]*model.StatsReferrerSummary, error)
	// Summary
	GetWebsiteSummary(ctx context.Context, filter Filter) (*model.StatsSummarySingle, error)
	// Time
	GetWebsiteTime(ctx context.Context, filter Filter) ([]*model.StatsTime, error)
	GetWebsiteTimeSummary(ctx context.Context, filter Filter) ([]*model.StatsTimeSummary, error)
	// Types
	GetWebsiteBrowsers(ctx context.Context, filter Filter) ([]*model.StatsBrowsers, error)
	GetWebsiteBrowsersSummary(ctx context.Context, filter Filter) ([]*model.StatsBrowserSummary, error)
	GetWebsiteOS(ctx context.Context, filter Filter) ([]*model.StatsOS, error)
	GetWebsiteDevices(ctx context.Context, filter Filter) ([]*model.StatsDevices, error)
	// UTM
	GetWebsiteUTMSources(ctx context.Context, filter Filter) ([]*model.StatsUTMSources, error)
	GetWebsiteUTMMediums(ctx context.Context, filter Filter) ([]*model.StatsUTMMediums, error)
	GetWebsiteUTMCampaigns(ctx context.Context, filter Filter) ([]*model.StatsUTMCampaigns, error)
}
