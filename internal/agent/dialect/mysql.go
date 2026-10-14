package dialect

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// MySQL 目标库方言。
type MySQL struct{}

// QuoteIdent MySQL 用反引号，`type` 等保留字必须引用。
func (MySQL) QuoteIdent(s string) string {
	return "`" + s + "`"
}

type mysqlColumnRow struct {
	Table    string  `db:"TABLE_NAME"`
	Column   string  `db:"COLUMN_NAME"`
	Type     string  `db:"COLUMN_TYPE"`
	Nullable string  `db:"IS_NULLABLE"`
	Default  *string `db:"COLUMN_DEFAULT"`
	Key      string  `db:"COLUMN_KEY"`
	Comment  string  `db:"COLUMN_COMMENT"`
}

type mysqlIndexRow struct {
	Table     string `db:"TABLE_NAME"`
	Index     string `db:"INDEX_NAME"`
	NonUnique int    `db:"NON_UNIQUE"`
	Column    string `db:"COLUMN_NAME"`
}

// LoadColumns 读取 information_schema.COLUMNS。schema 为空时用 DATABASE()。
func (MySQL) LoadColumns(ctx context.Context, db *sqlx.DB, schema string) (map[string]map[string]ColumnInfo, error) {
	const query = "SELECT TABLE_NAME, COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, " +
		"COLUMN_DEFAULT, COLUMN_KEY, COLUMN_COMMENT " +
		"FROM information_schema.COLUMNS " +
		"WHERE TABLE_SCHEMA = COALESCE(NULLIF(?, ''), DATABASE()) " +
		"ORDER BY TABLE_NAME, ORDINAL_POSITION"

	var rows []mysqlColumnRow
	if err := db.SelectContext(ctx, &rows, query, schema); err != nil {
		return nil, fmt.Errorf("query mysql columns: %w", err)
	}

	columns := make(map[string]map[string]ColumnInfo, len(rows))
	for _, row := range rows {
		info := ColumnInfo{
			Name:     row.Column,
			Type:     row.Type,
			Nullable: row.Nullable == "YES",
			Default:  row.Default,
			Key:      row.Key,
			Comment:  row.Comment,
		}
		if columns[row.Table] == nil {
			columns[row.Table] = make(map[string]ColumnInfo)
		}
		columns[row.Table][row.Column] = info
	}
	return columns, nil
}

// LoadIndexes 读取 information_schema.STATISTICS，忽略主键。
func (MySQL) LoadIndexes(ctx context.Context, db *sqlx.DB, schema string) (map[string]map[string]IndexInfo, error) {
	const query = "SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME " +
		"FROM information_schema.STATISTICS " +
		"WHERE TABLE_SCHEMA = COALESCE(NULLIF(?, ''), DATABASE()) " +
		"ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX"

	var rows []mysqlIndexRow
	if err := db.SelectContext(ctx, &rows, query, schema); err != nil {
		return nil, fmt.Errorf("query mysql indexes: %w", err)
	}

	indexes := make(map[string]map[string]IndexInfo)
	for _, row := range rows {
		if row.Index == "PRIMARY" {
			continue
		}
		if indexes[row.Table] == nil {
			indexes[row.Table] = make(map[string]IndexInfo)
		}
		info, ok := indexes[row.Table][row.Index]
		if !ok {
			info = IndexInfo{Name: row.Index, Unique: row.NonUnique == 0}
		}
		info.Columns = append(info.Columns, row.Column)
		indexes[row.Table][row.Index] = info
	}
	return indexes, nil
}
