package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/analyzer"
	"github.com/mia0x75/halo/internal/domain"
)

func TestStatementServiceSubmit(t *testing.T) {
	ctx := context.Background()
	statements := newFakeStatementRepo()
	a, err := analyzer.New(domain.ClusterTypeMySQL, nil)
	require.NoError(t, err)

	svc := NewStatementService(&fakeTx{repos: Repos{Statement: statements, RuleHit: newFakeRuleHitRepo()}}, statements, nil, nil, newFakeRuleHitRepo(), nil, a)

	first, err := svc.Submit(ctx, 1, "select 1")
	require.NoError(t, err)
	require.Equal(t, uint16(1), first.Sequence)
	require.Equal(t, domain.StatementTypeSelect, first.Type)
	require.Equal(t, domain.StatementStatusWaitingForMrv, first.Status)

	second, err := svc.Submit(ctx, 1, "insert into t values (1)")
	require.NoError(t, err)
	require.Equal(t, uint16(2), second.Sequence)
	require.Equal(t, domain.StatementTypeInsert, second.Type)

	list, err := svc.ListByTicket(ctx, 1)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, first.UUID, list[0].UUID)
	require.Equal(t, second.UUID, list[1].UUID)
}

func TestStatementServiceUpdateResult(t *testing.T) {
	ctx := context.Background()
	statements := newFakeStatementRepo()
	svc := NewStatementService(&fakeTx{repos: Repos{Statement: statements, RuleHit: newFakeRuleHitRepo()}}, statements, nil, nil, newFakeRuleHitRepo(), nil, nil)

	_, err := svc.Submit(ctx, 1, "select 1")
	require.NoError(t, err)

	err = svc.UpdateResult(ctx, 1, 1, StatementResult{Results: "ok", RowsAffected: 3})
	require.NoError(t, err)

	got, err := statements.GetByTicketSequence(ctx, 1, 1)
	require.NoError(t, err)
	require.NotNil(t, got.Results)
	require.Equal(t, "ok", *got.Results)
	require.NotNil(t, got.RowsAffected)
	require.Equal(t, uint(3), *got.RowsAffected)
}

func TestStatementServiceSubmitWithoutAnalyzer(t *testing.T) {
	ctx := context.Background()
	statements := newFakeStatementRepo()
	svc := NewStatementService(&fakeTx{repos: Repos{Statement: statements, RuleHit: newFakeRuleHitRepo()}}, statements, nil, nil, newFakeRuleHitRepo(), nil, nil)

	st, err := svc.Submit(ctx, 7, "whatever")
	require.NoError(t, err)
	require.Equal(t, domain.StatementTypeSelect, st.Type)
	require.Equal(t, domain.StatementStatusPending, st.Status)
}

func TestStatementServiceParseFailure(t *testing.T) {
	ctx := context.Background()
	statements := newFakeStatementRepo()
	a, err := analyzer.New(domain.ClusterTypeMySQL, nil)
	require.NoError(t, err)

	svc := NewStatementService(&fakeTx{repos: Repos{Statement: statements, RuleHit: newFakeRuleHitRepo()}}, statements, nil, nil, newFakeRuleHitRepo(), nil, a)
	_, err = svc.Submit(ctx, 1, "@@@ not sql @@@")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}
