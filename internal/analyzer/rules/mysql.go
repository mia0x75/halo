package rules

import "github.com/mia0x75/halo/internal/analyzer/ir"

// CheckMySQL applies MySQL-specific rules. The MVP relies on CheckCommon; this
// entry point keeps MySQL-only rules separate from the dialect-agnostic set.
// TODO(phase2): add MySQL-specific dangerous-operation rules.
func CheckMySQL(st *ir.Statement) []ir.Violation {
	return nil
}
