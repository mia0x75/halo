package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mia0x75/halo/internal/domain"
)

func TestTicketServiceReadPaths(t *testing.T) {
	ctx := context.Background()
	tickets := newFakeTicketRepo()
	_, err := tickets.Create(ctx, &domain.Ticket{TicketID: 1, UUID: "t1", Status: domain.TicketStatusLgtm, Version: 1})
	require.NoError(t, err)

	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets}}, tickets)

	got, err := svc.Get(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, "t1", got.UUID)

	byUUID, err := svc.GetByUUID(ctx, "t1")
	require.NoError(t, err)
	require.Equal(t, uint(1), byUUID.TicketID)

	_, err = svc.GetByCronID(ctx, 1)
	require.ErrorIs(t, err, domain.ErrNotFound)

	_, err = svc.List(ctx, 0, 10)
	require.NoError(t, err)
	_, err = svc.ListByUser(ctx, 5)
	require.NoError(t, err)
}

func TestTicketServiceReject(t *testing.T) {
	ctx := context.Background()
	tickets, logs := ticketFixture(t, domain.TicketStatusWaitingForMrv)
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)

	require.NoError(t, svc.Reject(ctx, 1, 42, "needs work"))
	tk, err := tickets.GetByID(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, domain.TicketStatusMrvFailure, tk.Status)
	require.Equal(t, uint(42), tk.ReviewerID)
	require.Contains(t, logs.entries[0].Operation, "ticket.reject")
	require.Contains(t, logs.entries[0].Operation, "needs work")
}

func TestTicketServiceMarkExecutedAttributesSubmitter(t *testing.T) {
	ctx := context.Background()

	tickets, logs := ticketFixture(t, domain.TicketStatusLgtm)
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)
	require.NoError(t, svc.MarkExecuted(ctx, 1))

	// The executor has no account, so the audit row belongs to the submitter.
	require.Len(t, logs.entries, 1)
	require.Equal(t, uint(5), logs.entries[0].UserID)
}

func TestLogServiceWritesWithoutOperator(t *testing.T) {
	// Only the wizard path skips auditing now, and it bypasses the service;
	// the service itself still records whatever it is handed.
	logs := &fakeLogRepo{}
	svc := NewLogService(logs)
	require.NoError(t, svc.Log(context.Background(), LogEntry{UserID: 0, Operation: "x"}))
	require.Len(t, logs.entries, 1)
}

func TestTicketServiceCloseAndMark(t *testing.T) {
	ctx := context.Background()

	tickets, logs := ticketFixture(t, domain.TicketStatusWaitingForVld)
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)
	require.NoError(t, svc.Close(ctx, 1, 9))
	tk, _ := tickets.GetByID(ctx, 1)
	require.Equal(t, domain.TicketStatusClosed, tk.Status)

	tickets, logs = ticketFixture(t, domain.TicketStatusLgtm)
	svc = NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)
	require.NoError(t, svc.MarkExecuted(ctx, 1))
	tk, _ = tickets.GetByID(ctx, 1)
	require.Equal(t, domain.TicketStatusDone, tk.Status)
	require.Equal(t, "ticket.done", logs.entries[0].Operation[:11])

	tickets, logs = ticketFixture(t, domain.TicketStatusLgtm)
	svc = NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)
	require.NoError(t, svc.MarkFailed(ctx, 1, "boom"))
	tk, _ = tickets.GetByID(ctx, 1)
	require.Equal(t, domain.TicketStatusExecFailure, tk.Status)
}

func TestTicketServiceRemove(t *testing.T) {
	ctx := context.Background()

	// Only the author may remove.
	tickets, logs := ticketFixture(t, domain.TicketStatusWaitingForVld)
	statements := newFakeStatementRepo()
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Statement: statements, Log: logs}}, tickets)
	_, err := svc.Remove(ctx, "t1", 999)
	require.ErrorIs(t, err, domain.ErrForbidden)

	// Terminal tickets are not removable.
	tickets, logs = ticketFixture(t, domain.TicketStatusDone)
	svc = NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Statement: newFakeStatementRepo(), Log: logs}}, tickets)
	_, err = svc.Remove(ctx, "t1", 5)
	require.ErrorIs(t, err, domain.ErrInvalidInput)

	// Happy path removes statements and the ticket.
	tickets, logs = ticketFixture(t, domain.TicketStatusWaitingForVld)
	statements = newFakeStatementRepo()
	require.NoError(t, statements.Create(ctx, &domain.Statement{TicketID: 1, Sequence: 1, UUID: "s1"}))
	svc = NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Statement: statements, Log: logs}}, tickets)
	removed, err := svc.Remove(ctx, "t1", 5)
	require.NoError(t, err)
	require.Equal(t, "t1", removed.UUID)
	_, err = tickets.GetByID(ctx, 1)
	require.ErrorIs(t, err, domain.ErrNotFound)
	_, err = statements.GetByTicketSequence(ctx, 1, 1)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestTicketServicePatchStatus(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name    string
		from    domain.TicketStatus
		target  domain.TicketStatus
		wantTo  domain.TicketStatus
		wantErr bool
		kind    TransitionKind
	}{
		{name: "terminal", from: domain.TicketStatusDone, target: domain.TicketStatusClosed, wantErr: true, kind: TransitionTerminal},
		{name: "close from vld", from: domain.TicketStatusWaitingForVld, target: domain.TicketStatusClosed, wantErr: true, kind: TransitionClose},
		{name: "approve is not allowed here", from: domain.TicketStatusVldWarning, target: domain.TicketStatusLgtm, wantErr: true, kind: TransitionInvalidTarget},
		{name: "reject is not allowed here", from: domain.TicketStatusWaitingForMrv, target: domain.TicketStatusMrvFailure, wantErr: true, kind: TransitionInvalidTarget},
		{name: "invalid target", from: domain.TicketStatusWaitingForVld, target: domain.TicketStatusWaitingForVld, wantErr: true, kind: TransitionInvalidTarget},
		{name: "warning to close", from: domain.TicketStatusVldWarning, target: domain.TicketStatusClosed, wantTo: domain.TicketStatusClosed},
		{name: "mrv to close", from: domain.TicketStatusWaitingForMrv, target: domain.TicketStatusClosed, wantTo: domain.TicketStatusClosed},
		{name: "vld failure to close", from: domain.TicketStatusVldFailure, target: domain.TicketStatusClosed, wantTo: domain.TicketStatusClosed},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tickets, logs := ticketFixture(t, tc.from)
			svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)

			// 5 是工单发起人，见 ticketFixture。
			got, err := svc.PatchStatus(ctx, "t1", tc.target, 5, false)
			if tc.wantErr {
				var te *TransitionError
				require.ErrorAs(t, err, &te)
				require.Equal(t, tc.kind, te.Kind)
				require.ErrorIs(t, err, domain.ErrInvalidInput)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantTo, got.Status)
			require.Equal(t, 2, got.Version)
		})
	}
}

