package model

import "go.jetpack.io/typeid"

type User struct {
	ID          typeid.TypeID `json:"id" db:"id"`
	Email       string        `json:"email" db:"email"`
	Password    string        `json:"password" db:"password"`
	Language    string        `json:"language" db:"language"`
	DateCreated int64         `json:"date_created" db:"date_created"`
	DateUpdated int64         `json:"date_updated" db:"date_updated"`
}
