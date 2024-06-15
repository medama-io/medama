package model

type Website struct {
	UserID   string `json:"user_id" db:"user_id"`
	Hostname string `json:"hostname" db:"hostname"`

	DateCreated int64 `json:"date_created" db:"date_created"`
	DateUpdated int64 `json:"date_updated" db:"date_updated"`
}

// NewWebsite returns a new instance of Website with the given values.
func NewWebsite(userID string, hostname string, dateCreated int64, dateUpdated int64) *Website {
	return &Website{
		UserID:   userID,
		Hostname: hostname,

		DateCreated: dateCreated,
		DateUpdated: dateUpdated,
	}
}
