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

type pgTableRow struct {
	Schema  string  `db:"table_schema"`
	Table   string  `db:"table_name"`
	Comment *string `db:"comment"`
}

type pgColumnRow struct {
	Table     string  `db:"table_name"`
	Column    string  `db:"column_name"`
	Type      string  `db:"udt_name"`
	Schema    string  `db:"data_type"`
	Nullable  string  `db:"is_nullable"`
	Default   *string `db:"column_default"`
	Collation *string `db:"collation_name"`
}

type pgIndexRow struct {
	Table  string `db:"table_name"`
	Index  string `db:"index_name"`
	Unique bool   `db:"is_unique"`
	Column string `db:"column_name"`
}

// LoadTables 读取 information_schema.tables。Engine/Charset 无 PG 对应概念，留空。
func (Postgres) LoadTables(ctx context.Context, db *sqlx.DB, schema string) ([]*Table, error) {
	const query = "SELECT c.table_schema, c.table_name, " +
		"obj_description((quote_ident(c.table_schema) || '.' || quote_ident(c.table_name))::regclass) AS comment " +
		"FROM information_schema.tables c " +
		"WHERE c.table_schema = COALESCE(NULLIF($1, ''), current_schema()) " +
		"AND c.table_type = 'BASE TABLE' " +
		"ORDER BY c.table_name"

	var rows []pgTableRow
	if err := db.SelectContext(ctx, &rows, query, schema); err != nil {
		return nil, fmt.Errorf("query postgres tables: %w", err)
	}

	tables := make([]*Table, 0, len(rows))
	for _, row := range rows {
		tables = append(tables, &Table{
			Name:    row.Table,
			Comment: deref(row.Comment),
			Columns: make(map[string]*Column),
			Indexes: make(map[string]*Index),
		})
	}
	return tables, nil
}

// LoadColumns 读取 information_schema.columns，类型取 udt_name 走统一分类。
func (Postgres) LoadColumns(ctx context.Context, db *sqlx.DB, schema string) (map[string]map[string]*Column, error) {
	const query = "SELECT table_name, column_name, udt_name, data_type, " +
		"is_nullable, column_default, collation_name " +
		"FROM information_schema.columns " +
		"WHERE table_schema = COALESCE(NULLIF($1, ''), current_schema()) " +
		"ORDER BY table_name, ordinal_position"

	var rows []pgColumnRow
	if err := db.SelectContext(ctx, &rows, query, schema); err != nil {
		return nil, fmt.Errorf("query postgres columns: %w", err)
	}

	columns := make(map[string]map[string]*Column, len(rows))
	for _, row := range rows {
		col := &Column{
			Name:     row.Column,
			SQLType:  SQLType{Name: normalizePostgresType(row.Schema, row.Type)},
			Nullable: row.Nullable == "YES",
			Default:  deref(row.Default),
		}
		if columns[row.Table] == nil {
			columns[row.Table] = make(map[string]*Column)
		}
		columns[row.Table][col.Name] = col
	}
	return columns, nil
}

// LoadIndexes 走 pg_catalog，按索引内列序保留顺序，忽略主键。
func (Postgres) LoadIndexes(ctx context.Context, db *sqlx.DB, schema string) (map[string]map[string]*Index, error) {
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

	indexes := make(map[string]map[string]*Index)
	for _, row := range rows {
		if indexes[row.Table] == nil {
			indexes[row.Table] = make(map[string]*Index)
		}
		idx, ok := indexes[row.Table][row.Index]
		if !ok {
			idx = &Index{Name: row.Index, Unique: row.Unique}
		}
		idx.Cols = append(idx.Cols, row.Column)
		indexes[row.Table][idx.Name] = idx
	}
	return indexes, nil
}

// LoadSchema 组合装配。PG 无 Engine/Charset 概念，留空。
func (p Postgres) LoadSchema(ctx context.Context, db *sqlx.DB, schema string) ([]*Table, error) {
	tables, err := p.LoadTables(ctx, db, schema)
	if err != nil {
		return nil, err
	}
	columns, err := p.LoadColumns(ctx, db, schema)
	if err != nil {
		return nil, err
	}
	indexes, err := p.LoadIndexes(ctx, db, schema)
	if err != nil {
		return nil, err
	}
	assemble(tables, columns, indexes)
	return tables, nil
}

// ListSchemas 只覆盖当前连接的库：PG 的 schema 位于单个 database 内，
// 跨库需另开连接，此处不跨。
func (Postgres) ListSchemas(ctx context.Context, db *sqlx.DB) ([]string, error) {
	const query = "SELECT schema_name FROM information_schema.schemata " +
		"WHERE schema_name <> 'information_schema' AND schema_name NOT LIKE 'pg\\_%' " +
		"ORDER BY schema_name"

	var names []string
	if err := db.SelectContext(ctx, &names, query); err != nil {
		return nil, fmt.Errorf("query postgres schemas: %w", err)
	}
	return names, nil
}

// LoadDatabases 返回当前库内的用户 schema。PG 无库级默认字符集/排序规则，留空。
func (Postgres) LoadDatabases(ctx context.Context, db *sqlx.DB) ([]*Database, error) {
	names, err := Postgres{}.ListSchemas(ctx, db)
	if err != nil {
		return nil, err
	}
	dbs := make([]*Database, 0, len(names))
	for _, name := range names {
		dbs = append(dbs, &Database{Name: name})
	}
	return dbs, nil
}

// normalizePostgresType maps a few information_schema data_type spellings onto
// the canonical names the shared sqlKinds table classifies; udt_name is upper
// cased otherwise.
func normalizePostgresType(dataType, udtName string) string {
	switch dataType {
	case "character varying":
		return "VARCHAR"
	case "character":
		return "CHAR"
	case "bytea":
		return "BYTEA"
	case "text":
		return "TEXT"
	case "ARRAY":
		// element type keeps the classification meaningful
		return strings.ToUpper(udtName)
	}
	return strings.ToUpper(udtName)
}
