package service

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/domain"
)

// MessageService delivers the in-app notifications.
type MessageService struct {
	messages domain.MessageRepository
}

func NewMessageService(messages domain.MessageRepository) *MessageService {
	return &MessageService{messages: messages}
}

// Send stores one message for its recipient.
func (s *MessageService) Send(ctx context.Context, m *domain.Message) error {
	return s.messages.Create(ctx, m)
}

// MessageFilter narrows one message page.
type MessageFilter struct {
	// BeforeID is the exclusive upper bound on message_id, 0 means no bound.
	BeforeID uint
	// UnreadOnly keeps only the messages that were never read.
	UnreadOnly bool
	// Category filters by category, empty matches every category.
	Category string
	// Keyword matches the title or the content, empty matches everything.
	Keyword string
	// Limit is the page size.
	Limit int
}

// ListByUser returns one page of the messages of one user, newest first.
func (s *MessageService) ListByUser(ctx context.Context, userID uint, filter MessageFilter) ([]*domain.Message, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	return s.messages.ListByUser(
		ctx, userID, filter.BeforeID, filter.UnreadOnly, filter.Category, filter.Keyword, filter.Limit)
}

// CountByUser returns how many messages match the filter.
func (s *MessageService) CountByUser(ctx context.Context, userID uint, filter MessageFilter) (uint, error) {
	return s.messages.CountByUser(ctx, userID, filter.UnreadOnly, filter.Category, filter.Keyword)
}

// Delete removes one message of the given user.
func (s *MessageService) Delete(ctx context.Context, uuid string, userID uint) error {
	return s.messages.Delete(ctx, uuid, userID)
}

// CountUnread returns the number of unread messages of one user.
func (s *MessageService) CountUnread(ctx context.Context, userID uint) (uint, error) {
	return s.messages.CountUnread(ctx, userID)
}

// MarkRead marks one message of the given user as read.
func (s *MessageService) MarkRead(ctx context.Context, uuid string, userID uint) error {
	return s.messages.MarkRead(ctx, uuid, userID)
}

// MarkAllRead marks every unread message of the user as read.
func (s *MessageService) MarkAllRead(ctx context.Context, userID uint) (int64, error) {
	if userID == 0 {
		return 0, fmt.Errorf("mark all messages read: missing user")
	}
	return s.messages.MarkAllRead(ctx, userID)
}
