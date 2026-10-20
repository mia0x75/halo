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
