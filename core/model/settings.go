package model

type SettingsKey string

const (
	// SettingsKeyLanguage is the key for the user's language setting.
	SettingsKeyLanguage SettingsKey = "language"
	// SettingsKeyScriptType is the key for the user's script type setting.
	SettingsKeyScriptType SettingsKey = "script_type"
)

type UserSettings struct {
	// Account
	Language string `json:"language" db:"language"`

	// Tracker
	ScriptType string `json:"script_type" db:"script_type"`
}

type WebsiteSettings struct{}

// NewSettings returns a new instance of Settings with default values.
func NewDefaultSettings() *UserSettings {
	return &UserSettings{
		Language:   "en",
		ScriptType: "default",
	}
}
