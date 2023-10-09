package model

import "time"

const (
	// SessionCookieName is the name of the session cookie.
	SessionCookieName = "_me_sess"
	// SessionDuration is the duration of a session.
	// TODO: Make this configurable.
	SessionDuration = 12 * time.Hour
)
