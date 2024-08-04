package model

type SettingsKey string

const (
	// SettingsKeyLanguage is the key for the user's language setting.
	SettingsKeyLanguage SettingsKey = "language"
	// SettingsKeyScriptType is the key for the user's script type setting.
	SettingsKeyScriptType SettingsKey = "script_type"
	// SettingsKeyThreads is the key for the user's threads setting.
	SettingsKeyThreads SettingsKey = "threads"
	// SettingsKeyMemoryLimit is the key for the user's memory limit setting.
	SettingsKeyMemoryLimit SettingsKey = "memory_limit"
)

type DuckDBSettings struct {
	// Usage
	Threads     int    `json:"threads" db:"threads"`
	MemoryLimit string `json:"memory_limit" db:"memory_limit"`
}

type UserSettings struct {
	// Account
	Language string `json:"language" db:"language"`

	// Tracker
	ScriptType string `json:"script_type" db:"script_type"`
}

type Settings struct {
	DuckDBSettings
	UserSettings
}
