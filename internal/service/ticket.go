package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mia0x75/halo/internal/domain"
)

// ticketTransitions is the allowed workflow graph over mm_tickets.status. It
// uses the real nine-state enum from domain/ticket.go; Closed is terminal.
var ticketTransitions = map[domain.TicketStatus][]domain.TicketStatus{
	domain.TicketStatusWaitingForVld: {
		domain.TicketStatusVldFailure,
		domain.TicketStatusVldWarning,
		domain.TicketStatusWaitingForMrv,
		domain.TicketStatusClosed,
	},
	domain.TicketStatusVldFailure: {
		domain.TicketStatusWaitingForVld,
		domain.TicketStatusClosed,
	},
	domain.TicketStatusVldWarning: {
		domain.TicketStatusWaitingForMrv,
		domain.TicketStatusClosed,
	},
	domain.TicketStatusWaitingForMrv: {
		domain.TicketStatusMrvFailure,
		domain.TicketStatusLgtm,
		domain.TicketStatusClosed,
	},
	domain.TicketStatusMrvFailure: {
		domain.TicketStatusWaitingForVld,
		domain.TicketStatusClosed,
	},
	domain.TicketStatusLgtm: {
		domain.TicketStatusDone,
		domain.TicketStatusExecFailure,
		domain.TicketStatusClosed,
	},
	domain.TicketStatusDone:        {domain.TicketStatusClosed},
	domain.TicketStatusExecFailure: {domain.TicketStatusClosed},
	domain.TicketStatusClosed:      {},
}

// canTransition reports whether the workflow permits from -> to.
func canTransition(from, to domain.TicketStatus) bool {
	for _, allowed := range ticketTransitions[from] {
		if allowed == to {
			return true
		}
	}
	return false
}

// TicketService owns the ticket workflow: submission, review outcome and the
// terminal transitions driven by execution. Every status change is audited in
// the same transaction as the row update.
type TicketService struct {
	tx      TxManager
	tickets domain.TicketRepository
}

func NewTicketService(tx TxManager, tickets domain.TicketRepository) *TicketService {
	return &TicketService{tx: tx, tickets: tickets}
}

// Get returns one ticket by id (read path, no transaction needed).
func (s *TicketService) Get(ctx context.Context, id uint) (*domain.Ticket, error) {
	return s.tickets.GetByID(ctx, int64(id))
}

// GetByUUID returns one ticket by its stable uuid (read path).
func (s *TicketService) GetByUUID(ctx context.Context, uuid string) (*domain.Ticket, error) {
	return s.tickets.GetByUUID(ctx, uuid)
}

// List returns a keyset page of tickets.
func (s *TicketService) List(ctx context.Context, cursor int64, limit int) ([]*domain.Ticket, error) {
	return s.tickets.List(ctx, cursor, limit)
}

// SubmitRequest is the input for submitting a new ticket.
type SubmitRequest struct {
	ClusterID uint
	Database  string
	Subject   string
	Content   string
	UserID    uint
	CronID    *uint
}

