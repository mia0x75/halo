import type { TextareaHTMLAttributes } from "react"

import { cn } from "@/lib/utils"

/** 多行输入：与 Input 同款边框与字号，高度随内容行数变化。 */
export function Textarea({ className, ...props }: TextareaHTMLAttributes<HTMLTextAreaElement>) {
  return (
    <textarea
      className={cn(
        "w-full rounded-md border border-input bg-card px-4 py-2.5 text-sm leading-5 text-foreground transition-colors",
        "placeholder:text-muted-foreground",
        "focus-visible:outline-none focus-visible:ring-3 focus-visible:ring-ring/50",
        "disabled:cursor-not-allowed disabled:opacity-50",
        className,
      )}
      {...props}
    />
  )
}
