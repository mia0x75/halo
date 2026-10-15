package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
)

// TicketFromDomain converts a domain entity to its legacy model. The nullable
// cron id is boxed into the legacy sql.NullInt64.
func TicketFromDomain(t *domain.Ticket) *models.Ticket {
	if t == nil {
		return nil
	}
	return &models.Ticket{
		TicketID:   t.TicketID,
		UUID:       t.UUID,
		ClusterID:  t.ClusterID,
		Database:   t.Database,
		Subject:    t.Subject,
		Content:    t.Content,
		Status:     uint8(t.Status),
		UserID:     t.UserID,
		ReviewerID: t.ReviewerID,
		CronID:     nullInt64(t.CronID),
		Version:    t.Version,
		UpdateAt:   deref(t.UpdateAt),
		CreateAt:   t.CreateAt,
	}
}

// TicketToDomain converts a legacy model to its domain entity.
func TicketToDomain(t *models.Ticket) *domain.Ticket {
	if t == nil {
		return nil
	}
	return &domain.Ticket{
		TicketID:   t.TicketID,
		UUID:       t.UUID,
		ClusterID:  t.ClusterID,
		Database:   t.Database,
		Subject:    t.Subject,
		Content:    t.Content,
		Status:     domain.TicketStatus(t.Status),
		UserID:     t.UserID,
		ReviewerID: t.ReviewerID,
		CronID:     uintFromNullInt64(t.CronID),
		Version:    t.Version,
		UpdateAt:   uintPtr(t.UpdateAt),
		CreateAt:   t.CreateAt,
	}
}
