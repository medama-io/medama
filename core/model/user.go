package model

type User struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`

	Settings    *UserSettings `json:"settings" db:"settings"`
	DateCreated int64         `json:"date_created" db:"date_created"`
	DateUpdated int64         `json:"date_updated" db:"date_updated"`
}

// NewUser returns a new instance of User with the given values.
func NewUser(id string, username string, password string, settings *UserSettings, dateCreated int64, dateUpdated int64) *User {
	return &User{
		ID:       id,
		Username: username,
		Password: password,

		Settings:    settings,
		DateCreated: dateCreated,
		DateUpdated: dateUpdated,
	}
}
