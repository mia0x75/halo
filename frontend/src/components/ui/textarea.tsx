import type { TextareaHTMLAttributes } from "react"

import { cn } from "@/lib/utils"

export interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {}

/** 多行输入：与输入框同一样式。 */
export function Textarea({ className, ...props }: TextareaProps) {
  return (
    <textarea
      className={cn(
        "w-full min-h-32 rounded-md border border-input bg-card px-4 py-2.5 text-sm leading-5 text-foreground transition-colors",
        "placeholder:text-muted-foreground",
        "focus-visible:outline-none focus-visible:ring-3 focus-visible:ring-ring/50",
        "disabled:cursor-not-allowed disabled:opacity-50",
        className,
      )}
      {...props}
    />
  )
}
