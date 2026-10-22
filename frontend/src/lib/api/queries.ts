/**
 * 后端查询文档与响应类型。字段名保持服务的 PascalCase 契约，
 * 页面层再转换成自己的 camelCase 模型。
 */

export interface MeUser {
  UUID: string
  Name: string | null
  Email: string
  Status: number
  Roles: { UUID: string; Name: string }[] | null
}

export const ME_QUERY = `
  query Me {
    me {
      UUID
      Name
      Email
      Status
      Roles { UUID Name }
    }
  }
`

export interface MetricNode {
  UUID: string
  Name: string
  Scope: number
  ScopeID: string | null
  Type: number
  Bucket: string
  Value: number
}

export const METRICS_QUERY = `
  query Metrics($Scope: UInt8, $Names: [String!]) {
    metrics(Scope: $Scope, Names: $Names) {
      UUID
      Name
      Scope
      ScopeID
      Type
      Bucket
      Value
    }
  }
`

export interface RuleHitCountNode {
  RuleName: string | null
  Level: number
  Count: number
}

export const RULE_HIT_LEVELS_QUERY = `
  query RuleHitLevels($from: UInt, $to: UInt) {
    ruleHitLevels(from: $from, to: $to) {
      RuleName
      Level
      Count
    }
  }
`

export const RULE_HIT_RULES_QUERY = `
  query RuleHitRules($from: UInt, $to: UInt, $first: Int) {
    ruleHitRules(from: $from, to: $to, first: $first) {
      RuleName
      Level
      Count
    }
  }
`

export interface MessageNode {
  UUID: string
  Category: string
  Title: string
  Content: string
  Link: string
  ReadAt: number | null
  CreateAt: number
}

export interface MessageConnectionNode {
  totalCount: number
  pageInfo: { hasNextPage: boolean; hasPreviousPage: boolean; startCursor: string; endCursor: string }
  edges: { cursor: string; node: MessageNode }[]
}

export const MESSAGES_QUERY = `
  query Messages($unreadOnly: Boolean, $category: String, $first: Int, $after: String) {
    messages(unreadOnly: $unreadOnly, category: $category, first: $first, after: $after) {
      totalCount
      pageInfo {
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
      edges {
        cursor
        node {
          UUID
          Category
          Title
          Content
          Link
          ReadAt
          CreateAt
        }
      }
    }
  }
`

export const REMOVE_MESSAGE_MUTATION = `
  mutation RemoveMessage($uuid: ID!) {
    removeMessage(uuid: $uuid)
  }
`

export const UNREAD_MESSAGE_COUNT_QUERY = `
  query UnreadMessageCount {
    unreadMessageCount
  }
`

export const MARK_MESSAGE_READ_MUTATION = `
  mutation MarkMessageRead($uuid: ID!) {
    markMessageRead(uuid: $uuid)
  }
`

export const MARK_ALL_MESSAGES_READ_MUTATION = `
  mutation MarkAllMessagesRead {
    markAllMessagesRead
  }
`

export interface GroupNode {
  UUID: string
  Name: string
  Description: string
  Status: number
  CreateAt: number
}

export const CREATE_GROUP_MUTATION = `
  mutation CreateGroup($input: CreateGroupInput!) {
    createGroup(input: $input) {
      UUID
      Name
      Description
      Status
      CreateAt
    }
  }
`

export const UPDATE_GROUP_MUTATION = `
  mutation UpdateGroup($input: UpdateGroupInput!) {
    updateGroup(input: $input) {
      UUID
      Name
      Description
      Status
      CreateAt
    }
  }
`

export const REMOVE_GROUP_MUTATION = `
  mutation RemoveGroup($id: ID!) {
    removeGroup(id: $id)
  }
`

export const GRANT_GROUPS_MUTATION = `
  mutation GrantGroups($input: GrantGroupsInput!) {
    grantGroups(input: $input)
  }
`

export const GRANT_CLUSTERS_MUTATION = `
  mutation GrantClusters($input: GrantClustersInput!) {
    grantClusters(input: $input)
  }
`

export const GROUPS_QUERY = `
  query Groups {
    groups {
      UUID
      Name
      Description
      Status
      CreateAt
    }
  }
`

export interface TicketNode {
  UUID: string
  /** 自增 ID，命中记录按它关联工单 */
  TicketID: number
  Database: string
  Subject: string
  Status: number
  CreateAt: number
  User: { UUID: string; Name: string | null } | null
  Reviewer: { UUID: string; Name: string | null } | null
  Cluster: {
    UUID: string
    Alias: string
    Host: string
    IP: string
    Port: number
    Type: string
    Group: string | null
  } | null
}

