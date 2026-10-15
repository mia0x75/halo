package agent

import (
	"context"

	"github.com/mia0x75/halo/internal/agent/dialect"
	"github.com/mia0x75/halo/internal/domain"
)

// RemoteExecutor 后续实现：中心服务调用目标机上的 Agent，
// 目标库凭据不下发到目标机之外的链路。首版不启用。
type RemoteExecutor struct {
	Endpoint string
}

// Execute 暂未实现。
func (RemoteExecutor) Execute(ctx context.Context, c *domain.Cluster, password, query string) (*ExecResult, error) {
	return nil, ErrNotImplemented
}

// Query 暂未实现。
func (RemoteExecutor) Query(ctx context.Context, c *domain.Cluster, password, query string) ([]map[string]string, error) {
	return nil, ErrNotImplemented
}

// LoadSchema 暂未实现。
func (RemoteExecutor) LoadSchema(ctx context.Context, c *domain.Cluster, password string) (*Schema, error) {
	return nil, ErrNotImplemented
}

// LoadAllSchemas 暂未实现。
func (RemoteExecutor) LoadAllSchemas(ctx context.Context, c *domain.Cluster, password string) (map[string][]*dialect.Table, error) {
	return nil, ErrNotImplemented
}

// Databases 暂未实现。
func (RemoteExecutor) Databases(ctx context.Context, c *domain.Cluster, password string) ([]*dialect.Database, error) {
	return nil, ErrNotImplemented
}
