# ADR 0001: 目标库执行架构 — 首版直连，后续 Agent

## 状态

已接受

## 背景

平台需将审核通过的 SQL 变更应用到目标库（MySQL / PostgreSQL）。
两类架构：

- A. 中心服务直连目标库
- B. 目标机装 Agent，中心调 Agent

行业最佳实践（Bytebase / Yearning）是 B：Agent 持有 DB 凭据，中心不碰生产库。
但 B 增加部署难度——每台目标机需装 Agent。

首版求快，MVP 功能简单。

## 决策

**首版采用 A（中心直连）**，但**接口预留 B**。

实现：
- `internal/agent/agent.go` 定义 `Executor` 接口
- 首版实现 `LocalExecutor`（中心直连目标库）
- `RemoteExecutor` placeholder（返回 `ErrNotImplemented`）
- `main.go` 按配置选实现，默认 `LocalExecutor`

## 后果

正面：
- 首版部署简单：单二进制
- 接口已抽象，后续切 Agent 不改上层
- 上层（service）只依赖 `Executor` 接口

负面：
- 首版**目标库凭据集中在中心服务**——安全风险
- 中心服务被攻破 = 所有目标库暴露

缓解：
- 凭据**加密存储**（红线 F7）
- 执行必经**审计**（红线 F8）
- 后续切 Agent 后恢复"中心不直连"

## 备选

- **首版直接上 Agent**：安全但部署复杂，MVP 阶段不取。
- **混合**：中心直连测试库，Agent 连生产库——复杂度高，不取。

## 参考

- Bytebase: https://github.com/bytebase/bytebase
- Yearning: https://github.com/cookieY/Yearning
