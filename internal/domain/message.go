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
	// ListByUser returns one page of the messages of a user, newest first.
	//
	// beforeID is the exclusive upper bound on message_id: the keyset cursor
	// the API hands out. unreadOnly keeps only the messages that were never
	// read, and an empty category matches every category.
	ListByUser(ctx context.Context, userID uint, beforeID uint, unreadOnly bool, category string, keyword string, limit int) ([]*Message, error)
	// CountByUser returns how many messages match the same filters.
	CountByUser(ctx context.Context, userID uint, unreadOnly bool, category string, keyword string) (uint, error)
	// CountUnread returns the number of unread messages of one user.
	CountUnread(ctx context.Context, userID uint) (uint, error)
	// Delete removes one message of the given user.
	Delete(ctx context.Context, uuid string, userID uint) error
	// MarkRead marks one message of the given user as read.
	MarkRead(ctx context.Context, uuid string, userID uint) error
	// MarkAllRead marks every unread message of the user as read and returns
	// how many rows were updated.
	MarkAllRead(ctx context.Context, userID uint) (int64, error)
}
