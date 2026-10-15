# ADR 0004: 移除 xorm，元数据与执行统一走 agent

## 状态

已接受

## 背景

平台早期用 `xorm.io/xorm` + `go-xorm/core` 直连目标库：`models.Cluster`
上的 `Connect`/`Metadata`/`Databases`/`Stat` 直接向目标库发出 `SHOW`/`DESCRIBE`
等语句，并借助 `core.Table`/`core.Column`/`core.Index` 做元数据建模。

随着 2A/2B/2C 落地，元数据库已全面改为 sqlc + sqlx；目标库访问也收归
`internal/agent`（`Executor` 接口 + 双方言 `dialect`）。此时 xorm 只剩两处残留价值：

- `models.Cluster` 的元数据方法（`Connect`/`Metadata`/`Databases`/`Stat`/`Repack`）
- `core.*` 元数据类型，被 `validate/` 与 Metadata resolver 消费

问题：
- xorm 违背"元数据库无 ORM"的既定方向（ADR 0001 精神、AGENTS.md 红线 F5 相邻）
- 目标库访问分散在 `models` 层，绕过 `internal/agent`，破坏分层
- `go-xorm/core` 经 `replace` 指向社区 fork，维护来源非官方
- 双解析/双元数据模型并存，心智负担与类型漂移风险

## 决策

**彻底移除 xorm 与 `go-xorm/core`**，用 `internal/agent/meta` 的中立类型 +
`agent.Executor` 取代。

具体：

1. 新增 `internal/agent/meta`：定义 dialect-neutral 的
   `Table`/`Column`/`Index`/`Database`/`ColumnType`，保留 `validate/` 真正用到的
   字段与 `ColumnType.IsText/IsBlob`、`Table.GetColumn` 语义。
2. `agent.Executor` 收敛为：`Execute`、`Query`、`LoadTables`（可变库名参数：
   不传时拉取实例全部用户库，承接 validate 的 `"*"` 语义）、`Databases`；
   `LocalExecutor` 直连实现，`RemoteExecutor` 预留（返回 `ErrNotImplemented`）。
   `LoadTables` 全程复用同一连接读取各库——PG 的 schema 位于单个 database 内，
   逐库另开连接会失效。
3. `validate/`、Metadata/Query/Execute resolver、`cli/execute.go`、
   `validate.Run` 全部改走 `agent`。
4. 删除 `models.Cluster` 的 `Connect`/`Execute`/`Metadata`/`Databases`/`Stat`/
   `Repack` 及私有 `columns`/`indexes`；`models.Table` 删除。
   `models.Database` 保留为对外 DTO。
5. `go mod tidy` 清出 xorm 及其传递依赖（builder、go-xorm/core、go-json、
   snappy、goleveldb），并删除 `replace` 指令。

## 后果

正面：
- 目标库访问单一入口（`internal/agent`），分层清晰，符合红线 F16 精神
- 元数据库不再有 ORM；依赖树更小
- 分类/字段语义经 `meta` 显式固化，不再依赖第三方内部结构
- 双方言（MySQL/PostgreSQL）差异收敛在 `dialect` 实现内

负面：
- `meta.Table` 与旧 `core.Table` 字段名/形态有差异（如 `Columns` 由
  `[]*Column` + 方法变为 `map`），迁移期需逐一核对
- Metadata resolver 输出的 JSON 形状随 `core.Table` → `meta.Table` 变化
  （仓库内未发现前端消费点，风险可控）

缓解：
- 类型迁移为纯编译期改名（`core.*` → `meta.*`），无反射/断言
- `Resolver`/`validate` 调用点全部编译期替换
- 运行时行为在 B.8 全量验收（build/vet/lint/test）中回归

## 架构

```
internal/agent/
├── agent.go        # Executor 接口 + ExecResult
├── meta/
│   └── meta.go     # Table/Column/Index/Database/ColumnType（中立类型）
├── dialect/
│   ├── dialect.go  # Dialect 接口（LoadTables/Databases/QuoteIdent）
│   ├── mysql.go    # 实现 + 私有 loadTableHeaders/loadColumns/loadIndexes
│   └── postgres.go # 同上，PG 实现
├── conn.go         # openTarget（DSN + 方言选择）
├── tables.go       # Target.LoadTables（复用单连接逐库读取）
├── query.go        # Target.Databases/Query
├── execute.go      # Target.Execute
├── local.go        # LocalExecutor 直连
└── remote.go       # RemoteExecutor 预留
```

## 备选

- 保留 `go-xorm/core` 仅作元数据类型：类型仍绑第三方 fork，且与"去 ORM"方向
  相悖，否决。
- 直接复用 `dialect` 之外的公共库做元数据建模：无必要，字段子集极小。

## 参考

- `docs/TECH_STACK.md`
- `docs/DATABASE_DESIGN.md`
- ADR 0001（目标库执行）
- AGENTS.md（红线 F5/F16）
