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
	LoadSchema(ctx context.Context, c *domain.Cluster, password string) (*Schema, error)
}

// ExecResult 单条语句在目标库上的执行结果。
type ExecResult struct {
	RowsAffected int64
	Duration     time.Duration
	Error        string
}

// Schema 目标库的元数据快照。外层键统一为表名。
type Schema struct {
	Tables  map[string]TableInfo
	Columns map[string]map[string]dialect.ColumnInfo
	Indexes map[string]map[string]dialect.IndexInfo
}

// TableInfo 表级元数据。
type TableInfo struct {
	Name    string
	Engine  string
	Comment string
	Charset string
	Collate string
}
