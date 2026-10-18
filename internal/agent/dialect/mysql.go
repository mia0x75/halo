package dialect

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/mia0x75/halo/internal/agent/meta"
)

// MySQL 目标库方言。
type MySQL struct{}

// QuoteIdent MySQL 用反引号，`type` 等保留字必须引用。
func (MySQL) QuoteIdent(s string) string {
	return "`" + s + "`"
}

type mysqlTableRow struct {
	Schema    string  `db:"TABLE_SCHEMA"`
	Table     string  `db:"TABLE_NAME"`
	Engine    *string `db:"ENGINE"`
	Collation *string `db:"TABLE_COLLATION"`
	Comment   *string `db:"TABLE_COMMENT"`
}

type mysqlColumnRow struct {
	Table    string  `db:"TABLE_NAME"`
	Column   string  `db:"COLUMN_NAME"`
	Type     string  `db:"COLUMN_TYPE"`
	Nullable string  `db:"IS_NULLABLE"`
	Default  *string `db:"COLUMN_DEFAULT"`
	Extra    *string `db:"EXTRA"`
	Comment  *string `db:"COLUMN_COMMENT"`
}

type mysqlIndexRow struct {
	Table     string `db:"TABLE_NAME"`
	Index     string `db:"INDEX_NAME"`
	NonUnique int    `db:"NON_UNIQUE"`
	Column    string `db:"COLUMN_NAME"`
}

// LoadTables 组合 TABLES/COLUMNS/STATISTICS 装配完整元数据。
func (m MySQL) LoadTables(ctx context.Context, db *sqlx.DB, database string) ([]*meta.Table, error) {
	tables, err := m.loadTableHeaders(ctx, db, database)
	if err != nil {
		return nil, err
	}
	columns, err := m.loadColumns(ctx, db, database)
	if err != nil {
		return nil, err
	}
	indexes, err := m.loadIndexes(ctx, db, database)
	if err != nil {
		return nil, err
	}
	assemble(tables, columns, indexes)
	return tables, nil
}

// loadTableHeaders 读取 information_schema.TABLES。database 为空时用 DATABASE()。
func (MySQL) loadTableHeaders(ctx context.Context, db *sqlx.DB, database string) ([]*meta.Table, error) {
	const query = "SELECT TABLE_SCHEMA, TABLE_NAME, ENGINE, TABLE_COLLATION, TABLE_COMMENT " +
		"FROM information_schema.TABLES " +
		"WHERE TABLE_SCHEMA = COALESCE(NULLIF(?, ''), DATABASE()) " +
		"AND TABLE_TYPE = 'BASE TABLE' " +
		"ORDER BY TABLE_NAME"

	var rows []mysqlTableRow
	if err := db.SelectContext(ctx, &rows, query, database); err != nil {
		return nil, fmt.Errorf("query mysql tables: %w", err)
	}

	tables := make([]*meta.Table, 0, len(rows))
	for _, row := range rows {
		t := &meta.Table{
			Name:    row.Table,
			Comment: deref(row.Comment),
			Engine:  deref(row.Engine),
			Collate: deref(row.Collation),
			Columns: make(map[string]*meta.Column),
			Indexes: make(map[string]*meta.Index),
		}
		t.Charset = strings.SplitN(t.Collate, "_", 2)[0]
		tables = append(tables, t)
	}
	return tables, nil
}

