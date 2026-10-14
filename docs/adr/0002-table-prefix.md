# ADR 0002: 表名前缀 `mm_` 例外 — 保留前缀，代码无前缀

## 状态

已接受

## 背景

AGENTS.md 要求代码命名与 DB 表名 1:1 对齐。
但现有 `halodb.sql` 的表名带 `mm_` 前缀：

```
mm_avatars, mm_clusters, mm_comments, mm_crons, mm_glossaries,
mm_logs, mm_options, mm_queries, mm_edges, mm_roles, mm_rules,
mm_statements, mm_statistics, mm_tickets, mm_users, mm_templates
```

`mm_` 来自历史命名（旧项目）。

## 决策

**保留 `mm_` 前缀**，但**代码命名无前缀**，映射集中在 repository 层。

- 表 `mm_tickets` → `domain.Ticket` / `TicketRepo` / `TicketService` / `TicketResolver`
- SQL 里显式写 `mm_tickets`
- DTO / API / 前端**不感知** `mm_`

## 理由

- **不迁移表名**——迁移风险高，收益低
- **代码去前缀**——符合"命名对齐"精神（1:1 指实体概念，非字面前缀）
- **映射集中**——`internal/repository/mysql/` 的 SQL 显式写 `mm_`，上层不感知

## 后果

正面：
- 不动 `halodb.sql`，零迁移风险
- 代码命名干净
- 映射单点（repository 层）

负面：
- **违反 AGENTS.md 的"1:1 对齐"字面要求**——本 ADR 即例外说明
- 新人需知道 `mm_` 存在

缓解：
- AGENTS.md 明确写例外
- repository 层 SQL 里 `mm_` 一目了然

## 备选

- **迁移去前缀**：`ALTER TABLE mm_tickets RENAME TO tickets;` × 16——风险高，不取。
- **代码也用 `mm_`**：`domain.MMTicket`——丑，不取。

## 参考

- AGENTS.md 命名章节
- `scripts/halodb.sql`
