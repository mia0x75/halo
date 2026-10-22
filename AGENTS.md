# AGENTS.md

> AI 编码协作契约。指令式，非建议。违反红线一律打回。
>
> 阅读顺序：本文件 → `docs/adr/` → `scripts/halodb_schema_mysql.sql`。

---

## 项目速览

**SQL 变更工单与审计平台** — SQL 变更的提交、审核、执行、审计一站式平台。

**战略前提（不得违背）**：
- 多组别用户，组间权限隔离
- 多角色（提交人 / 审核人 / 管理员）
- 执行必经审计
- 目标库：MySQL + PostgreSQL
- 元数据库：MySQL + PostgreSQL

**当前阶段**：v0.1 MVP。本文件为 AI 协作契约，`docs/adr/` 为架构决策权威源。

---

## 核心原则

- **KISS**：优先简单、可维护的方案。不做过度设计。
- **事实优先**：先看源码、配置、日志，再引用历史结论。禁止凭记忆断言。
- **可验证目标**：优化必须有测量依据。
- **中文回复**：所有用户可见回复、方案和任务清单用中文。

---

## 技术栈（不得偏离）

| 层 | 选型 |
|----|------|
| 后端语言 | Go 1.26+ |
| Web 框架 | Echo v5 |
| API | GraphQL（gqlgen） |
| 元数据库 | MySQL + PostgreSQL |
| 数据访问 | sqlc + sqlx |
| 目标库 | MySQL + PostgreSQL |
| 认证 | JWT + RBAC |
| 前端框架 | React 19 |
| 前端构建 | Vite |
| 前端状态 | urql（GraphQL）+ Zustand |
| 前端路由 | React Router 7 |
| 前端 UI | shadcn/ui + Radix + Tailwind 4 |
| 前端表单 | react-hook-form + zod |
| 前端表格 | TanStack Table |
| 前端 SQL 编辑 | CodeMirror 6（@uiw/react-codemirror） |
| 前端测试 | Vitest + Testing Library |

**换栈唯一途径**：改本文件技术栈表并由人类确认。AI 不得单方面换栈。

---

## 工作流程

- **复杂任务**：调研 → 方案 → 确认 → 分解 → 执行 → 验证。
- **简单任务**：直接执行 + 必要验证。
- **外部依赖**：涉及外部库、API、CLI 或可能变化的行为时，查阅最新官方文档和源码。禁止凭记忆调用。
- **状态优先级**：当前源码与运行状态 > 历史 session / memory / 截图。

---

## 目录结构

```
halo/
├── cmd/server/main.go
├── internal/
│   ├── domain/           # 实体 + 仓储接口（零外部依赖）
│   ├── service/          # 业务编排
│   ├── repository/       # 数据访问（唯一 import database/sql）
│   ├── graphql/          # GraphQL resolver（薄）
│   ├── agent/            # 目标库执行器（首版 LocalExecutor 直连；RemoteExecutor 预留）
│   ├── analyzer/         # SQL 语法分析（MySQL: mia0x75/parser；PG: pgplex/pgparser；统一 IR）
│   ├── platform/         # db / log / config
│   ├── auth/             # JWT + RBAC
│   └── http/             # 中间件
├── migrations/           # schema 迁移
├── frontend/             # React 19 项目
├── docs/
│   ├── adr/              # 架构决策记录
│   └── ...
└── testdata/
```

**依赖方向（轻量 Clean）**：
- `transport/* → service → domain ← repository`。禁止反向；`transport/auth` 不依赖 GraphQL 生成层。
- `domain` 不 import 任何外部库（无 `database/sql` / `echo` / `zap` / `json`），仅标准库。
- 接口定义在 `domain`，实现在 `repository`。
- `repository` 是唯一 import `database/sql` 的地方。
- `models` 是 GraphQL 契约 DTO，`bridge` 只做 models↔domain 转换（被 resolvers / caches 等消费），两者不写业务逻辑。
- `cmd/server/main.go` 负责组装与依赖注入。
- 不做：UseCase 类、Presenter、Mapper 层、Event Bus、CQRS。

---

## 编码规范

### Go

