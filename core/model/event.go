package model

import (
	ua "github.com/medama-io/go-useragent"
)

type RequestKey string

const (
	// RequestKeyBody is the key used to store the request in the context.
	RequestKeyBody RequestKey = "request"
)

// BrowserName - The name of the user's browser, stored as an integer to save
// space in the database.
type BrowserName uint8

const (
	Chrome    BrowserName = 1
	Edge      BrowserName = 2
	Firefox   BrowserName = 3
	IE        BrowserName = 4
	Opera     BrowserName = 5
	OperaMini BrowserName = 6
	Safari    BrowserName = 7
	Vivaldi   BrowserName = 8
	Samsung   BrowserName = 9
	Nintendo  BrowserName = 10
)

// NewBrowserName converts the browser name to a BrowserName integer.
func NewBrowserName(browser string) BrowserName {
	switch browser {
	case ua.Chrome:
		return Chrome
	case ua.Edge:
		return Edge
	case ua.Firefox:
		return Firefox
	case ua.IE:
		return IE
	case ua.Opera:
		return Opera
	case ua.OperaMini:
		return OperaMini
	case ua.Safari:
		return Safari
	case ua.Vivaldi:
		return Vivaldi
	case ua.Samsung:
		return Samsung
	case ua.Nintendo:
		return Nintendo
	default:
		return 0
	}
}

// String converts the browser name to a string.
func (b BrowserName) String() string {
	switch b {
	case Chrome:
		return ua.Chrome
	case Edge:
		return ua.Edge
	case Firefox:
		return ua.Firefox
	case IE:
		return ua.IE
	case Opera:
		return ua.Opera
	case OperaMini:
		return ua.OperaMini
	case Safari:
		return ua.Safari
	case Vivaldi:
		return ua.Vivaldi
	case Samsung:
		return ua.Samsung
	case Nintendo:
		return ua.Nintendo
	default:
		return ""
	}
}

// OSName - The name of the user's operating system, stored as an integer to
// save space in the database.
type OSName uint8

const (
	Android  OSName = 1
	ChromeOS OSName = 2
	IOS      OSName = 3
	Linux    OSName = 4
	MacOS    OSName = 5
	Windows  OSName = 6
)

// NewOSName converts the OS name to a OSName integer.
func NewOSName(os string) OSName {
	switch os {
	case ua.Android:
		return Android
	case ua.ChromeOS:
		return ChromeOS
	case ua.IOS:
		return IOS
	case ua.Linux:
		return Linux
	case ua.MacOS:
		return MacOS
	case ua.Windows:
		return Windows
	default:
		return 0
	}
}

// String converts the OS name to a string.
func (o OSName) String() string {
	switch o {
	case Android:
		return ua.Android
	case ChromeOS:
		return ua.ChromeOS
	case IOS:
		return ua.IOS
	case Linux:
		return ua.Linux
	case MacOS:
		return ua.MacOS
	case Windows:
		return ua.Windows
	default:
		return ""
	}
}

// DeviceType - The type of device the user is using, stored as an integer to
// save space in the database.
type DeviceType uint8

const (
	IsDesktop DeviceType = 1
	IsMobile  DeviceType = 2
	IsTablet  DeviceType = 3
	IsTV      DeviceType = 4
)

// NewDeviceType converts the device type to a DeviceType integer.
func NewDeviceType(desktop bool, mobile bool, tablet bool, tv bool) DeviceType {
	switch {
	case desktop:
		return IsDesktop
	case mobile:
		return IsMobile
	case tablet:
		return IsTablet
	case tv:
		return IsTV
	default:
		return 0
	}
}

// String converts the device type to a string.
func (d DeviceType) String() string {
	switch d {
	case IsDesktop:
		return ua.Desktop
	case IsMobile:
		return ua.Mobile
	case IsTablet:
		return ua.Tablet
	case IsTV:
		return ua.TV
	default:
		return ""
	}
}

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
	// Referrer - The referrer of the page view.
	Referrer string `db:"referrer"`
	// Title - The page title of the page view.
	Title string `db:"title"`
	// CountryCode - The country code associated with the user's timezone.
	CountryCode string `db:"country_code"`
	// Language - The language associated with the user's browser.
	Language string `db:"language"`

	// BrowserName - The name of the user's browser.
	BrowserName BrowserName `db:"ua_browser"`
	// BrowserVersion - The version of the user's browser.
	BrowserVersion string `db:"ua_version"`
	// OS - The operating system the user is using.
	OS OSName `db:"ua_os"`
	// DeviceType - The type of device the user is using.
	DeviceType DeviceType `db:"ua_device_type"`
	// RawUserAgent - The user agent of the user's browser. Only included if the
	// user agent was unable to be parsed.
	RawUserAgent string `db:"ua_raw"`

	// ScreenWidth - The width of the user's screen.
	ScreenWidth uint16 `json:"screen_width" db:"screen_width"`
	// ScreenHeight - The height of the user's screen.
	ScreenHeight uint16 `json:"screen_height" db:"screen_height"`

	// UTMSource - The UTM source of the page view.
	UTMSource string `db:"utm_source"`
	// UTMMedium - The UTM medium of the page view.
	UTMMedium string `db:"utm_medium"`
	// UTMCampaign - The UTM campaign of the page view.
	UTMCampaign string `db:"utm_campaign"`

	// DateCreated - Creation timestamp in UNIX.
	DateCreated int64 `json:"date_created" db:"date_created"`
}

type PageViewUpdate struct {
	// Beacon ID - Used to determine if multiple event types are
	// associated with a single page view.
	BID string `json:"bid" db:"bid"`
	// DurationMs - How long the user has been on the page in milliseconds.
	DurationMs int `json:"duration_ms" db:"duration_ms"`
}
