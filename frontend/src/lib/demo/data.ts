/**
 * 演示数据种子。演示环境每天重置一次数据，所以时间都相对「今天」生成，
 * 界面看起来始终是新的。所有结构都按 GraphQL 契约命名，直接喂给页面。
 */

import type {
  ClusterNode,
  CommentNode,
  CronNode,
  GroupNode,
  LogNode,
  MessageNode,
  OptionNode,
  RoleNode,
  RuleHitNode,
  RuleNode,
  StatementNode,
  TaskExecutionNode,
  TicketNode,
  UserNode,
} from "@/lib/api/queries"

const HOUR = 3600
const DAY = 86400
const DAYS = 365

/** 以天为粒度的基准时间，保证同一天里生成的种子完全一致。 */
export const today = Math.floor(Date.now() / 1000 / DAY) * DAY

/** 固定种子的伪随机，避免每次刷新演示数据跳变。 */
function seeded(seed: number): () => number {
  let state = seed % 2147483647
  if (state <= 0) state += 2147483646
  return () => {
    state = (state * 16807) % 2147483647
    return (state - 1) / 2147483646
  }
}

const random = seeded(20261021)

/** 形如 UUID 的稳定标识，日志页会截前 8 位展示。 */
export function demoUUID(prefix: number, id: number): string {
  return `${(prefix + id).toString(16).padStart(8, "0")}-0000-4000-8000-${id.toString(16).padStart(12, "0")}`
}

function pick<T>(items: T[], index: number): T {
  return items[index % items.length]!
}

export interface DemoTicket extends TicketNode {
  Content: string
  UpdateAt: number | null
  Statements: StatementNode[]
  Comments: CommentNode[]
}

export interface DemoState {
  groups: GroupNode[]
  clusters: ClusterNode[]
  databases: Record<string, { Name: string; Charset: string; Collate: string }[]>
  users: UserNode[]
  roles: RoleNode[]
  tickets: DemoTicket[]
  rules: RuleNode[]
  hits: RuleHitNode[]
  logs: LogNode[]
  crons: CronNode[]
  tasks: CronNode[]
  executions: TaskExecutionNode[]
  options: OptionNode[]
  messages: MessageNode[]
}

const groupSeeds = [
  ["核心交易", "订单、支付相关库"],
  ["用户中心", "账号与权限库"],
  ["报表分析", "只读报表库"],
  ["离线数据", "归档与备份库"],
  ["灰度环境", "预发布验证库"],
]

const roleSeeds = [
  ["Root", "平台最高权限"],
  ["Admin", "管理员，可配置平台与用户"],
  ["Developer", "提交并跟踪自己的工单"],
  ["Reviewer", "审核他人提交的工单"],
  ["User", "只读访问"],
  ["Viewer", "仅查看统计与报表"],
]

const userSeeds: [string, string, number, number[]][] = [
  ["演示账号", "demo@halo.local", 1, [0, 2, 3]],
  ["张三", "zhangsan@example.com", 1, [2]],
  ["李四", "lisi@example.com", 1, [3]],
  ["王五", "wangwu@example.com", 1, [1]],
  ["赵六", "zhaoliu@example.com", 1, [4]],
  ["运维值班", "oncall@example.com", 2, [3]],
]

const clusterSeeds = [
  ["mysql-prod-01", "prod-trade-01", "10.20.1.11", 3306, "mysql", 0],
  ["mysql-prod-02", "prod-user-02", "10.20.1.12", 3306, "mysql", 1],
  ["pg-analytics-01", "prod-report-01", "10.20.2.21", 5432, "postgres", 2],
  ["pg-archive-01", "archive-01", "10.20.3.31", 5432, "postgres", 3],
  ["mysql-stage-01", "stage-trade-01", "10.30.1.41", 3306, "mysql", 4],
  ["pg-stage-01", "stage-report-01", "10.30.2.51", 5432, "postgres", 4],
] as const

const databaseSeeds = [
  ["app_order", "app_payment", "app_settlement", "app_refund"],
  ["app_user", "app_auth", "app_profile"],
  ["dw_report", "dw_dimension", "dw_fact"],
  ["archive_2024", "archive_2025", "archive_2026"],
  ["stage_order", "stage_payment"],
  ["stage_report", "stage_dimension"],
]

