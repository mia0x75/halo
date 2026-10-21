package postgresrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.MessageRepository = (*MessageRepo)(nil)

type MessageRepo struct {
	q *Queries
}

func NewMessageRepo(q *Queries) *MessageRepo {
	return &MessageRepo{q: q}
}

func messageFromRow(r MmMessage) *domain.Message {
	return &domain.Message{
		MessageID: uint(r.MessageID),
		UUID:      r.Uuid,
		UserID:    uint(r.UserID),
		Category:  r.Category,
		Title:     r.Title,
		Content:   r.Content,
		Link:      r.Link,
		ReadAt:    timeFromNull(r.ReadAt),
		CreateAt:  timeToUnix(r.CreateAt),
	}
}

func (r *MessageRepo) Create(ctx context.Context, m *domain.Message) error {
	return wrap("message.Create", r.q.CreateMessage(ctx, CreateMessageParams{
		Uuid:     m.UUID,
		UserID:   int64(m.UserID),
		Category: m.Category,
		Title:    m.Title,
		Content:  m.Content,
		Link:     m.Link,
	}))
}

func (r *MessageRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Message, error) {
	row, err := r.q.GetMessageByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("message.GetByUUID", err)
	}
	return messageFromRow(row), nil
}

func (r *MessageRepo) ListByUser(ctx context.Context, userID uint, unreadOnly bool, limit int) ([]*domain.Message, error) {
	var rows []MmMessage
	var err error
	if unreadOnly {
		rows, err = r.q.ListUnreadMessagesByUser(ctx, ListUnreadMessagesByUserParams{
			UserID: int64(userID),
			Limit:  clampInt32(limit),
		})
	} else {
		rows, err = r.q.ListMessagesByUser(ctx, ListMessagesByUserParams{
			UserID: int64(userID),
			Limit:  clampInt32(limit),
		})
	}
	if err != nil {
		return nil, wrap("message.ListByUser", err)
	}
	out := make([]*domain.Message, 0, len(rows))
	for _, row := range rows {
		out = append(out, messageFromRow(row))
	}
	return out, nil
}

func (r *MessageRepo) CountUnread(ctx context.Context, userID uint) (uint, error) {
	count, err := r.q.CountUnreadMessages(ctx, int64(userID))
	if err != nil {
		return 0, wrap("message.CountUnread", err)
	}
	return uint(count), nil
}

func (r *MessageRepo) MarkAllRead(ctx context.Context, userID uint) (int64, error) {
	rows, err := r.q.MarkAllMessagesRead(ctx, int64(userID))
	if err != nil {
		return 0, wrap("message.MarkAllRead", err)
	}
	return rows, nil
}

func (r *MessageRepo) MarkRead(ctx context.Context, uuid string, userID uint) error {
	return wrap("message.MarkRead", r.q.MarkMessageRead(ctx, MarkMessageReadParams{
		Uuid:   uuid,
		UserID: int64(userID),
	}))
}
