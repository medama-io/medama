package main

import (
	"math/rand/v2"

	"github.com/medama-io/medama/model"
)

const (
	// Up to 10000ms (10s) of duration can be randomly generated.
	DURATION_RANGE = 10000
)

func generatePageViewHits(r *rand.Rand, count int, hostname string) []*model.PageViewHit {
	paths := []string{"/", "/about", "/contact"}
	booleanValues := []bool{true, false}
	referrers := []string{"", "medama.io", "google.com"}
	countryCodes := []string{"GB", "US", "JP"}
	languages := []string{"en", "jp"}
	browserNames := []model.BrowserName{model.ChromeBrowser, model.FirefoxBrowser, model.SafariBrowser}
	oses := []model.OSName{model.WindowsOS, model.LinuxOS, model.IOS}
	deviceTypes := []model.DeviceType{model.DesktopDevice, model.MobileDevice, model.TabletDevice}
	utmSources := []string{"", "bing", "twitter"}
	utmMediums := []string{"", "cpc", "organic"}
	utmCampaigns := []string{"", "summer", "winter"}

	pageViewHits := make([]*model.PageViewHit, count)

	for i := range count {
		pageViewHits[i] = &model.PageViewHit{
			Hostname:     hostname,
			Pathname:     paths[r.IntN(len(paths))],
			IsUniqueUser: booleanValues[r.IntN(len(booleanValues))],
			IsUniquePage: booleanValues[r.IntN(len(booleanValues))],
			Referrer:     referrers[r.IntN(len(referrers))],
			CountryCode:  countryCodes[r.IntN(len(countryCodes))],
			Language:     languages[r.IntN(len(languages))],
			BrowserName:  browserNames[r.IntN(len(browserNames))],
			OS:           oses[r.IntN(len(oses))],
			DeviceType:   deviceTypes[r.IntN(len(deviceTypes))],
			UTMSource:    utmSources[r.IntN(len(utmSources))],
			UTMMedium:    utmMediums[r.IntN(len(utmMediums))],
			UTMCampaign:  utmCampaigns[r.IntN(len(utmCampaigns))],
		}
	}

	return pageViewHits
}

func generatePageViewDurations(r *rand.Rand, hits []*model.PageViewHit, count int) []*model.PageViewDuration {
	durations := make([]*model.PageViewDuration, count)

	for i := range count {
		durations[i] = &model.PageViewDuration{
			PageViewHit: *hits[i],
			DurationMs:  r.IntN(DURATION_RANGE),
		}
	}

	return durations
}
