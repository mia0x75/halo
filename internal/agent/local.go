package agent

import (
	"context"

	"github.com/mia0x75/halo/internal/agent/meta"
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

// LoadTables 直连目标库拉取表级元数据。
func (LocalExecutor) LoadTables(ctx context.Context, c *domain.Cluster, password string, databases ...string) (map[string][]*meta.Table, error) {
	target, err := openTarget(ctx, c, password)
	if err != nil {
		return nil, err
	}
	defer func() { _ = target.Close() }()

	return target.LoadTables(ctx, databases...)
}

// LoadObjects 直连目标库拉取用户对象清单。
func (LocalExecutor) LoadObjects(ctx context.Context, c *domain.Cluster, password string, databases ...string) (map[string]*meta.Objects, error) {
	target, err := openTarget(ctx, c, password)
	if err != nil {
		return nil, err
	}
	defer func() { _ = target.Close() }()

	return target.LoadObjects(ctx, databases...)
}

// Query 直连目标库执行只读查询。
func (LocalExecutor) Query(ctx context.Context, c *domain.Cluster, password, query string) ([]map[string]string, error) {
	target, err := openTarget(ctx, c, password)
	if err != nil {
		return nil, err
	}
	defer func() { _ = target.Close() }()

	return target.Query(ctx, query)
}

// Databases 直连目标库列举用户库。
func (LocalExecutor) Databases(ctx context.Context, c *domain.Cluster, password string) ([]*meta.Database, error) {
	target, err := openTarget(ctx, c, password)
	if err != nil {
		return nil, err
	}
	defer func() { _ = target.Close() }()

	return target.Databases(ctx)
}
