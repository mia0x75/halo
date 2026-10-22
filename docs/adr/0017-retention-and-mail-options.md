# ADR 0017: 保留期与邮件配置落库

## 状态

已接受

## 背景

ADR 0015 只覆盖了指标桶与规则命中明细的清理，`mm_logs`、`mm_messages`、
`mm_tasks` 仍然无上限增长；ADR 0014 的邮件通道则只读配置文件里的 `mail` 段，
改一次 SMTP 就要改文件并重启服务。

两者都是「运行期需要管理员调整的策略」，继续留在文件里与系统设置页的定位不符。

## 决策

**把保留期与 SMTP 配置都放进 `mm_options`，由管理员在系统设置页维护。**

1. **保留期选项**：`retention.metrics-days`、`retention.hits-days`、
   `retention.logs-days`、`retention.messages-days`、`retention.tasks-days`，
   取值 1–3650，缺省分别为 90/90/180/90/30 天。
2. **统一清理作业**：新增 `CronJobEnum.RETENTION_CLEANUP`，固定映射到
   `halocli retention cleanup`；作业自身不接受天数参数，窗口一律读选项，
   所以调整策略不需要重建作业。
3. **邮件选项**：`mail.enabled`、`mail.addr`、`mail.user`、`mail.password`、
   `mail.insecure-skip-verify`。`EmailNotifier` 在发送时通过 provider 读取，
   改完立即生效；配置文件里的 `mail` 段只在该选项从未写入时兜底。
4. **密码加密**：`mail.password` 用 `secret.crypto` 做 AES 加密后 base64 存储
   （密文不是合法 UTF-8，不能直接进文本列）。API 回显固定掩码 `********`，
   原样提交掩码表示「不修改」。
5. **测试邮件**：`sendTestMail(email)` 仅管理员可用，走同一条 notifier 链路；
   邮件未启用时直接报错，不伪装成功。

## 后果

正面：

- 保留期与 SMTP 配置都有审计（`option.update`）与统一入口；
- 调整保留期只需改选项，不必修改 cron 定义；
- 邮件密码不再出现在配置文件与进程参数里。

负面：

- 缺少选项行时行为取决于配置文件，部署初期存在两处来源，需在升级说明里讲清；
- 清理作业按表逐条执行，单次运行时间与过期数据量线性相关；
- 掩码语义下无法通过 API 清空密码，只能改成新的非空值。

## 后续

- 若引入第二套邮件通道（DingTalk/WeChat）的配置，沿用同一 options 前缀；
- 多实例部署时保留期清理仍是每个实例各跑一次，依赖删除操作的幂等性。
