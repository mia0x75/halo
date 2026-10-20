import { useState } from "react"

import { AppShell } from "@/components/layout/app-shell"
import type { NavKey } from "@/components/layout/top-nav"
import { DashboardPage } from "@/pages/dashboard-page"
import { TicketsPage } from "@/pages/tickets-page"

export default function App() {
  const [page, setPage] = useState<NavKey>("dashboard")

  return (
    <AppShell active={page} onSelect={setPage}>
      {page === "dashboard" ? <DashboardPage /> : <TicketsPage />}
    </AppShell>
  )
}
