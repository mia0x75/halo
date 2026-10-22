package resolvers

import (
	"context"
	"time"

	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/transport/graphql/gqlapi"
)

// RuleHits 查询规则命中明细
func (r *queryRootResolver) RuleHits(ctx context.Context, ticketID *uint, ruleName *string, level *uint8, first *int) (L []*models.RuleHit, err error) {
	limit := 100
	if first != nil && *first > 0 {
		limit = *first
	}

	hits, err := r.app.Services.RuleHit.List(ctx, ticketIDValue(ticketID), limit)
	if err != nil {
		return nil, err
	}

	L = make([]*models.RuleHit, 0, len(hits))
	for _, hit := range hits {
		if ruleName != nil && hit.RuleName != *ruleName {
			continue
		}
		if level != nil && hit.Level != *level {
			continue
		}
		L = append(L, bridge.RuleHitFromDomain(hit))
	}
	return L, nil
}

// RuleHitLevels 按严重程度统计命中次数
func (r *queryRootResolver) RuleHitLevels(ctx context.Context, from *uint, to *uint) (L []*gqlapi.RuleHitCount, err error) {
	start, end := hitWindow(from, to)
	counts, err := r.app.Services.RuleHit.CountByLevel(ctx, start, end)
	if err != nil {
		return nil, err
	}
	L = make([]*gqlapi.RuleHitCount, 0, len(counts))
	for _, count := range counts {
		L = append(L, &gqlapi.RuleHitCount{Level: count.Level, Count: uint(count.Count)})
	}
	return L, nil
}

// RuleHitRules 按规则统计命中次数
func (r *queryRootResolver) RuleHitRules(ctx context.Context, from *uint, to *uint, first *int) (L []*gqlapi.RuleHitCount, err error) {
	start, end := hitWindow(from, to)
	limit := 10
	if first != nil && *first > 0 {
		limit = *first
	}
	counts, err := r.app.Services.RuleHit.CountByRule(ctx, start, end, limit)
	if err != nil {
		return nil, err
	}
	L = make([]*gqlapi.RuleHitCount, 0, len(counts))
	for _, count := range counts {
		L = append(L, &gqlapi.RuleHitCount{RuleName: &count.Rule, Count: uint(count.Count)})
	}
	return L, nil
}

// ticketIDValue returns the requested ticket or 0, which means "all tickets".
func ticketIDValue(ticketID *uint) uint {
	if ticketID == nil {
		return 0
	}
	return *ticketID
}

// hitWindow returns the requested time window, defaulting to the last 7 days.
func hitWindow(from *uint, to *uint) (uint, uint) {
	end := uint(time.Now().Unix())
	start := end - 7*24*3600
	if to != nil && *to > 0 {
		end = *to
	}
	if from != nil && *from > 0 {
		start = *from
	}
	return start, end
}
