package duckdb

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/model"
)

func (c *Client) GetSettingsUsage(ctx context.Context) (*model.GetSettingsUsage, error) {
	queryThreads := `--sql
		SELECT
			value AS threads
		FROM
			duckdb_settings()
		WHERE name = 'threads';`

	queryMemory := `--sql
		SELECT
			value AS memory_limit
		FROM
			duckdb_settings()
		WHERE name = 'memory_limit';`

	var usage model.GetSettingsUsage

	// Get the number of threads.
	err := c.GetContext(ctx, &usage, queryThreads)
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	// Get the memory limit.
	err = c.GetContext(ctx, &usage, queryMemory)
	if err != nil {
		return nil, errors.Wrap(err, "db")
	}

	return &usage, nil
}

func (c *Client) PatchSettingsUsage(ctx context.Context, settings *model.GetSettingsUsage) error {
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
			return errors.Wrap(err, "invalid memory limit")
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
		return "", errors.New("invalid memory limit format")
	}

	return limit, nil
}
