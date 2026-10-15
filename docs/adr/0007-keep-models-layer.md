# ADR 0007: 保留 models 层

## 状态

已接受

## 背景

2D-C-B 移除 xorm 后，`models` 层职责混杂：
- 实体（与 domain 重复）
- GraphQL 输入 DTO
- 遗留 ORM 钩子（已无用）
- validate 对外 DTO（Database / Violations）

候选：
- A. 删 models / bridge，全用 domain
- B. 保留 models 为 DTO / bridge 层

## 决策

采用 B：保留 models 层。

理由：
- 全删改动面大（16 实体 + 全项目调用点），风险高
- GraphQL 输出契约与 domain 尚未对齐
- models.Violations / models.Database 仍被 validate 使用
- 分阶段：前端联调后再评估 2D-C-D

## 后果

- models / bridge 保留
- 后续如做 2D-C-D，需整体迁移 + 测试

## 参考

- ADR 0002（表名前缀例外）
- ADR 0004（移除 xorm）

## 2D-C-D 重估（保留结论不变）

- `bridge` 被 caches / cli / crons / events / resolvers / validate 共 6 包使用，
  独立包合理，不内联；命名已是 `XxxToDomain` / `XxxFromDomain`，无需统一；
  删 5 个零引用函数（Avatar/Glossary/Log/Role/Statistic 的 `ToDomain`）。
- `models` 的 `BeforeInsert` / `BeforeUpdate` 不是死代码：承载 UUID 生成、
  TrimSpace、指纹计算、类型缺省，有 20+ 显式调用（resolvers、crons），
  且结果进 GraphQL 响应与事件。搬到 service 会改变响应/事件载荷，
  违背"不改业务逻辑"，故保留不动。
- `models.Violations` 不是被动 DTO：是带 `sync.Mutex` 的并发累加器
  （`Add`/`Append`/`Marshal`/`Clauses`，`Marshal` 内部用 `encoding/json`），
  搬入 `domain` 会改变 domain 的性质，故保留。
- `bridge` 泛型化：评估后不做（32 个函数形状接近但字段映射各异，收益低）。
