# 规则清单（RULES.md）

> 本文件由规则种子、注册表与处理器实现三方交叉核对生成，是当前 SQL 审核规则的权威索引。
>
> 数据来源：`scripts/halodb_seed_mysql.sql`（权威）、`scripts/halodb_seed_postgres.sql`、`internal/validate/registry.go`、`internal/validate/validate_*.go`。
>
> 规则总数 **248**，已实现 **240**，未实现 **8**（原因见「未实现规则与原因」）。

## 一、规则模型与分类维度

### 1.1 命名

规则名格式：`<前缀>-L<级别>-<序号>`，例如 `MTB-L2-006`。

- **前缀**：对象 / 语句类型（`CDB` 建库、`MTB` 改表、`SEL` 查询 ……）。
- **L 级别**：规则分类，**不是严重程度**（见 1.2）。
- **序号**：同前缀同级别内的编号。

### 1.2 L 级别（规则分类）

| 级别 | 含义 | 典型规则 |
|------|------|----------|
| `L0` | 工单级 / 跨语句 | `MSC-L0-001` 是否拆分工单 |
| `L1` | 危险或需成组禁止的语句 | `MSC-L1-001` 禁止 LOCK TABLE |
| `L2` | 对象定义、命名、结构约束 | `CTB-L2-006` 表名必须小写 |
| `L3` | 目标存在性、数量与边界 | `INS-L3-002` 目标表必须已存在 |

> ⚠️ L 级别与 `mm_rules.level`（严重程度）无一一对应关系：例如 `CDB-L2-007` 的严重程度是「错误」，`MTB-L3-005` 是「警告」。请勿混用。

### 1.3 `mm_rules` 关键字段

| 字段 | 含义 |
|------|------|
| `name` | 规则名（主键） |
| `group` | 业务分组号 |
| `level` | 严重程度：`0` 提示 / `1` 错误 / `2` 警告（见 `domain.RuleLevel`） |
| `vldr_group` | 校验器分组（对应 `internal/validate` 的 `XxxVldr`） |
| `operator` / `parameters` | 阈值或正则（`lte` 上限、`regexp` 正则、`in` 白名单、`not-in` 黑名单、`none`/`nil` 无参数） |
| `bitwise` | 启用位掩码；**偶数表示默认禁用**，奇数表示启用 |
| `func` | 处理器方法名（`(vldr_group, func)` 决定调用哪个方法） |
| `message` | 命中时的提示模板 |
| `element` | 前端渲染元素类型 |

### 1.4 状态图例

| 标记 | 含义 |
|------|------|
| ✅ | 已实现（处理器包含真正的校验逻辑） |
| ⚪ | 未实现（空函数体，仅打印调试日志或直接返回） |
| 默认「启用 / 禁用」 | 依据 `bitwise` 奇偶 |

### 1.5 统计总览

- 规则总数：248（MySQL / PG 种子一致）
- 已实现：240；未实现：8
- 默认启用：186；默认禁用：62
- 严重程度分布：错误 99、警告 149、提示 0

## 二、按语句类型的规则清单

### CREATE DATABASE（`CDB` · vldr_group 100）

共 7 条，已实现 7 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `CDB-L2-001` | 警告 | 启用 | `AvailableCharsets` ✅ | 新建数据库时允许的字符集 |
| `CDB-L2-002` | 警告 | 启用 | `AvailableCollates` ✅ | 新建数据库时允许的排序规则 |
| `CDB-L2-003` | 警告 | 启用 | `CharsetCollateMustMatch` ✅ | 新建数据库时字符集与排序规则必须匹配 |
| `CDB-L2-004` | 警告 | 启用 | `DatabaseNameQualified` ✅ | 库名规则 |
| `CDB-L2-005` | 警告 | 启用 | `DatabaseNameLowerCaseRequired` ✅ | 库名必须小写 |
| `CDB-L2-006` | 警告 | 启用 | `DatabaseNameMaxLength` ✅ | 库名最大长度 |
| `CDB-L2-007` | 错误 | 启用 | `TargetDatabaseDoesNotExist` ✅ | 新建数据库时目标库必须不存在 |

### ALTER DATABASE（`MDB` · vldr_group 101）

共 4 条，已实现 4 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `MDB-L2-001` | 警告 | 启用 | `AvailableCharsets` ✅ | 修改数据库时允许的字符集 |
| `MDB-L2-002` | 警告 | 启用 | `AvailableCollates` ✅ | 修改数据库时允许的排序规则 |
| `MDB-L2-003` | 警告 | 启用 | `CharsetCollateMustMatch` ✅ | 修改数据库时字符集与排序规则必须匹配 |
| `MDB-L2-004` | 错误 | 启用 | `TargetDatabaseDoesNotExist` ✅ | 修改数据库时目标库必须已存在 |

### DROP DATABASE（`DDB` · vldr_group 102）

共 1 条，已实现 1 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `DDB-L2-001` | 错误 | 启用 | `TargetDatabaseDoesNotExist` ✅ | 删除数据库时目标库必须已存在 |

### CREATE TABLE（`CTB` · vldr_group 110）

