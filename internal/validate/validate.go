package validate

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/mia0x75/halo/internal/agent"
	"github.com/mia0x75/halo/internal/agent/meta"
	"github.com/mia0x75/halo/internal/bridge"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/domain"
	"github.com/mia0x75/halo/internal/models"
	"github.com/mia0x75/halo/internal/platform/config"
	"github.com/mia0x75/halo/internal/platform/crypto"
	"github.com/mia0x75/parser/ast"
	log "github.com/sirupsen/logrus"
)

// 接口确认，如果有报错表示没有实现接口的全部方法
var (
	_ Validator = &DatabaseCreateVldr{}
	_ Validator = &DatabaseAlterVldr{}
	_ Validator = &DatabaseDropVldr{}
	_ Validator = &TableCreateVldr{}
	_ Validator = &TableAlterVldr{}
	_ Validator = &TableRenameVldr{}
	_ Validator = &TableDropVldr{}
	_ Validator = &IndexCreateVldr{}
	_ Validator = &IndexDropVldr{}
	_ Validator = &InsertVldr{}
	_ Validator = &UpdateVldr{}
	_ Validator = &DeleteVldr{}
	_ Validator = &ViewCreateVldr{}
	_ Validator = &ViewCreateVldr{}
	_ Validator = &ViewCreateVldr{}
	_ Validator = &FuncCreateVldr{}
	_ Validator = &FuncAlterVldr{}
	_ Validator = &FuncDropVldr{}
	_ Validator = &TriggerCreateVldr{}
	_ Validator = &TriggerAlterVldr{}
	_ Validator = &TriggerDropVldr{}
	_ Validator = &EventCreateVldr{}
	_ Validator = &EventAlterVldr{}
	_ Validator = &EventDropVldr{}
	_ Validator = &ProcCreateVldr{}
	_ Validator = &ProcCreateVldr{}
	_ Validator = &ProcCreateVldr{}
	_ Validator = &MiscVldr{}
)

// Context 在不同的验证组中共享内容
type Context struct {
	Subjects  []*Subject               // 全部等待审核的数据
	Ticket    *domain.Ticket           // 等待审核的工单
	Cluster   *domain.Cluster          // 目标群集
	Databases []models.Database        // 目标群集所有的用户数据库
	Tables    map[string][]*meta.Table // 目标群集上目标数据库的元数据
	// Objects 是各库的用户对象清单（视图/函数/存储过程/触发器/事件），供存在性规则
	// 使用。目标库未返回元数据时对应库为 nil，表示「未知」而非「不存在」。
	Objects map[string]*meta.Objects
	// Rules is the rule set the engine runs against. Run fills it once from
	// the cache; RunRegistry distributes it to each group.
	Rules []*Rule
}

// 相当于注册表
var registry = map[uint16]Validator{
	100: &DatabaseCreateVldr{},
	101: &DatabaseAlterVldr{},
	102: &DatabaseDropVldr{},
	110: &TableCreateVldr{},
	120: &TableAlterVldr{},
	130: &TableRenameVldr{},
	140: &TableDropVldr{},
	150: &IndexCreateVldr{},
	151: &IndexDropVldr{},
	160: &InsertVldr{},
	161: &ReplaceVldr{},
	162: &UpdateVldr{},
	163: &DeleteVldr{},
	170: &SelectVldr{},
	180: &ViewCreateVldr{},
	181: &ViewAlterVldr{},
	182: &ViewDropVldr{},
	190: &FuncCreateVldr{},
	191: &FuncAlterVldr{},
	192: &FuncDropVldr{},
	200: &TriggerCreateVldr{},
	201: &TriggerAlterVldr{},
	202: &TriggerDropVldr{},
	210: &EventCreateVldr{},
	211: &EventAlterVldr{},
	212: &EventDropVldr{},
	220: &ProcCreateVldr{},
	221: &ProcAlterVldr{},
	222: &ProcDropVldr{},
	230: &MiscVldr{},
}

