package validate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mia0x75/parser/ast"

	"github.com/mia0x75/halo/internal/agent/meta"
)

// tableRef is one base table reachable from a DML statement, keyed by the name
// its columns qualify with (alias when present, else table name).
type tableRef struct {
	table *meta.Table
}

// tableRefs resolves the base tables of a DML table-refs clause. Derived tables
// and subqueries are skipped because they carry no catalog metadata.
func (v *vldr) tableRefs(clause *ast.TableRefsClause) map[string]tableRef {
	out := map[string]tableRef{}
	if clause == nil {
		return out
	}
	v.collectTableRefs(clause.TableRefs, out)
	return out
}

func (v *vldr) collectTableRefs(node ast.ResultSetNode, out map[string]tableRef) {
	switch x := node.(type) {
	case *ast.Join:
		v.collectTableRefs(x.Left, out)
		if x.Right != nil {
			v.collectTableRefs(x.Right, out)
		}
	case *ast.TableSource:
		tn, ok := x.Source.(*ast.TableName)
		if !ok {
			return
		}
		table := v.TableInfo(tn.Schema.O, tn.Name.O)
		if table == nil {
			return
		}
		key := tn.Name.O
		if x.AsName.O != "" {
			key = x.AsName.O
		}
		out[key] = tableRef{table: table}
	}
}

// lookup resolves a possibly qualified column to its base table. An unqualified
// column only resolves when exactly one base table is in scope, so ambiguous
// multi-table statements are skipped rather than mis-attributed.
func lookupColumn(refs map[string]tableRef, c *ast.ColumnName) (*meta.Table, bool) {
	if c == nil {
		return nil, false
	}
	if c.Table.O != "" {
		ref, ok := refs[c.Table.O]
		if !ok {
			return nil, false
		}
		return ref.table, true
	}
	var found *meta.Table
	for _, ref := range refs {
		if found != nil && found != ref.table {
			return nil, false
		}
		found = ref.table
	}
	return found, found != nil
}

// collectedColumns returns every column reference in an AST subtree. Wildcards
// are not column references and are therefore absent.
func collectedColumns(node ast.Node) []*ast.ColumnName {
	if node == nil {
		return nil
	}
	c := &columnCollector{}
	node.Accept(c)
	return c.cols
}

type columnCollector struct {
	cols []*ast.ColumnName
}

func (c *columnCollector) Enter(n ast.Node) (ast.Node, bool) {
	if e, ok := n.(*ast.ColumnNameExpr); ok && e.Name != nil {
		c.cols = append(c.cols, e.Name)
	}
	return n, false
}

func (c *columnCollector) Leave(n ast.Node) (ast.Node, bool) {
	return n, true
}

// limitRows reads the constant LIMIT count. ok is false when the clause is
// absent or the count is not a literal.
func limitRows(limit *ast.Limit) (int64, bool) {
	if limit == nil || limit.Count == nil {
		return 0, false
	}
	val, ok := limit.Count.(ast.ValueExpr)
	if !ok {
		return 0, false
	}
	switch n := val.GetValue().(type) {
	case int64:
		return n, true
	case uint64:
		return int64(n), true
	case int:
		return int64(n), true
	}
	return 0, false
}

// appendViolation records one clause whose message is formatted with args.
func appendViolation(s *Subject, r *Rule, args ...interface{}) {
	s.Violations.Append(&Clause{
		Description: fmt.Sprintf(r.Message, args...),
		Level:       uint8(r.Level),
	})
}

// missingColumn reports the bare column name when c cannot be resolved to an
// existing column of a base table in refs. Unresolvable references (ambiguous,
// derived table, unknown table) are treated as out of scope and not reported.
func missingColumn(refs map[string]tableRef, c *ast.ColumnName) (string, bool) {
	if c == nil {
		return "", false
	}
	tbl, ok := lookupColumn(refs, c)
	if !ok || tbl.GetColumn(c.Name.O) != nil {
		return "", false
	}
	return c.Name.O, true
}

// indexColNames flattens the column list of an index definition.
func indexColNames(cols []*ast.IndexColName) []string {
	names := make([]string, 0, len(cols))
	for _, c := range cols {
		if c != nil && c.Column != nil {
			names = append(names, c.Column.Name.O)
		}
	}
	return names
}

// indexOverlay reports whether a and b share the same leading columns, i.e.
// one index's column list is a prefix of the other. MySQL treats the shorter
// index as redundant because the longer one serves every lookup it serves.
func indexOverlay(a, b []string) bool {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	if n == 0 {
		return false
	}
	for i := 0; i < n; i++ {
		if !strings.EqualFold(a[i], b[i]) {
			return false
		}
	}
	return true
}

// sortedIndexNames returns the index names of t in a stable order so the rules
// that scan every index emit violations deterministically.
func sortedIndexNames(t *meta.Table) []string {
	names := make([]string, 0, len(t.Indexes))
	for name := range t.Indexes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// insertTarget resolves the single base table an INSERT/REPLACE writes to, nil
// when the target is a join/derived table or absent from metadata.
func (v *vldr) insertTarget(id *ast.InsertStmt) *meta.Table {
	if id == nil || id.Table == nil {
		return nil
	}
	refs := v.tableRefs(id.Table)
	if len(refs) != 1 {
		return nil
	}
	for _, ref := range refs {
		return ref.table
	}
	return nil
}

// checkTargetInsertColumns flags INSERT/REPLACE columns absent from the target
// table. Both the column list and the SET form are covered.
func checkTargetInsertColumns(v *vldr, id *ast.InsertStmt, s *Subject, r *Rule) {
	target := v.insertTarget(id)
	if target == nil {
		return
	}
	seen := map[string]bool{}
	check := func(name string) {
		if name == "" || seen[name] {
			return
		}
		if target.GetColumn(name) == nil {
			seen[name] = true
			appendViolation(s, r, name)
		}
	}
	for _, c := range id.Columns {
		check(c.Name.O)
	}
	for _, a := range id.Setlist {
		if a.Column != nil {
			check(a.Column.Name.O)
		}
	}
}

// checkNotNullInsertColumns flags NOT NULL columns that are neither provided nor
// defaulted. A statement without an explicit column list populates every column,
// so it is skipped.
func checkNotNullInsertColumns(v *vldr, id *ast.InsertStmt, s *Subject, r *Rule) {
	target := v.insertTarget(id)
	if target == nil || (len(id.Columns) == 0 && len(id.Setlist) == 0) {
		return
	}
	provided := map[string]bool{}
	for _, c := range id.Columns {
		provided[c.Name.O] = true
	}
	for _, a := range id.Setlist {
		if a.Column != nil {
			provided[a.Column.Name.O] = true
		}
	}
	names := make([]string, 0, len(target.Columns))
	for name, col := range target.Columns {
		if col.Nullable || col.AutoIncrement || col.Default != "" {
			continue
		}
		if !provided[name] {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		appendViolation(s, r, name)
	}
}
