package model

type RequestKey string

const (
	// RequestKeyBody is the key used to store the request in the context.
	RequestKeyBody RequestKey = "request"
)

type EventHit struct {
	// BatchID - Used to group together multiple properties of the same event.
	BatchID string `db:"batch_id"`
	// Group - The group name of the event, typically the hostname.
	Group string `db:"group_name"`
	// Name - The name of the event.
	Name string `db:"name"`
	// Value - The value of the event.
	Value string `db:"value"`
}

type PageViewHit struct {
	// Beacon ID - Used to determine if multiple event types are
	// associated with a single page view.
	BID string `db:"bid"`

	// Hostname - The hostname of the page view.
	Hostname string `db:"hostname"`
	// Pathname - The pathname of the associated URL linked to the page view.
	Pathname string `db:"pathname"`

	// IsUniqueUser - Whether or not the page view is from a unique user.
	IsUniqueUser bool `db:"is_unique_user"`
	// IsUniquePage - Whether or not the user has visited the page before.
	IsUniquePage bool `db:"is_unique_page"`
	// ReferrerHost - The referrer hostname of the page view.
	ReferrerHost string `db:"referrer_host"`
	// ReferrerGroup - The referrer group of the page view. e.g. Google
	ReferrerGroup string `db:"referrer_group"`
	// Country - The country name associated with the user's timezone.
	Country string `db:"country"`
	// LanguageBase - The base language associated with the user's browser.
	LanguageBase string `db:"language_base"`
	// LanguageDialect - The dialect of the user's language. e.g. British English.
	LanguageDialect string `db:"language_dialect"`

	// BrowserName - The name of the user's browser.
	BrowserName string `db:"ua_browser"`
	// OS - The operating system the user is using.
	OS string `db:"ua_os"`
	// DeviceType - The type of device the user is using.
	DeviceType string `db:"ua_device_type"`

	// UTMSource - The UTM source of the page view.
	UTMSource string `db:"utm_source"`
	// UTMMedium - The UTM medium of the page view.
	UTMMedium string `db:"utm_medium"`
	// UTMCampaign - The UTM campaign of the page view.
	UTMCampaign string `db:"utm_campaign"`
}

type PageViewDuration struct {
	// Beacon ID - Used to determine if multiple event types are
	// associated with a single page view.
	BID string `db:"bid"`
	// DurationMs - How long the user has been on the page in milliseconds.
	DurationMs int `db:"duration_ms"`
}
