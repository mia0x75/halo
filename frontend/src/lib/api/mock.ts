import {
  ClusterStatus,
  TicketStatus,
  type AuditLog,
  type Cluster,
  type DashboardSummary,
  type SeverityBreakdown,
  type Ticket,
  type TicketTrendPoint,
  type UserRef,
} from "./types"

/**
 * Typed fixtures for the mock-first phase. Every shape mirrors lib/api/types.ts,
 * which mirrors the GraphQL contract, so swapping this module for the real
 * client is a one-line change at the call site.
 */

const NOW = Math.floor(Date.now() / 1000)

const clusters: Cluster[] = [
  {
    uuid: "c1a0f5e2-0f4a-4d1f-9f0b-2a4d0a1f0001",
    host: "db-prod-01.internal",
    alias: "prod-01",
    ip: "10.20.30.11",
    port: 3306,
    user: "halo_ops",
    status: ClusterStatus.Normal,
    createAt: NOW - 86400 * 180,
  },
  {
    uuid: "c1a0f5e2-0f4a-4d1f-9f0b-2a4d0a1f0002",
    host: "db-prod-02.internal",
    alias: "prod-02",
    ip: "10.20.30.12",
    port: 3306,
    user: "halo_ops",
    status: ClusterStatus.Normal,
    createAt: NOW - 86400 * 180,
  },
  {
    uuid: "c1a0f5e2-0f4a-4d1f-9f0b-2a4d0a1f0003",
    host: "db-stg-01.internal",
    alias: "staging-01",
    ip: "10.20.31.21",
    port: 3306,
    user: "halo_ops",
    status: ClusterStatus.Normal,
    createAt: NOW - 86400 * 90,
  },
]

const users: UserRef[] = [
  { uuid: "u0000000-0000-4000-8000-000000000001", email: "submitter@example.com", name: "提交人" },
  { uuid: "u0000000-0000-4000-8000-000000000002", email: "reviewer@example.com", name: "审核人" },
  { uuid: "u0000000-0000-4000-8000-000000000003", email: "admin@example.com", name: "管理员" },
]

const tickets: Ticket[] = [
  {
    uuid: "t0000000-0000-4000-8000-000000000001",
    cluster: clusters[0]!,
    database: "app_order",
    subject: "订单表新增归档字段",
    content: "ALTER TABLE orders ADD COLUMN archived_at DATETIME NULL;",
    status: TicketStatus.WaitingForMrv,
    user: users[0]!,
    reviewer: users[1]!,
    createAt: NOW - 3600 * 2,
    updateAt: NOW - 3600,
  },
  {
    uuid: "t0000000-0000-4000-8000-000000000002",
    cluster: clusters[1]!,
    database: "app_user",
    subject: "用户表索引优化",
    content: "CREATE INDEX idx_users_phone ON users (phone);",
    status: TicketStatus.Done,
    user: users[0]!,
    reviewer: users[1]!,
    createAt: NOW - 86400,
    updateAt: NOW - 86400 + 600,
  },
  {
    uuid: "t0000000-0000-4000-8000-000000000003",
    cluster: clusters[2]!,
    database: "app_report",
    subject: "报表视图重建",
    content: "CREATE OR REPLACE VIEW v_daily_report AS SELECT day, sum(amount) FROM payment GROUP BY day;",
    status: TicketStatus.VldFailure,
    user: users[2]!,
    reviewer: null,
    createAt: NOW - 86400 * 3,
    updateAt: NOW - 86400 * 3 + 120,
  },
  {
    uuid: "t0000000-0000-4000-8000-000000000004",
    cluster: clusters[0]!,
    database: "app_order",
    subject: "清理历史数据",
    content: "DELETE FROM orders WHERE created_at < '2024-01-01';",
    status: TicketStatus.ExecFailure,
    user: users[0]!,
    reviewer: users[1]!,
    createAt: NOW - 86400 * 5,
    updateAt: NOW - 86400 * 5 + 300,
  },
]

const logs: AuditLog[] = [
  {
    uuid: "l0000000-0000-4000-8000-000000000001",
    logId: 101,
    user: users[1]!,
    operation: "审核通过工单 t0000000-0000-4000-8000-000000000001",
    createAt: NOW - 3600,
  },
  {
    uuid: "l0000000-0000-4000-8000-000000000002",
    logId: 102,
    user: users[0]!,
    operation: "提交工单 t0000000-0000-4000-8000-000000000002",
    createAt: NOW - 86400,
  },
  {
    uuid: "l0000000-0000-4000-8000-000000000003",
    logId: 103,
    user: users[2]!,
    operation: "执行工单 t0000000-0000-4000-8000-000000000004 失败",
    createAt: NOW - 86400 * 5 + 300,
  },
]

const trend: TicketTrendPoint[] = [
  { date: "10-14", submitted: 12, executed: 9 },
  { date: "10-15", submitted: 18, executed: 14 },
  { date: "10-16", submitted: 9, executed: 12 },
  { date: "10-17", submitted: 21, executed: 17 },
  { date: "10-18", submitted: 15, executed: 13 },
  { date: "10-19", submitted: 7, executed: 6 },
  { date: "10-20", submitted: 24, executed: 19 },
]

const severity: SeverityBreakdown[] = [
  { severity: "critical", count: 3 },
  { severity: "high", count: 8 },
  { severity: "medium", count: 14 },
  { severity: "low", count: 21 },
]

const pendingStatuses: TicketStatus[] = [
  TicketStatus.WaitingForVld,
  TicketStatus.WaitingForMrv,
]

const failedStatuses: TicketStatus[] = [TicketStatus.VldFailure, TicketStatus.ExecFailure]

export const mockDashboardSummary: DashboardSummary = {
  pendingTickets: tickets.filter((t) => pendingStatuses.includes(t.status)).length,
  executedToday: trend[trend.length - 1]!.executed,
  failedTickets: tickets.filter((t) => failedStatuses.includes(t.status)).length,
  clusters: clusters.length,
  trend,
  severity,
  recentTickets: tickets,
  recentLogs: logs,
}

export const mockTickets: Ticket[] = tickets
