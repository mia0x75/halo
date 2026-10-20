package analyzer

import (
	"context"
	"errors"
	"testing"

	"github.com/mia0x75/halo/internal/analyzer/ir"
	"github.com/mia0x75/halo/internal/domain"
)

// fakeValidator records the request it receives so delegation can be asserted
// without depending on the concrete rule engine.
type fakeValidator struct {
	got        domain.ValidateRequest
	violations []domain.Violation
	err        error
}

func (f *fakeValidator) Validate(_ context.Context, req domain.ValidateRequest) ([]domain.Violation, error) {
	f.got = req
	return f.violations, f.err
}

func TestNewUnknownType(t *testing.T) {
	if _, err := New(domain.ClusterType("oracle"), nil); err == nil {
		t.Fatal("expected error for unknown cluster type")
	}
	for _, ct := range []domain.ClusterType{domain.ClusterTypeMySQL, domain.ClusterTypePostgres} {
		if _, err := New(ct, nil); err != nil {
			t.Fatalf("New(%q): %v", ct, err)
		}
	}
}

func TestMySQLParseKinds(t *testing.T) {
	a, err := New(domain.ClusterTypeMySQL, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	cases := []struct {
		sql  string
		kind ir.StatementKind
	}{
		{"SELECT * FROM users", ir.KindSelect},
		{"INSERT INTO users VALUES (1)", ir.KindInsert},
		{"UPDATE users SET name = 'x'", ir.KindUpdate},
		{"DELETE FROM users WHERE id = 1", ir.KindDelete},
		{"TRUNCATE TABLE users", ir.KindDDL},
		{"DROP TABLE users", ir.KindDDL},
	}
	for _, tc := range cases {
		st, err := a.Parse(context.Background(), tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.sql, err)
		}
		if st.Kind != tc.kind {
			t.Errorf("Parse(%q): kind = %q, want %q", tc.sql, st.Kind, tc.kind)
		}
	}
}

func TestPostgresParseDialect(t *testing.T) {
	a, err := New(domain.ClusterTypePostgres, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	st, err := a.Parse(context.Background(), "SELECT 1")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if st.Dialect != domain.ClusterTypePostgres {
		t.Fatalf("dialect = %q, want postgres", st.Dialect)
	}
}

func TestAnalyzeDelegatesToValidator(t *testing.T) {
	const sql = "UPDATE users SET name = 'x'"
	cluster := &domain.Cluster{Alias: "target-1"}
	fv := &fakeValidator{violations: []domain.Violation{{Message: "boom", Level: 2}}}
	a, err := New(domain.ClusterTypeMySQL, fv)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	res, err := a.Analyze(context.Background(), sql, cluster, "testdb")
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if len(res.Violations) != 1 || res.Violations[0].Message != "boom" {
		t.Fatalf("violations = %+v, want one boom", res.Violations)
	}
	if fv.got.SQL != sql {
		t.Errorf("validator SQL = %q, want %q", fv.got.SQL, sql)
	}
	if fv.got.Cluster != cluster {
		t.Errorf("validator cluster = %p, want %p", fv.got.Cluster, cluster)
	}
}

func TestAnalyzeNilValidatorSkipsRules(t *testing.T) {
	a, err := New(domain.ClusterTypeMySQL, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	res, err := a.Analyze(context.Background(), "DELETE FROM users", nil, "testdb")
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}
	if res.Statement.Kind != ir.KindDelete {
		t.Errorf("kind = %q, want delete", res.Statement.Kind)
	}
	if len(res.Violations) != 0 {
		t.Errorf("violations = %+v, want none", res.Violations)
	}
}

func TestAnalyzeValidatorErrorIsWrapped(t *testing.T) {
	sentinel := errors.New("engine down")
	a, err := New(domain.ClusterTypeMySQL, &fakeValidator{err: sentinel})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := a.Analyze(context.Background(), "SELECT 1", nil, "testdb"); !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want wrapping %v", err, sentinel)
	}
}

func TestParseErrorIsWrapped(t *testing.T) {
	a, err := New(domain.ClusterTypeMySQL, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := a.Parse(context.Background(), "SELECT FROM WHERE"); !errors.Is(err, ErrParseFailed) {
		t.Fatalf("err = %v, want ErrParseFailed", err)
	}
}
