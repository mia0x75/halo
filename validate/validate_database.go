package validate

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/mia0x75/parser/ast"
	"github.com/mia0x75/parser/charset"

	"github.com/mia0x75/halo/models"
)

// DatabaseCreateVldr 创建数据库语句相关的审核规则
type DatabaseCreateVldr struct {
	vldr

	cd *ast.CreateDatabaseStmt
}

// Call 利用反射方法动态调用审核函数
func (v *DatabaseCreateVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *DatabaseCreateVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *DatabaseCreateVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *DatabaseCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if cd, ok := node.(*ast.CreateDatabaseStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.cd = cd
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// DatabaseCreateAvailableCharsets 建库允许的字符集
// RULE: CDB-L2-001
func (v *DatabaseCreateVldr) DatabaseCreateAvailableCharsets(s *models.Statement, r *models.Rule) {
	var charsets []string
	json.Unmarshal([]byte(r.Values), &charsets)
	var useCharset = "<empty>"
	valid := false
	for _, db := range v.cd.Options {
		if db.Tp == ast.DatabaseOptionCharset {
			useCharset = db.Value
			for _, charset := range charsets {
				if strings.EqualFold(useCharset, charset) {
					valid = true
					break
				}
			}
			break
		}
	}
	if !valid {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, useCharset, charsets),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// DatabaseCreateAvailableCollates 建库允许的排序规则
// RULE: CDB-L2-002
func (v *DatabaseCreateVldr) DatabaseCreateAvailableCollates(s *models.Statement, r *models.Rule) {
	var collates []string
	json.Unmarshal([]byte(r.Values), &collates)

	var useCollate = "<empty>"
	valid := false
	for _, db := range v.cd.Options {
		if db.Tp == ast.DatabaseOptionCollate {
			useCollate = db.Value
			for _, collate := range collates {
				if strings.EqualFold(useCollate, collate) {
					valid = true
					break
				}
			}
			break
		}
	}
	if !valid {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, useCollate, collates),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// DatabaseCreateCharsetCollateMatch 建库的字符集与排序规则必须匹配
// RULE: CDB-L2-003
func (v *DatabaseCreateVldr) DatabaseCreateCharsetCollateMatch(s *models.Statement, r *models.Rule) {
	useCharset := ""
	useCollate := ""
	for _, opt := range v.cd.Options {
		if opt.Tp == ast.DatabaseOptionCharset {
			useCharset = opt.Value
		}
		if opt.Tp == ast.DatabaseOptionCollate {
			useCollate = opt.Value
		}
	}

	if len(useCharset) == 0 || len(useCollate) == 0 {
		return
	}

	if !charset.ValidCharsetAndCollation(useCharset, useCollate) {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, useCharset, useCollate),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// DatabaseCreateDatabaseNameQualified 库名标识符规则
// RULE: CDB-L2-004
func (v *DatabaseCreateVldr) DatabaseCreateDatabaseNameQualified(s *models.Statement, r *models.Rule) {
	dbName := v.cd.Name
	if err := Match(r, dbName, dbName, r.Values); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// DatabaseCreateDatabaseNameLowerCaseRequired 库名大小写规则
// RULE: CDB-L2-005
func (v *DatabaseCreateVldr) DatabaseCreateDatabaseNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	dbName := v.cd.Name
	if err := Match(r, dbName, dbName); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// DatabaseCreateDatabaseNameMaxLength 库名长度规则
// RULE: CDB-L2-006
func (v *DatabaseCreateVldr) DatabaseCreateDatabaseNameMaxLength(s *models.Statement, r *models.Rule) {
	dbName := v.cd.Name
	threshold, _ := strconv.Atoi(r.Values)
	if len(dbName) > threshold {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, dbName, threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// DatabaseAlterVldr 修改数据库语句相关的审核规则
type DatabaseAlterVldr struct {
	vldr

	rd *ast.AlterDatabaseStmt
}

// Call 利用反射方法动态调用审核函数
func (v *DatabaseAlterVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *DatabaseAlterVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *DatabaseAlterVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *DatabaseAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if rd, ok := node.(*ast.AlterDatabaseStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.rd = rd
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// DatabaseAlterAvailableCharsets 改库允许的字符集
// RULE: MDB-L2-001
func (v *DatabaseAlterVldr) DatabaseAlterAvailableCharsets(s *models.Statement, r *models.Rule) {
	var charsets []string
	json.Unmarshal([]byte(r.Values), &charsets)

	var useCharset string
	valid := true
	for _, db := range v.rd.Options {
		switch db.Tp {
		case ast.DatabaseOptionCharset:
			valid = false
			useCharset = db.Value
			for _, charset := range charsets {
				if strings.EqualFold(useCharset, charset) {
					valid = true
					break
				}
			}
		}
	}

	if !valid {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, useCharset, charsets),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// DatabaseAlterAvailableCollates 改库允许的排序规则
// RULE: MDB-L2-002
func (v *DatabaseAlterVldr) DatabaseAlterAvailableCollates(s *models.Statement, r *models.Rule) {
	var collates []string
	json.Unmarshal([]byte(r.Values), &collates)

	var useCollate string
	valid := true
	for _, db := range v.rd.Options {
		switch db.Tp {
		case ast.DatabaseOptionCollate:
			valid = false
			useCollate = db.Value
			for _, collate := range collates {
				if strings.EqualFold(useCollate, collate) {
					valid = true
					break
				}
			}
		}
	}

	if !valid {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, useCollate, collates),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// DatabaseAlterCharsetCollateMatch 改库的字符集与排序规则必须匹配
// RULE: MDB-L2-003
func (v *DatabaseAlterVldr) DatabaseAlterCharsetCollateMatch(s *models.Statement, r *models.Rule) {
	useCharset := ""
	useCollate := ""
	for _, db := range v.rd.Options {
		switch db.Tp {
		case ast.DatabaseOptionCharset:
			useCharset = db.Value
		case ast.DatabaseOptionCollate:
			useCollate = db.Value
		}
	}

	if len(useCharset) == 0 && len(useCollate) == 0 {
		return
	}

	if len(useCharset) != 0 && len(useCollate) != 0 {
		if !charset.ValidCharsetAndCollation(useCharset, useCollate) {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, useCharset, useCollate),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	} else {
		// TODO:
		// 从数据库获取字符集和排序规则
	}
}

// DatabaseAlterTargetDatabaseExists 目标库不存在
// RULE: MDB-L3-001
func (v *DatabaseAlterVldr) DatabaseAlterTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// DatabaseDropVldr 删除数据库语句相关的审核规则
type DatabaseDropVldr struct {
	vldr

	dd *ast.DropDatabaseStmt
}

// Call 利用反射方法动态调用审核函数
func (v *DatabaseDropVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *DatabaseDropVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *DatabaseDropVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *DatabaseDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if dd, ok := node.(*ast.DropDatabaseStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.dd = dd
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// DatabaseDropTargetDatabaseExists 目标库不存在
// RULE: DDB-L3-001
func (v *DatabaseDropVldr) DatabaseDropTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}
