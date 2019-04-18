package validate

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/mia0x75/parser/ast"

	"github.com/mia0x75/halo/models"
)

// InsertVldr 插入数据语句相关的审核规则
type InsertVldr struct {
	vldr

	id *ast.InsertStmt
}

// Call 利用反射方法动态调用审核函数
func (v *InsertVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *InsertVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *InsertVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *InsertVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if id, ok := node.(*ast.InsertStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.id = id
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// InsertExplicitColumnRequired 是否要求显式列申明
// RULE: INS-L2-001
func (v *InsertVldr) InsertExplicitColumnRequired(s *models.Statement, r *models.Rule) {
	if v.id.IsReplace {
		return
	}

	if v.id.Columns == nil || len(v.id.Columns) == 0 {
		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// InsertUsingSelectEnabled 是否允许INSERT...SELECT
// RULE: INS-L2-002
func (v *InsertVldr) InsertUsingSelectEnabled(s *models.Statement, r *models.Rule) {
	if v.id.IsReplace {
		return
	}

	if v.id.Select != nil {
		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// InsertRowsLimit 单语句允许操作的最大行数
// RULE: INS-L2-004
func (v *InsertVldr) InsertRowsLimit(s *models.Statement, r *models.Rule) {
	thredhold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	if len(v.id.Lists) > thredhold {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, thredhold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// InsertColumnValueMatch 列类型、值是否匹配
// RULE: INS-L2-005
func (v *InsertVldr) InsertColumnValueMatch(s *models.Statement, r *models.Rule) {
	columnNum := len(v.id.Columns)
	for _, list := range v.id.Lists {
		if len(list) != columnNum {
			c := &models.Clause{
				Description: r.Message,
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// InsertTargetDatabaseExists 目标库是否存在
// RULE: INS-L3-001
func (v *InsertVldr) InsertTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// InsertTargetTableExists 目标表是否存在
// RULE: INS-L3-002
func (v *InsertVldr) InsertTargetTableExists(s *models.Statement, r *models.Rule) {
}

// InsertTargetColumnExists 目标列是否存在
// RULE: INS-L3-003
func (v *InsertVldr) InsertTargetColumnExists(s *models.Statement, r *models.Rule) {
}

// InsertValueForNotNullColumnRequired 非空列是否有值
// RULE: INS-L3-004
func (v *InsertVldr) InsertValueForNotNullColumnRequired(s *models.Statement, r *models.Rule) {
}

// ReplaceVldr 插入数据语句相关的审核规则
type ReplaceVldr struct {
	vldr

	id *ast.InsertStmt
}

// Call 利用反射方法动态调用审核函数
func (v *ReplaceVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *ReplaceVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *ReplaceVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *ReplaceVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if id, ok := node.(*ast.InsertStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.id = id
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}
