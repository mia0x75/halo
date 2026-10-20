import {
  IconChevronsLeft,
  IconChevronsRight,
  IconDatabase,
  IconLogout,
  IconSettings,
  IconUserCircle,
} from "@tabler/icons-react"

import { navSections, type NavKey } from "@/components/layout/nav-items"
import { Avatar } from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { cn } from "@/lib/utils"

export interface SideNavProps {
  active: NavKey
  onSelect: (key: NavKey) => void
  collapsed: boolean
  onToggle: () => void
}

/** 侧边栏：15rem 宽（收起后 4rem 只留图标），顶部是当前用户，底部是帮助入口。 */
export function SideNav({ active, onSelect, collapsed, onToggle }: SideNavProps) {
  return (
    <aside
      className={cn(
        "flex h-svh shrink-0 flex-col border-r border-border bg-card transition-[width] duration-200",
        collapsed ? "w-16" : "w-60",
      )}
    >
      <div
        className={cn(
          "flex h-14 items-center gap-2 border-b border-border",
          collapsed ? "justify-center px-0" : "px-4",
        )}
      >
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button
              type="button"
              className={cn(
                "flex min-w-0 items-center gap-2 rounded-md px-1 py-1 text-left transition-colors hover:bg-surface-secondary",
                collapsed && "px-0",
              )}
              aria-label="当前用户"
            >
              <Avatar label="张" size="md" />
              {collapsed ? null : (
                <span className="flex min-w-0 flex-col leading-tight">
                  <span className="truncate text-sm font-medium">张三</span>
                  <span className="truncate text-sm text-muted-foreground">管理员 · 核心交易</span>
                </span>
              )}
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" className="w-48">
            <DropdownMenuLabel>张三 · 管理员</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <IconUserCircle size={16} stroke={1.5} />
              个人资料
            </DropdownMenuItem>
            <DropdownMenuItem>
              <IconSettings size={16} stroke={1.5} />
              账号设置
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <IconLogout size={16} stroke={1.5} />
              退出登录
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        {collapsed ? null : (
          <button
            type="button"
            onClick={onToggle}
            aria-label="收起导航栏"
            className="ml-auto flex size-8 shrink-0 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-surface-secondary hover:text-foreground"
          >
            <IconChevronsLeft size={16} stroke={1.5} />
          </button>
        )}
      </div>

      {collapsed ? (
        <button
          type="button"
          onClick={onToggle}
          aria-label="展开导航栏"
          className="mx-auto mt-2 flex size-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-surface-secondary hover:text-foreground"
        >
          <IconChevronsRight size={16} stroke={1.5} />
        </button>
      ) : null}

      <nav className="flex-1 overflow-y-auto px-2 py-3" aria-label="主导航">
        {navSections.map((section) => (
          <div key={section.title} className="mb-4 last:mb-0">
            {collapsed ? null : (
              <p className="mt-4 px-2 pt-1 pb-1 text-xs font-medium text-muted-foreground first:mt-0">
                {section.title}
              </p>
            )}
            <ul className="flex flex-col gap-0.5">
              {section.items.map((item) => {
                const isActive = item.key === active
                const IconComponent = item.icon
                return (
                  <li key={item.key}>
                    <button
                      type="button"
                      title={item.label}
                      aria-current={isActive ? "page" : undefined}
                      onClick={() => onSelect(item.key)}
                      className={cn(
                        "flex w-full items-center gap-2 rounded-md px-4 py-2 text-sm transition-colors",
                        collapsed && "justify-center px-0",
                        isActive
                          ? "bg-primary font-medium text-primary-foreground"
                          : "text-foreground hover:bg-surface-secondary",
                      )}
                    >
                      <IconComponent size={20} stroke={1.5} className="shrink-0" />
                      {collapsed ? null : <span className="truncate">{item.label}</span>}
                      {item.badge !== undefined && !collapsed ? (
                        <span
                          className={cn(
                            "ml-auto rounded-md px-1 py-0.5 text-[0.625rem] font-medium",
                            isActive
                              ? "bg-primary-foreground/20 text-primary-foreground"
                              : "bg-surface-secondary text-muted-foreground",
                          )}
                        >
                          {item.badge}
                        </span>
                      ) : null}
                    </button>
                  </li>
                )
              })}
            </ul>
          </div>
        ))}
      </nav>

      {collapsed ? null : (
        <div className="border-t border-border p-2">
          <button
            type="button"
            className="flex w-full items-center gap-2 rounded-md px-2 py-2 text-sm text-muted-foreground transition-colors hover:bg-surface-secondary hover:text-foreground"
          >
            <IconDatabase size={20} stroke={1.5} />
            帮助 / 文档
          </button>
        </div>
      )}
    </aside>
  )
}
