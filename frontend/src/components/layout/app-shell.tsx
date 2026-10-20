import { useState, type ReactNode } from "react"

import { SideNav } from "@/components/layout/side-nav"
import { TopNav } from "@/components/layout/top-nav"
import type { NavKey } from "@/components/layout/nav-items"

export interface AppShellProps {
  active: NavKey
  onSelect: (key: NavKey) => void
  onConnectionChange: () => void
  children: ReactNode
}

/** 应用外壳：左侧可收起导航 + 顶栏 + 内容区，页脚通栏。 */
export function AppShell({ active, onSelect, onConnectionChange, children }: AppShellProps) {
  const [collapsed, setCollapsed] = useState(false)

  return (
    <div className="flex min-h-svh flex-col bg-background text-foreground">
      <div className="flex min-h-0 flex-1">
        <SideNav
          active={active}
          onSelect={onSelect}
          collapsed={collapsed}
          onToggle={() => setCollapsed((value) => !value)}
        />
        <div className="flex min-w-0 flex-1 flex-col">
          <TopNav onConnectionChange={onConnectionChange} />
          <main className="flex-1">
            <div className="mx-auto w-full max-w-[1320px] px-4 py-6">
              <div className="flex flex-col gap-4">{children}</div>
            </div>
          </main>
        </div>
      </div>

      <footer className="border-t border-border bg-card">
        <div
          className="flex items-center justify-between py-3 pr-4 text-sm text-muted-foreground"
          // 页脚通栏，但内容与右侧内容区左对齐：左侧留出导航栏的宽度。
          style={{ paddingLeft: (collapsed ? 64 : 240) + 16 }}
        >
          <span>Halo · SQL 变更工单与审计平台</span>
          <span>v0.1 · 演示数据</span>
        </div>
      </footer>
    </div>
  )
}
