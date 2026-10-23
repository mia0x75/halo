# ADR 0005: 元数据库时间列改用 DATETIME，读写由数据库托管

## 状态

已接受

## 背景

元数据库的 `create_at`/`update_at` 原来是 `INT UNSIGNED` 存 Unix 秒，
由应用在 service 层 `time.Now().Unix()` 写入。问题：

- MySQL `TIMESTAMP` 有 2038 上限（`DATETIME` 无此问题）；
- 应用写时间戳导致同一语义在 6 个 service 文件里重复，且与
  `models.BeforeUpdate()` 各写各的；
- 为后续 PostgreSQL 元数据库（`TIMESTAMPTZ`）做准备，需要统一的
  UTC 存取约定。

## 决策

1. **类型**：MySQL 用 `DATETIME`（`create_at NOT NULL DEFAULT CURRENT_TIMESTAMP`，
   `update_at NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP`）；
   PG 对应 `TIMESTAMPTZ`（后续 Effort）。
2. **写入托管**：INSERT 不再带 `create_at`/`update_at`（靠列 DEFAULT）；
   MySQL UPDATE 不再带 `update_at`（靠 `ON UPDATE CURRENT_TIMESTAMP`）；
   PG UPDATE 显式 `update_at = NOW()`（后续 Effort）。
3. **读出转换**：storage 层 `sql.NullTime`/`time.Time` ↔ domain 的
   `uint`/`*uint`（Unix 秒），GraphQL 对外仍是 `UInt`，API 零破坏。
4. **UTC**：DSN 统一 `loc=UTC&time_zone='%2B00%3A00'`（引号不可省略，
   否则 go-sql-driver 拼出的 `SET time_zone=+00:00` 报 1064），
   会话时区固定 UTC。

## 后果

- service 层不再设置 `CreateAt`/`UpdateAt`（删 6 文件相关代码）；
  内存实体的写入后时间戳以 DB 为准，不再回填。
- 存量时间列由 UNIX 秒转 DATETIME 的迁移已折叠进 `0001_init.sql` 基线
  （项目未发布，基线直接按 DATETIME 建列）：转换用
  `FROM_UNIXTIME(NULLIF(col,0))`（Unix 0 表示"未更新"→NULL），
  `create_at` 遇 0 回退 `UTC_TIMESTAMP()` 以满足 NOT NULL；
  0 值在严格 sql_mode 下直接 `FROM_UNIXTIME(0)` 会报 1292。
- sqlc 重新生成，`Create*/Update*Params` 不再含时间字段。
