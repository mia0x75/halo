package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestReviewServiceValidatorError(t *testing.T) {
	ctx := context.Background()
	tickets, logs := ticketFixture(t, domain.TicketStatusWaitingForVld)
	svc := NewReviewService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets, &fakeValidator{err: errors.New("engine down")})

	_, err := svc.Review(ctx, 1, 5)
	require.Error(t, err)
}

func TestReviewServiceMissingTicket(t *testing.T) {
	ctx := context.Background()
	tickets := newFakeTicketRepo()
	svc := NewReviewService(&fakeTx{repos: Repos{Ticket: tickets, Log: &fakeLogRepo{}}}, tickets, &fakeValidator{})

	_, err := svc.Review(ctx, 404, 5)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

// TestReviewServiceKeepsNonInitialStatus checks that reviewing a ticket that has
// already left the initial state records the audit row but does not move it.
func TestReviewServiceKeepsNonInitialStatus(t *testing.T) {
	ctx := context.Background()
	tickets, logs := ticketFixture(t, domain.TicketStatusWaitingForMrv)
	svc := NewReviewService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets, &fakeValidator{})

	res, err := svc.Review(ctx, 1, 5)
	require.NoError(t, err)
	require.True(t, res.CanApprove)

	tk, err := tickets.GetByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, domain.TicketStatusWaitingForMrv, tk.Status)
	require.Equal(t, 1, tk.Version)
}
