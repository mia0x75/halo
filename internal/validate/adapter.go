package validate

import (
	"context"
	"errors"

	"github.com/mia0x75/halo/internal/agent/meta"
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// ErrNotImplemented marks the 2D-A adapter skeleton. The translation from the
// legacy engine output to domain.Violation lands in 2D-B together with the
// schema source.
var ErrNotImplemented = errors.New("validate: adapter not implemented")

// RuleSource loads the validation rules from mm_rules.
type RuleSource interface {
	Load(ctx context.Context) ([]*domain.Rule, error)
}

// SchemaSource loads the target database metadata. It returns *meta.Table so
// the rule engine keeps inspecting charset/collate/engine exactly as before;
// the neutral type preserves the same fields the engine reads.
type SchemaSource interface {
	Load(ctx context.Context, cluster *domain.Cluster) (map[string][]*meta.Table, error)
}

// Adapter implements domain.Validator on top of the legacy rule engine.
type Adapter struct {
	rules  RuleSource
	schema SchemaSource
}

var _ domain.Validator = (*Adapter)(nil)

func NewAdapter(r RuleSource, s SchemaSource) *Adapter {
	return &Adapter{rules: r, schema: s}
}

// Validate runs the rule engine over one statement and returns the hits.
//
// The legacy engine still owns the IO (rule set from the cache, target
// metadata from the agent); the adapter only translates the outcome into
// domain.Violation. A request without a target cluster cannot be checked,
// because the existence and charset rules read the target metadata.
func (a *Adapter) Validate(ctx context.Context, req domain.ValidateRequest) ([]domain.Violation, error) {
	if req.Cluster == nil {
		return nil, nil
	}

	stmt := &models.Statement{
		Content:    req.SQL,
		Violations: &models.Violations{},
	}
	cluster := bridge.ClusterFromDomain(req.Cluster)
	ticket := &models.Ticket{Database: req.Database}

	Run([]*models.Statement{stmt}, cluster, ticket)

	violations := make([]domain.Violation, 0)
	for _, clause := range stmt.Violations.Clauses() {
		violations = append(violations, domain.Violation{
			Rule:    clause.Rule,
			Message: clause.Description,
			Level:   clause.Level,
		})
	}
	return violations, nil
}
