package model

import "errors"

var (
	// General
	// ErrInternalServerError is the default generic error.
	ErrInternalServerError = errors.New("internal server error")

	// Users
	// ErrUserExists is returned when a user already exists.
	ErrUserExists = errors.New("user already exists")
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")
)
