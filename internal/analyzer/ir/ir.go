// Package ir defines the dialect-neutral intermediate representation shared by
// the MySQL and PostgreSQL analyzers. Backends lower their own ASTs into these
// types so rule checking and auditing never depend on a backend AST.
package ir

import "github.com/mia0x75/halo/internal/domain"

// StatementKind is the coarse statement category used by rule matching.
type StatementKind string

const (
	KindSelect StatementKind = "select"
	KindInsert StatementKind = "insert"
	KindUpdate StatementKind = "update"
	KindDelete StatementKind = "delete"
	KindDDL    StatementKind = "ddl"
	KindOther  StatementKind = "other"
)

// Statement is one parsed SQL statement reduced to what auditing needs.
type Statement struct {
	Kind    StatementKind
	Dialect domain.ClusterType
	Raw     string
	Tables  []TableRef
	Columns []ColumnUsage
	Where   []Condition
	Joins   []JoinRef
	DDL     *DDLInfo
}

// TableRef names a table, qualified by schema and/or alias when available.
type TableRef struct {
	Schema string
	Name   string
	Alias  string
}

// ColumnUsageKind tells how a statement touches a column.
type ColumnUsageKind string

const (
	ColumnRead  ColumnUsageKind = "read"
	ColumnWrite ColumnUsageKind = "write"
	ColumnWhere ColumnUsageKind = "where"
	ColumnJoin  ColumnUsageKind = "join"
)

// ColumnRef identifies a column, table-qualified when the SQL provides it.
type ColumnRef struct {
	Table string
	Name  string
}

// ColumnUsage pairs a column with the way the statement uses it.
type ColumnUsage struct {
	Column ColumnRef
	Usage  ColumnUsageKind
}

// Condition is a predicate kept as text; a structured form is deferred.
type Condition struct {
	Raw string
}

// JoinRef describes a single join.
type JoinRef struct {
	Type  string
	Table TableRef
	On    string
}

// DDLInfo captures the object and action of a DDL statement.
type DDLInfo struct {
	Action  string
	Object  string
	Targets []TableRef
}

// Result is the outcome of analyzing one statement. Violations are owned by
// domain so the IR stays independent of any particular rule engine.
type Result struct {
	Statement  *Statement
	Violations []domain.Violation
}
