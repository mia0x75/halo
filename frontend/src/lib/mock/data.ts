/**
 * 演示用的假数据。界面尚未接入后端，数据形态按 GraphQL 契约的 camelCase 镜像，
 * 接入时只需把这里的取值换成接口返回。
 */

export type TicketStatus =
  | "waiting-vld"
  | "vld-failure"
  | "vld-warning"
  | "waiting-mrv"
  | "mrv-failure"
  | "ltgm"
  | "done"
  | "exec-failure"
  | "closed"

export interface Group {
  id: number
  uuid: string
  name: string
  description: string
  status: "normal" | "blocked"
}

export interface Cluster {
  id: number
  uuid: string
  alias: string
  host: string
  ip: string
  port: number
  type: "mysql" | "postgres"
  group: string
  status: "normal" | "blocked"
  latencyMs: number
  version: string
}

export interface User {
  id: number
  name: string
  email: string
  role: "提交人" | "审核人" | "管理员"
  group: string
}

export interface Ticket {
  id: number
  uuid: string
  subject: string
  database: string
  cluster: string
  status: TicketStatus
  submitter: string
  reviewer: string
  createdAt: string
  updatedAt: string
  statementCount: number
  violationCount: number
}

export interface Rule {
  id: number
  name: string
  group: string
  level: "提示" | "错误" | "警告"
  enabled: boolean
  description: string
  hitCount: number
}

export interface Violation {
  id: number
  rule: string
  ticket: string
  level: "提示" | "错误" | "警告"
  message: string
  createdAt: string
}

export interface AuditLog {
  id: number
  actor: string
  operation: string
  target: string
  ip: string
  createdAt: string
}

export interface TrendPoint {
  date: string
  submitted: number
  executed: number
  failed: number
}

export interface SeverityPoint {
  severity: "致命" | "高危" | "中危" | "低危"
  count: number
}

export const groups: Group[] = [
  { id: 1, uuid: "g-1", name: "核心交易", description: "订单、支付相关库", status: "normal" },
  { id: 2, uuid: "g-2", name: "用户中心", description: "账号与权限库", status: "normal" },
  { id: 3, uuid: "g-3", name: "报表分析", description: "只读报表库", status: "normal" },
  { id: 4, uuid: "g-4", name: "离线数据", description: "归档与备份库", status: "blocked" },
  { id: 5, uuid: "g-5", name: "灰度环境", description: "预发布验证库", status: "normal" },
]

export const clusters: Cluster[] = [
  { id: 1, uuid: "c-1", alias: "prod-trade-01", host: "db-trade-01.internal", ip: "10.20.30.11", port: 3306, type: "mysql", group: "核心交易", status: "normal", latencyMs: 3, version: "8.0.36" },
  { id: 2, uuid: "c-2", alias: "prod-trade-02", host: "db-trade-02.internal", ip: "10.20.30.12", port: 3306, type: "mysql", group: "核心交易", status: "normal", latencyMs: 4, version: "8.0.36" },
  { id: 3, uuid: "c-3", alias: "prod-user-01", host: "db-user-01.internal", ip: "10.20.31.21", port: 3306, type: "mysql", group: "用户中心", status: "normal", latencyMs: 2, version: "8.0.36" },
  { id: 4, uuid: "c-4", alias: "prod-user-pg", host: "db-user-pg.internal", ip: "10.20.31.22", port: 5432, type: "postgres", group: "用户中心", status: "normal", latencyMs: 6, version: "16.4" },
  { id: 5, uuid: "c-5", alias: "prod-report-01", host: "db-report-01.internal", ip: "10.20.32.31", port: 3306, type: "mysql", group: "报表分析", status: "normal", latencyMs: 9, version: "8.0.36" },
  { id: 6, uuid: "c-6", alias: "prod-report-pg", host: "db-report-pg.internal", ip: "10.20.32.32", port: 5432, type: "postgres", group: "报表分析", status: "normal", latencyMs: 11, version: "16.4" },
  { id: 7, uuid: "c-7", alias: "prod-archive-01", host: "db-archive-01.internal", ip: "10.20.33.41", port: 3306, type: "mysql", group: "离线数据", status: "blocked", latencyMs: 38, version: "8.0.36" },
  { id: 8, uuid: "c-8", alias: "staging-trade-01", host: "db-staging-01.internal", ip: "10.20.40.11", port: 3306, type: "mysql", group: "灰度环境", status: "normal", latencyMs: 1, version: "8.0.36" },
]

