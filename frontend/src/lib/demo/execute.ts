/**
 * 演示模式的本地执行器。把前端发出的 GraphQL 文档按操作名分发到内存数据，
 * 返回与真实接口同形的结果，因此页面代码无需区分真假接口。
 */

import type {
  ClusterNode,
  CommentNode,
  GroupNode,
  StatementNode,
  TicketNode,
  UserNode,
} from "@/lib/api/queries"
import { buildDemoState, demoUser, demoUUID, metricSeeds, today, type DemoState, type DemoTicket } from "@/lib/demo/data"

const HOUR = 3600

let state: DemoState = buildDemoState()

export function resetDemoData(): void {
  state = buildDemoState()
}

type Variables = Record<string, unknown>

function num(value: unknown, fallback: number): number {
  return typeof value === "number" && Number.isFinite(value) ? value : fallback
}

function text(value: unknown): string | undefined {
  return typeof value === "string" ? value : undefined
}

function flag(value: unknown): boolean {
  return value === true
}

function record(value: unknown): Variables {
  return typeof value === "object" && value !== null ? (value as Variables) : {}
}

function encodeCursor(index: number): string {
  return btoa(`cursor${index}`)
}

function decodeCursor(cursor: string | undefined): number | null {
  if (!cursor) return null
  try {
    const matched = /^cursor(\d+)$/.exec(atob(cursor))
    return matched ? Number(matched[1]) : null
  } catch {
    return null
  }
}

interface PageInfo {
  hasNextPage: boolean
  hasPreviousPage: boolean
  startCursor: string
  endCursor: string
}

interface Connection<T> {
  totalCount: number
  edges: { cursor: string; node: T }[]
  pageInfo: PageInfo
}

/** 与后端一致：游标是列表下标的不透明编码，after 指向已返回的最后一条。 */
function paginate<T>(items: T[], variables: Variables): Connection<T> {
  const size = Math.max(1, num(variables.first, 20))
  const after = decodeCursor(text(variables.after))
  const start = after === null ? 0 : after + 1
  const slice = items.slice(start, start + size)
  const edges = slice.map((node, offset) => ({ cursor: encodeCursor(start + offset), node }))
  return {
    totalCount: items.length,
    edges,
    pageInfo: {
      hasNextPage: start + slice.length < items.length,
      hasPreviousPage: start > 0,
      startCursor: edges[0]?.cursor ?? "",
      endCursor: edges[edges.length - 1]?.cursor ?? "",
    },
  }
}

function ticketNode(ticket: DemoTicket): TicketNode {
  return {
    UUID: ticket.UUID,
    TicketID: ticket.TicketID,
    Database: ticket.Database,
    Subject: ticket.Subject,
    Status: ticket.Status,
    CreateAt: ticket.CreateAt,
    User: ticket.User,
    Reviewer: ticket.Reviewer,
    Cluster: ticket.Cluster,
  }
}

function ticketById(id: string | undefined): DemoTicket | undefined {
  if (!id) return undefined
  return (
    state.tickets.find((item) => item.UUID === id) ??
    state.tickets.find((item) => String(item.TicketID) === id)
  )
}

function newTicketID(): number {
  return state.tickets.reduce((max, item) => Math.max(max, item.TicketID), 0) + 1
}

function statementNodes(): StatementNode[] {
  return state.tickets
    .flatMap((ticket) => ticket.Statements)
    .sort((a, b) => b.CreateAt - a.CreateAt || b.TicketID - a.TicketID)
}

function buildRuleHitLevels(): { RuleName: string | null; Level: number; Count: number }[] {
  const buckets = new Map<string, { RuleName: string | null; Level: number; Count: number }>()
  for (const hit of state.hits) {
    const key = `${hit.RuleName}:${hit.Level}`
    const row = buckets.get(key)
    if (row) {
      row.Count += 1
    } else {
      buckets.set(key, { RuleName: hit.RuleName, Level: hit.Level, Count: 1 })
    }
  }
  return [...buckets.values()]
}

function buildRuleHitRules(first: number): { RuleName: string | null; Level: number; Count: number }[] {
  return buildRuleHitLevels()
    .sort((a, b) => b.Count - a.Count)
    .slice(0, Math.max(1, first))
}

