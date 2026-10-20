package resolvers

import (
	"context"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/collector/env"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// Environments 系统运行环境状态信息（纯 Go 实现，无 CGO 依赖）
func (r *queryRootResolver) Environments(ctx context.Context) (envs *gqlapi.Environments, err error) {
	envs = &gqlapi.Environments{}
	envs.CPUStats = bridge.CPUStatsFromDomain(env.GetCPUStats())
	envs.MemStats = bridge.MemStatsFromDomain(env.GetMemStats())
	envs.HostInfos = bridge.HostInfosFromDomain(env.GetHostInfos())
	envs.ProcessStats = bridge.ProcessStatsFromDomain(env.GetProcessStats())
	return
}
