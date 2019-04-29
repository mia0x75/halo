package validate

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/mia0x75/parser/ast"
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/models"
)

// IndexCreateVldr 创建索引语句相关的审核规则
type IndexCreateVldr struct {
	vldr

	ci *ast.CreateIndexStmt
}

// Call 利用反射方法动态调用审核函数
func (v *IndexCreateVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *IndexCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *IndexCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Stmts {
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
			v.Call(r.Func, s, r)
		}
	}
}

// IndexMaxColumnLimit 组合索引允许的最大列数
// RULE: CIX-L2-001
func (v *IndexCreateVldr) IndexMaxColumnLimit(s *models.Statement, r *models.Rule) {
	threshold, _ := strconv.Atoi(r.Values)
	if threshold < len(v.ci.IndexColNames) {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, v.ci.IndexName, threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// IndexNameQualified 索引名标识符规则
// RULE: CIX-L2-002
func (v *IndexCreateVldr) IndexNameQualified(s *models.Statement, r *models.Rule) {
	indexName := v.ci.IndexName

	if err := Match(r, indexName, indexName, r.Values); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// IndexNameLowerCaseRequired 索引名大小写规则
// RULE: CIX-L2-003
func (v *IndexCreateVldr) IndexNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	indexName := v.ci.IndexName

	if err := Match(r, indexName, indexName); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// IndexNameMaxLength 索引名长度规则
// RULE: CIX-L2-004
func (v *IndexCreateVldr) IndexNameMaxLength(s *models.Statement, r *models.Rule) {
	threshold, _ := strconv.Atoi(r.Values)
	indexName := v.ci.IndexName
	if len(indexName) > threshold {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, indexName, threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
		return
	}
}

// IndexNamePrefixRequired 索引名前缀规则
// RULE: CIX-L2-005
func (v *IndexCreateVldr) IndexNamePrefixRequired(s *models.Statement, r *models.Rule) {
	indexName := v.ci.IndexName

	if err := Match(r, indexName, indexName, r.Values); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// DuplicateIndexColumn 组合索引中是否有重复列
// RULE: CIX-L2-006
func (v *IndexCreateVldr) DuplicateIndexColumn(s *models.Statement, r *models.Rule) {
	m := make(map[string]int)
	for _, k := range v.ci.IndexColNames {
		if _, ok := m[k.Column.Name.L]; ok {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, v.ci.IndexName),
				Level:       r.Level,
			}
			s.Violations.Append(c)
			return
		}
		m[k.Column.Name.L] = 1
	}
}

// TargetDatabaseDoesNotExist 添加索引的表所属库是否存在
// RULE: CIX-L3-001
func (v *IndexCreateVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
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

// TargetTableDoesNotExist 条件索引的表是否存在
// RULE: CIX-L3-002
func (v *IndexCreateVldr) TargetTableDoesNotExist(s *models.Statement, r *models.Rule) {
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

// TargetColumnDoesNotExist 添加索引的列是否存在
// RULE: CIX-L3-003
func (v *IndexCreateVldr) TargetColumnDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.ci.Table.Schema.O, v.ci.Table.Name.O)
	if ti == nil {
		return
	}
	for _, col := range v.ci.IndexColNames {
		if ti.GetColumn(col.Column.Name.O) == nil {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, col.Column.Name.O),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// IndexOverlay 索引内容是否重复
// RULE: CIX-L3-004
func (v *IndexCreateVldr) IndexOverlay(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
}

// IndexNameDuplicate 索引名是否重复
// RULE: CIX-L3-005
func (v *IndexCreateVldr) IndexNameDuplicate(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.ci.Table.Schema.O, v.ci.Table.Name.O)
	if ti == nil {
		return
	}
	if v.IndexInfo(v.ci.Table.Schema.O, v.ci.Table.Name.O, v.ci.IndexName) != nil {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, v.ci.IndexName, v.ci.Table.Name.O),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// IndexCountLimit 最多能建多少个索引
// RULE: CIX-L3-006
func (v *IndexCreateVldr) IndexCountLimit(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.ci.Table.Schema.O, v.ci.Table.Name.O)
	if ti == nil {
		return
	}
	// 索引数量
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}

	if len(ti.Indexes)+1 > threshold {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// IndexOnBlobColumnNotAllowed 是否允许在BLOB/TEXT列上建索引
// RULE: CIX-L3-007
func (v *IndexCreateVldr) IndexOnBlobColumnNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.ci.Table.Schema.O, v.ci.Table.Name.O)
	if ti == nil {
		return
	}
	for _, col := range v.ci.IndexColNames {
		ci := ti.GetColumn(col.Column.Name.O)
		if ci != nil {
			if ci.SQLType.IsText() || ci.SQLType.IsBlob() {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, ci.Name),
					Level:       r.Level,
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

// Call 利用反射方法动态调用审核函数
func (v *IndexDropVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *IndexDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *IndexDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Stmts {
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
			v.Call(r.Func, s, r)
		}
	}
}

// TargetDatabaseDoesNotExist 目标库是否存在
// RULE: RIX-L3-001
func (v *IndexDropVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
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

// TargetTableDoesNotExist 目标表是否存在
// RULE: RIX-L3-002
func (v *IndexDropVldr) TargetTableDoesNotExist(s *models.Statement, r *models.Rule) {
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

// TargetIndexDoesNotExist 目标索引是否存在
// RULE: RIX-L3-003
func (v *IndexDropVldr) TargetIndexDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.di.Table.Schema.O, v.di.Table.Name.O)
	if ti == nil {
		return
	}
	if v.IndexInfo(v.di.Table.Schema.O, v.di.Table.Name.O, v.di.IndexName) == nil {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, v.di.IndexName),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}
