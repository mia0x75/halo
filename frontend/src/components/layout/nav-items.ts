import {
  IconAlertTriangle,
  IconCalendarTime,
  IconChartBar,
  IconDatabase,
  IconFileText,
  IconHistory,
  IconInbox,
  IconLayoutDashboard,
  IconListCheck,
  IconRocket,
  IconServer,
  IconShieldLock,
  IconShieldCheck,
  IconTicket,
  IconUsers,
  type Icon,
} from "@tabler/icons-react"

export type NavKey =
  | "dashboard"
  | "messages"
  | "my-tickets"
  | "pending-review"
  | "tickets"
  | "execution"
  | "crons"
  | "clusters"
  | "databases"
  | "users"
  | "rules"
  | "violations"
  | "logs"
  | "settings"
  | "stats"

export interface NavEntry {
  key: NavKey
  label: string
  icon: Icon
  /** 仅管理员可见的入口 */
  adminOnly?: boolean
}

export interface NavSection {
  title: string
  items: NavEntry[]
}

/** 侧边栏分组：与后端能力一一对应，徽标由页面按接口总数实时提供。 */
export const navSections: NavSection[] = [
  {
    title: "工作台",
    items: [
      { key: "dashboard", label: "概览", icon: IconLayoutDashboard },
      { key: "messages", label: "我的消息", icon: IconInbox },
      { key: "my-tickets", label: "我的工单", icon: IconTicket },
      { key: "pending-review", label: "待我审核", icon: IconListCheck },
    ],
  },
  {
    title: "变更管理",
    items: [
      { key: "tickets", label: "工单列表", icon: IconFileText },
      { key: "execution", label: "执行计划", icon: IconRocket },
      { key: "crons", label: "定时任务", icon: IconCalendarTime },
    ],
  },
  {
    title: "资源",
    items: [
      { key: "clusters", label: "群集", icon: IconServer },
      { key: "databases", label: "数据库", icon: IconDatabase },
      { key: "users", label: "用户与角色", icon: IconUsers },
    ],
  },
  {
    title: "规则与审计",
    items: [
      { key: "rules", label: "规则库", icon: IconShieldCheck },
      { key: "violations", label: "命中记录", icon: IconAlertTriangle },
      { key: "logs", label: "操作日志", icon: IconHistory },
    ],
  },
  {
    title: "系统",
    items: [
      { key: "settings", label: "系统设置", icon: IconShieldLock, adminOnly: true },
      { key: "stats", label: "统计报表", icon: IconChartBar },
    ],
  },
]
