package model

type StatsSummary struct {
	Uniques   int `db:"uniques"`
	Pageviews int `db:"pageviews"`
	Bounces   int `db:"bounces"`
	Duration  int `db:"duration"`
	Active    int `db:"active"`
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
	Title     string `db:"title"`
	Pageviews int    `db:"pageviews"`
	Bounces   int    `db:"bounces"`
	Uniques   int    `db:"uniques"`
}

type StatsReferrers struct {
	Referrer         string  `db:"referrer"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
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

type StatsBrowsersSummary struct {
	Browser          string  `db:"browser"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
}

type StatsBrowsers struct {
	StatsBrowsersSummary
	Version string `db:"version"`
}

type StatsOS struct {
	OS               string  `db:"os"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
}

type StatsDevices struct {
	Device           string  `db:"device"`
	Uniques          int     `db:"uniques"`
	UniquePercentage float32 `db:"unique_percentage"`
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
