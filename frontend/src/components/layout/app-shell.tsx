import type { ReactNode } from "react"

import { TopNav, type NavKey } from "@/components/layout/top-nav"

export interface AppShellProps {
  active: NavKey
  onSelect: (key: NavKey) => void
  children: ReactNode
}

export function AppShell({ active, onSelect, children }: AppShellProps) {
  return (
    <div className="flex min-h-svh flex-col bg-background text-foreground">
      <TopNav active={active} onSelect={onSelect} />
      <main className="mx-auto w-full max-w-[1320px] flex-1 px-4 py-6">
        <div className="space-y-4">{children}</div>
      </main>
      <footer className="border-t border-border bg-card">
        <div className="mx-auto flex max-w-[1320px] items-center justify-between px-4 py-3 text-xs text-muted-foreground">
          <span>Halo · SQL 变更工单与审计平台</span>
          <span>v0.1</span>
        </div>
      </footer>
    </div>
  )
}
