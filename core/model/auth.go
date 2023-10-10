package model

import "time"

type ContextKey string

const (
	// ContextKeyUserID is the key used to store the user ID in the context.
	ContextKeyUserID ContextKey = "userId"
	// SessionCookieName is the name of the session cookie.
	SessionCookieName = "_me_sess"
	// SessionDuration is the duration of a session.
	// TODO: Make this configurable.
	SessionDuration = 12 * time.Hour
)
