package services

import (
	"context"
	"strings"

	"github.com/medama-io/medama/api"
	"github.com/medama-io/medama/model"
	"github.com/medama-io/medama/util/logger"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func (h *Handler) GetSettingsUsage(ctx context.Context, params api.GetSettingsUsageParams) (api.GetSettingsUsageRes, error) {
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
	metadata, err := h.analyticsDB.GetSettingsUsage(ctx)
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

func (h *Handler) PatchSettingsUsage(ctx context.Context, req *api.SettingsUsagePatch, params api.PatchSettingsUsageParams) (api.PatchSettingsUsageRes, error) {
	log := logger.Get()
	if h.auth.IsDemoMode {
		log.Debug().Msg("patch settings rejected in demo mode")
		return ErrForbidden(model.ErrDemoMode), nil
	}

	// Update the database settings.
	settings := &model.GetSettingsUsage{
		Threads:     req.GetThreads().Value,
		MemoryLimit: req.GetMemoryLimit().Value,
	}

	log = log.With().Int("threads", settings.Threads).Str("memory_limit", settings.MemoryLimit).Logger()

	err := h.analyticsDB.PatchSettingsUsage(ctx, settings)
	if err != nil {
		log.Error().Err(err).Msg("failed to update the usage settings")
		return nil, err
	}

	log.Warn().Msg("updated the usage settings")
	return &api.PatchSettingsUsageCreated{}, nil
}
