package analyzer

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/analyzer/ir"
	"github.com/mia0x75/halo/internal/domain"
)

// dialectParser lowers one dialect's SQL into the shared IR. Each backend
// implements only this; rule checking is shared at the analyzer level.
type dialectParser interface {
	Parse(ctx context.Context, sql string) (*ir.Statement, error)
}

// Analyzer parses SQL for one dialect and delegates rule checking to a
// domain.Validator. The validator is injected so this package never depends on
// the concrete rule engine.
type Analyzer struct {
	parser    dialectParser
	validator domain.Validator
}

// New returns an analyzer for a target cluster type. validator may be nil until
// the rule engine adapter is wired in phase 2D; Analyze then returns the parsed
// statement without violations.
func New(t domain.ClusterType, v domain.Validator) (*Analyzer, error) {
	var p dialectParser
	switch t {
	case domain.ClusterTypeMySQL:
		p = MySQL{}
	case domain.ClusterTypePostgres:
		p = Postgres{}
	default:
		return nil, fmt.Errorf("analyzer: unsupported cluster type %q", t)
	}
	return &Analyzer{parser: p, validator: v}, nil
}

// Parse lowers sql into the shared IR for the configured dialect.
func (a *Analyzer) Parse(ctx context.Context, sql string) (*ir.Statement, error) {
	return a.parser.Parse(ctx, sql)
}

// Analyze parses sql then asks the validator to check it. cluster and database
// are forwarded to the validator because the rules inspect the target
// metadata, which is database scoped; they may be empty before the engine is
// wired.
func (a *Analyzer) Analyze(ctx context.Context, sql string, cluster *domain.Cluster, database string) (*ir.Result, error) {
	st, err := a.Parse(ctx, sql)
	if err != nil {
		return nil, err
	}
	res := &ir.Result{Statement: st}
	if a.validator == nil {
		return res, nil
	}
	vs, err := a.validator.Validate(ctx, domain.ValidateRequest{SQL: sql, Cluster: cluster, Database: database})
	if err != nil {
		return nil, fmt.Errorf("analyzer: validate: %w", err)
	}
	res.Violations = vs
	return res, nil
}
