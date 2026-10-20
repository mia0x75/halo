import type { HTMLAttributes } from "react"

import { cn } from "@/lib/utils"

export interface AvatarProps extends HTMLAttributes<HTMLSpanElement> {
  /** 头像文字，通常是姓名首字 */
  label: string
  size?: "sm" | "md" | "lg"
}

const sizeMap = {
  sm: "size-6 text-[0.625rem]",
  md: "size-8 text-xs",
  lg: "size-10 text-sm",
} as const

/** 头像：圆形，首字居中。 */
export function Avatar({ className, label, size = "md", ...props }: AvatarProps) {
  return (
    <span
      className={cn(
        "inline-flex shrink-0 items-center justify-center rounded-full bg-surface-secondary font-medium text-foreground",
        sizeMap[size],
        className,
      )}
      {...props}
    >
      {label}
    </span>
  )
}
