/**
 * 页面模型。接口节点（camelCase 之外的 GraphQL 契约）在 mappers 里转成这些结构，
 * 组件只认这里定义的类型。
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

export interface TrendPoint {
  date: string
  submitted: number
  executed: number
  failed: number
}

export interface TopRulePoint {
  rule: string
  count: number
}

export interface SeverityPoint {
  severity: string
  count: number
}
