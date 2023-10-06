package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/medama-io/medama/model"
)

type Handler interface {
	CreateUser(user *model.User) error
	GetUser(id string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id string) error
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
