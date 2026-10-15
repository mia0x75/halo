package dialect

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Dialect 目标库方言。仅承载元数据读取与标识符引用，语句执行由 Target 统一负责。
type Dialect interface {
	// LoadTables 返回单库全部表级元数据。schema 为空时取当前库。
	LoadTables(ctx context.Context, db *sqlx.DB, schema string) ([]*Table, error)
	// LoadColumns 返回 表名 -> 列名 -> 列信息。schema 为空时取当前库。
	LoadColumns(ctx context.Context, db *sqlx.DB, schema string) (map[string]map[string]*Column, error)
	// LoadIndexes 返回 表名 -> 索引名 -> 索引信息。
	LoadIndexes(ctx context.Context, db *sqlx.DB, schema string) (map[string]map[string]*Index, error)
	// LoadSchema 组合以上三者并装配成完整 Table 列表。
	LoadSchema(ctx context.Context, db *sqlx.DB, schema string) ([]*Table, error)
	// ListSchemas 返回实例内的用户库/模式名，已过滤系统库。
	ListSchemas(ctx context.Context, db *sqlx.DB) ([]string, error)
	// LoadDatabases 返回用户库/模式及其默认字符集与排序规则。
	LoadDatabases(ctx context.Context, db *sqlx.DB) ([]*Database, error)
	// QuoteIdent 按方言加标识符引用，规避保留字冲突。
	QuoteIdent(s string) string
}
