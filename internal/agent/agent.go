package agent

import (
	"context"
	"time"

	"github.com/mia0x75/halo/internal/agent/dialect"
	"github.com/mia0x75/halo/internal/domain"
)

// Executor 目标库执行器抽象。首版由 LocalExecutor 直连目标库实现，
// 后续 RemoteExecutor 改为调用目标机上的 Agent。
type Executor interface {
	Execute(ctx context.Context, c *domain.Cluster, password, query string) (*ExecResult, error)
	// Query 执行只读查询并返回按列名索引的字符串结果集。
	Query(ctx context.Context, c *domain.Cluster, password, query string) ([]map[string]string, error)
	LoadSchema(ctx context.Context, c *domain.Cluster, password string) (*Schema, error)
	// LoadAllSchemas 拉取目标实例全部用户库的元数据，外层键为库名。
	LoadAllSchemas(ctx context.Context, c *domain.Cluster, password string) (map[string][]*dialect.Table, error)
	// Databases 列举目标实例的用户库及其默认字符集/排序规则。
	Databases(ctx context.Context, c *domain.Cluster, password string) ([]*dialect.Database, error)
}

// ExecResult 单条语句在目标库上的执行结果。
type ExecResult struct {
	RowsAffected int64
	Duration     time.Duration
	Error        string
}

// Schema 单库的元数据快照。
type Schema struct {
	Tables []*dialect.Table
}
