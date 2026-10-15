package agent

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/agent/dialect"
)

// LoadSchema 拉取目标库单库的完整元数据（表/列/索引）。
func (t *Target) LoadSchema(ctx context.Context) (*Schema, error) {
	tables, err := t.dialect.LoadSchema(ctx, t.db, "")
	if err != nil {
		return nil, err
	}
	return &Schema{Tables: tables}, nil
}

// LoadAllSchemas 拉取目标实例全部用户库的元数据。MySQL 走 SHOW DATABASES
// 后逐库加载；validate/ 的 "*" 语义由它承接。
func (t *Target) LoadAllSchemas(ctx context.Context) (map[string][]*dialect.Table, error) {
	schemas, err := t.listSchemas(ctx)
	if err != nil {
		return nil, fmt.Errorf("list schemas: %w", err)
	}

	all := make(map[string][]*dialect.Table, len(schemas))
	for _, schema := range schemas {
		tables, err := t.dialect.LoadSchema(ctx, t.db, schema)
		if err != nil {
			return nil, fmt.Errorf("load schema %s: %w", schema, err)
		}
		all[schema] = tables
	}
	return all, nil
}

// listSchemas 枚举实例内的用户库/模式名。
func (t *Target) listSchemas(ctx context.Context) ([]string, error) {
	return t.dialect.ListSchemas(ctx, t.db)
}
