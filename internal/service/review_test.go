package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestReviewService(t *testing.T) {
	cases := []struct {
		name       string
		violations []domain.Violation
		want       domain.TicketStatus
		canApprove bool
	}{
		{
			name:       "clean",
			violations: nil,
			want:       domain.TicketStatusWaitingForMrv,
			canApprove: true,
		},
		{
			name:       "warning",
			violations: []domain.Violation{{Message: "w", Level: uint8(domain.RuleLevelWarning)}},
			want:       domain.TicketStatusVldWarning,
			canApprove: true,
		},
		{
			name:       "error",
			violations: []domain.Violation{{Message: "e", Level: uint8(domain.RuleLevelError)}},
			want:       domain.TicketStatusVldFailure,
			canApprove: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			tickets := newFakeTicketRepo()
			_, err := tickets.Create(ctx, &domain.Ticket{
				UUID:    "u",
				Status:  domain.TicketStatusWaitingForVld,
				Content: "select 1",
				Version: 1,
			})
			require.NoError(t, err)

			svc := NewReviewService(
				&fakeTx{repos: Repos{Ticket: tickets, Log: &fakeLogRepo{}, Cluster: newFakeClusterRepo()}},
				tickets,
				&fakeValidator{violations: tc.violations},
			)

			res, err := svc.Review(ctx, 1, 5)
			require.NoError(t, err)
			require.Equal(t, tc.canApprove, res.CanApprove)
			require.Len(t, res.Violations, len(tc.violations))

			got, err := tickets.GetByID(ctx, 1)
			require.NoError(t, err)
			require.Equal(t, tc.want, got.Status)
		})
	}
}

func TestReviewServiceNilValidator(t *testing.T) {
	ctx := context.Background()
	tickets := newFakeTicketRepo()
	_, err := tickets.Create(ctx, &domain.Ticket{UUID: "u", Status: domain.TicketStatusWaitingForVld, Version: 1})
	require.NoError(t, err)

	svc := NewReviewService(&fakeTx{repos: Repos{Ticket: tickets, Log: &fakeLogRepo{}, Cluster: newFakeClusterRepo()}}, tickets, nil)
	res, err := svc.Review(ctx, 1, 1)
	require.NoError(t, err)
	require.True(t, res.CanApprove)
	require.Empty(t, res.Violations)
}
