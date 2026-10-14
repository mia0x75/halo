package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/models"
)

// CommentFromDomain converts a domain entity to its legacy model.
func CommentFromDomain(c *domain.Comment) *models.Comment {
	if c == nil {
		return nil
	}
	return &models.Comment{
		CommentID: c.CommentID,
		UUID:      c.UUID,
		Content:   c.Content,
		TicketID:  c.TicketID,
		UserID:    c.UserID,
		Version:   c.Version,
		UpdateAt:  deref(c.UpdateAt),
		CreateAt:  c.CreateAt,
	}
}

// CommentToDomain converts a legacy model to its domain entity.
func CommentToDomain(c *models.Comment) *domain.Comment {
	if c == nil {
		return nil
	}
	return &domain.Comment{
		CommentID: c.CommentID,
		UUID:      c.UUID,
		Content:   c.Content,
		TicketID:  c.TicketID,
		UserID:    c.UserID,
		Version:   c.Version,
		UpdateAt:  uintPtr(c.UpdateAt),
		CreateAt:  c.CreateAt,
	}
}
