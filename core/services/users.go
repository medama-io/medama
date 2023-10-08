package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
	"go.jetpack.io/typeid"
)

func (h *Handler) GetUsersUserId(ctx context.Context, params api.GetUsersUserIdParams) (api.GetUsersUserIdRes, error) {
	user, err := h.db.GetUser(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return ErrNotFound(err), nil
		}

		slog.Log(ctx, slog.LevelError, "get user error", "error", err)
		return nil, err
	}

	return &api.UserGet{
		ID:          user.ID,
		Email:       user.Email,
		Language:    user.Language,
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}, nil
}

func (h *Handler) PatchUsersUserId(ctx context.Context, req api.OptUserPatch, params api.PatchUsersUserIdParams) (api.PatchUsersUserIdRes, error) {
	user, err := h.db.GetUser(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return ErrNotFound(err), nil
		}

		slog.Log(ctx, slog.LevelError, "get user error", "error", err)
		return nil, model.ErrInternalServerError
	}

	// Update values
	dateUpdated := time.Now().Unix()
	email := req.Value.Email.Value
	if email != "" {
		user.Email = email
		err = h.db.UpdateUserEmail(ctx, user.ID, email, dateUpdated)
		if err != nil {
			return nil, err
		}
	}

	password := req.Value.Password.Value
	if password != "" {
		pwdHash, err := h.auth.HashPassword(password)
		if err != nil {
			return nil, err
		}

		err = h.db.UpdateUserPassword(ctx, user.ID, pwdHash, dateUpdated)
		if err != nil {
			return nil, err
		}
	}

	return &api.UserGet{
		ID:          user.ID,
		Email:       user.Email,
		Language:    user.Language,
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}, nil
}

func (h *Handler) PostUser(ctx context.Context, req api.OptUserCreate) (api.PostUserRes, error) {
	// UUIDv7 id generation
	typeid, err := typeid.New("user")
	if err != nil {
		return nil, err
	}
	id := typeid.String()

	// Validate language as an accepted enum
	err = req.Value.Language.Value.Validate()
	if err != nil {
		return nil, err
	}
	language := string(req.Value.Language.Value)

	// Hash password
	pwdHash, err := h.auth.HashPassword(req.Value.Password)
	if err != nil {
		return nil, err
	}

	dateCreated := time.Now().Unix()
	dateUpdated := dateCreated

	user := &model.User{
		ID:          id,
		Email:       req.Value.Email,
		Password:    pwdHash,
		Language:    language,
		DateCreated: dateCreated,
		DateUpdated: dateUpdated,
	}

	attributes := []slog.Attr{
		slog.String("id", id),
		slog.String("email", req.Value.Email),
		slog.String("language", language),
		slog.Int64("date_created", dateCreated),
		slog.Int64("date_updated", dateUpdated),
	}
	slog.LogAttrs(ctx, slog.LevelDebug, "creating user", attributes...)

	err = h.db.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &api.UserGet{
		ID:          id,
		Email:       req.Value.Email,
		Language:    language,
		DateCreated: dateCreated,
		DateUpdated: dateUpdated,
	}, nil
}
