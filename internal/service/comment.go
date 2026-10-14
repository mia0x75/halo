package service

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/domain"
)

type CommentService struct {
	tx       TxManager
	comments domain.CommentRepository
}

func NewCommentService(tx TxManager, comments domain.CommentRepository) *CommentService {
	return &CommentService{tx: tx, comments: comments}
}

func (s *CommentService) Get(ctx context.Context, id uint) (*domain.Comment, error) {
	return s.comments.GetByID(ctx, int64(id))
}

func (s *CommentService) GetByUUID(ctx context.Context, uuid string) (*domain.Comment, error) {
	return s.comments.GetByUUID(ctx, uuid)
}

func (s *CommentService) List(ctx context.Context, cursor int64, limit int) ([]*domain.Comment, error) {
	return s.comments.List(ctx, cursor, limit)
}

func (s *CommentService) Create(ctx context.Context, c *domain.Comment, operatorID uint) (uint, error) {
	var id uint
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		newID, err := repos.Comment.Create(ctx, c)
		if err != nil {
			return err
		}
		id = uint(newID)
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "comment.create", Target: "comment:" + c.UUID})
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *CommentService) Update(ctx context.Context, c *domain.Comment, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Comment.Update(ctx, c); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "comment.update", Target: "comment:" + c.UUID})
	})
}

func (s *CommentService) Delete(ctx context.Context, id uint, operatorID uint) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		if err := repos.Comment.Delete(ctx, int64(id)); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: operatorID, Operation: "comment.delete", Target: fmt.Sprintf("comment:%d", id)})
	})
}
