package mysqlrepo

import (
	"context"
	"time"

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
		UserID:   uint32(m.UserID),
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

func (r *MessageRepo) ListByUser(
	ctx context.Context,
	userID uint,
	beforeID uint,
	unreadOnly bool,
	category string,
	keyword string,
	limit int,
) ([]*domain.Message, error) {
	rows, err := r.q.ListMessagesByUser(ctx, ListMessagesByUserParams{
		UserID:     uint32(userID),
		BeforeID:   uint32(beforeID),
		UnreadOnly: boolToInt(unreadOnly),
		Category:   category,
		Keyword:    keyword,
		Limit:      clampInt32(limit),
	})
	if err != nil {
		return nil, wrap("message.ListByUser", err)
	}
	out := make([]*domain.Message, 0, len(rows))
	for _, row := range rows {
		out = append(out, messageFromRow(row))
	}
	return out, nil
}

func (r *MessageRepo) CountByUser(ctx context.Context, userID uint, unreadOnly bool, category string, keyword string) (uint, error) {
	count, err := r.q.CountMessagesByUser(ctx, CountMessagesByUserParams{
		UserID:     uint32(userID),
		UnreadOnly: boolToInt(unreadOnly),
		Category:   category,
		Keyword:    keyword,
	})
	if err != nil {
		return 0, wrap("message.CountByUser", err)
	}
	return uint(count), nil
}

func (r *MessageRepo) Delete(ctx context.Context, uuid string, userID uint) error {
	rows, err := r.q.DeleteMessageByUUID(ctx, DeleteMessageByUUIDParams{
		Uuid:   uuid,
		UserID: uint32(userID),
	})
	if err != nil {
		return wrap("message.Delete", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// DeleteBefore removes the messages created before the given time.
func (r *MessageRepo) DeleteBefore(ctx context.Context, before uint) (int64, error) {
	rows, err := r.q.DeleteMessagesBefore(ctx, time.Unix(int64(before), 0).UTC())
	if err != nil {
		return 0, wrap("message.DeleteBefore", err)
	}
	return rows, nil
}

func (r *MessageRepo) CountUnread(ctx context.Context, userID uint) (uint, error) {
	count, err := r.q.CountUnreadMessages(ctx, uint32(userID))
	if err != nil {
		return 0, wrap("message.CountUnread", err)
	}
	return uint(count), nil
}

func (r *MessageRepo) MarkAllRead(ctx context.Context, userID uint) (int64, error) {
	rows, err := r.q.MarkAllMessagesRead(ctx, uint32(userID))
	if err != nil {
		return 0, wrap("message.MarkAllRead", err)
	}
	return rows, nil
}

func (r *MessageRepo) MarkRead(ctx context.Context, uuid string, userID uint) error {
	return wrap("message.MarkRead", r.q.MarkMessageRead(ctx, MarkMessageReadParams{
		Uuid:   uuid,
		UserID: uint32(userID),
	}))
}

// boolToInt renders a flag for the SQL layer, which compares it to 0.
func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
