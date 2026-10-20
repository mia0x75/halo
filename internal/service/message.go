package service

import (
	"context"

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

// ListByUser returns the messages of one user, newest first.
func (s *MessageService) ListByUser(ctx context.Context, userID uint, unreadOnly bool, limit int) ([]*domain.Message, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.messages.ListByUser(ctx, userID, unreadOnly, limit)
}

// CountUnread returns the number of unread messages of one user.
func (s *MessageService) CountUnread(ctx context.Context, userID uint) (uint, error) {
	return s.messages.CountUnread(ctx, userID)
}

// MarkRead marks one message of the given user as read.
func (s *MessageService) MarkRead(ctx context.Context, uuid string, userID uint) error {
	return s.messages.MarkRead(ctx, uuid, userID)
}