interface SubjectSeed {
  subject: string
  database: string
  type: string
  sql: string
  report: { Level: number; Description: string; Rule: string }[]
}

const subjectSeeds: SubjectSeed[] = [
  {
    subject: "订单表新增归档字段",
    database: "app_order",
    type: "ALTER_TABLE",
    sql: "ALTER TABLE orders ADD COLUMN archived_at DATETIME NULL;\nCREATE INDEX idx_orders_archived ON orders (archived_at);",
    report: [{ Level: 0, Description: "新增可空列不会锁表，索引建议在低峰期创建。", Rule: "CTB-L3-006" }],
  },
  {
    subject: "补齐支付流水索引",
    database: "app_payment",
    type: "CREATE_INDEX",
    sql: "CREATE INDEX idx_payment_trade_no ON payment_flow (trade_no);",
    report: [],
  },
  {
    subject: "清理测试环境过期数据",
    database: "stage_order",
    type: "DELETE",
    sql: "DELETE FROM orders WHERE created_at < '2025-01-01' AND status = 9;",
    report: [{ Level: 2, Description: "DELETE 命中范围较大，建议先备份。", Rule: "DEL-L2-001" }],
  },
  {
    subject: "用户表增加最近登录时间",
    database: "app_user",
    type: "ALTER_TABLE",
    sql: "ALTER TABLE users ADD COLUMN last_login_at DATETIME NULL;\nUPDATE users SET last_login_at = updated_at WHERE last_login_at IS NULL;",
    report: [{ Level: 1, Description: "UPDATE 缺少 WHERE 主键条件，会全表扫描。", Rule: "UPD-L1-002" }],
  },
  {
    subject: "报表视图重建",
    database: "dw_report",
    type: "CREATE_VIEW",
    sql: "DROP VIEW IF EXISTS v_order_daily;\nCREATE VIEW v_order_daily AS SELECT date_trunc('day', created_at) AS day, count(*) FROM orders GROUP BY 1;",
    report: [],
  },
  {
    subject: "结算表分区调整",
    database: "app_settlement",
    type: "ALTER_TABLE",
    sql: "ALTER TABLE settlement ADD PARTITION (PARTITION p2027 VALUES LESS THAN (2028));",
    report: [{ Level: 0, Description: "建议确认分区键与历史数据一致。", Rule: "CTB-L3-002" }],
  },
  {
    subject: "归档库补建唯一索引",
    database: "archive_2025",
    type: "CREATE_INDEX",
    sql: "CREATE UNIQUE INDEX uk_archive_no ON archive_orders (archive_no);",
    report: [],
  },
  {
    subject: "退款表新增审核状态",
    database: "app_refund",
    type: "ALTER_TABLE",
    sql: "ALTER TABLE refund ADD COLUMN audit_status TINYINT NOT NULL DEFAULT 0;",
    report: [],
  },
  {
    subject: "统计维度表补字段注释",
    database: "dw_dimension",
    type: "ALTER_TABLE",
    sql: "ALTER TABLE dim_region MODIFY COLUMN region_name VARCHAR(64) COMMENT '地区名称';",
    report: [],
  },
  {
    subject: "灰度环境重建订单表",
    database: "stage_order",
    type: "DROP_TABLE",
    sql: "DROP TABLE IF EXISTS orders_old;",
    report: [{ Level: 2, Description: "DROP TABLE 不可恢复，请确认已备份。", Rule: "DTB-L2-001" }],
  },
  {
    subject: "权限表补充角色映射",
    database: "app_auth",
    type: "INSERT",
    sql: "INSERT INTO role_map (role_id, user_id) VALUES (3, 1001), (3, 1002);",
    report: [],
  },
  {
    subject: "逻辑删除订单明细",
    database: "app_order",
    type: "UPDATE",
    sql: "UPDATE order_items SET deleted = 1 WHERE order_id IN (SELECT id FROM orders WHERE status = 9);",
    report: [{ Level: 1, Description: "子查询没有索引支撑。", Rule: "SEL-L1-004" }],
  },
]

