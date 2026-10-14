package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/agent"
	"github.com/mia0x75/halo/internal/domain"
)

func setupExecutableTicket(t *testing.T) (*fakeTicketRepo, *fakeStatementRepo, *fakeClusterRepo, *fakeLogRepo) {
	t.Helper()
	ctx := context.Background()

	clusters := newFakeClusterRepo()
	_, err := clusters.Create(ctx, &domain.Cluster{ClusterID: 1, UUID: "c1", Password: []byte("pw")})
	require.NoError(t, err)

	tickets := newFakeTicketRepo()
	_, err = tickets.Create(ctx, &domain.Ticket{TicketID: 1, UUID: "t1", ClusterID: 1, Status: domain.TicketStatusLgtm, Version: 1})
	require.NoError(t, err)

	statements := newFakeStatementRepo()
	err = statements.Create(ctx, &domain.Statement{TicketID: 1, Sequence: 1, UUID: "s1", Content: "alter table t add c int", Version: 1})
	require.NoError(t, err)

	return tickets, statements, clusters, &fakeLogRepo{}
}

func TestAgentServiceExecuteSuccess(t *testing.T) {
	ctx := context.Background()
	tickets, statements, clusters, logs := setupExecutableTicket(t)

	exec := &fakeExecutor{}
	svc := NewAgentService(&fakeTx{repos: Repos{Ticket: tickets, Statement: statements, Cluster: clusters, Log: logs}}, exec)

	require.NoError(t, svc.Execute(ctx, 1))
	require.Equal(t, 1, exec.calls)

	tk, err := tickets.GetByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, domain.TicketStatusDone, tk.Status)

	st, err := statements.GetByTicketSequence(ctx, 1, 1)
	require.NoError(t, err)
	require.NotNil(t, st.Results)
	require.NotNil(t, st.RowsAffected)
	require.Equal(t, uint(1), *st.RowsAffected)

	// statement.execute + ticket.execute
	require.Len(t, logs.entries, 2)
}

func TestAgentServiceExecuteFailure(t *testing.T) {
	ctx := context.Background()
	tickets, statements, clusters, logs := setupExecutableTicket(t)

	exec := &fakeExecutor{results: map[string]*agent.ExecResult{
		"alter table t add c int": {Error: "syntax error"},
	}}
	svc := NewAgentService(&fakeTx{repos: Repos{Ticket: tickets, Statement: statements, Cluster: clusters, Log: logs}}, exec)

	require.NoError(t, svc.Execute(ctx, 1))

	tk, err := tickets.GetByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, domain.TicketStatusExecFailure, tk.Status)
}

func TestAgentServiceRejectsUnapprovedTicket(t *testing.T) {
	ctx := context.Background()
	tickets := newFakeTicketRepo()
	_, err := tickets.Create(ctx, &domain.Ticket{TicketID: 1, UUID: "t1", Status: domain.TicketStatusWaitingForVld, Version: 1})
	require.NoError(t, err)

	svc := NewAgentService(&fakeTx{repos: Repos{Ticket: tickets, Log: &fakeLogRepo{}}}, &fakeExecutor{})
	err = svc.Execute(ctx, 1)
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}
