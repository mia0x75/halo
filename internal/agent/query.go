package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/mia0x75/halo/internal/agent/dialect"
)

// Databases 列举目标库的用户库/模式及默认字符集。
func (t *Target) Databases(ctx context.Context) ([]*dialect.Database, error) {
	return t.dialect.LoadDatabases(ctx, t.db)
}

// Query 执行只读查询，结果按列名转成字符串，保持旧 QueryString 的形态。

func (t *Target) Query(ctx context.Context, query string) ([]map[string]string, error) {
	rows, err := t.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query target: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]map[string]string, 0)
	for rows.Next() {
		raw := make(map[string]interface{})
		if err := rows.MapScan(raw); err != nil {
			return nil, fmt.Errorf("scan target row: %w", err)
		}
		row := make(map[string]string, len(raw))
		for k, v := range raw {
			row[k] = stringify(v)
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate target rows: %w", err)
	}
	return out, nil
}

// stringify renders a scanned value the way the legacy driver did: bytes and
// strings verbatim, everything else via fmt.
func stringify(v interface{}) string {
	switch x := v.(type) {
	case nil:
		return ""
	case []byte:
		return string(x)
	case string:
		return x
	case time.Time:
		return x.Format("2006-01-02 15:04:05")
	default:
		return fmt.Sprint(x)
	}
}
