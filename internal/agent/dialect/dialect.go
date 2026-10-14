package dialect

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Dialect 目标库方言。仅承载元数据读取与标识符引用，语句执行由 Target 统一负责。
type Dialect interface {
	// LoadColumns 返回 表名 -> 列名 -> 列信息。schema 为空时取当前库。
	LoadColumns(ctx context.Context, db *sqlx.DB, schema string) (map[string]map[string]ColumnInfo, error)
	// LoadIndexes 返回 表名 -> 索引名 -> 索引信息。
	LoadIndexes(ctx context.Context, db *sqlx.DB, schema string) (map[string]map[string]IndexInfo, error)
	// QuoteIdent 按方言加标识符引用，规避保留字冲突。
	QuoteIdent(s string) string
}

// ColumnInfo 列级元数据。
type ColumnInfo struct {
	Name     string
	Type     string
	Nullable bool
	Default  *string
	Key      string
	Comment  string
}

// IndexInfo 索引级元数据。
type IndexInfo struct {
	Name    string
	Unique  bool
	Columns []string
}
