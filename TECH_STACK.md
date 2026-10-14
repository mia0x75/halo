# 技术栈

> 版本：v0.1.0
> 最后更新：2026-10-14
> 权威性：本文档是技术选型 / 依赖 / 构建 / 部署的权威源。
> 去重：数据库 schema 见 `DATABASE_DESIGN.md`；需求见 `REQUIREMENTS.md`；API 见 `docs/api/`。

---

## 一、后端

| 层 | 选型 | 版本 | 说明 |
|----|------|------|------|
| 语言 | Go | 1.24+ | `go.mod` 指令 |
| Web 框架 | Echo | v5 | 零分配、`net/http` 兼容 |
| API | gqlgen | latest | GraphQL schema-first |
| 元数据库 | MySQL | 8.0+ | 多角色 / 并发 / 审计 |
| 数据访问 | sqlc + sqlx | latest | sqlc 生成类型安全代码；sqlx 做 scan |
| 目标库 driver | go-sql-driver/mysql + pgx | latest | MySQL + PG |
| 认证 | golang-jwt/jwt | v5 | JWT |
| 日志 | go.uber.org/zap | latest | 结构化日志 |
| SQL parser（MySQL） | mia0x75/parser | latest | TiDB fork |
| 规则引擎 | mia0x75/yql | latest | 规则匹配 |
| 测试 | stretchr/testify | latest | assert / require / mock |
| 热重载 | air-verse/air | latest | 开发期 |

**PG parser 说明**：`pgplex/pgparser` 是 `gram.y` 的 Go 移植，AST 1:1 对应 PG 内部结构。

---

## 二、前端

| 层 | 选型 | 版本 | 说明 |
|----|------|------|------|
| 框架 | React | 19 | 函数式 + hooks |
| 语言 | TypeScript | 5.x | 全量 TS |
| 构建 | Vite | latest | 生态成熟 |
| 包管理 | Bun | latest | install 快 |
| 路由 | React Router | 7 | |
| GraphQL 客户端 | urql | latest | 轻，TS 友好 |
| 客户端状态 | Zustand | latest | 纯 UI 状态 |
| UI | shadcn/ui + Radix | latest | 源码入库，可改 |
| 样式 | Tailwind CSS | 4 | |
| 表单 | react-hook-form + zod | latest | |
| 表格 | TanStack Table | latest | |
| SQL 编辑器 | monaco-editor | latest | |
| diff | react-diff-viewer-continued | latest | |
| 代码生成 | GraphQL Code Generator | latest | 生成 TS 类型 + hooks |
| 测试 | Vitest + Testing Library | latest | |

---

## 三、构建与交付

| 层 | 选型 | 说明 |
|----|------|------|
| 前端构建 | Vite | `npm run build` |
| Go 嵌入 | `//go:embed all:dist` | 前端产物嵌入 Go 二进制 |
| SPA Fallback | 自定义 handler | 真实文件直接返回；否则 `index.html` |
| 开发模式 | Vite dev server + Go `air` | Vite 代理 `/api` → `:8080` |
| 生产模式 | `make build` | 单二进制 |

### 交付流程

```bash
# 1. 构建
make build                            # → bin/halo

# 或分步：
cd frontend && bun run build          # → internal/web/dist/
go build -ldflags="-s -w" -o halo ./cmd/server

# 2. 运行
./halo
```

---

## 四、运维与部署

| 层 | 选型 | 说明 |
|----|------|------|
| 部署形态 | Docker + docker-compose | |
| 数据持久化 | 卷挂载 | MySQL 数据 |
| 备份 | `mysqldump` / `pg_dump` | 待定 |
| 监控 | `/health` 端点 | |
| 日志 | 结构化 JSON | 单机场景 |

---

## 五、明确不引入的技术

| 技术 | 不引入理由 |
|------|-----------|
| **ORM（GORM / xorm）** | SQL 全部手写 + sqlc |
| **Redis** | 单机场景，无缓存层需求 |
| **消息队列** | 单机场景，应用层队列覆盖 |
| **Kubernetes** | Docker 单容器交付 |
| **SSR / RSC** | 纯 API 客户端 |
| **pg_query_go** | CGO；改用 pgplex/pgparser（纯 Go） |
| **SQLite 目标库** | 语法分析只 MySQL + PG |

---

## 六、已知约束

| 约束 | 说明 | 应对 |
|------|------|------|
| 元数据库 MySQL | 首版只 MySQL | 后续评估 SQLite / PG |
| 目标库 MySQL + PG | SQLite 不做 | 语法分析只双方言 |
| PG parser | 两套 AST 不兼容 | 统一 IR（`internal/analyzer/ir.go`） |
| CGO | 中心服务禁 CGO | Agent 端可 |

---

## 七、环境搭建

```bash
# 后端
go mod init github.com/mia0x75/halo
go get github.com/labstack/echo/v5
go get github.com/99designs/gqlgen
go get github.com/go-sql-driver/mysql
go get github.com/jackc/pgx/v5
go get github.com/jmoiron/sqlx
go get github.com/golang-jwt/jwt/v5
go get go.uber.org/zap
go get github.com/stretchr/testify
go get github.com/pgplex/pgparser
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/air-verse/air@latest

# 前端
bun create vite frontend --template react-ts
cd frontend
bun add urql graphql react-router@7 zustand
bun add -d @graphql-codegen/cli @graphql-codegen/typescript @graphql-codegen/typescript-urql
bun add -d tailwindcss@4
bunx --bun shadcn@latest init
```

**目标**：`go run ./cmd/server` 起 Echo + MySQL，`curl localhost:8080/health` 返回 `{"status":"ok"}`。

---

## 八、变更记录

| 日期 | 变更 |
|------|------|
| v0.1.0 | 基线 |