- 所有导出函数第一个参数 `context.Context`
- 错误必须 `fmt.Errorf("...: %w", err)` 包装
- 禁止业务错误用 `panic`（仅 `main` 初始化失败可用）
- 禁止全局变量存状态（配置除外）
- 禁止 `init()` 做重活
- IO 密集场景优先异步实现

### TypeScript / React

- 组件用函数式 + hooks
- Props 用 `interface`，显式声明类型
- 禁止组件内直接 `fetch`（必须走 `src/api/`）
- 禁止 `any` 和 `as any`
- 禁止 `v-html`（React 对应 `dangerouslySetInnerHTML`）渲染未净化内容

### 注释

- **语言：英文**
- **只写"为什么"，不写"是什么"**
- **不复述函数签名**
- **不写显而易见的**
- **TODO / FIXME 要具体**

```go
// ✅
// DSN format differs by driver.
switch c.Type {

// ❌
// This function checks the cluster type and then
// decides which DSN format to use.
```

### 命名

| 场景 | 约定 |
|------|------|
| Go 文件 | `snake_case.go` |
| Go 接口 | 单方法用 `-er` 后缀 |
| DB 表 / 列 | `snake_case`（**保留 `mm_` 前缀**） |
| API 路径 | `/api/kebab-case` |
| React 组件 | `PascalCase.tsx` |
| Zustand store | `useXxxStore` |

**表名前缀例外（ADR 0002）**：DB 表名保留 `mm_` 前缀，代码命名无前缀，映射集中在 repository 层。

- 表 `mm_tickets` → `domain.Ticket` / `TicketRepo` / `TicketService` / `TicketResolver`
- SQL 里显式写 `mm_tickets`

---

## 枚举规范

- **存储值用整数或字符串常量**：状态列优先 `INTEGER`，类型列可用 `VARCHAR`（如 `clusters.type`）
- 枚举常量定义在 `domain`
- **代码禁止裸数字/裸字符串**，一律引用常量
- 显示值走前端映射

---

## 数据库

### 元数据库（MySQL / PostgreSQL）

- **参数化查询**，禁止字符串拼接 SQL
- **禁止 `SELECT *`**，显式列出字段
- **禁止无 `WHERE` 的 `UPDATE` / `DELETE`**
- **禁止 N+1 查询**，必须批量化
- **Row Scan 用 sqlx**：`db.GetContext` / `db.SelectContext`
- **SQL 全部手写**：sqlc 生成，无 ORM
- **双引擎并行包**：`internal/repository/mysql` 与 `internal/repository/postgres`
  实现同一套 `domain` 仓储接口，`database.driver`（`mysql`/`postgres`，缺省按
  DSN 自动识别）选择；禁止跨包共享查询与转换代码
- **PG 方言**：占位符 `$n`、自增取号用 `RETURNING`、upsert 用 `ON CONFLICT`、
  保留字双引号、时间列 `TIMESTAMPTZ`、UPDATE 显式 `update_at = NOW()`
 （MySQL 侧靠 `ON UPDATE CURRENT_TIMESTAMP`）

### 目标库（MySQL / PostgreSQL）

- **只通过 `internal/agent` 访问**
- 凭据加密存储
- 执行必经审计

### Schema 迁移

- 文件：`migrations/<engine>/NNNN_name.sql`（MySQL / PostgreSQL 各一份，同号同内容）
- 台账：`mm_migrations`（元数据库），由 `halocli migrate up` 写入
- 执行：`make migrate-up`（需 `HALO_CFG`）；**服务启动不自动迁移**
- 迁移前暂停写入 + 备份
- 迁移在事务内，失败回滚（PostgreSQL）；MySQL 的 DDL 隐式提交，失败需人工修复
- 已应用的迁移文件不可修改（checksum 校验），改错只能追加新迁移
- 基线与权威 schema 脚本同步：`make check-migrations`
- 详见 `docs/adr/0011-schema-migration-ledger.md`

---

## API 约定

### GraphQL

- Schema 在 `gqlapi/schema.graphql`
- Resolver 只做参数转换 + 调 service
- **禁止 resolver 写业务逻辑**
- 分页用 **cursor**（keyset），不用 offset

### 响应格式

