package mysqlrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.CommentRepository = (*CommentRepo)(nil)

type CommentRepo struct {
	q *Queries
}

func NewCommentRepo(q *Queries) *CommentRepo {
	return &CommentRepo{q: q}
}

func commentFromRow(r MmComment) *domain.Comment {
	return &domain.Comment{
		CommentID: uint(r.CommentID),
		UUID:      r.Uuid,
		Content:   r.Content,
		TicketID:  uint(r.TicketID),
		UserID:    uint(r.UserID),
		Version:   int(r.Version),
		UpdateAt:  timeFromNull(r.UpdateAt),
		CreateAt:  timeToUnix(r.CreateAt),
	}
}

func (r *CommentRepo) Create(ctx context.Context, c *domain.Comment) (int64, error) {
	id, err := r.q.CreateComment(ctx, CreateCommentParams{
		Uuid:     c.UUID,
		Content:  c.Content,
		TicketID: uint32(c.TicketID),
		UserID:   uint32(c.UserID),
		Version:  uint32(c.Version),
	})
	if err != nil {
		return 0, wrap("comment.Create", err)
	}
	return id, nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id int64) (*domain.Comment, error) {
	row, err := r.q.GetCommentByID(ctx, uint32(id))
	if err != nil {
		return nil, wrap("comment.GetByID", err)
	}
	return commentFromRow(row), nil
}

func (r *CommentRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Comment, error) {
	row, err := r.q.GetCommentByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("comment.GetByUUID", err)
	}
	return commentFromRow(row), nil
}

func (r *CommentRepo) Update(ctx context.Context, c *domain.Comment) error {
	err := r.q.UpdateComment(ctx, UpdateCommentParams{
		Uuid:      c.UUID,
		Content:   c.Content,
		TicketID:  uint32(c.TicketID),
		UserID:    uint32(c.UserID),
		Version:   uint32(c.Version),
		CommentID: uint32(c.CommentID),
	})
	return wrap("comment.Update", err)
}

func (r *CommentRepo) Delete(ctx context.Context, id int64) error {
	return wrap("comment.Delete", r.q.DeleteCommentByID(ctx, uint32(id)))
}

func (r *CommentRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("comment.DeleteByUUID", r.q.DeleteCommentByUUID(ctx, uuid))
}

func (r *CommentRepo) List(ctx context.Context, cursor int64, limit int) ([]*domain.Comment, error) {
	rows, err := r.q.ListComments(ctx, ListCommentsParams{
		CommentID: uint32(cursor),
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, wrap("comment.List", err)
	}
	out := make([]*domain.Comment, 0, len(rows))
	for _, row := range rows {
		out = append(out, commentFromRow(row))
	}
	return out, nil
}

// ListByTicket returns every comment on a ticket, oldest first.
func (r *CommentRepo) ListByTicket(ctx context.Context, ticketID uint) ([]*domain.Comment, error) {
	rows, err := r.q.ListCommentsByTicket(ctx, uint32(ticketID))
	if err != nil {
		return nil, wrap("comment.ListByTicket", err)
	}
	out := make([]*domain.Comment, 0, len(rows))
	for _, row := range rows {
		out = append(out, commentFromRow(row))
	}
	return out, nil
}