export const users: User[] = [
  { id: 1, name: "张三", email: "zhangsan@example.com", role: "管理员", group: "核心交易" },
  { id: 2, name: "李四", email: "lisi@example.com", role: "提交人", group: "用户中心" },
  { id: 3, name: "王五", email: "wangwu@example.com", role: "审核人", group: "核心交易" },
  { id: 4, name: "赵六", email: "zhaoliu@example.com", role: "提交人", group: "报表分析" },
]

export const tickets: Ticket[] = [
  { id: 1024, uuid: "t-1024", subject: "订单表新增归档字段", database: "app_order", cluster: "prod-trade-01", status: "waiting-mrv", submitter: "李四", reviewer: "王五", createdAt: "2026-10-20 09:12", updatedAt: "2026-10-20 09:48", statementCount: 2, violationCount: 0 },
  { id: 1023, uuid: "t-1023", subject: "用户表增加手机号索引", database: "app_user", cluster: "prod-user-01", status: "done", submitter: "李四", reviewer: "王五", createdAt: "2026-10-19 14:03", updatedAt: "2026-10-19 15:20", statementCount: 1, violationCount: 0 },
  { id: 1022, uuid: "t-1022", subject: "重建日报视图", database: "app_report", cluster: "prod-report-pg", status: "vld-failure", submitter: "赵六", reviewer: "王五", createdAt: "2026-10-19 10:26", updatedAt: "2026-10-19 10:28", statementCount: 1, violationCount: 2 },
  { id: 1021, uuid: "t-1021", subject: "清理历史订单", database: "app_order", cluster: "prod-trade-02", status: "exec-failure", submitter: "李四", reviewer: "王五", createdAt: "2026-10-18 16:41", updatedAt: "2026-10-18 16:46", statementCount: 1, violationCount: 1 },
  { id: 1020, uuid: "t-1020", subject: "支付表拆分字段", database: "app_pay", cluster: "prod-trade-01", status: "ltgm", submitter: "张三", reviewer: "王五", createdAt: "2026-10-18 11:05", updatedAt: "2026-10-18 17:30", statementCount: 3, violationCount: 0 },
  { id: 1019, uuid: "t-1019", subject: "账号表新增状态列", database: "app_account", cluster: "prod-user-pg", status: "vld-warning", submitter: "李四", reviewer: "王五", createdAt: "2026-10-18 09:32", updatedAt: "2026-10-18 09:33", statementCount: 1, violationCount: 1 },
  { id: 1018, uuid: "t-1018", subject: "报表库加只读账号", database: "app_report", cluster: "prod-report-01", status: "closed", submitter: "赵六", reviewer: "王五", createdAt: "2026-10-17 15:18", updatedAt: "2026-10-17 16:02", statementCount: 2, violationCount: 0 },
  { id: 1017, uuid: "t-1017", subject: "灰度库同步结构", database: "app_order", cluster: "staging-trade-01", status: "waiting-vld", submitter: "张三", reviewer: "王五", createdAt: "2026-10-17 10:09", updatedAt: "2026-10-17 10:09", statementCount: 1, violationCount: 0 },
  { id: 1016, uuid: "t-1016", subject: "归档表分区调整", database: "app_archive", cluster: "prod-archive-01", status: "mrv-failure", submitter: "赵六", reviewer: "王五", createdAt: "2026-10-16 14:55", updatedAt: "2026-10-16 16:40", statementCount: 2, violationCount: 3 },
  { id: 1015, uuid: "t-1015", subject: "用户标签表新增字段", database: "app_user", cluster: "prod-user-01", status: "done", submitter: "李四", reviewer: "王五", createdAt: "2026-10-16 09:21", updatedAt: "2026-10-16 11:03", statementCount: 1, violationCount: 0 },
  { id: 1014, uuid: "t-1014", subject: "订单表字段注释补全", database: "app_order", cluster: "prod-trade-02", status: "done", submitter: "张三", reviewer: "王五", createdAt: "2026-10-15 13:47", updatedAt: "2026-10-15 15:12", statementCount: 5, violationCount: 0 },
  { id: 1013, uuid: "t-1013", subject: "移除无用索引", database: "app_order", cluster: "prod-trade-01", status: "done", submitter: "李四", reviewer: "王五", createdAt: "2026-10-15 08:30", updatedAt: "2026-10-15 09:02", statementCount: 1, violationCount: 0 },
]

