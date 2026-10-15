package agent

import (
	"context"
	"time"

	"github.com/mia0x75/halo/internal/agent/meta"
	"github.com/mia0x75/halo/internal/domain"
)

// Executor 目标库执行器抽象。首版由 LocalExecutor 直连目标库实现，
// 后续 RemoteExecutor 改为调用目标机上的 Agent。
type Executor interface {
	Execute(ctx context.Context, c *domain.Cluster, password, query string) (*ExecResult, error)
	// Query 执行只读查询并返回按列名索引的字符串结果集。
	Query(ctx context.Context, c *domain.Cluster, password, query string) ([]map[string]string, error)
	// LoadTables 拉取表级元数据，外层键为库名。不传库名时拉取实例全部用户库，
	// validate/ 的 "*" 语义由它承接。
	LoadTables(ctx context.Context, c *domain.Cluster, password string, databases ...string) (map[string][]*meta.Table, error)
	// Databases 列举目标实例的用户库及其默认字符集/排序规则。
	Databases(ctx context.Context, c *domain.Cluster, password string) ([]*meta.Database, error)
}

// ExecResult 单条语句在目标库上的执行结果。
type ExecResult struct {
	RowsAffected int64
	Duration     time.Duration
	Error        string
}
