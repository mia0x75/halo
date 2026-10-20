import {
  ClusterStatus,
  TicketStatus,
  clusterStatusByValue,
  ticketStatusByValue,
  type Ticket,
  type UserRef,
} from "./types"

/**
 * Thin GraphQL transport for the Go backend.
 *
 * The backend serves GraphQL at /api/query; the Vite dev server proxies that
 * path (see vite.config.ts). No cache layer is installed on purpose — the
 * decision record keeps server state handling out of scope until a concrete
 * need appears.
 */

const ENDPOINT = "/api/query"

interface GraphQLResponse<TData> {
  data?: TData
  errors?: { message: string }[]
}

export class ApiError extends Error {
  constructor(message: string) {
    super(message)
    this.name = "ApiError"
  }
}

export async function graphqlRequest<TData, TVariables extends Record<string, unknown> = Record<string, never>>(
  query: string,
  variables?: TVariables,
): Promise<TData> {
  const response = await fetch(ENDPOINT, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ query, variables }),
  })

  if (!response.ok) {
    throw new ApiError(`GraphQL 请求失败：HTTP ${response.status}`)
  }

  const payload = (await response.json()) as GraphQLResponse<TData>
  if (payload.errors?.length) {
    throw new ApiError(payload.errors.map((e) => e.message).join("; "))
  }
  if (!payload.data) {
    throw new ApiError("GraphQL 响应缺少 data 字段")
  }
  return payload.data
}

interface RawUserRef {
  UUID: string
  Email: string
  Name: string | null
}

interface RawTicket {
  UUID: string
  Database: string
  Subject: string
  Content: string
  Status: number
  CreateAt: number
  UpdateAt: number | null
  Cluster: {
    UUID: string
    Host: string
    Alias: string
    IP: string
    Port: number
    User: string
    Status: number
    CreateAt: number
  }
  User: RawUserRef
  Reviewer: RawUserRef | null
}

const TICKETS_QUERY = `
  query Tickets($first: Int) {
    tickets(first: $first) {
      totalCount
      edges {
        node {
          UUID
          Database
          Subject
          Content
          Status
          CreateAt
          UpdateAt
          Cluster { UUID Host Alias IP Port User Status CreateAt }
          User { UUID Email Name }
          Reviewer { UUID Email Name }
        }
      }
    }
  }
`

function mapUser(raw: RawUserRef): UserRef {
  return { uuid: raw.UUID, email: raw.Email, name: raw.Name }
}

function mapTicket(raw: RawTicket): Ticket {
  return {
    uuid: raw.UUID,
    cluster: {
      uuid: raw.Cluster.UUID,
      host: raw.Cluster.Host,
      alias: raw.Cluster.Alias,
      ip: raw.Cluster.IP,
      port: raw.Cluster.Port,
      user: raw.Cluster.User,
      status: clusterStatusByValue[raw.Cluster.Status] ?? ClusterStatus.Normal,
      createAt: raw.Cluster.CreateAt,
    },
    database: raw.Database,
    subject: raw.Subject,
    content: raw.Content,
    status: ticketStatusByValue[raw.Status] ?? TicketStatus.WaitingForVld,
    user: mapUser(raw.User),
    reviewer: raw.Reviewer ? mapUser(raw.Reviewer) : null,
    createAt: raw.CreateAt,
    updateAt: raw.UpdateAt,
  }
}

export async function listTickets(first = 20): Promise<{ totalCount: number; tickets: Ticket[] }> {
  const data = await graphqlRequest<
    { tickets: { totalCount: number; edges: { node: RawTicket }[] } },
    { first: number }
  >(TICKETS_QUERY, { first })

  return {
    totalCount: data.tickets.totalCount,
    tickets: data.tickets.edges.map((edge) => mapTicket(edge.node)),
  }
}