共 53 条，已实现 53 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `CTB-L2-001` | 警告 | 启用 | `AvailableCharsets` ✅ | 允许的字符集 |
| `CTB-L2-002` | 警告 | 启用 | `AvailableCollates` ✅ | 允许的排序规则 |
| `CTB-L2-003` | 警告 | 启用 | `TableCharsetCollateMustMatch` ✅ | 字符集与排序规则必须匹配 |
| `CTB-L2-004` | 警告 | 启用 | `AvailableEngines` ✅ | 允许的存储引擎 |
| `CTB-L2-005` | 警告 | 启用 | `TableNameQualified` ✅ | 表名规则 |
| `CTB-L2-006` | 警告 | 启用 | `TableNameLowerCaseRequired` ✅ | 表名必须小写 |
| `CTB-L2-007` | 警告 | 启用 | `TableNameMaxLength` ✅ | 表名最大长度 |
| `CTB-L2-008` | 警告 | 启用 | `TableCommentRequired` ✅ | 表必须有注释 |
| `CTB-L2-009` | 警告 | 启用 | `CreateTableFromSelectNotAllowed` ✅ | 禁止使用CREATE TABLE ... SELECT ...建表 |
| `CTB-L2-010` | 警告 | 启用 | `ColumnNameQualified` ✅ | 列名规则 |
| `CTB-L2-011` | 警告 | 启用 | `ColumnNameLowerCaseRequired` ✅ | 列名必须小写 |
| `CTB-L2-012` | 警告 | 启用 | `ColumnNameMaxLength` ✅ | 列名最大长度 |
| `CTB-L2-013` | 警告 | 启用 | `ColumnNameDuplicate` ✅ | 列名是否重复 |
| `CTB-L2-014` | 警告 | 启用 | `MaxAllowedColumnCount` ✅ | 表允许的最大列数 |
| `CTB-L2-015` | 警告 | 启用 | `ColumnTypesDoesNotExpect` ✅ | 列禁用的数据类型 |
| `CTB-L2-016` | 警告 | 启用 | `ColumnCommentRequired` ✅ | 列必须有注释 |
| `CTB-L2-017` | 警告 | 启用 | `ColumnAvailableCharsets` ✅ | 列允许的字符集 |
| `CTB-L2-018` | 警告 | 启用 | `ColumnAvailableCollates` ✅ | 列允许的排序规则 |
| `CTB-L2-019` | 警告 | 启用 | `ColumnCharsetCollateMustMatch` ✅ | 列字符集与排序规则必须匹配 |
| `CTB-L2-020` | 警告 | 启用 | `ColumnNotNullWithDefaultRequired` ✅ | 非空列必须有默认值 |
| `CTB-L2-021` | 警告 | 启用 | `ColumnAutoIncAvailableTypes` ✅ | 自增列允许的数据类型 |
| `CTB-L2-022` | 警告 | 启用 | `ColumnAutoIncUnsignedRequired` ✅ | 自增列必须是无符号 |
| `CTB-L2-023` | 警告 | 启用 | `ColumnAutoIncMustPrimaryKey` ✅ | 自增列必须是主键 |
| `CTB-L2-024` | 警告 | 启用 | `MaxAllowedTimestampCount` ✅ | 仅允许一个时间戳类型的列 |
| `CTB-L2-025` | 警告 | 启用 | `MaxAllowedIndexColumnCount` ✅ | 单一索引最大列数 |
| `CTB-L2-026` | 警告 | 启用 | `PrimaryKeyRequired` ✅ | 必须有主键 |
| `CTB-L2-027` | 警告 | 禁用 | `PrimaryKeyNameExplicit` ✅ | 主键是否显式命名 |
| `CTB-L2-028` | 警告 | 禁用 | `PrimaryKeyNameQualified` ✅ | 主键名规则 |
| `CTB-L2-029` | 警告 | 禁用 | `PrimryKeyLowerCaseRequired` ✅ | 主键名必须小写 |
| `CTB-L2-030` | 警告 | 禁用 | `PrimryKeyMaxLength` ✅ | 主键名最大长度 |
| `CTB-L2-031` | 警告 | 禁用 | `PrimryKeyPrefixRequired` ✅ | 主键名前缀规则 |
| `CTB-L2-032` | 警告 | 启用 | `IndexNameExplicit` ✅ | 索引必须命名 |
| `CTB-L2-033` | 警告 | 启用 | `IndexNameQualified` ✅ | 索引名规则 |
| `CTB-L2-034` | 警告 | 启用 | `IndexNameLowerCaseRequired` ✅ | 索引名必须小写 |
| `CTB-L2-035` | 警告 | 启用 | `IndexNameMaxLength` ✅ | 索引名最大长度 |
| `CTB-L2-036` | 警告 | 启用 | `IndexNamePrefixRequired` ✅ | 索引名前缀规则 |
| `CTB-L2-037` | 警告 | 启用 | `UniqueNameExplicit` ✅ | 唯一索引必须命名 |
| `CTB-L2-038` | 警告 | 启用 | `UniqueNameQualified` ✅ | 唯一索引索名规则 |
| `CTB-L2-039` | 警告 | 启用 | `UniqueNameLowerCaseRequired` ✅ | 唯一索引名必须小写 |
| `CTB-L2-040` | 警告 | 启用 | `UniqueNameMaxLength` ✅ | 唯一索引名最大长度 |
| `CTB-L2-041` | 警告 | 启用 | `UniqueNamePrefixRequired` ✅ | 唯一索引名前缀规则 |
| `CTB-L2-042` | 警告 | 启用 | `ForeignKeyNotAllowed` ✅ | 禁止外键 |
| `CTB-L2-043` | 警告 | 启用 | `ForeignKeyNameExplicit` ✅ | 外键是否显式命名 |
| `CTB-L2-044` | 警告 | 启用 | `ForeignKeyNameQualified` ✅ | 外键名规则 |
| `CTB-L2-045` | 警告 | 启用 | `ForeignKeyNameLowerCaseRequired` ✅ | 外键名必须小写 |
| `CTB-L2-046` | 警告 | 启用 | `ForeignKeyNameMaxLength` ✅ | 外键名最大长度 |
| `CTB-L2-047` | 警告 | 启用 | `ForeignKeyNamePrefixRequired` ✅ | 外键名前缀规则 |
| `CTB-L2-048` | 警告 | 启用 | `MaxAllowedIndexCount` ✅ | 表中最多可建多少个索引 |
| `CTB-L2-049` | 警告 | 启用 | `CreateTableUseLikeNotAllowed` ✅ | 禁止使用CREATE TABLE ... LIKE ...建表 |
| `CTB-L2-050` | 警告 | 启用 | `AutoIncColumnDuplicate` ✅ | 仅允许定义一个自增列 |
| `CTB-L2-051` | 警告 | 启用 | `PrimaryKeyDuplicate` ✅ | 仅允许定义一个主键 |
| `CTB-L3-001` | 错误 | 启用 | `TargetDatabaseDoesNotExist` ✅ | 目标库必须已存在 |
| `CTB-L3-002` | 错误 | 启用 | `TargetTableDoesNotExist` ✅ | 目标表必须不存在 |

### ALTER TABLE（`MTB` · vldr_group 120）

