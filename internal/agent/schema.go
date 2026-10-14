package agent

import (
	"context"
	"fmt"
)

// LoadSchema 拉取目标库当前库的列与索引元数据。
func (t *Target) LoadSchema(ctx context.Context) (*Schema, error) {
	columns, err := t.dialect.LoadColumns(ctx, t.db, "")
	if err != nil {
		return nil, fmt.Errorf("load columns: %w", err)
	}
	indexes, err := t.dialect.LoadIndexes(ctx, t.db, "")
	if err != nil {
		return nil, fmt.Errorf("load indexes: %w", err)
	}

	tables := make(map[string]TableInfo, len(columns))
	for name := range columns {
		tables[name] = TableInfo{Name: name}
	}

	return &Schema{Tables: tables, Columns: columns, Indexes: indexes}, nil
}
