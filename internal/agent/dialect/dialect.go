package dialect

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/mia0x75/halo/internal/agent/meta"
)

// Dialect 目标库方言。仅承载元数据读取与标识符引用，语句执行由 Target 统一负责。
type Dialect interface {
	// LoadTables 返回单库完整的表级元数据，database 为空时取当前库。
	LoadTables(ctx context.Context, db *sqlx.DB, database string) ([]*meta.Table, error)
	// LoadObjects 返回单库用户对象（视图/函数/存储过程/触发器/事件）的名字集合，
	// 供存在性规则使用。database 为空时取当前库。
	LoadObjects(ctx context.Context, db *sqlx.DB, database string) (*meta.Objects, error)
	// Databases 返回实例内的用户库/模式及其默认字符集与排序规则。
	Databases(ctx context.Context, db *sqlx.DB) ([]*meta.Database, error)
	// QuoteIdent 按方言加标识符引用，规避保留字冲突。
	QuoteIdent(s string) string
}
