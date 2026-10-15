package postgresrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestStatementRepo(t *testing.T) {
	repo := NewStatementRepo(setupTestDB(t))
	ctx := context.Background()

	s := &domain.Statement{
		TicketID: newID(),
		Sequence: 1,
		UUID:     newUUID(),
		Content:  "alter table t add column c int",
		Type:     domain.StatementTypeAlterTable,
		Status:   domain.StatementStatusPending,
		Report:   "ok",
		Version:  1,
		CreateAt: nowUnix(),
	}
	require.NoError(t, repo.Create(ctx, s))

	got, err := repo.GetByUUID(ctx, s.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.StatementTypeAlterTable, got.Type)
	require.Nil(t, got.Plan)
	require.Nil(t, got.Results)

	byKey, err := repo.GetByTicketSequence(ctx, s.TicketID, s.Sequence)
	require.NoError(t, err)
	require.Equal(t, s.UUID, byKey.UUID)

	plan := "plan"
	rows := uint(5)
	updated := nowUnix()
	s.Status = domain.StatementStatusVldWarning
	s.Plan = &plan
	s.RowsAffected = &rows
	s.Version = 2
	s.UpdateAt = &updated
	require.NoError(t, repo.Update(ctx, s))

	got, err = repo.GetByUUID(ctx, s.UUID)
	require.NoError(t, err)
	require.Equal(t, domain.StatementStatusVldWarning, got.Status)
	require.NotNil(t, got.Plan)
	require.Equal(t, "plan", *got.Plan)

	list, err := repo.List(ctx, s.UUID+"~", 100)
	require.NoError(t, err)
	require.True(t, hasStatement(list, s.UUID))

	require.NoError(t, repo.DeleteByTicketSequence(ctx, s.TicketID, s.Sequence))
	_, err = repo.GetByUUID(ctx, s.UUID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func hasStatement(list []*domain.Statement, uuid string) bool {
	for _, it := range list {
		if it.UUID == uuid {
			return true
		}
	}
	return false
}
