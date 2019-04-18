package validate

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/mia0x75/parser/ast"

	"github.com/mia0x75/halo/models"
)

// IndexCreateVldr 创建索引语句相关的审核规则
type IndexCreateVldr struct {
	vldr

	ic *ast.CreateIndexStmt
}

// Call 利用反射方法动态调用审核函数
func (v *IndexCreateVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *IndexCreateVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *IndexCreateVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *IndexCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if ic, ok := node.(*ast.CreateIndexStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.ic = ic
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// IndexCreateIndexMaxColumnLimit 组合索引允许的最大列数
// RULE: CIX-L2-001
func (v *IndexCreateVldr) IndexCreateIndexMaxColumnLimit(s *models.Statement, r *models.Rule) {
	threshold, _ := strconv.Atoi(r.Values)
	if threshold < len(v.ic.IndexColNames) {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, v.ic.IndexName, threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// IndexCreateIndexNameQualified 索引名标识符规则
// RULE: CIX-L2-002
func (v *IndexCreateVldr) IndexCreateIndexNameQualified(s *models.Statement, r *models.Rule) {
	indexName := v.ic.IndexName

	if err := Match(r, indexName, indexName, r.Values); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// IndexCreateIndexNameLowerCaseRequired 索引名大小写规则
// RULE: CIX-L2-003
func (v *IndexCreateVldr) IndexCreateIndexNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	indexName := v.ic.IndexName

	if err := Match(r, indexName, indexName); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// IndexCreateIndexNameMaxLength 索引名长度规则
// RULE: CIX-L2-004
func (v *IndexCreateVldr) IndexCreateIndexNameMaxLength(s *models.Statement, r *models.Rule) {
	threshold, _ := strconv.Atoi(r.Values)
	if len(v.ic.IndexName) > threshold {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, v.ic.IndexName, threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
		return
	}
}

// IndexCreateIndexNamePrefixRequired 索引名前缀规则
// RULE: CIX-L2-005
func (v *IndexCreateVldr) IndexCreateIndexNamePrefixRequired(s *models.Statement, r *models.Rule) {
	indexName := v.ic.IndexName

	if err := Match(r, indexName, indexName, r.Values); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// IndexCreateDuplicateIndexColumn 组合索引中是否有重复列
// RULE: CIX-L2-006
func (v *IndexCreateVldr) IndexCreateDuplicateIndexColumn(s *models.Statement, r *models.Rule) {
	m := make(map[string]int)
	for _, k := range v.ic.IndexColNames {
		if _, ok := m[k.Column.Name.L]; ok {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, v.ic.IndexName),
				Level:       r.Level,
			}
			s.Violations.Append(c)
			return
		}
		m[k.Column.Name.L] = 1
	}
}

// IndexCreateTargetDatabaseExists 添加索引的表所属库是否存在
// RULE: CIX-L3-001
func (v *IndexCreateVldr) IndexCreateTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// IndexCreateTargetTableExists 条件索引的表是否存在
// RULE: CIX-L3-002
func (v *IndexCreateVldr) IndexCreateTargetTableExists(s *models.Statement, r *models.Rule) {
}

// IndexCreateTargetColumnExists 添加索引的列是否存在
// RULE: CIX-L3-003
func (v *IndexCreateVldr) IndexCreateTargetColumnExists(s *models.Statement, r *models.Rule) {
}

// IndexCreateTargetIndexExists 索引内容是否重复
// RULE: CIX-L3-004
func (v *IndexCreateVldr) IndexCreateTargetIndexExists(s *models.Statement, r *models.Rule) {
}

// IndexCreateTargetNameExists 索引名是否重复
// RULE: CIX-L3-005
func (v *IndexCreateVldr) IndexCreateTargetNameExists(s *models.Statement, r *models.Rule) {
}

// IndexCreateIndexCountLimit 最多能建多少个索引
// RULE: CIX-L3-006
func (v *IndexCreateVldr) IndexCreateIndexCountLimit(s *models.Statement, r *models.Rule) {
}

// IndexCreateIndexBlobColumnEnabled 是否允许在BLOB/TEXT列上建索引
// RULE: CIX-L3-007
func (v *IndexCreateVldr) IndexCreateIndexBlobColumnEnabled(s *models.Statement, r *models.Rule) {
}

// IndexDropVldr 删除索引语句相关的审核规则
type IndexDropVldr struct {
	vldr

	id *ast.DropIndexStmt
}

// Call 利用反射方法动态调用审核函数
func (v *IndexDropVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *IndexDropVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *IndexDropVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *IndexDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if id, ok := node.(*ast.DropIndexStmt); !ok {
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

// IndexDropTargetDatabaseExists 目标库是否存在
// RULE: RIX-L3-001
func (v *IndexDropVldr) IndexDropTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// IndexDropTargetTableExists 目标表是否存在
// RULE: RIX-L3-002
func (v *IndexDropVldr) IndexDropTargetTableExists(s *models.Statement, r *models.Rule) {
}

// IndexDropTargetIndexExists 目标索引是否存在
// RULE: RIX-L3-003
func (v *IndexDropVldr) IndexDropTargetIndexExists(s *models.Statement, r *models.Rule) {
}
