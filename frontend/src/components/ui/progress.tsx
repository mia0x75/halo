import type { HTMLAttributes } from "react"

import { cn } from "@/lib/utils"

export interface ProgressProps extends HTMLAttributes<HTMLDivElement> {
  /** 0-100 */
  value: number
  tone?: "primary" | "success" | "warning" | "danger"
}

const toneMap = {
  primary: "bg-primary",
  success: "bg-success",
  warning: "bg-warning",
  danger: "bg-danger",
} as const

/** 进度条：8px 高，圆角裁切。 */
export function Progress({ className, value, tone = "primary", ...props }: ProgressProps) {
  const clamped = Math.min(100, Math.max(0, value))
  return (
    <div
      className={cn("h-2 w-full overflow-hidden rounded-full bg-surface-secondary", className)}
      role="progressbar"
      aria-valuenow={clamped}
      aria-valuemin={0}
      aria-valuemax={100}
      {...props}
    >
      <div className={cn("h-full rounded-full", toneMap[tone])} style={{ width: `${clamped}%` }} />
    </div>
  )
}
