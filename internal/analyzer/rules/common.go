package rules

import (
	"regexp"
	"strings"

	"github.com/mia0x75/halo/internal/analyzer/ir"
)

// whereRe is a fallback predicate scan used until backends populate IR.Where.
var whereRe = regexp.MustCompile(`(?i)\bwhere\b`)

// CheckCommon applies dialect-agnostic dangerous-statement rules to st. It is
// the safe fallback while backends cannot yet produce a detailed IR.
func CheckCommon(st *ir.Statement) []ir.Violation {
	if st == nil {
		return nil
	}

	var violations []ir.Violation
	upper := strings.ToUpper(strings.TrimSpace(st.Raw))
	switch {
	case strings.HasPrefix(upper, "TRUNCATE"):
		violations = append(violations, ir.Violation{
			Rule:     "truncate",
			Message:  "TRUNCATE 会清空整张表，请确认已备份。",
			Severity: ir.SeverityCritical,
		})
	case strings.HasPrefix(upper, "DROP TABLE"):
		violations = append(violations, ir.Violation{
			Rule:     "drop-table",
			Message:  "DROP TABLE 会删除整张表及其数据。",
			Severity: ir.SeverityCritical,
		})
	case effectiveKind(st) == ir.KindUpdate && !hasWhere(st):
		violations = append(violations, ir.Violation{
			Rule:     "update-without-where",
			Message:  "UPDATE 缺少 WHERE 条件，将更新全表。",
			Severity: ir.SeverityCritical,
		})
	case effectiveKind(st) == ir.KindDelete && !hasWhere(st):
		violations = append(violations, ir.Violation{
			Rule:     "delete-without-where",
			Message:  "DELETE 缺少 WHERE 条件，将删除全表。",
			Severity: ir.SeverityCritical,
		})
	}
	return violations
}

// effectiveKind trusts the backend IR kind and falls back to the raw text when
// the backend only produced KindOther (e.g. the PostgreSQL placeholder).
func effectiveKind(st *ir.Statement) ir.StatementKind {
	if st.Kind != ir.KindOther {
		return st.Kind
	}
	u := strings.ToUpper(strings.TrimSpace(st.Raw))
	switch {
	case strings.HasPrefix(u, "SELECT"):
		return ir.KindSelect
	case strings.HasPrefix(u, "INSERT"):
		return ir.KindInsert
	case strings.HasPrefix(u, "UPDATE"):
		return ir.KindUpdate
	case strings.HasPrefix(u, "DELETE"):
		return ir.KindDelete
	}
	return ir.KindOther
}

// hasWhere prefers the structured IR. TODO(phase2): drop the raw scan once
// backends populate IR.Where.
func hasWhere(st *ir.Statement) bool {
	if len(st.Where) > 0 {
		return true
	}
	return whereRe.MatchString(st.Raw)
}
