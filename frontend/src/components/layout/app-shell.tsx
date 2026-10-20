import type { ReactNode } from "react"

import { SideNav } from "@/components/layout/side-nav"
import { TopNav } from "@/components/layout/top-nav"
import type { NavKey } from "@/components/layout/nav-items"

export interface AppShellProps {
  active: NavKey
  onSelect: (key: NavKey) => void
  children: ReactNode
}

/** 应用外壳：左侧固定导航 + 顶栏 + 内容区。 */
export function AppShell({ active, onSelect, children }: AppShellProps) {
  return (
    <div className="flex min-h-svh bg-background text-foreground">
      <SideNav active={active} onSelect={onSelect} />
      <div className="flex min-w-0 flex-1 flex-col">
        <TopNav />
        <main className="flex-1">
          <div className="mx-auto w-full max-w-[1320px] px-4 py-6">
            <div className="flex flex-col gap-4">{children}</div>
          </div>
        </main>
        <footer className="border-t border-border bg-card">
          <div className="mx-auto flex max-w-[1320px] items-center justify-between px-4 py-3 text-xs text-muted-foreground">
            <span>Halo · SQL 变更工单与审计平台</span>
            <span>v0.1 · 演示数据</span>
          </div>
        </footer>
      </div>
    </div>
  )
}
