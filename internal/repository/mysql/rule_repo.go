package mysqlrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.RuleRepository = (*RuleRepo)(nil)

type RuleRepo struct {
	q *Queries
}

func NewRuleRepo(q *Queries) *RuleRepo {
	return &RuleRepo{q: q}
}

func ruleFromRow(r MmRule) *domain.Rule {
	return &domain.Rule{
		Name:        r.Name,
		UUID:        r.Uuid,
		Group:       r.Group,
		Description: r.Description,
		Level:       domain.RuleLevel(r.Level),
		VldrGroup:   r.VldrGroup,
		Operator:    r.Operator,
		Values:      r.Values,
		Bitwise:     r.Bitwise,
		Func:        r.Func,
		Message:     r.Message,
		Element:     r.Element,
		Version:     int(r.Version),
		UpdateAt:    uintFromNull(r.UpdateAt),
		CreateAt:    uint(r.CreateAt),
	}
}

func (r *RuleRepo) Create(ctx context.Context, rule *domain.Rule) error {
	return wrap("rule.Create", r.q.CreateRule(ctx, CreateRuleParams{
		Name:        rule.Name,
		Uuid:        rule.UUID,
		Group:       rule.Group,
		Description: rule.Description,
		Level:       uint8(rule.Level),
		VldrGroup:   rule.VldrGroup,
		Operator:    rule.Operator,
		Values:      rule.Values,
		Bitwise:     rule.Bitwise,
		Func:        rule.Func,
		Message:     rule.Message,
		Element:     rule.Element,
		Version:     uint32(rule.Version),
		UpdateAt:    nullFromUint(rule.UpdateAt),
		CreateAt:    uint32(rule.CreateAt),
	}))
}

func (r *RuleRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Rule, error) {
	row, err := r.q.GetRuleByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("rule.GetByUUID", err)
	}
	return ruleFromRow(row), nil
}

// GetByName resolves the row by its natural key.
func (r *RuleRepo) GetByName(ctx context.Context, name string) (*domain.Rule, error) {
	row, err := r.q.GetRuleByName(ctx, name)
	if err != nil {
		return nil, wrap("rule.GetByName", err)
	}
	return ruleFromRow(row), nil
}

func (r *RuleRepo) Update(ctx context.Context, rule *domain.Rule) error {
	return wrap("rule.Update", r.q.UpdateRule(ctx, UpdateRuleParams{
		Group:       rule.Group,
		Description: rule.Description,
		Level:       uint8(rule.Level),
		VldrGroup:   rule.VldrGroup,
		Operator:    rule.Operator,
		Values:      rule.Values,
		Bitwise:     rule.Bitwise,
		Func:        rule.Func,
		Message:     rule.Message,
		Element:     rule.Element,
		Version:     uint32(rule.Version),
		UpdateAt:    nullFromUint(rule.UpdateAt),
		Uuid:        rule.UUID,
	}))
}

func (r *RuleRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("rule.DeleteByUUID", r.q.DeleteRuleByUUID(ctx, uuid))
}

func (r *RuleRepo) List(ctx context.Context, cursor string, limit int) ([]*domain.Rule, error) {
	rows, err := r.q.ListRules(ctx, ListRulesParams{
		Uuid:  cursor,
		Limit: int32(limit),
	})
	if err != nil {
		return nil, wrap("rule.List", err)
	}
	out := make([]*domain.Rule, 0, len(rows))
	for _, row := range rows {
		out = append(out, ruleFromRow(row))
	}
	return out, nil
}
