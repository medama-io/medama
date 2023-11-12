package model

type RequestKey string

const (
	// RequestKeyBody is the key used to store the request in the context.
	RequestKeyBody RequestKey = "request"
)

type PageView struct {
	// Beacon ID - Used to determine if multiple event types are
	// associated with a single page view.
	BID string `json:"bid" db:"bid"`
	// Event Type - Represents the type of page view event: pagehide, unload,
	// load, visibilitychange, hidden, visible. This originates from the browser
	// API.
	EventType string `json:"event_type" db:"event_type"`
	// Hostname - The hostname of the page view.
	Hostname string `json:"hostname" db:"hostname"`
	// Pathname - The pathname of the associated URL linked to the page view.
	Pathname string `json:"pathname" db:"pathname"`

	// Referrer - The referrer of the page view.
	Referrer string `json:"referrer" db:"referrer"`
	// Title - The page title of the page view.
	Title string `json:"title" db:"title"`
	// Timezone - The timezone associated with the user's browser, allowing us
	// to determine the country of the user's location without compromising
	// their privacy with usage of IP addresses.
	Timezone string `json:"timezone" db:"timezone"`

	// ScreenWidth - The width of the user's screen.
	ScreenWidth int64 `json:"screen_width" db:"screen_width"`
	// ScreenHeight - The height of the user's screen.
	ScreenHeight int64 `json:"screen_height" db:"screen_height"`

	// DurationMs - How long the user has been on the page in milliseconds.
	DurationMs int64 `json:"duration_ms" db:"duration_ms"`
	// DateCreated - Creation timestamp in UNIX.
	DateCreated int64 `json:"date_created" db:"date_created"`
}