// Run 调用入口。签名保持不变以兼容 resolvers。
//
// It performs the IO (target DB metadata) and reads the rule set once, then
// delegates the pure rule execution to RunRegistry. Violations are written back
// in place onto each statement, exactly as before.
func Run(stmts []*models.Statement, cluster *models.Cluster, ticket *models.Ticket) {
	subjects := make([]*Subject, 0, len(stmts))
	for _, m := range stmts {
		subjects = append(subjects, subjectFromModels(m))
	}
	ctx := &Context{
		Cluster:  bridge.ClusterToDomain(cluster),
		Ticket:   bridge.TicketToDomain(ticket),
		Subjects: subjects,
	}
	target := ctx.Cluster
	// PG connections are database-scoped, so bind the target to the ticket's
	// database before reading metadata. MySQL ignores this at connect time but
	// the field is harmless there.
	if ctx.Ticket != nil {
		target.Database = ctx.Ticket.Database
	}
	passwd := func(c *domain.Cluster) string {
		bs, _ := crypto.DecryptAES(c.Password, config.Config().Secret.Crypto)
		return string(bs)
	}
	exec := agent.LocalExecutor{}
	var err error
	if ctx.Tables, err = exec.LoadTables(context.Background(), target, passwd(target)); err != nil {
		fmt.Println(err)
	}
	if ctx.Objects, err = exec.LoadObjects(context.Background(), target, passwd(target)); err != nil {
		fmt.Println(err)
	}
	if databases, derr := exec.Databases(context.Background(), target, passwd(target)); derr != nil {
		fmt.Println(derr)
	} else {
		ctx.Databases = make([]models.Database, 0, len(databases))
		for _, db := range databases {
			ctx.Databases = append(ctx.Databases, models.Database{
				Name:    db.Name,
				Charset: db.Charset,
				Collate: db.Collate,
			})
		}
	}
	cachedRules := cache.RulesMap.All()
	ctx.Rules = make([]*Rule, 0, len(cachedRules))
	for _, m := range cachedRules {
		ctx.Rules = append(ctx.Rules, ruleFromDomain(bridge.RuleToDomain(m)))
	}

	RunRegistry(ctx)
}

// Violation pairs a rule hit with the statement it was raised on.
type Violation struct {
	Statement *Subject
	Clause    *Clause
}

// RunRegistry is the pure rule engine: it runs every enabled validator group
// over c and returns the aggregated hits. It performs no IO and reads no
// caches; each group only sees the rules carried by c.
func RunRegistry(c *Context) []Violation {
	wg := &sync.WaitGroup{}
	for g, v := range registry {
		if v.Enabled() {
			wg.Add(1)
			v.SetContext(c)
			v.SetRules(rulesByGroup(c.Rules, g))
			go v.Validate(wg)
		}
	}
	wg.Wait()
	return collectViolations(c.Subjects)
}

// rulesByGroup picks the rules addressed to a single validator group.
func rulesByGroup(rules []*Rule, gid uint16) []*Rule {
	out := make([]*Rule, 0, len(rules))
	for _, r := range rules {
		if uint16(r.VldrGroup) == gid {
			out = append(out, r)
		}
	}
	return out
}

// collectViolations flattens the per-statement hits into a single slice.
func collectViolations(stmts []*Subject) []Violation {
	var out []Violation
	for _, s := range stmts {
		if s.Violations == nil {
			continue
		}
		for _, clause := range s.Violations.Clauses() {
			out = append(out, Violation{Statement: s, Clause: clause})
		}
	}
	return out
}

// Validator 接口定义
type Validator interface {
	SetRules(rules []*Rule) // 设定分组，获得规则列表
	GetRules() []*Rule
	Enabled() bool               // 分组规则是否启用
	SetContext(ctx *Context)     // 设定上下文
	Validate(wg *sync.WaitGroup) // 验证
}

// VisitInfo 访问信息
type VisitInfo struct {
	Database string
	Table    *TableEntry
	Column   *ColumnEntry
	Index    *IndexEntry
}

func (v *vldr) appendVisitInfo(vi VisitInfo) {
	v.Vi = append(v.Vi, vi)
}

// TableEntry 抽象语法树需要访问的表
type TableEntry struct {
	Name   string
	Alias  string
	IsView bool
}

// ColumnEntry 抽象语法树需要访问的列
type ColumnEntry struct {
	Table TableEntry
	Name  string
	Alias string
}

// IndexEntry 抽象语法树需要访问的索引
type IndexEntry struct {
	Table TableEntry
	Name  string
}

type vldr struct {
	Rules []*Rule  // 全部适用的规则
	Ctx   *Context // 上下文，保存检测报告？
	Vi    []VisitInfo
}

