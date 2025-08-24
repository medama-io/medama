package model

type User struct {
	ID       string `db:"id"       json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`

	Settings    *UserSettings `db:"settings"     json:"settings"`
	DateCreated int64         `db:"date_created" json:"date_created"`
	DateUpdated int64         `db:"date_updated" json:"date_updated"`
}

// NewUser returns a new instance of User with the given values.
func NewUser(
	id string,
	username string,
	password string,
	settings *UserSettings,
	dateCreated int64,
	dateUpdated int64,
) *User {
	return &User{
		ID:       id,
		Username: username,
		Password: password,

		Settings:    settings,
		DateCreated: dateCreated,
		DateUpdated: dateUpdated,
	}
}
