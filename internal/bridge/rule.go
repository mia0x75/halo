package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/transport/graphql/models"
)

// RuleFromDomain converts a domain entity to its legacy model. Field sets match;
// only the numeric widths differ.
func RuleFromDomain(r *domain.Rule) *models.Rule {
	if r == nil {
		return nil
	}
	return &models.Rule{
		Name:        r.Name,
		UUID:        r.UUID,
		Group:       r.Group,
		Description: r.Description,
		Level:       uint8(r.Level),
		VldrGroup:   uint16(r.VldrGroup),
		Operator:    r.Operator,
		Values:      r.Values,
		Bitwise:     r.Bitwise,
		Func:        r.Func,
		Message:     r.Message,
		Element:     r.Element,
		Version:     r.Version,
		UpdateAt:    deref(r.UpdateAt),
		CreateAt:    r.CreateAt,
	}
}

// RuleToDomain converts a legacy model to its domain entity.
func RuleToDomain(r *models.Rule) *domain.Rule {
	if r == nil {
		return nil
	}
	return &domain.Rule{
		Name:        r.Name,
		UUID:        r.UUID,
		Group:       r.Group,
		Description: r.Description,
		Level:       domain.RuleLevel(r.Level),
		VldrGroup:   uint8(r.VldrGroup),
		Operator:    r.Operator,
		Values:      r.Values,
		Bitwise:     r.Bitwise,
		Func:        r.Func,
		Message:     r.Message,
		Element:     r.Element,
		Version:     r.Version,
		UpdateAt:    uintPtr(r.UpdateAt),
		CreateAt:    r.CreateAt,
	}
}
