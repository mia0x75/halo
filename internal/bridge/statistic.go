package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// StatisticFromDomain converts a domain entity to its legacy model.
func StatisticFromDomain(s *domain.Statistic) *models.Statistic {
	if s == nil {
		return nil
	}
	return &models.Statistic{
		Group:    s.Group,
		Key:      s.Key,
		UUID:     s.UUID,
		Value:    s.Value,
		Version:  s.Version,
		UpdateAt: deref(s.UpdateAt),
		CreateAt: s.CreateAt,
	}
}

// StatisticToDomain converts a legacy model to its domain entity.
