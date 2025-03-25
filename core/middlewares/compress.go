package middlewares

import (
	"fmt"
	"net/http"

	"github.com/CAFxX/httpcompression"
)

func Compress() func(next http.Handler) http.Handler {
	compress, err := httpcompression.DefaultAdapter()
	if err != nil {
		panic(fmt.Errorf("failed to create compression adapter: %w", err))
	}

	return compress
}
