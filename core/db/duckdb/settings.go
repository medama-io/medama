package duckdb

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-faster/errors"
	qb "github.com/medama-io/medama/db/duckdb/query"
	"github.com/medama-io/medama/model"
)

func (c *Client) GetDuckDBSettings(ctx context.Context) (*model.DuckDBSettings, error) {
	queryThreads := qb.New().
		Select("value AS threads").
		From("duckdb_settings()").
		Where("name = 'threads'")

	queryMemory := qb.New().
		Select("value AS memory_limit").
		From("duckdb_settings()").
		Where("name = 'memory_limit'")

	var usage model.DuckDBSettings

	// Get the number of threads.
	err := c.GetContext(ctx, &usage, queryThreads.Build())
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	// Get the memory limit.
	err = c.GetContext(ctx, &usage, queryMemory.Build())
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	return &usage, nil
}

func (c *Client) SetDuckDBSettings(ctx context.Context, settings *model.DuckDBSettings) error {
	// SET does not support binding parameters, so we need to sanitize the input manually
	if settings.Threads > 0 {
		exec := fmt.Sprintf("SET threads = %s", strconv.Itoa(settings.Threads))
		_, err := c.ExecContext(ctx, exec)
		if err != nil {
			return errors.Wrap(err, "db")
		}
	}
	if settings.MemoryLimit != "" {
		sanitized, err := sanitizeMemoryLimit(settings.MemoryLimit)
		if err != nil {
			return errors.Wrap(err, "db")
		}
		exec := fmt.Sprintf("SET memory_limit = '%s'", sanitized)
		_, err = c.ExecContext(ctx, exec)
		if err != nil {
			return errors.Wrap(err, "db")
		}
	}

	return nil
}

func sanitizeMemoryLimit(limit string) (string, error) {
	// Regex pattern explanation:
	// ^ - Start of string
	// (\d+) - One or more digits (captured in group 1)
	// (MB|GB|TB|MiB|GiB|TiB) - One of the valid units (captured in group 2)
	// $ - End of string
	pattern := `^(\d+(?:\.\d+)?)(MB|GB|TB|MiB|GiB|TiB)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(limit)
	if matches == nil {
		return "", errors.New("invalid memory limit format " + limit)
	}

	return limit, nil
}
