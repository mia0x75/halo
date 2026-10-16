# ADR 0009: validate 分层与 Violations 例外

## 状态

已接受

## 背景

架构审核发现 `internal/validate`（约 7600 行规则引擎）内部直接依赖 `internal/models`
（GraphQL 契约 DTO）。依赖方向应为 `validate → domain`，`models` 只应出现在
`Run` 边界。ADR 0007 保留 models 层，但未界定 validate 内部的类型边界。

约束：

- `Run` 签名冻结：`Run(stmts []*models.Statement, cluster *models.Cluster, ticket *models.Ticket)`
- 不改规则逻辑；不改 `mm_rules` 表
- `domain` 零外部依赖（不能引入 `ast` / `regexp` / `json`）

## 决策

- validate 内部类型改 `domain`：`Rule` / `Cluster` / `Ticket`
- 新增 `validate.Subject` 取代引擎内的 `models.Statement`：
  - 持久化字段（TicketID / Sequence / UUID / Content / Type / Status）
  - 运行期字段 `StmtNode ast.StmtNode`、`Violations *models.Violations`
  - `domain.Statement` **不能**承载 `StmtNode`（`ast` 是外部包，破坏 domain 零依赖），
    故由 validate 自有类型承载
- 新增 `validate.Rule`（嵌入 `domain.Rule`）+ 懒解析 `Int` / `StringSlice` / `Regexp`
  （带 `sync.Mutex` 保护的 memo cache，供并发 validator 组使用）。`domain.Rule` 不引
  `regexp` / `json`
- 启动校验 `ValidateRules([]*domain.Rule)`：校验 Name / Func 非空、Name 去重、按
  `Operator`（`in` / `regexp`）试解析 `Values`。在 `cmd/server/main.go` 的 `cache.Init`
  之后调用，失败即 `Fatalf`
- `SetGroup` / `Run` 在边界处用 `bridge.RuleToDomain` 等把 `models.Rule` 转成引擎类型
- `RuleSource.Load` 直接返回 `[]*domain.Rule`（仓储行本就是 domain），删除
  domain→models→domain 的往返转换

### Violations 例外

`models.Violations`（含 `sync.Mutex`）与其元素 `models.Clause` 保留在 models：

- 并发累加器，搬入 `domain` 会改变 domain 性质（ADR 0007）
- `Violations` 由 `Subject` 持有；`Run` 边界处 `Subject.Violations` 与
  `models.Statement.Violations` **共享同一指针**，规则写入即回写，无需显式拷贝
- `Clause` 在 validate 内以**类型别名** `type Clause = models.Clause` 暴露，使 172 处
  规则构造点无需 import `models`

结果：`internal/models` 在 validate 内只出现于 `subject.go`（Subject/Violations/Clause）
与 `validate.go`（Run 边界 + Database 例外）。

### Database 例外

`models.Database` 为 validate 对外读取目标库元数据所用，保留（ADR 0007）。

## 后果

- `validate` 规则文件（`validate_*.go`）不再 import `models`
- `domain` 保持零外部依赖
- 去反射（V.2）与事件总线（P.3）不在本 ADR 范围，另行决策
- `Context.Rules` 仍是遗留的全局缓存通道，待 2D-B 移除

## 参考

- ADR 0007（保留 models 层）
- `docs/review/architecture-review-2026-10-16.md`（问题 A1）