function createTicket(input: Variables): DemoTicket {
  const id = newTicketID()
  const cluster = state.clusters.find((item) => item.UUID === text(input.ClusterUUID)) ?? state.clusters[0]!
  const reviewerUUID = text(input.ReviewerUUID)
  const reviewer = state.users.find((item) => item.UUID === reviewerUUID)
  const content = text(input.Content) ?? ""
  const createAt = Math.floor(Date.now() / 1000)
  const statements: StatementNode[] = content
    .split(";")
    .filter((chunk) => chunk.trim().length > 0)
    .map((chunk, index) => ({
      UUID: demoUUID(0x60000000, id * 10 + index),
      TicketID: id,
      Sequence: index + 1,
      Content: `${chunk.trim()};`,
      TypeDesc: "UNKNOWN",
      Status: 1,
      Report: "[]",
      RowsAffected: null,
      CreateAt: createAt,
      UpdateAt: null,
    }))

  const ticket: DemoTicket = {
    UUID: demoUUID(0x10000000, id),
    TicketID: id,
    Database: text(input.Database) ?? "",
    Subject: text(input.Subject) ?? "未命名工单",
    Status: 1,
    CreateAt: createAt,
    User: { UUID: demoUser.UUID, Name: demoUser.Name },
    Reviewer: reviewer
      ? { UUID: reviewer.UUID, Name: reviewer.Name }
      : { UUID: demoUUID(0x30000000, 3), Name: "李四" },
    Cluster: {
      UUID: cluster.UUID,
      Alias: cluster.Alias,
      Host: cluster.Host,
      IP: cluster.IP,
      Port: cluster.Port,
      Type: cluster.Type,
      Group: cluster.Group,
    },
    Content: content,
    UpdateAt: null,
    Statements: statements,
    Comments: [],
  }
  state.tickets.unshift(ticket)
  return ticket
}

function updateTicket(input: Variables): DemoTicket | undefined {
  const ticket = ticketById(text(input.TicketUUID) ?? text(input.UUID) ?? text(input.id))
  if (!ticket) return undefined
  if (typeof input.Subject === "string") ticket.Subject = input.Subject
  if (typeof input.Content === "string") ticket.Content = input.Content
  if (typeof input.Database === "string") ticket.Database = input.Database
  if (typeof input.Status === "number") ticket.Status = input.Status
  ticket.UpdateAt = Math.floor(Date.now() / 1000)
  return ticket
}

function clusterInput(input: Variables, existing?: ClusterNode): ClusterNode {
  // GroupUUID 传空串表示「未分组」，与「不传该字段」要区分开。
  const groupGiven = "GroupUUID" in input
  const groupUUID = text(input.GroupUUID)
  const group = state.groups.find((item) => item.UUID === groupUUID)
  return {
    UUID: existing?.UUID ?? `cluster-${state.clusters.length + 1}`,
    Alias: text(input.Alias) ?? existing?.Alias ?? "cluster",
    Host: text(input.Host) ?? existing?.Host ?? "",
    IP: text(input.IP) ?? existing?.IP ?? "",
    Port: num(input.Port, existing?.Port ?? 3306),
    User: text(input.User) ?? existing?.User ?? "root",
    Status: num(input.Status, existing?.Status ?? 1),
    Type: text(input.Type) ?? existing?.Type ?? "mysql",
    Group: groupGiven ? (group?.Name ?? null) : (existing?.Group ?? null),
    GroupUUID: groupGiven ? (group?.UUID ?? null) : (existing?.GroupUUID ?? null),
    CreateAt: existing?.CreateAt ?? Math.floor(Date.now() / 1000),
  }
}

/**
 * 某个用户的群集权限：直接授权，加上角色被授予的分组里的群集。
 * 返回每个群集的生效状态与来源，供管理面区分两层权限。
 */
