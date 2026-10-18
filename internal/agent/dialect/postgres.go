package dialect

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/mia0x75/halo/internal/agent/meta"
)

// Postgres 目标库方言。
type Postgres struct{}

// QuoteIdent Postgres 用双引号，内部双引号需转义。
func (Postgres) QuoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

type postgresTableRow struct {
	Schema  string  `db:"table_schema"`
	Table   string  `db:"table_name"`
	Comment *string `db:"comment"`
}

type postgresColumnRow struct {
	Table     string  `db:"table_name"`
	Column    string  `db:"column_name"`
	Type      string  `db:"udt_name"`
	Schema    string  `db:"data_type"`
	Nullable  string  `db:"is_nullable"`
	Default   *string `db:"column_default"`
	Identity  string  `db:"is_identity"`
	Collation *string `db:"collation_name"`
}

type postgresIndexRow struct {
	Table  string `db:"table_name"`
	Index  string `db:"index_name"`
	Unique bool   `db:"is_unique"`
	Column string `db:"column_name"`
}

// LoadTables 组合表/列/索引，装配成完整元数据。PG 无 Engine/Charset 概念，留空。
func (p Postgres) LoadTables(ctx context.Context, db *sqlx.DB, database string) ([]*meta.Table, error) {
	tables, err := p.loadTableHeaders(ctx, db, database)
	if err != nil {
		return nil, err
	}
	columns, err := p.loadColumns(ctx, db, database)
	if err != nil {
		return nil, err
	}
	indexes, err := p.loadIndexes(ctx, db, database)
	if err != nil {
		return nil, err
	}
	assemble(tables, columns, indexes)
	return tables, nil
}

// loadTableHeaders 读取当前库内用户 schema 的 information_schema.tables。
// A PG connection is bound to a single database, so `database` only records
// which key the caller expects; tables come from every non-system schema the
// connection can see. Engine/Charset 无 PG 对应概念，留空。
func (Postgres) loadTableHeaders(ctx context.Context, db *sqlx.DB, database string) ([]*meta.Table, error) {
	const query = "SELECT c.table_schema, c.table_name, " +
		"obj_description((quote_ident(c.table_schema) || '.' || quote_ident(c.table_name))::regclass) AS comment " +
		"FROM information_schema.tables c " +
		"WHERE c.table_catalog = current_database() " +
		"AND c.table_schema NOT IN ('information_schema') " +
		"AND c.table_schema NOT LIKE 'pg\\_%' " +
		"AND c.table_type = 'BASE TABLE' " +
		"ORDER BY c.table_name"

	var rows []postgresTableRow
	if err := db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("query postgres tables: %w", err)
	}

	tables := make([]*meta.Table, 0, len(rows))
	for _, row := range rows {
		tables = append(tables, &meta.Table{
			Name:    row.Table,
			Comment: deref(row.Comment),
			Columns: make(map[string]*meta.Column),
			Indexes: make(map[string]*meta.Index),
		})
	}
	return tables, nil
}

// loadColumns 读取当前库内用户 schema 的 information_schema.columns，类型取
// udt_name 走统一分类。
func (Postgres) loadColumns(ctx context.Context, db *sqlx.DB, database string) (map[string]map[string]*meta.Column, error) {
	const query = "SELECT table_name, column_name, udt_name, data_type, " +
		"is_nullable, column_default, is_identity, collation_name " +
		"FROM information_schema.columns " +
		"WHERE table_catalog = current_database() " +
		"AND table_schema NOT IN ('information_schema') " +
		"AND table_schema NOT LIKE 'pg\\_%' " +
		"ORDER BY table_name, ordinal_position"

	var rows []postgresColumnRow
	if err := db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("query postgres columns: %w", err)
	}

	columns := make(map[string]map[string]*meta.Column, len(rows))
	for _, row := range rows {
		col := &meta.Column{
			Name:          row.Column,
			Type:          meta.ColumnType(normalizePostgresType(row.Schema, row.Type)),
			Nullable:      row.Nullable == "YES",
			Default:       deref(row.Default),
			AutoIncrement: row.Identity == "YES" || strings.HasPrefix(deref(row.Default), "nextval("),
		}
		if columns[row.Table] == nil {
			columns[row.Table] = make(map[string]*meta.Column)
		}
		columns[row.Table][col.Name] = col
	}
	return columns, nil
}

// loadIndexes 走 pg_catalog，按索引内列序保留顺序，忽略主键。
func (Postgres) loadIndexes(ctx context.Context, db *sqlx.DB, database string) (map[string]map[string]*meta.Index, error) {
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

	var rows []postgresIndexRow
	if err := db.SelectContext(ctx, &rows, query, database); err != nil {
		return nil, fmt.Errorf("query postgres indexes: %w", err)
	}

	indexes := make(map[string]map[string]*meta.Index)
	for _, row := range rows {
		if indexes[row.Table] == nil {
			indexes[row.Table] = make(map[string]*meta.Index)
		}
		idx, ok := indexes[row.Table][row.Index]
		if !ok {
			idx = &meta.Index{Name: row.Index, Unique: row.Unique}
		}
		idx.Cols = append(idx.Cols, row.Column)
		indexes[row.Table][idx.Name] = idx
	}
	return indexes, nil
}

// Databases 返回实例内的用户库名。A PG connection is bound to one database, so
// callers must not expect every name to be reachable over this connection; the
// list exists so ticket/query validation can confirm the target database name.
// PG 无库级默认字符集/排序规则，留空。
func (Postgres) Databases(ctx context.Context, db *sqlx.DB) ([]*meta.Database, error) {
	const query = "SELECT datname FROM pg_database WHERE datistemplate = false ORDER BY datname"

	var names []string
	if err := db.SelectContext(ctx, &names, query); err != nil {
		return nil, fmt.Errorf("query postgres databases: %w", err)
	}
	dbs := make([]*meta.Database, 0, len(names))
	for _, name := range names {
		dbs = append(dbs, &meta.Database{Name: name})
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
