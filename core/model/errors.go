package model

import "errors"

var (
	// General
	// ErrInternalServerError is the default generic error.
	ErrInternalServerError = errors.New("internal server error")

	// Users
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")
)
