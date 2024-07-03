package services

import (
	"context"

	"github.com/medama-io/medama/api"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func (h *Handler) GetSettingsResource(ctx context.Context, params api.GetSettingsResourceParams) (api.GetSettingsResourceRes, error) {
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

	return &api.SettingsResource{
		CPU: api.SettingsResourceCPU{
			Usage:   float32(cpuUsage),
			Cores:   cpuCores,
			Threads: cpuThreads,
		},
		Memory: api.SettingsResourceMemory{
			Used:  int64(vmStat.Used),
			Total: int64(vmStat.Total),
		},
		Disk: api.SettingsResourceDisk{
			Used:  int64(diskStat.Used),
			Total: int64(diskStat.Total),
		},
	}, nil
}
