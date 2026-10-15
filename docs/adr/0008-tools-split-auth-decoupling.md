# ADR 0008: tools 拆分与鉴权断边

## 状态

Accepted

## 背景

`tools` 是根目录杂烩包，9 个文件承载 6 类无关职责（加密、邮件、进程、密码、JWT、调试），
且 `tools/token.go` 反向依赖 `caches`/`g`/`gqlapi`/`models`，使 `tools` 拖入近全应用依赖。

## 决策

- `internal/platform/crypto`：保留 `EncryptAES`/`DecryptAES`（删 5 个零引用哈希/填充函数）
- `internal/auth`：纯 JWT 签发/验签，**零 mia0x75 依赖**（`go list -deps` 仅 jwt-go + stdlib）
- `internal/transport/auth`：token 内存 store + `Credential` + 用户/角色解析（依赖 caches/models/gqlapi）
- `internal/platform/mail`：`MailDto`（零引用，迁后保留）
- `internal/tools`：`PanicDetail`/`TimeoutedExec`/`Contains`（`Contains` 泛型化）
- 删 `tools/secure.go`/`string.go`/`tools.go`

## 关键约束

- `caches` 已 import `service`（`caches.Init(r service.Repos)`），故 `service` 不能反向依赖
  `caches`。任务书建议"store 搬 service"会成环，因此解析逻辑归 `internal/transport/auth`。
- `JWT.issuer`/`secret` 从未从配置接线（历史潜在 bug）。空 secret 保持旧行为不变，避免
  token 全失效（行为变更），待后续接线。

## 后果

- `internal/auth` 成为可独立测试的纯 JWT 层
- `Credential` 值类型在 context 中统一传递（middleware 存 `*credential` 解引用值）
- `FromContext` 零引用已删
