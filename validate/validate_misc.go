package validate

import (
	"fmt"
	"sync"

	"github.com/mia0x75/parser/ast"

	"github.com/mia0x75/halo/models"
	"github.com/mia0x75/halo/tools"
)

// MiscVldr 其他各种规则，包含跨语句相关的审核规则，检测一个工单不同语句之间的可能的问题
type MiscVldr struct {
	vldr
}

// Call 利用反射方法动态调用审核函数
func (v *MiscVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *MiscVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *MiscVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, r := range v.Rules {
		if r.Bitwise&1 != 1 {
			continue
		}
		v.Call(r.Func, r)
	}
}

// LockTableProhibited 是否允许LOCK TABLE
// RULE: MSC-L1-001
func (v *MiscVldr) LockTableProhibited(r *models.Rule) {
	for _, s := range v.Ctx.Stmts {
		if _, ok := s.StmtNode.(*ast.LockTableStmt); !ok {
			continue
		}

		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// FlushTableProhibited 是否允许FLUSH TABLE => TODO: FLUSH语法支持不全
// RULE: MSC-L1-002
func (v *MiscVldr) FlushTableProhibited(r *models.Rule) {
	for _, s := range v.Ctx.Stmts {
		ct, ok := s.StmtNode.(*ast.FlushStmt)
		if !ok {
			continue
		}

		switch ct.Tp {
		case ast.FlushTables, ast.FlushStatus, ast.FlushPrivileges, ast.FlushNone:
			c := &models.Clause{
				Description: r.Message,
				Level:       r.Level,
			}
			s.Violations.Append(c)
		default:
			continue // 继续下一轮for循环
		}
	}
}

// TruncateTableProhibited 是否允许TRUNCATE TABLE
// RULE: MSC-L1-003
func (v *MiscVldr) TruncateTableProhibited(r *models.Rule) {
	for _, s := range v.Ctx.Stmts {
		if _, ok := s.StmtNode.(*ast.TruncateTableStmt); !ok {
			continue
		}

		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// MergeRequired 合并同一个表的操作，这个对大表的ALTER很重要
// RULE: MSC-L1-004
func (v *MiscVldr) MergeRequired(r *models.Rule) {
	slice := []string{}
	name := ""
	for _, s := range v.Ctx.Stmts {
		switch s.StmtNode.(type) {
		case *ast.CreateIndexStmt:
			n := s.StmtNode.(*ast.CreateIndexStmt)
			dbName := n.Table.Schema.O
			if dbName == "" {
				dbName = v.Ctx.Ticket.Database
			}
			name = fmt.Sprintf("`%s`.`%s`", dbName, n.Table.Name.L)
		case *ast.CreateTableStmt:
			n := s.StmtNode.(*ast.CreateTableStmt)
			dbName := n.Table.Schema.O
			if dbName == "" {
				dbName = v.Ctx.Ticket.Database
			}
			name = fmt.Sprintf("`%s`.`%s`", dbName, n.Table.Name.L)
		case *ast.AlterTableStmt:
			n := s.StmtNode.(*ast.AlterTableStmt)
			dbName := n.Table.Schema.O
			if dbName == "" {
				dbName = v.Ctx.Ticket.Database
			}
			name = fmt.Sprintf("`%s`.`%s`", dbName, n.Table.Name.L)
		case *ast.DropIndexStmt:
			n := s.StmtNode.(*ast.DropIndexStmt)
			dbName := n.Table.Schema.O
			if dbName == "" {
				dbName = v.Ctx.Ticket.Database
			}
			name = fmt.Sprintf("`%s`.`%s`", dbName, n.Table.Name.L)
		default:
			continue // 继续下一轮for循环
		}
		if tools.Contains(slice, name) {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, name),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		} else {
			if name != "" {
				slice = append(slice, name)
			}
		}
	}
}

// PurgeLogsProhibited 是否允许PURGE LOG
// RULE: MSC-L1-005
func (v *MiscVldr) PurgeLogsProhibited(r *models.Rule) {
	for _, s := range v.Ctx.Stmts {
		if _, ok := s.StmtNode.(*ast.PurgeStmt); !ok {
			continue
		}

		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// UnlockTableProhibited 是否允许UNLOCK TABLES
// RULE: MSC-L1-006
func (v *MiscVldr) UnlockTableProhibited(r *models.Rule) {
	for _, s := range v.Ctx.Stmts {
		if _, ok := s.StmtNode.(*ast.UnlockTableStmt); !ok {
			continue
		}

		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// KillProhibited 是否允许KILL线程
// RULE: MSC-L1-007
func (v *MiscVldr) KillProhibited(r *models.Rule) {
	for _, s := range v.Ctx.Stmts {
		if _, ok := s.StmtNode.(*ast.KillStmt); !ok {
			continue
		}

		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// SplitRequired 是否允许同时出现DDL、DML
// RULE: MSC-L1-008
func (v *MiscVldr) SplitRequired(r *models.Rule) {
	method := uint8(0)
	for _, s := range v.Ctx.Stmts {
		switch s.StmtNode.(type) {
		case *ast.AlterTableStmt, *ast.AlterDatabaseStmt, *ast.CreateDatabaseStmt, *ast.CreateIndexStmt, *ast.CreateTableStmt, *ast.CreateViewStmt, *ast.DropDatabaseStmt, *ast.DropIndexStmt, *ast.DropTableStmt, *ast.RenameTableStmt, *ast.TruncateTableStmt:
			if method != 0 && method != 1 {
				c := &models.Clause{
					Description: r.Message,
					Level:       r.Level,
				}
				s.Violations.Append(c)
			} else {
				method = 1
			}
		case *ast.DeleteStmt, *ast.InsertStmt, *ast.UnionStmt, *ast.UpdateStmt, *ast.SelectStmt, *ast.ShowStmt, *ast.LoadDataStmt:
			if method != 0 && method != 2 {
				c := &models.Clause{
					Description: r.Message,
					Level:       r.Level,
				}
				s.Violations.Append(c)
			} else {
				method = 2
			}
		}
	}
}
