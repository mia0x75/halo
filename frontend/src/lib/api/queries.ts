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

export const MESSAGES_QUERY = `
  query Messages($unreadOnly: Boolean, $first: Int) {
    messages(unreadOnly: $unreadOnly, first: $first) {
      UUID
      Category
      Title
      Content
      Link
      ReadAt
      CreateAt
    }
  }
`

export const UNREAD_MESSAGE_COUNT_QUERY = `
  query UnreadMessageCount {
    unreadMessageCount
  }
`

export interface GroupNode {
  UUID: string
  Name: string
  Description: string
  Status: number
  CreateAt: number
}

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
  Cluster: { UUID: string; Alias: string; Host: string; IP: string; Port: number } | null
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
  Sequence: number
  Content: string
  TypeDesc: string
  Status: number
  Report: string
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