const statusPool = [1, 4, 7, 1, 4, 7, 3, 2, 5, 6, 8, 9]

/** 数据库落在哪个群集上，别的数据源看起来才合理。 */
const databaseCluster: Record<string, number> = {
  app_order: 0,
  app_payment: 0,
  app_settlement: 0,
  app_refund: 0,
  app_user: 1,
  app_auth: 1,
  app_profile: 1,
  dw_report: 2,
  dw_dimension: 2,
  dw_fact: 2,
  archive_2024: 3,
  archive_2025: 3,
  archive_2026: 3,
  stage_order: 4,
  stage_payment: 4,
  stage_report: 5,
  stage_dimension: 5,
}

const commentSeeds = [
  "已在低峰期窗口执行，注意回滚脚本也要一起提交。",
  "这条语句会锁表，建议分批执行。",
  "确认过索引选择性，可以上线。",
  "请补充回滚方案后再提交审核。",
  "灰度环境已验证通过。",
]

function buildTickets(users: UserNode[], clusters: ClusterNode[]): DemoTicket[] {
  const tickets: DemoTicket[] = []
  for (let index = 0; index < 28; index += 1) {
    const subject = pick(subjectSeeds, index)
    const cluster = clusters[databaseCluster[subject.database] ?? 0]!
    const submitter = pick(users, index % 3 === 0 ? 0 : index + 1)
    const createAt = today - index * 6 * HOUR - 2 * HOUR
    const status = pick(statusPool, index)
    // 待人工审核的工单都挂在演示账号名下，审核队列才有内容。
    const reviewer = status === 4 ? users[0]! : pick(users, index + 2)
    const content = subject.sql

    const statements: StatementNode[] = content.split(";").filter(Boolean).map((chunk, seq) => {
      const text = `${chunk.trim()};`
      const report = seq === 0 ? subject.report : []
      return {
        UUID: demoUUID(0x60000000, index * 10 + seq),
        TicketID: index + 1,
        Sequence: seq + 1,
        Content: text,
        TypeDesc: subject.type,
        Status: status === 8 ? 3 : 1,
        Report: JSON.stringify(report),
        RowsAffected: status === 7 ? random() * 200 | 0 : null,
        CreateAt: createAt,
        UpdateAt: null,
      }
    })

    const comments: CommentNode[] = status >= 4
      ? Array.from({ length: status === 9 ? 1 : 2 }, (_, seq) => ({
          UUID: demoUUID(0x70000000, index * 10 + seq),
          Content: pick(commentSeeds, index + seq),
          CreateAt: createAt + (seq + 1) * HOUR,
          User: { UUID: reviewer.UUID, Name: reviewer.Name },
        }))
      : []

    tickets.push({
      UUID: demoUUID(0x10000000, index + 1),
      TicketID: index + 1,
      Database: subject.database,
      Subject: subject.subject,
      Status: status,
      CreateAt: createAt,
      User: { UUID: submitter.UUID, Name: submitter.Name },
      Reviewer: { UUID: reviewer.UUID, Name: reviewer.Name },
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
      UpdateAt: comments.length > 0 ? createAt + HOUR : null,
      Statements: statements,
      Comments: comments,
    } as DemoTicket & { Group?: string | null })
  }
  return tickets
}

const rulePrefixes = [
  ["UPD", 1],
  ["DEL", 2],
  ["SEL", 1],
  ["INS", 0],
  ["CTB", 0],
  ["DTB", 2],
  ["CIX", 0],
  ["RIX", 0],
  ["CVW", 1],
  ["DVW", 2],
  ["MSC", 0],
] as const

const ruleDescriptions: Record<string, string> = {
  UPD: "UPDATE 语句的变更检查",
  DEL: "DELETE 语句的删除范围检查",
  SEL: "SELECT 语句的索引检查",
  INS: "INSERT 语句的字段检查",
  CTB: "建表语句的规范检查",
  DTB: "删表语句的风险检查",
  CIX: "建索引语句的规范检查",
  RIX: "删索引语句的风险检查",
  CVW: "建视图语句的规范检查",
  DVW: "删视图语句的风险检查",
  MSC: "其它风险检查",
}

