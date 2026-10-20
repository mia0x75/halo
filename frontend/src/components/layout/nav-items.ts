import {
  IconAlertTriangle,
  IconCalendarTime,
  IconChartBar,
  IconDatabase,
  IconFileText,
  IconHistory,
  IconLayoutDashboard,
  IconListCheck,
  IconServer,
  IconSettings,
  IconShieldCheck,
  IconTicket,
  IconUsers,
  type Icon,
} from "@tabler/icons-react"

export type NavKey =
  | "dashboard"
  | "my-tickets"
  | "pending-review"
  | "tickets"
  | "statements"
  | "crons"
  | "clusters"
  | "databases"
  | "users"
  | "rules"
  | "violations"
  | "logs"
  | "options"
  | "stats"

export interface NavEntry {
  key: NavKey
  label: string
  icon: Icon
  badge?: number
}

export interface NavSection {
  title: string
  items: NavEntry[]
}

/** 侧边栏分组：与后端能力一一对应，徽标数字来自假数据。 */
export const navSections: NavSection[] = [
  {
    title: "工作台",
    items: [
      { key: "dashboard", label: "概览", icon: IconLayoutDashboard },
      { key: "my-tickets", label: "我的工单", icon: IconTicket, badge: 7 },
      { key: "pending-review", label: "待我审核", icon: IconListCheck, badge: 12 },
    ],
  },
  {
    title: "变更管理",
    items: [
      { key: "tickets", label: "工单列表", icon: IconFileText },
      { key: "statements", label: "语句管理", icon: IconDatabase },
      { key: "crons", label: "定时任务", icon: IconCalendarTime },
    ],
  },
  {
    title: "资源",
    items: [
      { key: "clusters", label: "群集", icon: IconServer, badge: 8 },
      { key: "databases", label: "数据库", icon: IconDatabase },
      { key: "users", label: "用户与角色", icon: IconUsers },
    ],
  },
  {
    title: "规则与审计",
    items: [
      { key: "rules", label: "规则库", icon: IconShieldCheck, badge: 248 },
      { key: "violations", label: "命中记录", icon: IconAlertTriangle },
      { key: "logs", label: "操作日志", icon: IconHistory },
    ],
  },
  {
    title: "系统",
    items: [
      { key: "options", label: "系统选项", icon: IconSettings },
      { key: "stats", label: "统计报表", icon: IconChartBar },
    ],
  },
]