// SetContext 设置上下文
func (v *vldr) SetContext(ctx *Context) {
	v.Ctx = ctx
}

// GetRules 设置验证分组，并初始化验证规则
func (v *vldr) GetRules() []*Rule {
	return v.Rules
}

// SetRules 设置验证分组适用的规则列表
func (v *vldr) SetRules(rules []*Rule) {
	v.Rules = rules
}

// IndexInfo 获取表信息
func (v *vldr) IndexInfo(database, table, index string) *meta.Index {
	if ti := v.TableInfo(database, table); ti != nil {
		for _, idx := range ti.Indexes {
			if idx.Name == index {
				return idx
			}
		}
	}
	return nil
}

// ColumnInfo 获取表信息
func (v *vldr) ColumnInfo(database, table, column string) *meta.Column {
	if ti := v.TableInfo(database, table); ti != nil {
		return ti.GetColumn(column)
	}
	return nil
}

// TableInfo 获取表信息
func (v *vldr) TableInfo(database, table string) *meta.Table {
	if database == "" {
		database = v.Ctx.Ticket.Database
	}
	if tables, ok := v.Ctx.Tables[database]; ok {
		for _, elem := range tables {
			if elem.Name == table {
				return elem
			}
		}
	}
	return nil
}

// Stack LIFO栈的简单实现
type Stack struct {
	sync.Mutex // you don't have to do this if you don't want thread safety
	s          []int
}

// NewStack 初始化堆栈
func NewStack() *Stack {
	return &Stack{
		s: make([]int, 0),
	}
}

// Push 压栈
func (s *Stack) Push(v int) {
	s.Lock()
	defer s.Unlock()

	s.s = append(s.s, v)
}