function buildRules(): RuleNode[] {
  const rules = rulePrefixes.flatMap(([prefix, level], index) => {
    const count = prefix === "UPD" || prefix === "DEL" ? 3 : 2
    return Array.from({ length: count }, (_, seq) => {
      const name = `${prefix}-L${level + 1}-${String(seq + 1).padStart(3, "0")}`
      return {
        UUID: `rule-${index}-${seq}`,
        Name: name,
        Group: 1,
        Level: level,
        Description: ruleDescriptions[prefix] ?? "风险检查",
        Parameters: JSON.stringify({ enabled: true, threshold: (seq + 1) * 100 }),
        Bitwise: prefix === "SEL" ? 4 : prefix === "UPD" ? 2 : 3,
        Func: "internal.analyzer",
        Message: `${name} 命中了待确认的写法`,
      }
    })
  })
  // The real service lists rules by rule ID, so the demo keeps the same order.
  return rules.sort((a, b) => a.Name.localeCompare(b.Name))
}

function buildHits(rules: RuleNode[], tickets: DemoTicket[]): RuleHitNode[] {
  const hits: RuleHitNode[] = []
  for (let index = 0; index < 26; index += 1) {
    const ticket = pick(tickets, index)
    const rule = pick(rules, index * 3)
    hits.push({
      UUID: demoUUID(0x80000000, index),
      TicketID: ticket.TicketID,
      Sequence: (index % 2) + 1,
      RuleName: rule.Name,
      Level: rule.Level,
      Message: `${ticket.Subject}：${rule.Description}`,
      CreateAt: ticket.CreateAt + 600,
    })
  }
  return hits
}

function buildLogs(users: UserNode[], tickets: DemoTicket[]): LogNode[] {
  const actions = [
    "ticket.create",
    "ticket.update",
    "ticket.execute",
    "cluster.create",
    "cluster.update",
    "user.update",
    "option.update",
    "user.login",
  ]
  return Array.from({ length: 42 }, (_, index) => {
    const action = pick(actions, index)
    const ticket = pick(tickets, index)
    const user = pick(users, index)
    const target = action.startsWith("ticket")
      ? `ticket:${ticket.UUID}`
      : action.startsWith("cluster")
        ? `cluster:${ticket.Cluster?.UUID ?? "cluster"}`
        : action.startsWith("user")
          ? `user:${user.UUID}`
          : "option:site-name"
    return {
      UUID: demoUUID(0x90000000, index),
      LogID: 2000 - index,
      User: { Name: user.Name, Email: user.Email },
      Operation: action === "user.login" ? `${user.Email} 登录系统。` : `${action}|${target}|`,
      CreateAt: today - index * 1800,
    }
  })
}

function buildCrons(): { crons: CronNode[]; tasks: CronNode[]; executions: TaskExecutionNode[] } {
  const crons: CronNode[] = [
    {
      UUID: "cron-1",
      Name: "每日指标汇总",
      Cmd: "METRICS_ROLLUP",
      Params: "{}",
      Interval: "24h",
      Duration: "3s",
      LastRun: new Date((today - DAY) * 1000).toISOString(),
      NextRun: new Date((today + DAY) * 1000).toISOString(),
      Recurrent: 1,
      Status: "PENDING",
      CreateAt: today - 30 * DAY,
    },
    {
      UUID: "cron-2",
      Name: "过期指标清理",
      Cmd: "METRICS_CLEANUP",
      Params: '{"days":90}',
      Interval: "24h",
      Duration: "1s",
      LastRun: new Date((today - DAY) * 1000).toISOString(),
      NextRun: new Date((today + DAY) * 1000).toISOString(),
      Recurrent: 1,
      Status: "PENDING",
      CreateAt: today - 30 * DAY,
    },
    {
      UUID: "cron-3",
      Name: "任务状态巡检",
      Cmd: "TASK_INSPECT",
      Params: "{}",
      Interval: "1h",
      Duration: "0s",
      LastRun: new Date((today - HOUR) * 1000).toISOString(),
      NextRun: new Date((today + HOUR) * 1000).toISOString(),
      Recurrent: 1,
      Status: "DONE",
      CreateAt: today - 7 * DAY,
    },
  ]

  const tasks: CronNode[] = [
    { ...crons[0]!, UUID: "task-1", Name: "等待执行的指标汇总", Status: "PENDING" },
    { ...crons[1]!, UUID: "task-2", Name: "等待执行的指标清理", Status: "PENDING" },
    { ...crons[2]!, UUID: "task-3", Name: "等待执行的任务巡检", Status: "RUNNING" },
  ]

  const executions: TaskExecutionNode[] = Array.from({ length: 16 }, (_, index) => {
    const status = index === 0 ? 2 : index % 7 === 3 ? 3 : 1
    const scheduled = today - index * 3 * HOUR
    return {
      UUID: demoUUID(0xb0000000, index),
      Name: index % 2 === 0 ? "每日指标汇总" : "过期指标清理",
      Type: index % 2 === 0 ? "METRICS_ROLLUP" : "METRICS_CLEANUP",
      Payload: "{}",
      Status: status,
      Attempt: status === 3 ? 3 : 1,
      CronID: index % 2 === 0 ? 1 : 2,
      ScheduledAt: scheduled,
      StartedAt: status === 1 || status === 3 ? scheduled + 5 : scheduled,
      FinishedAt: status === 2 ? null : scheduled + 40,
      Error: status === 3 ? "连接元数据库超时" : "",
      CreateAt: scheduled,
    }
  })

  return { crons, tasks, executions }
}

