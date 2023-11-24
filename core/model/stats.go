package model

type StatsSummary struct {
	Uniques   int `db:"uniques"`
	Pageviews int `db:"pageviews"`
	Bounces   int `db:"bounces"`
	Duration  int `db:"duration"`
	Active    int `db:"active"`
}

type StatsPagesSummary struct {
	Pathname          string  `db:"pathname"`
	Uniques           int     `db:"uniques"`
	UniquesPercentage float32 `db:"uniques_percentage"`
}

type StatsPages struct {
	StatsPagesSummary
	Title     string `db:"title"`
	Pageviews int    `db:"pageviews"`
	Bounces   int    `db:"bounces"`
	Duration  int    `db:"duration"`
}
