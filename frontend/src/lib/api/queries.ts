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
          Database
          Subject
          Status
          CreateAt
          User { UUID Name }
          Reviewer { UUID Name }
          Cluster { UUID Alias Host IP Port }
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
            Cluster { UUID Alias Host IP Port }
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
            Cluster { UUID Alias Host IP Port }
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
  Status: number
  Group: string | null
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
          Status
          Group
          CreateAt
        }
      }
    }
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
      Cluster { UUID Alias Host IP Port }
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