function buildOptions(): OptionNode[] {
  const seeds = [
    ["site-name", "Halo", 1, "站点名称", "input"],
    ["ticket-auto-close-days", "30", 1, "工单自动关闭天数", "input"],
    ["ticket-max-statements", "200", 1, "单工单语句上限", "input"],
    ["execute-timeout", "600", 1, "执行超时（秒）", "input"],
    ["review-required", "true", 0, "是否强制人工审核", "switch"],
    ["metrics-retention", "90", 0, "指标保留天数", "input"],
    ["mail-notify", "false", 1, "邮件通知开关", "switch"],
    ["danger-keywords", "DROP,TRUNCATE", 1, "危险关键字", "input"],
    ["backup-before-execute", "true", 0, "执行前备份", "switch"],
    ["session-timeout", "7200", 1, "登录会话有效期（秒）", "input"],
    ["cluster-default-port", "3306", 1, "群集默认端口", "input"],
    ["audit-log-retention", "180", 0, "审计日志保留天数", "input"],
  ] as const
  return seeds.map(([name, parameter, writable, description, element], index) => ({
    UUID: demoUUID(0xc0000000, index),
    Name: name,
    Parameter: parameter,
    Writable: writable,
    Description: description,
    Element: element,
  }))
}

function buildMetrics(): { Name: string; Bucket: string; Value: number }[] {
  const names = [
    "tickets-created",
    "tickets-executed",
    "tickets-failed",
    "violations-error",
    "violations-warning",
    "violations-notice",
  ]
  const rows: { Name: string; Bucket: string; Value: number }[] = []
  for (let offset = DAYS - 1; offset >= 0; offset -= 1) {
    const stamp = today - offset * DAY
    const date = new Date(stamp * 1000)
    const weekday = date.getDay()
    const busy = weekday === 0 || weekday === 6 ? 0.3 : 1
    const bucket = date.toISOString().slice(0, 10)
    for (const name of names) {
      const base = name.startsWith("tickets") ? 6 : 14
      const value = Math.round((base + random() * base) * busy * (name.endsWith("failed") || name.endsWith("error") ? 0.25 : 1))
      rows.push({ Name: name, Bucket: bucket, Value: value })
    }
  }
  return rows
}

