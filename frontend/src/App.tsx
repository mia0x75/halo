import { useCallback, useState } from "react"
import { AppShell } from "@/components/layout/app-shell"
import type { NavKey } from "@/components/layout/nav-items"
import { ClusterFormPage } from "@/pages/cluster-form-page"
import { DatabasesPage } from "@/pages/databases-page"
import { CronCreatePage } from "@/pages/cron-create-page"
import { CronsPage } from "@/pages/crons-page"
import { ClustersPage } from "@/pages/clusters-page"
import { DashboardPage } from "@/pages/dashboard-page"
import { LogsPage } from "@/pages/logs-page"
import { OptionsPage } from "@/pages/options-page"
import { StatementsPage } from "@/pages/statements-page"
import { StatsPage } from "@/pages/stats-page"
import { RulesPage } from "@/pages/rules-page"
import { TicketCreatePage } from "@/pages/ticket-create-page"
import { UsersPage } from "@/pages/users-page"
import { TicketDetailPage } from "@/pages/ticket-detail-page"
import { TicketEditPage } from "@/pages/ticket-edit-page"
import { TicketQueuePage } from "@/pages/ticket-queue-page"
import { TicketsPage } from "@/pages/tickets-page"
import { ViolationsPage } from "@/pages/violations-page"
import { LoginPage } from "@/pages/login-page"
import { useSession } from "@/hooks/use-session"
import { clearToken } from "@/lib/api/session"

type View =
  | { kind: "nav"; key: NavKey }
  | { kind: "ticket"; uuid: string }
  | { kind: "ticket-create" }
  | { kind: "ticket-edit"; uuid: string }
  | { kind: "cron-create" }
  | { kind: "cluster-form"; mode: "create" | "edit"; uuid?: string }

export default function App() {
  const session = useSession()
  const [view, setView] = useState<View>({ kind: "nav", key: "dashboard" })
  const [searchKeyword, setSearchKeyword] = useState("")

  const signOut = useCallback(() => {
    clearToken()
    setView({ kind: "nav", key: "dashboard" })
  }, [])

  const openTicket = useCallback((uuid: string) => setView({ kind: "ticket", uuid }), [])
  const createTicket = useCallback(() => setView({ kind: "ticket-create" }), [])
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

  const page = view.kind === "nav" ? view.key : "tickets"

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
        active="tickets"
        onSelect={(key) => setView({ kind: "nav", key })}
        user={session.user}
        onSignOut={signOut}
        onSearch={searchTickets}
      >
        <TicketEditPage
          uuid={view.uuid}
          onBack={() => setView({ kind: "ticket", uuid: view.uuid })}
          onSaved={(uuid) => setView({ kind: "ticket", uuid })}
        />
      </AppShell>
    )
  }

  if (view.kind === "ticket-create") {
    return (
      <AppShell
        active="tickets"
        onSelect={(key) => setView({ kind: "nav", key })}
        user={session.user}
        onSignOut={signOut}
        onSearch={searchTickets}
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
        active="clusters"
        onSelect={(key) => setView({ kind: "nav", key })}
        user={session.user}
        onSignOut={signOut}
        onSearch={searchTickets}
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

  if (view.kind === "ticket") {
    return (
      <AppShell
        active="tickets"
        onSelect={(key) => setView({ kind: "nav", key })}
        user={session.user}
        onSignOut={signOut}
        onSearch={searchTickets}
      >
        <TicketDetailPage
          uuid={view.uuid}
          onBack={() => setView({ kind: "nav", key: "tickets" })}
          onEdit={() => setView({ kind: "ticket-edit", uuid: view.uuid })}
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
    >
      {page === "dashboard" ? (
        <DashboardPage
          onOpenTicket={openTicket}
          onCreateTicket={createTicket}
          onViewTickets={viewTickets}
        />
      ) : null}
      {page === "tickets" ? (
        <TicketsPage
          key={searchKeyword}
          initialKeyword={searchKeyword}
          onOpenTicket={openTicket}
          onCreateTicket={createTicket}
        />
      ) : null}
      {page === "clusters" ? (
        <ClustersPage onCreate={createCluster} onEdit={editCluster} />
      ) : null}
      {page === "rules" ? <RulesPage /> : null}
      {page === "violations" ? <ViolationsPage /> : null}
      {page === "logs" ? <LogsPage /> : null}
      {page === "my-tickets" ? (
        <TicketQueuePage mode="mine" onOpenTicket={openTicket} onCreateTicket={createTicket} />
      ) : null}
      {page === "pending-review" ? (
        <TicketQueuePage mode="review" onOpenTicket={openTicket} />
      ) : null}
      {page === "statements" ? <StatementsPage /> : null}
      {page === "crons" ? (
        <CronsPage onCreate={() => setView({ kind: "cron-create" })} />
      ) : null}
      {page === "databases" ? <DatabasesPage /> : null}
      {page === "users" ? <UsersPage /> : null}
      {page === "options" ? <OptionsPage /> : null}
      {page === "stats" ? <StatsPage /> : null}
    </AppShell>
  )
}
