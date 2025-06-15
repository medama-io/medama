package referrer

import (
	"fmt"

	"github.com/medama-io/go-referrer-parser"
)

type Parser struct {
	parser *referrer.Parser
}

// NewParser creates a new referrer Parser instance.
func NewParser() (*Parser, error) {
	referrerParser, err := referrer.NewParser()
	if err != nil {
		return nil, fmt.Errorf("failed to create referrer parser: %w", err)
	}

	return &Parser{
		parser: referrerParser,
	}, nil
}
