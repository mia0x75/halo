import { IconDatabase, IconShieldCheck } from "@tabler/icons-react"

import { ThemeToggle } from "@/components/layout/theme-toggle"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

export type NavKey = "dashboard" | "tickets"

interface NavItem {
  key: NavKey
  label: string
}

const navItems: NavItem[] = [
  { key: "dashboard", label: "概览" },
  { key: "tickets", label: "工单" },
]

export interface TopNavProps {
  active: NavKey
  onSelect: (key: NavKey) => void
}

export function TopNav({ active, onSelect }: TopNavProps) {
  return (
    <header className="sticky top-0 z-20 border-b border-border bg-card/95 backdrop-blur">
      <div className="mx-auto flex h-14 max-w-[1320px] items-center gap-6 px-4">
        <div className="flex items-center gap-2">
          <span className="flex size-7 items-center justify-center rounded-md bg-primary text-primary-foreground">
            <IconDatabase size={16} stroke={1.5} />
          </span>
          <span className="text-sm font-semibold tracking-tight">Halo</span>
        </div>

        <nav className="flex items-center gap-1" aria-label="主导航">
          {navItems.map((item) => (
            <Button
              key={item.key}
              variant="ghost"
              size="sm"
              aria-current={active === item.key ? "page" : undefined}
              onClick={() => onSelect(item.key)}
              className={cn(
                "text-sm",
                active === item.key ? "bg-accent text-accent-foreground" : "text-muted-foreground",
              )}
            >
              {item.label}
            </Button>
          ))}
        </nav>

        <div className="ml-auto flex items-center gap-1">
          <Button variant="ghost" size="sm" aria-label="规则与审计">
            <IconShieldCheck size={16} stroke={1.5} />
            <span className="hidden sm:inline">规则与审计</span>
          </Button>
          <ThemeToggle />
        </div>
      </div>
    </header>
  )
}
