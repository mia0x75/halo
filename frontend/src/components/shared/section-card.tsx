import type { ReactNode } from "react"

import { cn } from "@/lib/utils"

export interface SectionCardProps {
  title: string
  description?: string
  action?: ReactNode
  children: ReactNode
  className?: string
  bodyClassName?: string
}

/** 区块卡片：头部 1.25rem 内边距 + 1px 边线，内容区可自定义内边距。 */
export function SectionCard({
  title,
  description,
  action,
  children,
  className,
  bodyClassName,
}: SectionCardProps) {
  return (
    <section
      className={cn(
        "flex flex-col overflow-hidden rounded-lg border border-border bg-card shadow-xs",
        className,
      )}
    >
      <div className="flex min-h-14 items-center justify-between gap-2 border-b border-border px-5 py-4">
        <div className="flex flex-col">
          <h3 className="text-base font-medium text-strong">{title}</h3>
          {description ? (
            <p className="text-xs text-muted-foreground">{description}</p>
          ) : null}
        </div>
        {action}
      </div>
      <div className={cn("p-5", bodyClassName)}>{children}</div>
    </section>
  )
}
