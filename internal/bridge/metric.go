package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// MetricFromDomain converts a domain entity to its model.
func MetricFromDomain(m *domain.Metric) *models.Metric {
	if m == nil {
		return nil
	}
	return &models.Metric{
		MetricID: m.MetricID,
		UUID:     m.UUID,
		Name:     m.Name,
		Scope:    m.Scope,
		ScopeID:  m.ScopeID,
		Type:     m.Type,
		Bucket:   m.Bucket,
		Value:    m.Value,
		Version:  m.Version,
		UpdateAt: deref(m.UpdateAt),
		CreateAt: m.CreateAt,
	}
}

// MetricToDomain converts a model to its domain entity.
func MetricToDomain(m *models.Metric) *domain.Metric {
	if m == nil {
		return nil
	}
	updateAt := m.UpdateAt
	return &domain.Metric{
		MetricID: m.MetricID,
		UUID:     m.UUID,
		Name:     m.Name,
		Scope:    m.Scope,
		ScopeID:  m.ScopeID,
		Type:     m.Type,
		Bucket:   m.Bucket,
		Value:    m.Value,
		Version:  m.Version,
		UpdateAt: &updateAt,
		CreateAt: m.CreateAt,
	}
}
