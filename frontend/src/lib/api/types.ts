/**
 * Client-side mirror of the GraphQL contract in
 * internal/transport/graphql/gqlapi/schema.graphql.
 *
 * The server speaks PascalCase (gqlgen default); the client uses camelCase
 * and converts at the network boundary (see lib/api/client.ts). Enum values
 * keep the server integer encoding so no mapping table is needed.
 */

export const TicketStatus = {
  WaitingForVld: 1,
  VldFailure: 2,
  VldWarning: 3,
  WaitingForMrv: 4,
  MrvFailure: 5,
  Ltgm: 6,
  Done: 7,
  ExecFailure: 8,
  Closed: 9,
} as const

export type TicketStatus = (typeof TicketStatus)[keyof typeof TicketStatus]

export const ClusterStatus = {
  Normal: 1,
  Blocked: 2,
} as const

export type ClusterStatus = (typeof ClusterStatus)[keyof typeof ClusterStatus]

export const UserStatus = {
  Normal: 1,
  Blocked: 2,
  Pending: 3,
} as const

export type UserStatus = (typeof UserStatus)[keyof typeof UserStatus]

/** Numeric wire values mapped to the enum members, for decoding responses. */
export const ticketStatusByValue: Record<number, TicketStatus> = {
  1: TicketStatus.WaitingForVld,
  2: TicketStatus.VldFailure,
  3: TicketStatus.VldWarning,
  4: TicketStatus.WaitingForMrv,
  5: TicketStatus.MrvFailure,
  6: TicketStatus.Ltgm,
  7: TicketStatus.Done,
  8: TicketStatus.ExecFailure,
  9: TicketStatus.Closed,
}

export const clusterStatusByValue: Record<number, ClusterStatus> = {
  1: ClusterStatus.Normal,
  2: ClusterStatus.Blocked,
}

export interface Cluster {
  uuid: string
  host: string
  alias: string
  ip: string
  port: number
  user: string
  status: ClusterStatus
  createAt: number
}

export interface UserRef {
  uuid: string
  email: string
  name: string | null
}

export interface Ticket {
  uuid: string
  cluster: Cluster
  database: string
  subject: string
  content: string
  status: TicketStatus
  user: UserRef
  reviewer: UserRef | null
  createAt: number
  updateAt: number | null
}

export interface AuditLog {
  uuid: string
  logId: number
  user: UserRef
  operation: string
  createAt: number
}

export interface TicketTrendPoint {
  date: string
  submitted: number
  executed: number
}

export interface SeverityBreakdown {
  severity: "critical" | "high" | "medium" | "low"
  count: number
}

export interface DashboardSummary {
  pendingTickets: number
  executedToday: number
  failedTickets: number
  clusters: number
  trend: TicketTrendPoint[]
  severity: SeverityBreakdown[]
  recentTickets: Ticket[]
  recentLogs: AuditLog[]
}
