# ADR 0003: PostgreSQL SQL Parser 选型

## 状态

已接受

## 背景

平台需对 PostgreSQL 目标库的 SQL 做语法分析（危险操作检测、元数据提取、审计）。
现有 MySQL 分析用 `mia0x75/parser`（TiDB fork），只支持 MySQL。
PG 需独立 parser。

约束：
- 中心服务禁 CGO（红线 F3）
- 纯 Go 优先
- AST 兼容 PG 官方语法

## 候选

| 项目 | 语言 | 维护 | CGO | PG 兼容 | 评价 |
|------|------|------|-----|---------|------|
| ValkDB/postgresparser | Go (ANTLR) | 活跃 | ✅ | 非官方语法 | IR 提取为主，非官方语法 |
| aipescience/queryparser | Python | 活跃 | — | 旧 PG | 不适用 |
| **pgplex/pgparser** | Go (goyacc) | 活跃 | ✅ | **100% PG 17.7** | **官方 gram.y 移植** |
| auxten/postgresql-parser | Go (goyacc) | 低 | ✅ | CockroachDB 风格 | 非纯 PG，重 |

## 决策

采用 **`pgplex/pgparser`**。

理由：
- 纯 Go，无 CGO
- 官方 `gram.y` 移植，AST 1:1 对应 PG 内部结构
- Bytebase 背书（同领域：SQL 治理）
- 回归测试 45,000 条，99.6% 通过
- 线程安全

## 后果

正面：
- 无 CGO——符合红线 F3
- 官方语法——权威
- 中心服务可引（都纯 Go）

负面：
- 两套 AST（MySQL / PG）差异大——需统一 IR
- PG parser 较新——成熟度待观察

缓解：
- 统一 IR：`internal/analyzer/ir.go`，上层只依赖 IR
- PG IR 转换首版 placeholder，规则检测走通用规则
- 出问题回退到通用规则检测（关键字 / 正则）

## 架构

```
internal/analyzer/
├── analyzer.go      # Analyzer 接口 + 方言分派
├── ir.go            # 统一 IR
├── mysql.go         # mia0x75/parser → IR
├── postgres.go      # pgplex/pgparser → IR（首版 placeholder）
├── rules/
│   ├── common.go
│   ├── mysql.go
│   └── postgres.go
└── errors.go
```

## 备选

若 `pgplex/pgparser` 不稳定：
- 短期：通用规则检测（MVP 够用）
- 长期：评估 `auxten/postgresql-parser`（重，但成熟）

## 参考

- https://github.com/pgplex/pgparser
- https://github.com/auxten/postgresql-parser
- https://github.com/ValkDB/postgresparser
- TECH_STACK.md 后端章节