// Submit creates a ticket in the initial validation state and audits it.
func (s *TicketService) Submit(ctx context.Context, req SubmitRequest) (uint, error) {
	var id uint
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		tk := &domain.Ticket{
			UUID:      uuid.NewString(),
			ClusterID: req.ClusterID,
			Database:  req.Database,
			Subject:   req.Subject,
			Content:   req.Content,
			Status:    domain.TicketStatusWaitingForVld,
			UserID:    req.UserID,
			CronID:    req.CronID,
			Version:   1,
			CreateAt:  uint(time.Now().Unix()),
		}
		newID, err := repos.Ticket.Create(ctx, tk)
		if err != nil {
			return err
		}
		id = uint(newID)
		return NewLogService(repos.Log).Log(ctx, LogEntry{
			UserID:    req.UserID,
			Operation: "ticket.create",
			Target:    "ticket:" + tk.UUID,
		})
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Approve moves a manually-reviewed ticket to Lgtm.
func (s *TicketService) Approve(ctx context.Context, ticketID, reviewerID uint) error {
	return s.transition(ctx, ticketID, domain.TicketStatusLgtm, reviewerID, "ticket.approve", "", func(tk *domain.Ticket) {
		tk.ReviewerID = reviewerID
	})
}

// Reject records a manual-review rejection and its reason.
func (s *TicketService) Reject(ctx context.Context, ticketID, reviewerID uint, reason string) error {
	return s.transition(ctx, ticketID, domain.TicketStatusMrvFailure, reviewerID, "ticket.reject", reason, func(tk *domain.Ticket) {
		tk.ReviewerID = reviewerID
	})
}

// Close moves a ticket to the terminal state.
func (s *TicketService) Close(ctx context.Context, ticketID, userID uint) error {
	return s.transition(ctx, ticketID, domain.TicketStatusClosed, userID, "ticket.close", "", nil)
}

// MarkExecuted is called by AgentService after a successful run.
func (s *TicketService) MarkExecuted(ctx context.Context, ticketID uint) error {
	return s.transition(ctx, ticketID, domain.TicketStatusDone, 0, "ticket.done", "", nil)
}

// MarkFailed is called by AgentService after a failed run.
func (s *TicketService) MarkFailed(ctx context.Context, ticketID uint, reason string) error {
	return s.transition(ctx, ticketID, domain.TicketStatusExecFailure, 0, "ticket.exec_failure", reason, nil)
}

// Remove deletes a ticket and its statements after checking ownership and that
// the workflow allows it. It returns the removed ticket so the caller can fire
// its audit event.
func (s *TicketService) Remove(ctx context.Context, ticketUUID string, userID uint) (*domain.Ticket, error) {
	var removed *domain.Ticket
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		tk, err := repos.Ticket.GetByUUID(ctx, ticketUUID)
		if err != nil {
			return err
		}
		if tk.UserID != userID {
			return fmt.Errorf("%w: only the author may remove ticket %s", domain.ErrForbidden, ticketUUID)
		}
		switch tk.Status {
		case domain.TicketStatusClosed, domain.TicketStatusDone, domain.TicketStatusExecFailure:
			return fmt.Errorf("%w: ticket %s is not removable in status %d", domain.ErrInvalidInput, ticketUUID, tk.Status)
		}
		if err := repos.Statement.DeleteByTicket(ctx, tk.TicketID); err != nil {
			return err
		}
		if err := repos.Ticket.Delete(ctx, int64(tk.TicketID)); err != nil {
			return err
		}
		removed = tk
		return NewLogService(repos.Log).Log(ctx, LogEntry{UserID: userID, Operation: "ticket.remove", Target: "ticket:" + tk.UUID})
	})
	if err != nil {
		return nil, err
	}
	return removed, nil
}

// transition loads the ticket, validates the edge and applies it with its audit
// row in one transaction.
func (s *TicketService) transition(ctx context.Context, ticketID uint, to domain.TicketStatus, actor uint, op, detail string, apply func(*domain.Ticket)) error {
	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		tk, err := repos.Ticket.GetByID(ctx, int64(ticketID))
		if err != nil {
			return err
		}
		if !canTransition(tk.Status, to) {
			return fmt.Errorf("%w: ticket %d: %d -> %d not allowed", domain.ErrInvalidInput, ticketID, tk.Status, to)
		}
		tk.Status = to
		tk.Version++
		now := uint(time.Now().Unix())
		tk.UpdateAt = &now
		if apply != nil {
			apply(tk)
		}
		if err := repos.Ticket.Update(ctx, tk); err != nil {
			return err
		}
		return NewLogService(repos.Log).Log(ctx, LogEntry{
			UserID:    actor,
			Operation: op,
			Target:    "ticket:" + tk.UUID,
			Detail:    detail,
		})
	})
}