共 58 条，已实现 58 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `MTB-L2-001` | 警告 | 启用 | `AvailableCharsets` ✅ | 改表允许的字符集 |
| `MTB-L2-002` | 警告 | 启用 | `AvailableCollates` ✅ | 改表允许的校验规则 |
| `MTB-L2-003` | 警告 | 启用 | `TableCharsetCollateMustMatch` ✅ | 表的字符集与排序规则必须匹配 |
| `MTB-L2-004` | 警告 | 启用 | `AvailableEngines` ✅ | 改表允许的存储引擎 |
| `MTB-L2-005` | 警告 | 启用 | `ColumnNameQualified` ✅ | 列名必须符合命名规范 |
| `MTB-L2-006` | 警告 | 启用 | `ColumnNameLowerCaseRequired` ✅ | 列名必须小写 |
| `MTB-L2-007` | 警告 | 启用 | `ColumnNameMaxLength` ✅ | 列名最大长度 |
| `MTB-L2-008` | 警告 | 启用 | `ColumnUnwantedTypes` ✅ | 列禁用的数据类型 |
| `MTB-L2-009` | 警告 | 启用 | `ColumnCommentRequired` ✅ | 列必须有注释 |
| `MTB-L2-010` | 警告 | 启用 | `ColumnAvailableCharsets` ✅ | 列允许的字符集 |
| `MTB-L2-011` | 警告 | 启用 | `ColumnAvailableCollates` ✅ | 列允许的排序规则 |
| `MTB-L2-012` | 警告 | 启用 | `ColumnCharsetCollateMustMatch` ✅ | 列的字符集与排序规则必须匹配 |
| `MTB-L2-013` | 警告 | 启用 | `ColumnNotNullWithDefaultRequired` ✅ | 非空列必须有默认值 |
| `MTB-L2-014` | 警告 | 启用 | `IndexNameExplicit` ✅ | 索引必须命名 |
| `MTB-L2-015` | 警告 | 启用 | `IndexNameQualified` ✅ | 索引名标识符必须满足规则 |
| `MTB-L2-016` | 警告 | 启用 | `IndexNameLowerCaseRequired` ✅ | 索引名必须小写 |
| `MTB-L2-017` | 警告 | 启用 | `IndexNameMaxLength` ✅ | 索引名最大长度 |
| `MTB-L2-018` | 警告 | 启用 | `IndexNamePrefixRequired` ✅ | 索引名前缀规则 |
| `MTB-L2-019` | 警告 | 启用 | `UniqueNameExplicit` ✅ | 唯一索引必须命名 |
| `MTB-L2-020` | 警告 | 启用 | `UniqueNameQualified` ✅ | 唯一索引索名标识符必须符合规则 |
| `MTB-L2-021` | 警告 | 启用 | `UniqueNameLowerCaseRequired` ✅ | 唯一索引名必须小写 |
| `MTB-L2-022` | 警告 | 启用 | `UniqueNameMaxLength` ✅ | 唯一索引名不能超过最大长度 |
| `MTB-L2-023` | 警告 | 启用 | `UniqueNamePrefixRequired` ✅ | 唯一索引名前缀必须符合规则 |
| `MTB-L2-024` | 警告 | 启用 | `ForeignKeyNotAllowed` ✅ | 禁止外键 |
| `MTB-L2-025` | 警告 | 启用 | `ForeignKeyNameExplicit` ✅ | 外键是否显式命名 |
| `MTB-L2-026` | 警告 | 启用 | `ForeignKeyNameQualified` ✅ | 外键名标识符规则 |
| `MTB-L2-027` | 警告 | 启用 | `ForeignKeyNameLowerCaseRequired` ✅ | 外键名必须小写 |
| `MTB-L2-028` | 警告 | 启用 | `ForeignKeyNameMaxLength` ✅ | 外键名最大长度 |
| `MTB-L2-029` | 警告 | 启用 | `ForeignKeyNamePrefixRequired` ✅ | 外键名前缀规则 |
| `MTB-L2-030` | 警告 | 启用 | `NewTableNameQualified` ✅ | 更名新表规则 |
| `MTB-L2-031` | 警告 | 启用 | `NewTableNameLowerCaseRequired` ✅ | 更名新表必须小写 |
| `MTB-L2-032` | 警告 | 启用 | `NewTableNameMaxLength` ✅ | 更名新表最大长度 |
| `MTB-L2-033` | 警告 | 启用 | `FullTextIndexNotAllowed` ✅ | 禁用全文索引 |
| `MTB-L2-034` | 警告 | 启用 | `FullTextIndexNameExplicit` ✅ | 索引必须命名 |
| `MTB-L2-035` | 警告 | 启用 | `FullTextIndexNameQualified` ✅ | 索引名标识符规则 |
| `MTB-L2-036` | 警告 | 启用 | `FullTextIndexNameLowerCaseRequired` ✅ | 索引名必须小写 |
| `MTB-L2-037` | 警告 | 启用 | `FullTextIndexNameMaxLength` ✅ | 索引名不能超过最大长度 |
| `MTB-L2-038` | 警告 | 启用 | `FullTextIndexNamePrefixRequired` ✅ | 索引名前缀必须匹配规则 |
| `MTB-L2-039` | 警告 | 启用 | `MaxAllowedIndexColumnCount` ✅ | 单一索引最大列数 |
| `MTB-L3-001` | 错误 | 启用 | `TargetDatabaseDoesNotExist` ✅ | 目标库必须已存在 |
| `MTB-L3-002` | 错误 | 启用 | `TargetTableDoesNotExist` ✅ | 目标表必须已存在 |
| `MTB-L3-004` | 错误 | 启用 | `PositionColumnDoesNotExist` ✅ | 位置标记列必须已存在 |
| `MTB-L3-005` | 警告 | 启用 | `ColumnNameDuplicate` ✅ | 列名是否重复 |
| `MTB-L3-006` | 警告 | 禁用 | `MaxAllowedColumnCount` ✅ | 表允许的最大列数 |
| `MTB-L3-007` | 警告 | 禁用 | `MaxAllowedTimestampCount` ✅ | 仅允许一个时间戳类型的列 |
| `MTB-L3-008` | 错误 | 启用 | `ColumnNameDoesNotExist` ✅ | 删除列时目标列必须已存在 |
| `MTB-L3-013` | 错误 | 启用 | `IndexNameDuplicate` ✅ | 添加索引时索引必须不存在 |
| `MTB-L3-014` | 错误 | 启用 | `IndexOverlayNotAllowed` ✅ | 覆盖索引检查 |
| `MTB-L3-015` | 错误 | 启用 | `IndexColumnDoesNotExist` ✅ | 同名外键检查 |
| `MTB-L3-016` | 错误 | 启用 | `IndexOnBlobColumnNotAllowed` ✅ | 添加外键时外键必须不存在 |
| `MTB-L3-017` | 错误 | 启用 | `IndexDoesNotExist` ✅ | 启用禁用KEY时KEY必须已存在 |
| `MTB-L3-018` | 错误 | 启用 | `PrimaryKeyDoesNotExist` ✅ | 删主键时主键必须存在 |
| `MTB-L3-020` | 错误 | 启用 | `ForeignKeyDoesNotExist` ✅ | 删外键时外键必须存在 |
| `MTB-L3-021` | 错误 | 启用 | `TargetTableDuplicate` ✅ | 改名时目标表已存在 |
| `MTB-L3-022` | 错误 | 启用 | `FullTextIndexNameDuplicate` ✅ | 全文索引必须不存在 |
| `MTB-L3-023` | 错误 | 禁用 | `FullTextIndexNameDoesNotExist` ✅ | 删全文索引时索引必须存在 |
| `MTB-L3-024` | 错误 | 启用 | `CommentRequired` ✅ | 修改表注解信息时不可为空 |
| `MTB-L3-025` | 错误 | 启用 | `IndexCountLimit` ✅ | 索引数量超过允许的阈值 |