export interface TicketConnectionNode {
  totalCount: number
  edges: { cursor: string; node: TicketNode }[]
  pageInfo: {
    hasNextPage: boolean
    hasPreviousPage: boolean
    startCursor: string
    endCursor: string
  }
}

export const TICKETS_QUERY = `
  query Tickets($first: Int, $after: String, $status: UInt8, $keyword: String) {
    tickets(first: $first, after: $after, status: $status, keyword: $keyword) {
      totalCount
      pageInfo {
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
      edges {
        cursor
        node {
          UUID
          TicketID
          Database
          Subject
          Status
          CreateAt
          User { UUID Name }
          Reviewer { UUID Name }
          Cluster { UUID Alias Host IP Port Type Group }
        }
      }
    }
  }
`

export interface LoginPayloadNode {
  Me: { UUID: string; Name: string | null; Email: string }
  Token: string
}

export const LOGIN_MUTATION = `
  mutation Login($input: UserLoginInput!) {
    login(input: $input) {
      Me { UUID Name Email }
      Token
    }
  }
`

export const MY_TICKETS_QUERY = `
  query MyTickets($first: Int, $after: String) {
    me {
      MyTickets(first: $first, after: $after) {
        totalCount
        pageInfo {
          hasNextPage
          hasPreviousPage
          startCursor
          endCursor
        }
        edges {
          cursor
          node {
            UUID
            Database
            Subject
            Status
            CreateAt
            User { UUID Name }
            Reviewer { UUID Name }
            Cluster { UUID Alias Host IP Port Type Group }
          }
        }
      }
    }
  }
`

export const PENDING_REVIEWS_QUERY = `
  query PendingReviews($first: Int, $after: String) {
    me {
      PendingReviews(first: $first, after: $after) {
        totalCount
        pageInfo {
          hasNextPage
          hasPreviousPage
          startCursor
          endCursor
        }
        edges {
          cursor
          node {
            UUID
            Database
            Subject
            Status
            CreateAt
            User { UUID Name }
            Reviewer { UUID Name }
            Cluster { UUID Alias Host IP Port Type Group }
          }
        }
      }
    }
  }
`

export interface ClusterNode {
  UUID: string
  Alias: string
  Host: string
  IP: string
  Port: number
  User: string
  Status: number
  /** 目标库类型：mysql / postgres */
  Type: string
  Group: string | null
  GroupUUID: string | null
  CreateAt: number
}

export const CLUSTERS_QUERY = `
  query Clusters($first: Int, $after: String) {
    clusters(first: $first, after: $after) {
      totalCount
      pageInfo {
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
      edges {
        cursor
        node {
          UUID
          Alias
          Host
          IP
          Port
          User
          Status
          Type
          Group
          GroupUUID
          CreateAt
        }
      }
    }
  }
`

export const CREATE_CLUSTER_MUTATION = `
  mutation CreateCluster($input: CreateClusterInput!) {
    createCluster(input: $input) {
      UUID
      Alias
      Host
      IP
      Port
      Status
      Type
      Group
    }
  }
`

export const UPDATE_CLUSTER_MUTATION = `
  mutation UpdateCluster($input: UpdateClusterInput!) {
    updateCluster(input: $input) {
      UUID
      Alias
      Host
      IP
      Port
      Status
      Type
      Group
    }
  }
`

export const REMOVE_CLUSTER_MUTATION = `
  mutation RemoveCluster($id: ID!) {
    removeCluster(id: $id)
  }
`

export interface StatementNode {
  UUID: string
  /** 所属工单的自增 ID */
  TicketID: number
  Sequence: number
  Content: string
  TypeDesc: string
  Status: number
  Report: string
  RowsAffected: number | null
  CreateAt: number
  UpdateAt: number | null
}

export interface CommentNode {
  UUID: string
  Content: string
  CreateAt: number
  User: { UUID: string; Name: string | null }
}

export const CREATE_TICKET_MUTATION = `
  mutation CreateTicket($input: CreateTicketInput!) {
    createTicket(input: $input) {
      UUID
      Subject
      Status
    }
  }
`

export const UPDATE_TICKET_MUTATION = `
  mutation UpdateTicket($input: UpdateTicketInput!) {
    updateTicket(input: $input) {
      UUID
      Subject
      Status
    }
  }
`

export const REVIEW_TICKET_MUTATION = `
  mutation ReviewTicket($input: ReviewTicketInput!) {
    reviewTicket(input: $input) {
      UUID
      Status
    }
  }
`

export const CREATE_COMMENT_MUTATION = `
  mutation CreateComment($input: CreateCommentInput!) {
    createComment(input: $input) {
      UUID
      Content
      CreateAt
      User { UUID Name }
    }
  }
`

