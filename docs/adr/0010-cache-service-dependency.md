# ADR 0010: cache 依赖 service 的折中

## 状态

已接受

## 背景

`internal/cache` 的 `Init(repos service.Repos)` 依赖 `service`。
而 `service` 若依赖 `cache`（读缓存），则成环。

架构审核（P.1）标记此为边界模糊。

## 决策

保留 `cache → service.Repos` 依赖。

理由：
- `Repos` 是接口集合（`domain.XxxRepository`），不是 `service` 逻辑。
- 打破 `cache → service → cache` 环的唯一干净方式：`service` 不读 `cache`。
- 当前 `service` 不读 `cache`（仅 `resolvers` 层读），依赖方向 OK。

## 后果

- `cache` 与 `service` 同层（都在 `internal/`）。
- 若未来 `service` 需读 `cache`——需重新设计（如 cache 变 repository 的装饰器）。

## 参考

- `docs/review/architecture-review-2026-10-16.md`（A2）
- ADR 0007（models / bridge 保留）