### RENAME TABLE（`RTB` · vldr_group 130）

共 8 条，已实现 8 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `RTB-L2-001` | 错误 | 启用 | `TablesIdentical` ✅ | 目标表跟源表是同一个表 |
| `RTB-L2-002` | 警告 | 启用 | `TargetTableNameQualified` ✅ | 目标表名规则 |
| `RTB-L2-003` | 警告 | 启用 | `TargetTableNameLowerCaseRequired` ✅ | 目标表名必须小写 |
| `RTB-L2-004` | 警告 | 启用 | `TargetTableNameMaxLength` ✅ | 目标表名最大长度 |
| `RTB-L3-001` | 错误 | 禁用 | `SourceTableDoesNotExist` ✅ | 源库必须已存在 |
| `RTB-L3-002` | 错误 | 禁用 | `SourceDatabaseDoesNotExist` ✅ | 源表必须已存在 |
| `RTB-L3-003` | 错误 | 禁用 | `TargetTableDoesNotExist` ✅ | 目标库必须已存在 |
| `RTB-L3-004` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 目标表必须不存在 |

### DROP TABLE（`DTB` · vldr_group 140）

共 2 条，已实现 2 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `DTB-L3-001` | 错误 | 启用 | `TableDropTargetDatabaseDoesNotExist` ✅ | 目标库必须已存在 |
| `DTB-L3-002` | 错误 | 启用 | `TableDropTargetTableDoesNotExist` ✅ | 目标表必须已存在 |

### CREATE INDEX（`CIX` · vldr_group 150）

共 13 条，已实现 13 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `CIX-L2-001` | 警告 | 启用 | `MaxAllowedIndexColumnCount` ✅ | 组合索引允许的最大列数 |
| `CIX-L2-002` | 警告 | 启用 | `IndexNameQualified` ✅ | 索引名规则 |
| `CIX-L2-003` | 警告 | 启用 | `IndexNameLowerCaseRequired` ✅ | 索引名必须小写 |
| `CIX-L2-004` | 警告 | 启用 | `IndexNameMaxLength` ✅ | 索引名最大长度 |
| `CIX-L2-005` | 警告 | 启用 | `IndexNamePrefixRequired` ✅ | 索引名前缀规则 |
| `CIX-L2-006` | 警告 | 启用 | `IndexColumnDuplicate` ✅ | 组合索引中是否有重复列 |
| `CIX-L3-001` | 错误 | 启用 | `TargetDatabaseDoesNotExist` ✅ | 目标库必须已存在 |
| `CIX-L3-002` | 错误 | 启用 | `TargetTableDoesNotExist` ✅ | 目标表必须已存在 |
| `CIX-L3-003` | 错误 | 启用 | `TargetColumnDoesNotExist` ✅ | 索引列必须已存在 |
| `CIX-L3-004` | 错误 | 禁用 | `IndexOverlayNotAllowed` ✅ | 索引内容是否重复 |
| `CIX-L3-005` | 错误 | 启用 | `IndexNameDuplicate` ✅ | 索引名是否重复 |
| `CIX-L3-006` | 错误 | 禁用 | `MaxAllowedIndexCount` ✅ | 最多能建多少个索引 |
| `CIX-L3-007` | 警告 | 启用 | `IndexOnBlobColumnNotAllowed` ✅ | 禁止在BLOB/TEXT列上建索引 |

### DROP INDEX（`RIX` · vldr_group 151）

共 3 条，已实现 3 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `RIX-L3-001` | 错误 | 启用 | `TargetDatabaseDoesNotExist` ✅ | 目标库必须已存在 |
| `RIX-L3-002` | 错误 | 启用 | `TargetTableDoesNotExist` ✅ | 目标表必须已存在 |
| `RIX-L3-003` | 错误 | 启用 | `TargetIndexDoesNotExist` ✅ | 目标索引必须已存在 |

### INSERT（`INS` · vldr_group 160）

共 9 条，已实现 8 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `INS-L2-001` | 警告 | 启用 | `ExplicitColumnRequired` ✅ | INSERT时强制显式列申明 |
| `INS-L2-002` | 警告 | 启用 | `UsingSelectNotAllowed` ✅ | 禁止INSERT...SELECT |
| `INS-L2-003` | 警告 | 禁用 | `MergeRequired` ⚪ | INSERT时是否需要合并（未实现） |
| `INS-L2-005` | 警告 | 启用 | `ColumnsValuesMustMatch` ✅ | INSERT时列类型、值是否匹配 |
| `INS-L3-001` | 错误 | 启用 | `TargetDatabaseDoesNotExist` ✅ | INSERT时目标库必须已存在 |
| `INS-L3-002` | 错误 | 启用 | `TargetTableDoesNotExist` ✅ | INSERT时目标表必须已存在 |
| `INS-L3-003` | 错误 | 禁用 | `TargetColumnDoesNotExist` ✅ | INSERT时目标列必须已存在 |
| `INS-L3-004` | 错误 | 禁用 | `ValueForNotNullColumnRequired` ✅ | INSERT时非空列是否有值 |
| `INS-L3-006` | 错误 | 启用 | `RowsLimit` ✅ | 单语句允许操作的最大行数 |

