package dialect

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Postgres 目标库方言。
type Postgres struct{}

// QuoteIdent Postgres 用双引号，内部双引号需转义。
func (Postgres) QuoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

type pgColumnRow struct {
	Table    string  `db:"table_name"`
	Column   string  `db:"column_name"`
	Type     string  `db:"data_type"`
	Nullable string  `db:"is_nullable"`
	Default  *string `db:"column_default"`
}

type pgIndexRow struct {
	Table  string `db:"table_name"`
	Index  string `db:"index_name"`
	Unique bool   `db:"is_unique"`
	Column string `db:"column_name"`
}

// LoadColumns 读取 information_schema.columns。schema 为空时用 current_schema()。
func (Postgres) LoadColumns(ctx context.Context, db *sqlx.DB, schema string) (map[string]map[string]ColumnInfo, error) {
	const query = "SELECT table_name, column_name, data_type, is_nullable, column_default " +
		"FROM information_schema.columns " +
		"WHERE table_schema = COALESCE(NULLIF($1, ''), current_schema()) " +
		"ORDER BY table_name, ordinal_position"

	var rows []pgColumnRow
	if err := db.SelectContext(ctx, &rows, query, schema); err != nil {
		return nil, fmt.Errorf("query postgres columns: %w", err)
	}

	columns := make(map[string]map[string]ColumnInfo, len(rows))
	for _, row := range rows {
		info := ColumnInfo{
			Name:     row.Column,
			Type:     row.Type,
			Nullable: row.Nullable == "YES",
			Default:  row.Default,
		}
		if columns[row.Table] == nil {
			columns[row.Table] = make(map[string]ColumnInfo)
		}
		columns[row.Table][row.Column] = info
	}
	return columns, nil
}

// LoadIndexes 走 pg_catalog，按索引内列序保留顺序，忽略主键。
func (Postgres) LoadIndexes(ctx context.Context, db *sqlx.DB, schema string) (map[string]map[string]IndexInfo, error) {
	const query = "SELECT t.relname AS table_name, i.relname AS index_name, " +
		"ix.indisunique AS is_unique, a.attname AS column_name " +
		"FROM pg_index ix " +
		"JOIN pg_class t ON t.oid = ix.indrelid " +
		"JOIN pg_class i ON i.oid = ix.indexrelid " +
		"JOIN pg_namespace n ON n.oid = t.relnamespace " +
		"JOIN LATERAL generate_subscripts(ix.indkey, 1) AS s ON true " +
		"JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ix.indkey[s] " +
		"WHERE n.nspname = COALESCE(NULLIF($1, ''), current_schema()) " +
		"AND NOT ix.indisprimary " +
		"ORDER BY t.relname, i.relname, s"

	var rows []pgIndexRow
	if err := db.SelectContext(ctx, &rows, query, schema); err != nil {
		return nil, fmt.Errorf("query postgres indexes: %w", err)
	}

	indexes := make(map[string]map[string]IndexInfo)
	for _, row := range rows {
		if indexes[row.Table] == nil {
			indexes[row.Table] = make(map[string]IndexInfo)
		}
		info, ok := indexes[row.Table][row.Index]
		if !ok {
			info = IndexInfo{Name: row.Index, Unique: row.Unique}
		}
		info.Columns = append(info.Columns, row.Column)
		indexes[row.Table][row.Index] = info
	}
	return indexes, nil
}
