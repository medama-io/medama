package model

type GetUser struct {
	ID          string `json:"id" db:"id"`
	Email       string `json:"email" db:"email"`
	Language    string `json:"language" db:"language"`
	DateCreated int64  `json:"date_created" db:"date_created"`
	DateUpdated int64  `json:"date_updated" db:"date_updated"`
}

type User struct {
	ID       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`

	Language    string `json:"language" db:"language"`
	DateCreated int64  `json:"date_created" db:"date_created"`
	DateUpdated int64  `json:"date_updated" db:"date_updated"`
}

// NewUser returns a new instance of User with the given values.
func NewUser(id string, email string, password string, language string, dateCreated int64, dateUpdated int64) *User {
	return &User{
		ID:       id,
		Email:    email,
		Password: password,

		Language:    language,
		DateCreated: dateCreated,
		DateUpdated: dateUpdated,
	}
}
