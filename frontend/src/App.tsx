import { useCallback, useEffect, useState } from "react"
import { AppShell } from "@/components/layout/app-shell"
import type { NavKey } from "@/components/layout/nav-items"
import { ClusterFormPage } from "@/pages/cluster-form-page"
import { DatabasesPage } from "@/pages/databases-page"
import { CronCreatePage } from "@/pages/cron-create-page"
import { CronsPage } from "@/pages/crons-page"
import { ClustersPage } from "@/pages/clusters-page"
import { DashboardPage } from "@/pages/dashboard-page"
import { LogsPage } from "@/pages/logs-page"
import { MessagesPage } from "@/pages/messages-page"
import { SettingsPage, type SettingsTab } from "@/pages/settings-page"
import { GroupFormPage } from "@/pages/group-form-page"
import { StatsPage } from "@/pages/stats-page"
import { RulesPage } from "@/pages/rules-page"
import { TicketCreatePage } from "@/pages/ticket-create-page"
import { UsersPage } from "@/pages/users-page"
import { UserFormPage } from "@/pages/user-form-page"
import { TicketDetailPage } from "@/pages/ticket-detail-page"
import { TicketRejectPage } from "@/pages/ticket-reject-page"
import { TicketEditPage } from "@/pages/ticket-edit-page"
import { TicketQueuePage } from "@/pages/ticket-queue-page"
import { TicketsPage } from "@/pages/tickets-page"
import { ViolationsPage } from "@/pages/violations-page"
import { LoginPage } from "@/pages/login-page"
import { WelcomePage } from "@/pages/welcome-page"
import { useSession } from "@/hooks/use-session"
import { clearToken } from "@/lib/api/session"
import { fetchSetupStatus } from "@/lib/api/setup"

type View =
  | { kind: "nav"; key: NavKey }
  // Ticket screens remember the list they were opened from, so the sidebar
  // keeps highlighting it and 返回列表 goes back there instead of always
  // landing on the full ticket list.
  | { kind: "ticket"; uuid: string; from?: NavKey }
  | { kind: "ticket-create"; from?: NavKey }
  | { kind: "ticket-edit"; uuid: string; from?: NavKey }
  | { kind: "cron-create" }
  | { kind: "cluster-form"; mode: "create" | "edit"; uuid?: string }
  | { kind: "group-form"; mode: "create" | "edit"; uuid?: string }
  | { kind: "user-form"; mode: "create" | "edit"; uuid?: string }
  | { kind: "ticket-reject"; uuid: string; from?: NavKey }

type SetupState =
  | { status: "checking" }
  | { status: "ready" }
  | { status: "needed"; configFile: string }

