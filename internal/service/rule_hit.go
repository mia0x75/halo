package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

// RuleHitService keeps the rule hit detail: one row per hit, so the dashboard
// can group hits by severity or by rule.
type RuleHitService struct {
	hits domain.RuleHitRepository
}

func NewRuleHitService(hits domain.RuleHitRepository) *RuleHitService {
	return &RuleHitService{hits: hits}
}

// Record writes the hits of one statement.
func (s *RuleHitService) Record(ctx context.Context, hits []*domain.RuleHit) error {
	if len(hits) == 0 {
		return nil
	}
	return s.hits.CreateBatch(ctx, hits)
}

// ListByTicket returns the hits of one ticket.
func (s *RuleHitService) ListByTicket(ctx context.Context, ticketID uint) ([]*domain.RuleHit, error) {
	return s.hits.ListByTicket(ctx, ticketID)
}

// List returns the most recent hits; a ticket id of 0 lists every ticket.
func (s *RuleHitService) List(ctx context.Context, ticketID uint, limit int) ([]*domain.RuleHit, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.hits.List(ctx, ticketID, limit)
}

// CountByLevel counts the hits per severity inside a time window.
func (s *RuleHitService) CountByLevel(ctx context.Context, from, to uint) ([]domain.RuleHitCount, error) {
	return s.hits.CountByLevel(ctx, from, to)
}

// CountByRule counts the hits per rule inside a time window, most frequent
// first.
func (s *RuleHitService) CountByRule(ctx context.Context, from, to uint, limit int) ([]domain.RuleHitCount, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.hits.CountByRule(ctx, from, to, limit)
}
