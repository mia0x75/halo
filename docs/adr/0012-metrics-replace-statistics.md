# ADR 0012: 指标表合并统计表

## 状态

已接受

## 背景

`mm_statistics` 用「分组（group）+ 键（key）」表达指标，写入路径只有一种：

```sql
ON DUPLICATE KEY UPDATE value = value + VALUES(value)
```

于是 counter 与 gauge 混在一张表里：

- `total-tickets` 之类的 counter 语义正确（原子累加，不丢并发更新）；
- `latest-login` 之类的 gauge 却被**累加**——`internal/event/handlers.go` 里
  `bump(ctx, repo, e.User.UUID, "latest-login", float64(time.Now().Unix()))`
  会把登录时间戳越加越大，而不是记录最近一次登录时间。

同时这张表承载不了新 Dashboard 需要的查询：

- 分组列混装了三种语义（`overall` / `tickets-daily` / 用户 UUID），无法按维度筛选；
- 按天的指标把日期塞进键里（`tickets-daily` + `2026-10-20`），没有独立的时间桶，
  做不了通用的时间窗口聚合；
- 表里没有 counter / gauge 的类型标记，写入路径无法区分。

## 决策

**用 `mm_metrics` 取代 `mm_statistics`，把「分组 + 键」拆成「名称 + 作用域 + 作用域实体 + 时间桶」，并显式区分 counter 与 gauge。**

`mm_metrics` 的自然键是 `(name, scope, scope_id, bucket)`：

| 维度 | 取值 | 说明 |
|------|------|------|
| `name` | 如 `total-tickets`、`total-logins`、`latest-login` | 指标名 |
| `scope` | 0 全局 / 1 按天 / 2 用户 / 3 群集 / 4 分组 | 作用域类型 |
| `scope_id` | 作用域实体的 UUID，全局作用域为空串 | 作用域实例 |
| `bucket` | 时间桶，如 `2026-10-20`，非按天作用域为空串 | 时间维度 |
| `type` | 0 counter / 1 gauge | 写入语义 |

写入路径随之分开：

- counter：`value = value + delta`（原子 UPSERT，并发不丢更新）；
- gauge：`value = new_value`（覆盖）。

`latest-login` 因此改为 gauge 写入，登录时间不再被累加。

配套变化：

- 领域层 `domain.Statistic` → `domain.Metric`，`StatisticRepository` → `MetricRepository`，
  仓储接口新增 `Increment` 与 `SetGauge` 两条语义明确的写入方法；
- GraphQL 契约 `type Statistic` → `type Metric`，`User.Statistics` → `User.Metrics`，
  `QueryRoot.statistics(groups)` → `QueryRoot.metrics(scope, names)`（未发布，可破坏性变更）；
- CLI `halocli statistics` → `halocli metrics`；
- 迁移 `0004_metrics_replace_statistics.sql` 在两个引擎各一份：建表 → 搬迁旧数据 →
  删旧表，可重复执行（目标库已按权威 schema 建好时只补台账）。

## 后果

正面：

- counter 与 gauge 语义分离，`latest-login` 的累加 bug 被修掉；
- 作用域与时间桶成为一等维度，Dashboard 的「近 7 天趋势」「按群集统计」可以直接查询；
- 旧的「分组」字符串不再承载多义，减少了一层隐式约定。

负面：

- 旧数据搬迁时 `tickets-daily` / `queries-daily` 的指标名被规范化成
  `total-tickets` / `total-queries`，历史行里 `group` 与 `key` 的组合语义不再保留；
- GraphQL 契约变更，前端与外部调用方需要同步；
- `mm_statistics` 不再存在，历史 review 文档里的表名仅作历史记录。

## 后续

- 指标明细（规则命中明细）落在 `mm_violations`，聚合计数落在 `mm_metrics`：
  明细可追到工单/语句/规则，聚合只保留计数，两者不重复；
- 指标的时序桶需要保留期与清理策略，挂到定时任务执行。
