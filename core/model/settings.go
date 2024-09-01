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
	Language string `db:"language" json:"language"`

	// Tracker
	ScriptType string `db:"script_type" json:"script_type"`
}

type WebsiteSettings struct{}

// NewSettings returns a new instance of Settings with default values.
func NewDefaultSettings() *UserSettings {
	return &UserSettings{
		Language:   "en",
		ScriptType: "default",
	}
}
