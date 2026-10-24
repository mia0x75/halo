package service

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/agent"
	"github.com/mia0x75/halo/internal/domain"
)

// AgentService executes an approved ticket against its target cluster. The
// executor call is blocking network IO, so it deliberately runs outside any
// transaction; only the status changes and audit rows are transactional.
type AgentService struct {
	tx       TxManager
	executor agent.Executor
}

func NewAgentService(tx TxManager, executor agent.Executor) *AgentService {
	return &AgentService{tx: tx, executor: executor}
}

// Execute runs every statement of an Lgtm ticket and records the outcome.
func (s *AgentService) Execute(ctx context.Context, ticketID uint) error {
	var cluster *domain.Cluster
	var stmts []*domain.Statement
	// The executor acts on behalf of the submitter, so the audit trail keeps a
	// human owner even though the transition itself is automatic.
	var ownerID uint

	err := s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		tk, err := repos.Ticket.GetByID(ctx, int64(ticketID))
		if err != nil {
			return err
		}
		ownerID = tk.UserID
		if tk.Status != domain.TicketStatusLgtm {
			return fmt.Errorf("%w: ticket %d is not approved (status %d)", domain.ErrInvalidInput, ticketID, tk.Status)
		}
		cluster, err = repos.Cluster.GetByID(ctx, int64(tk.ClusterID))
		if err != nil {
			return err
		}
		stmts, err = listStatementsByTicket(ctx, repos.Statement, ticketID)
		return err
	})
	if err != nil {
		return err
	}

	type outcome struct {
		stmt *domain.Statement
		res  *agent.ExecResult
		err  error
	}
	outcomes := make([]outcome, 0, len(stmts))
	failed := false
	password := string(cluster.Password)
	for _, st := range stmts {
		res, err := s.executor.Execute(ctx, cluster, password, st.Content)
		if err != nil || (res != nil && res.Error != "") {
			failed = true
		}
		outcomes = append(outcomes, outcome{stmt: st, res: res, err: err})
	}

	return s.tx.WithTx(ctx, func(ctx context.Context, repos Repos) error {
		log := NewLogService(repos.Log)
		for _, o := range outcomes {
			st := o.stmt
			results := ""
			switch {
			case o.err != nil:
				results = o.err.Error()
			case o.res != nil && o.res.Error != "":
				results = o.res.Error
			case o.res != nil:
				results = fmt.Sprintf("ok: %d rows affected", o.res.RowsAffected)
			}
			st.Results = &results
			if o.res != nil {
				rows := uint(o.res.RowsAffected)
				st.RowsAffected = &rows
			}
			st.Version++
			if err := repos.Statement.Update(ctx, st); err != nil {
				return err
			}
			_ = log.Log(ctx, LogEntry{
				Operation: "statement.execute",
				Target:    fmt.Sprintf("statement:%d:%d", st.TicketID, st.Sequence),
				UserID:    ownerID,
			})
		}

		tk, err := repos.Ticket.GetByID(ctx, int64(ticketID))
		if err != nil {
			return err
		}
		to := domain.TicketStatusDone
		if failed {
			to = domain.TicketStatusExecFailure
		}
		if !canTransition(tk.Status, to) {
			return fmt.Errorf("%w: ticket %d: %d -> %d not allowed", domain.ErrInvalidInput, ticketID, tk.Status, to)
		}
		tk.Status = to
		tk.Version++
		if err := repos.Ticket.Update(ctx, tk); err != nil {
			return err
		}
		return log.Log(ctx, LogEntry{
			Operation: "ticket.execute",
			Target:    "ticket:" + tk.UUID,
			UserID:    ownerID,
		})
	})
}
