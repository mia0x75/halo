import { cn } from "@/lib/utils"

export interface TabItem {
  key: string
  label: string
  badge?: number
}

export interface TabsProps {
  items: TabItem[]
  value: string
  onChange: (key: string) => void
  className?: string
}

/** 标签页：选中项用 2px 主色下划线标记。 */
export function Tabs({ items, value, onChange, className }: TabsProps) {
  return (
    <div className={cn("flex items-center gap-4 border-b border-border", className)} role="tablist">
      {items.map((item) => {
        const active = item.key === value
        return (
          <button
            key={item.key}
            type="button"
            role="tab"
            aria-selected={active}
            onClick={() => onChange(item.key)}
            className={cn(
              "-mb-px flex items-center gap-2 border-b-2 px-1 pb-3 text-sm transition-colors",
              active
                ? "border-primary font-medium text-foreground"
                : "border-transparent text-muted-foreground hover:text-foreground",
            )}
          >
            {item.label}
            {item.badge !== undefined ? (
              <span className="rounded-md bg-surface-secondary px-1.5 py-0.5 text-[0.625rem] font-medium text-muted-foreground">
                {item.badge}
              </span>
            ) : null}
          </button>
        )
      })}
    </div>
  )
}
