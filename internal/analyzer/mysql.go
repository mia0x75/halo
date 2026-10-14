package analyzer

import (
	"context"
	"fmt"

	"github.com/mia0x75/halo/internal/analyzer/ir"
	"github.com/mia0x75/halo/internal/domain"

	"github.com/mia0x75/parser"
	"github.com/mia0x75/parser/ast"
	_ "github.com/mia0x75/parser/driver"
)

// MySQL analyzes MySQL-dialect SQL with mia0x75/parser.
type MySQL struct{}

// Parse lowers the first statement of sql into the shared IR.
func (MySQL) Parse(ctx context.Context, sql string) (*ir.Statement, error) {
	p := parser.New()
	stmts, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseFailed, err)
	}
	if len(stmts) == 0 {
		return nil, ErrUnsupportedStmt
	}
	return mysqlToIR(sql, stmts[0]), nil
}

// mysqlToIR lowers a MySQL AST node. Only kind, dialect and raw text are filled
// for now; structural extraction lands with the phase-2 analyzer work.
func mysqlToIR(sql string, node ast.StmtNode) *ir.Statement {
	// TODO(phase2): lower Tables, Columns, Where, Joins and DDL detail.
	return &ir.Statement{
		Kind:    mysqlKind(node),
		Dialect: domain.ClusterTypeMySQL,
		Raw:     sql,
	}
}

// mysqlKind classifies a node; DDL is detected via the marker interface so new
// DDL statements are covered without editing a type switch.
func mysqlKind(node ast.StmtNode) ir.StatementKind {
	switch node.(type) {
	case *ast.SelectStmt:
		return ir.KindSelect
	case *ast.InsertStmt:
		return ir.KindInsert
	case *ast.UpdateStmt:
		return ir.KindUpdate
	case *ast.DeleteStmt:
		return ir.KindDelete
	}
	if _, ok := node.(ast.DDLNode); ok {
		return ir.KindDDL
	}
	return ir.KindOther
}
