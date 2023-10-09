package model

import "errors"

var (
	// General
	// ErrUnauthorised is returned when a user is not authorised.
	ErrUnauthorised = errors.New("unauthorised session")
	// ErrInternalServerError is the default generic error.
	ErrInternalServerError = errors.New("internal server error")

	// Authentication
	// ErrInvalidSession is returned when a session is invalid.
	ErrInvalidSession = errors.New("invalid session")
	// ErrSessionNotFound is returned when a session is not found.
	ErrSessionNotFound = errors.New("session not found")

	// Users
	// ErrUserExists is returned when a user already exists.
	ErrUserExists = errors.New("user already exists")
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")

	// Websites
	// ErrWebsiteExists is returned when a website already exists.
	ErrWebsiteExists = errors.New("website already exists")
	// ErrWebsiteNotFound is returned when a website is not found.
	ErrWebsiteNotFound = errors.New("website not found")
)
