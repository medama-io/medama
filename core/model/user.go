package model

type User struct {
	ID          string `json:"id" db:"id"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	Language    string `json:"language" db:"language"`
	DateCreated int64  `json:"date_created" db:"date_created"`
	DateUpdated int64  `json:"date_updated" db:"date_updated"`
}
