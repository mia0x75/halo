package analyzer

import (
	"context"
	"errors"
	"testing"

	"github.com/mia0x75/halo/internal/analyzer/ir"
	"github.com/mia0x75/halo/internal/domain"
)

func TestForUnknownType(t *testing.T) {
	if _, err := For(domain.ClusterType("oracle")); err == nil {
		t.Fatal("expected error for unknown cluster type")
	}
	for _, ct := range []domain.ClusterType{domain.ClusterTypeMySQL, domain.ClusterTypePostgres} {
		if _, err := For(ct); err != nil {
			t.Fatalf("For(%q): %v", ct, err)
		}
	}
}

func TestMySQLAnalyzeCommonRules(t *testing.T) {
	a := MySQL{}
	cases := []struct {
		sql    string
		kind   ir.StatementKind
		rule   string
		hasHit bool
	}{
		{"UPDATE users SET name = 'x'", ir.KindUpdate, "update-without-where", true},
		{"DELETE FROM users WHERE id = 1", ir.KindDelete, "", false},
		{"DELETE FROM users", ir.KindDelete, "delete-without-where", true},
		{"TRUNCATE TABLE users", ir.KindDDL, "truncate", true},
		{"DROP TABLE users", ir.KindDDL, "drop-table", true},
		{"SELECT * FROM users", ir.KindSelect, "", false},
	}
	for _, tc := range cases {
		res, err := a.Analyze(context.Background(), tc.sql)
		if err != nil {
			t.Fatalf("Analyze(%q): %v", tc.sql, err)
		}
		if res.Statement.Kind != tc.kind {
			t.Errorf("Analyze(%q): kind = %q, want %q", tc.sql, res.Statement.Kind, tc.kind)
		}
		if tc.hasHit {
			if !containsRule(res.Violations, tc.rule) {
				t.Errorf("Analyze(%q): missing violation %q, got %+v", tc.sql, tc.rule, res.Violations)
			}
			continue
		}
		if len(res.Violations) != 0 {
			t.Errorf("Analyze(%q): unexpected violations %+v", tc.sql, res.Violations)
		}
	}
}

func TestPostgresAnalyzeFallsBackToCommonRules(t *testing.T) {
	a := Postgres{}
	res, err := a.Analyze(context.Background(), "UPDATE users SET name = 'x'")
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if res.Statement.Dialect != domain.ClusterTypePostgres {
		t.Fatalf("dialect = %q, want postgres", res.Statement.Dialect)
	}
	if !containsRule(res.Violations, "update-without-where") {
		t.Fatalf("expected update-without-where, got %+v", res.Violations)
	}

	res, err = a.Analyze(context.Background(), "SELECT 1")
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if len(res.Violations) != 0 {
		t.Fatalf("unexpected violations %+v", res.Violations)
	}
}

func TestParseErrorIsWrapped(t *testing.T) {
	_, err := (MySQL{}).Parse(context.Background(), "SELECT FROM WHERE")
	if !errors.Is(err, ErrParseFailed) {
		t.Fatalf("err = %v, want ErrParseFailed", err)
	}
}

func containsRule(vs []ir.Violation, rule string) bool {
	for _, v := range vs {
		if v.Rule == rule {
			return true
		}
	}
	return false
}
