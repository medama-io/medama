package services

import (
	"context"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
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
		Username: user.Username,
		Settings: api.UserSettings{
			Language:   api.NewOptUserSettingsLanguage(api.UserSettingsLanguage(user.Settings.Language)),
			ScriptType: api.NewOptUserSettingsScriptType(api.UserSettingsScriptType(user.Settings.ScriptType)),
		},
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}, nil
}

func (h *Handler) GetUserUsage(ctx context.Context, params api.GetUserUsageParams) (api.GetUserUsageRes, error) {
	// CPU statistics.
	cpuCores, err := cpu.CountsWithContext(ctx, false)
	if err != nil {
		return nil, err
	}
	cpuThreads, err := cpu.CountsWithContext(ctx, true)
	if err != nil {
		return nil, err
	}

	cpuUsageArr, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		return nil, err
	}
	// Get the average CPU usage.
	cpuUsage := 0.0
	for _, v := range cpuUsageArr {
		cpuUsage += v
	}
	cpuUsage /= float64(len(cpuUsageArr))

	// Memory statistics.
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}

	// Disk statistics.
	diskStat, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		return nil, err
	}

	// Get the current database settings.
	metadata, err := h.analyticsDB.GetDuckDBSettings(ctx)
	if err != nil {
		return nil, err
	}

	return &api.SettingsUsageGet{
		CPU: api.SettingsUsageGetCPU{
			Usage:   float32(cpuUsage),
			Cores:   cpuCores,
			Threads: cpuThreads,
		},
		Memory: api.SettingsUsageGetMemory{
			Used:  int64(vmStat.Used),
			Total: int64(vmStat.Total),
		},
		Disk: api.SettingsUsageGetDisk{
			Used:  int64(diskStat.Used),
			Total: int64(diskStat.Total),
		},
		Metadata: api.SettingsUsageGetMetadata{
			Threads:     api.NewOptInt(metadata.Threads),
			MemoryLimit: api.NewOptString(strings.ReplaceAll(metadata.MemoryLimit, " ", "")),
		},
	}, nil
}

func (h *Handler) PatchUser(ctx context.Context, req *api.UserPatch, params api.PatchUserParams) (api.PatchUserRes, error) {
	log := logger.Get()
	if h.auth.IsDemoMode {
		log.Debug().Msg("patch user rejected in demo mode")
		return ErrForbidden(model.ErrDemoMode), nil
	}

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

	if req.Username.IsSet() {
		username := req.Username.Value
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

	if req.Password.IsSet() {
		password := req.Password.Value
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

	// Settings
	if req.Settings.IsSet() {
		settings := user.Settings
		if req.Settings.Value.Language.IsSet() {
			settings.Language = string(req.Settings.Value.Language.Value)
		}
		if req.Settings.Value.ScriptType.IsSet() {
			settings.ScriptType = string(req.Settings.Value.ScriptType.Value)
		}

		err = h.db.UpdateSettings(ctx, user.ID, settings, dateUpdated)
		if err != nil {
			log.Error().Err(err).Msg("failed to update user settings")
			return nil, errors.Wrap(err, "services")
		}
	}

	return &api.UserGet{
		Username: user.Username,
		Settings: api.UserSettings{
			Language:   api.NewOptUserSettingsLanguage(api.UserSettingsLanguage(user.Settings.Language)),
			ScriptType: api.NewOptUserSettingsScriptType(api.UserSettingsScriptType(user.Settings.ScriptType)),
		},
		DateCreated: user.DateCreated,
		DateUpdated: user.DateUpdated,
	}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, params api.DeleteUserParams) (api.DeleteUserRes, error) {
	log := logger.Get()
	if h.auth.IsDemoMode {
		log.Debug().Msg("delete user rejected in demo mode")
		return ErrForbidden(model.ErrDemoMode), nil
	}

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
			Int64("date_created", user.DateCreated).
			Int64("date_updated", user.DateUpdated).
			Err(err).
			Msg("failed to delete user")
		return nil, errors.Wrap(err, "services")
	}

	return &api.DeleteUserNoContent{}, nil
}
