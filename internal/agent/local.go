package agent

import (
	"context"

	"github.com/mia0x75/halo/internal/agent/dialect"
	"github.com/mia0x75/halo/internal/domain"
)

// LocalExecutor 直连目标库执行。首版默认实现。
type LocalExecutor struct{}

// Execute 直连目标库执行单条语句。
func (LocalExecutor) Execute(ctx context.Context, c *domain.Cluster, password, query string) (*ExecResult, error) {
	target, err := openTarget(ctx, c, password)
	if err != nil {
		return nil, err
	}
	defer func() { _ = target.Close() }()

	return target.Execute(ctx, query)
}

// LoadSchema 直连目标库拉取元数据。
func (LocalExecutor) LoadSchema(ctx context.Context, c *domain.Cluster, password string) (*Schema, error) {
	target, err := openTarget(ctx, c, password)
	if err != nil {
		return nil, err
	}
	defer func() { _ = target.Close() }()

	return target.LoadSchema(ctx)
}

// LoadAllSchemas 直连目标库拉取实例内全部用户库的元数据。
func (LocalExecutor) LoadAllSchemas(ctx context.Context, c *domain.Cluster, password string) (map[string][]*dialect.Table, error) {
	target, err := openTarget(ctx, c, password)
	if err != nil {
		return nil, err
	}
	defer func() { _ = target.Close() }()

	return target.LoadAllSchemas(ctx)
}
