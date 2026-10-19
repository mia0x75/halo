// Package meta holds the dialect-neutral metadata model shared by the target-db
// dialects and their consumers (validate rules, resolvers).
package meta

// ColumnType names a column type after both dialects normalize it, so the
// IsText/IsBlob classification used by index rules stays dialect-independent.
type ColumnType string

// IsText reports whether the type belongs to the text family.
func (t ColumnType) IsText() bool { return sqlKinds[t] == kindText }

// IsBlob reports whether the type belongs to the binary family.
func (t ColumnType) IsBlob() bool { return sqlKinds[t] == kindBlob }

// IsLargeObject reports whether the type stores arbitrary-length payloads (the
// TEXT/BLOB families). Unlike IsText it excludes CHAR/VARCHAR, so rules that
// forbid returning large payloads do not reject ordinary string columns.
func (t ColumnType) IsLargeObject() bool {
	switch t {
	case "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT",
		"TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB":
		return true
	}
	return false
}

type sqlKind int

const (
	kindOther sqlKind = iota
	kindText
	kindBlob
)

// sqlKinds classifies the canonical (upper-cased) type names both dialects
// normalize to. Unknown types fall back to kindOther.
var sqlKinds = map[ColumnType]sqlKind{
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

// Column is column-level metadata.
type Column struct {
	Name     string
	Type     ColumnType
	Nullable bool
	Default  string
	Comment  string
	// AutoIncrement marks identity/auto_increment columns so column-level rules
	// do not require an explicit value for them.
	AutoIncrement bool
}

// Index is index-level metadata.
type Index struct {
	Name   string
	Unique bool
	Cols   []string
}

// Table is table-level metadata. Fields without a counterpart on the target
// dialect (Engine/Charset on PostgreSQL) stay empty.
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

// Database is database-level metadata. Charset/Collate are dialect defaults and
// stay empty on PostgreSQL.
type Database struct {
	Name    string
	Charset string
	Collate string
}

// ObjectKind identifies a user-defined object whose existence the table/column
// model does not cover.
type ObjectKind string

const (
	ObjectView      ObjectKind = "view"
	ObjectFunction  ObjectKind = "function"
	ObjectProcedure ObjectKind = "procedure"
	ObjectTrigger   ObjectKind = "trigger"
	ObjectEvent     ObjectKind = "event"
)

// Objects holds the names of user-defined objects in one database. Names are
// schema-scoped, matching how the review rules and the MySQL parser address
// these objects (DROP TRIGGER carries no table, so triggers are keyed by name).
type Objects struct {
	Views      map[string]struct{}
	Functions  map[string]struct{}
	Procedures map[string]struct{}
	Triggers   map[string]struct{}
	Events     map[string]struct{}
}

// NewObjects returns an Objects with every set initialised.
func NewObjects() *Objects {
	return &Objects{
		Views:      map[string]struct{}{},
		Functions:  map[string]struct{}{},
		Procedures: map[string]struct{}{},
		Triggers:   map[string]struct{}{},
		Events:     map[string]struct{}{},
	}
}

// Has reports whether an object of the given kind exists.
func (o *Objects) Has(kind ObjectKind, name string) bool {
	if o == nil {
		return false
	}
	var set map[string]struct{}
	switch kind {
	case ObjectView:
		set = o.Views
	case ObjectFunction:
		set = o.Functions
	case ObjectProcedure:
		set = o.Procedures
	case ObjectTrigger:
		set = o.Triggers
	case ObjectEvent:
		set = o.Events
	default:
		return false
	}
	_, ok := set[name]
	return ok
}

// Add records an object name under the given kind.
func (o *Objects) Add(kind ObjectKind, name string) {
	switch kind {
	case ObjectView:
		o.Views[name] = struct{}{}
	case ObjectFunction:
		o.Functions[name] = struct{}{}
	case ObjectProcedure:
		o.Procedures[name] = struct{}{}
	case ObjectTrigger:
		o.Triggers[name] = struct{}{}
	case ObjectEvent:
		o.Events[name] = struct{}{}
	}
}