### REPLACE（`RPL` · vldr_group 161）

共 9 条，已实现 8 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `RPL-L2-001` | 警告 | 启用 | `ExplicitColumnRequired` ✅ | REPLACE时强制显式列申明 |
| `RPL-L2-002` | 警告 | 启用 | `UsingSelectNotAllowed` ✅ | 禁止REPLACE...SELECT |
| `RPL-L2-003` | 警告 | 禁用 | `MergeRequired` ⚪ | REPLACE时是否需要合并（未实现） |
| `RPL-L2-004` | 警告 | 启用 | `ReplaceRowsLimit` ✅ | REPLACE时单语句允许操作的最大行数 |
| `RPL-L2-005` | 警告 | 启用 | `ColumnsValuesMustMatch` ✅ | REPLACE时列类型、值是否匹配 |
| `RPL-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | REPLACE时目标库必须已存在 |
| `RPL-L3-002` | 错误 | 禁用 | `TargetTableDoesNotExist` ✅ | REPLACE时目标表必须已存在 |
| `RPL-L3-003` | 错误 | 禁用 | `TargetColumnDoesNotExist` ✅ | REPLACE时目标列必须已存在 |
| `RPL-L3-004` | 错误 | 禁用 | `ValueForNotNullColumnRequired` ✅ | REPLACE时非空列是否有值 |

### UPDATE（`UPD` · vldr_group 162）

共 6 条，已实现 5 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `UPD-L2-001` | 错误 | 启用 | `WithoutWhereNotAllowed` ✅ | 禁止没有WHERE的更新 |
| `UPD-L3-001` | 错误 | 启用 | `TargetDatabaseDoesNotExist` ✅ | 目标库必须已存在 |
| `UPD-L3-002` | 错误 | 启用 | `TargetTableDoesNotExist` ✅ | 目标表必须已存在 |
| `UPD-L3-003` | 错误 | 禁用 | `TargetColumnDoesNotExist` ✅ | 目标列必须已存在 |
| `UPD-L3-004` | 错误 | 禁用 | `FilterColumnExists` ⚪ | 条件限定列必须已存在（未实现） |
| `UPD-L3-005` | 错误 | 禁用 | `MaxAllowedUpdateRows` ✅ | 允许单次更新的最大行数 |

### DELETE（`DEL` · vldr_group 163）

共 5 条，已实现 5 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `DEL-L2-001` | 警告 | 启用 | `WithoutWhereNotAllowed` ✅ | 禁止没有WHERE的删除 |
| `DEL-L3-001` | 错误 | 启用 | `MaxAllowedDeleteRows` ✅ | 单次删除的最大行数 |
| `DEL-L3-002` | 错误 | 启用 | `TargetDatabaseDoesNotExist` ✅ | 目标库必须已存在 |
| `DEL-L3-003` | 错误 | 启用 | `TargetTableDoesNotDoesNotExist` ✅ | 目标表必须已存在 |
| `DEL-L3-004` | 错误 | 禁用 | `TargetColumnDoesNotDoesNotExist` ✅ | 条件过滤列必须已存在 |

### SELECT（`SEL` · vldr_group 170）

共 8 条，已实现 8 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `SEL-L2-001` | 警告 | 启用 | `WithoutWhereNotAllowed` ✅ | 禁止没有WHERE的查询 |
| `SEL-L2-002` | 警告 | 启用 | `WithoutLimitNotAllowed` ✅ | 禁止没有LIMIT的查询 |
| `SEL-L2-003` | 警告 | 启用 | `UseWildcardNotAllowed` ✅ | 禁止SELECT STAR |
| `SEL-L2-004` | 警告 | 启用 | `UseExplicitLockNotAllowed` ✅ | 禁止SELECT FOR UPDATE |
| `SEL-L3-001` | 错误 | 启用 | `TargetDatabaseDoesNotExist` ✅ | 目标数据库必须已存在 |
| `SEL-L3-002` | 错误 | 启用 | `TargetTableDoesNotExist` ✅ | 目标表必须已存在 |
| `SEL-L3-003` | 错误 | 禁用 | `TargetColumnDoesNotDoesNotExist` ✅ | 目标列必须已存在 |
| `SEL-L3-004` | 错误 | 禁用 | `ReturnBlobOrTextNotAllowed` ✅ | 是否允许返回BLOB/TEXT列 |

### CREATE VIEW（`CVW` · vldr_group 180）

共 6 条，已实现 6 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `CVW-L2-001` | 警告 | 启用 | `ViewNameQualified` ✅ | 新建视图时视图名规则 |
| `CVW-L2-002` | 警告 | 启用 | `ViewNameLowerCaseRequired` ✅ | 新建视图时视图名必须小写 |
| `CVW-L2-003` | 警告 | 启用 | `ViewNameMaxLength` ✅ | 新建视图时视图名最大长度 |
| `CVW-L2-004` | 警告 | 启用 | `ViewNamePrefixRequired` ✅ | 新建视图时视图名前缀规则 |
| `CVW-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 新建视图时目标库必须已存在 |
| `CVW-L3-002` | 错误 | 禁用 | `TargetViewDoesNotExist` ✅ | 新建视图时目标视图必须不存在 |

### ALTER VIEW（`MVW` · vldr_group 181）

共 2 条，已实现 2 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `MVW-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 修改视图时目标库必须已存在 |
| `MVW-L3-002` | 错误 | 禁用 | `TargetViewDoesNotExist` ✅ | 修改视图时目标视图必须已存在 |

### DROP VIEW（`DVW` · vldr_group 182）

共 2 条，已实现 2 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `DVW-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 删除视图时目标库必须已存在 |
| `DVW-L3-002` | 错误 | 禁用 | `TargetViewDoesNotExist` ✅ | 删除视图时目标视图必须已存在 |

### CREATE FUNCTION（`CFU` · vldr_group 190）

