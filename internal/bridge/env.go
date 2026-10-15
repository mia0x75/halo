package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

func CPUStatsFromDomain(d *domain.CPUStats) *models.CPUStats {
	if d == nil {
		return nil
	}
	return &models.CPUStats{
		User:      d.User,
		Kernel:    d.Kernel,
		Idle:      d.Idle,
		IOWait:    d.IOWait,
		Swap:      d.Swap,
		Nice:      d.Nice,
		LoadMin1:  d.LoadMin1,
		LoadMin5:  d.LoadMin5,
		LoadMin15: d.LoadMin15,
	}
}

func MemStatsFromDomain(d *domain.MemStats) *models.MemStats {
	if d == nil {
		return nil
	}
	return &models.MemStats{
		Total:     d.Total,
		Free:      d.Free,
		Used:      d.Used,
		Cache:     d.Cache,
		SwapTotal: d.SwapTotal,
		SwapUsed:  d.SwapUsed,
		SwapFree:  d.SwapFree,
	}
}

func HostInfosFromDomain(d *domain.HostInfos) *models.HostInfos {
	if d == nil {
		return nil
	}
	return &models.HostInfos{
		NCPUs:     d.NCPUs,
		MaxCPUs:   d.MaxCPUs,
		BitWidth:  d.BitWidth,
		OSName:    d.OSName,
		OSRelease: d.OSRelease,
		OSVersion: d.OSVersion,
		Platform:  d.Platform,
		HostName:  d.HostName,
	}
}

func ProcessStatsFromDomain(d *domain.ProcessStats) *models.ProcessStats {
	if d == nil {
		return nil
	}
	return &models.ProcessStats{
		Total:    d.Total,
		Running:  d.Running,
		Sleeping: d.Sleeping,
		Stopped:  d.Stopped,
		Zombie:   d.Zombie,
	}
}
