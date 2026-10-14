package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestTicketServiceSubmit(t *testing.T) {
	tickets := newFakeTicketRepo()
	logs := &fakeLogRepo{}
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)

	id, err := svc.Submit(context.Background(), SubmitRequest{
		ClusterID: 1,
		Database:  "halo",
		Subject:   "s",
		Content:   "select 1",
		UserID:    9,
	})
	require.NoError(t, err)
	require.NotZero(t, id)

	tk, err := tickets.GetByID(context.Background(), int64(id))
	require.NoError(t, err)
	require.Equal(t, domain.TicketStatusWaitingForVld, tk.Status)
	require.NotEmpty(t, tk.UUID)

	require.Len(t, logs.entries, 1)
	require.Equal(t, "ticket.create|ticket:"+tk.UUID+"|", logs.entries[0].Operation)
}

func TestTicketServiceSubmitLogFailure(t *testing.T) {
	tickets := newFakeTicketRepo()
	logs := &fakeLogRepo{err: errors.New("boom")}
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)

	_, err := svc.Submit(context.Background(), SubmitRequest{Content: "x", UserID: 1})
	require.Error(t, err)
}

func TestTicketServiceTransitions(t *testing.T) {
	ctx := context.Background()
	tickets := newFakeTicketRepo()
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: &fakeLogRepo{}}}, tickets)

	_, err := tickets.Create(ctx, &domain.Ticket{
		UUID:    "u1",
		Status:  domain.TicketStatusWaitingForMrv,
		Version: 1,
	})
	require.NoError(t, err)

	require.NoError(t, svc.Approve(ctx, 1, 42))
	got, err := tickets.GetByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, domain.TicketStatusLgtm, got.Status)
	require.Equal(t, uint(42), got.ReviewerID)
	require.Equal(t, 2, got.Version)
}

func TestTicketServiceRejectsIllegalTransition(t *testing.T) {
	ctx := context.Background()
	tickets := newFakeTicketRepo()
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: &fakeLogRepo{}}}, tickets)

	_, err := tickets.Create(ctx, &domain.Ticket{
		UUID:    "u2",
		Status:  domain.TicketStatusClosed,
		Version: 1,
	})
	require.NoError(t, err)

	err = svc.Approve(ctx, 1, 1)
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestCanTransition(t *testing.T) {
	require.True(t, canTransition(domain.TicketStatusWaitingForMrv, domain.TicketStatusLgtm))
	require.True(t, canTransition(domain.TicketStatusLgtm, domain.TicketStatusDone))
	require.False(t, canTransition(domain.TicketStatusClosed, domain.TicketStatusLgtm))
	require.False(t, canTransition(domain.TicketStatusDone, domain.TicketStatusLgtm))
}
