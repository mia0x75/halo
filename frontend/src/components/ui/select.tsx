import { IconChevronDown } from "@tabler/icons-react"
import type { SelectHTMLAttributes } from "react"

import { cn } from "@/lib/utils"

export interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {}

/** 下拉选择：与输入框同高同内边距，右侧固定箭头。 */
export function Select({ className, children, ...props }: SelectProps) {
  return (
    <div className="relative">
      <select
        className={cn(
          "h-[2.375rem] w-full appearance-none rounded-md border border-input bg-card py-[0.5625rem] pr-9 pl-4 text-sm leading-5 text-foreground transition-colors",
          "focus-visible:outline-none focus-visible:ring-3 focus-visible:ring-ring/50",
          "disabled:cursor-not-allowed disabled:opacity-50",
          className,
        )}
        {...props}
      >
        {children}
      </select>
      <IconChevronDown
        size={16}
        stroke={1.5}
        className="pointer-events-none absolute top-1/2 right-3 -translate-y-1/2 text-muted-foreground"
      />
    </div>
  )
}
