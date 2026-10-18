package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

// TestAgentServiceExecutorError covers the transport failure branch: the
// executor returns an error instead of an ExecResult.
func TestAgentServiceExecutorError(t *testing.T) {
	ctx := context.Background()
	tickets, statements, clusters, logs := setupExecutableTicket(t)

	exec := &fakeExecutor{err: errors.New("dial timeout")}
	svc := NewAgentService(&fakeTx{repos: Repos{Ticket: tickets, Statement: statements, Cluster: clusters, Log: logs}}, exec)

	require.NoError(t, svc.Execute(ctx, 1))
	tk, err := tickets.GetByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, domain.TicketStatusExecFailure, tk.Status)

	st, err := statements.GetByTicketSequence(ctx, 1, 1)
	require.NoError(t, err)
	require.NotNil(t, st.Results)
	require.Equal(t, "dial timeout", *st.Results)
	require.Nil(t, st.RowsAffected)
}

// TestAgentServiceRejectsMissingCluster covers a target cluster that vanished
// between approval and execution.
func TestAgentServiceRejectsMissingCluster(t *testing.T) {
	ctx := context.Background()
	tickets := newFakeTicketRepo()
	_, err := tickets.Create(ctx, &domain.Ticket{TicketID: 1, UUID: "t1", ClusterID: 99, Status: domain.TicketStatusLgtm, Version: 1})
	require.NoError(t, err)

	svc := NewAgentService(&fakeTx{repos: Repos{Ticket: tickets, Statement: newFakeStatementRepo(), Cluster: newFakeClusterRepo(), Log: &fakeLogRepo{}}}, &fakeExecutor{})
	err = svc.Execute(ctx, 1)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestAgentServiceNoStatements(t *testing.T) {
	ctx := context.Background()
	tickets := newFakeTicketRepo()
	_, err := tickets.Create(ctx, &domain.Ticket{TicketID: 1, UUID: "t1", ClusterID: 1, Status: domain.TicketStatusLgtm, Version: 1})
	require.NoError(t, err)
	clusters := newFakeClusterRepo()
	_, err = clusters.Create(ctx, &domain.Cluster{ClusterID: 1, UUID: "c1"})
	require.NoError(t, err)
	logs := &fakeLogRepo{}
	svc := NewAgentService(&fakeTx{repos: Repos{Ticket: tickets, Statement: newFakeStatementRepo(), Cluster: clusters, Log: logs}}, &fakeExecutor{})

	require.NoError(t, svc.Execute(ctx, 1))
	tk, _ := tickets.GetByID(ctx, 1)
	require.Equal(t, domain.TicketStatusDone, tk.Status)
	require.Len(t, logs.entries, 1)
	require.Equal(t, "ticket.execute|ticket:t1|", logs.entries[0].Operation)
}