共 6 条，已实现 6 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `CFU-L2-001` | 警告 | 启用 | `FuncNameQuilified` ✅ | 新建函数时函数名规则 |
| `CFU-L2-002` | 警告 | 启用 | `FuncNameLowerCaseRequired` ✅ | 新建函数时函数名必须小写 |
| `CFU-L2-003` | 警告 | 启用 | `FuncNameMaxLength` ✅ | 新建函数时函数名最大长度 |
| `CFU-L2-004` | 警告 | 启用 | `FuncNamePrefixRequired` ✅ | 新建函数时函数名前缀规则 |
| `CFU-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 新建函数时目标库必须已存在 |
| `CFU-L3-002` | 错误 | 禁用 | `TargetFuncDoesNotExist` ✅ | 新建函数时目标函数必须不存在 |

### ALTER FUNCTION（`MFU` · vldr_group 191）

共 2 条，已实现 2 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `MFU-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 修改函数时目标库必须已存在 |
| `MFU-L3-002` | 错误 | 禁用 | `TargetFuncDoesNotExist` ✅ | 修改函数时目标函数必须已存在 |

### DROP FUNCTION（`DFU` · vldr_group 192）

共 2 条，已实现 2 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `DFU-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 删除函数时目标库必须已存在 |
| `DFU-L3-002` | 错误 | 禁用 | `TargetFuncDoesNotExist` ✅ | 删除函数时目标函数必须已存在 |

### CREATE TRIGGER（`CTG` · vldr_group 200）

共 7 条，已实现 7 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `CTG-L2-001` | 警告 | 启用 | `TriggerNameQualified` ✅ | 新建触发器时触发器名规则 |
| `CTG-L2-002` | 警告 | 启用 | `TriggerNameLowerCaseRequired` ✅ | 新建触发器时触发器名必须小写 |
| `CTG-L2-003` | 警告 | 启用 | `TriggerNameMaxLength` ✅ | 新建触发器时触发器名最大长度 |
| `CTG-L2-004` | 警告 | 启用 | `TriggerPrefixRequired` ✅ | 新建触发器时触发器名前缀规则 |
| `CTG-L3-001` | 警告 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 新建触发器时目标库必须已存在 |
| `CTG-L3-002` | 警告 | 禁用 | `TargetTableDoesNotExist` ✅ | 新建触发器时目标表必须已存在 |
| `CTG-L3-003` | 警告 | 禁用 | `TargetTriggerDoesNotExist` ✅ | 新建触发器时目标触发器必须不存在 |

### ALTER TRIGGER（`MTG` · vldr_group 201）

共 3 条，已实现 0 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `MTG-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ⚪ | 修改触发器时目标库必须已存在（未实现） |
| `MTG-L3-002` | 错误 | 禁用 | `TargetTableDoesNotExist` ⚪ | 修改触发器时目标表必须已存在（未实现） |
| `MTG-L3-003` | 错误 | 禁用 | `TargetTriggerDoesNotExist` ⚪ | 修改触发器时目标触发器必须已存在（未实现） |

### DROP TRIGGER（`DTG` · vldr_group 202）

共 3 条，已实现 2 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `DTG-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 删除触发器时目标库必须已存在 |
| `DTG-L3-002` | 错误 | 禁用 | `TargetTableDoesNotExist` ⚪ | 删除触发器时目标表必须已存在（未实现） |
| `DTG-L3-003` | 错误 | 禁用 | `TargetTriggerDoesNotExist` ✅ | 删除触发器时目标触发器必须已存在 |

### CREATE EVENT（`CEV` · vldr_group 210）

共 6 条，已实现 6 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `CEV-L2-001` | 警告 | 启用 | `EventNameQualified` ✅ | 创建事件时事件名规则 |
| `CEV-L2-002` | 警告 | 启用 | `EventNameLowerCaseRequired` ✅ | 创建事件时事件名必须小写 |
| `CEV-L2-003` | 警告 | 启用 | `EventNameMaxLength` ✅ | 创建事件时事件名最大长度 |
| `CEV-L2-004` | 警告 | 启用 | `EventNamePrefixRequired` ✅ | 创建事件时事件名前缀规则 |
| `CEV-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 创建事件时目标库必须已存在 |
| `CEV-L3-002` | 错误 | 禁用 | `TargetEventDoesNotExist` ✅ | 创建事件时目标事件必须已存在 |

### ALTER EVENT（`MEV` · vldr_group 211）

共 2 条，已实现 2 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `MEV-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 修改事件时目标库必须已存在 |
| `MEV-L3-002` | 错误 | 禁用 | `TargetEventDoesNotExist` ✅ | 修改事件时目标事件必须已存在 |

### DROP EVENT（`DEV` · vldr_group 212）

共 2 条，已实现 2 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `DEV-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 删除事件时目标库必须已存在 |
| `DEV-L3-002` | 错误 | 禁用 | `TargetEventDoesNotExist` ✅ | 删除事件时目标事件必须已存在 |

### CREATE PROCEDURE（`CSP` · vldr_group 220）

共 6 条，已实现 6 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `CSP-L2-001` | 警告 | 启用 | `ProcNameQualified` ✅ | 新建存储过程时存储过程名规则 |
| `CSP-L2-002` | 警告 | 启用 | `ProcNameLowerCaseRequired` ✅ | 新建存储过程时存储过程名必须小写 |
| `CSP-L2-003` | 警告 | 启用 | `ProcNameMaxLength` ✅ | 新建存储过程时存储过程名最大长度 |
| `CSP-L2-004` | 警告 | 启用 | `ProcNamePrefixRequired` ✅ | 新建存储过程时存储过程名前缀规则 |
| `CSP-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 新建存储过程时目标库必须已存在 |
| `CSP-L3-002` | 错误 | 禁用 | `TargetProcDoesNotExist` ✅ | 新建存储过程时目标存储过程必须不存在 |

### ALTER PROCEDURE（`MSP` · vldr_group 221）

共 2 条，已实现 2 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `MSP-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 修改存储过程时目标库必须已存在 |
| `MSP-L3-002` | 错误 | 禁用 | `TargetProcDoesNotExist` ✅ | 修改存储过程时目标存储过程必须已存在 |

### DROP PROCEDURE（`DSP` · vldr_group 222）

