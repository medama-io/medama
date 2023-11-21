package model

type StatsSummary struct {
	Uniques   int `db:"uniques"`
	Pageviews int `db:"pageviews"`
	Bounces   int `db:"bounces"`
	Duration  int `db:"duration"`
}
