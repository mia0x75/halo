package postgresrepo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestRuleHitRepo(t *testing.T) {
	repo := NewRuleHitRepo(setupTestDB(t))
	ctx := context.Background()

	hit := &domain.RuleHit{
		UUID:     newUUID(),
		TicketID: uint(newID()),
		Sequence: 1,
		RuleName: "DEL-L2-001",
		Level:    1,
		Message:  "禁止没有WHERE的删除",
	}
	require.NoError(t, repo.Create(ctx, hit))

	got, err := repo.ListByTicket(ctx, hit.TicketID)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "DEL-L2-001", got[0].RuleName)
	require.Equal(t, uint8(1), got[0].Level)
}

func TestRuleHitRepoAggregates(t *testing.T) {
	repo := NewRuleHitRepo(setupTestDB(t))
	ctx := context.Background()

	ticketID := uint(newID())
	for i, rule := range []string{"DEL-L2-001", "DEL-L2-001", "UPD-L2-001"} {
		require.NoError(t, repo.Create(ctx, &domain.RuleHit{
			UUID:     newUUID(),
			TicketID: ticketID,
			Sequence: uint16(i + 1),
			RuleName: rule,
			Level:    1,
			Message:  "命中",
		}))
	}

	from := uint(time.Now().Add(-time.Hour).Unix())
	to := uint(time.Now().Add(time.Hour).Unix())

	byLevel, err := repo.CountByLevel(ctx, from, to)
	require.NoError(t, err)
	require.NotEmpty(t, byLevel)
	require.Equal(t, uint(3), byLevel[0].Count)

	byRule, err := repo.CountByRule(ctx, from, to, 10)
	require.NoError(t, err)
	require.Len(t, byRule, 2)
	require.Equal(t, "DEL-L2-001", byRule[0].Rule)
	require.Equal(t, uint(2), byRule[0].Count)
}
