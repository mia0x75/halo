package validate

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sync"

	"github.com/go-xorm/core"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/g"
	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
	"github.com/mia0x75/parser/ast"
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
	Stmts     []*models.Statement      // 全部等待审核的数据
	Ticket    *models.Ticket           // 等待审核的工单
	Cluster   *models.Cluster          // 目标群集
	Databases []models.Database        // 目标群集所有的用户数据库
	Tables    map[string][]*core.Table // 目标群集上目标数据库的元数据
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

// Run 调用入口
func Run(stmts []*models.Statement, cluster *models.Cluster, ticket *models.Ticket) {
	wg := &sync.WaitGroup{}
	ctx := &Context{
		Cluster: cluster,
		Ticket:  ticket,
		Stmts:   stmts,
	}
	passwd := func(c *models.Cluster) []byte {
		bs, _ := tools.DecryptAES(c.Password, g.Config().Secret.Crypto)
		return bs
	}
	var err error
	if ctx.Tables, err = cluster.Metadata("*", passwd); err != nil {
		fmt.Println(err)
	}
	if ctx.Databases, err = cluster.Databases(passwd); err != nil {
		fmt.Println(err)
	}

	for g, v := range registry {
		if v.Enabled() {
			wg.Add(1)
			v.SetGroup(g)
			v.SetContext(ctx)
			go v.Validate(wg)
		}
	}
	wg.Wait()
}

// Validator 接口定义
type Validator interface {
	SetGroup(gid uint16) // 设定分组，获得规则列表
	GetRules() []*models.Rule
	Enabled() bool
	SetContext(ctx *Context)     // 分组规则是否启用
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
	Rules []*models.Rule // 全部适用的规则
	Ctx   *Context       // 上下文，保存检测报告？
	Vi    []VisitInfo
}

// SetContext 设置上下文
func (v *vldr) SetContext(ctx *Context) {
	v.Ctx = ctx
}

// GetRules 设置验证分组，并初始化验证规则
func (v *vldr) GetRules() []*models.Rule {
	return v.Rules
}

// SetGroup 设置验证分组，并初始化验证规则
func (v *vldr) SetGroup(gid uint16) {
	rules := caches.RulesMap.Filter(func(elem *models.Rule) bool {
		if elem.VldrGroup == gid {
			return true
		}
		return false
	})
	v.Rules = rules
}

// IndexInfo 获取表信息
func (v *vldr) IndexInfo(database, table, index string) *core.Index {
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
func (v *vldr) ColumnInfo(database, table, column string) *core.Column {
	if ti := v.TableInfo(database, table); ti != nil {
		return ti.GetColumn(column)
	}
	return nil
}

// TableInfo 获取表信息
func (v *vldr) TableInfo(database, table string) *core.Table {
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
		return 0, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

// WalkSubquery 递归检查自查询有效性并获取访问信息
func (v *vldr) WalkSubquery(node ast.SubqueryExpr, stack *Stack) {
	return
}

// Walk 语法树分析
func (v *vldr) Walk(node ast.Node) {
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
	case *ast.UnionStmt:
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

// Walk TODO
func Walk(ctx *Context, node ast.Node) []VisitInfo { //may have duplicate db name
	var vi []VisitInfo

	if node == nil {
		return vi
	}

	switch x := node.(type) {
	case *ast.TableSource:
		Walk(ctx, x.Source)
	case *ast.UnionStmt:
		if x.SelectList != nil {
			for _, sel := range x.SelectList.Selects {
				Walk(ctx, sel)
			}
		}
	case *ast.SelectStmt:
		if x.From != nil {
			Walk(ctx, x.From.TableRefs)
		}
		if x.GroupBy != nil {
			for _, item := range x.OrderBy.Items {
				fmt.Printf("x.GroupBy.Items: %T, %+v\n", item, item)
			}
		}
		if x.OrderBy != nil {
			for _, item := range x.OrderBy.Items {
				fmt.Printf("x.OrderBy.Items: %T, %+v\n", item, item)
			}
		}
		if x.Where != nil {
			fmt.Printf("x.Where: %T, %+v", x.Where, x.Where)
			switch expr := x.Where.(type) {
			case *ast.AggregateFuncExpr:
			case *ast.WindowFuncExpr:
			case *ast.BetweenExpr:
			case *ast.BinaryOperationExpr:
				fmt.Printf("ast.BinaryOperationExpr, L:%T, %+v\n", expr.L, expr.L)
				fmt.Printf("ast.BinaryOperationExpr, R:%T, %+v\n", expr.R, expr.R)
			case *ast.CaseExpr:
			case *ast.ColumnNameExpr:
			case *ast.CompareSubqueryExpr:
			case *ast.DefaultExpr:
			case *ast.ExistsSubqueryExpr:
			case *ast.FuncCallExpr:
			case *ast.FuncCastExpr:
			case *ast.IsNullExpr:
			case *ast.IsTruthExpr:
			case *ast.ParenthesesExpr:
			case *ast.PatternInExpr:
			case *ast.PatternLikeExpr:
			case *ast.PatternRegexpExpr:
			case *ast.PositionExpr:
			case *ast.RowExpr:
			case *ast.SubqueryExpr:
			case *ast.UnaryOperationExpr:
			case *ast.ValuesExpr:
			case *ast.VariableExpr:
			}
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
		Walk(ctx, x.Left)
		Walk(ctx, x.Right)
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

// Call 方法反射
func Call(object interface{}, method string, params ...interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("错误代码: 1500, 错误信息: %s\n", err)
			return
		}
	}()

	v := reflect.ValueOf(object)
	// Check if the passed interface is a pointer
	if v.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface, so we have a pointer to work with
		v = reflect.New(reflect.TypeOf(object))
	}
	// Get the method by name
	f := v.MethodByName(method)
	if !f.IsValid() {
		return fmt.Errorf("Couldn't find method `%s` in interface `%s`, is it exported?", method, v.Type())
	}
	// 'dereference' with Elem() and get the field by name
	// f := ValueIface.Elem().FieldByName(FieldName)
	// if !f.IsValid() {
	// 	return fmt.Errorf("Interface `%s` does not have the field `%s`", v.Type(), f)
	// }
	// 判断函数参数和传入的参数是否相等
	if len(params) != f.Type().NumIn() {
		return fmt.Errorf("")
	}
	if f.Type().NumIn() == 0 {
		f.Call(nil)
	} else {
		// 然后将传入参数转为反射类型切片
		in := make([]reflect.Value, len(params))
		for k, param := range params {
			in[k] = reflect.ValueOf(param)
		}
		// 利用函数反射对象的call方法调用函数.
		f.Call(in)
	}
	return nil
}

// Match 正则匹配
func Match(r *models.Rule, src string, params ...interface{}) (err error) {
	match := false
	match, err = regexp.MatchString(r.Values, src)
	if err != nil {
		err = fmt.Errorf("`%s`不是一个有的效正则表达式。", r.Values)
		return
	}

	if match {
		return nil
	}
	err = fmt.Errorf(r.Message, params...)
	return
}