func TestTicketServicePatchStatusRejectsOtherUser(t *testing.T) {
	ctx := context.Background()
	tickets, logs := ticketFixture(t, domain.TicketStatusWaitingForMrv)
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)

	_, err := svc.PatchStatus(ctx, "t1", domain.TicketStatusClosed, 99, false)
	require.ErrorIs(t, err, domain.ErrForbidden)
}

func TestTicketServicePatchStatusAllowsAdmin(t *testing.T) {
	ctx := context.Background()
	tickets, logs := ticketFixture(t, domain.TicketStatusWaitingForMrv)
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)

	got, err := svc.PatchStatus(ctx, "t1", domain.TicketStatusClosed, 99, true)
	require.NoError(t, err)
	require.Equal(t, domain.TicketStatusClosed, got.Status)
}

func TestTransitionErrorMessageAndIs(t *testing.T) {
	err := &TransitionError{Current: domain.TicketStatusDone, Kind: TransitionTerminal}
	require.Equal(t, "ticket transition rejected", err.Error())
	require.ErrorIs(t, err, domain.ErrInvalidInput)
	require.False(t, errors.Is(err, domain.ErrNotFound))
}

func TestTicketServiceUpdateCron(t *testing.T) {
	ctx := context.Background()
	tickets, logs := ticketFixture(t, domain.TicketStatusWaitingForVld)
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Log: logs}}, tickets)

	require.NoError(t, svc.UpdateCron(ctx, "t1", 77))
	tk, err := tickets.GetByID(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, tk.CronID)
	require.Equal(t, uint(77), *tk.CronID)
	require.Equal(t, 2, tk.Version)
}

func TestTicketServiceCreateFillsUUIDs(t *testing.T) {
	ctx := context.Background()
	tickets := newFakeTicketRepo()
	statements := newFakeStatementRepo()
	logs := &fakeLogRepo{}
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Statement: statements, Log: logs}}, tickets)

	st := &domain.Statement{Sequence: 1, Content: "select 1"}
	id, err := svc.Create(ctx, &domain.Ticket{ClusterID: 2, UserID: 8}, []*domain.Statement{st})
	require.NoError(t, err)
	require.NotZero(t, id)

	tk, err := tickets.GetByID(ctx, int64(id))
	require.NoError(t, err)
	require.NotEmpty(t, tk.UUID)
	require.NotEmpty(t, st.UUID)
	require.Equal(t, tk.TicketID, st.TicketID)
	require.Equal(t, "ticket.create|ticket:"+tk.UUID+"|", logs.entries[0].Operation)
}

func TestTicketServiceUpdateReplacesStatements(t *testing.T) {
	ctx := context.Background()
	tickets := newFakeTicketRepo()
	_, err := tickets.Create(ctx, &domain.Ticket{TicketID: 1, UUID: "t1", Status: domain.TicketStatusWaitingForVld, Version: 1})
	require.NoError(t, err)
	statements := newFakeStatementRepo()
	require.NoError(t, statements.Create(ctx, &domain.Statement{TicketID: 1, Sequence: 1, UUID: "old"}))
	logs := &fakeLogRepo{}
	svc := NewTicketService(&fakeTx{repos: Repos{Ticket: tickets, Statement: statements, Log: logs}}, tickets)

	st := &domain.Statement{Sequence: 1, Content: "select 1"}
	require.NoError(t, svc.Update(ctx, &domain.Ticket{TicketID: 1, UUID: "t1", UserID: 8}, []*domain.Statement{st}))

	_, err = statements.GetByUUID(ctx, "old")
	require.ErrorIs(t, err, domain.ErrNotFound)
	got, err := statements.GetByTicketSequence(ctx, 1, 1)
	require.NoError(t, err)
	require.NotEmpty(t, got.UUID)
	require.Equal(t, "select 1", got.Content)
	require.Equal(t, "ticket.update|ticket:t1|", logs.entries[0].Operation)
}

// ticketFixture seeds one ticket in the given status owned by user 5.
func ticketFixture(t *testing.T, status domain.TicketStatus) (*fakeTicketRepo, *fakeLogRepo) {
	t.Helper()
	tickets := newFakeTicketRepo()
	_, err := tickets.Create(context.Background(), &domain.Ticket{
		TicketID: 1,
		UUID:     "t1",
		UserID:   5,
		Status:   status,
		Version:  1,
	})
	require.NoError(t, err)
	return tickets, &fakeLogRepo{}
}
