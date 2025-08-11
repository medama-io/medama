package model

type SettingsKey string

const (
	// SettingsKeyLanguage is the key for the user's language setting.
	SettingsKeyLanguage SettingsKey = "language"
	// SettingsKeyScriptType is the key for the user's script type setting.
	SettingsKeyScriptType SettingsKey = "script_type"
	// SettingsKeyBlockAbusiveIPs is the key for blocking abusive IPs setting.
	SettingsKeyBlockAbusiveIPs SettingsKey = "block_abusive_ips"
	// SettingsKeyBlockTorExitNodes is the key for blocking Tor exit nodes setting.
	SettingsKeyBlockTorExitNodes SettingsKey = "block_tor_exit_nodes"
	// SettingsKeyBlockedIPs is the key for the manually blocked IPs setting.
	SettingsKeyBlockedIPs SettingsKey = "blocked_ips"
)

type UserSettings struct {
	// Account
	Language string `db:"language" json:"language"`

	// Tracker
	ScriptType string `db:"script_type" json:"script_type"`

	// Spam Protection
	BlockAbusiveIPs   string `db:"block_abusive_ips"    json:"block_abusive_ips"`
	BlockTorExitNodes string `db:"block_tor_exit_nodes" json:"block_tor_exit_nodes"`
	BlockedIPs        string `db:"blocked_ips"          json:"blocked_ips"`
}

type WebsiteSettings struct{}

// NewSettings returns a new instance of Settings with default values.
func NewDefaultSettings() *UserSettings {
	return &UserSettings{
		Language:          "en",
		ScriptType:        "default",
		BlockAbusiveIPs:   "true",
		BlockTorExitNodes: "true",
		BlockedIPs:        "",
	}
}
