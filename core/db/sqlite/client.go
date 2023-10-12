package sqlite

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/medama-io/medama/model"
)

type Handler interface {
	// Users
	// CreateUser adds a new user to the database.
	CreateUser(ctx context.Context, user *model.User) error
	// GetUser retrieves a user from the database by id.
	GetUser(ctx context.Context, id string) (*model.User, error)
	// GetUserByEmail retrieves a user from the database by email.
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	// GetUserCount retrieves the total number of users from the database.
	GetUserCount(ctx context.Context) (int64, error)
	// UpdateUserEmail updates a user's email in the database.
	UpdateUserEmail(ctx context.Context, id string, email string, dateUpdated int64) error
	// UpdateUserPassword updates a user's password in the database.
	UpdateUserPassword(ctx context.Context, id string, password string, dateUpdated int64) error
	// DeleteUser deletes a user from the database.
	DeleteUser(ctx context.Context, id string) error

	// Websites
	// CreateWebsite adds a new website to the database.
	CreateWebsite(ctx context.Context, website *model.Website) error
	// ListWebsites retrieves a list of websites from the database.
	ListWebsites(ctx context.Context, userID string) ([]*model.Website, error)
	// UpdateWebsite updates a website in the database.
	UpdateWebsite(ctx context.Context, website *model.Website) error
	// GetWebsite retrieves a website from the database by id.
	GetWebsite(ctx context.Context, id string) (*model.Website, error)
	// DeleteWebsite deletes a website from the database.
	DeleteWebsite(ctx context.Context, id string) error
}

type Client struct {
	*sqlx.DB
}

// Compile time check for Handler.
var _ Handler = (*Client)(nil)

// NewClient returns a new instance of Client with the given configuration.
func NewClient(host string) (*Client, error) {
	// Enable foreign key support in sqlite
	db, err := sqlx.Connect("sqlite3", fmt.Sprintf("file:%s", host))
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}
