package domain

import "context"

// Message is a row of mm_messages: an in-app notification for one user.
type Message struct {
	MessageID uint
	UUID      string
	UserID    uint
	Category  string
	Title     string
	Content   string
	Link      string
	ReadAt    *uint
	CreateAt  uint
}

// MessageRepository stores the in-app notifications.
type MessageRepository interface {
	Create(ctx context.Context, m *Message) error
	GetByUUID(ctx context.Context, uuid string) (*Message, error)
	// ListByUser returns the messages of one user, newest first. When
	// unreadOnly is true only the messages that were never read are returned.
	ListByUser(ctx context.Context, userID uint, unreadOnly bool, limit int) ([]*Message, error)
	// CountUnread returns the number of unread messages of one user.
	CountUnread(ctx context.Context, userID uint) (uint, error)
	// MarkRead marks one message of the given user as read.
	MarkRead(ctx context.Context, uuid string, userID uint) error
	// MarkAllRead marks every unread message of the user as read and returns
	// how many rows were updated.
	MarkAllRead(ctx context.Context, userID uint) (int64, error)
}