共 2 条，已实现 2 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `DSP-L3-001` | 错误 | 禁用 | `TargetDatabaseDoesNotExist` ✅ | 删除存储过程时目标库必须已存在 |
| `DSP-L3-002` | 错误 | 禁用 | `TargetProcDoesNotExist` ✅ | 删除存储过程时目标存储过程必须已存在 |

### 跨语句 / 工单（`MSC` · vldr_group 230）

共 9 条，已实现 8 条。

| 规则 | 严重 | 默认 | 处理器 | 目的 |
|------|------|------|--------|------|
| `MSC-L0-001` | 错误 | 禁用 | `SplitTicket` ⚪ | 是否拆分工单（未实现） |
| `MSC-L1-001` | 错误 | 启用 | `LockTableProhibited` ✅ | 禁止LOCK TABLE |
| `MSC-L1-002` | 错误 | 启用 | `FlushTableProhibited` ✅ | 禁止FLUSH TABLE |
| `MSC-L1-003` | 错误 | 启用 | `TruncateTableProhibited` ✅ | 禁止TRUNCATE TABLE |
| `MSC-L1-004` | 错误 | 启用 | `MergeRequired` ✅ | 对同一个表/库的操作需要合并 |
| `MSC-L1-005` | 错误 | 启用 | `PurgeLogsProhibited` ✅ | 禁止PURGE LOG |
| `MSC-L1-006` | 错误 | 启用 | `UnlockTableProhibited` ✅ | 禁止UNLOCK TABLE |
| `MSC-L1-007` | 错误 | 启用 | `KillProhibited` ✅ | 禁止KILL |
| `MSC-L1-008` | 错误 | 启用 | `SplitRequired` ✅ | 禁止同时出现DDL、DML |

## 三、未实现规则与原因

共 **8** 条未实现，分三类：

### 3.1 MySQL 无该语句：ALTER TRIGGER（3 条）

MySQL **没有** `ALTER TRIGGER` 语法，`mia0x75/parser` 也没有对应 AST，因此 `MTG` 组收不到任何语句。处理器保留为文档化的空实现。

| 语句类型 | 未实现规则 | 条数 |
|----------|------------|------|
| ALTER TRIGGER `MTG` | `MTG-L3-001..003` | 3 |

**实现条件**：当目标方言（如 PG）引入 `ALTER TRIGGER` 且有对应 AST 时再补处理器。

**状态：待日后处理**（本轮无法实现，保留空实现，不删除种子行）。

### 3.2 AST 不携带目标表：DROP TRIGGER（1 条）

| 规则 | 说明 |
|------|------|
| `DTG-L3-002` 删除触发器时目标表必须已存在 | MySQL `DROP TRIGGER` 只写触发器名（`*ast.DropTriggerStmt` 仅含 `Name`），不带表引用，无法判定目标表。 |

**实现条件**：需在元数据中维护「触发器 → 所属表」映射，或目标方言的 `DROP TRIGGER` 携带表名。

**状态：待日后处理**（本轮无法实现，保留空实现，不删除种子行）。

### 3.3 语义已由其它规则覆盖（4 条）

保留这些规则的处理器是为了让种子里的规则能解析到处理器（否则 `ValidateRegistry` 会报「缺 handler」），实际校验交由等价规则完成：

| 规则 | 目的 | 由谁覆盖 |
|------|------|----------|
| `INS-L2-003` `MergeRequired` | INSERT 是否需要合并 | `MSC-L1-004` 对同一表/库的操作需要合并 |
| `RPL-L2-003` `MergeRequired` | REPLACE 是否需要合并 | `MSC-L1-004` |
| `UPD-L3-004` `FilterColumnExists` | 条件限定列必须已存在 | `UPD-L3-003` 已同时校验 SET 与 WHERE 列 |
| `MSC-L0-001` `SplitTicket` | 是否拆分工单 | `MSC-L1-008` 禁止同时出现 DDL、DML |

### 3.4 已补齐（本轮）

自 `mia0x75/parser` 升级到 `v0.2.0` 后，以下 43 条已实现：

- 存储程序命名 + 存在性：`CFU`/`MFU`/`DFU`（函数）、`CSP`/`MSP`/`DSP`（存储过程）、`CTG`/`DTG`（触发器）、`CEV`/`MEV`/`DEV`（事件）。
- `ALTER VIEW`：`MVW-L3-001/002`。
- 视图存在性：`CVW-L3-002`、`DVW-L3-002`（`internal/agent` 新增 `LoadObjects`，经 `information_schema` / `pg_proc` / `pg_trigger` 读取对象清单）。
- SELECT 过滤列：`SEL-L3-003` 由仅校验选择列表扩展为同时校验 `WHERE` / `JOIN ON` / `GROUP BY` / `ORDER BY` / `HAVING`（承接上游被移除的 `SEL-L3-005`）。

## 四、缺失规则分析

### 4.1 与上游规则集的差异

对照上游 `e80ef5e:scripts/halodb.sql`（251 条），当前种子删除 11 条、新增 8 条。

**删除的 11 条（含去向）：**

| 上游规则 | 现状 |
|----------|------|
| `CDB-L3-001` DatabaseCreateTargetDatabaseExists | → `CDB-L2-007`（改名 + 级别 L3→L2） |
| `DDB-L3-001` DatabaseDropTargetDatabaseExists | → `DDB-L2-001` |
| `MDB-L3-001` DatabaseAlterTargetDatabaseExists | → `MDB-L2-004` |
| `INS-L2-004` InsertRowsLimit（上限 10000） | → `INS-L3-006` RowsLimit（上限 1000） |
| `MTB-L3-003` 新增列时目标列必须不存在 | 合并进 `MTB-L3-005`（列名是否重复） |
| `MTB-L3-009` 修改列时目标列必须已存在 | 合并进 `MTB-L3-008` |
| `MTB-L3-010` 修改列时目标列必须不存在 | 合并进 `MTB-L3-005` |
| `MTB-L3-011` 仅允许一个时间戳列 | → `MTB-L3-007` MaxAllowedTimestampCount |
| `MTB-L3-012` 列名是否重复 | → `MTB-L3-005` ColumnNameDuplicate |
| `MTB-L3-019` 删索引时索引必须存在 | 合并进 `MTB-L3-017` IndexDoesNotExist |
| `SEL-L3-005` SelectFilterColumnExists | → 能力并入 `SEL-L3-003`（本轮扩展过滤列校验） |

