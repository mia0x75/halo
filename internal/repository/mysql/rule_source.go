package mysqlrepo

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
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

// Load returns every rule as the legacy models.Rule the engine consumes.
func (s *RuleSource) Load(ctx context.Context) ([]*models.Rule, error) {
	rows, err := s.rules.List(ctx, ruleSourceCursor, ruleSourceLimit)
	if err != nil {
		return nil, fmt.Errorf("load rules: %w", err)
	}
	out := make([]*models.Rule, 0, len(rows))
	for _, r := range rows {
		out = append(out, legacyRuleFromDomain(r))
	}
	return out, nil
}

// legacyRuleFromDomain bridges domain.Rule to the legacy models.Rule. The field
// sets are identical; only the numeric widths differ, so every cast is lossless.
// UpdateAt is never read by the engine, so a nil maps to 0.
func legacyRuleFromDomain(r *domain.Rule) *models.Rule {
	var updateAt uint
	if r.UpdateAt != nil {
		updateAt = *r.UpdateAt
	}
	return &models.Rule{
		Name:        r.Name,
		UUID:        r.UUID,
		Group:       r.Group,
		Level:       uint8(r.Level),
		VldrGroup:   uint16(r.VldrGroup),
		Operator:    r.Operator,
		Values:      r.Values,
		Bitwise:     r.Bitwise,
		Func:        r.Func,
		Message:     r.Message,
		Description: r.Description,
		Element:     r.Element,
		Version:     r.Version,
		UpdateAt:    updateAt,
		CreateAt:    r.CreateAt,
	}
}
