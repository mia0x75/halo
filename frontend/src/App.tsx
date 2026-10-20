import { useState } from "react"
import {
  IconCalendarTime,
  IconChartBar,
  IconDatabase,
  IconListCheck,
  IconSettings,
  IconTicket,
  IconUsers,
} from "@tabler/icons-react"

import { AppShell } from "@/components/layout/app-shell"
import type { NavKey } from "@/components/layout/nav-items"
import { ClustersPage } from "@/pages/clusters-page"
import { DashboardPage } from "@/pages/dashboard-page"
import { LogsPage } from "@/pages/logs-page"
import { PlaceholderPage } from "@/pages/placeholder-page"
import { RulesPage } from "@/pages/rules-page"
import { TicketsPage } from "@/pages/tickets-page"
import { ViolationsPage } from "@/pages/violations-page"

export default function App() {
  const [page, setPage] = useState<NavKey>("dashboard")

  return (
    <AppShell active={page} onSelect={setPage}>
      {page === "dashboard" ? <DashboardPage /> : null}
      {page === "tickets" ? <TicketsPage /> : null}
      {page === "clusters" ? <ClustersPage /> : null}
      {page === "rules" ? <RulesPage /> : null}
      {page === "violations" ? <ViolationsPage /> : null}
      {page === "logs" ? <LogsPage /> : null}
      {page === "my-tickets" ? (
        <PlaceholderPage
          title="我的工单"
          description="我提交与参与审核的工单"
          icon={IconTicket}
        />
      ) : null}
      {page === "pending-review" ? (
        <PlaceholderPage
          title="待我审核"
          description="等待我处理审核的工单"
          icon={IconListCheck}
        />
      ) : null}
      {page === "statements" ? (
        <PlaceholderPage
          title="语句管理"
          description="工单下的 SQL 语句与执行计划"
          icon={IconDatabase}
        />
      ) : null}
      {page === "crons" ? (
        <PlaceholderPage
          title="定时任务"
          description="系统作业与工单预约执行"
          icon={IconCalendarTime}
        />
      ) : null}
      {page === "databases" ? (
        <PlaceholderPage
          title="数据库"
          description="目标库清单与字符集"
          icon={IconDatabase}
        />
      ) : null}
      {page === "users" ? (
        <PlaceholderPage
          title="用户与角色"
          description="账号、角色与分组授权"
          icon={IconUsers}
        />
      ) : null}
      {page === "options" ? (
        <PlaceholderPage
          title="系统选项"
          description="平台级配置项"
          icon={IconSettings}
        />
      ) : null}
      {page === "stats" ? (
        <PlaceholderPage
          title="统计报表"
          description="按天汇总的工单与命中指标"
          icon={IconChartBar}
        />
      ) : null}
    </AppShell>
  )
}
