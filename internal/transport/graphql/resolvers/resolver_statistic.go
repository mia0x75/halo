package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/collector/env"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/platform/cache"
	"github.com/mia0x75/halo/internal/tools"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// Statistics 根据分组获取统计信息
func (r *queryRootResolver) Statistics(ctx context.Context, groups []string) (L []*models.Statistic, err error) {
	L = cache.StatisticsMap.Filter(func(elem *models.Statistic) bool {
		return tools.Contains(groups, elem.Group)
	})

	// TODO: 测试代码，发布时删除
	id := uuid.New().String()
	notifyTicketStatusChange("e70e78bb-9d08-405d-a0ed-266ec703de19", id, fmt.Sprintf("Ticket (uuid=%s) status changed.", id))

	return
}

// Environments 系统运行环境状态信息（纯 Go 实现，无 CGO 依赖）
func (r *queryRootResolver) Environments(ctx context.Context) (envs *gqlapi.Environments, err error) {
	envs = &gqlapi.Environments{}
	envs.CPUStats = bridge.CPUStatsFromDomain(env.GetCPUStats())
	envs.MemStats = bridge.MemStatsFromDomain(env.GetMemStats())
	envs.HostInfos = bridge.HostInfosFromDomain(env.GetHostInfos())
	envs.ProcessStats = bridge.ProcessStatsFromDomain(env.GetProcessStats())
	return
}
