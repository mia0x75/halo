package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// OptionFromDomain converts a domain entity to its legacy model.
func OptionFromDomain(o *domain.Option) *models.Option {
	if o == nil {
		return nil
	}
	return &models.Option{
		Name:        o.Name,
		UUID:        o.UUID,
		Parameter:   o.Parameter,
		Description: o.Description,
		Element:     o.Element,
		Writable:    o.Writable,
		GroupName:   o.GroupName,
		SortOrder:   o.SortOrder,
		DependsOn:   o.DependsOn,
		MinValue:    o.MinValue,
		MaxValue:    o.MaxValue,
		Version:     o.Version,
		UpdateAt:    deref(o.UpdateAt),
		CreateAt:    o.CreateAt,
	}
}

// OptionToDomain converts a legacy model to its domain entity.
func OptionToDomain(o *models.Option) *domain.Option {
	if o == nil {
		return nil
	}
	return &domain.Option{
		Name:        o.Name,
		UUID:        o.UUID,
		Parameter:   o.Parameter,
		Description: o.Description,
		Element:     o.Element,
		Writable:    o.Writable,
		GroupName:   o.GroupName,
		SortOrder:   o.SortOrder,
		DependsOn:   o.DependsOn,
		MinValue:    o.MinValue,
		MaxValue:    o.MaxValue,
		Version:     o.Version,
		UpdateAt:    uintPtr(o.UpdateAt),
		CreateAt:    o.CreateAt,
	}
}
