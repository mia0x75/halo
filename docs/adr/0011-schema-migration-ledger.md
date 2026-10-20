# ADR 0011: 元数据库迁移台账与迁移执行器

## 状态

已接受

## 背景

仓库里已有 `migrations/mysql/`、`migrations/postgres/` 两个目录，但：

- **没有任何代码应用这些文件**（全仓 grep 不到调用方），它们是孤儿文件；
- 没有版本台账，无法知道某个库的 schema 处于哪个版本；
- 两个引擎的迁移集不对齐：PostgreSQL 只有 `0001_init.sql`，MySQL 缺 `0001`，
  且 `0001_init.sql` 里 `mm_rules` 的列名是 `value`/`values`，而代码与权威
  schema 用的是 `parameter`/`parameters`——这份基线在 PostgreSQL 上根本跑不通；
- `migrations/postgres/0001_init.sql` 的内容与 `scripts/halodb_schema_postgres.sql`
  已经漂移。

同时，`scripts/halodb_schema_mysql.sql` 与 `scripts/halodb_schema_postgres.sql`
是 AGENTS.md 指定的表结构权威源，测试环境用它们初始化容器。

## 决策

**引入 `mm_migrations` 台账表 + 迁移执行器，迁移文件用 `go:embed` 打进单二进制，
由 `halocli migrate` 子命令驱动，不在服务启动时自动执行。**

要点：

1. **台账表** `mm_migrations`（两个元数据库各一份，列一致）：
   `version`（主键）、`name`（唯一）、`checksum`（文件 sha256）、`duration_ms`、
   `success`、`apply_at`。表由迁移器自身以 `CREATE TABLE IF NOT EXISTS` 建立，
   同时也写进权威 schema 脚本，保证新库开箱即用。
2. **迁移文件**：`migrations/<engine>/NNNN_name.sql`，按版本号排序执行；
   文件用 `go:embed` 打包进二进制（保持单二进制交付，不引外部工具）。
3. **漂移检测**：已应用的迁移文件内容变更（checksum 不一致）、或台账中的版本在
   迁移目录里缺失，直接报错拒绝执行。已应用的迁移文件不可修改。
4. **并发保护**：MySQL 用 `GET_LOCK`，PostgreSQL 用 `pg_advisory_lock`，
   同一元数据库上多个实例不会重复执行。
5. **失败语义**：
   - PostgreSQL 支持事务性 DDL，整个文件在一个事务内执行，失败回滚且不写台账，
     迁移保持「待执行」可安全重试；
   - MySQL 的 DDL 隐式提交，无法回滚，失败时写入 `success=0` 的台账记录，
     后续执行会拒绝继续，必须人工修复数据库并清理该记录。
6. **基线折叠**：项目尚未发布，`0002_add_cluster_type.sql`、
   `0003_timestamp_to_datetime.sql` 的变更已经包含在基线快照里，从未在任何库
   执行过。这两个文件加上 `-- migrate:folded-into-baseline` 标记，迁移器只登记
   台账、不执行语句——既保留迁移历史（每个文件一条台账记录），又不会重复执行。
7. **基线同步校验**：`scripts/check_migrations_sync.sh` 比较基线与权威 schema
   脚本（忽略 DROP 语句、`IF NOT EXISTS` 差异与注释），防止两者分叉。
8. **不自动迁移**：`db.New` 的注释已经写明「caller 自行迁移，避免 schema 漂移
   悄悄启动服务」；迁移只在 `halocli migrate up` 时执行，符合 AGENTS.md
   「迁移前暂停写入 + 备份」的要求。

## 后果

正面：

- 迁移从「没人执行的文件」变成有台账、可校验、可并发安全的正式流程；
- PostgreSQL 基线里 `value`/`values` 与代码不一致的问题被修正；
- 后续每个 schema 变更都有明确的编号、校验和与执行记录。

负面：

- 台账一旦写入，迁移文件就不能再改（含注释），改错只能追加新迁移；
- MySQL 迁移失败需要人工介入，不能自动重试；
- 权威 schema 脚本与迁移链是两份表示，需要 `check_migrations_sync.sh` 卡住漂移。

## 后续

- 新增表/字段一律：先写 `migrations/<engine>/NNNN_*.sql`（两引擎各一份），
  再同步 `scripts/halodb_schema_*.sql`，最后跑 `scripts/check_migrations_sync.sh`。
- 目标库（用户业务库）的 schema 演进不由 Halo 接管，不在目标库建台账。
