package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/mia0x75/halo/internal/analyzer"
	"github.com/mia0x75/halo/internal/analyzer/ir"
	"github.com/mia0x75/halo/internal/domain"
)

// StatementService splits a ticket's SQL into statements, validates each one and
// persists the result. It is the only service that touches the analyzer.
type StatementService struct {
	statements domain.StatementRepository
	log        *LogService
	analyzer   *analyzer.Analyzer
}

func NewStatementService(statements domain.StatementRepository, log *LogService, a *analyzer.Analyzer) *StatementService {
	return &StatementService{statements: statements, log: log, analyzer: a}
}

// Submit parses, validates and stores one statement for a ticket. The sequence
// is the next free slot under the (ticket_id, sequence) primary key.
func (s *StatementService) Submit(ctx context.Context, ticketID uint, sql string) (*domain.Statement, error) {
	seq, err := s.nextSequence(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	st := &domain.Statement{
		TicketID: ticketID,
		Sequence: seq,
		UUID:     uuid.NewString(),
		Content:  sql,
		Type:     domain.StatementTypeSelect,
		Status:   domain.StatementStatusPending,
		Version:  1,
		CreateAt: uint(time.Now().Unix()),
	}

	if s.analyzer != nil {
		res, err := s.analyzer.Analyze(ctx, sql, nil)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
		}
		st.Type = statementTypeFromIR(res.Statement)
		st.Status = statementStatusFromViolations(res.Violations)
		st.Report = joinViolations(res.Violations)
	}

	if err := s.statements.Create(ctx, st); err != nil {
		return nil, err
	}
	if s.log != nil {
		_ = s.log.Log(ctx, LogEntry{
			Operation: "statement.create",
			Target:    fmt.Sprintf("statement:%d:%d", st.TicketID, st.Sequence),
		})
	}
	return st, nil
}

// ListByTicket returns the statements of one ticket in sequence order.
func (s *StatementService) ListByTicket(ctx context.Context, ticketID uint) ([]*domain.Statement, error) {
	return listStatementsByTicket(ctx, s.statements, ticketID)
}

// GetByUUID returns one statement by its stable uuid.
func (s *StatementService) GetByUUID(ctx context.Context, uuid string) (*domain.Statement, error) {
	return s.statements.GetByUUID(ctx, uuid)
}

// StatementResult is the outcome written back after execution.
type StatementResult struct {
	Results      string
	RowsAffected uint
}

// UpdateResult writes an execution result back onto a statement.
func (s *StatementService) UpdateResult(ctx context.Context, ticketID uint, sequence uint16, r StatementResult) error {
	st, err := s.statements.GetByTicketSequence(ctx, ticketID, sequence)
	if err != nil {
		return err
	}
	st.Results = &r.Results
	st.RowsAffected = &r.RowsAffected
	st.Version++
	if err := s.statements.Update(ctx, st); err != nil {
		return err
	}
	return nil
}

// nextSequence finds the first free sequence for a ticket by probing the
// composite key until the repository reports not-found.
func (s *StatementService) nextSequence(ctx context.Context, ticketID uint) (uint16, error) {
	for seq := uint16(1); ; seq++ {
		if seq == 0 {
			return 0, fmt.Errorf("%w: ticket %d has no free statement sequence", domain.ErrInvalidInput, ticketID)
		}
		_, err := s.statements.GetByTicketSequence(ctx, ticketID, seq)
		if errors.Is(err, domain.ErrNotFound) {
			return seq, nil
		}
		if err != nil {
			return 0, err
		}
	}
}

// listStatementsByTicket walks sequences from 1 until the first missing slot.
// Statements are numbered contiguously, so this terminates on the ticket's set.
func listStatementsByTicket(ctx context.Context, repo domain.StatementRepository, ticketID uint) ([]*domain.Statement, error) {
	var out []*domain.Statement
	for seq := uint16(1); ; seq++ {
		if seq == 0 {
			break
		}
		st, err := repo.GetByTicketSequence(ctx, ticketID, seq)
		if errors.Is(err, domain.ErrNotFound) {
			break
		}
		if err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, nil
}

// statementTypeFromIR maps the coarse IR kind to the granular statement type.
// The MySQL lowerer only fills Kind for now, so DDL falls back to a generic
// marker until structural extraction lands (see internal/analyzer TODO).
func statementTypeFromIR(st *ir.Statement) domain.StatementType {
	if st == nil {
		return domain.StatementTypeSelect
	}
	switch st.Kind {
	case ir.KindSelect:
		return domain.StatementTypeSelect
	case ir.KindInsert:
		return domain.StatementTypeInsert
	case ir.KindUpdate:
		return domain.StatementTypeUpdate
	case ir.KindDelete:
		return domain.StatementTypeDelete
	case ir.KindDDL:
		return statementTypeFromDDL(st)
	default:
		return domain.StatementTypeSelect
	}
}

func statementTypeFromDDL(st *ir.Statement) domain.StatementType {
	if st.DDL == nil {
		return domain.StatementTypeAlterTable
	}
	switch strings.ToLower(st.DDL.Action) {
	case "create_table":
		return domain.StatementTypeCreateTable
	case "drop_table":
		return domain.StatementTypeDropTable
	case "alter_table":
		return domain.StatementTypeAlterTable
	case "create_index":
		return domain.StatementTypeCreateIndex
	case "drop_index":
		return domain.StatementTypeDropIndex
	case "create_database":
		return domain.StatementTypeCreateDatabase
	case "drop_database":
		return domain.StatementTypeDropDatabase
	case "truncate_table":
		return domain.StatementTypeTruncateTable
	case "create_view":
		return domain.StatementTypeCreateView
	default:
		return domain.StatementTypeAlterTable
	}
}

// statementStatusFromViolations derives the statement state from rule hits.
func statementStatusFromViolations(violations []domain.Violation) domain.StatementStatus {
	status := domain.StatementStatusWaitingForMrv
	for _, v := range violations {
		if v.Level == uint8(domain.RuleLevelError) {
			return domain.StatementStatusVldFailure
		}
		if v.Level == uint8(domain.RuleLevelWarning) {
			status = domain.StatementStatusVldWarning
		}
	}
	return status
}

func joinViolations(violations []domain.Violation) string {
	if len(violations) == 0 {
		return ""
	}
	msgs := make([]string, 0, len(violations))
	for _, v := range violations {
		msgs = append(msgs, v.Message)
	}
	return strings.Join(msgs, "; ")
}
