package sqlite

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/medama-io/medama/model"
)

type Handler interface {
	// CreateUser adds a new user to the database.
	CreateUser(ctx context.Context, user *model.User) error
	// GetUser retrieves a user from the database by id.
	GetUser(ctx context.Context, id string) (*model.User, error)
	// GetUserByEmail retrieves a user from the database by email.
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	// UpdateUser updates a user in the database.
	UpdateUser(ctx context.Context, user *model.User) error
	// DeleteUser deletes a user from the database.
	DeleteUser(ctx context.Context, id string) error
}

type Client struct {
	*sqlx.DB
}

// Compile time check for Handler.
var _ Handler = (*Client)(nil)

// NewClient returns a new instance of Client with the given configuration.
func NewClient(host string) (*Client, error) {
	db, err := sqlx.Connect("sqlite3", host)
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}
