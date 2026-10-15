package agent

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/agent/meta"
)

// LoadTables 拉取表级元数据。未指定库名时先枚举用户库再逐库读取，全程复用
// 同一连接——PG 的 schema 位于单个 database 内，逐库另开连接会失效。
func (t *Target) LoadTables(ctx context.Context, databases ...string) (map[string][]*meta.Table, error) {
	if len(databases) == 0 {
		dbs, err := t.dialect.Databases(ctx, t.db)
		if err != nil {
			return nil, fmt.Errorf("list databases: %w", err)
		}
		for _, db := range dbs {
			databases = append(databases, db.Name)
		}
	}

	all := make(map[string][]*meta.Table, len(databases))
	for _, database := range databases {
		tables, err := t.dialect.LoadTables(ctx, t.db, database)
		if err != nil {
			return nil, fmt.Errorf("load tables %s: %w", database, err)
		}
		all[database] = tables
	}
	return all, nil
}
