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

/**
 * Purpose-built shell for table/chart cards. shadcn cards drive spacing from a
 * CSS variable, which is fragile to override with utilities (see the decision
 * record §10), so dense cards use this explicit wrapper instead.
 */
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
      <div className="flex min-h-14 items-center justify-between gap-2 border-b border-border px-5">
        <div className="flex flex-col gap-0.5">
          <h3 className="text-base font-semibold tracking-tight">{title}</h3>
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
