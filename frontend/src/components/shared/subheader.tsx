import type { HTMLAttributes } from "react"

import { cn } from "@/lib/utils"

export interface SubheaderProps extends HTMLAttributes<HTMLDivElement> {
  children: string
}

/** 区块小标题：12px 大写、字距放宽。 */
export function Subheader({ className, children, ...props }: SubheaderProps) {
  return (
    <div
      className={cn(
        "text-xs font-semibold uppercase tracking-[0.04em] text-muted-foreground",
        className,
      )}
      {...props}
    >
      {children}
    </div>
  )
}