export default function App() {
  const session = useSession()
  const [setup, setSetup] = useState<SetupState>({ status: "checking" })
  const [view, setView] = useState<View>({ kind: "nav", key: "dashboard" })
  const [searchKeyword, setSearchKeyword] = useState("")
  const [settingsTab, setSettingsTab] = useState<SettingsTab>("options")

  // 服务器没配置元数据库时先走安装向导。拿不到状态就按已配置处理，
  // 免得演示环境与旧版本后端被挡住。
  useEffect(() => {
    let cancelled = false
    fetchSetupStatus()
      .then((status) => {
        if (cancelled) return
        setSetup(
          status.configured
            ? { status: "ready" }
            : { status: "needed", configFile: status.configFile },
        )
      })
      .catch(() => {
        if (!cancelled) setSetup({ status: "ready" })
      })
    return () => {
      cancelled = true
    }
  }, [])

  const signOut = useCallback(() => {
    clearToken()
    setView({ kind: "nav", key: "dashboard" })
  }, [])

  const openTicket = useCallback(
    (uuid: string, from?: NavKey) => setView({ kind: "ticket", uuid, from }),
    [],
  )
  const createTicket = useCallback((from?: NavKey) => setView({ kind: "ticket-create", from }), [])
  const createCluster = useCallback(() => setView({ kind: "cluster-form", mode: "create" }), [])
  const editCluster = useCallback(
    (uuid: string) => setView({ kind: "cluster-form", mode: "edit", uuid }),
    [],
  )
  const searchTickets = useCallback((keyword: string) => {
    setSearchKeyword(keyword)
    setView({ kind: "nav", key: "tickets" })
  }, [])
  const viewTickets = useCallback(() => setView({ kind: "nav", key: "tickets" }), [])
  const viewMessages = useCallback(() => setView({ kind: "nav", key: "messages" }), [])
  const createUser = useCallback(() => setView({ kind: "user-form", mode: "create" }), [])
  const editUser = useCallback(
    (uuid: string) => setView({ kind: "user-form", mode: "edit", uuid }),
    [],
  )
  const createGroup = useCallback(() => {
    setSettingsTab("groups")
    setView({ kind: "group-form", mode: "create" })
  }, [])
  const editGroup = useCallback((uuid: string) => {
    setSettingsTab("groups")
    setView({ kind: "group-form", mode: "edit", uuid })
  }, [])

  // Ticket screens carry the list they came from; the other detail screens
  // keep highlighting their own section.
  const page: NavKey =
    view.kind === "nav"
      ? view.key
      : view.kind === "ticket" ||
          view.kind === "ticket-create" ||
          view.kind === "ticket-edit" ||
          view.kind === "ticket-reject"
        ? (view.from ?? "tickets")
        : "tickets"

  if (setup.status === "checking") {
    return (
      <div className="flex min-h-svh items-center justify-center text-sm text-muted-foreground">
        正在加载…
      </div>
    )
  }

  if (setup.status === "needed") {
    return <WelcomePage onDone={() => window.location.reload()} />
  }

  if (session.status !== "authenticated") {
    return <LoginPage />
  }

  if (view.kind === "cron-create") {
    return (
      <AppShell
        active="crons"
        onSelect={(key) => setView({ kind: "nav", key })}
        user={session.user}
        onSignOut={signOut}
        onSearch={searchTickets}
        onViewMessages={viewMessages}
        onOpenTicket={openTicket}
      >
        <CronCreatePage
          onBack={() => setView({ kind: "nav", key: "crons" })}
          onSaved={() => setView({ kind: "nav", key: "crons" })}
        />
      </AppShell>
    )
  }

  if (view.kind === "ticket-edit") {
    return (
      <AppShell
        active={page}
        onSelect={(key) => setView({ kind: "nav", key })}
        user={session.user}
        onSignOut={signOut}
        onSearch={searchTickets}
        onViewMessages={viewMessages}
        onOpenTicket={openTicket}
      >
        <TicketEditPage
          uuid={view.uuid}
          onBack={() => setView({ kind: "ticket", uuid: view.uuid, from: view.from })}
          onSaved={(uuid) => setView({ kind: "ticket", uuid, from: view.from })}
        />
      </AppShell>
    )
  }

  if (view.kind === "ticket-create") {
    return (
      <AppShell
        active={page}
        onSelect={(key) => setView({ kind: "nav", key })}
        user={session.user}
        onSignOut={signOut}
        onSearch={searchTickets}
        onViewMessages={viewMessages}
        onOpenTicket={openTicket}
      >
        <TicketCreatePage
          onBack={() => setView({ kind: "nav", key: page })}
          onCreated={(uuid) => {
            setView({ kind: "ticket", uuid, from: view.from })
          }}
        />
      </AppShell>
    )
  }

  if (view.kind === "cluster-form") {
    return (
      <AppShell
        active="clusters"
        onSelect={(key) => setView({ kind: "nav", key })}
        user={session.user}
        onSignOut={signOut}
        onSearch={searchTickets}
        onViewMessages={viewMessages}
        onOpenTicket={openTicket}
      >
        <ClusterFormPage
          mode={view.mode}
          uuid={view.uuid}
          onBack={() => setView({ kind: "nav", key: "clusters" })}
          onSaved={() => setView({ kind: "nav", key: "clusters" })}
        />
      </AppShell>
    )
  }

  if (view.kind === "ticket-reject") {
    return (
      <AppShell
        active={page}
        onSelect={(key) => setView({ kind: "nav", key })}
        user={session.user}
        onSignOut={signOut}
        onSearch={searchTickets}
        onViewMessages={viewMessages}
        onOpenTicket={openTicket}
      >
        <TicketRejectPage
          uuid={view.uuid}
          onBack={() => setView({ kind: "ticket", uuid: view.uuid, from: view.from })}
          onDone={() => setView({ kind: "ticket", uuid: view.uuid, from: view.from })}
        />
      </AppShell>
    )
  }

  if (view.kind === "user-form") {
    return (
      <AppShell
        active="users"
        onSelect={(key) => setView({ kind: "nav", key })}
        user={session.user}
        onSignOut={signOut}
        onSearch={searchTickets}
        onViewMessages={viewMessages}
        onOpenTicket={openTicket}
      >
        <UserFormPage
          mode={view.mode}
          uuid={view.uuid}
          onBack={() => setView({ kind: "nav", key: "users" })}
          onSaved={() => setView({ kind: "nav", key: "users" })}
        />
      </AppShell>
    )
  }

  if (view.kind === "group-form") {
    return (
      <AppShell
        active="settings"
        onSelect={(key) => setView({ kind: "nav", key })}
        user={session.user}
        onSignOut={signOut}
        onSearch={searchTickets}
        onViewMessages={viewMessages}
        onOpenTicket={openTicket}
      >
        <GroupFormPage
          mode={view.mode}
          uuid={view.uuid}
          onBack={() => setView({ kind: "nav", key: "settings" })}
          onSaved={() => setView({ kind: "nav", key: "settings" })}
        />
      </AppShell>
    )
  }

  if (view.kind === "ticket") {
    return (
      <AppShell
        active={page}
        onSelect={(key) => setView({ kind: "nav", key })}
        user={session.user}
        onSignOut={signOut}
        onSearch={searchTickets}
        onViewMessages={viewMessages}
        onOpenTicket={openTicket}
      >
        <TicketDetailPage
          uuid={view.uuid}
          onBack={() => setView({ kind: "nav", key: page })}
          onEdit={() => setView({ kind: "ticket-edit", uuid: view.uuid, from: view.from })}
          viewerUUID={session.user.UUID}
          onReject={(uuid) => setView({ kind: "ticket-reject", uuid, from: view.from })}
        />
      </AppShell>
    )
  }

  return (
    <AppShell
      active={page}
      onSelect={(key) => setView({ kind: "nav", key })}
      user={session.user}
      onSignOut={signOut}
      onSearch={searchTickets}
      onViewMessages={viewMessages}
      onOpenTicket={(uuid) => openTicket(uuid, page)}
    >
      {page === "dashboard" ? (
        <DashboardPage
          onOpenTicket={(uuid) => openTicket(uuid, "dashboard")}
          onCreateTicket={() => createTicket("dashboard")}
          onViewTickets={viewTickets}
        />
      ) : null}
      {page === "messages" ? (
        <MessagesPage onOpenTicket={(uuid) => openTicket(uuid, "messages")} />
      ) : null}
      {page === "tickets" ? (
        <TicketsPage
          key={searchKeyword}
          initialKeyword={searchKeyword}
          onOpenTicket={(uuid) => openTicket(uuid, "tickets")}
          onCreateTicket={() => createTicket("tickets")}
        />
      ) : null}
      {page === "clusters" ? (
        <ClustersPage onCreate={createCluster} onEdit={editCluster} />
      ) : null}
      {page === "rules" ? <RulesPage /> : null}
      {page === "violations" ? <ViolationsPage /> : null}
      {page === "logs" ? <LogsPage /> : null}
      {page === "my-tickets" ? (
        <TicketQueuePage
          mode="mine"
          onOpenTicket={(uuid) => openTicket(uuid, "my-tickets")}
          onCreateTicket={() => createTicket("my-tickets")}
        />
      ) : null}
      {page === "pending-review" ? (
        <TicketQueuePage
          mode="review"
          onOpenTicket={(uuid) => openTicket(uuid, "pending-review")}
        />
      ) : null}
      {page === "crons" ? (
        <CronsPage onCreate={() => setView({ kind: "cron-create" })} />
      ) : null}
      {page === "databases" ? <DatabasesPage /> : null}
      {page === "users" ? (
        <UsersPage onCreateUser={createUser} onEditUser={editUser} />
      ) : null}
      {page === "settings" ? (
        <SettingsPage
          tab={settingsTab}
          onTabChange={setSettingsTab}
          onCreateGroup={createGroup}
          onEditGroup={editGroup}
        />
      ) : null}
      {page === "stats" ? <StatsPage /> : null}
    </AppShell>
  )
}
