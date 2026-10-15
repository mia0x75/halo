package service

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

// ReviewService runs the automated rule engine over a ticket's SQL and advances
// the ticket through the validation states. Manual approval stays in
// TicketService.Approve; this service only reports whether approval is safe.
type ReviewService struct {
	tx        TxManager
	tickets   domain.TicketRepository
	validator domain.Validator
}

func NewReviewService(tx TxManager, tickets domain.TicketRepository, validator domain.Validator) *ReviewService {
	return &ReviewService{tx: tx, tickets: tickets, validator: validator}
}

// ReviewResult is the outcome of an automated review.
type ReviewResult struct {
	Violations []domain.Violation
	CanApprove bool
}

// Review validates the ticket content and moves a freshly submitted ticket to
// WaitingForMrv, VldWarning or VldFailure depending on the rule hits.
func (s *ReviewService) Review(ctx context.Context, ticketID, reviewerID uint) (*ReviewResult, error) {
	var result *ReviewResult
	err := s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		tk, err := repos.Ticket.GetByID(ctx, int64(ticketID))
		if err != nil {
			return err
		}

		var violations []domain.Violation
		if s.validator != nil {
			violations, err = s.validator.Validate(ctx, domain.ValidateRequest{SQL: tk.Content})
			if err != nil {
				return err
			}
		}

		canApprove := true
		hasWarning := false
		for _, v := range violations {
			switch v.Level {
			case uint8(domain.RuleLevelError):
				canApprove = false
			case uint8(domain.RuleLevelWarning):
				hasWarning = true
			}
		}
		result = &ReviewResult{Violations: violations, CanApprove: canApprove}

		if tk.Status == domain.TicketStatusWaitingForVld {
			to := domain.TicketStatusWaitingForMrv
			switch {
			case !canApprove:
				to = domain.TicketStatusVldFailure
			case hasWarning:
				to = domain.TicketStatusVldWarning
			}
			if canTransition(tk.Status, to) {
				tk.Status = to
				tk.Version++
				if err := repos.Ticket.Update(ctx, tk); err != nil {
					return err
				}
			}
		}

		return NewLogService(repos.Log).Log(ctx, LogEntry{
			UserID:    reviewerID,
			Operation: "ticket.review",
			Target:    "ticket:" + tk.UUID,
		})
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
