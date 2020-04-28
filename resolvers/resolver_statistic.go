package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/gqlapi"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
)

// Statistics 根据分组获取统计信息
func (r queryRootResolver) Statistics(ctx context.Context, groups []string) (L []*models.Statistic, err error) {
	L = caches.StatisticsMap.Filter(func(elem *models.Statistic) bool {
		return tools.Contains(groups, elem.Group)
	})

	// TODO: 测试代码，发布时删除
	id := uuid.New().String()
	notifyTicketStatusChange("e70e78bb-9d08-405d-a0ed-266ec703de19", id, fmt.Sprintf("Ticket (uuid=%s) status changed.", id))

	return
}

// Environments 系统运行环境状态信息
func (r queryRootResolver) Environments(ctx context.Context) (env *gqlapi.Environments, err error) {
	env = &gqlapi.Environments{}
	cs := g.GlobalStat.CPUStats()
	ms := g.GlobalStat.MemStats()
	hi := g.GlobalStat.HostInfos()
	ps := g.GlobalStat.ProcessStats()
	env.CPUStats = cs
	env.HostInfos = hi
	env.ProcessStats = ps
	env.MemStats = ms
	return
}

type statisticResolver struct{ *Resolver }
