package model

import "time"

const (
	// DateFormat is the format used for date and time values (YYYY-MM-DD).
	DateFormat = time.RFC3339
)

type StatsSummarySingle struct {
	Visitors   int     `db:"visitors"`
	Pageviews  int     `db:"pageviews"`
	BounceRate float32 `db:"bounce_rate"`
	Duration   int     `db:"duration"`
}

type StatsSummary struct {
	Current  StatsSummarySingle
	Previous StatsSummarySingle
}

type StatsSummaryLast24Hours struct {
	Visitors int `db:"visitors"`
}

type StatsIntervals struct {
	Interval   string  `db:"interval"`
	Visitors   int     `db:"visitors"`
	Pageviews  int     `db:"pageviews"`
	BounceRate float32 `db:"bounce_rate"`
	Duration   int     `db:"duration"`
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
	BounceRate          float32 `db:"bounce_rate"`
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
}

type StatsReferrerSummary struct {
	Referrer           string  `db:"referrer"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsReferrers struct {
	StatsReferrerSummary
	BounceRate float32 `db:"bounce_rate"`
	Duration   int     `db:"duration"`
}

type StatsUTMSourcesSummary struct {
	Source             string  `db:"source"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsUTMSources struct {
	StatsUTMSourcesSummary
	BounceRate float32 `db:"bounce_rate"`
	Duration   int     `db:"duration"`
}

type StatsUTMMediumsSummary struct {
	Medium             string  `db:"medium"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsUTMMediums struct {
	StatsUTMMediumsSummary
	BounceRate float32 `db:"bounce_rate"`
	Duration   int     `db:"duration"`
}

type StatsUTMCampaignsSummary struct {
	Campaign           string  `db:"campaign"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsUTMCampaigns struct {
	StatsUTMCampaignsSummary
	BounceRate float32 `db:"bounce_rate"`
	Duration   int     `db:"duration"`
}

type StatsBrowsersSummary struct {
	Browser            string  `db:"browser"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsBrowsers struct {
	StatsBrowsersSummary
	BounceRate float32 `db:"bounce_rate"`
	Duration   int     `db:"duration"`
}

type StatsOSSummary struct {
	OS                 string  `db:"os"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsOS struct {
	StatsOSSummary
	BounceRate float32 `db:"bounce_rate"`
	Duration   int     `db:"duration"`
}

type StatsDevicesSummary struct {
	Device             string  `db:"device"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsDevices struct {
	StatsDevicesSummary
	BounceRate float32 `db:"bounce_rate"`
	Duration   int     `db:"duration"`
}

type StatsCountriesSummary struct {
	Country            string  `db:"country"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsCountries struct {
	StatsCountriesSummary
	BounceRate float32 `db:"bounce_rate"`
	Duration   int     `db:"duration"`
}

type StatsLanguagesSummary struct {
	Language           string  `db:"language"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsLanguages struct {
	StatsLanguagesSummary
	BounceRate float32 `db:"bounce_rate"`
	Duration   int     `db:"duration"`
}

type StatsCustomProperties struct {
	Name     string `db:"name"`
	Value    string `db:"value"`
	Events   int    `db:"events"`
	Visitors int    `db:"visitors"`
}
