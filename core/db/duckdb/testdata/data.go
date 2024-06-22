package main

import (
	"math/rand/v2"
	"strconv"

	"github.com/medama-io/medama/model"
)

func generatePageViewHits(r *rand.Rand, count int, hostname string) []*model.PageViewHit {
	paths := []string{"/", "/about", "/contact"}
	booleanValues := []bool{true, false}
	referrers := []string{"", "medama.io", "google.com"}
	countries := []string{"United Kingdom", "United States", "Japan"}
	languagesBase := []string{"English", "Japanese"}
	languagesDialects := []string{"British English", "American English"}
	browserNames := []string{"Chrome", "Firefox", "Safari"}
	oses := []string{"Windows", "Linux", "iOS"}
	deviceTypes := []string{"Desktop", "Mobile", "Tablet"}
	utmSources := []string{"", "bing", "twitter"}
	utmMediums := []string{"", "cpc", "organic"}
	utmCampaigns := []string{"", "summer", "winter"}

	pageViewHits := make([]*model.PageViewHit, count)

	for i := range count {
		languageBase := languagesBase[r.IntN(len(languagesBase))]
		var languageDialect string
		if languageBase == "English" {
			languageDialect = languagesDialects[r.IntN(len(languagesDialects))]
		} else {
			languageDialect = "Japanese"
		}

		pageViewHits[i] = &model.PageViewHit{
			Hostname:        hostname,
			BID:             strconv.Itoa(i),
			Pathname:        paths[r.IntN(len(paths))],
			IsUniqueUser:    booleanValues[r.IntN(len(booleanValues))],
			IsUniquePage:    booleanValues[r.IntN(len(booleanValues))],
			ReferrerHost:    referrers[r.IntN(len(referrers))],
			Country:         countries[r.IntN(len(countries))],
			LanguageBase:    languageBase,
			LanguageDialect: languageDialect,
			BrowserName:     browserNames[r.IntN(len(browserNames))],
			OS:              oses[r.IntN(len(oses))],
			DeviceType:      deviceTypes[r.IntN(len(deviceTypes))],
			UTMSource:       utmSources[r.IntN(len(utmSources))],
			UTMMedium:       utmMediums[r.IntN(len(utmMediums))],
			UTMCampaign:     utmCampaigns[r.IntN(len(utmCampaigns))],
		}
	}

	return pageViewHits
}

func generatePageViewDurations(r *rand.Rand, hits []*model.PageViewHit, count int) []*model.PageViewDuration {
	durations := make([]*model.PageViewDuration, count)

	for i := range count {
		durations[i] = &model.PageViewDuration{
			BID:        hits[i].BID,
			DurationMs: r.IntN(DURATION_RANGE),
		}
	}

	return durations
}
