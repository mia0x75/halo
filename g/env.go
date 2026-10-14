package g

import (
	"context"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"

	"github.com/mia0x75/halo/models"
)

const envCollectTimeout = 2 * time.Second

// GetCPUStats 纯 Go 实现的 CPU 状态采集，替代 statgo（CGO + libstatgrab）。
func GetCPUStats() *models.CPUStats {
	stats := &models.CPUStats{}
	ctx, cancel := context.WithTimeout(context.Background(), envCollectTimeout)
	defer cancel()

	if percents, err := cpu.PercentWithContext(ctx, 500*time.Millisecond, false); err == nil && len(percents) == 1 {
		// 仅有一个总体使用率时，无法细分，退化为 User 字段，保证仪表盘有数据。
		stats.User = percents[0]
		stats.Idle = 100 - percents[0]
	}
	if times, err := cpu.TimesWithContext(ctx, false); err == nil && len(times) == 1 {
		t := times[0]
		stats.User = t.User
		stats.Kernel = t.System
		stats.Idle = t.Idle
		stats.IOWait = t.Iowait
		stats.Nice = t.Nice
		// gopsutil 无 Swap CPU 时间概念，保持与旧行为兼容置 0。
		stats.Swap = 0
	}
	if avg, err := load.AvgWithContext(ctx); err == nil && avg != nil {
		stats.LoadMin1 = avg.Load1
		stats.LoadMin5 = avg.Load5
		stats.LoadMin15 = avg.Load15
	}
	return stats
}

// GetMemStats 纯 Go 实现的内存状态采集。
func GetMemStats() *models.MemStats {
	stats := &models.MemStats{}
	ctx, cancel := context.WithTimeout(context.Background(), envCollectTimeout)
	defer cancel()

	if vm, err := mem.VirtualMemoryWithContext(ctx); err == nil && vm != nil {
		stats.Total = int(vm.Total)
		stats.Free = int(vm.Free)
		stats.Used = int(vm.Used)
		stats.Cache = int(vm.Cached)
	}
	if sm, err := mem.SwapMemoryWithContext(ctx); err == nil && sm != nil {
		stats.SwapTotal = int(sm.Total)
		stats.SwapUsed = int(sm.Used)
		stats.SwapFree = int(sm.Free)
	}
	return stats
}

// GetHostInfos 纯 Go 实现的主机信息采集。
func GetHostInfos() *models.HostInfos {
	info := &models.HostInfos{
		NCPUs:    runtime.NumCPU(),
		MaxCPUs:  runtime.NumCPU(),
		BitWidth: strconv.IntSize,
	}
	ctx, cancel := context.WithTimeout(context.Background(), envCollectTimeout)
	defer cancel()

	if h, err := host.InfoWithContext(ctx); err == nil && h != nil {
		info.OSName = h.OS
		info.OSRelease = h.KernelVersion
		info.OSVersion = h.PlatformVersion
		info.Platform = h.Platform
		info.HostName = h.Hostname
	}
	if n, err := cpu.CountsWithContext(ctx, true); err == nil && n > 0 {
		info.NCPUs = n
		info.MaxCPUs = n
	}
	return info
}

// GetProcessStats 纯 Go 实现的进程状态采集。失败时返回零值，不抛错，保证 GraphQL 查询可用。
func GetProcessStats() *models.ProcessStats {
	stats := &models.ProcessStats{}
	ctx, cancel := context.WithTimeout(context.Background(), envCollectTimeout)
	defer cancel()

	pids, err := process.PidsWithContext(ctx)
	if err != nil {
		return stats
	}
	stats.Total = len(pids)
	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			continue
		}
		statuses, err := p.StatusWithContext(ctx)
		if err != nil || len(statuses) == 0 {
			continue
		}
		status := statuses[0]
		switch status {
		case process.Running, "R":
			stats.Running++
		case process.Sleep, process.Idle, "S", "D", "I":
			stats.Sleeping++
		case process.Stop, "T", "t":
			stats.Stopped++
		case process.Zombie, "Z":
			stats.Zombie++
		default:
			// 未知状态计入 Total 即可
		}
	}
	return stats
}
