package analyzer

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/analyzer/ir"
	"github.com/mia0x75/halo/internal/analyzer/rules"
	"github.com/mia0x75/halo/internal/domain"

	pgparser "github.com/pgplex/pgparser/parser"
)

// Postgres analyzes PostgreSQL-dialect SQL with pgplex/pgparser.
type Postgres struct{}

// Parse parses sql and lowers it into the shared IR. Lowering is not built yet,
// so the statement is classified as KindOther and only dialect/raw are filled.
func (Postgres) Parse(ctx context.Context, sql string) (*ir.Statement, error) {
	stmts, err := pgparser.Parse(sql)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseFailed, err)
	}
	if stmts == nil || len(stmts.Items) == 0 {
		return nil, ErrUnsupportedStmt
	}
	return postgresToIR(sql), nil
}

// Analyze parses sql and runs the common rule set over the resulting IR.
func (p Postgres) Analyze(ctx context.Context, sql string) (*ir.Result, error) {
	st, err := p.Parse(ctx, sql)
	if err != nil {
		return nil, err
	}
	return &ir.Result{Statement: st, Violations: rules.CheckCommon(st)}, nil
}

func postgresToIR(sql string) *ir.Statement {
	// TODO(phase2): walk nodes.List to lower Tables, Columns, Where, Joins, DDL.
	return &ir.Statement{
		Kind:    ir.KindOther,
		Dialect: domain.ClusterTypePostgres,
		Raw:     sql,
	}
}
