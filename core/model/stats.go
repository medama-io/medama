package model

import "time"

const (
	// DateFormat is the format used for date and time values (YYYY-MM-DD).
	DateFormat = time.RFC3339
)

type StatsSummarySingle struct {
	Visitors  int `db:"visitors"`
	Pageviews int `db:"pageviews"`
	Bounces   int `db:"bounces"`
	Duration  int `db:"duration"`
}

type StatsSummary struct {
	Current  StatsSummarySingle
	Previous StatsSummarySingle
}

type StatsSummaryLast24Hours struct {
	Visitors int `db:"visitors"`
}

type StatsIntervals struct {
	Interval  string `db:"interval"`
	Visitors  int    `db:"visitors"`
	Pageviews int    `db:"pageviews"`
}

type StatsPagesSummary struct {
	Pathname           string  `db:"pathname"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsPages struct {
	StatsPagesSummary
	Pageviews           int     `db:"pageviews"`
	PageviewsPercentage float32 `db:"pageviews_percentage"`
	Bounces             int     `db:"bounces"`
	Duration            int     `db:"duration"`
}

type StatsTimeSummary struct {
	Pathname           string  `db:"pathname"`
	Duration           int     `db:"duration"`
	DurationPercentage float32 `db:"duration_percentage"`
}

type StatsTime struct {
	StatsTimeSummary
	DurationUpperQuartile int `db:"duration_upper_quartile"`
	DurationLowerQuartile int `db:"duration_lower_quartile"`
	Pageviews             int `db:"pageviews"`
	Visitors              int `db:"visitors"`
	Bounces               int `db:"bounces"`
}

type StatsReferrerSummary struct {
	Referrer           string  `db:"referrer_host"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsReferrers struct {
	StatsReferrerSummary
	Bounces  int `db:"bounces"`
	Duration int `db:"duration"`
}

type StatsUTMSourcesSummary struct {
	Source             string  `db:"source"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsUTMSources struct {
	StatsUTMSourcesSummary
	Bounces  int `db:"bounces"`
	Duration int `db:"duration"`
}

type StatsUTMMediumsSummary struct {
	Medium             string  `db:"medium"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsUTMMediums struct {
	StatsUTMMediumsSummary
	Bounces  int `db:"bounces"`
	Duration int `db:"duration"`
}

type StatsUTMCampaignsSummary struct {
	Campaign           string  `db:"campaign"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsUTMCampaigns struct {
	StatsUTMCampaignsSummary
	Bounces  int `db:"bounces"`
	Duration int `db:"duration"`
}

type StatsBrowsersSummary struct {
	Browser            BrowserName `db:"browser"`
	Visitors           int         `db:"visitors"`
	VisitorsPercentage float32     `db:"visitors_percentage"`
}

type StatsBrowsers struct {
	StatsBrowsersSummary
	Bounces  int `db:"bounces"`
	Duration int `db:"duration"`
}

type StatsOSSummary struct {
	OS                 OSName  `db:"os"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsOS struct {
	StatsOSSummary
	Bounces  int `db:"bounces"`
	Duration int `db:"duration"`
}

type StatsDevicesSummary struct {
	Device             DeviceType `db:"device"`
	Visitors           int        `db:"visitors"`
	VisitorsPercentage float32    `db:"visitors_percentage"`
}

type StatsDevices struct {
	StatsDevicesSummary
	Bounces  int `db:"bounces"`
	Duration int `db:"duration"`
}

type StatsCountriesSummary struct {
	Country            string  `db:"country"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsCountries struct {
	StatsCountriesSummary
	Bounces  int `db:"bounces"`
	Duration int `db:"duration"`
}

type StatsLanguagesSummary struct {
	Language           string  `db:"language_base"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsLanguages struct {
	StatsLanguagesSummary
	Bounces  int `db:"bounces"`
	Duration int `db:"duration"`
}
