import { IconTrendingDown, IconTrendingUp } from "@tabler/icons-react"
import type { ReactNode } from "react"

import { Sparkline } from "@/components/charts/sparkline"
import { cn } from "@/lib/utils"

export interface StatCardProps {
  label: string
  value: string
  /** 环比变化，正数为上升 */
  delta?: number
  hint?: string
  trend?: number[]
  icon?: ReactNode
}

/** 指标卡：标签 + 数值 + 环比 + 迷你趋势图。 */
export function StatCard({ label, value, delta, hint, trend, icon }: StatCardProps) {
  const rising = (delta ?? 0) >= 0

  return (
    <div className="flex flex-col rounded-lg border border-border bg-card p-5 shadow-xs">
      <div className="flex items-start justify-between gap-3">
        <p className="text-sm font-medium text-muted-foreground">
          {label}
        </p>
        {icon ? <span className="text-muted-foreground [&_svg]:size-5">{icon}</span> : null}
      </div>

      <div className="mt-3 flex items-end justify-between gap-3">
        <div>
          <p className="text-2xl font-semibold tracking-tight text-strong">{value}</p>
          <div className="mt-1 flex items-center gap-2 text-sm">
            {delta !== undefined ? (
              <span
                className={cn(
                  "inline-flex items-center gap-0.5 font-medium",
                  rising ? "text-success" : "text-danger",
                )}
              >
                {rising ? (
                  <IconTrendingUp size={14} stroke={1.5} />
                ) : (
                  <IconTrendingDown size={14} stroke={1.5} />
                )}
                {Math.abs(delta)}%
              </span>
            ) : null}
            {hint ? <span className="text-muted-foreground">{hint}</span> : null}
          </div>
        </div>
        {trend ? <Sparkline data={trend} className="w-24" /> : null}
      </div>
    </div>
  )
}
