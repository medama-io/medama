package duckdb

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/medama-io/medama/model"
)

type Handler interface {
	// Views
	AddPageView(ctx context.Context, event *model.PageView) error
	UpdatePageView(ctx context.Context, event *model.PageViewUpdate) error
	// Stats
	GetWebsiteSummary(ctx context.Context, hostname string) (*model.StatsSummary, error)
}

type Client struct {
	*sqlx.DB
}

// Compile time check for Handler.
var _ Handler = (*Client)(nil)

// NewClient returns a new instance of Client with the given configuration.
func NewClient(host string) (*Client, error) {
	// Enable foreign key support in sqlite
	db, err := sqlx.Connect("duckdb", host)
	if err != nil {
		return nil, err
	}

	return &Client{
		DB: db,
	}, nil
}