export const TICKET_DETAIL_QUERY = `
  query Ticket($id: ID!) {
    ticket(id: $id) {
      UUID
      Subject
      Content
      Database
      Status
      CreateAt
      UpdateAt
      Cluster { UUID Alias Host IP Port Type Group }
      User { UUID Name Email }
      Reviewer { UUID Name }
      Statements(first: 20) {
        edges { node { UUID Sequence Content TypeDesc Status Report } }
      }
      Comments(first: 20) {
        edges { node { UUID Content CreateAt User { UUID Name } } }
      }
    }
  }
`

export const ME_CLUSTERS_QUERY = `
  query MeClusters($first: Int) {
    me {
      Clusters(first: $first) {
        edges {
          node {
            UUID
            Alias
            Host
            IP
            Port
            Type
            Group
          }
        }
      }
    }
  }
`

export const ME_REVIEWERS_QUERY = `
  query MeReviewers {
    me {
      Reviewers {
        UUID
        Name
        Email
      }
    }
  }
`

export const TEST_CLUSTER_QUERY = `
  query TestCluster($input: ValidateConnectionInput!) {
    testCluster(input: $input) {
      Success
      Message
      Version
      LatencyMs
      Databases
    }
  }
`

export interface ClusterTestResult {
  Success: boolean
  Message: string
  Version: string
  LatencyMs: number
  Databases: string[]
}

export const RULES_QUERY = `
  query Rules {
    rules {
      UUID
      Name
      Group
      Level
      Description
      Parameters
      Bitwise
      Func
      Message
    }
  }
`

export const RULE_HITS_QUERY = `
  query RuleHits($first: Int) {
    ruleHits(first: $first) {
      UUID
      TicketID
      Sequence
      RuleName
      Level
      Message
      CreateAt
    }
  }
`

export const LOGS_QUERY = `
  query Logs($first: Int, $after: String) {
    logs(first: $first, after: $after) {
      totalCount
      pageInfo {
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
      edges {
        cursor
        node {
          UUID
          LogID
          User { Name Email }
          Operation
          CreateAt
        }
      }
    }
  }
`

export interface RuleNode {
  UUID: string
  Name: string
  Group: number
  /** 严重级别：0 提示 1 错误 2 警告 */
  Level: number
  Description: string
  Parameters: string
  Bitwise: number
  Func: string
  Message: string
}

export interface RuleHitNode {
  UUID: string
  TicketID: number
  Sequence: number
  RuleName: string
  /** 严重级别：0 提示 1 错误 2 警告 */
  Level: number
  Message: string
  CreateAt: number
}

export interface LogNode {
  UUID: string
  LogID: number
  User: { Name: string | null; Email: string }
  /** 形如 "ticket.update|ticket:<uuid>|" 的结构化日志 */
  Operation: string
  CreateAt: number
}

export interface LogConnectionNode {
  totalCount: number
  edges: { cursor: string; node: LogNode }[]
  pageInfo: {
    hasNextPage: boolean
    hasPreviousPage: boolean
    startCursor: string
    endCursor: string
  }
}

export const CRONS_QUERY = `
  query Crons($first: Int, $after: String) {
    crons(first: $first, after: $after) {
      totalCount
      pageInfo {
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
      edges {
        cursor
        node {
          UUID
          Name
          Cmd
          Params
          Interval
          Duration
          LastRun
          NextRun
          Recurrent
          Status
          CreateAt
        }
      }
    }
  }
`

export const PENDING_TASKS_QUERY = `
  query PendingTasks {
    tasks {
      UUID
      Name
      Cmd
      Params
      Interval
      LastRun
      NextRun
      Recurrent
      Status
    }
  }
`

export const TASK_EXECUTIONS_QUERY = `
  query TaskExecutions($status: UInt8, $first: Int) {
    taskExecutions(status: $status, first: $first) {
      UUID
      Name
      Type
      Payload
      Status
      Attempt
      CronID
      ScheduledAt
      StartedAt
      FinishedAt
      Error
      CreateAt
    }
    taskExecutionSummary {
      Status
      Count
    }
  }
`

export const CREATE_CRON_MUTATION = `
  mutation CreateCron($input: CreateCronInput!) {
    createCron(input: $input) {
      UUID
      Name
      Interval
      Status
    }
  }
`

export const REMOVE_CRON_MUTATION = `
  mutation RemoveCron($id: ID!) {
    removeCron(id: $id)
  }
`

export const DATABASES_QUERY = `
  query Databases($clusterUUID: String!) {
    databases(ClusterUUID: $clusterUUID) {
      Name
      Charset
      Collate
    }
  }
`