// loadColumns 读取 information_schema.COLUMNS。database 为空时用 DATABASE()。
func (MySQL) loadColumns(ctx context.Context, db *sqlx.DB, database string) (map[string]map[string]*meta.Column, error) {
	const query = "SELECT TABLE_NAME, COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, " +
		"COLUMN_DEFAULT, EXTRA, COLUMN_COMMENT " +
		"FROM information_schema.COLUMNS " +
		"WHERE TABLE_SCHEMA = COALESCE(NULLIF(?, ''), DATABASE()) " +
		"ORDER BY TABLE_NAME, ORDINAL_POSITION"

	var rows []mysqlColumnRow
	if err := db.SelectContext(ctx, &rows, query, database); err != nil {
		return nil, fmt.Errorf("query mysql columns: %w", err)
	}

	columns := make(map[string]map[string]*meta.Column, len(rows))
	for _, row := range rows {
		col := &meta.Column{
			Name:          row.Column,
			Type:          meta.ColumnType(normalizeMySQLType(row.Type)),
			Nullable:      row.Nullable == "YES",
			Default:       deref(row.Default),
			Comment:       deref(row.Comment),
			AutoIncrement: strings.Contains(deref(row.Extra), "auto_increment"),
		}
		if columns[row.Table] == nil {
			columns[row.Table] = make(map[string]*meta.Column)
		}
		columns[row.Table][col.Name] = col
	}
	return columns, nil
}

// loadIndexes 读取 information_schema.STATISTICS，忽略主键。
func (MySQL) loadIndexes(ctx context.Context, db *sqlx.DB, database string) (map[string]map[string]*meta.Index, error) {
	const query = "SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME " +
		"FROM information_schema.STATISTICS " +
		"WHERE TABLE_SCHEMA = COALESCE(NULLIF(?, ''), DATABASE()) " +
		"ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX"

	var rows []mysqlIndexRow
	if err := db.SelectContext(ctx, &rows, query, database); err != nil {
		return nil, fmt.Errorf("query mysql indexes: %w", err)
	}

	indexes := make(map[string]map[string]*meta.Index)
	for _, row := range rows {
		if row.Index == "PRIMARY" {
			continue
		}
		if indexes[row.Table] == nil {
			indexes[row.Table] = make(map[string]*meta.Index)
		}
		idx, ok := indexes[row.Table][row.Index]
		if !ok {
			idx = &meta.Index{Name: row.Index, Unique: row.NonUnique == 0}
		}
		idx.Cols = append(idx.Cols, row.Column)
		indexes[row.Table][idx.Name] = idx
	}
	return indexes, nil
}

type mysqlDatabaseRow struct {
	Name    string `db:"SCHEMA_NAME"`
	Charset string `db:"DEFAULT_CHARACTER_SET_NAME"`
	Collate string `db:"DEFAULT_COLLATION_NAME"`
}

// Databases 读取 information_schema.SCHEMATA，过滤系统库。
func (MySQL) Databases(ctx context.Context, db *sqlx.DB) ([]*meta.Database, error) {
	const query = "SELECT SCHEMA_NAME, DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME " +
		"FROM information_schema.SCHEMATA " +
		"WHERE SCHEMA_NAME NOT IN ('information_schema', 'performance_schema', 'mysql', 'sys') " +
		"ORDER BY SCHEMA_NAME"

	var rows []mysqlDatabaseRow
	if err := db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("query mysql databases: %w", err)
	}

	dbs := make([]*meta.Database, 0, len(rows))
	for _, row := range rows {
		dbs = append(dbs, &meta.Database{Name: row.Name, Charset: row.Charset, Collate: row.Collate})
	}
	return dbs, nil
}

// normalizeMySQLType strips display width/enum members so the canonical name
// (e.g. VARCHAR, TEXT, INT UNSIGNED) matches the sqlKinds table.
func normalizeMySQLType(columnType string) string {
	name := columnType
	if i := strings.IndexByte(name, '('); i >= 0 {
		name = name[:i]
	}
	return strings.ToUpper(strings.TrimSpace(name))
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// assemble merges column and index metadata into the table descriptors. It is
// shared by both dialects because the merge logic is dialect-neutral.
func assemble(tables []*meta.Table, columns map[string]map[string]*meta.Column, indexes map[string]map[string]*meta.Index) {
	for _, t := range tables {
		t.Columns = columns[t.Name]
		t.Indexes = indexes[t.Name]
	}
}
