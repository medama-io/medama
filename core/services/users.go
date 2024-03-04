package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
)

func (h *Handler) GetUser(ctx context.Context, params api.GetUserParams) (api.GetUserRes, error) {
	// Get user id from request context and check if user exists
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}
	attributes := []slog.Attr{
		slog.String("id", userId),
	}
	slog.LogAttrs(ctx, slog.LevelDebug, "getting user...", attributes...)

	user, err := h.db.GetUser(ctx, userId)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))

		if errors.Is(err, model.ErrUserNotFound) {
			slog.LogAttrs(ctx, slog.LevelDebug, "user not found", attributes...)
			return ErrNotFound(err), nil
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to get user", attributes...)
		return nil, errors.Wrap(err, "services")
	}

	return &api.UserGet{
		Username:    user.Username,
		Language:    api.UserGetLanguage(user.Language),
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}, nil
}

func (h *Handler) PatchUser(ctx context.Context, req *api.UserPatch, params api.PatchUserParams) (api.PatchUserRes, error) {
	// Get user id from request context and check if user exists
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}

	attributes := []slog.Attr{
		slog.String("id", userId),
	}
	slog.LogAttrs(ctx, slog.LevelDebug, "getting user...", attributes...)

	user, err := h.db.GetUser(ctx, userId)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))

		if errors.Is(err, model.ErrUserNotFound) {
			slog.LogAttrs(ctx, slog.LevelDebug, "user not found", attributes...)
			return ErrNotFound(err), nil
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to get user", attributes...)
		return nil, errors.Wrap(err, "services")
	}

	// Update values
	dateUpdated := time.Now().Unix()
	username := req.Username.Value
	if username != "" {
		user.Username = username
		err = h.db.UpdateUserUsername(ctx, user.ID, username, dateUpdated)
		if err != nil {
			if errors.Is(err, model.ErrUserExists) {
				slog.LogAttrs(ctx, slog.LevelDebug, "email to patch already exists", attributes...)
				return ErrConflict(err), nil
			}

			if errors.Is(err, model.ErrUserNotFound) {
				slog.LogAttrs(ctx, slog.LevelDebug, "user not found", attributes...)
				return ErrNotFound(err), nil
			}

			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to update user email", attributes...)
			return nil, errors.Wrap(err, "services")
		}
	}

	password := req.Password.Value
	if password != "" {
		pwdHash, err := h.auth.HashPassword(password)
		if err != nil {
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to hash password", attributes...)
			return nil, errors.Wrap(err, "services")
		}

		err = h.db.UpdateUserPassword(ctx, user.ID, pwdHash, dateUpdated)
		if err != nil {
			attributes = append(attributes, slog.String("error", err.Error()))
			slog.LogAttrs(ctx, slog.LevelError, "failed to update user password", attributes...)
			return nil, errors.Wrap(err, "services")
		}
	}

	return &api.UserGet{
		Username:    user.Username,
		Language:    api.UserGetLanguage(user.Language),
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, params api.DeleteUserParams) (api.DeleteUserRes, error) {
	// Get user id from request context and check if user exists
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}

	attributes := []slog.Attr{
		slog.String("id", userId),
	}
	slog.LogAttrs(ctx, slog.LevelDebug, "getting user...", attributes...)

	user, err := h.db.GetUser(ctx, userId)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))

		if errors.Is(err, model.ErrUserNotFound) {
			slog.LogAttrs(ctx, slog.LevelDebug, "user not found", attributes...)
			return ErrNotFound(err), nil
		}

		slog.LogAttrs(ctx, slog.LevelError, "failed to get user", attributes...)
		return nil, errors.Wrap(err, "services")
	}

	attributes = append(attributes,
		slog.String("email", user.Username),
		slog.String("language", user.Language),
		slog.Int64("date_created", user.DateCreated),
		slog.Int64("date_updated", user.DateUpdated),
	)
	slog.LogAttrs(ctx, slog.LevelDebug, "deleting user", attributes...)

	err = h.db.DeleteUser(ctx, user.ID)
	if err != nil {
		attributes = append(attributes, slog.String("error", err.Error()))
		slog.LogAttrs(ctx, slog.LevelError, "failed to delete user", attributes...)
		return nil, errors.Wrap(err, "services")
	}

	return &api.DeleteUserOK{}, nil
}
