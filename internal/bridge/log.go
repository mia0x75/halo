package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/transport/graphql/models"
)

// LogFromDomain converts a domain entity to its legacy model.
func LogFromDomain(l *domain.Log) *models.Log {
	if l == nil {
		return nil
	}
	return &models.Log{
		LogID:     int(l.LogID),
		UUID:      l.UUID,
		UserID:    l.UserID,
		Operation: l.Operation,
		Version:   l.Version,
		CreateAt:  l.CreateAt,
	}
}

// LogToDomain converts a legacy model to its domain entity.
