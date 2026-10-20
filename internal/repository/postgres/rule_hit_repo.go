package postgresrepo

import (
	"context"
	"time"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.RuleHitRepository = (*RuleHitRepo)(nil)

type RuleHitRepo struct {
	q *Queries
}

func NewRuleHitRepo(q *Queries) *RuleHitRepo {
	return &RuleHitRepo{q: q}
}

func ruleHitFromRow(r MmViolation) *domain.RuleHit {
	return &domain.RuleHit{
		HitID:    uint(r.ViolationID),
		UUID:     r.Uuid,
		TicketID: uint(r.TicketID),
		Sequence: uint16(r.Sequence),
		RuleName: r.RuleName,
		Level:    uint8(r.Level),
		Message:  r.Message,
		CreateAt: uint(r.CreateAt.Unix()),
	}
}

func (r *RuleHitRepo) Create(ctx context.Context, hit *domain.RuleHit) error {
	return wrap("ruleHit.Create", r.q.CreateRuleHit(ctx, CreateRuleHitParams{
		Uuid:     hit.UUID,
		TicketID: int64(hit.TicketID),
		Sequence: int32(hit.Sequence),
		RuleName: hit.RuleName,
		Level:    int16(hit.Level),
		Message:  hit.Message,
	}))
}

// CreateBatch writes the hits one by one: the caller already runs inside a
// transaction, so a single multi-row statement would not buy anything.
func (r *RuleHitRepo) CreateBatch(ctx context.Context, hits []*domain.RuleHit) error {
	for _, hit := range hits {
		if err := r.Create(ctx, hit); err != nil {
			return err
		}
	}
	return nil
}

func (r *RuleHitRepo) ListByTicket(ctx context.Context, ticketID uint) ([]*domain.RuleHit, error) {
	rows, err := r.q.ListRuleHitsByTicket(ctx, int64(ticketID))
	if err != nil {
		return nil, wrap("ruleHit.ListByTicket", err)
	}
	out := make([]*domain.RuleHit, 0, len(rows))
	for _, row := range rows {
		out = append(out, ruleHitFromRow(row))
	}
	return out, nil
}

// List returns the most recent hits. A ticket id of 0 lists every ticket.
func (r *RuleHitRepo) List(ctx context.Context, ticketID uint, limit int) ([]*domain.RuleHit, error) {
	if ticketID == 0 {
		rows, err := r.q.ListRuleHits(ctx, clampInt32(limit))
		if err != nil {
			return nil, wrap("ruleHit.List", err)
		}
		out := make([]*domain.RuleHit, 0, len(rows))
		for _, row := range rows {
			out = append(out, ruleHitFromRow(row))
		}
		return out, nil
	}
	return r.ListByTicket(ctx, ticketID)
}

func (r *RuleHitRepo) CountByLevel(ctx context.Context, from, to uint) ([]domain.RuleHitCount, error) {
	rows, err := r.q.CountRuleHitsByLevel(ctx, CountRuleHitsByLevelParams{
		CreateAt:   time.Unix(int64(from), 0).UTC(),
		CreateAt_2: time.Unix(int64(to), 0).UTC(),
	})
	if err != nil {
		return nil, wrap("ruleHit.CountByLevel", err)
	}
	out := make([]domain.RuleHitCount, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.RuleHitCount{Level: uint8(row.Level), Count: uint(row.Count)})
	}
	return out, nil
}

func (r *RuleHitRepo) CountByRule(ctx context.Context, from, to uint, limit int) ([]domain.RuleHitCount, error) {
	rows, err := r.q.CountRuleHitsByRule(ctx, CountRuleHitsByRuleParams{
		CreateAt:   time.Unix(int64(from), 0).UTC(),
		CreateAt_2: time.Unix(int64(to), 0).UTC(),
		Limit:      clampInt32(limit),
	})
	if err != nil {
		return nil, wrap("ruleHit.CountByRule", err)
	}
	out := make([]domain.RuleHitCount, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.RuleHitCount{Rule: row.RuleName, Count: uint(row.Count)})
	}
	return out, nil
}
