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
  /** 右上角操作区 */
  action?: ReactNode
}

/**
 * 指标卡。结构对齐设计稿的 Revenue / Active subscriptions：
 * 第一行标签与右上角操作，第二行大数值与环比，第三行图表独占一块贴底。
 */
export function StatCard({ label, value, delta, hint, trend, action }: StatCardProps) {
  const rising = (delta ?? 0) >= 0

  return (
    <div className="flex flex-col overflow-hidden rounded-lg border border-border bg-card shadow-xs">
      <div className="p-5 pb-0">
        <div className="flex items-center gap-2">
          <p className="text-sm font-medium text-strong">{label}</p>
          {action ? <div className="ml-auto">{action}</div> : null}
        </div>

        <div className="mt-2 flex items-baseline gap-2">
          <span className="text-2xl font-semibold leading-8 tracking-tight text-strong">
            {value}
          </span>
          {delta !== undefined ? (
            <span
              className={cn(
                "inline-flex items-center gap-1 text-sm leading-4 font-medium",
                rising ? "text-success" : "text-danger",
              )}
            >
              {rising ? (
                <IconTrendingUp size={16} stroke={1.5} />
              ) : (
                <IconTrendingDown size={16} stroke={1.5} />
              )}
              {Math.abs(delta)}%
            </span>
          ) : null}
          {hint ? <span className="text-sm text-muted-foreground">{hint}</span> : null}
        </div>
      </div>

      {trend ? (
        <div className="mt-3 px-2 pb-2">
          <Sparkline data={trend} height={56} />
        </div>
      ) : (
        <div className="mt-3" />
      )}
    </div>
  )
}
