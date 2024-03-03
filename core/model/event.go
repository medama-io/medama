package model

type RequestKey string

// BrowserName - The name of the user's browser.
type BrowserName uint8

// OSName - The name of the user's operating system.
type OSName uint8

// DeviceType - The type of device the user is using.
type DeviceType uint8

/**
 * These values are synced with go-useragent: https://github.com/medama-io/go-useragent/blob/main/match.go
 */

const (
	// RequestKeyBody is the key used to store the request in the context.
	RequestKeyBody RequestKey = "request"

	unknown = "Unknown"
	na      = "N/A"

	UnknownBrowser          BrowserName = 0
	ChromeBrowser           BrowserName = 1
	EdgeBrowser             BrowserName = 2
	FirefoxBrowser          BrowserName = 3
	InternetExplorerBrowser BrowserName = 4
	OperaBrowser            BrowserName = 5
	OperaMiniBrowser        BrowserName = 6
	SafariBrowser           BrowserName = 7
	VivaldiBrowser          BrowserName = 8
	SamsungBrowser          BrowserName = 9
	NintendoBrowser         BrowserName = 10

	UnknownOS OSName = 0
	AndroidOS OSName = 1
	ChromeOS  OSName = 2
	IOS       OSName = 3
	LinuxOS   OSName = 4
	MacOS     OSName = 5
	WindowsOS OSName = 6

	UnknownDevice DeviceType = 0
	DesktopDevice DeviceType = 1
	MobileDevice  DeviceType = 2
	TabletDevice  DeviceType = 3
	TVDevice      DeviceType = 4
)

func NewBrowserName(name string) BrowserName {
	switch name {
	case "Chrome":
		return ChromeBrowser
	case "Edge":
		return EdgeBrowser
	case "Firefox":
		return FirefoxBrowser
	case "InternetExplorer":
		return InternetExplorerBrowser
	case "Opera":
		return OperaBrowser
	case "OperaMini":
		return OperaMiniBrowser
	case "Safari":
		return SafariBrowser
	case "Vivaldi":
		return VivaldiBrowser
	case "SamsungBrowser":
		return SamsungBrowser
	case "NintendoBrowser":
		return NintendoBrowser
	default:
		return UnknownBrowser
	}
}

func (b BrowserName) String() string {
	switch b {
	case ChromeBrowser:
		return "Chrome"
	case EdgeBrowser:
		return "Edge"
	case FirefoxBrowser:
		return "Firefox"
	case InternetExplorerBrowser:
		return "Internet Explorer"
	case OperaBrowser:
		return "Opera"
	case OperaMiniBrowser:
		return "Opera Mini"
	case SafariBrowser:
		return "Safari"
	case VivaldiBrowser:
		return "Vivaldi"
	case SamsungBrowser:
		return "Samsung Browser"
	case NintendoBrowser:
		return "Nintendo Browser"
	case UnknownBrowser:
		return unknown
	default:
		return na
	}
}

func NewOSName(name string) OSName {
	switch name {
	case "Android":
		return AndroidOS
	case "ChromeOS":
		return ChromeOS
	case "iOS":
		return IOS
	case "Linux":
		return LinuxOS
	case "MacOS":
		return MacOS
	case "Windows":
		return WindowsOS
	default:
		return UnknownOS
	}
}

func (o OSName) String() string {
	switch o {
	case AndroidOS:
		return "Android"
	case ChromeOS:
		return "ChromeOS"
	case IOS:
		return "iOS"
	case LinuxOS:
		return "Linux"
	case MacOS:
		return "MacOS"
	case WindowsOS:
		return "Windows"
	case UnknownOS:
		return unknown
	default:
		return na
	}
}

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

func NewDeviceTypeString(name string) DeviceType {
	switch name {
	case "Desktop":
		return DesktopDevice
	case "Mobile":
		return MobileDevice
	case "Tablet":
		return TabletDevice
	case "TV":
		return TVDevice
	default:
		return UnknownDevice
	}
}

func (d DeviceType) String() string {
	switch d {
	case DesktopDevice:
		return "Desktop"
	case MobileDevice:
		return "Mobile"
	case TabletDevice:
		return "Tablet"
	case TVDevice:
		return "TV"
	case UnknownDevice:
		return unknown
	default:
		return na
	}
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
	// Referrer - The referrer URL of the page view.
	Referrer string `db:"referrer"`
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

type PageViewDuration struct {
	// Beacon ID - Used to determine if multiple event types are
	// associated with a single page view.
	BID string `db:"bid"`
	// DurationMs - How long the user has been on the page in milliseconds.
	DurationMs int `db:"duration_ms"`
}