export const rules: Rule[] = [
  { id: 1, name: "DEL-L2-001", group: "删除", level: "警告", enabled: true, description: "禁止没有 WHERE 的删除", hitCount: 42 },
  { id: 2, name: "UPD-L2-001", group: "更新", level: "错误", enabled: true, description: "禁止没有 WHERE 的更新", hitCount: 31 },
  { id: 3, name: "SEL-L2-001", group: "查询", level: "警告", enabled: true, description: "禁止没有 WHERE 的查询", hitCount: 26 },
  { id: 4, name: "CTB-L2-006", group: "建表", level: "错误", enabled: true, description: "表名必须小写", hitCount: 18 },
  { id: 5, name: "INS-L3-002", group: "插入", level: "错误", enabled: true, description: "目标表必须已存在", hitCount: 15 },
  { id: 6, name: "MTG-L1-001", group: "杂项", level: "错误", enabled: true, description: "禁止 LOCK TABLE", hitCount: 9 },
  { id: 7, name: "CTB-L2-008", group: "建表", level: "警告", enabled: true, description: "表必须有注释", hitCount: 8 },
  { id: 8, name: "CVW-L3-002", group: "视图", level: "错误", enabled: true, description: "修改的视图必须已存在", hitCount: 6 },
  { id: 9, name: "CTL-L2-003", group: "列", level: "警告", enabled: false, description: "禁止使用可空列", hitCount: 4 },
  { id: 10, name: "MTB-L3-015", group: "表", level: "警告", enabled: true, description: "必须指定存储引擎", hitCount: 3 },
  { id: 11, name: "SEL-L3-003", group: "查询", level: "错误", enabled: true, description: "过滤列必须存在", hitCount: 2 },
  { id: 12, name: "MSC-L1-002", group: "杂项", level: "错误", enabled: true, description: "禁止 FLUSH 语句", hitCount: 1 },
]

export const violations: Violation[] = [
  { id: 1, rule: "DEL-L2-001", ticket: "清理历史订单", level: "警告", message: "删除语句没有 WHERE 条件，可能全表删除", createdAt: "2026-10-18 16:41" },
  { id: 2, rule: "INS-L3-002", ticket: "重建日报视图", level: "错误", message: "目标表 app_report.daily_report 不存在", createdAt: "2026-10-19 10:26" },
  { id: 3, rule: "CTB-L2-006", ticket: "重建日报视图", level: "错误", message: "表名 Daily_Report 必须全部小写", createdAt: "2026-10-19 10:26" },
  { id: 4, rule: "CTL-L2-003", ticket: "账号表新增状态列", level: "警告", message: "列 status 允许为空，建议设置 NOT NULL", createdAt: "2026-10-18 09:32" },
  { id: 5, rule: "MTB-L3-015", ticket: "归档表分区调整", level: "警告", message: "建表语句没有指定存储引擎", createdAt: "2026-10-16 14:55" },
  { id: 6, rule: "UPD-L2-001", ticket: "归档表分区调整", level: "错误", message: "更新语句没有 WHERE 条件", createdAt: "2026-10-16 14:56" },
  { id: 7, rule: "CVW-L3-002", ticket: "归档表分区调整", level: "错误", message: "视图 v_archive 不存在", createdAt: "2026-10-16 14:56" },
  { id: 8, rule: "MTG-L1-001", ticket: "移除无用索引", level: "错误", message: "禁止 LOCK TABLE 语句", createdAt: "2026-10-15 08:30" },
  { id: 9, rule: "SEL-L3-003", ticket: "支付表拆分字段", level: "错误", message: "过滤列 pay_type 不存在", createdAt: "2026-10-18 11:05" },
  { id: 10, rule: "CTB-L2-008", ticket: "用户标签表新增字段", level: "警告", message: "新建表缺少注释", createdAt: "2026-10-16 09:21" },
  { id: 11, rule: "MSC-L1-002", ticket: "灰度库同步结构", level: "错误", message: "禁止 FLUSH NO_WRITE_TO_BINLOG 语句", createdAt: "2026-10-17 10:09" },
  { id: 12, rule: "SEL-L2-001", ticket: "订单表字段注释补全", level: "警告", message: "查询语句没有 WHERE 条件", createdAt: "2026-10-15 13:47" },
]

