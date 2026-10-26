package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestTicketServiceUpdateAndClearCron(t *testing.T) {
	ctx := context.Background()
	tickets, logs := ticketFixture(t, domain.TicketStatusLgtm)
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)

	cronID := uint(77)
	require.NoError(t, svc.UpdateCron(ctx, "t1", cronID))
	tk, err := tickets.GetByID(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, tk.CronID)
	require.Equal(t, cronID, *tk.CronID)
	// 状态保持 LGTM，工单只是排期，没有发生状态流转。
	require.Equal(t, domain.TicketStatusLgtm, tk.Status)

	require.NoError(t, svc.ClearCron(ctx, "t1"))
	tk, err = tickets.GetByID(ctx, 1)
	require.NoError(t, err)
	require.Nil(t, tk.CronID)
	// 取消预约后工单回到待上线执行队列，可重新预约或立即执行。
	require.Equal(t, domain.TicketStatusLgtm, tk.Status)

	require.ErrorIs(t, svc.ClearCron(ctx, "missing"), domain.ErrNotFound)
}

func TestTicketServiceListExecutionTickets(t *testing.T) {
	ctx := context.Background()
	tickets := newFakeTicketRepo()
	seed := func(id uint, uuid string, reviewer uint, status domain.TicketStatus, cronID *uint) {
		_, err := tickets.Create(ctx, &domain.Ticket{
			TicketID:   id,
			UUID:       uuid,
			ReviewerID: reviewer,
			Status:     status,
			CronID:     cronID,
			Version:    1,
		})
		require.NoError(t, err)
	}
	scheduled := uint(9)
	seed(1, "mine", 7, domain.TicketStatusLgtm, nil)
	seed(2, "mine-scheduled", 7, domain.TicketStatusLgtm, &scheduled)
	seed(3, "theirs", 8, domain.TicketStatusLgtm, nil)
	seed(4, "not-lgtm", 7, domain.TicketStatusWaitingForMrv, nil)

	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets}}, tickets)

	// 审核人只看分派给自己的 LGTM 工单
	mine, err := svc.ListExecutionTickets(ctx, 7, false, 0, 10)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"mine", "mine-scheduled"}, ticketUUIDs(mine))

	// reviewerID=0 是管理员视角，看全部 LGTM
	all, err := svc.ListExecutionTickets(ctx, 0, false, 0, 10)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"mine", "mine-scheduled", "theirs"}, ticketUUIDs(all))

	// 已预约列表只包含带 cron 的工单
	scheduledList, err := svc.ListExecutionTickets(ctx, 0, true, 0, 10)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"mine-scheduled"}, ticketUUIDs(scheduledList))

	total, err := svc.CountExecutionTickets(ctx, 7, false)
	require.NoError(t, err)
	require.Equal(t, int64(2), total)

	total, err = svc.CountExecutionTickets(ctx, 0, true)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
}

func ticketUUIDs(tickets []*domain.Ticket) []string {
	out := make([]string, 0, len(tickets))
	for _, tk := range tickets {
		out = append(out, tk.UUID)
	}
	return out
}
