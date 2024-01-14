package model

const (
	// DateFormat is the format used for date and time values (YYYY-MM-DD).
	DateFormat = "2006-01-02"
)

type StatsSummarySingle struct {
	Uniques   int `db:"uniques"`
	Pageviews int `db:"pageviews"`
	Bounces   int `db:"bounces"`
	Duration  int `db:"duration"`
	Active    int `db:"active"`
}

type StatsSummary struct {
	Current  StatsSummarySingle
	Previous StatsSummarySingle
}

type StatsPagesSummary struct {
	Pathname         string  `db:"pathname"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
}

type StatsPages struct {
	StatsPagesSummary
	Title     string `db:"title"`
	Pageviews int    `db:"pageviews"`
	Bounces   int    `db:"bounces"`
	Duration  int    `db:"duration"`
}

type StatsTimeSummary struct {
	Pathname           string  `db:"pathname"`
	Duration           int     `db:"duration"`
	DurationPercentage float32 `db:"duration_percentage"`
}

type StatsTime struct {
	StatsTimeSummary
	DurationUpperQuartile int    `db:"duration_upper_quartile"`
	DurationLowerQuartile int    `db:"duration_lower_quartile"`
	Title                 string `db:"title"`
	Pageviews             int    `db:"pageviews"`
	Bounces               int    `db:"bounces"`
	Uniques               int    `db:"uniques"`
}

type StatsReferrerSummary struct {
	ReferrerHostname string  `db:"referrer_hostname"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
}

type StatsReferrers struct {
	StatsReferrerSummary
	ReferrerPathname string `db:"referrer_pathname"`
	Bounces          int    `db:"bounces"`
	Duration         int    `db:"duration"`
}

type StatsUTMSources struct {
	Source           string  `db:"source"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
}

type StatsUTMMediums struct {
	Medium           string  `db:"medium"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
}

type StatsUTMCampaigns struct {
	Campaign         string  `db:"campaign"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
}

type StatsBrowserSummary struct {
	Browser          BrowserName `db:"browser"`
	Uniques          int         `db:"uniques"`
	UniquePercentage float32     `db:"unique_percentage"`
}

type StatsBrowsers struct {
	StatsBrowserSummary
	Version string `db:"version"`
}

type StatsOS struct {
	OS               OSName  `db:"os"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
}

type StatsDevices struct {
	Device           DeviceType `db:"device"`
	Uniques          int        `db:"uniques"`
	UniquePercentage float32    `db:"unique_percentage"`
}

type StatsScreens struct {
	Screen           string  `db:"screen"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
}

type StatsCountries struct {
	Country          string  `db:"country"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
}

type StatsLanguages struct {
	Language         string  `db:"language"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
}