// Pop 弹出
func (s *Stack) Pop() (int, error) {
	s.Lock()
	defer s.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("empty stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

// WalkSubquery 递归检查自查询有效性并获取访问信息
func (v *vldr) WalkSubquery(node ast.SubqueryExpr, stack *Stack) {
}

// Walk 语法树分析
func (v *vldr) Walk(node ast.Node) {
	// Validators are package singletons, so VisitInfo would otherwise leak from
	// the previous statement (and the previous Run) into this one.
	v.Vi = v.Vi[:0]
	switch x := node.(type) {
	case *ast.DeallocateStmt:
		return
	case *ast.DeleteStmt:
		if x.TableRefs != nil {
			for _, elem := range Walk(v.Ctx, x.TableRefs.TableRefs) {
				v.appendVisitInfo(elem)
			}
		}
	case *ast.ExecuteStmt:
		// TODO:
	case *ast.ExplainStmt:
		// TODO:
	case *ast.ExplainForStmt:
		// TODO:
	case *ast.InsertStmt:
		for _, elem := range Walk(v.Ctx, x.Table.TableRefs) {
			v.appendVisitInfo(elem)
		}
		for _, elem := range Walk(v.Ctx, x.Select) {
			v.appendVisitInfo(elem)
		}
	case *ast.LoadDataStmt:
		// TODO:
	case *ast.PrepareStmt:
		// TODO:
	case *ast.SelectStmt:
		for _, elem := range Walk(v.Ctx, x) {
			v.appendVisitInfo(elem)
		}
	case *ast.SetOprStmt:
		for _, ss := range x.SelectList.Selects {
			for _, elem := range Walk(v.Ctx, ss) {
				v.appendVisitInfo(elem)
			}
		}
	case *ast.UpdateStmt:
		if x.TableRefs != nil {
			for _, elem := range Walk(v.Ctx, x.TableRefs.TableRefs) {
				v.appendVisitInfo(elem)
			}
		}
	case *ast.ShowStmt:
		// TODO:
	case *ast.AnalyzeTableStmt:
		for _, elem := range x.TableNames {
			table := &TableEntry{
				Name: elem.Name.O,
			}
			database := elem.Schema.O
			if database == "" {
				database = v.Ctx.Ticket.Database
			}
			vi := VisitInfo{
				Database: database,
				Table:    table,
			}
			v.appendVisitInfo(vi)
		}
		// TODO: 索引
	case *ast.BinlogStmt, *ast.FlushStmt, *ast.UseStmt,
		*ast.BeginStmt, *ast.CommitStmt, *ast.RollbackStmt, *ast.CreateUserStmt, *ast.SetPwdStmt,
		*ast.GrantStmt, *ast.DropUserStmt, *ast.AlterUserStmt, *ast.RevokeStmt, *ast.KillStmt, *ast.DropStatsStmt,
		*ast.GrantRoleStmt, *ast.RevokeRoleStmt, *ast.SetRoleStmt, *ast.SetDefaultRoleStmt:
		// TODO:
	case *ast.AlterTableStmt:
		table := &TableEntry{
			Name: x.Table.Name.O,
		}
		database := x.Table.Schema.O
		if database == "" {
			database = v.Ctx.Ticket.Database
		}
		vi := VisitInfo{
			Database: database,
			Table:    table,
		}
		v.appendVisitInfo(vi)
	case *ast.CreateDatabaseStmt:
		vi := VisitInfo{
			Database: x.Name,
		}
		v.appendVisitInfo(vi)
	case *ast.CreateIndexStmt:
		table := &TableEntry{
			Name: x.Table.Name.O,
		}
		database := x.Table.Schema.O
		if database == "" {
			database = v.Ctx.Ticket.Database
		}
		vi := VisitInfo{
			Database: database,
			Table:    table,
		}
		v.appendVisitInfo(vi)
	case *ast.CreateTableStmt:
		table := &TableEntry{
			Name: x.Table.Name.O,
		}
		database := x.Table.Schema.O
		if database == "" {
			database = v.Ctx.Ticket.Database
		}
		vi := VisitInfo{
			Database: database,
			Table:    table,
		}
		v.appendVisitInfo(vi)
		if x.ReferTable != nil {
			table := &TableEntry{
				Name: x.ReferTable.Name.O,
			}
			database := x.ReferTable.Schema.O
			if database == "" {
				database = v.Ctx.Ticket.Database
			}
			vi := VisitInfo{
				Database: database,
				Table:    table,
			}
			v.appendVisitInfo(vi)
		}
	case *ast.CreateViewStmt:
		// The target of a view rule is the view itself, not the tables its
		// SELECT reads from.
		database := x.ViewName.Schema.O
		if database == "" {
			database = v.Ctx.Ticket.Database
		}
		vi := VisitInfo{
			Database: database,
			Table:    &TableEntry{Name: x.ViewName.Name.O},
		}
		v.appendVisitInfo(vi)
	case *ast.AlterViewStmt:
		// Same rationale as CREATE VIEW: the target is the view itself.
		database := x.ViewName.Schema.O
		if database == "" {
			database = v.Ctx.Ticket.Database
		}
		vi := VisitInfo{
			Database: database,
			Table:    &TableEntry{Name: x.ViewName.Name.O},
		}
		v.appendVisitInfo(vi)
	case *ast.DropDatabaseStmt:
		vi := VisitInfo{
			Database: x.Name,
		}
		v.appendVisitInfo(vi)
	case *ast.DropIndexStmt:
		table := &TableEntry{
			Name: x.Table.Name.O,
		}
		database := x.Table.Schema.O
		if database == "" {
			database = v.Ctx.Ticket.Database
		}
		vi := VisitInfo{
			Database: database,
			Table:    table,
		}
		v.appendVisitInfo(vi)
	case *ast.DropTableStmt:
		for _, elem := range x.Tables {
			table := &TableEntry{
				Name: elem.Name.O,
			}
			database := elem.Schema.O
			if database == "" {
				database = v.Ctx.Ticket.Database
			}
			vi := VisitInfo{
				Database: database,
				Table:    table,
			}
			v.appendVisitInfo(vi)
		}
	case *ast.TruncateTableStmt:
		table := &TableEntry{
			Name: x.Table.Name.O,
		}
		database := x.Table.Schema.O
		if database == "" {
			database = v.Ctx.Ticket.Database
		}
		vi := VisitInfo{
			Database: database,
			Table:    table,
		}
		v.appendVisitInfo(vi)
	case *ast.RenameTableStmt:
		tables := []*ast.TableName{x.OldTable, x.NewTable}
		for _, elem := range tables {
			table := &TableEntry{
				Name: elem.Name.O,
			}
			database := elem.Schema.O
			if database == "" {
				database = v.Ctx.Ticket.Database
			}
			vi := VisitInfo{
				Database: database,
				Table:    table,
			}
			v.appendVisitInfo(vi)
		}
	}
}

// Walk flattens a table-ref AST node into the VisitInfo list the schema-aware
// rules resolve against. Results may repeat when a table is referenced twice.
func Walk(ctx *Context, node ast.Node) []VisitInfo {
	var vi []VisitInfo

	if node == nil {
		return vi
	}

	switch x := node.(type) {
	case *ast.TableSource:
		vi = append(vi, Walk(ctx, x.Source)...)
	case *ast.SetOprStmt:
		if x.SelectList != nil {
			for _, sel := range x.SelectList.Selects {
				vi = append(vi, Walk(ctx, sel)...)
			}
		}
	case *ast.SetOprSelectList:
		for _, sel := range x.Selects {
			vi = append(vi, Walk(ctx, sel)...)
		}
	case *ast.SelectStmt:
		if x.From != nil {
			vi = append(vi, Walk(ctx, x.From.TableRefs)...)
		}
	case *ast.TableName:
		table := &TableEntry{
			Name: x.Name.O,
		}
		database := x.Schema.O
		if database == "" {
			database = ctx.Ticket.Database
		}
		elem := VisitInfo{
			Database: database,
			Table:    table,
		}
		vi = append(vi, elem)
	case *ast.Join:
		vi = append(vi, Walk(ctx, x.Left)...)
		if x.Right != nil {
			vi = append(vi, Walk(ctx, x.Right)...)
		}
	}

	return vi
}

// func WalkExpr(s ast.ExprNode) (L []TableInfo) {
// 	nodes := []ast.ExprNode{}
// 	switch s.(type) {
// 	case *ast.BetweenExpr:
// 		expr := s.(*ast.BetweenExpr)
// 		nodes = []ast.ExprNode{expr.Left, expr.Right}
// 	case *ast.BinaryOperationExpr:
// 		expr := s.(*ast.BinaryOperationExpr)
// 		nodes = []ast.ExprNode{expr.L, expr.R}
// 	case *ast.CaseExpr:
// 		expr := s.(*ast.CaseExpr)
// 		nodes = []ast.ExprNode{expr.Value, expr.ElseClause}
// 		for _, w := range expr.WhenClauses {
// 			nodes = append(nodes, w.Expr)
// 			nodes = append(nodes, w.Result)
// 		}
// 	case *ast.ColumnNameExpr:
// 		// expr := s.(*ast.ColumnNameExpr)
// 		// fmt.Printf("*ast.ColumnNameExpr: %T, %+v\n", expr, expr)
// 	case *ast.DefaultExpr:
// 		// expr := s.(*ast.DefaultExpr)
// 		// fmt.Printf("*ast.DefaultExpr: %T, %+v\n", expr, expr)
// 	case *ast.CompareSubqueryExpr:
// 		expr := s.(*ast.CompareSubqueryExpr)
// 		nodes = []ast.ExprNode{expr.L, expr.R}
// 	case *ast.ExistsSubqueryExpr:
// 		expr := s.(*ast.ExistsSubqueryExpr)
// 		nodes = []ast.ExprNode{expr.Sel}
// 	case *ast.IsNullExpr:
// 		expr := s.(*ast.IsNullExpr)
// 		nodes = []ast.ExprNode{expr.Expr}
// 	case *ast.IsTruthExpr:
// 		expr := s.(*ast.IsTruthExpr)
// 		nodes = []ast.ExprNode{expr.Expr}
// 	case *ast.ParenthesesExpr:
// 		expr := s.(*ast.ParenthesesExpr)
// 		nodes = []ast.ExprNode{expr.Expr}
// 	case *ast.PatternInExpr:
// 		expr := s.(*ast.PatternInExpr)
// 		nodes = []ast.ExprNode{expr.Expr}
// 		if expr.Sel != nil {
// 			nodes = append(nodes, expr.Sel)
// 		}
// 		nodes = append(nodes, expr.List...)
// 	case *ast.PatternLikeExpr:
// 		expr := s.(*ast.PatternLikeExpr)
// 		nodes = []ast.ExprNode{expr.Expr, expr.Pattern}
// 	case *ast.PatternRegexpExpr:
// 		expr := s.(*ast.PatternRegexpExpr)
// 		nodes = []ast.ExprNode{expr.Expr, expr.Pattern}
// 	case *ast.PositionExpr:
// 		expr := s.(*ast.PositionExpr)
// 		nodes = []ast.ExprNode{expr.P}
// 	case *ast.RowExpr:
// 		expr := s.(*ast.RowExpr)
// 		nodes = append(nodes, expr.Values...)
// 	case *ast.SubqueryExpr:
// 		expr := s.(*ast.SubqueryExpr)
// 		stmt := expr.Query.(*ast.SelectStmt)
// 		L = append(L, WalkJoin(stmt.From.TableRefs)...)
// 		nodes = []ast.ExprNode{stmt.Where}
// 	case *ast.UnaryOperationExpr:
// 		expr := s.(*ast.UnaryOperationExpr)
// 		nodes = []ast.ExprNode{expr.V}
// 	case *ast.ValuesExpr:
// 		expr := s.(*ast.ValuesExpr)
// 		L = append(L, TableInfo{expr.Column.Name.Schema.O, expr.Column.Name.Table.O, ""})
// 	case *ast.VariableExpr:
// 		expr := s.(*ast.VariableExpr)
// 		nodes = []ast.ExprNode{expr.Value}
// 	}

// 	for _, node := range nodes {
// 		if node == nil {
// 			continue
// 		}
// 		L = append(L, WalkExpr(node)...)
// 	}

// 	return
// }

// func WalkJoin(s *ast.Join) (L []TableInfo) {
// 	nodes := []ast.ResultSetNode{s.Left, s.Right}
// 	for _, node := range nodes {
// 		if node == nil {
// 			continue
// 		}
// 		switch node.(type) {
// 		case *ast.TableSource:
// 			ts := node.(*ast.TableSource)
// 			switch ts.Source.(type) {
// 			case *ast.TableName:
// 				tn := ts.Source.(*ast.TableName)
// 				L = append(L, TableInfo{Schema: tn.Schema.O, Name: tn.Name.O})
// 			}
// 		case *ast.Join:
// 			L = append(L, WalkJoin(node.(*ast.Join))...)
// 		}
// 	}
// 	return
// }

// DatabaseInfo 根据数据库名称获取数据库信息
func (v *vldr) DatabaseInfo(name string) *models.Database {
	if name == "" {
		name = v.Ctx.Ticket.Database
	}
	for _, database := range v.Ctx.Databases {
		if database.Name == name {
			return &database
		}
	}
	return nil
}

// objectDB resolves the database a DDL object belongs to, falling back to the
// ticket database when the statement omits the schema qualifier.
func (v *vldr) objectDB(schema string) string {
	if schema != "" {
		return schema
	}
	if v.Ctx.Ticket != nil {
		return v.Ctx.Ticket.Database
	}
	return ""
}

// lookupObject reports whether name exists as kind in the already-resolved
// database. known is false when the target object metadata was not loaded,
// letting callers skip the check instead of treating an unknown object as absent.
func (v *vldr) lookupObject(database string, kind meta.ObjectKind, name string) (exists, known bool) {
	objects, ok := v.Ctx.Objects[database]
	if !ok || objects == nil {
		return false, false
	}
	return objects.Has(kind, name), true
}

// objectRef is the common shape of the routine/view DDL nodes: a possibly
// schema-qualified name.
type objectRef struct {
	Schema string
	Name   string
}

// refOf extracts the schema and bare name from a parser TableName, tolerating
// a nil node.
func refOf(name *ast.TableName) objectRef {
	if name == nil {
		return objectRef{}
	}
	return objectRef{Schema: name.Schema.O, Name: name.Name.O}
}

// Call runs the rule handler registered for (group, name) against validator v.
// Unknown handlers are logged and skipped, matching the previous behaviour where
// reflection failures were swallowed.
func Call(v Validator, group uint16, name string, s *Subject, r *Rule) {
	fn, ok := lookup(group, name)
	if !ok {
		log.Errorf("validate: unknown rule func: group=%d func=%s", group, name)
		return
	}
	fn(v, s, r)
}

// Match 正则匹配
func Match(r *Rule, src string, params ...interface{}) (err error) {
	match := false
	match, err = regexp.MatchString(r.Parameters, src)
	if err != nil {
		err = fmt.Errorf("`%s`不是一个有的效正则表达式。", r.Parameters)
		return
	}

	if match {
		return nil
	}
	err = fmt.Errorf(r.Message, params...)
	return
}
