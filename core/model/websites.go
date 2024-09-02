package model

type Website struct {
	UserID   string `db:"user_id"  json:"user_id"`
	Hostname string `db:"hostname" json:"hostname"`

	DateCreated int64 `db:"date_created" json:"date_created"`
	DateUpdated int64 `db:"date_updated" json:"date_updated"`
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