GraphQL 标准 `{ data, errors }`。

### 中间件顺序

`Recover → Logger → CORS → JWT → Handler`

---

## 执行代理（Agent）

**首版无 Agent**：中心服务直连目标库（`LocalExecutor`）。

**后续**：目标机装 Agent，中心调 Agent（`RemoteExecutor`）。

**接口**：`internal/agent/agent.go` 的 `Executor`。

**红线**：
- 目标库凭据**加密存储**
- 执行必经**审计**
- 中心服务**不 import** `internal/agent`（仅 service 层经接口调用）

---

## SQL 语法分析

- **MySQL**：`mia0x75/parser`（现有）
- **PostgreSQL**：`pgplex/pgparser`（纯 Go，100% PG 17.7 兼容，Bytebase 支持）
- **危险操作检测**：通用规则（关键字 / 正则）
- **规则引擎**：`yql`（保留）
- **统一 IR**：`internal/analyzer/ir/ir.go`（两套 AST 输出同一结构）
- **选型依据**：ADR 0003

---

## 测试

| 层 | 覆盖率 |
|----|--------|
| domain | ≥ 80% |
| repository | 关键路径 100% |
| service | 分支覆盖 |
| graphql | 关键 resolver 100% |

---

## Git 提交

```
<type>(<scope>): <subject>

[body]

[footer]
```

**type**：`feat` / `fix` / `refactor` / `perf` / `test` / `docs` / `chore`

**禁止**：
- 一个 commit 干多件事
- `git commit -m "update"`
- 提交未通过的测试
- 提交含调试代码 / `fmt.Println`

---

## 任务执行流程

### 1. 任务理解

- 确认对应需求文档章节
- 与需求冲突 → 先问人类
- 需要新表 / 新字段 → 先确认 `scripts/halodb_schema_mysql.sql`（PG 同步改 `scripts/halodb_schema_postgres.sql`）

### 2. 方案设计

- 多文件改动 → 先输出方案（≤ 200 字）
- 算法改动 → 明确输入输出 + 边界 + 测试用例
- Schema 改动 → 先写迁移脚本

### 3. 编码

- 按上述规范
- 写测试
- 组织 commit

### 4. 自检（必做）

- [ ] `go build ./...` 通过
- [ ] `go test ./...` 通过
- [ ] `go vet ./...` 无警告
- [ ] `golangci-lint run` 无错误
- [ ] 前端 `npm run build` 通过
- [ ] 关键路径有测试
- [ ] 无 `TODO` / `FIXME` 遗留
- [ ] 无调试代码
- [ ] 无 `any` / `as any`
- [ ] commit message 符合规范

### 5. 交付

- 输出：改动清单 + 测试结果 + 关联需求编号
- 引入 ADR → 写入 `docs/adr/`
- 发现需求硬 bug → 先报告，不擅自改

---

## 红线（违反即打回）

| # | 禁止 | 理由 |
|---|------|------|
| F1 | 换技术栈 | 技术栈表已冻结 |
| F2 | 改需求文档 | 已冻结；硬 bug 先报告 |
| F3 | 引入 CGO 依赖（中心服务） | 破坏单二进制；Agent 端可 |
| F4 | 引入多租户 / `tenant_id` | 无多租户 |
| F5 | 引入 Redis / MQ / ORM | 与极简目标冲突 |
| F6 | 字符串拼接 SQL | 参数化查询 |
| F7 | 目标库凭据明文存储 | 安全底线 |
| F8 | 执行绕过审计 | 审计强制 |
| F9 | 前端裸 `fetch` | 违反 API 封装 |
| F10 | 提交调试代码 / 硬编码密钥 | 安全底线 |
| F11 | `panic` 用于业务错误 | 破坏稳定性 |
| F12 | 无测试的关键逻辑改动 | 行为变更须有测试 |
| F13 | TypeScript 使用 `any` / `as any` | 破坏类型安全 |
| F14 | 凭记忆调用外部 API | 版本漂移 |
| F15 | resolver 写业务逻辑 | 分层破坏 |
| F16 | 中心服务 import `internal/agent` | 分层破坏 |
| F17 | 用 `any` 类型断言 PG parser AST（用统一 IR） | 破坏类型安全 |

