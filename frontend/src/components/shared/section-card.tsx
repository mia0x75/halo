import type { ReactNode } from "react"

import { cn } from "@/lib/utils"

export interface SectionCardProps {
  /** 区块标题。页面标题已经说明区块时可省略。 */
  title?: string
  /** 区块说明。 */
  description?: string
  action?: ReactNode
  children: ReactNode
  className?: string
  bodyClassName?: string
}

/**
 * 区块卡片：标题可选，内容区可自定义内边距。
 *
 * 页面标题已经说明区块时不必再写一遍标题，此时卡片只剩内容；没有标题也没有
 * 右侧操作时连头部都不渲染。
 */
export function SectionCard({
  title,
  description,
  action,
  children,
  className,
  bodyClassName,
}: SectionCardProps) {
  const hasHeading = title !== undefined || description !== undefined
  return (
    <section
      className={cn(
        "flex flex-col overflow-hidden rounded-lg border border-border bg-card shadow-xs",
        className,
      )}
    >
      {hasHeading || action ? (
        <div
          className={cn(
            "flex min-h-14 items-center gap-2 px-5 py-4",
            hasHeading ? "border-b border-border" : "justify-end",
          )}
        >
          {hasHeading ? (
            <div className="flex flex-1 flex-col">
              {title ? <h3 className="text-base font-medium text-strong">{title}</h3> : null}
              {description ? (
                <p className="text-sm text-muted-foreground">{description}</p>
              ) : null}
            </div>
          ) : null}
          {action}
        </div>
      ) : null}
      <div className={cn("p-5", bodyClassName)}>{children}</div>
    </section>
  )
}
