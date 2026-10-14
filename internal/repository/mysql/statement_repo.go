package mysqlrepo

import (
	"context"

	"github.com/mia0x75/halo/internal/domain"
)

var _ domain.StatementRepository = (*StatementRepo)(nil)

type StatementRepo struct {
	q *Queries
}

func NewStatementRepo(q *Queries) *StatementRepo {
	return &StatementRepo{q: q}
}

func statementFromRow(r MmStatement) *domain.Statement {
	return &domain.Statement{
		TicketID:     uint(r.TicketID),
		Sequence:     r.Sequence,
		UUID:         r.Uuid,
		Content:      r.Content,
		Type:         domain.StatementType(r.Type),
		Status:       domain.StatementStatus(r.Status),
		Report:       r.Report,
		Plan:         stringFromNull(r.Plan),
		Results:      stringFromNull(r.Results),
		RowsAffected: uintFromNull(r.RowsAffected),
		Version:      int(r.Version),
		UpdateAt:     uintFromNull(r.UpdateAt),
		CreateAt:     uint(r.CreateAt),
	}
}

func (r *StatementRepo) Create(ctx context.Context, s *domain.Statement) error {
	return wrap("statement.Create", r.q.CreateStatement(ctx, CreateStatementParams{
		TicketID:     uint32(s.TicketID),
		Sequence:     s.Sequence,
		Uuid:         s.UUID,
		Content:      s.Content,
		Type:         uint8(s.Type),
		Status:       uint8(s.Status),
		Report:       s.Report,
		Plan:         nullFromString(s.Plan),
		Results:      nullFromString(s.Results),
		RowsAffected: nullFromUint(s.RowsAffected),
		Version:      uint32(s.Version),
		UpdateAt:     nullFromUint(s.UpdateAt),
		CreateAt:     uint32(s.CreateAt),
	}))
}

func (r *StatementRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Statement, error) {
	row, err := r.q.GetStatementByUUID(ctx, uuid)
	if err != nil {
		return nil, wrap("statement.GetByUUID", err)
	}
	return statementFromRow(row), nil
}

// GetByTicketSequence resolves the row by its natural composite key.
func (r *StatementRepo) GetByTicketSequence(ctx context.Context, ticketID uint, sequence uint16) (*domain.Statement, error) {
	row, err := r.q.GetStatementByTicketSequence(ctx, GetStatementByTicketSequenceParams{
		TicketID: uint32(ticketID),
		Sequence: sequence,
	})
	if err != nil {
		return nil, wrap("statement.GetByTicketSequence", err)
	}
	return statementFromRow(row), nil
}

func (r *StatementRepo) Update(ctx context.Context, s *domain.Statement) error {
	return wrap("statement.Update", r.q.UpdateStatement(ctx, UpdateStatementParams{
		Content:      s.Content,
		Type:         uint8(s.Type),
		Status:       uint8(s.Status),
		Report:       s.Report,
		Plan:         nullFromString(s.Plan),
		Results:      nullFromString(s.Results),
		RowsAffected: nullFromUint(s.RowsAffected),
		Version:      uint32(s.Version),
		UpdateAt:     nullFromUint(s.UpdateAt),
		Uuid:         s.UUID,
	}))
}

func (r *StatementRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	return wrap("statement.DeleteByUUID", r.q.DeleteStatementByUUID(ctx, uuid))
}

func (r *StatementRepo) DeleteByTicketSequence(ctx context.Context, ticketID uint, sequence uint16) error {
	return wrap("statement.DeleteByTicketSequence", r.q.DeleteStatementByTicketSequence(ctx, DeleteStatementByTicketSequenceParams{
		TicketID: uint32(ticketID),
		Sequence: sequence,
	}))
}

// DeleteByTicket removes every statement of a ticket in one shot.
func (r *StatementRepo) DeleteByTicket(ctx context.Context, ticketID uint) error {
	return wrap("statement.DeleteByTicket", r.q.DeleteStatementsByTicket(ctx, uint32(ticketID)))
}

func (r *StatementRepo) List(ctx context.Context, cursor string, limit int) ([]*domain.Statement, error) {
	rows, err := r.q.ListStatements(ctx, ListStatementsParams{
		Uuid:  cursor,
		Limit: int32(limit),
	})
	if err != nil {
		return nil, wrap("statement.List", err)
	}
	out := make([]*domain.Statement, 0, len(rows))
	for _, row := range rows {
		out = append(out, statementFromRow(row))
	}
	return out, nil
}
