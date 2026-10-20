import type { HTMLAttributes } from "react"

import { cn } from "@/lib/utils"

export interface SubheaderProps extends HTMLAttributes<HTMLDivElement> {
  children: string
}

/** Section header in the dense admin style: small, uppercase, muted. */
export function Subheader({ className, children, ...props }: SubheaderProps) {
  return (
    <div
      className={cn(
        "text-xs font-medium uppercase tracking-[0.04em] text-muted-foreground",
        className,
      )}
      {...props}
    >
      {children}
    </div>
  )
}
