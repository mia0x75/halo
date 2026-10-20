import type { ReactNode } from "react"

import { Sparkline } from "@/components/charts/sparkline"
import { cn } from "@/lib/utils"

export interface KpiCardProps {
  label: string
  value: string
  hint?: ReactNode
  trend?: number[]
  className?: string
}

export function KpiCard({ label, value, hint, trend, className }: KpiCardProps) {
  return (
    <div
      className={cn(
        "flex flex-col rounded-lg border border-border bg-card p-5 shadow-xs",
        className,
      )}
    >
      <p className="text-sm font-medium text-muted-foreground">{label}</p>
      <div className="mt-2 flex items-end justify-between gap-3">
        <div className="flex flex-col gap-1">
          <span className="text-2xl font-semibold tracking-tight text-strong">{value}</span>
          {hint ? <span className="text-sm text-muted-foreground">{hint}</span> : null}
        </div>
        {trend ? <Sparkline data={trend} className="w-24" /> : null}
      </div>
    </div>
  )
}
