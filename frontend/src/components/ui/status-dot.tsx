import type { HTMLAttributes } from "react"

import { cn } from "@/lib/utils"

export interface StatusDotProps extends HTMLAttributes<HTMLSpanElement> {
  tone?: "primary" | "success" | "warning" | "danger" | "muted"
  label?: string
}

const toneMap = {
  primary: "bg-primary",
  success: "bg-success",
  warning: "bg-warning",
  danger: "bg-danger",
  muted: "bg-muted",
} as const

/** 状态点：6px 圆点，附带可读文本。 */
export function StatusDot({ className, tone = "primary", label, ...props }: StatusDotProps) {
  return (
    <span className={cn("inline-flex items-center gap-2", className)} {...props}>
      <span className={cn("size-1.5 rounded-full", toneMap[tone])} />
      {label ? <span className="text-sm">{label}</span> : null}
    </span>
  )
}
