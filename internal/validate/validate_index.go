package validate

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/mia0x75/parser/ast"
	log "github.com/sirupsen/logrus"
)

// IndexCreateVldr 创建索引语句相关的审核规则
type IndexCreateVldr struct {
	vldr

	ci *ast.CreateIndexStmt
}

// Enabled 当前规则组是否生效
func (v *IndexCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *IndexCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if ci, ok := node.(*ast.CreateIndexStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.ci = ci
			v.Walk(v.ci)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// MaxAllowedIndexColumnCount 组合索引允许的最大列数
// RULE: CIX-L2-001
func (v *IndexCreateVldr) MaxAllowedIndexColumnCount(s *Subject, r *Rule) {
	threshold, _ := strconv.Atoi(r.Parameters)
	if threshold < len(v.ci.IndexColNames) {
		c := &Clause{
			Description: fmt.Sprintf(r.Message, v.ci.IndexName, threshold),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// IndexNameQualified 索引名标识符规则
// RULE: CIX-L2-002
func (v *IndexCreateVldr) IndexNameQualified(s *Subject, r *Rule) {
	indexName := v.ci.IndexName

	if err := Match(r, indexName, indexName, r.Parameters); err != nil {
		c := &Clause{
			Description: err.Error(),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// IndexNameLowerCaseRequired 索引名大小写规则
// RULE: CIX-L2-003
func (v *IndexCreateVldr) IndexNameLowerCaseRequired(s *Subject, r *Rule) {
	indexName := v.ci.IndexName

	if err := Match(r, indexName, indexName); err != nil {
		c := &Clause{
			Description: err.Error(),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// IndexNameMaxLength 索引名长度规则
// RULE: CIX-L2-004
func (v *IndexCreateVldr) IndexNameMaxLength(s *Subject, r *Rule) {
	threshold, _ := strconv.Atoi(r.Parameters)
	indexName := v.ci.IndexName
	if len(indexName) > threshold {
		c := &Clause{
			Description: fmt.Sprintf(r.Message, indexName, threshold),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
		return
	}
}

// IndexNamePrefixRequired 索引名前缀规则
// RULE: CIX-L2-005
func (v *IndexCreateVldr) IndexNamePrefixRequired(s *Subject, r *Rule) {
	indexName := v.ci.IndexName

	if err := Match(r, indexName, indexName, r.Parameters); err != nil {
		c := &Clause{
			Description: err.Error(),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// IndexColumnDuplicate 组合索引中是否有重复列
// RULE: CIX-L2-006
func (v *IndexCreateVldr) IndexColumnDuplicate(s *Subject, r *Rule) {
	m := make(map[string]int)
	for _, k := range v.ci.IndexColNames {
		if _, ok := m[k.Column.Name.L]; ok {
			c := &Clause{
				Description: fmt.Sprintf(r.Message, v.ci.IndexName),
				Level:       uint8(r.Level),
			}
			s.Violations.Append(c)
			return
		}
		m[k.Column.Name.L] = 1
	}
}

// TargetDatabaseDoesNotExist 添加索引的表所属库是否存在
// RULE: CIX-L3-001
func (v *IndexCreateVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if v.DatabaseInfo(ti.Database) == nil {
			c := &Clause{
				Description: fmt.Sprintf(r.Message, ti.Database),
				Level:       uint8(r.Level),
			}
			s.Violations.Append(c)
		}
	}
}

// TargetTableDoesNotExist 条件索引的表是否存在
// RULE: CIX-L3-002
func (v *IndexCreateVldr) TargetTableDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if ti.Table == nil {
			continue
		}
		if v.TableInfo(ti.Database, ti.Table.Name) == nil {
			c := &Clause{
				Description: fmt.Sprintf(r.Message, fmt.Sprintf("`%s`.`%s`", ti.Database, ti.Table.Name)),
				Level:       uint8(r.Level),
			}
			s.Violations.Append(c)
		}
	}
}

// TargetColumnDoesNotExist 添加索引的列是否存在
// RULE: CIX-L3-003
func (v *IndexCreateVldr) TargetColumnDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.ci.Table.Schema.O, v.ci.Table.Name.O)
	if ti == nil {
		return
	}
	for _, col := range v.ci.IndexColNames {
		if ti.GetColumn(col.Column.Name.O) == nil {
			c := &Clause{
				Description: fmt.Sprintf(r.Message, col.Column.Name.O),
				Level:       uint8(r.Level),
			}
			s.Violations.Append(c)
		}
	}
}

// IndexOverlayNotAllowed 索引内容是否重复
// RULE: CIX-L3-004
func (v *IndexCreateVldr) IndexOverlayNotAllowed(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.ci.Table.Schema.O, v.ci.Table.Name.O)
	if ti == nil {
		return
	}
	cols := indexColNames(v.ci.IndexColNames)
	for _, name := range sortedIndexNames(ti) {
		if name == v.ci.IndexName {
			continue
		}
		if indexOverlay(cols, ti.Indexes[name].Cols) {
			c := &Clause{
				Description: fmt.Sprintf(r.Message, v.ci.IndexName, name),
				Level:       uint8(r.Level),
			}
			s.Violations.Append(c)
		}
	}
}

// IndexNameDuplicate 索引名是否重复
// RULE: CIX-L3-005
func (v *IndexCreateVldr) IndexNameDuplicate(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.ci.Table.Schema.O, v.ci.Table.Name.O)
	if ti == nil {
		return
	}
	if v.IndexInfo(v.ci.Table.Schema.O, v.ci.Table.Name.O, v.ci.IndexName) != nil {
		c := &Clause{
			Description: fmt.Sprintf(r.Message, v.ci.IndexName, v.ci.Table.Name.O),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// MaxAllowedIndexCount 最多能建多少个索引
// RULE: CIX-L3-006
func (v *IndexCreateVldr) MaxAllowedIndexCount(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.ci.Table.Schema.O, v.ci.Table.Name.O)
	if ti == nil {
		return
	}
	// 索引数量
	threshold, err := strconv.Atoi(r.Parameters)
	if err != nil {
		return
	}

	if len(ti.Indexes)+1 > threshold {
		c := &Clause{
			Description: fmt.Sprintf(r.Message, threshold),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}

// IndexOnBlobColumnNotAllowed 是否允许在BLOB/TEXT列上建索引
// RULE: CIX-L3-007
func (v *IndexCreateVldr) IndexOnBlobColumnNotAllowed(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.ci.Table.Schema.O, v.ci.Table.Name.O)
	if ti == nil {
		return
	}
	for _, col := range v.ci.IndexColNames {
		ci := ti.GetColumn(col.Column.Name.O)
		if ci != nil {
			if ci.Type.IsText() || ci.Type.IsBlob() {
				c := &Clause{
					Description: fmt.Sprintf(r.Message, ci.Name),
					Level:       uint8(r.Level),
				}
				s.Violations.Append(c)
			}
		}
	}
}

// IndexDropVldr 删除索引语句相关的审核规则
type IndexDropVldr struct {
	vldr

	di *ast.DropIndexStmt
}

// Enabled 当前规则组是否生效
func (v *IndexDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *IndexDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Subjects {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if di, ok := node.(*ast.DropIndexStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.di = di
			v.Walk(v.di)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			Call(v, uint16(r.VldrGroup), r.Func, s, r)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: RIX-L3-001
func (v *IndexDropVldr) TargetDatabaseDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if v.DatabaseInfo(ti.Database) == nil {
			c := &Clause{
				Description: fmt.Sprintf(r.Message, ti.Database),
				Level:       uint8(r.Level),
			}
			s.Violations.Append(c)
		}
	}
}

// TargetTableDoesNotExist 目标表是否存在
// RULE: RIX-L3-002
func (v *IndexDropVldr) TargetTableDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, ti := range v.Vi {
		if ti.Table == nil {
			continue
		}
		if v.TableInfo(ti.Database, ti.Table.Name) == nil {
			c := &Clause{
				Description: fmt.Sprintf(r.Message, fmt.Sprintf("`%s`.`%s`", ti.Database, ti.Table.Name)),
				Level:       uint8(r.Level),
			}
			s.Violations.Append(c)
		}
	}
}

// TargetIndexDoesNotExist 目标索引是否存在
// RULE: RIX-L3-003
func (v *IndexDropVldr) TargetIndexDoesNotExist(s *Subject, r *Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.di.Table.Schema.O, v.di.Table.Name.O)
	if ti == nil {
		return
	}
	if v.IndexInfo(v.di.Table.Schema.O, v.di.Table.Name.O, v.di.IndexName) == nil {
		c := &Clause{
			Description: fmt.Sprintf(r.Message, v.di.IndexName),
			Level:       uint8(r.Level),
		}
		s.Violations.Append(c)
	}
}