**新增的 8 条：** `CDB-L2-007`、`MDB-L2-004`、`DDB-L2-001`、`INS-L2-003`、`INS-L3-006`、`RPL-L2-003`、`MSC-L0-001`、`MTB-L3-025`。

> 结论：删除项都有等价承接，不构成能力缺失。

### 4.2 真实缺口（建议补齐）

1. **PostgreSQL 语句完全未纳入规则引擎（高优先级）**
   需求前提是 MySQL + PostgreSQL 双目标库，但 `internal/validate` 仅处理 MySQL AST（`mia0x75/parser`）；PG 语句（经 `pgplex/pgparser` 解析）不经过任何规则。当前只有语法分析层（`internal/analyzer`）支持 PG。
   *建议*：明确 PG 审核范围，或补充 PG 版校验器。

2. **ALTER TRIGGER（`MTG`，3 条）与 `DTG-L3-002` 不可达**（见 3.1 / 3.2）

3. **管理/会话类语句无规则**
   以下语句有 AST 但无对应规则：`CREATE/ALTER/DROP USER`、`GRANT/REVOKE`、`SET`、`PREPARE/EXECUTE/DEALLOCATE`、`LOAD DATA`、`ANALYZE TABLE`、`BEGIN/COMMIT/ROLLBACK`。其中高风险项已由 `MSC-L1-*` 禁止（`LOCK/FLUSH/TRUNCATE/PURGE/UNLOCK/KILL`）。

### 4.3 对象级覆盖核对

| 对象 | 命名规则 | 存在性规则 | 结论 |
|------|----------|------------|------|
| 数据库 | `CDB-L2-004..006` | `CDB-L2-007`/`MDB-L2-004`/`DDB-L2-001` | 完整 |
| 表 | `CTB-L2-005..007`/`MTB-L2-030..032`/`RTB-L2-002..004` | `CTB-L3-002`/`MTB-L3-002`/`RTB-L3-*`/`DTB-L3-*` | 完整 |
| 列 | `CTB-L2-010..012`/`MTB-L2-005..007` | `INS/UPD/RPL-L3-003`、`MTB-L3-008` | 完整 |
| 索引 | `CIX-L2-002..005`/`CTB-L2-032..036` | `CIX-L3-005`、`MTB-L3-013`、`RIX-L3-003` | 完整 |
| 唯一索引 / 外键 / 全文索引 / 主键 | `CTB/MTB-L2-02x..03x` | `MTB-L3-018/020/022/023` | 完整 |
| 视图 | `CVW-L2-001..004` | `CVW-L3-002`/`MVW-L3-002`/`DVW-L3-002` | 完整 |
| 存储程序 | `CFU/CTG/CEV/CSP-L2-*` | `CFU/CTG/CEV/CSP-L3-*` 等 | 完整（MTG/DTG-L3-002 除外，见 3.1/3.2） |
| SELECT 过滤列 | — | `SEL-L3-003`（含 WHERE/JOIN/GROUP/ORDER/HAVING） | 完整 |

## 五、已知数据问题

> 以下为核对中发现的种子/消息缺陷。按协作约定先记录、报告，不擅自改动行为。

1. **`UPD-L3-004` 严重程度越界（已修复）**
   原 `level=3` 超出 `RuleLevel`（仅 0/1/2）。已在本轮修正为 `1`（错误），MySQL / PG 种子同步。

2. **`RTB-L3-001..004` 描述与 `func` 语义错位（上游遗留）**

   | 规则 | 处理器（真实语义） | 描述 / 消息（当前） |
   |------|--------------------|---------------------|
   | `RTB-L3-001` | `SourceTableDoesNotExist` 源表必须已存在 | 写作「源库」 |
   | `RTB-L3-002` | `SourceDatabaseDoesNotExist` 源库必须已存在 | 写作「源表」 |
   | `RTB-L3-003` | `TargetTableDoesNotExist` 目标表必须不存在 | 写作「目标库」 |
   | `RTB-L3-004` | `TargetDatabaseDoesNotExist` 目标库必须已存在 | 写作「目标表」 |

3. **`MTB-L3-015` / `MTB-L3-016` 描述、消息与实现不符，且消息占位符数与处理器参数不匹配**

   | 规则 | 实现 | 描述 / 消息（当前） | 问题 |
   |------|------|---------------------|------|
   | `MTB-L3-015` | `IndexColumnDoesNotExist`：校验新增约束引用的列是否存在 | 描述「同名外键检查」；消息「目标表"%s"上已存在外键"%s"。」 | 消息 2 个占位符，处理器只传 1 个（列名）→ 命中时输出 `%!s(MISSING)` |
   | `MTB-L3-016` | `IndexOnBlobColumnNotAllowed`：禁止在 TEXT/BLOB 列加索引 | 描述「添加外键时外键必须不存在」；消息同上 | 同上 |

4. **`CTB-L2-027` 消息多余参数**
   处理器 `PrimaryKeyNameExplicit` 传入表名，但消息「主键没有提供名称。」无占位符 → 命中时输出 `%!(EXTRA string=...)`。

5. **描述/消息重复或与实现不符**（示例）
   `MTB-L3-014` 描述「覆盖索引检查」而消息为「目标表"%s"上已存在索引"%s"。」；`MTB-L2-034..038`（全文索引）与 `CIX` 命名族描述雷同，容易误读。

6. **拼写错误（上游遗留，为保持兼容未改名）**
   `PrimryKeyLowerCaseRequired`、`PrimryKeyMaxLength`、`PrimryKeyPrefixRequired`、`FuncNameQuilified`、`TargetColumnDoesNotDoesNotExist`、`TargetTableDoesNotDoesNotExist`、`TableDropTargetTableDoesNotExist`。

## 六、维护约定

- 新增/修改规则时，**同时**改 MySQL 与 PG 两个种子文件，保持条数一致。
- `func` 必须在 `internal/validate/registry.go` 的对应 `vldr_group` 中注册，否则规则不会被调度。
- 无法实现的规则保留空函数体并写明原因（英文注释），不要删除种子行。
- 严重程度只能用 `0/1/2`。
- 修改后运行：`go test ./internal/validate/...`（需 `HALO_TEST_CONFIG=$(pwd)/bin/config.test.json`）与 `golangci-lint run ./internal/validate/...`。
