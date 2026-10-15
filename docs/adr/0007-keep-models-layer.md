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
