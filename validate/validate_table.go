package validate

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/mia0x75/parser/ast"
	"github.com/mia0x75/parser/charset"
	"github.com/mia0x75/parser/mysql"
	"github.com/mia0x75/parser/types"
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/halo/models"
)

// TableCreateVldr 创建表语句相关的审核规则
type TableCreateVldr struct {
	vldr

	ct *ast.CreateTableStmt
}

// Call 利用反射方法动态调用审核函数
func (v *TableCreateVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *TableCreateVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *TableCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if ct, ok := node.(*ast.CreateTableStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.ct = ct
			v.Walk(v.ct)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// AvailableCharsets 建表允许的字符集
// RULE: CTB-L2-001
func (v *TableCreateVldr) AvailableCharsets(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// CREATE ... LIKE ... 或者 CREATE ... SELECT ...
	if v.ct.ReferTable != nil || v.ct.Select != nil {
		return
	}

	var charsets []string
	// 将字符串反解析为结构体
	json.Unmarshal([]byte(r.Values), &charsets)

	useCharset := "<empty>"
	valid := false
	for _, option := range v.ct.Options {
		if option.Tp == ast.TableOptionCharset {
			useCharset = option.StrValue
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

// AvailableCollates 建表允许的排序规则
// RULE: CTB-L2-002
func (v *TableCreateVldr) AvailableCollates(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// CREATE ... LIKE ... 或者 CREATE ... SELECT ...
	if v.ct.ReferTable != nil || v.ct.Select != nil {
		return
	}

	var collates []string
	// 将字符串反解析为结构体
	json.Unmarshal([]byte(r.Values), &collates)

	useCollate := "<empty>"
	valid := false
	for _, opt := range v.ct.Options {
		if opt.Tp == ast.TableOptionCollate {
			useCollate = opt.StrValue
			for _, collate := range collates {
				if strings.EqualFold(opt.StrValue, collate) {
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

// TableCharsetCollateMustMatch 建表是校验规则与字符集必须匹配
// RULE: CTB-L2-003
func (v *TableCreateVldr) TableCharsetCollateMustMatch(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// TODO 需要判断后续字符集与排序规则的关系,考虑读字典表
	// CREATE ... LIKE ... 或者 CREATE ... SELECT ...
	if v.ct.ReferTable != nil || v.ct.Select != nil {
		return
	}

	useCharset := ""
	useCollate := ""
	for _, opt := range v.ct.Options {
		if opt.Tp == ast.TableOptionCharset {
			useCharset = opt.StrValue
		}
		if opt.Tp == ast.TableOptionCollate {
			useCollate = opt.StrValue
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

// AvailableEngines 建表允许的存储引擎
// RULE: CTB-L2-004
func (v *TableCreateVldr) AvailableEngines(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// CREATE ... LIKE ... 或者 CREATE ... SELECT ...
	if v.ct.ReferTable != nil || v.ct.Select != nil {
		return
	}

	var engines []string
	json.Unmarshal([]byte(r.Values), &engines)

	useEngine := "<empty>"
	valid := false
	for _, table := range v.ct.Options {
		if table.Tp == ast.TableOptionEngine {
			useEngine = table.StrValue
			for _, engine := range engines {
				if strings.EqualFold(useEngine, engine) {
					valid = true
					break
				}
			}
		}
	}

	if !valid {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, useEngine, engines),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TableNameQualified 表名必须符合命名规范
// RULE: CTB-L2-005
func (v *TableCreateVldr) TableNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	tableName := v.ct.Table.Name.O
	if err := Match(r, tableName, tableName, r.Values); err != nil {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, tableName, r.Values),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TableNameLowerCaseRequired 表名是否允许大写
// RULE: CTB-L2-006
func (v *TableCreateVldr) TableNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	tableName := strings.TrimSpace(v.ct.Table.Name.O)
	if err := Match(r, tableName, tableName); err != nil {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, tableName),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TableNameMaxLength 表名最大长度
// RULE: CTB-L2-007
func (v *TableCreateVldr) TableNameMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	tableName := strings.TrimSpace(v.ct.Table.Name.O)
	if len(tableName) > threshold {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, tableName, threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TableCommentRequired 表必须有注释
// RULE: CTB-L2-008
func (v *TableCreateVldr) TableCommentRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// CREATE ... LIKE ... 或者 CREATE ... SELECT ...
	if v.ct.ReferTable != nil || v.ct.Select != nil {
		return
	}

	valid := false
	tableName := strings.TrimSpace(v.ct.Table.Name.O)
	for _, op := range v.ct.Options {
		if op.Tp == ast.TableOptionComment {
			// 判断是否为非空注释
			if len(strings.TrimSpace(op.StrValue)) > 0 {
				valid = true
			}
			break
		}
	}
	if !valid {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, tableName),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// CreateTableFromSelectNotAllowed 是否允许查询语句建表
// RULE: CTB-L2-009
func (v *TableCreateVldr) CreateTableFromSelectNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// 规则允许执行，表示禁止新建表
	if v.ct.Select == nil {
		return
	}
	if _, ok := v.ct.Select.(*ast.SelectStmt); ok {
		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// ColumnNameQualified 列名必须符合命名规范
// RULE: CTB-L2-010
func (v *TableCreateVldr) ColumnNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, col := range v.ct.Cols {
		colName := col.Name.Name.O
		if len(strings.TrimSpace(colName)) == 0 {
			colName = "<empty>"
		}
		if err := Match(r, colName, colName, r.Values); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// ColumnNameLowerCaseRequired 列名是否允许大写
// RULE: CTB-L2-011
func (v *TableCreateVldr) ColumnNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, col := range v.ct.Cols {
		colName := strings.TrimSpace(col.Name.Name.O)
		if len(colName) == 0 {
			continue
		}
		if err := Match(r, colName, colName); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// ColumnNameMaxLength 列名最大长度
// RULE: CTB-L2-012
func (v *TableCreateVldr) ColumnNameMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	for _, col := range v.ct.Cols {
		colName := strings.TrimSpace(col.Name.Name.O)
		if len(colName) > threshold {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, colName, threshold),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// ColumnNameDuplicate 列名是否重复
// RULE: CTB-L2-013
func (v *TableCreateVldr) ColumnNameDuplicate(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	tableName := v.ct.Table.Name.O
	// 字段名
	columnMap := map[string]int{}
	for _, col := range v.ct.Cols {
		colName := col.Name.Name.L // 不关心大小写
		if _, ok := columnMap[colName]; ok {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, tableName, colName),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		} else {
			columnMap[colName] = 1
		}
	}
}

// MaxAllowedColumnCount 表允许的最大列数
// RULE: CTB-L2-014
func (v *TableCreateVldr) MaxAllowedColumnCount(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	tableName := v.ct.Table.Name.O
	// 字段数量
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}

	if len(v.ct.Cols) > threshold {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, tableName, len(v.ct.Cols), threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// ColumnTypesDoesNotExpect 列不允许的数据类型
// RULE: CTB-L2-015
func (v *TableCreateVldr) ColumnTypesDoesNotExpect(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	var availableTypes []string
	json.Unmarshal([]byte(r.Values), &availableTypes)

	if len(availableTypes) == 0 {
		return
	}

	for _, col := range v.ct.Cols {
		colType := types.TypeStr(col.Tp.Tp)
		for _, t := range availableTypes {
			if strings.EqualFold(colType, t) {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, col.Name, colType, availableTypes),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ColumnCommentRequired 列必须有注释
// RULE: CTB-L2-016
func (v *TableCreateVldr) ColumnCommentRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, col := range v.ct.Cols {
		// 排除含有注释字段
		hasComment := false
		for _, op := range col.Options {
			if op.Tp == ast.ColumnOptionComment {
				comment := op.Expr.(ast.ValueExpr).GetValue()
				if comment.(string) != "" {
					hasComment = true
				}
				break
			}
		}
		// 格式化输出错误信息
		if !hasComment {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, col.Name.Name.O),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// ColumnAvailableCharsets 列允许的字符集
// RULE: CTB-L2-017
func (v *TableCreateVldr) ColumnAvailableCharsets(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// 字符集
	var charsets []string
	json.Unmarshal([]byte(r.Values), &charsets)

	for _, col := range v.ct.Cols {
		colCharset := col.Tp.Charset
		// 字符集允许隐式声明
		if len(colCharset) == 0 {
			continue
		}

		valid := false
		for _, charset := range charsets {
			if strings.EqualFold(charset, colCharset) {
				valid = true
				break
			}
		}
		if !valid {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, col.Name, colCharset, charsets),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// ColumnAvailableCollates 列允许的排序规则
// RULE: CTB-L2-018
func (v *TableCreateVldr) ColumnAvailableCollates(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// 排序规则
	var collates []string
	json.Unmarshal([]byte(r.Values), &collates)

	for _, col := range v.ct.Cols {
		colCollate := col.Tp.Collate
		// 排序规则允许隐式声明
		if len(colCollate) == 0 {
			continue
		}
		valid := false
		for _, collate := range collates {
			if strings.EqualFold(collate, colCollate) {
				valid = true
				break
			}
		}
		if !valid {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, col.Name, colCollate, collates),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// ColumnCharsetCollateMustMatch 列的字符集与排序规则必须匹配
// RULE: CTB-L2-019
func (v *TableCreateVldr) ColumnCharsetCollateMustMatch(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, col := range v.ct.Cols {
		// TODO 需要判断后续字符集与排序规则的关系,考虑读字典表
		colCharset := col.Tp.Charset
		colCollate := col.Tp.Collate

		if len(colCollate) == 0 || len(colCharset) == 0 {
			continue
		}

		if len(colCollate) > 0 && len(colCharset) > 0 {
			if !charset.ValidCharsetAndCollation(colCharset, colCollate) {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, col.Name, colCharset, colCollate),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ColumnNotNullWithDefaultRequired 非空列是否有默认值
// RULE: CTB-L2-020
func (v *TableCreateVldr) ColumnNotNullWithDefaultRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// 非空字段 默认值字段
	nameMap := make(map[string]int)

	for _, col := range v.ct.Cols {
		for _, op := range col.Options {
			if op.Tp == ast.ColumnOptionNotNull {
				nameMap[col.Name.Name.O] = 1 // 只会有一个该选项
				break
			}
		}
	}
	for _, col := range v.ct.Cols {
		for _, op := range col.Options {
			switch op.Tp {
			case ast.ColumnOptionDefaultValue:
				if _, ok := nameMap[col.Name.Name.O]; ok {
					nameMap[col.Name.Name.O]++
				}
			case ast.ColumnOptionAutoIncrement:
				if _, ok := nameMap[col.Name.Name.O]; ok {
					nameMap[col.Name.Name.O]++
				}
			}
		}
	}

	for name, val := range nameMap {
		if val == 1 {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, name),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// ColumnAutoIncAvailableTypes 自增列允许的数据类型
// RULE: CTB-L2-021
func (v *TableCreateVldr) ColumnAutoIncAvailableTypes(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	var availableTypes []string
	json.Unmarshal([]byte(r.Values), &availableTypes)

	for _, col := range v.ct.Cols {
		hasAutoInc := false
		for _, op := range col.Options {
			// 只有一个自增列，需要优化退出逻辑
			if op.Tp == ast.ColumnOptionAutoIncrement {
				hasAutoInc = true
				valid := false
				colType := types.TypeStr(col.Tp.Tp)
				for _, t := range availableTypes {
					if strings.EqualFold(colType, t) {
						valid = true
						break
					}
				}
				if !valid {
					c := &models.Clause{
						Description: fmt.Sprintf(r.Message, col.Name, colType, availableTypes),
						Level:       r.Level,
					}
					s.Violations.Append(c)
				}
				break
			}
		}
		if hasAutoInc {
			break
		}
	}
}

// ColumnAutoIncUnsignedRequired 自增列必须是无符号
// RULE: CTB-L2-022
func (v *TableCreateVldr) ColumnAutoIncUnsignedRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, col := range v.ct.Cols {
		hasAutoInc := false
		for _, op := range col.Options {
			if op.Tp == ast.ColumnOptionAutoIncrement {
				hasAutoInc = true
				if !mysql.HasUnsignedFlag(col.Tp.Flag) {
					c := &models.Clause{
						Description: fmt.Sprintf(r.Message, col.Name),
						Level:       r.Level,
					}
					s.Violations.Append(c)
				}
				break
			}
		}
		if hasAutoInc {
			break
		}
	}
}

// ColumnAutoIncMustPrimaryKey 自增列必须是主键
// RULE: CTB-L2-023
func (v *TableCreateVldr) ColumnAutoIncMustPrimaryKey(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, col := range v.ct.Cols {
		for _, op := range col.Options {
			if op.Tp == ast.ColumnOptionAutoIncrement {
				if mysql.HasPriKeyFlag(col.Tp.Flag) {
					return
				}
				for _, csts := range v.ct.Constraints {
					if csts.Tp == ast.ConstraintPrimaryKey {
						if len(csts.Keys) == 1 && csts.Keys[0].Column.String() == col.Name.String() {
							return
						}
					}
				}
			}
		}
	}

	autoIncMap := map[string]int{}
	pkMap := map[string]int{}
	for _, col := range v.ct.Cols {
		for _, op := range col.Options {
			switch op.Tp {
			case ast.ColumnOptionAutoIncrement:
				autoIncMap[col.Name.Name.String()] = 1
			case ast.ColumnOptionPrimaryKey:
				pkMap[col.Name.Name.String()] = 1
			}
		}
	}
	if len(autoIncMap) == 0 {
		return
	}
	for colName := range autoIncMap {
		if pkMap[colName] != 1 {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, colName),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// MaxAllowedTimestampCount 仅允许一个时间戳类型的列
// RULE: CTB-L2-024
func (v *TableCreateVldr) MaxAllowedTimestampCount(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	tableName := v.ct.Table.Name.O
	count := 0
	for _, col := range v.ct.Cols {
		colType := types.TypeStr(col.Tp.Tp)
		if strings.EqualFold(colType, "timestamp") {
			count++
		}
	}

	if count > 1 {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, tableName),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// MaxAllowedIndexColumnCount 单一索引最大列数
// RULE: CTB-L2-025
func (v *TableCreateVldr) MaxAllowedIndexColumnCount(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	threshold, _ := strconv.Atoi(r.Values)
	for _, c := range v.ct.Constraints {
		switch c.Tp {
		case ast.ConstraintIndex, ast.ConstraintPrimaryKey, ast.ConstraintKey, ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
			if len(c.Keys) > threshold {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, c.Name, threshold),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// PrimaryKeyRequired 必须有主键
// RULE: CTB-L2-026
func (v *TableCreateVldr) PrimaryKeyRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// CREATE ... LIKE ... 或者 CREATE ... SELECT ...
	if v.ct.ReferTable != nil || v.ct.Select != nil {
		return
	}

	for _, col := range v.ct.Cols {
		for _, op := range col.Options {
			if op.Tp == ast.ColumnOptionPrimaryKey {
				return
			}
		}
	}
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintPrimaryKey {
			return
		}
	}

	c := &models.Clause{
		Description: r.Message,
		Level:       r.Level,
	}
	s.Violations.Append(c)
}

// PrimaryKeyNameExplicit 主键是否显式命名
// RULE: CTB-L2-027
func (v *TableCreateVldr) PrimaryKeyNameExplicit(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	tableName := v.ct.Table.Name.O
	valid := false

	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintPrimaryKey {
			// 显式命名
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) > 0 {
				valid = true
				break
			}
		}
	}
	if !valid {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, tableName),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// PrimaryKeyNameQualified 主键名标识符规则
// RULE: CTB-L2-028
func (v *TableCreateVldr) PrimaryKeyNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if len(v.ct.Constraints) == 0 {
		return
	}
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintPrimaryKey {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				continue
			}
			if err := Match(r, keyName, keyName, r.Values); err != nil {
				c := &models.Clause{
					Description: err.Error(),
					Level:       r.Level,
				}
				s.Violations.Append(c)
				// TODO:
				return
			}
		}
	}
}

// PrimryKeyLowerCaseRequired 主键名大小写规则
// RULE: CTB-L2-029
func (v *TableCreateVldr) PrimryKeyLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintPrimaryKey {
			if len(c.Name) > 0 {
				if err := Match(r, c.Name, c.Name); err != nil {
					c := &models.Clause{
						Description: err.Error(),
						Level:       r.Level,
					}
					s.Violations.Append(c)
				}
			}
			break
		}
	}
}

// PrimryKeyMaxLength 主键名长度规则
// RULE: CTB-L2-030
func (v *TableCreateVldr) PrimryKeyMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}

	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintPrimaryKey {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) > threshold {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, keyName, threshold),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
			break
		}
	}
}

// PrimryKeyPrefixRequired 主键名前缀规则
// RULE: CTB-L2-031
func (v *TableCreateVldr) PrimryKeyPrefixRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	if len(v.ct.Constraints) == 0 {
		return
	}
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintPrimaryKey {
			if len(c.Name) > 0 {
				if err := Match(r, c.Name, r.Values); err != nil {
					c := &models.Clause{
						Description: fmt.Sprintf(r.Message, c.Name, r.Values),
						Level:       r.Level,
					}
					s.Violations.Append(c)
				}
			}
		}
	}
}

// IndexNameExplicit 索引必须命名
// RULE: CTB-L2-032
func (v *TableCreateVldr) IndexNameExplicit(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintIndex {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				c := &models.Clause{
					Description: r.Message,
					Level:       r.Level,
				}
				s.Violations.Append(c)
				break
			}
		}
	}
}

// IndexNameQualified 索引名标识符规则
// RULE: CTB-L2-033
func (v *TableCreateVldr) IndexNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintIndex {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				continue
			}
			if err := Match(r, keyName, keyName, r.Values); err != nil {
				c := &models.Clause{
					Description: err.Error(),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// IndexNameLowerCaseRequired 索引名大小写规则
// RULE: CTB-L2-034
func (v *TableCreateVldr) IndexNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintIndex {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				continue
			}
			if err := Match(r, keyName, keyName); err != nil {
				c := &models.Clause{
					Description: err.Error(),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// IndexNameMaxLength 索引名长度规则
// RULE: CTB-L2-035
func (v *TableCreateVldr) IndexNameMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintIndex {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) > threshold {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, keyName, threshold),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// IndexNamePrefixRequired 索引名前缀规则
// RULE: CTB-L2-036
func (v *TableCreateVldr) IndexNamePrefixRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintIndex {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				continue
			}
			if err := Match(r, keyName, keyName, r.Values); err != nil {
				c := &models.Clause{
					Description: err.Error(),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// UniqueNameExplicit 唯一索引必须命名
// RULE: CTB-L2-037
func (v *TableCreateVldr) UniqueNameExplicit(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintUniq {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				c := &models.Clause{
					Description: r.Message,
					Level:       r.Level,
				}
				s.Violations.Append(c)
				break
			}
		}
	}
}

// UniqueNameQualified 唯一索引索名标识符规则
// RULE: CTB-L2-038
func (v *TableCreateVldr) UniqueNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintUniq {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				continue
			}
			if err := Match(r, keyName, keyName, r.Values); err != nil {
				c := &models.Clause{
					Description: err.Error(),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// UniqueNameLowerCaseRequired 唯一索引名大小写规则
// RULE: CTB-L2-039
func (v *TableCreateVldr) UniqueNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintUniq {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				continue
			}
			if err := Match(r, keyName, keyName); err != nil {
				c := &models.Clause{
					Description: err.Error(),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// UniqueNameMaxLength 唯一索引名长度规则
// RULE: CTB-L2-040
func (v *TableCreateVldr) UniqueNameMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintUniq {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) > threshold {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, keyName, threshold),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// UniqueNamePrefixRequired 唯一索引名前缀规则
// RULE: CTB-L2-041
func (v *TableCreateVldr) UniqueNamePrefixRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintUniq {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				continue
			}
			if err := Match(r, keyName, keyName, r.Values); err != nil {
				c := &models.Clause{
					Description: err.Error(),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ForeignKeyNotAllowed 是否允许外键
// RULE: CTB-L2-042
func (v *TableCreateVldr) ForeignKeyNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	valid := true
	for i := 0; i < len(v.ct.Cols) && valid; i++ {
		col := v.ct.Cols[i]
		for _, op := range col.Options {
			if op.Tp == ast.ColumnOptionReference {
				valid = false
				break
			}
		}
	}

	for _, constraint := range v.ct.Constraints {
		if constraint.Tp == ast.ConstraintForeignKey {
			valid = false
			break
		}
	}

	if !valid {
		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// ForeignKeyNameExplicit 是否配置外键名称
// RULE: CTB-L2-043
func (v *TableCreateVldr) ForeignKeyNameExplicit(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintForeignKey {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				c := &models.Clause{
					Description: r.Message,
					Level:       r.Level,
				}
				s.Violations.Append(c)
				break
			}
		}
	}
}

// ForeignKeyNameQualified 外键名标识符规则
// RULE: CTB-L2-044
func (v *TableCreateVldr) ForeignKeyNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintForeignKey {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				keyName = "<empty>"
			}
			if err := Match(r, keyName, keyName, r.Values); err != nil {
				c := &models.Clause{
					Description: err.Error(),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ForeignKeyNameLowerCaseRequired 外键名大小写规则
// RULE: CTB-L2-045
func (v *TableCreateVldr) ForeignKeyNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintForeignKey {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				continue
			}
			if err := Match(r, keyName, keyName); err != nil {
				c := &models.Clause{
					Description: err.Error(),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ForeignKeyNameMaxLength 外键名长度规则
// RULE: CTB-L2-046
func (v *TableCreateVldr) ForeignKeyNameMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintForeignKey {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) > threshold {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, keyName, threshold),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ForeignKeyNamePrefixRequired 外键名前缀规则
// RULE: CTB-L2-047
func (v *TableCreateVldr) ForeignKeyNamePrefixRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, c := range v.ct.Constraints {
		if c.Tp == ast.ConstraintForeignKey {
			keyName := strings.TrimSpace(c.Name)
			if len(keyName) == 0 {
				continue
			}
			if err := Match(r, keyName, keyName, r.Values); err != nil {
				c := &models.Clause{
					Description: err.Error(),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// MaxAllowedIndexCount 表中最多可建多少个索引
// RULE: CTB-L2-048
func (v *TableCreateVldr) MaxAllowedIndexCount(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	tableName := v.ct.Table.Name.O
	count := 0
	// 索引字段集
	for _, c := range v.ct.Constraints {
		switch c.Tp {
		case ast.ConstraintIndex:
			fallthrough
		case ast.ConstraintPrimaryKey:
			fallthrough
		case ast.ConstraintKey:
			fallthrough
		case ast.ConstraintUniq:
			fallthrough
		case ast.ConstraintUniqKey:
			fallthrough
		case ast.ConstraintUniqIndex:
			count++
		}
	}

	for _, col := range v.ct.Cols {
		for _, op := range col.Options {
			if op.Tp == ast.ColumnOptionPrimaryKey || op.Tp == ast.ColumnOptionUniqKey {
				count++
				break
			}
		}
	}
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	if count > threshold {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, tableName, count, threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// CreateTableUseLikeNotAllowed 禁止允许LIKE方式建表
// RULE: CTB-L2-049
func (v *TableCreateVldr) CreateTableUseLikeNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// 规则允许执行，表示禁止新建表
	if v.ct.ReferTable != nil {
		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// AutoIncColumnDuplicate 只允许一个自增列
// RULE: CTB-L2-050
func (v *TableCreateVldr) AutoIncColumnDuplicate(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	count := 0
	for _, col := range v.ct.Cols {
		for _, op := range col.Options {
			if op.Tp == ast.ColumnOptionAutoIncrement {
				count++
				break
			}
		}
	}
	if count > 1 {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, v.ct.Table.Name.O),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// PrimaryKeyDuplicate 只允许一个主键
// RULE: CTB-L2-051
func (v *TableCreateVldr) PrimaryKeyDuplicate(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	count := 0
	for _, col := range v.ct.Cols {
		for _, op := range col.Options {
			if op.Tp == ast.ColumnOptionPrimaryKey {
				count++
				break
			}
		}
	}
	for _, op := range v.ct.Constraints {
		if op.Tp == ast.ConstraintPrimaryKey {
			count++
		}
	}
	if count > 1 {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, v.ct.Table.Name.O),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TargetDatabaseDoesNotExist 建表时检查目标库是否存在
// RULE: CTB-L3-001
func (v *TableCreateVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	database := v.ct.Table.Schema.O
	if len(database) == 0 {
		database = v.Ctx.Ticket.Database
	}
	if v.DatabaseInfo(database) == nil {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, database),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TargetTableDoesNotExist 建表时检查表是否已经存在
// RULE: CTB-L3-002
func (v *TableCreateVldr) TargetTableDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	database := v.ct.Table.Schema.O
	if len(database) == 0 {
		database = v.Ctx.Ticket.Database
	}
	if nil != v.TableInfo(database, v.ct.Table.Name.O) {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, fmt.Sprintf("`%s`.`%s`", database, v.ct.Table.Name.O)),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TableAlterVldr 修改表语句相关的审核规则
type TableAlterVldr struct {
	vldr

	at *ast.AlterTableStmt
}

// Call 利用反射方法动态调用审核函数
func (v *TableAlterVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *TableAlterVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *TableAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if at, ok := node.(*ast.AlterTableStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.at = at
			v.Walk(v.at)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// AvailableCharsets 改表允许的字符集
// RULE: MTB-L2-001
func (v *TableAlterVldr) AvailableCharsets(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	var charsets []string
	// 将字符串反解析为结构体
	json.Unmarshal([]byte(r.Values), &charsets)

	useCharset := "<empty>"
	valid := true
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableOption {
			continue
		}

		for _, option := range spec.Options {
			if option.Tp == ast.TableOptionCharset {
				useCharset = option.StrValue
				valid = false
				for _, charset := range charsets {
					if strings.EqualFold(useCharset, charset) {
						valid = true
					}
				}
				break
			}
		}
		break
	}

	if !valid {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, useCharset, charsets),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// AvailableCollates 改表允许的校验规则
// RULE: MTB-L2-002
func (v *TableAlterVldr) AvailableCollates(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	var collates []string
	// 将字符串反解析为结构体
	json.Unmarshal([]byte(r.Values), &collates)

	useCollate := "<empty>"
	valid := true
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableOption {
			continue
		}

		for _, option := range spec.Options {
			if option.Tp == ast.TableOptionCollate {
				useCollate = option.StrValue
				valid = false
				for _, collate := range collates {
					if strings.EqualFold(useCollate, collate) {
						valid = true
						break
					}
				}
				break
			}
		}
		break
	}

	if !valid {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, useCollate, collates),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TableCharsetCollateMustMatch 表的字符集与排序规则必须匹配
// RULE: MTB-L2-003
func (v *TableAlterVldr) TableCharsetCollateMustMatch(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	useCharset := ""
	useCollate := ""
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableOption {
			continue
		}

		for _, opt := range spec.Options {
			if opt.Tp == ast.TableOptionCharset {
				useCharset = opt.StrValue
			}
			if opt.Tp == ast.TableOptionCollate {
				useCollate = opt.StrValue
			}
		}
		break
	}

	if len(useCharset) == 0 && len(useCollate) == 0 {
		return
	}
	if len(useCharset) == 0 {
		useCharset = ti.Charset
	}
	if len(useCollate) == 0 {
		useCollate = ti.Collate
	}
	if !charset.ValidCharsetAndCollation(useCharset, useCollate) {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, useCharset, useCollate),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// AvailableEngines 改表允许的存储引擎
// RULE: MTB-L2-004
func (v *TableAlterVldr) AvailableEngines(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	var engines []string
	// 将字符串反解析为结构体
	json.Unmarshal([]byte(r.Values), &engines)

	useEngine := "<empty>"
	valid := true
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableOption {
			continue
		}

		for _, option := range spec.Options {
			if option.Tp == ast.TableOptionEngine {
				useEngine = option.StrValue
				valid = false
				for _, engine := range engines {
					if strings.EqualFold(useEngine, engine) {
						valid = true
						break
					}
				}
				break
			}
		}
		break
	}

	if !valid {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, useEngine, engines),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// CommentRequired 如果修改表的注解信息，那么不可以为空
// RULE: MTB-L2-005
func (v *TableAlterVldr) CommentRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableOption {
			continue
		}
		for _, opt := range spec.Options {
			if opt.Tp != ast.TableOptionComment {
				continue
			}
			if strings.TrimSpace(opt.StrValue) == "" {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ColumnNameQualified 列名必须符合命名规范
// RULE: MTB-L2-005
func (v *TableAlterVldr) ColumnNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		// AlterTableModifyColumn 更改已有的列，不处理
		if spec.Tp != ast.AlterTableAddColumns &&
			spec.Tp != ast.AlterTableChangeColumn {
			continue
		}

		for _, col := range spec.NewColumns {
			colName := col.Name.Name.O
			if len(strings.TrimSpace(colName)) == 0 {
				colName = "<empty>"
			}
			if err := Match(r, colName, colName, r.Values); err != nil {
				c := &models.Clause{
					Description: err.Error(),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ColumnNameLowerCaseRequired 列名必须小写
// RULE: MTB-L2-006
func (v *TableAlterVldr) ColumnNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		// AlterTableModifyColumn 更改已有的列，不处理
		if spec.Tp != ast.AlterTableAddColumns &&
			spec.Tp != ast.AlterTableChangeColumn {
			continue
		}

		for _, col := range spec.NewColumns {
			colName := strings.TrimSpace(col.Name.Name.O)
			if len(colName) == 0 {
				continue
			}
			if err := Match(r, colName, colName); err != nil {
				c := &models.Clause{
					Description: err.Error(),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ColumnNameMaxLength 列名最大长度
// RULE: MTB-L2-007
func (v *TableAlterVldr) ColumnNameMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	for _, spec := range v.at.Specs {
		// AlterTableModifyColumn 更改已有的列，不处理
		if spec.Tp != ast.AlterTableAddColumns &&
			spec.Tp != ast.AlterTableChangeColumn {
			continue
		}

		for _, col := range spec.NewColumns {
			colName := strings.TrimSpace(col.Name.Name.O)
			if len(colName) > threshold {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, colName, threshold),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ColumnUnwantedTypes 列允许的数据类型
// RULE: MTB-L2-008
func (v *TableAlterVldr) ColumnUnwantedTypes(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	var availableTypes []string
	json.Unmarshal([]byte(r.Values), &availableTypes)

	if len(availableTypes) == 0 {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddColumns &&
			spec.Tp != ast.AlterTableChangeColumn &&
			spec.Tp != ast.AlterTableModifyColumn {
			continue
		}
		for _, col := range spec.NewColumns {
			colType := types.TypeStr(col.Tp.Tp)
			for _, t := range availableTypes {
				if strings.EqualFold(colType, t) {
					c := &models.Clause{
						Description: fmt.Sprintf(r.Message, col.Name, colType, availableTypes),
						Level:       r.Level,
					}
					s.Violations.Append(c)
				}
			}
		}
	}
}

// ColumnCommentRequired 列必须有注释
// RULE: MTB-L2-009
func (v *TableAlterVldr) ColumnCommentRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddColumns &&
			spec.Tp != ast.AlterTableChangeColumn &&
			spec.Tp != ast.AlterTableModifyColumn {
			continue
		}
		for _, col := range spec.NewColumns {
			// 排除含有注释字段
			hasComment := false
			for _, op := range col.Options {
				if op.Tp == ast.ColumnOptionComment {
					comment := op.Expr.(ast.ValueExpr).GetValue()
					if strings.TrimSpace(comment.(string)) != "" {
						hasComment = true
					}
					break
				}
			}
			// 格式化输出错误信息
			if !hasComment {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, col.Name.Name.O),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ColumnAvailableCharsets 列允许的字符集
// RULE: MTB-L2-010
func (v *TableAlterVldr) ColumnAvailableCharsets(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	// 字符集
	var charsets []string
	json.Unmarshal([]byte(r.Values), &charsets)

	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddColumns &&
			spec.Tp != ast.AlterTableChangeColumn &&
			spec.Tp != ast.AlterTableModifyColumn {
			continue
		}
		for _, col := range spec.NewColumns {
			colCharset := col.Tp.Charset
			// 字符集允许隐式声明
			if len(colCharset) == 0 {
				continue
			}
			valid := false
			for _, charset := range charsets {
				if strings.EqualFold(charset, colCharset) {
					valid = true
					break
				}
			}
			if !valid {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, col.Name, colCharset, charsets),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ColumnAvailableCollates 列允许的排序规则
// RULE: MTB-L2-011
func (v *TableAlterVldr) ColumnAvailableCollates(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	// 排序规则
	var collates []string
	json.Unmarshal([]byte(r.Values), &collates)
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddColumns &&
			spec.Tp != ast.AlterTableChangeColumn &&
			spec.Tp != ast.AlterTableModifyColumn {
			continue
		}
		for _, col := range spec.NewColumns {
			colCollate := col.Tp.Collate

			// 排序规则允许隐式声明
			if len(colCollate) == 0 {
				continue
			}
			valid := false
			for _, collate := range collates {
				if strings.EqualFold(collate, colCollate) {
					valid = true
					break
				}
			}
			if !valid {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, col.Name, colCollate, collates),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// ColumnCharsetCollateMustMatch 列的字符集与排序规则必须匹配
// RULE: MTB-L2-012
// TODO: 测试各种语法，了解字符集和排序规则的规则
func (v *TableAlterVldr) ColumnCharsetCollateMustMatch(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddColumns &&
			spec.Tp != ast.AlterTableChangeColumn &&
			spec.Tp != ast.AlterTableModifyColumn {
			continue
		}
		for _, col := range spec.NewColumns {
			colCharset := col.Tp.Charset
			colCollate := col.Tp.Collate

			// 排序规则允许隐式声明
			if len(colCharset) == 0 && len(colCollate) == 0 {
				continue
			}

			if len(colCharset) != 0 && len(colCollate) != 0 {
				if !charset.ValidCharsetAndCollation(colCharset, colCollate) {
					c := &models.Clause{
						Description: fmt.Sprintf(r.Message, col.Name, colCharset, colCollate),
						Level:       r.Level,
					}
					s.Violations.Append(c)
				}
			} else {
				// TODO:
				// 从数据库获取字符集和排序规则
			}
		}
	}
}

// PositionColumnDoesNotExist 位置标记列必须存在
// RULE: MTB-L2-013
func (v *TableAlterVldr) PositionColumnDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddColumns &&
			spec.Tp != ast.AlterTableChangeColumn &&
			spec.Tp != ast.AlterTableModifyColumn {
			continue
		}
		// Position 永远有值
		if spec.Position.Tp != ast.ColumnPositionAfter {
			continue
		}
		if ti.GetColumn(spec.Position.RelativeColumn.Name.O) != nil {
			continue
		}
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// ColumnNotNullWithDefaultRequired 非空列必须有默认值，应该做为警告级别
// RULE: MTB-L2-013
func (v *TableAlterVldr) ColumnNotNullWithDefaultRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	isNotNull := false
	hasDefaultValue := false

	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddColumns &&
			spec.Tp != ast.AlterTableChangeColumn &&
			spec.Tp != ast.AlterTableModifyColumn {
			continue
		}
		for _, col := range spec.NewColumns {
			if col.Options == nil {
				continue
			}
			for _, opt := range col.Options {
				if opt.Tp == ast.ColumnOptionNotNull {
					isNotNull = true
				}
				if opt.Tp == ast.ColumnOptionDefaultValue {
					hasDefaultValue = true
				}
			}
			if isNotNull && !hasDefaultValue {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, col.Name),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// IndexNameExplicit 索引必须命名
// RULE: MTB-L2-014
func (v *TableAlterVldr) IndexNameExplicit(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		keyName := ""
		switch spec.Tp {
		case ast.AlterTableAddConstraint:
			// ALTER TABLE ADD {INDEX|KEY|CONSTRAINT}
			c := spec.Constraint
			if c.Tp != ast.ConstraintKey && c.Tp != ast.ConstraintIndex {
				continue
			}
			keyName = strings.TrimSpace(c.Name)
		case ast.AlterTableRenameIndex:
			// ALTER TABLE RENAME {INDEX|KEY}
			keyName = strings.TrimSpace(spec.ToKey.O)
		default:
			continue
		}
		if len(keyName) == 0 {
			c := &models.Clause{
				Description: r.Message,
				Level:       r.Level,
			}
			s.Violations.Append(c)
			// 因为提示信息不会有名称，所以只要有违反规则，就不再检测其他的
			break
		}
	}
}

// IndexNameQualified 索引名标识符必须满足规则
// RULE: MTB-L2-015
func (v *TableAlterVldr) IndexNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		keyName := ""
		switch spec.Tp {
		case ast.AlterTableAddConstraint:
			// ALTER TABLE ADD {INDEX|KEY|CONSTRAINT}
			c := spec.Constraint
			if c.Tp != ast.ConstraintKey && c.Tp != ast.ConstraintIndex {
				continue
			}
			keyName = strings.TrimSpace(c.Name)
		case ast.AlterTableRenameIndex:
			// ALTER TABLE RENAME {INDEX|KEY}
			keyName = strings.TrimSpace(spec.ToKey.O)
		default:
			continue
		}
		if len(keyName) == 0 {
			continue // 如果系统允许匿名索引，则不进行正则匹配
		}
		if err := Match(r, keyName, keyName, r.Values); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// IndexNameLowerCaseRequired 索引名必须小写
// RULE: MTB-L2-016
func (v *TableAlterVldr) IndexNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		keyName := ""
		switch spec.Tp {
		case ast.AlterTableAddConstraint:
			// ALTER TABLE ADD {INDEX|KEY|CONSTRAINT}
			c := spec.Constraint
			if c.Tp != ast.ConstraintKey && c.Tp != ast.ConstraintIndex {
				continue
			}
			keyName = strings.TrimSpace(c.Name)
		case ast.AlterTableRenameIndex:
			// ALTER TABLE RENAME {INDEX|KEY}
			keyName = strings.TrimSpace(spec.ToKey.O)
		default:
			continue
		}
		if len(keyName) == 0 {
			continue // 如果系统允许匿名索引，则不进行正则匹配
		}

		if err := Match(r, keyName, keyName); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// IndexNameMaxLength 索引名最大长度
// RULE: MTB-L2-017
func (v *TableAlterVldr) IndexNameMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}

	for _, spec := range v.at.Specs {
		keyName := ""
		switch spec.Tp {
		case ast.AlterTableAddConstraint:
			// ALTER TABLE ADD {INDEX|KEY|CONSTRAINT}
			c := spec.Constraint
			if c.Tp != ast.ConstraintKey && c.Tp != ast.ConstraintIndex {
				continue
			}
			keyName = strings.TrimSpace(c.Name)
		case ast.AlterTableRenameIndex:
			// ALTER TABLE RENAME {INDEX|KEY}
			keyName = strings.TrimSpace(spec.ToKey.O)
		default:
			continue
		}
		if len(keyName) == 0 {
			continue // 如果系统允许匿名索引，则不进行正则匹配
		}

		if len(keyName) > threshold {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, keyName, threshold),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// IndexNamePrefixRequired 索引名前缀规则
// RULE: MTB-L2-018
func (v *TableAlterVldr) IndexNamePrefixRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		keyName := ""
		switch spec.Tp {
		case ast.AlterTableAddConstraint:
			// ALTER TABLE ADD {INDEX|KEY|CONSTRAINT}
			c := spec.Constraint
			if c.Tp != ast.ConstraintKey && c.Tp != ast.ConstraintIndex {
				continue
			}
			keyName = strings.TrimSpace(c.Name)
		case ast.AlterTableRenameIndex:
			// ALTER TABLE RENAME {INDEX|KEY}
			keyName = strings.TrimSpace(spec.ToKey.O)
		default:
			continue
		}
		if len(keyName) == 0 {
			continue // 如果系统允许匿名索引，则不进行正则匹配
		}
		if err := Match(r, keyName, keyName, r.Values); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// UniqueNameExplicit 唯一索引必须命名
// RULE: MTB-L2-019
func (v *TableAlterVldr) UniqueNameExplicit(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintUniq &&
			c.Tp != ast.ConstraintUniqIndex &&
			c.Tp != ast.ConstraintUniqKey {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			c := &models.Clause{
				Description: r.Message,
				Level:       r.Level,
			}
			s.Violations.Append(c)
			// 因为提示信息不会有名称，所以只要有违反规则，就不再检测其他的
			break
		}
	}
}

// UniqueNameQualified 唯一索引索名标识符必须符合规则
// RULE: MTB-L2-020
func (v *TableAlterVldr) UniqueNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintUniq && c.Tp != ast.ConstraintUniqIndex && c.Tp != ast.ConstraintUniqKey {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			continue // 如果系统允许匿名索引，则不进行正则匹配
		}
		if err := Match(r, keyName, keyName, r.Values); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// UniqueNameLowerCaseRequired 唯一索引名必须小写
// RULE: MTB-L2-021
func (v *TableAlterVldr) UniqueNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintUniq && c.Tp != ast.ConstraintUniqIndex && c.Tp != ast.ConstraintUniqKey {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			continue // 如果系统允许匿名索引，则不进行正则匹配
		}
		if err := Match(r, keyName, keyName); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// UniqueNameMaxLength 唯一索引名不能超过最大长度
// RULE: MTB-L2-022
func (v *TableAlterVldr) UniqueNameMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}

	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintUniq && c.Tp != ast.ConstraintUniqIndex && c.Tp != ast.ConstraintUniqKey {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) > threshold {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, keyName, threshold),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// UniqueNamePrefixRequired 唯一索引名前缀必须符合规则
// RULE: MTB-L2-023
func (v *TableAlterVldr) UniqueNamePrefixRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintUniq && c.Tp != ast.ConstraintUniqIndex && c.Tp != ast.ConstraintUniqKey {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			continue // 如果系统允许匿名索引，则不进行正则匹配
		}
		if err := Match(r, keyName, keyName, r.Values); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// ForeignKeyNotAllowed 禁止外键
// RULE: MTB-L2-024
func (v *TableAlterVldr) ForeignKeyNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint

		if c.Tp == ast.ConstraintForeignKey {
			c := &models.Clause{
				Description: r.Message,
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// ForeignKeyNameExplicit 外键是否显式命名
// RULE: MTB-L2-025
func (v *TableAlterVldr) ForeignKeyNameExplicit(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintForeignKey {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			c := &models.Clause{
				Description: r.Message,
				Level:       r.Level,
			}
			s.Violations.Append(c)
			break
		}
	}
}

// ForeignKeyNameQualified 外键名标识符规则
// RULE: MTB-L2-026
func (v *TableAlterVldr) ForeignKeyNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint

		if c.Tp != ast.ConstraintForeignKey {
			continue
		}
		keyName := c.Name
		if len(keyName) == 0 {
			continue
		}
		if err := Match(r, keyName, keyName, r.Values); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// ForeignKeyNameLowerCaseRequired 外键名必须小写
// RULE: MTB-L2-027
func (v *TableAlterVldr) ForeignKeyNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintForeignKey {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			continue
		}
		if err := Match(r, keyName, keyName); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// ForeignKeyNameMaxLength 外键名最大长度
// RULE: MTB-L2-028
func (v *TableAlterVldr) ForeignKeyNameMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}

	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintForeignKey {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) > threshold {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, keyName, threshold),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// ForeignKeyNamePrefixRequired 外键名前缀规则
// RULE: MTB-L2-029
func (v *TableAlterVldr) ForeignKeyNamePrefixRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintForeignKey {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			continue
		}
		if err := Match(r, keyName, keyName, r.Values); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// NewTableNameQualified 更名新表标识符规则
// RULE: MTB-L2-030
func (v *TableAlterVldr) NewTableNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableRenameTable || spec.NewTable == nil {
			continue
		}
		newName := spec.NewTable.Name.O
		if len(strings.TrimSpace(newName)) == 0 {
			newName = "<empty>"
		}
		if err := Match(r, newName, newName, r.Values); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// NewTableNameLowerCaseRequired 更名新表必须小写
// RULE: MTB-L2-031
func (v *TableAlterVldr) NewTableNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableRenameTable {
			continue
		}
		newName := strings.TrimSpace(spec.NewTable.Name.O)
		if len(newName) == 0 {
			continue
		}
		if err := Match(r, newName, newName); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// NewTableNameMaxLength 更名新表最大长度
// RULE: MTB-L2-032
func (v *TableAlterVldr) NewTableNameMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableRenameTable {
			continue
		}
		newName := strings.TrimSpace(spec.NewTable.Name.O)
		if len(newName) > threshold {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, newName, threshold),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// FullTextIndexNotAllowed 禁用全文索引
// RULE: MTB-L2-033
func (v *TableAlterVldr) FullTextIndexNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintFulltext {
			continue
		}
		if err := Match(r, c.Name); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// FullTextIndexExplicit 索引必须命名
// RULE: MTB-L2-034
func (v *TableAlterVldr) FullTextIndexExplicit(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintForeignKey {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			c := &models.Clause{
				Description: r.Message,
				Level:       r.Level,
			}
			s.Violations.Append(c)
			break
		}
	}
}

// FullTextIndexNameQualified 索引名标识符规则
// RULE: MTB-L2-035
func (v *TableAlterVldr) FullTextIndexNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintFulltext {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			continue
		}
		if err := Match(r, keyName, keyName, r.Values); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// FullTextIndexNameLowerCaseRequired 索引名必须小写
// RULE: MTB-L2-036
func (v *TableAlterVldr) FullTextIndexNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintFulltext {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			continue
		}
		if err := Match(r, keyName, keyName); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// FullTextIndexNameMaxLength 索引名不能超过最大长度
// RULE: MTB-L2-037
func (v *TableAlterVldr) FullTextIndexNameMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintFulltext {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) > threshold {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, keyName, threshold),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// FullTextIndexNamePrefixRequired 索引名前缀必须匹配规则
// RULE: MTB-L2-038
func (v *TableAlterVldr) FullTextIndexNamePrefixRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintFulltext {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			continue
		}
		if err := Match(r, keyName, keyName, r.Values); err != nil {
			c := &models.Clause{
				Description: err.Error(),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// MaxAllowedIndexColumnCount 单一索引最大列数
// RULE: MTB-L2-039
func (v *TableAlterVldr) MaxAllowedIndexColumnCount(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			keyName = "anonymous"
		}
		if len(c.Keys) > threshold {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, keyName, threshold),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// TargetDatabaseDoesNotExist 单一索引最大列数
// RULE: MTB-L2-039
func (v *TableAlterVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
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

// TargetTableDoesNotExist 单一索引最大列数
// RULE: MTB-L2-039
func (v *TableAlterVldr) TargetTableDoesNotExist(s *models.Statement, r *models.Rule) {
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

// ColumnNameDuplicate 添加列 - 数据库检查
// RULE: MTB
func (v *TableAlterVldr) ColumnNameDuplicate(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		// 添加列和更改列时，不能和现有列重复
		if spec.Tp != ast.AlterTableAddColumns &&
			spec.Tp != ast.AlterTableChangeColumn {
			continue
		}
		for _, col := range spec.NewColumns {
			ci := ti.GetColumn(col.Name.Name.O)
			if ci != nil {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, v.at.Table.Name.O, col.Name.Name.O),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// MaxAllowedColumnCount 列的数量是否超过阈值
// RULE: MTB-L3-006
func (v *TableAlterVldr) MaxAllowedColumnCount(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// 字段数量
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	// AlterTableDropColumn+AlterTableAddColumns影响总列数
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	diff := 0
	for _, spec := range v.at.Specs {
		if spec.Tp == ast.AlterTableDropColumn {
			diff--
		}
		if spec.Tp == ast.AlterTableAddColumns {
			diff += len(spec.NewColumns)
		}
	}

	total := len(ti.Columns()) + diff
	if total > threshold {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, v.at.Table.Name.O, total, threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// ColumnNameDoesNotExist 修改删除列 - 数据库检查
func (v *TableAlterVldr) ColumnNameDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableModifyColumn &&
			spec.Tp != ast.AlterTableChangeColumn &&
			spec.Tp != ast.AlterTableAlterColumn &&
			spec.Tp != ast.AlterTableDropColumn {
			continue
		}
		// 只有CHANGE/DROP语法涉及到的原始列使用OldColumnName保存
		if spec.Tp == ast.AlterTableChangeColumn ||
			spec.Tp == ast.AlterTableDropColumn {
			ci := ti.GetColumn(spec.OldColumnName.Name.O)
			if ci == nil {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, spec.OldColumnName.Name.O),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
			continue
		}
		for _, col := range spec.NewColumns {
			ci := ti.GetColumn(col.Name.Name.O)
			if ci == nil {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, col.Name.Name.O),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}
		}
	}
}

// IndexNameDuplicate 添加索引时名称不可重复
func (v *TableAlterVldr) IndexNameDuplicate(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	// TODO: 如果一次添加多个索引，这个索引之间也不可以存在问题
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		constraint := spec.Constraint
		keyName := strings.TrimSpace(constraint.Name)
		if v.IndexInfo(v.at.Table.Schema.O, v.at.Table.Name.O, keyName) == nil {
			continue
		}
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, keyName, v.at.Table.Name.O),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// IndexCountLimit 索引数量限制
func (v *TableAlterVldr) IndexCountLimit(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	// 字段数量
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	diff := 0
	for _, spec := range v.at.Specs {
		if spec.Tp == ast.AlterTableDropIndex ||
			spec.Tp == ast.AlterTableDropPrimaryKey {
			diff--
		}
		if spec.Tp == ast.AlterTableAddConstraint {
			diff++
		}
	}

	total := len(ti.Indexes) + diff
	if total > threshold {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, v.at.Table.Name.O, total, threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// IndexDoesNotExist 修改删除列 - 数据库检查
// TODO:
func (v *TableAlterVldr) IndexDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		keyName := ""
		switch spec.Tp {
		case ast.AlterTableDropIndex:
			keyName = spec.Name
		case ast.AlterTableRenameIndex:
			keyName = spec.FromKey.O
		default:
			continue
		}
		if v.IndexInfo(v.at.Table.Schema.O, v.at.Table.Name.O, keyName) == nil {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, keyName),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// IndexColumnDoesNotExist 索引的目标字段必须存在
func (v *TableAlterVldr) IndexColumnDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		for _, col := range c.Keys {
			if ti.GetColumn(col.Column.Name.O) != nil {
				continue
			}
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, col.Column.Name.O),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// TargetTableDuplicate 修改表名，如果新名称已经存在
func (v *TableAlterVldr) TargetTableDuplicate(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableRenameTable {
			continue
		}

		if v.TableInfo(spec.NewTable.Schema.O, spec.NewTable.Name.O) == nil {
			database := strings.TrimSpace(spec.NewTable.Schema.O)
			if len(database) == 0 {
				database = v.Ctx.Ticket.Database
			}

			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, fmt.Sprintf("`%s`.`%s`", database, spec.NewTable.Name.O)),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// PrimaryKeyDoesNotExist TODO
func (v *TableAlterVldr) PrimaryKeyDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableDropPrimaryKey {
			continue
		}
	}
}

// ForeignKeyDoesNotExist TODO
func (v *TableAlterVldr) ForeignKeyDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableDropForeignKey {
			continue
		}
	}
}

// IndexOnBlobColumnNotAllowed 禁止在BLOB/TEXT上索引
func (v *TableAlterVldr) IndexOnBlobColumnNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	ti := v.TableInfo(v.at.Table.Schema.O, v.at.Table.Name.O)
	if ti == nil {
		return
	}
	for _, spec := range v.at.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		for _, col := range c.Keys {
			ci := ti.GetColumn(col.Column.Name.O)
			if ci == nil {
				// 目标字段不存在，这个有另外函数处理
				// 这里只处理索引字段数据类型问题
				continue
			}
			// 既不是TEXT，也不是BLOB
			if !ci.SQLType.IsText() &&
				!ci.SQLType.IsBlob() {
				continue
			}
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, ci.Name),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
}

// MaxAllowedTimestampCount 仅允许一个事件戳类型的列
// RULE: MTB-L3-007
func (v *TableAlterVldr) MaxAllowedTimestampCount(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
}

// IndexOverlayNotAllowed 不允许索引覆盖
// RULE: MTB-L3-007
func (v *TableAlterVldr) IndexOverlayNotAllowed(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
}

// TableRenameVldr 改名表语句相关的审核规则
type TableRenameVldr struct {
	vldr

	rt *ast.RenameTableStmt
}

// Call 利用反射方法动态调用审核函数
func (v *TableRenameVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *TableRenameVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *TableRenameVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if rt, ok := node.(*ast.RenameTableStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.rt = rt
			v.Walk(v.rt)
		}
		for _, r := range v.Rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// TableRenameTablesIdentical 目标表跟源表是同一个表
// RULE: RTB-L2-001
func (v *TableRenameVldr) TableRenameTablesIdentical(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	tableName := v.rt.OldTable.Schema.O + v.rt.OldTable.Name.O
	newName := v.rt.NewTable.Schema.O + v.rt.NewTable.Name.O
	if strings.TrimSpace(tableName) == strings.TrimSpace(newName) {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, tableName, newName),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TableRenameTargetTableNameQualified 目标表名标识符规则
// RULE: RTB-L2-002
func (v *TableRenameVldr) TableRenameTargetTableNameQualified(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	newName := v.rt.NewTable.Name.O
	if err := Match(r, newName, newName, r.Values); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}

}

// TableRenameTargetTableNameLowerCaseRequired 目标表名大小写规则
// RULE: RTB-L2-003
func (v *TableRenameVldr) TableRenameTargetTableNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	newName := v.rt.NewTable.Name.O
	if err := Match(r, newName, newName); err != nil {
		c := &models.Clause{
			Description: err.Error(),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TableRenameTargetTableNameMaxLength 目标表名长度规则
// RULE: RTB-L2-004
func (v *TableRenameVldr) TableRenameTargetTableNameMaxLength(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	tableName := strings.TrimSpace(v.rt.NewTable.Name.O)
	count := len(tableName)
	if count > threshold {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, tableName, threshold),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}

}

// TableRenameSourceTableDoesNotExist 源库是否存在
// RULE: RTB-L3-001
func (v *TableRenameVldr) TableRenameSourceTableDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
}

// TableRenameSourceDatabaseDoesNotExist 源表是否存在
// RULE: RTB-L3-002
func (v *TableRenameVldr) TableRenameSourceDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
}

// TableRenameTargetTableDoesNotExist 目标库是否存在
// RULE: RTB-L3-003
func (v *TableRenameVldr) TableRenameTargetTableDoesNotExist(s *models.Statement, r *models.Rule) {
	log.Debugf("[D] RULE: %s, %s", r.Name, r.Func)
}

// TableDropVldr 删除表语句相关的审核规则
type TableDropVldr struct {
	vldr

	dt *ast.DropTableStmt
}

// Call 利用反射方法动态调用审核函数
func (v *TableDropVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *TableDropVldr) Enabled() bool {
	return true
}

// Validate 规则组的审核入口
func (v *TableDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.Ctx.Stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if dt, ok := node.(*ast.DropTableStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.dt = dt
			v.Walk(v.dt)
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
// RULE: DTB-L3-001
func (v *TableDropVldr) TargetDatabaseDoesNotExist(s *models.Statement, r *models.Rule) {
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
// RULE: DTB-L3-002
func (v *TableDropVldr) TargetTableDoesNotExist(s *models.Statement, r *models.Rule) {
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
