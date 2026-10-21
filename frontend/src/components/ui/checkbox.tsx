import { Checkbox as CheckboxPrimitive } from "radix-ui"
import { IconCheck } from "@tabler/icons-react"
import type { ComponentProps } from "react"

import { cn } from "@/lib/utils"

export interface CheckboxProps extends ComponentProps<typeof CheckboxPrimitive.Root> {
  label?: string
}

/** 复选框：主色勾选框，配合 label 使用。 */
export function Checkbox({ className, label, children, ...props }: CheckboxProps) {
  const box = (
    <CheckboxPrimitive.Root
      className={cn(
        "peer size-4 shrink-0 rounded border border-border bg-background transition-colors",
        "data-[state=checked]:border-primary data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40",
        "disabled:cursor-not-allowed disabled:opacity-50",
        className,
      )}
      {...props}
    >
      <CheckboxPrimitive.Indicator className="flex items-center justify-center">
        <IconCheck size={12} stroke={2.5} />
      </CheckboxPrimitive.Indicator>
    </CheckboxPrimitive.Root>
  )

  if (!label && !children) return box

  return (
    <label className="flex cursor-pointer items-center gap-2 text-sm">
      {box}
      {children ?? label}
    </label>
  )
}
