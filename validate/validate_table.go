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

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *TableCreateVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *TableCreateVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if ct, ok := node.(*ast.CreateTableStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.ct = ct
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// TableCreateAvailableCharsets 建表允许的字符集
// RULE: CTB-L2-001
func (v *TableCreateVldr) TableCreateAvailableCharsets(s *models.Statement, r *models.Rule) {
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

// TableCreateAvailableCollates 建表允许的排序规则
// RULE: CTB-L2-002
func (v *TableCreateVldr) TableCreateAvailableCollates(s *models.Statement, r *models.Rule) {
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

// TableCreateTableCharsetCollateMatch 建表是校验规则与字符集必须匹配
// RULE: CTB-L2-003
func (v *TableCreateVldr) TableCreateTableCharsetCollateMatch(s *models.Statement, r *models.Rule) {
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

// TableCreateAvailableEngines 建表允许的存储引擎
// RULE: CTB-L2-004
func (v *TableCreateVldr) TableCreateAvailableEngines(s *models.Statement, r *models.Rule) {
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

// TableCreateTableNameQualified 表名必须符合命名规范
// RULE: CTB-L2-005
func (v *TableCreateVldr) TableCreateTableNameQualified(s *models.Statement, r *models.Rule) {
	tableName := v.ct.Table.Name.O
	if err := Match(r, tableName, tableName, r.Values); err != nil {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, tableName, r.Values),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TableCreateTableNameLowerCaseRequired 表名是否允许大写
// RULE: CTB-L2-006
func (v *TableCreateVldr) TableCreateTableNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	tableName := strings.TrimSpace(v.ct.Table.Name.O)
	if err := Match(r, tableName, tableName); err != nil {
		c := &models.Clause{
			Description: fmt.Sprintf(r.Message, tableName),
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TableCreateTableNameMaxLength 表名最大长度
// RULE: CTB-L2-007
func (v *TableCreateVldr) TableCreateTableNameMaxLength(s *models.Statement, r *models.Rule) {
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

// TableCreateTableCommentRequired 表必须有注释
// RULE: CTB-L2-008
func (v *TableCreateVldr) TableCreateTableCommentRequired(s *models.Statement, r *models.Rule) {
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

// TableCreateUseSelectEnabled 是否允许查询语句建表
// RULE: CTB-L2-009
func (v *TableCreateVldr) TableCreateUseSelectEnabled(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnNameQualified 列名必须符合命名规范
// RULE: CTB-L2-010
func (v *TableCreateVldr) TableCreateColumnNameQualified(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnNameLowerCaseRequired 列名是否允许大写
// RULE: CTB-L2-011
func (v *TableCreateVldr) TableCreateColumnNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnNameMaxLength 列名最大长度
// RULE: CTB-L2-012
func (v *TableCreateVldr) TableCreateColumnNameMaxLength(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnNameDuplicate 列名是否重复
// RULE: CTB-L2-013
func (v *TableCreateVldr) TableCreateColumnNameDuplicate(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnCountLimit 表允许的最大列数
// RULE: CTB-L2-014
func (v *TableCreateVldr) TableCreateColumnCountLimit(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnUnwantedTypes 列不允许的数据类型
// RULE: CTB-L2-015
func (v *TableCreateVldr) TableCreateColumnUnwantedTypes(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnCommentRequired 列必须有注释
// RULE: CTB-L2-016
func (v *TableCreateVldr) TableCreateColumnCommentRequired(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnAvailableCharsets 列允许的字符集
// RULE: CTB-L2-017
func (v *TableCreateVldr) TableCreateColumnAvailableCharsets(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnAvailableCollates 列允许的排序规则
// RULE: CTB-L2-018
func (v *TableCreateVldr) TableCreateColumnAvailableCollates(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnCharsetCollateMatch 列的字符集与排序规则必须匹配
// RULE: CTB-L2-019
func (v *TableCreateVldr) TableCreateColumnCharsetCollateMatch(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnNotNullWithDefaultRequired 非空列是否有默认值
// RULE: CTB-L2-020
func (v *TableCreateVldr) TableCreateColumnNotNullWithDefaultRequired(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnAutoIncAvailableTypes 自增列允许的数据类型
// RULE: CTB-L2-021
func (v *TableCreateVldr) TableCreateColumnAutoIncAvailableTypes(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnAutoIncIsUnsigned 自增列必须是无符号
// RULE: CTB-L2-022
func (v *TableCreateVldr) TableCreateColumnAutoIncIsUnsigned(s *models.Statement, r *models.Rule) {
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

// TableCreateColumnAutoIncMustPrimaryKey 自增列必须是主键
// RULE: CTB-L2-023
func (v *TableCreateVldr) TableCreateColumnAutoIncMustPrimaryKey(s *models.Statement, r *models.Rule) {
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

// TableCreateTimestampColumnCountLimit 仅允许一个时间戳类型的列
// RULE: CTB-L2-024
func (v *TableCreateVldr) TableCreateTimestampColumnCountLimit(s *models.Statement, r *models.Rule) {
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

// TableCreateIndexMaxColumnLimit 单一索引最大列数
// RULE: CTB-L2-025
func (v *TableCreateVldr) TableCreateIndexMaxColumnLimit(s *models.Statement, r *models.Rule) {
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

// TableCreatePrimaryKeyRequired 必须有主键
// RULE: CTB-L2-026
func (v *TableCreateVldr) TableCreatePrimaryKeyRequired(s *models.Statement, r *models.Rule) {
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

// TableCreatePrimaryKeyNameExplicit 主键是否显式命名
// RULE: CTB-L2-027
func (v *TableCreateVldr) TableCreatePrimaryKeyNameExplicit(s *models.Statement, r *models.Rule) {
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

// TableCreatePrimaryKeyNameQualified 主键名标识符规则
// RULE: CTB-L2-028
func (v *TableCreateVldr) TableCreatePrimaryKeyNameQualified(s *models.Statement, r *models.Rule) {
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

// TableCreatePrimryKeyLowerCaseRequired 主键名大小写规则
// RULE: CTB-L2-029
func (v *TableCreateVldr) TableCreatePrimryKeyLowerCaseRequired(s *models.Statement, r *models.Rule) {
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

// TableCreatePrimryKeyMaxLength 主键名长度规则
// RULE: CTB-L2-030
func (v *TableCreateVldr) TableCreatePrimryKeyMaxLength(s *models.Statement, r *models.Rule) {
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

// TableCreatePrimryKeyPrefixRequired 主键名前缀规则
// RULE: CTB-L2-031
func (v *TableCreateVldr) TableCreatePrimryKeyPrefixRequired(s *models.Statement, r *models.Rule) {
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

// TableCreateIndexNameExplicit 索引必须命名
// RULE: CTB-L2-032
func (v *TableCreateVldr) TableCreateIndexNameExplicit(s *models.Statement, r *models.Rule) {
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

// TableCreateIndexNameQualified 索引名标识符规则
// RULE: CTB-L2-033
func (v *TableCreateVldr) TableCreateIndexNameQualified(s *models.Statement, r *models.Rule) {
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

// TableCreateIndexNameLowerCaseRequired 索引名大小写规则
// RULE: CTB-L2-034
func (v *TableCreateVldr) TableCreateIndexNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
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

// TableCreateIndexNameMaxLength 索引名长度规则
// RULE: CTB-L2-035
func (v *TableCreateVldr) TableCreateIndexNameMaxLength(s *models.Statement, r *models.Rule) {
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

// TableCreateIndexNamePrefixRequired 索引名前缀规则
// RULE: CTB-L2-036
func (v *TableCreateVldr) TableCreateIndexNamePrefixRequired(s *models.Statement, r *models.Rule) {
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

// TableCreateUniqueNameExplicit 唯一索引必须命名
// RULE: CTB-L2-037
func (v *TableCreateVldr) TableCreateUniqueNameExplicit(s *models.Statement, r *models.Rule) {
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

// TableCreateUniqueNameQualified 唯一索引索名标识符规则
// RULE: CTB-L2-038
func (v *TableCreateVldr) TableCreateUniqueNameQualified(s *models.Statement, r *models.Rule) {
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

// TableCreateUniqueNameLowerCaseRequired 唯一索引名大小写规则
// RULE: CTB-L2-039
func (v *TableCreateVldr) TableCreateUniqueNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
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

// TableCreateUniqueNameMaxLength 唯一索引名长度规则
// RULE: CTB-L2-040
func (v *TableCreateVldr) TableCreateUniqueNameMaxLength(s *models.Statement, r *models.Rule) {
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

// TableCreateUniqueNamePrefixRequired 唯一索引名前缀规则
// RULE: CTB-L2-041
func (v *TableCreateVldr) TableCreateUniqueNamePrefixRequired(s *models.Statement, r *models.Rule) {
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

// TableCreateForeignKeyEnabled 是否允许外键
// RULE: CTB-L2-042
func (v *TableCreateVldr) TableCreateForeignKeyEnabled(s *models.Statement, r *models.Rule) {
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

// TableCreateForeignKeyNameExplicit 是否配置外键名称
// RULE: CTB-L2-043
func (v *TableCreateVldr) TableCreateForeignKeyNameExplicit(s *models.Statement, r *models.Rule) {
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

// TableCreateForeignKeyNameQualified 外键名标识符规则
// RULE: CTB-L2-044
func (v *TableCreateVldr) TableCreateForeignKeyNameQualified(s *models.Statement, r *models.Rule) {
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

// TableCreateForeignKeyNameLowerCaseRequired 外键名大小写规则
// RULE: CTB-L2-045
func (v *TableCreateVldr) TableCreateForeignKeyNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
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

// TableCreateForeignKeyNameMaxLength 外键名长度规则
// RULE: CTB-L2-046
func (v *TableCreateVldr) TableCreateForeignKeyNameMaxLength(s *models.Statement, r *models.Rule) {
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

// TableCreateForeignKeyNamePrefixRequired 外键名前缀规则
// RULE: CTB-L2-047
func (v *TableCreateVldr) TableCreateForeignKeyNamePrefixRequired(s *models.Statement, r *models.Rule) {
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

// TableCreateIndexCountLimit 表中最多可建多少个索引
// RULE: CTB-L2-048
func (v *TableCreateVldr) TableCreateIndexCountLimit(s *models.Statement, r *models.Rule) {
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

// TableCreateUseLikeEnabled 禁止允许LIKE方式建表
// RULE: CTB-L2-049
func (v *TableCreateVldr) TableCreateUseLikeEnabled(s *models.Statement, r *models.Rule) {
	// 规则允许执行，表示禁止新建表
	if v.ct.ReferTable != nil {
		c := &models.Clause{
			Description: r.Message,
			Level:       r.Level,
		}
		s.Violations.Append(c)
	}
}

// TableCreateAutoIncColumnCountLimit 只允许一个自增列
// RULE: CTB-L2-050
func (v *TableCreateVldr) TableCreateAutoIncColumnCountLimit(s *models.Statement, r *models.Rule) {
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

// TableCreatePrimaryKeyCountLimit 只允许一个主键
// RULE: CTB-L2-051
func (v *TableCreateVldr) TableCreatePrimaryKeyCountLimit(s *models.Statement, r *models.Rule) {
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

// TableCreateTargetDatabaseExists 建表时检查目标库是否存在
// RULE: CTB-L3-001
func (v *TableCreateVldr) TableCreateTargetDatabaseExists(s *models.Statement, r *models.Rule) {
}

// TableCreateTargetTableExists 建表时检查表是否已经存在
// RULE: CTB-L3-002
func (v *TableCreateVldr) TableCreateTargetTableExists(s *models.Statement, r *models.Rule) {
}

// TableAlterVldr 修改表语句相关的审核规则
type TableAlterVldr struct {
	vldr

	rt *ast.AlterTableStmt
}

// Call 利用反射方法动态调用审核函数
func (v *TableAlterVldr) Call(method string, params ...interface{}) {
	Call(v, method, params...)
}

// Enabled 当前规则组是否生效
func (v *TableAlterVldr) Enabled() bool {
	return true
}

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *TableAlterVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *TableAlterVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if rt, ok := node.(*ast.AlterTableStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.rt = rt
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// TableAlterAvailableCharsets 改表允许的字符集
// RULE: MTB-L2-001
func (v *TableAlterVldr) TableAlterAvailableCharsets(s *models.Statement, r *models.Rule) {
	var charsets []string
	// 将字符串反解析为结构体
	json.Unmarshal([]byte(r.Values), &charsets)

	useCharset := "<empty>"
	valid := true
	for _, spec := range v.rt.Specs {
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

// TableAlterAvailableCollates 改表允许的校验规则
// RULE: MTB-L2-002
func (v *TableAlterVldr) TableAlterAvailableCollates(s *models.Statement, r *models.Rule) {
	var collates []string
	// 将字符串反解析为结构体
	json.Unmarshal([]byte(r.Values), &collates)

	useCollate := "<empty>"
	valid := true
	for _, spec := range v.rt.Specs {
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

// TableAlterCharsetCollateMatch 表的字符集与排序规则必须匹配
// RULE: MTB-L2-003
func (v *TableAlterVldr) TableAlterCharsetCollateMatch(s *models.Statement, r *models.Rule) {
	useCharset := ""
	useCollate := ""
	for _, spec := range v.rt.Specs {
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

// TableAlterAvailableEngines 改表允许的存储引擎
// RULE: MTB-L2-004
func (v *TableAlterVldr) TableAlterAvailableEngines(s *models.Statement, r *models.Rule) {
	var engines []string
	// 将字符串反解析为结构体
	json.Unmarshal([]byte(r.Values), &engines)

	useEngine := "<empty>"
	valid := true
	for _, spec := range v.rt.Specs {
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

// TableAlterColumnNameQualified 列名必须符合命名规范
// RULE: MTB-L2-005
func (v *TableAlterVldr) TableAlterColumnNameQualified(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
		if spec.Tp != ast.AlterTableAddColumns && spec.Tp != ast.AlterTableChangeColumn {
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

// TableAlterColumnNameLowerCaseRequired 列名必须小写
// RULE: MTB-L2-006
func (v *TableAlterVldr) TableAlterColumnNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
		if spec.Tp != ast.AlterTableAddColumns && spec.Tp != ast.AlterTableChangeColumn {
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

// TableAlterColumnNameMaxLength 列名最大长度
// RULE: MTB-L2-007
func (v *TableAlterVldr) TableAlterColumnNameMaxLength(s *models.Statement, r *models.Rule) {
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	for _, spec := range v.rt.Specs {
		if spec.Tp != ast.AlterTableAddColumns && spec.Tp != ast.AlterTableChangeColumn {
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

// TableAlterColumnUnwantedTypes 列允许的数据类型
// RULE: MTB-L2-008
func (v *TableAlterVldr) TableAlterColumnUnwantedTypes(s *models.Statement, r *models.Rule) {
	var availableTypes []string
	json.Unmarshal([]byte(r.Values), &availableTypes)

	if len(availableTypes) == 0 {
		return
	}
	for _, spec := range v.rt.Specs {
		if spec.NewColumns == nil {
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

// TableAlterColumnCommentRequired 列必须有注释
// RULE: MTB-L2-009
func (v *TableAlterVldr) TableAlterColumnCommentRequired(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
		if spec.NewColumns == nil {
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

// TableAlterColumnAvailableCharsets 列允许的字符集
// RULE: MTB-L2-010
func (v *TableAlterVldr) TableAlterColumnAvailableCharsets(s *models.Statement, r *models.Rule) {
	// 字符集
	var charsets []string
	json.Unmarshal([]byte(r.Values), &charsets)

	for _, spec := range v.rt.Specs {
		if spec.NewColumns == nil {
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

// TableAlterColumnAvailableCollates 列允许的排序规则
// RULE: MTB-L2-011
func (v *TableAlterVldr) TableAlterColumnAvailableCollates(s *models.Statement, r *models.Rule) {
	// 排序规则
	var collates []string
	json.Unmarshal([]byte(r.Values), &collates)
	for _, spec := range v.rt.Specs {
		if spec.NewColumns == nil {
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

// TableAlterColumnCharsetCollateMatch 列的字符集与排序规则必须匹配
// RULE: MTB-L2-012
func (v *TableAlterVldr) TableAlterColumnCharsetCollateMatch(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
		if spec.NewColumns == nil {
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

// TableAlterColumnNotNullWithDefaultRequired 非空列必须有默认值
// RULE: MTB-L2-013
func (v *TableAlterVldr) TableAlterColumnNotNullWithDefaultRequired(s *models.Statement, r *models.Rule) {
	isNotNull := false
	hasDefaultValue := false

	for _, spec := range v.rt.Specs {
		if spec == nil || (spec.Tp != ast.AlterTableAddColumns && spec.Tp != ast.AlterTableModifyColumn && spec.Tp != ast.AlterTableChangeColumn) {
			continue
		}
		for _, col := range spec.NewColumns {
			name := col.Name
			if col.Options != nil {
				for _, opt := range col.Options {
					if opt.Tp == ast.ColumnOptionNotNull {
						isNotNull = true
					}
					if opt.Tp == ast.ColumnOptionDefaultValue {
						hasDefaultValue = true
					}

				}
			}
			if isNotNull && !hasDefaultValue {
				c := &models.Clause{
					Description: fmt.Sprintf(r.Message, name),
					Level:       r.Level,
				}
				s.Violations.Append(c)
			}

		}
	}

}

// TableAlterIndexNameExplicit 索引必须命名
// RULE: MTB-L2-014
func (v *TableAlterVldr) TableAlterIndexNameExplicit(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintKey && c.Tp != ast.ConstraintIndex {
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

// TableAlterIndexNameQualified 索引名标识符必须满足规则
// RULE: MTB-L2-015
func (v *TableAlterVldr) TableAlterIndexNameQualified(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintKey && c.Tp != ast.ConstraintIndex {
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

// TableAlterIndexNameLowerCaseRequired 索引名必须小写
// RULE: MTB-L2-016
func (v *TableAlterVldr) TableAlterIndexNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintKey && c.Tp != ast.ConstraintIndex {
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

// TableAlterIndexNameMaxLength 索引名最大长度
// RULE: MTB-L2-017
func (v *TableAlterVldr) TableAlterIndexNameMaxLength(s *models.Statement, r *models.Rule) {
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}

	for _, spec := range v.rt.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintKey && c.Tp != ast.ConstraintIndex {
			continue
		}
		keyName := strings.TrimSpace(c.Name)
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

// TableAlterIndexNamePrefixRequired 索引名前缀规则
// RULE: MTB-L2-018
func (v *TableAlterVldr) TableAlterIndexNamePrefixRequired(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintKey && c.Tp != ast.ConstraintIndex {
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

// TableAlterUniqueNameExplicit 唯一索引必须命名
// RULE: MTB-L2-019
func (v *TableAlterVldr) TableAlterUniqueNameExplicit(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		if c.Tp != ast.ConstraintUniq && c.Tp != ast.ConstraintUniqIndex && c.Tp != ast.ConstraintUniqKey {
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

// TableAlterUniqueNameQualified 唯一索引索名标识符必须符合规则
// RULE: MTB-L2-020
func (v *TableAlterVldr) TableAlterUniqueNameQualified(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterUniqueNameLowerCaseRequired 唯一索引名必须小写
// RULE: MTB-L2-021
func (v *TableAlterVldr) TableAlterUniqueNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterUniqueNameMaxLength 唯一索引名不能超过最大长度
// RULE: MTB-L2-022
func (v *TableAlterVldr) TableAlterUniqueNameMaxLength(s *models.Statement, r *models.Rule) {
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}

	for _, spec := range v.rt.Specs {
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

// TableAlterUniqueNamePrefixRequired 唯一索引名前缀必须符合规则
// RULE: MTB-L2-023
func (v *TableAlterVldr) TableAlterUniqueNamePrefixRequired(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterForeignKeyEnabled 禁止外键
// RULE: MTB-L2-024
func (v *TableAlterVldr) TableAlterForeignKeyEnabled(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterForeignKeyNameExplicit 外键是否显式命名
// RULE: MTB-L2-025
func (v *TableAlterVldr) TableAlterForeignKeyNameExplicit(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterForeignKeyNameQualified 外键名标识符规则
// RULE: MTB-L2-026
func (v *TableAlterVldr) TableAlterForeignKeyNameQualified(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterForeignKeyNameLowerCaseRequired 外键名必须小写
// RULE: MTB-L2-027
func (v *TableAlterVldr) TableAlterForeignKeyNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterForeignKeyNameMaxLength 外键名最大长度
// RULE: MTB-L2-028
func (v *TableAlterVldr) TableAlterForeignKeyNameMaxLength(s *models.Statement, r *models.Rule) {
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}

	for _, spec := range v.rt.Specs {
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

// TableAlterForeignKeyNamePrefixRequired 外键名前缀规则
// RULE: MTB-L2-029
func (v *TableAlterVldr) TableAlterForeignKeyNamePrefixRequired(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterNewTableNameQualified 更名新表标识符规则
// RULE: MTB-L2-030
func (v *TableAlterVldr) TableAlterNewTableNameQualified(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterNewTableNameLowerCaseRequired 更名新表必须小写
// RULE: MTB-L2-031
func (v *TableAlterVldr) TableAlterNewTableNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterNewTableNameMaxLength 更名新表最大长度
// RULE: MTB-L2-032
func (v *TableAlterVldr) TableAlterNewTableNameMaxLength(s *models.Statement, r *models.Rule) {
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	for _, spec := range v.rt.Specs {
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

// TableAlterFullTextEnabled 禁用全文索引
// RULE: MTB-L2-033
func (v *TableAlterVldr) TableAlterFullTextEnabled(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterFullTextNameExplicit 索引必须命名
// RULE: MTB-L2-034
func (v *TableAlterVldr) TableAlterFullTextNameExplicit(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterFullTextNameQualified 索引名标识符规则
// RULE: MTB-L2-035
func (v *TableAlterVldr) TableAlterFullTextNameQualified(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterFullTextNameLowerCaseRequired 索引名必须小写
// RULE: MTB-L2-036
func (v *TableAlterVldr) TableAlterFullTextNameLowerCaseRequired(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterFullTextNameMaxLength 索引名不能超过最大长度
// RULE: MTB-L2-037
func (v *TableAlterVldr) TableAlterFullTextNameMaxLength(s *models.Statement, r *models.Rule) {
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	for _, spec := range v.rt.Specs {
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

// TableAlterFullTextNamePrefixRequired 索引名前缀必须匹配规则
// RULE: MTB-L2-038
func (v *TableAlterVldr) TableAlterFullTextNamePrefixRequired(s *models.Statement, r *models.Rule) {
	for _, spec := range v.rt.Specs {
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

// TableAlterIndexMaxColumnLimit 单一索引最大列数
// RULE: MTB-L2-039
func (v *TableAlterVldr) TableAlterIndexMaxColumnLimit(s *models.Statement, r *models.Rule) {
	threshold, err := strconv.Atoi(r.Values)
	if err != nil {
		return
	}
	for _, spec := range v.rt.Specs {
		if spec.Tp != ast.AlterTableAddConstraint {
			continue
		}
		c := spec.Constraint
		keyName := strings.TrimSpace(c.Name)
		if len(keyName) == 0 {
			continue
		}
		if len(c.Keys) > threshold {
			c := &models.Clause{
				Description: fmt.Sprintf(r.Message, c.Name, threshold),
				Level:       r.Level,
			}
			s.Violations.Append(c)
		}
	}
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

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *TableRenameVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *TableRenameVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if rt, ok := node.(*ast.RenameTableStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.rt = rt
		}
		for _, r := range v.rules {
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

// TableRenameSourceTableExists 源库是否存在
// RULE: RTB-L3-001
func (v *TableRenameVldr) TableRenameSourceTableExists(s *models.Statement, r *models.Rule) {
}

// TableRenameSourceDatabaseExists 源表是否存在
// RULE: RTB-L3-002
func (v *TableRenameVldr) TableRenameSourceDatabaseExists(s *models.Statement, r *models.Rule) {
}

// TableRenameTargetTableExists 目标库是否存在
// RULE: RTB-L3-003
func (v *TableRenameVldr) TableRenameTargetTableExists(s *models.Statement, r *models.Rule) {
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

// SetContext 在不同的规则组之间共享信息，这个可能暂时没用
func (v *TableDropVldr) SetContext(ctx Context) {
}

// Validate 规则组的审核入口
func (v *TableDropVldr) Validate(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, s := range v.stmts {
		// 该方法不能放到结构体vldr是因为，反射时找不到子类的方法
		node := s.StmtNode
		if dt, ok := node.(*ast.DropTableStmt); !ok {
			// 类型断言不成功
			continue
		} else {
			v.dt = dt
		}
		for _, r := range v.rules {
			if r.Bitwise&1 != 1 {
				continue
			}
			v.Call(r.Func, s, r)
		}
	}
}

// TableDropSourceDatabaseExists 目标库是否存在
// RULE: DTB-L3-001
func (v *TableDropVldr) TableDropSourceDatabaseExists(s *models.Statement, r *models.Rule) {
}

// TableDropSourceTableExists 目标表是否存在
// RULE: DTB-L3-002
func (v *TableDropVldr) TableDropSourceTableExists(s *models.Statement, r *models.Rule) {
}
