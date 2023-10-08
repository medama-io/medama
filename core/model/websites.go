package model

type Website struct {
	ID       string `json:"id" db:"id"`
	UserID   string `json:"user_id" db:"user_id"`
	Hostname string `json:"hostname" db:"hostname"`

	DateCreated int64 `json:"date_created" db:"date_created"`
	DateUpdated int64 `json:"date_updated" db:"date_updated"`
}
