package agent

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/agent/meta"
)

// LoadObjects 拉取各库的用户对象清单（视图/函数/存储过程/触发器/事件）。未指定
// 库名时先枚举用户库再逐库读取，与 LoadTables 保持一致。
func (t *Target) LoadObjects(ctx context.Context, databases ...string) (map[string]*meta.Objects, error) {
	if len(databases) == 0 {
		dbs, err := t.dialect.Databases(ctx, t.db)
		if err != nil {
			return nil, fmt.Errorf("list databases: %w", err)
		}
		for _, db := range dbs {
			databases = append(databases, db.Name)
		}
	}

	all := make(map[string]*meta.Objects, len(databases))
	for _, database := range databases {
		objects, err := t.dialect.LoadObjects(ctx, t.db, database)
		if err != nil {
			return nil, fmt.Errorf("load objects %s: %w", database, err)
		}
		all[database] = objects
	}
	return all, nil
}
