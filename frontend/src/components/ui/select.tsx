import { IconCheck, IconChevronDown } from "@tabler/icons-react"
import { Select as SelectPrimitive } from "radix-ui"
import type { ComponentProps } from "react"

import { cn } from "@/lib/utils"

export const Select = SelectPrimitive.Root
export const SelectGroup = SelectPrimitive.Group
export const SelectValue = SelectPrimitive.Value

/** 下拉触发器：与输入框同高，右侧固定箭头，展开后弹出圆角面板。 */
export function SelectTrigger({
  className,
  children,
  ...props
}: ComponentProps<typeof SelectPrimitive.Trigger>) {
  return (
    <SelectPrimitive.Trigger
      className={cn(
        "flex h-[2.375rem] w-full items-center justify-between gap-2 rounded-md border border-input bg-card px-4 py-2 text-sm leading-5 text-foreground transition-colors",
        "data-[placeholder]:text-muted-foreground",
        "focus-visible:border-ring focus-visible:outline-none focus-visible:ring-3 focus-visible:ring-ring/50",
        "disabled:cursor-not-allowed disabled:opacity-50",
        "[&>span]:line-clamp-1 [&>span]:text-left",
        className,
      )}
      {...props}
    >
      {children}
      <SelectPrimitive.Icon asChild>
        <IconChevronDown
          size={16}
          stroke={1.5}
          className="shrink-0 text-muted-foreground"
        />
      </SelectPrimitive.Icon>
    </SelectPrimitive.Trigger>
  )
}

/** 弹出面板：圆角卡片，选中项带勾。 */
export function SelectContent({
  className,
  children,
  position = "popper",
  ...props
}: ComponentProps<typeof SelectPrimitive.Content>) {
  return (
    <SelectPrimitive.Portal>
      <SelectPrimitive.Content
        position={position}
        className={cn(
          "relative z-50 max-h-[var(--radix-select-content-available-height)] min-w-[8rem] overflow-hidden rounded-lg border border-border bg-popover p-1 text-popover-foreground shadow-lg",
          "data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0",
          position === "popper" &&
            "data-[side=bottom]:translate-y-1 data-[side=top]:-translate-y-1",
          className,
        )}
        {...props}
      >
        <SelectPrimitive.ScrollUpButton className="flex cursor-default items-center justify-center py-1" />
        <SelectPrimitive.Viewport
          className={cn(
            "p-0",
            position === "popper" && "w-full min-w-[var(--radix-select-trigger-width)]",
          )}
        >
          {children}
        </SelectPrimitive.Viewport>
        <SelectPrimitive.ScrollDownButton className="flex cursor-default items-center justify-center py-1" />
      </SelectPrimitive.Content>
    </SelectPrimitive.Portal>
  )
}

/** 选项：圆角，hover 高亮。 */
export function SelectItem({
  className,
  children,
  ...props
}: ComponentProps<typeof SelectPrimitive.Item>) {
  return (
    <SelectPrimitive.Item
      className={cn(
        "relative flex w-full cursor-pointer select-none items-center gap-2 rounded-md py-1.5 pr-8 pl-2.5 text-sm outline-none",
        "focus:bg-accent focus:text-accent-foreground",
        "data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
        className,
      )}
      {...props}
    >
      <SelectPrimitive.ItemText>{children}</SelectPrimitive.ItemText>
      <span className="absolute right-2 flex size-3.5 items-center justify-center">
        <SelectPrimitive.ItemIndicator>
          <IconCheck size={14} stroke={2} />
        </SelectPrimitive.ItemIndicator>
      </span>
    </SelectPrimitive.Item>
  )
}

/** 分组标题，对应原生 select 的 optgroup。 */
export function SelectLabel({
  className,
  ...props
}: ComponentProps<typeof SelectPrimitive.Label>) {
  return (
    <SelectPrimitive.Label
      className={cn(
        "px-2.5 py-1.5 text-xs font-medium text-muted-foreground",
        className,
      )}
      {...props}
    />
  )
}
