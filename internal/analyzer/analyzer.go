package analyzer

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/analyzer/ir"
	"github.com/mia0x75/halo/internal/domain"
)

// Analyzer parses and rule-checks SQL for one target dialect.
type Analyzer interface {
	Parse(ctx context.Context, sql string) (*ir.Statement, error)
	Analyze(ctx context.Context, sql string) (*ir.Result, error)
}

// For returns the analyzer matching a target cluster type.
func For(t domain.ClusterType) (Analyzer, error) {
	switch t {
	case domain.ClusterTypeMySQL:
		return MySQL{}, nil
	case domain.ClusterTypePostgres:
		return Postgres{}, nil
	default:
		return nil, fmt.Errorf("analyzer: unsupported cluster type %q", t)
	}
}
