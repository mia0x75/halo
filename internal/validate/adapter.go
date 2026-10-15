package validate

import (
	"context"
	"errors"

	"github.com/mia0x75/halo/internal/agent/meta"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// ErrNotImplemented marks the 2D-A adapter skeleton. The translation from the
// legacy engine output to domain.Violation lands in 2D-B together with the
// schema source.
var ErrNotImplemented = errors.New("validate: adapter not implemented")

// RuleSource loads the validation rules from mm_rules.
type RuleSource interface {
	Load(ctx context.Context) ([]*models.Rule, error)
}

// SchemaSource loads the target database metadata. It returns *meta.Table so
// the rule engine keeps inspecting charset/collate/engine exactly as before;
// the neutral type preserves the same fields the engine reads.
type SchemaSource interface {
	Load(ctx context.Context, cluster *models.Cluster) (map[string][]*meta.Table, error)
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

// Validate is a 2D-A skeleton: the engine is not yet wired to domain types.
func (a *Adapter) Validate(ctx context.Context, req domain.ValidateRequest) ([]domain.Violation, error) {
	_ = a.rules
	_ = a.schema
	return nil, ErrNotImplemented
}
