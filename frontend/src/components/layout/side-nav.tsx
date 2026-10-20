import { IconBrandGithub, IconChevronsLeft } from "@tabler/icons-react"

import { navSections, type NavKey } from "@/components/layout/nav-items"
import { cn } from "@/lib/utils"

export interface SideNavProps {
  active: NavKey
  onSelect: (key: NavKey) => void
}

/** 侧边栏：15rem 宽，白底、右侧 1px 边线，选中项为主色底。 */
export function SideNav({ active, onSelect }: SideNavProps) {
  return (
    <aside className="flex h-svh w-60 shrink-0 flex-col border-r border-border bg-card">
      <div className="flex h-14 items-center gap-2 border-b border-border px-4">
        <span className="flex size-8 items-center justify-center rounded-md bg-primary text-primary-foreground">
          <IconBrandGithub size={18} stroke={1.5} />
        </span>
        <div className="flex flex-col leading-tight">
          <span className="text-sm font-semibold">Halo</span>
          <span className="text-[0.625rem] text-muted-foreground">SQL 变更与审计</span>
        </div>
        <button
          type="button"
          aria-label="收起侧边栏"
          className="ml-auto flex size-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-surface-secondary hover:text-foreground"
        >
          <IconChevronsLeft size={16} stroke={1.5} />
        </button>
      </div>

      <nav className="flex-1 overflow-y-auto px-2 py-3" aria-label="主导航">
        {navSections.map((section) => (
          <div key={section.title} className="mb-4 last:mb-0">
            <p className="mt-4 px-2 pt-1 pb-1 text-xs font-medium text-muted-foreground first:mt-0">
              {section.title}
            </p>
            <ul className="flex flex-col gap-0.5">
              {section.items.map((item) => {
                const isActive = item.key === active
                const IconComponent = item.icon
                return (
                  <li key={item.key}>
                    <button
                      type="button"
                      aria-current={isActive ? "page" : undefined}
                      onClick={() => onSelect(item.key)}
                      className={cn(
                        "flex w-full items-center gap-2 rounded-md px-4 py-2 text-sm transition-colors",
                        isActive
                          ? "bg-primary font-medium text-primary-foreground"
                          : "text-foreground hover:bg-surface-secondary",
                      )}
                    >
                      <IconComponent size={20} stroke={1.5} />
                      <span className="truncate">{item.label}</span>
                      {item.badge !== undefined ? (
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
    </aside>
  )
}
