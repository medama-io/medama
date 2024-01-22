package model

type RequestKey string

const (
	// RequestKeyBody is the key used to store the request in the context.
	RequestKeyBody RequestKey = "request"
)

// BrowserName - The name of the user's browser.
type BrowserName string

// OSName - The name of the user's operating system.
type OSName string

// DeviceType - The type of device the user is using.
type DeviceType string

func NewDeviceType(desktop bool, mobile bool, tablet bool, tv bool) DeviceType {
	switch {
	case desktop:
		return DesktopDevice
	case mobile:
		return MobileDevice
	case tablet:
		return TabletDevice
	case tv:
		return TVDevice
	default:
		return UnknownDevice
	}
}

const (
	UnknownBrowser          BrowserName = "Unknown"
	ChromeBrowser           BrowserName = "Chrome"
	EdgeBrowser             BrowserName = "Edge"
	FirefoxBrowser          BrowserName = "Firefox"
	InternetExplorerBrowser BrowserName = "Internet Explorer"
	OperaBrowser            BrowserName = "Opera"
	OperaMiniBrowser        BrowserName = "Opera Mini"
	SafariBrowser           BrowserName = "Safari"
	VivaldiBrowser          BrowserName = "Vivaldi"
	SamsungBrowser          BrowserName = "Samsung Browser"
	NintendoBrowser         BrowserName = "Nintendo Browser"

	UnknownOS OSName = "Unknown"
	AndroidOS OSName = "Android"
	ChromeOS  OSName = "Chrome OS"
	IOS       OSName = "iOS"
	LinuxOS   OSName = "Linux"
	MacOS     OSName = "Mac OS"
	WindowsOS OSName = "Windows"

	UnknownDevice DeviceType = "Unknown"
	DesktopDevice DeviceType = "Desktop"
	MobileDevice  DeviceType = "Mobile"
	TabletDevice  DeviceType = "Tablet"
	TVDevice      DeviceType = "TV"
)

type PageView struct {
	// Beacon ID - Used to determine if multiple event types are
	// associated with a single page view.
	BID string `db:"bid"`
	// Hostname - The hostname of the page view.
	Hostname string `db:"hostname"`
	// Pathname - The pathname of the associated URL linked to the page view.
	Pathname string `db:"pathname"`

	// IsUnique - Whether or not the page view is unique.
	IsUnique bool `db:"is_unique"`
	// ReferrerHostname - The hostname of the referrer of the page view.
	ReferrerHostname string `db:"referrer_hostname"`
	// ReferrerPathname - The pathname of the referrer of the page view.
	ReferrerPathname string `db:"referrer_pathname"`
	// CountryCode - The country code associated with the user's timezone.
	CountryCode string `db:"country_code"`
	// Language - The language associated with the user's browser.
	Language string `db:"language"`

	// BrowserName - The name of the user's browser.
	BrowserName BrowserName `db:"ua_browser"`
	// OS - The operating system the user is using.
	OS OSName `db:"ua_os"`
	// DeviceType - The type of device the user is using.
	DeviceType DeviceType `db:"ua_device_type"`

	// UTMSource - The UTM source of the page view.
	UTMSource string `db:"utm_source"`
	// UTMMedium - The UTM medium of the page view.
	UTMMedium string `db:"utm_medium"`
	// UTMCampaign - The UTM campaign of the page view.
	UTMCampaign string `db:"utm_campaign"`
}

type PageViewUpdate struct {
	// Beacon ID - Used to determine if multiple event types are
	// associated with a single page view.
	BID string `json:"bid" db:"bid"`
	// DurationMs - How long the user has been on the page in milliseconds.
	DurationMs int `json:"duration_ms" db:"duration_ms"`
}
