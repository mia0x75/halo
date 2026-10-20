package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestListPendingReviews(t *testing.T) {
	ctx := context.Background()
	repo := &fakeTicketRepo{}
	repo.reviewer = []*domain.Ticket{
		{TicketID: 1, Status: domain.TicketStatusWaitingForMrv},
		{TicketID: 2, Status: domain.TicketStatusWaitingForVld},
		{TicketID: 3, Status: domain.TicketStatusWaitingForMrv},
		{TicketID: 4, Status: domain.TicketStatusDone},
	}
	svc := NewTicketService(&fakeTx{}, repo)

	pending, err := svc.ListPendingReviews(ctx, 7)
	require.NoError(t, err)
	require.Len(t, pending, 2)
	require.Equal(t, uint(1), pending[0].TicketID)
	require.Equal(t, uint(3), pending[1].TicketID)
}