function demoClusterGrants(userUUID: string) {
  const user = state.users.find((item) => item.UUID === userUUID)
  const granted = new Set((user?.GrantedClusters ?? []).map((cluster) => cluster.UUID))

  const sourcesByGroup = new Map<string, string[]>()
  for (const role of state.roles) {
    const owned = (user?.Roles ?? []).some((item) => item.UUID === role.UUID)
    if (!owned) continue
    for (const group of role.Groups ?? []) {
      const list = sourcesByGroup.get(group.UUID) ?? []
      list.push(`${role.Name} / ${group.Name}`)
      sourcesByGroup.set(group.UUID, list)
    }
  }

  return state.clusters.map((cluster) => {
    const sources = cluster.GroupUUID ? (sourcesByGroup.get(cluster.GroupUUID) ?? []) : []
    return {
      Cluster: { UUID: cluster.UUID, Alias: cluster.Alias, Group: cluster.Group },
      Direct: granted.has(cluster.UUID),
      Inherited: sources.length > 0,
      Sources: sources,
    }
  })
}

/** 当前用户可见的群集：直接授权，加上角色被授予的分组里的群集。 */
function demoUserClusters(): ClusterNode[] {
  const grants = demoClusterGrants(demoUser.UUID)
  return state.clusters.filter((cluster) => {
    const grant = grants.find((item) => item.Cluster.UUID === cluster.UUID)
    return Boolean(grant && (grant.Direct || grant.Inherited))
  })
}

