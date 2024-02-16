package model

const (
	// DateFormat is the format used for date and time values (YYYY-MM-DD).
	DateFormat = "2006-01-02T15:04:05Z"
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
	Referrer           string  `db:"referrer"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsReferrers struct {
	StatsReferrerSummary
	Bounces  int `db:"bounces"`
	Duration int `db:"duration"`
}

type StatsUTMSources struct {
	Source             string  `db:"source"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsUTMMediums struct {
	Medium             string  `db:"medium"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsUTMCampaigns struct {
	Campaign           string  `db:"campaign"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}
type StatsBrowsers struct {
	Browser            BrowserName `db:"browser"`
	Visitors           int         `db:"visitors"`
	VisitorsPercentage float32     `db:"visitors_percentage"`
}

type StatsOS struct {
	OS                 OSName  `db:"os"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsDevices struct {
	Device             DeviceType `db:"device"`
	Visitors           int        `db:"visitors"`
	VisitorsPercentage float32    `db:"visitors_percentage"`
}

type StatsCountries struct {
	Country            string  `db:"country"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}

type StatsLanguages struct {
	Language           string  `db:"language"`
	Visitors           int     `db:"visitors"`
	VisitorsPercentage float32 `db:"visitors_percentage"`
}
