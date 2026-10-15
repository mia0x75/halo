package dialect

// Dialect-neutral metadata descriptors. They mirror the subset of
// go-xorm/core's Table/Column/Index that validate/ and the Metadata resolver
// actually consume, so the rule engine keeps its semantics after xorm removal.
//
// core.* -> dialect.* mapping:
//   core.Table.Name         -> Table.Name
//   core.Table.StoreEngine  -> Table.Engine
//   core.Table.Charset      -> Table.Charset
//   core.Table.Collate      -> Table.Collate
//   core.Table.Comment      -> Table.Comment
//   core.Table.columns      -> Table.Columns (map, was []*Column + GetColumn)
//   core.Table.Indexes      -> Table.Indexes
//   core.Column.Name        -> Column.Name
//   core.Column.SQLType     -> Column.SQLType (methods IsText/IsBlob kept)
//   core.Column.Nullable    -> Column.Nullable
//   core.Column.Default      -> Column.Default (string, "" when absent)
//   core.Column.Comment      -> Column.Comment
//   core.Index.Name          -> Index.Name
//   core.Index.Type          -> Index.Unique (UniqueType -> true)
//   core.Index.Cols          -> Index.Cols
//
// Fields with no dialect equivalent (Engine/Charset on PostgreSQL) stay empty.

// SQLType names a column type the way xorm's core.SQLType did, so the
// IsText/IsBlob classification used by index rules keeps working.
type SQLType struct {
	Name string
}

// IsText reports whether the type belongs to the text family.
func (s SQLType) IsText() bool { return sqlKinds[s.Name] == kindText }

// IsBlob reports whether the type belongs to the binary family.
func (s SQLType) IsBlob() bool { return sqlKinds[s.Name] == kindBlob }

type sqlKind int

const (
	kindOther sqlKind = iota
	kindText
	kindBlob
)

// sqlKinds classifies the canonical (upper-cased) type names both dialects
// normalize to. Unknown types fall back to kindOther, same as xorm did.
var sqlKinds = map[string]sqlKind{
	"ENUM":       kindText,
	"SET":        kindText,
	"JSON":       kindText,
	"JSONB":      kindText,
	"CHAR":       kindText,
	"VARCHAR":    kindText,
	"NCHAR":      kindText,
	"NVARCHAR":   kindText,
	"TINYTEXT":   kindText,
	"TEXT":       kindText,
	"MEDIUMTEXT": kindText,
	"LONGTEXT":   kindText,

	"BINARY":     kindBlob,
	"VARBINARY":  kindBlob,
	"TINYBLOB":   kindBlob,
	"BLOB":       kindBlob,
	"MEDIUMBLOB": kindBlob,
	"LONGBLOB":   kindBlob,
	"BYTEA":      kindBlob,
}

// Column 列级元数据。
type Column struct {
	Name     string
	SQLType  SQLType
	Nullable bool
	Default  string
	Comment  string
}

// Index 索引级元数据。
type Index struct {
	Name   string
	Unique bool
	Cols   []string
}

// Table 表级元数据。PG 无对应的字段（Engine/Charset）留空。
type Table struct {
	Name    string
	Engine  string
	Charset string
	Collate string
	Comment string
	Columns map[string]*Column
	Indexes map[string]*Index
}

// GetColumn returns the column metadata by name, nil when absent.
func (t *Table) GetColumn(name string) *Column {
	return t.Columns[name]
}

// Database database-level metadata. It replaces models.Database, which the
// resolver layer exposes to clients; Charset/Collate are dialect defaults and
// stay empty on PostgreSQL.
type Database struct {
	Name    string
	Charset string
	Collate string
}
