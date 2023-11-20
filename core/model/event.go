package model

import (
	ua "github.com/medama-io/go-useragent"
)

type RequestKey string

const (
	// RequestKeyBody is the key used to store the request in the context.
	RequestKeyBody RequestKey = "request"
)

type BrowserName string

const (
	Chrome    BrowserName = ua.Chrome
	Edge      BrowserName = ua.Edge
	Firefox   BrowserName = ua.Firefox
	IE        BrowserName = ua.IE
	Opera     BrowserName = ua.Opera
	OperaMini BrowserName = ua.OperaMini
	Safari    BrowserName = ua.Safari
	Vivaldi   BrowserName = ua.Vivaldi
	Samsung   BrowserName = ua.Samsung
	Nintendo  BrowserName = ua.Nintendo
)

type OSName string // Operating System

const (
	Android  OSName = ua.Android
	ChromeOS OSName = ua.ChromeOS
	IOS      OSName = ua.IOS
	Linux    OSName = ua.Linux
	MacOS    OSName = ua.MacOS
	Windows  OSName = ua.Windows
)

type DeviceType string

const (
	IsDesktop DeviceType = ua.Desktop
	IsMobile  DeviceType = ua.Mobile
	IsTablet  DeviceType = ua.Tablet
	IsTV      DeviceType = ua.TV
)

type PageView struct {
	// Beacon ID - Used to determine if multiple event types are
	// associated with a single page view.
	BID string `db:"bid"`
	// Hostname - The hostname of the page view.
	Hostname string `db:"hostname"`
	// Pathname - The pathname of the associated URL linked to the page view.
	Pathname string `db:"pathname"`

	// Referrer - The referrer of the page view.
	Referrer string `db:"referrer"`
	// Title - The page title of the page view.
	Title string `db:"title"`
	// Timezone - The timezone associated with the user's browser, allowing us
	// to determine the country of the user's location without compromising
	// their privacy with usage of IP addresses.
	Timezone string `db:"timezone"`
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
	ScreenWidth int `json:"screen_width" db:"screen_width"`
	// ScreenHeight - The height of the user's screen.
	ScreenHeight int `json:"screen_height" db:"screen_height"`

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
