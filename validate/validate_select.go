package validate

import (
	"fmt"
	"sync"

	"github.com/mia0x75/parser/ast"
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/models"
)

// SelectVldr 查询语句审核
type SelectVldr struct {
	vldr

	sd *ast.SelectStmt
}

// Call 利用反射方法动态调用审核函数
func (v *SelectVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *SelectVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *SelectVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if sd, ok := node.(*ast.SelectStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.sd = sd
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// WithoutWhereNotAllowed 禁止没有WHERE的查询
// RULE: SEL-L2-001
func (v *SelectVldr) WithoutWhereNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if v.sd.Where == nil {
		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// WithoutLimitNotAllowed 禁止没有LIMIT的查询
// RULE: SEL-L2-002
func (v *SelectVldr) WithoutLimitNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if v.sd.Limit == nil {
		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// UseWildcardNotAllowed 禁止SELECT STAR
// RULE: SEL-L2-003
func (v *SelectVldr) UseWildcardNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, field := range v.sd.Fields.Fields {
		if field.WildCard != nil {
			c := &models.Clause{
				Description: r.Message,
				Level:       r.Level,
			}
			s.Violations.Append(c)
			break
		}
	}
}

// UseExplicitLockNotAllowed 禁止指定锁的类型
// RULE: SEL-L2-004
func (v *SelectVldr) UseExplicitLockNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if v.sd.LockTp != ast.SelectLockNone {
		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TargetDatabaseDoesNotExist 目标数据库必须已存在
// RULE: SEL-L3-001
func (v *SelectVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if v.DatabaseInfo(ti.Database) == nil {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, ti.Database),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// TargetTableDoesNotExist 目标表必须已存在
// RULE: SEL-L3-002
func (v *SelectVldr) TargetTableDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if ti.Table == nil {
			continue
		}
		if v.TableInfo(ti.Database, ti.Table.Name) == nil {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, fmt.Sprintf("`%s`.`%s`", ti.Database, ti.Table.Name)),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// TargetColumnDoesNotDoesNotExist 目标列必须已存在，从抽象语法数解出来所有的关联字段信息，包括FieldList/Where/OrderBy/GroupBy
// RULE: SEL-L3-003
func (v *SelectVldr) TargetColumnDoesNotDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, field := range v.sd.Fields.Fields {
		if field.WildCard != nil {
			// 如果是SELECT *或者SELECT alias.*则不处理
			continue
		}
		// TODO:
		fmt.Printf("Expr: %T, %+v\n", field.Expr, field.Expr)
	}
}

// ReturnBlobOrTextNotAllowed 是否允许返回BLOB/TEXT列
// RULE: SEL-L3-004
func (v *SelectVldr) ReturnBlobOrTextNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	fields := []*ast.SelectField{}
	for _, field := range v.sd.Fields.Fields {
		if field.WildCard != nil {
			// 如果是SELECT *或者SELECT alias.*则需要展开对应表的所有字段
			// TODO:
		}
		fields = append(fields, field)
	}
	for _, field := range fields {
		// TODO:
		fmt.Printf("Expr: %T, %+v\n", field.Expr, field.Expr)
	}
}