function buildMessages(tickets: TicketNode[]): MessageNode[] {
  const approved = tickets.find((ticket) => ticket.Status === 5) ?? tickets[0]
  const pending = tickets.find((ticket) => ticket.Status === 4) ?? tickets[0]
  const seeds: [string, string, string, boolean][] = [
    ["ticket", "工单「订单表新增归档字段」审核通过", "可以执行，执行窗口为今晚 22:00。", false],
    ["review", "有一条工单待你审核", "「报表视图重建」等待人工审核。", false],
    ["ticket", "工单「清理测试环境过期数据」执行失败", "执行超时，已回滚。", false],
    ["system", "每日指标已汇总", "昨日指标已生成，可在统计报表中查看。", true],
    ["system", "演示数据已恢复", "演示环境每天零点重置一次数据。", true],
  ]
  return seeds.map(([category, title, content, read], index) => ({
    UUID: `message-${index}`,
    Category: category,
    Title: title,
    Content: content,
    Link:
      index === 0 && approved
        ? `/tickets/${approved.UUID}`
        : index === 1 && pending
          ? `/tickets/${pending.UUID}`
          : "",
    ReadAt: read ? today - HOUR * (index + 1) : null,
    CreateAt: today - HOUR * (index + 1),
  }))
}

export function buildDemoState(): DemoState {
  const groups: GroupNode[] = groupSeeds.map(([name, description], index) => ({
    UUID: demoUUID(0x40000000, index + 1),
    Name: name,
    Description: description,
    Status: 1,
    CreateAt: today - 120 * DAY,
  }))

  const roles: RoleNode[] = roleSeeds.map(([name, description], index) => ({
    UUID: demoUUID(0x50000000, index + 1),
    Name: name,
    Description: description,
    Groups: index <= 2 ? groups.map((group) => ({ UUID: group.UUID, Name: group.Name })) : null,
  }))

  const clusters: ClusterNode[] = clusterSeeds.map(([alias, host, ip, port, type, group], index) => ({
    UUID: demoUUID(0x20000000, index + 1),
    Alias: alias,
    Host: host,
    IP: ip,
    Port: port,
    User: type === "postgres" ? "postgres" : "root",
    Status: index === 5 ? 2 : 1,
    Type: type,
    Group: groups[group]!.Name,
    GroupUUID: groups[group]!.UUID,
    CreateAt: today - (60 - index) * DAY,
  }))

  const users: UserNode[] = userSeeds.map(([name, email, status, roleIndexes], index) => ({
    UUID: demoUUID(0x30000000, index + 1),
    Email: email,
    Name: name,
    Status: status,
    Roles: roleIndexes.map((roleIndex) => ({
      UUID: roles[roleIndex]!.UUID,
      Name: roles[roleIndex]!.Name,
    })),
    Reviewers: index === 0 ? [{ UUID: demoUUID(0x30000000, 3), Name: "李四", Email: "lisi@example.com" }] : [],
    // 演示账号持有全部群集，其余账号按序号分到两三个。
    GrantedClusters: clusters
      .filter((_, clusterIndex) => index === 0 || (clusterIndex + index) % 3 !== 0)
      .map((cluster) => ({ UUID: cluster.UUID, Alias: cluster.Alias })),
  }))

  const databases: DemoState["databases"] = Object.fromEntries(
    clusters.map((cluster, index) => [
      cluster.UUID,
      (databaseSeeds[index] ?? ["app_main"]).map((name) => ({
        Name: name,
        Charset: cluster.Type === "postgres" ? "UTF8" : "utf8mb4",
        Collate: cluster.Type === "postgres" ? "C" : "utf8mb4_general_ci",
      })),
    ]),
  )

  const tickets = buildTickets(users, clusters)
  const rules = buildRules()
  const hits = buildHits(rules, tickets)
  const logs = buildLogs(users, tickets)
  const { crons, tasks, executions } = buildCrons()

  return {
    groups,
    clusters,
    databases,
    users,
    roles,
    tickets,
    rules,
    hits,
    logs,
    crons,
    tasks,
    executions,
    options: buildOptions(),
    messages: buildMessages(tickets),
  }
}

export const metricSeeds = buildMetrics()

export const demoUser: UserNode = {
  UUID: demoUUID(0x30000000, 1),
  Email: "demo@halo.local",
  Name: "演示账号",
  Status: 1,
  Roles: [
    { UUID: demoUUID(0x50000000, 1), Name: "Root" },
    { UUID: demoUUID(0x50000000, 3), Name: "Developer" },
    { UUID: demoUUID(0x50000000, 4), Name: "Reviewer" },
  ],
  Reviewers: [{ UUID: demoUUID(0x30000000, 3), Name: "李四", Email: "lisi@example.com" }],
  GrantedClusters: [],
}
