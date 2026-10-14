package rules

import "github.com/mia0x75/halo/internal/analyzer/ir"

// CheckPostgres applies PostgreSQL-specific rules. The MVP relies on
// CheckCommon; this entry point keeps PG-only rules separate from the
// dialect-agnostic set.
// TODO(phase2): add PostgreSQL-specific dangerous-operation rules.
func CheckPostgres(st *ir.Statement) []ir.Violation {
	return nil
}
