package services

import (
	"context"
	"log/slog"

	"github.com/medama-io/medama/api"
)

func (h *Handler) NewError(ctx context.Context, err error) *api.InternalServerErrorStatusCode {
	slog.Error("Error W", err)
	return nil
}
