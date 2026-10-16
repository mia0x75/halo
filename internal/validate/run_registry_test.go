package validate

import (
	"testing"

	"github.com/mia0x75/parser"
	"github.com/mia0x75/parser/ast"
	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/agent/meta"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// TestRunRegistryEmptyContext exercises the pure engine with no statements: it
// must be a no-op and must not touch the target database.
func TestRunRegistryEmptyContext(t *testing.T) {
	require.Empty(t, RunRegistry(&Context{}))
}

// TestCollectViolations flattens per-statement hits.
func TestCollectViolations(t *testing.T) {
	s := &Subject{Violations: &models.Violations{}}
	s.Violations.Append(&models.Clause{Level: 1, Description: "boom"})

	got := collectViolations([]*Subject{s})
	require.Len(t, got, 1)
	require.Same(t, s, got[0].Statement)
	require.Equal(t, "boom", got[0].Clause.Description)

	require.Empty(t, collectViolations([]*Subject{{}}))
}

// TestRunRegistryVisitInfoIsolation guards against VisitInfo leaking across
// statements: validators are singletons, so a second statement must not inherit
// the tables visited by the first. The first statement points at a missing
// table (violation), the second at an existing one (clean); with the leak the
// second would wrongly inherit the missing-table visit and be flagged too.
func TestRunRegistryVisitInfoIsolation(t *testing.T) {
	p := parser.New()
	parse := func(s string) ast.StmtNode {
		node, err := p.ParseOneStmt(s, "", "")
		require.NoError(t, err)
		return node
	}

	first := &Subject{
		Content:    "CREATE INDEX i1 ON ghost(id);",
		Violations: &models.Violations{},
		StmtNode:   parse("CREATE INDEX i1 ON ghost(id);"),
	}
	second := &Subject{
		Content:    "CREATE INDEX i2 ON t1(id);",
		Violations: &models.Violations{},
		StmtNode:   parse("CREATE INDEX i2 ON t1(id);"),
	}

	ctx := &Context{
		Ticket:    &domain.Ticket{Database: "starwars"},
		Databases: []models.Database{{Name: "starwars"}},
		Tables: map[string][]*meta.Table{
			"starwars": {{Name: "t1", Columns: map[string]*meta.Column{"id": {Name: "id"}}}},
		},
		Subjects: []*Subject{first, second},
		Rules: []*Rule{{Rule: domain.Rule{
			Name:      "CIX-L3-002",
			VldrGroup: 150,
			Func:      "TargetTableDoesNotExist",
			Bitwise:   1,
			Level:     1,
			Message:   "table missing",
		}}},
	}

	RunRegistry(ctx)

	require.NotEmpty(t, first.Violations.Marshal(), "missing table must be flagged")
	require.Empty(t, second.Violations.Marshal(), "existing table must not inherit the first visit")
}
