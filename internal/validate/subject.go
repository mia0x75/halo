package validate

import (
	"github.com/mia0x75/parser/ast"

	"github.com/mia0x75/halo/internal/models"
)

// Subject is one SQL statement under validation: the persisted fields the
// engine reads plus the parsed AST and the violation accumulator. It replaces
// models.Statement inside the engine so rules no longer depend on the GraphQL
// DTO; only the Run boundary converts.
//
// Violations stays a models type on purpose (ADR 0009): the accumulator owns a
// mutex and is written back onto models.Statement unchanged. Clause is its
// element type, aliased here so rule files build hits without importing models.
type Subject struct {
	TicketID uint
	Sequence uint16
	UUID     string
	Content  string
	Type     uint8
	Status   uint8

	StmtNode   ast.StmtNode
	Violations *models.Violations
}

// Clause is a single rule hit's description. It aliases models.Clause so rule
// files construct hits through this package while the accumulator stays a
// models type (ADR 0009).
type Clause = models.Clause

// subjectFromModels converts a persisted statement into an engine subject. The
// Violations pointer is shared, so rule hits land on the original statement.
func subjectFromModels(m *models.Statement) *Subject {
	if m == nil {
		return nil
	}
	return &Subject{
		TicketID:   m.TicketID,
		Sequence:   m.Sequence,
		UUID:       m.UUID,
		Content:    m.Content,
		Type:       m.Type,
		Status:     m.Status,
		StmtNode:   m.StmtNode,
		Violations: m.Violations,
	}
}
