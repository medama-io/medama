package services

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
)

func (h *Handler) GetUser(ctx context.Context, params api.GetUserParams) (api.GetUserRes, error) {
	// Get user id from request context and check if user exists
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}

	user, err := h.db.GetUser(ctx, userId)
	if err != nil {
		log := logger.Get().With().Err(err).Logger()

		if errors.Is(err, model.ErrUserNotFound) {
			log.Debug().Msg("user not found")
			return ErrNotFound(err), nil
		}

		log.Error().Msg("failed to get user")
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
	log := logger.Get()
	// Get user id from request context and check if user exists
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}

	user, err := h.db.GetUser(ctx, userId)
	if err != nil {
		log := log.With().Err(err).Logger()

		if errors.Is(err, model.ErrUserNotFound) {
			log.Debug().Msg("user not found")
			return ErrNotFound(err), nil
		}

		log.Error().Msg("failed to get user")
		return nil, errors.Wrap(err, "services")
	}

	// Update values
	dateUpdated := time.Now().Unix()
	username := req.Username.Value
	if username != "" {
		user.Username = username
		err = h.db.UpdateUserUsername(ctx, user.ID, username, dateUpdated)
		if err != nil {
			log := log.With().Str("username", username).Err(err).Logger()

			if errors.Is(err, model.ErrUserExists) {
				log.Debug().Msg("username already exists")
				return ErrConflict(err), nil
			}

			if errors.Is(err, model.ErrUserNotFound) {
				log.Debug().Msg("user not found")
				return ErrNotFound(err), nil
			}

			log.Error().Msg("failed to update user email")
			return nil, errors.Wrap(err, "services")
		}
	}

	password := req.Password.Value
	if password != "" {
		pwdHash, err := h.auth.HashPassword(password)
		if err != nil {
			log.Error().Err(err).Msg("failed to hash password")
			return nil, errors.Wrap(err, "services")
		}

		err = h.db.UpdateUserPassword(ctx, user.ID, pwdHash, dateUpdated)
		if err != nil {
			log.Error().Err(err).Msg("failed to update user password")
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
	log := logger.Get()
	// Get user id from request context and check if user exists
	userId, ok := ctx.Value(model.ContextKeyUserID).(string)
	if !ok {
		return ErrUnauthorised(model.ErrSessionNotFound), nil
	}

	user, err := h.db.GetUser(ctx, userId)
	if err != nil {
		log = log.With().Err(err).Logger()

		if errors.Is(err, model.ErrUserNotFound) {
			log.Debug().Msg("user not found")
			return ErrNotFound(err), nil
		}

		log.Error().Msg("failed to get user")
		return nil, errors.Wrap(err, "services")
	}

	err = h.db.DeleteUser(ctx, user.ID)
	if err != nil {
		log.Error().
			Str("username", user.Username).
			Str("language", user.Language).
			Int64("date_created", user.DateCreated).
			Int64("date_updated", user.DateUpdated).
			Err(err).
			Msg("failed to delete user")
		return nil, errors.Wrap(err, "services")
	}

	return &api.DeleteUserOK{}, nil
}