---

## 常见陷阱

| 陷阱 | 正确做法 |
|------|----------|
| 循环里逐条查询 | 批量化 |
| `COALESCE` 对空结果集无效 | 外层包 `COALESCE((SELECT ...), 'default')` |
| 计算函数内部 save | 计算与持久化分离 |
| SQL 字符串拼接 | 参数化查询 |
| 前端 `any` 泛滥 | 显式类型定义 |
| 提交前不跑 `go test` | 每次 commit 前必跑 |
| 凭记忆用旧版 Echo / React API | 查官方文档 |
| `type` 是 SQL 保留字 | MySQL 反引号 / PG 双引号 |
| resolver 里写业务 | 搬到 service |
| 目标库查询在中心服务 | 走 `internal/agent` |
| 审计不写 | 写操作必经审计 |
| 注释写"是什么" | 写"为什么" |
| 注释用中文 | **英文** |
| 注释 AI 味 | 贴近人类 |
| 两套 parser AST 直接混用 | 走统一 IR（`internal/analyzer/ir/ir.go`） |

---

## Sensitive content guardrails

### 1. 禁止主动读取的文件与目录
- `.env`, `.env.*`, `*.key`, `*.pem`, `*.p12`, `*.jks`, `*.keystore`
- `id_rsa`, `id_ed25519`, `*.token`, `*.session`, `*.log`, `*.sql`, `*.dump`
- `fixtures/`, `mocks/`, `testdata/`, `i18n/`, `locales/`, `regions/`, `countries/`
- 任何包含真实用户手机号、邮箱、身份证、地址、订单号的测试数据文件

### 2. 禁止在回答中输出以下内容
- 完整密钥、API Key、Token、JWT、私钥、证书、数据库连接串
- 真实或疑似真实的手机号、身份证号、邮箱、地址、车牌号
- 完整国家/地区/时区/货币枚举表，尤其是涉及港澳台、争议地区、政治敏感地区的列表
- 内网 IP、域名、端口、账号密码、弱口令示例
- 任何可能被网关敏感词系统命中的政治、宗教、暴力、色情、违禁词内容

### 3. 遇到敏感内容的处理方式
- 如果代码或文件中包含密钥/Token/私钥，只说明“这里存在凭据，需要移到环境变量或密钥管理系统”，不要复述原文。
- 如果测试数据包含手机号、身份证、邮箱，统一替换为 `MOCK_PHONE`、`MOCK_ID`、`MOCK_EMAIL`。
- 如果项目包含 i18n/地区列表，不要展开枚举；只概括“这里有一份地区/语言配置”。
- 如果 shell 命令输出、日志、git log、配置文件可能携带真实数据，只引用路径和结构，不粘贴完整内容。
- 如果用户提问本身可能触发风控，先提示风险，再给出合规、抽象化的技术方案。

### 4. 工具调用约束
- 执行 `cat`、`grep`、`git log`、`ps`、`env`、`printenv` 前，先确认不会读取密钥、日志、测试数据或环境配置。
- 不要为了调试而读取 `.env`、日志文件、数据库导出、上传目录或用户上传文件。
- 不要将完整文件内容粘贴到回答里；优先给出路径、结构、关键行号或脱敏后的片段。

### 5. 输出风格
- 解释代码时，只保留必要的业务逻辑和结构说明。
- 不替用户生成、补全、猜测真实凭据。
- 不输出可能被内容安全系统误判的测试数据、地区列表、政治敏感示例。

---

## 参考文档

- `AGENTS.md` — 本文件（技术栈、规范、红线）
- `docs/adr/` — 架构决策
- `gqlapi/schema.graphql` — GraphQL schema
- `scripts/halodb_schema_mysql.sql` — MySQL 表结构权威源
- `scripts/halodb_schema_postgres.sql` — PG 表结构

---

> **核心原则**：能跑、能测、能提交的代码，比完美的设计文档更有价值。
>
> **不做**：不改需求、不换栈、不引新技术、不做无测试的算法改动、不用 `any`、不凭记忆调 API。
>
> **必做**：显式错误、参数化查询、审计、测试、自检清单全过。