export function demoExecute(query: string, variables?: Variables): Record<string, unknown> {
  const source = variables ?? {}
  const matched = /(?:query|mutation)\s+(\w+)/.exec(query)
  const operation = matched?.[1] ?? ""

  switch (operation) {
    case "Me": {
      const user = state.users.find((item) => item.UUID === demoUser.UUID)
      return {
        me: {
          ...demoUser,
          GrantedClusters: user?.GrantedClusters ?? [],
          ClusterGrants: demoClusterGrants(demoUser.UUID),
        },
      }
    }
    case "Metrics":
    case "MetricSummary": {
      const names = Array.isArray(source.Names) ? source.Names.filter((item): item is string => typeof item === "string") : []
      const scoped = names.length > 0 ? metricSeeds.filter((row) => names.includes(row.Name)) : metricSeeds
      // 页面用 find() 取最新值，所以按时间倒序返回。
      return {
        metrics: [...scoped].reverse().map((row, index) => ({
          UUID: `metric-${index}`,
          Name: row.Name,
          Scope: 1,
          ScopeID: null,
          Type: 0,
          Bucket: row.Bucket,
          Value: row.Value,
          CreateAt: today,
        })),
      }
    }
    case "RuleHitLevels": {
      return { ruleHitLevels: buildRuleHitLevels() }
    }
    case "RuleHitRules": {
      return { ruleHitRules: buildRuleHitRules(num(source.first, 10)) }
    }
    case "Messages": {
      const unreadOnly = flag(source.unreadOnly)
      const category = text(source.category) ?? ""
      const keyword = (text(source.keyword) ?? "").trim().toLowerCase()
      const matched = state.messages.filter((item) => {
        if (unreadOnly && item.ReadAt !== null) return false
        if (category && item.Category !== category) return false
        if (
          keyword &&
          !`${item.Title} ${item.Content}`.toLowerCase().includes(keyword)
        ) {
          return false
        }
        return true
      })
      const connection = paginate(matched, source)
      return { messages: connection }
    }
    case "UnreadMessageCount": {
      return { unreadMessageCount: state.messages.filter((item) => item.ReadAt === null).length }
    }
    case "MarkMessageRead": {
      const target = state.messages.find((item) => item.UUID === text(source.uuid))
      if (target && target.ReadAt === null) target.ReadAt = Math.floor(Date.now() / 1000)
      return { markMessageRead: true }
    }
    case "RemoveMessage": {
      const uuid = text(source.uuid)
      state.messages = state.messages.filter((item) => item.UUID !== uuid)
      return { removeMessage: true }
    }
    case "MarkAllMessagesRead": {
      const now = Math.floor(Date.now() / 1000)
      let rows = 0
      for (const item of state.messages) {
        if (item.ReadAt === null) {
          item.ReadAt = now
          rows++
        }
      }
      return { markAllMessagesRead: rows }
    }
    case "Groups": {
      return { groups: state.groups }
    }
    case "CreateGroup": {
      const input = record(source.input)
      const group: GroupNode = {
        UUID: `group-${state.groups.length + 1}`,
        Name: text(input.Name) ?? "新分组",
        Description: text(input.Description) ?? "",
        Status: num(input.Status, 1),
        CreateAt: Math.floor(Date.now() / 1000),
      }
      state.groups.push(group)
      return { createGroup: group }
    }
    case "UpdateGroup": {
      const input = record(source.input)
      const existing = state.groups.find((item) => item.UUID === text(input.GroupUUID))
      if (existing) {
        existing.Name = text(input.Name) ?? existing.Name
        existing.Description = text(input.Description) ?? existing.Description
        existing.Status = num(input.Status, existing.Status)
      }
      return { updateGroup: existing ?? null }
    }
    case "RemoveGroup": {
      const uuid = text(source.id)
      state.groups = state.groups.filter((item) => item.UUID !== uuid)
      for (const role of state.roles) {
        role.Groups = (role.Groups ?? []).filter((group) => group.UUID !== uuid)
      }
      return { removeGroup: true }
    }
    case "GrantGroups": {
      const input = record(source.input)
      const role = state.roles.find((item) => item.UUID === text(input.RoleUUID))
      const uuids = Array.isArray(input.GroupUUIDs)
        ? input.GroupUUIDs.filter((item): item is string => typeof item === "string")
        : []
      if (role) {
        role.Groups = state.groups.filter((group) => uuids.includes(group.UUID))
      }
      return { grantGroups: true }
    }
    case "GrantClusters": {
      const input = record(source.input)
      const user = state.users.find((item) => item.UUID === text(input.UserUUID))
      const uuids = Array.isArray(input.ClusterUUIDs)
        ? input.ClusterUUIDs.filter((item): item is string => typeof item === "string")
        : []
      if (user) {
        user.GrantedClusters = state.clusters
          .filter((cluster) => uuids.includes(cluster.UUID))
          .map((cluster) => ({ UUID: cluster.UUID, Alias: cluster.Alias }))
      }
      return { grantClusters: true }
    }
    case "Tickets": {
      // 0 表示「全部状态」，和后端过滤语义一致。
      const status = typeof source.status === "number" && source.status > 0 ? source.status : null
      const keyword = (text(source.keyword) ?? "").trim().toLowerCase()
      const rows = [...state.tickets]
        .sort((a, b) => b.TicketID - a.TicketID)
        .filter((item) => (status === null || item.Status === status) && (keyword === "" || `${item.Subject} ${item.Database}`.toLowerCase().includes(keyword)))
        .map(ticketNode)
      return { tickets: paginate(rows, source) }
    }
    case "Login": {
      return { login: { Me: { UUID: demoUser.UUID, Name: demoUser.Name, Email: demoUser.Email }, Token: "demo" } }
    }
    case "MyTickets": {
      const rows = state.tickets
        .filter((item) => item.User?.UUID === demoUser.UUID)
        .sort((a, b) => b.TicketID - a.TicketID)
        .map(ticketNode)
      return { me: { ...demoUser, MyTickets: paginate(rows, source) } }
    }
    case "PendingReviews": {
      const rows = state.tickets
        .filter((item) => item.Status === 4)
        .sort((a, b) => b.TicketID - a.TicketID)
        .map(ticketNode)
      return { me: { ...demoUser, PendingReviews: paginate(rows, source) } }
    }
    case "Clusters": {
      return { clusters: paginate(state.clusters, source) }
    }
    case "CreateCluster": {
      const cluster = clusterInput(record(source.input))
      state.clusters.unshift(cluster)
      return { createCluster: cluster }
    }
    case "UpdateCluster": {
      const input = record(source.input)
      const existing = state.clusters.find((item) => item.UUID === text(input.ClusterUUID))
      const cluster = clusterInput(input, existing)
      if (existing) Object.assign(existing, cluster)
      return { updateCluster: cluster }
    }
    case "RemoveCluster": {
      const id = text(source.id)
      state.clusters = state.clusters.filter((item) => item.UUID !== id)
      return { removeCluster: true }
    }
    case "CreateTicket": {
      const ticket = createTicket(record(source.input))
      return { createTicket: ticketNode(ticket) }
    }
    case "ReviewTicket": {
      const input = record(source.input)
      const ticket = ticketById(text(input.TicketUUID))
      if (!ticket) throw new Error("演示环境：工单不存在")
      if (ticket.Status !== 4) throw new Error("演示环境：工单不在人工审核阶段")
      const outcome = text(input.Outcome)
      if (outcome === "APPROVED") ticket.Status = 6
      else if (outcome === "REJECTED") ticket.Status = 5
      else throw new Error("演示环境：审核结论无效")
      ticket.UpdateAt = Math.floor(Date.now() / 1000)
      return { reviewTicket: { UUID: ticket.UUID, Status: ticket.Status } }
    }
    case "CreateComment": {
      const input = record(source.input)
      const ticket = ticketById(text(input.TicketUUID))
      if (!ticket) throw new Error("演示环境：工单不存在")
      const content = text(input.Content) ?? ""
      const comment: CommentNode = {
        UUID: `comment-${ticket.TicketID}-${ticket.Comments.length + 1}`,
        Content: content,
        CreateAt: Math.floor(Date.now() / 1000),
        User: { UUID: demoUser.UUID, Name: demoUser.Name },
      }
      ticket.Comments.push(comment)
      return { createComment: comment }
    }
    case "UpdateTicket": {
      const ticket = updateTicket(record(source.input))
      return { updateTicket: ticket ? ticketNode(ticket) : null }
    }
    case "Ticket": {
      const ticket = ticketById(text(source.id))
      if (!ticket) return { ticket: null }
      return {
        ticket: {
          ...ticket,
          Statements: paginate(ticket.Statements, { first: 20 }),
          Comments: paginate(ticket.Comments, { first: 20 }),
        },
      }
    }
    case "MeClusters": {
      return {
        me: {
          ...demoUser,
          GrantedClusters: demoUserClusters(),
          Clusters: paginate(demoUserClusters(), { ...source, after: undefined }),
        },
      }
    }
    case "MeReviewers": {
      return { me: { ...demoUser, Reviewers: demoUser.Reviewers ?? [] } }
    }
    case "TestCluster": {
      const input = record(source.input)
      const type = text(input.Type) ?? "mysql"
      const databases = Object.values(state.databases).flat()
      return {
        testCluster: {
          Success: true,
          Message: "连接成功",
          Version: type === "postgres" ? "PostgreSQL 17.2" : "MySQL 8.0.36",
          LatencyMs: 12,
          Databases: [...new Set(databases.map((item) => item.Name))].slice(0, 6),
        },
      }
    }
    case "Rules": {
      return { rules: state.rules }
    }
    case "RuleHits": {
      const first = num(source.first, 20)
      const rows = state.hits
        .slice()
        .sort((a, b) => b.CreateAt - a.CreateAt)
        .slice(0, first)
      return { ruleHits: rows }
    }
    case "Logs": {
      const rows = state.logs.slice().sort((a, b) => b.LogID - a.LogID)
      return { logs: paginate(rows, source) }
    }
    case "Crons": {
      return { crons: paginate(state.crons, source) }
    }
    case "PendingTasks":
    case "Tasks": {
      return { tasks: state.tasks }
    }
    case "TaskExecutions": {
      const status = typeof source.status === "number" && source.status > 0 ? source.status : null
      const first = num(source.first, 20)
      const rows = state.executions
        .filter((item) => status === null || item.Status === status)
        .slice(0, first)
      const summary = new Map<number, number>()
      for (const item of state.executions) {
        summary.set(item.Status, (summary.get(item.Status) ?? 0) + 1)
      }
      return {
        taskExecutions: rows,
        taskExecutionSummary: [...summary.entries()].map(([status_, count]) => ({ Status: status_, Count: count })),
      }
    }
    case "CreateCron": {
      const input = record(source.input)
      const cron = {
        UUID: `cron-${state.crons.length + 1}`,
        Name: text(input.Name) ?? "新作业",
        Cmd: text(input.Job) ?? "METRICS_ROLLUP",
        Params: "{}",
        Interval: text(input.Interval) ?? "24h",
        Duration: null,
        LastRun: null,
        NextRun: new Date((today + HOUR) * 1000).toISOString(),
        Recurrent: 1,
        Status: "PENDING",
        CreateAt: Math.floor(Date.now() / 1000),
      }
      state.crons.unshift(cron)
      return { createCron: cron }
    }
    case "RemoveCron": {
      const id = text(source.id)
      state.crons = state.crons.filter((item) => item.UUID !== id)
      return { removeCron: true }
    }
    case "Databases": {
      const clusterUUID = text(source.clusterUUID) ?? ""
      return { databases: state.databases[clusterUUID] ?? [] }
    }
    case "Users": {
      const rows = state.users.map((user) => ({
        ...user,
        ClusterGrants: demoClusterGrants(user.UUID),
      }))
      return { users: paginate(rows, source) }
    }
    case "CreateUser": {
      const input = record(source.input)
      const roleUUIDs = Array.isArray(input.RoleUUIDs)
        ? input.RoleUUIDs.filter((item): item is string => typeof item === "string")
        : []
      const clusterUUIDs = Array.isArray(input.ClusterUUIDs)
        ? input.ClusterUUIDs.filter((item): item is string => typeof item === "string")
        : []
      const reviewerUUIDs = Array.isArray(input.ReviewerUUIDs)
        ? input.ReviewerUUIDs.filter((item): item is string => typeof item === "string")
        : []
      const user: UserNode = {
        UUID: `user-${state.users.length + 1}`,
        Email: text(input.Email) ?? "user@example.com",
        Name: text(input.Name) ?? "",
        Status: num(input.Status, 1),
        Roles: state.roles
          .filter((role) => roleUUIDs.includes(role.UUID))
          .map((role) => ({ UUID: role.UUID, Name: role.Name })),
        Reviewers: state.users
          .filter((item) => reviewerUUIDs.includes(item.UUID))
          .map((item) => ({ UUID: item.UUID, Name: item.Name, Email: item.Email })),
        GrantedClusters: state.clusters
          .filter((cluster) => clusterUUIDs.includes(cluster.UUID))
          .map((cluster) => ({ UUID: cluster.UUID, Alias: cluster.Alias })),
      }
      state.users.push(user)
      return { createUser: user }
    }
    case "UpdateUser": {
      const input = record(source.input)
      const existing = state.users.find((item) => item.UUID === text(input.UserUUID))
      if (existing) {
        existing.Email = text(input.Email) ?? existing.Email
        existing.Name = text(input.Name) ?? existing.Name
        existing.Status = num(input.Status, existing.Status)
      }
      return { updateUser: existing ?? null }
    }
    case "PatchUserStatus": {
      const input = record(source.input)
      const existing = state.users.find((item) => item.UUID === text(input.UserUUID))
      if (existing) existing.Status = num(input.Status, existing.Status)
      return { patchUserStatus: true }
    }
    case "GrantRoles": {
      const input = record(source.input)
      const user = state.users.find((item) => item.UUID === text(input.UserUUID))
      const roleUUIDs = Array.isArray(input.RoleUUIDs)
        ? input.RoleUUIDs.filter((item): item is string => typeof item === "string")
        : []
      if (user) {
        user.Roles = state.roles
          .filter((role) => roleUUIDs.includes(role.UUID))
          .map((role) => ({ UUID: role.UUID, Name: role.Name }))
      }
      return { grantRoles: true }
    }
    case "GrantReviewers": {
      const input = record(source.input)
      const user = state.users.find((item) => item.UUID === text(input.UserUUID))
      const reviewerUUIDs = Array.isArray(input.ReviewerUUIDs)
        ? input.ReviewerUUIDs.filter((item): item is string => typeof item === "string")
        : []
      if (user) {
        user.Reviewers = state.users
          .filter((item) => reviewerUUIDs.includes(item.UUID))
          .map((item) => ({ UUID: item.UUID, Name: item.Name, Email: item.Email }))
      }
      return { grantReviewers: true }
    }
    case "Roles": {
      return { roles: state.roles }
    }
    case "Options": {
      return { options: state.options }
    }
    case "SystemInfo": {
      return {
        systemInfo: {
          Version: "0.1.0",
          Git: "demo",
          Compile: "2026-10-21",
          Branch: "main",
          Database: "mysql",
          MailEnabled: false,
          ClusterCount: state.clusters.length,
          UserCount: state.users.length,
        },
      }
    }
    case "PatchOption": {
      const input = record(source.input)
      const option = state.options.find((item) => item.UUID === text(input.OptionUUID))
      if (option) option.Parameter = text(input.Parameter) ?? option.Parameter
      return { patchOptionParameters: true }
    }
    case "Statements": {
      return { statements: paginate(statementNodes(), source) }
    }
    default:
      throw new Error(`演示环境暂不支持该操作：${operation || "未知"}`)
  }
}
