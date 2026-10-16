package validate

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/mia0x75/parser/ast"
	"github.com/mia0x75/parser/charset"
	log "github.com/sirupsen/logrus"
)

// DatabaseCreateVldr 创建数据库语句相关的审核规则
type DatabaseCreateVldr struct {
	vldr

	cd *ast.CreateDatabaseStmt
}

// Enabled 当前规则组是否生效
func (v *DatabaseCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *DatabaseCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if cd, ok := node.(*ast.CreateDatabaseStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.cd = cd
			v.Walk(v.cd)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// AvailableCharsets 建库允许的字符集
// RULE: CDB-L2-001
func (v *DatabaseCreateVldr) AvailableCharsets(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	var charsets []string
	_ = json.Unmarshal([]byte(r.Values), &charsets)
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
		c := &Clause{
			Description: fmt.Sprintf(r.Message, useCharset, charsets),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// AvailableCollates 建库允许的排序规则
// RULE: CDB-L2-002
func (v *DatabaseCreateVldr) AvailableCollates(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	var collates []string
	_ = json.Unmarshal([]byte(r.Values), &collates)

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
		c := &Clause{
			Description: fmt.Sprintf(r.Message, useCollate, collates),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// CharsetCollateMustMatch 建库的字符集与排序规则必须匹配
// RULE: CDB-L2-003
func (v *DatabaseCreateVldr) CharsetCollateMustMatch(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
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
		c := &Clause{
			Description: fmt.Sprintf(r.Message, useCharset, useCollate),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// DatabaseNameQualified 库名标识符规则
// RULE: CDB-L2-004
func (v *DatabaseCreateVldr) DatabaseNameQualified(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	dbName := v.cd.Name
	if err := Match(r, dbName, dbName, r.Values); err != nil {
		c := &Clause{
			Description: err.Error(),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// DatabaseNameLowerCaseRequired 库名大小写规则
// RULE: CDB-L2-005
func (v *DatabaseCreateVldr) DatabaseNameLowerCaseRequired(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	dbName := v.cd.Name
	if err := Match(r, dbName, dbName); err != nil {
		c := &Clause{
			Description: err.Error(),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// DatabaseNameMaxLength 库名长度规则
// RULE: CDB-L2-006
func (v *DatabaseCreateVldr) DatabaseNameMaxLength(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	dbName := v.cd.Name
	threshold, _ := strconv.Atoi(r.Values)
	if len(dbName) > threshold {
		c := &Clause{
			Description: fmt.Sprintf(r.Message, dbName, threshold),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// TargetDatabaseDoesNotExist 目标数据库已经存在
// RULE: CDB-L2-007
func (v *DatabaseCreateVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	database := v.cd.Name
	if v.DatabaseInfo(database) != nil {
		c := &Clause{
			Description: fmt.Sprintf(r.Message, database),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// DatabaseAlterVldr 修改数据库语句相关的审核规则
type DatabaseAlterVldr struct {
	vldr

	ad *ast.AlterDatabaseStmt
}

// Enabled 当前规则组是否生效
func (v *DatabaseAlterVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *DatabaseAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if ad, ok := node.(*ast.AlterDatabaseStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.ad = ad
			v.Walk(v.ad)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// AvailableCharsets 改库允许的字符集
// RULE: MDB-L2-001
func (v *DatabaseAlterVldr) AvailableCharsets(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	var charsets []string
	_ = json.Unmarshal([]byte(r.Values), &charsets)

	var useCharset string
	valid := true
	for _, db := range v.ad.Options {
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
		c := &Clause{
			Description: fmt.Sprintf(r.Message, useCharset, charsets),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// AvailableCollates 改库允许的排序规则
// RULE: MDB-L2-002
func (v *DatabaseAlterVldr) AvailableCollates(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	var collates []string
	_ = json.Unmarshal([]byte(r.Values), &collates)

	var useCollate string
	valid := true
	for _, db := range v.ad.Options {
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
		c := &Clause{
			Description: fmt.Sprintf(r.Message, useCollate, collates),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// CharsetCollateMustMatch 改库的字符集与排序规则必须匹配
// RULE: MDB-L2-003
func (v *DatabaseAlterVldr) CharsetCollateMustMatch(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	useCharset := ""
	useCollate := ""
	for _, opt := range v.ad.Options {
		switch opt.Tp {
		case ast.DatabaseOptionCharset:
			useCharset = opt.Value
		case ast.DatabaseOptionCollate:
			useCollate = opt.Value
		}
	}

	di := v.DatabaseInfo(v.ad.Name)
	if di == nil {
		return
	}

	if len(useCharset) == 0 {
		useCharset = di.Charset
	}

	if len(useCollate) == 0 {
		useCollate = di.Collate
	}

	if !charset.ValidCharsetAndCollation(useCharset, useCollate) {
		c := &Clause{
			Description: fmt.Sprintf(r.Message, useCharset, useCollate),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// TargetDatabaseDoesNotExist 目标库不存在
// RULE: MDB-L2-004
func (v *DatabaseAlterVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	database := v.ad.Name
	if v.DatabaseInfo(database) == nil {
		c := &Clause{
			Description: fmt.Sprintf(r.Message, database),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// DatabaseDropVldr 删除数据库语句相关的审核规则
type DatabaseDropVldr struct {
	vldr

	dd *ast.DropDatabaseStmt
}

// Enabled 当前规则组是否生效
func (v *DatabaseDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *DatabaseDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if dd, ok := node.(*ast.DropDatabaseStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.dd = dd
			v.Walk(v.dd)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库不存在
// RULE: DDB-L2-001
func (v *DatabaseDropVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	database := v.dd.Name
	if v.DatabaseInfo(database) == nil {
		c := &Clause{
			Description: fmt.Sprintf(r.Message, database),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}
