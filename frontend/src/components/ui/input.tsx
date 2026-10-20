import type { InputHTMLAttributes } from "react"

import { cn } from "@/lib/utils"

/** 表单控件：0.5625rem/1rem 内边距、14px 字号、6px 圆角。 */
export function Input({ className, type = "text", ...props }: InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      type={type}
      className={cn(
        "flex h-[2.375rem] w-full rounded-md border border-input bg-card px-4 py-[0.5625rem] text-sm leading-5 text-foreground transition-colors",
        "placeholder:text-muted-foreground",
        "focus-visible:outline-none focus-visible:ring-3 focus-visible:ring-ring/50",
        "disabled:cursor-not-allowed disabled:opacity-50",
        className,
      )}
      {...props}
    />
  )
}