export const USERS_QUERY = `
  query Users($first: Int, $after: String) {
    users(first: $first, after: $after) {
      totalCount
      pageInfo {
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
      edges {
        cursor
        node {
          UUID
          Email
          Name
          Status
          Roles { UUID Name }
          Reviewers { UUID Name Email }
          GrantedClusters { UUID Alias }
          ClusterGrants {
            Cluster { UUID Alias Group }
            Direct
            Inherited
            Sources
          }
        }
      }
    }
  }
`

export const CREATE_USER_MUTATION = `
  mutation CreateUser($input: CreateUserInput!) {
    createUser(input: $input) {
      UUID
      Email
      Name
      Status
    }
  }
`

export const UPDATE_USER_MUTATION = `
  mutation UpdateUser($input: UpdateUserInput!) {
    updateUser(input: $input) {
      UUID
      Email
      Name
      Status
    }
  }
`

export const PATCH_USER_STATUS_MUTATION = `
  mutation PatchUserStatus($input: PatchUserStatusInput!) {
    patchUserStatus(input: $input)
  }
`

export const GRANT_ROLES_MUTATION = `
  mutation GrantRoles($input: GrantRolesInput!) {
    grantRoles(input: $input)
  }
`

export const GRANT_REVIEWERS_MUTATION = `
  mutation GrantReviewers($input: GrantReviewersInput!) {
    grantReviewers(input: $input)
  }
`

export const ROLES_QUERY = `
  query Roles {
    roles {
      UUID
      Name
      Description
      Groups { UUID Name }
    }
  }
`

export interface SystemInfoNode {
  Version: string
  Git: string
  BuildAt: string
  GoVer: string
  Branch: string
  Database: string
  MailEnabled: boolean
  ClusterCount: number
  UserCount: number
}

export const SYSTEM_INFO_QUERY = `
  query SystemInfo {
    systemInfo {
      Version
      Git
      BuildAt
      GoVer
      Branch
      Database
      MailEnabled
      ClusterCount
      UserCount
    }
  }
`

export const OPTIONS_QUERY = `
  query Options {
    options {
      UUID
      Name
      Parameter
      Writable
      Description
      Element
    }
  }
`

export const PATCH_OPTION_MUTATION = `
  mutation PatchOption($input: PatchOptionParameterInput!) {
    patchOptionParameters(input: $input)
  }
`

export const SEND_TEST_MAIL_MUTATION = `
  mutation SendTestMail($email: String!) {
    sendTestMail(email: $email)
  }
`

export const METRIC_SUMMARY_QUERY = `
  query MetricSummary($Scope: UInt8, $Names: [String!]) {
    metrics(Scope: $Scope, Names: $Names) {
      Name
      Scope
      ScopeID
      Type
      Bucket
      Value
      CreateAt
    }
  }
`

export interface CronNode {
  UUID: string
  Name: string
  Cmd: string
  Params: string | null
  Interval: string | null
  Duration: string | null
  LastRun: string | null
  NextRun: string | null
  Recurrent: number
  Status: string
  CreateAt: number
}

export interface TaskExecutionNode {
  UUID: string
  Name: string
  Type: string
  Payload: string
  Status: number
  Attempt: number
  CronID: number
  ScheduledAt: number
  StartedAt: number | null
  FinishedAt: number | null
  Error: string
  CreateAt: number
}

export interface UserNode {
  UUID: string
  Email: string
  Name: string | null
  Status: number
  Roles: { UUID: string; Name: string }[] | null
  Reviewers: { UUID: string; Name: string | null; Email: string }[] | null
  /** 直接授权给该用户的群集，不含角色继承的分组权限 */
  GrantedClusters: { UUID: string; Alias: string }[] | null
  /** 对每个群集的生效权限与来源 */
  ClusterGrants?: ClusterGrantNode[] | null
}

export interface ClusterGrantNode {
  Cluster: { UUID: string; Alias: string; Group: string | null }
  /** 是否有直接授权 */
  Direct: boolean
  /** 是否通过角色继承的分组获得权限 */
  Inherited: boolean
  /** 继承来源，形如「角色 / 分组」 */
  Sources: string[]
}

export interface RoleNode {
  UUID: string
  Name: string
  Description: string
  Groups: { UUID: string; Name: string }[] | null
}

export interface OptionNode {
  UUID: string
  Name: string
  Parameter: string
  Writable: number
  Description: string
  Element: string
}

export const STATEMENTS_QUERY = `
  query Statements($first: Int, $after: String) {
    statements(first: $first, after: $after) {
      totalCount
      pageInfo {
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
      edges {
        cursor
        node {
          UUID
          TicketID
          Sequence
          Content
          TypeDesc
          Status
          Report
          RowsAffected
          CreateAt
          UpdateAt
        }
      }
    }
  }
`
