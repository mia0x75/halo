import { useCallback, useState } from "react"
import {
  IconCalendarTime,
  IconChartBar,
  IconDatabase,
  IconSettings,
  IconUsers,
} from "@tabler/icons-react"

import { AppShell } from "@/components/layout/app-shell"
import type { NavKey } from "@/components/layout/nav-items"
import { ClusterFormPage } from "@/pages/cluster-form-page"
import { ClustersPage } from "@/pages/clusters-page"
import { DashboardPage } from "@/pages/dashboard-page"
import { LogsPage } from "@/pages/logs-page"
import { PlaceholderPage } from "@/pages/placeholder-page"
import { RulesPage } from "@/pages/rules-page"
import { TicketCreatePage } from "@/pages/ticket-create-page"
import { TicketDetailPage } from "@/pages/ticket-detail-page"
import { TicketQueuePage } from "@/pages/ticket-queue-page"
import { TicketsPage } from "@/pages/tickets-page"
import { ViolationsPage } from "@/pages/violations-page"

type View =
  | { kind: "nav"; key: NavKey }
  | { kind: "ticket"; uuid: string }
  | { kind: "ticket-create" }
  | { kind: "cluster-form"; mode: "create" | "edit"; uuid?: string }

export default function App() {
  const [view, setView] = useState<View>({ kind: "nav", key: "dashboard" })
  const [nonce, setNonce] = useState(0)

  const refresh = useCallback(() => setNonce((value) => value + 1), [])
  const openTicket = useCallback((uuid: string) => setView({ kind: "ticket", uuid }), [])
  const createTicket = useCallback(() => setView({ kind: "ticket-create" }), [])
  const createCluster = useCallback(() => setView({ kind: "cluster-form", mode: "create" }), [])
  const editCluster = useCallback(
    (uuid: string) => setView({ kind: "cluster-form", mode: "edit", uuid }),
    [],
  )

  const page = view.kind === "nav" ? view.key : "tickets"

  if (view.kind === "ticket-create") {
    return (
      <AppShell
        key={nonce}
        active="tickets"
        onSelect={(key) => setView({ kind: "nav", key })}
        onConnectionChange={refresh}
      >
        <TicketCreatePage
          onBack={() => setView({ kind: "nav", key: "tickets" })}
          onCreated={(uuid) => {
            setView({ kind: "ticket", uuid })
          }}
        />
      </AppShell>
    )
  }

  if (view.kind === "cluster-form") {
    return (
      <AppShell
        key={nonce}
        active="clusters"
        onSelect={(key) => setView({ kind: "nav", key })}
        onConnectionChange={refresh}
      >
        <ClusterFormPage
          mode={view.mode}
          uuid={view.uuid}
          onBack={() => setView({ kind: "nav", key: "clusters" })}
          onSaved={() => {
            setView({ kind: "nav", key: "clusters" })
            refresh()
          }}
        />
      </AppShell>
    )
  }

  if (view.kind === "ticket") {
    return (
      <AppShell
        key={nonce}
        active="tickets"
        onSelect={(key) => setView({ kind: "nav", key })}
        onConnectionChange={refresh}
      >
        <TicketDetailPage
          uuid={view.uuid}
          onBack={() => setView({ kind: "nav", key: "tickets" })}
        />
      </AppShell>
    )
  }

  return (
    <AppShell
      key={nonce}
      active={page}
      onSelect={(key) => setView({ kind: "nav", key })}
      onConnectionChange={refresh}
    >
      {page === "dashboard" ? (
        <DashboardPage onOpenTicket={openTicket} onCreateTicket={createTicket} />
      ) : null}
      {page === "tickets" ? (
        <TicketsPage onOpenTicket={openTicket} onCreateTicket={createTicket} />
      ) : null}
      {page === "clusters" ? (
        <ClustersPage onCreate={createCluster} onEdit={editCluster} />
      ) : null}
      {page === "rules" ? <RulesPage /> : null}
      {page === "violations" ? <ViolationsPage /> : null}
      {page === "logs" ? <LogsPage /> : null}
      {page === "my-tickets" ? <TicketQueuePage mode="mine" onOpenTicket={openTicket} /> : null}
      {page === "pending-review" ? <TicketQueuePage mode="review" onOpenTicket={openTicket} /> : null}
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
