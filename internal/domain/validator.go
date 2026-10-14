package domain

import "context"

// Validator checks a SQL statement against business rules. Implementations live
// outside this package (the legacy validate/ engine provides one) so that domain
// stays free of parser, ORM and driver dependencies.
type Validator interface {
	Validate(ctx context.Context, req ValidateRequest) ([]Violation, error)
}

// ValidateRequest carries what a checker needs: the raw SQL and the target
// cluster whose metadata the rules inspect.
type ValidateRequest struct {
	SQL     string
	Cluster *Cluster
}

// Violation is a single rule hit. Level is the raw mm_rules.level value
// (0=message, 1=error, 2=warning); the current engine does not expose the rule
// name or vldr_group, so those fields are added only when the rule engine is
// ported into this layer.
type Violation struct {
	Message string
	Level   uint8
}
