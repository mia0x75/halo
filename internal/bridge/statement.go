package bridge

import (
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/transport/graphql/models"
)

// StatementFromDomain converts a domain entity to its legacy model. The legacy
// model also carries non-persisted AST (StmtNode) and Violations fields, which
// stay nil: they are populated by the validation pipeline, not by a DB read.
func StatementFromDomain(s *domain.Statement) *models.Statement {
	if s == nil {
		return nil
	}
	return &models.Statement{
		TicketID:     s.TicketID,
		Sequence:     s.Sequence,
		UUID:         s.UUID,
		Content:      s.Content,
		Type:         uint8(s.Type),
		Status:       uint8(s.Status),
		Report:       s.Report,
		Plan:         derefStr(s.Plan),
		Results:      derefStr(s.Results),
		RowsAffected: deref(s.RowsAffected),
		Version:      s.Version,
		UpdateAt:     deref(s.UpdateAt),
		CreateAt:     s.CreateAt,
	}
}

// StatementToDomain converts a legacy model to its domain entity.
func StatementToDomain(s *models.Statement) *domain.Statement {
	if s == nil {
		return nil
	}
	return &domain.Statement{
		TicketID:     s.TicketID,
		Sequence:     s.Sequence,
		UUID:         s.UUID,
		Content:      s.Content,
		Type:         domain.StatementType(s.Type),
		Status:       domain.StatementStatus(s.Status),
		Report:       s.Report,
		Plan:         strPtr(s.Plan),
		Results:      strPtr(s.Results),
		RowsAffected: uintPtr(s.RowsAffected),
		Version:      s.Version,
		UpdateAt:     uintPtr(s.UpdateAt),
		CreateAt:     s.CreateAt,
	}
}