export const auditLogs: AuditLog[] = [
  { id: 1, actor: "王五", operation: "审核通过工单 t-1023", target: "工单 1023", ip: "10.20.30.11", createdAt: "2026-10-19 15:20" },
  { id: 2, actor: "李四", operation: "提交工单 t-1024", target: "工单 1024", ip: "10.20.30.12", createdAt: "2026-10-20 09:12" },
  { id: 3, actor: "王五", operation: "驳回工单 t-1022", target: "工单 1022", ip: "10.20.30.11", createdAt: "2026-10-19 10:28" },
  { id: 4, actor: "张三", operation: "执行工单 t-1021 失败", target: "工单 1021", ip: "10.20.30.12", createdAt: "2026-10-18 16:46" },
  { id: 5, actor: "王五", operation: "预约工单 t-1020", target: "工单 1020", ip: "10.20.30.11", createdAt: "2026-10-18 11:10" },
  { id: 6, actor: "张三", operation: "新建群集 prod-archive-01", target: "群集 7", ip: "10.20.30.12", createdAt: "2026-10-17 08:40" },
  { id: 7, actor: "张三", operation: "修改规则 CTB-L2-008 参数", target: "规则 7", ip: "10.20.30.12", createdAt: "2026-10-16 20:11" },
  { id: 8, actor: "李四", operation: "关闭工单 t-1018", target: "工单 1018", ip: "10.20.30.12", createdAt: "2026-10-17 16:02" },
  { id: 9, actor: "王五", operation: "审核通过工单 t-1015", target: "工单 1015", ip: "10.20.30.11", createdAt: "2026-10-16 11:03" },
  { id: 10, actor: "李四", operation: "提交工单 t-1015", target: "工单 1015", ip: "10.20.30.12", createdAt: "2026-10-16 09:21" },
]

export const trend: TrendPoint[] = [
  { date: "10-14", submitted: 12, executed: 9, failed: 2 },
  { date: "10-15", submitted: 18, executed: 14, failed: 3 },
  { date: "10-16", submitted: 9, executed: 12, failed: 1 },
  { date: "10-17", submitted: 21, executed: 17, failed: 4 },
  { date: "10-18", submitted: 15, executed: 13, failed: 2 },
  { date: "10-19", submitted: 7, executed: 6, failed: 1 },
  { date: "10-20", submitted: 24, executed: 19, failed: 3 },
]

export const severity: SeverityPoint[] = [
  { severity: "致命", count: 3 },
  { severity: "高危", count: 8 },
  { severity: "中危", count: 14 },
  { severity: "低危", count: 21 },
]

export const ticketStatusMeta: Record<
  TicketStatus,
  { label: string; tone: "default" | "success" | "warning" | "danger" | "info" | "muted" }
> = {
  "waiting-vld": { label: "待系统审核", tone: "info" },
  "vld-failure": { label: "系统审核不通过", tone: "danger" },
  "vld-warning": { label: "系统审核告警", tone: "warning" },
  "waiting-mrv": { label: "待人工审核", tone: "info" },
  "mrv-failure": { label: "人工审核不通过", tone: "danger" },
  ltgm: { label: "待上线执行", tone: "default" },
  done: { label: "执行完成", tone: "success" },
  "exec-failure": { label: "执行失败", tone: "danger" },
  closed: { label: "已关闭", tone: "muted" },
}

export const periods = [7, 30, 180, 365] as const

export type Period = (typeof periods)[number]

export const periodLabels: Record<Period, string> = {
  7: "近 7 天",
  30: "近 30 天",
  180: "近 180 天",
  365: "近 365 天",
}

/**
 * 按时段生成趋势数据。用固定种子的线性同余，保证同一时段每次渲染一致，
 * 切换时段时曲线形态有变化但不随机跳动。
 */
export function trendFor(period: Period): TrendPoint[] {
  if (period === 7) {
    return trend
  }

  let seed = 20261020 + period
  const next = () => {
    seed = (seed * 1103515245 + 12345) % 2147483648
    return seed / 2147483648
  }

  const points: TrendPoint[] = []
  const step = period <= 30 ? 1 : period <= 180 ? 7 : 30
  for (let offset = period - step; offset >= 0; offset -= step) {
    const submitted = Math.round(6 + next() * 20)
    const executed = Math.round(submitted * (0.55 + next() * 0.35))
    const failed = Math.max(0, Math.round(submitted - executed) + Math.round(next() * 2))
    const date = new Date(2026, 9, 20)
    date.setDate(date.getDate() - offset)
    points.push({
      date: `${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")}`,
      submitted,
      executed,
      failed,
    })
  }
  return points
}

/** 接入群集数量的历史，用于第四张指标卡的趋势。 */
export const clusterTrend = [4, 5, 5, 6, 7, 7, 8]
