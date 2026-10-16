package postgresrepo

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/validate"
)

// ruleSourceLimit pages in every rule in one shot. The rule set is seed data
// measured in tens of rows, so a single large page avoids cursor bookkeeping.
const ruleSourceLimit = 1000

// ruleSourceCursor sorts after every UUID, so List's keyset predicate
// (uuid < cursor) yields the whole set on the first page. An empty cursor would
// match nothing.
const ruleSourceCursor = "ffffffff-ffff-ffff-ffff-ffffffffffff"

var _ validate.RuleSource = (*RuleSource)(nil)

// RuleSource adapts the metadata rule table to the validate engine's input.
type RuleSource struct {
	rules domain.RuleRepository
}

func NewRuleSource(rules domain.RuleRepository) *RuleSource {
	return &RuleSource{rules: rules}
}

// Load returns every rule as a domain entity for the validate engine.
func (s *RuleSource) Load(ctx context.Context) ([]*domain.Rule, error) {
	rows, err := s.rules.List(ctx, ruleSourceCursor, ruleSourceLimit)
	if err != nil {
		return nil, fmt.Errorf("load rules: %w", err)
	}
	return rows, nil
}
