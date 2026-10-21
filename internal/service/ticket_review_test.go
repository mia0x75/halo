package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func reviewTicket(t *testing.T) (*fakeTicketRepo, *fakeLogRepo, *TicketService) {
	t.Helper()
	tickets := newFakeTicketRepo()
	logs := &fakeLogRepo{}
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)

	_, err := tickets.Create(context.Background(), &domain.Ticket{
		UUID:       "tk-review",
		TicketID:   1,
		Subject:    "待审核工单",
		Status:     domain.TicketStatusWaitingForMrv,
		UserID:     9,
		ReviewerID: 5,
		Version:    1,
	})
	require.NoError(t, err)
	return tickets, logs, svc
}

func TestTicketServiceReviewApprove(t *testing.T) {
	_, logs, svc := reviewTicket(t)

	result, err := svc.Review(context.Background(), "tk-review", 5, ReviewApproved, "")
	require.NoError(t, err)
	require.Equal(t, ReviewApproved, result.Outcome)
	require.Equal(t, domain.TicketStatusLgtm, result.Ticket.Status)

	require.Len(t, logs.entries, 1)
	require.Equal(t, "ticket.approve|ticket:tk-review|", logs.entries[0].Operation)
}

func TestTicketServiceReviewRejectKeepsReason(t *testing.T) {
	_, logs, svc := reviewTicket(t)

	result, err := svc.Review(context.Background(), "tk-review", 5, ReviewRejected, "SQL 会全表扫描")
	require.NoError(t, err)
	require.Equal(t, domain.TicketStatusMrvFailure, result.Ticket.Status)

	require.Len(t, logs.entries, 1)
	require.Equal(t, "ticket.reject|ticket:tk-review|SQL 会全表扫描", logs.entries[0].Operation)
}

func TestTicketServiceReviewRejectsOtherReviewer(t *testing.T) {
	_, logs, svc := reviewTicket(t)

	_, err := svc.Review(context.Background(), "tk-review", 6, ReviewApproved, "")
	require.ErrorIs(t, err, domain.ErrForbidden)
	require.Empty(t, logs.entries)
}

func TestTicketServiceReviewRejectsWrongStatus(t *testing.T) {
	tickets := newFakeTicketRepo()
	logs := &fakeLogRepo{}
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)
	_, err := tickets.Create(context.Background(), &domain.Ticket{
		UUID:       "tk-done",
		TicketID:   1,
		Subject:    "已通过",
		Status:     domain.TicketStatusLgtm,
		UserID:     9,
		ReviewerID: 5,
		Version:    1,
	})
	require.NoError(t, err)

	_, err = svc.Review(context.Background(), "tk-done", 5, ReviewApproved, "")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
	require.Empty(t, logs.entries)
}
