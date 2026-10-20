import { cva, type VariantProps } from "class-variance-authority"
import type { ButtonHTMLAttributes } from "react"

import { cn } from "@/lib/utils"

// 尺寸与内边距按设计稿：主按钮 0.6875rem/1.75rem、小按钮 0.3125rem/0.5rem，字号 0.75rem。
const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90",
        outline: "border border-border bg-card text-foreground hover:bg-surface-secondary",
        ghost: "text-foreground hover:bg-surface-secondary",
        subtle: "bg-surface-secondary text-foreground hover:bg-surface-tertiary",
        danger: "bg-danger text-danger-foreground hover:bg-danger/90",
        link: "text-primary underline-offset-4 hover:underline",
      },
      size: {
        // 0.5625rem/1rem 内边距 + 1.25rem 行高 = 2.375rem，与设计稿一致
        default: "px-4 py-[0.5625rem]",
        sm: "rounded-sm px-2 py-[0.3125rem] text-xs",
        lg: "px-6 py-3 text-sm",
        icon: "size-9",
        "icon-sm": "size-8 rounded-sm",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
)

export interface ButtonProps
  extends ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {}

export function Button({ className, variant, size, type = "button", ...props }: ButtonProps) {
  return (
    <button type={type} className={cn(buttonVariants({ variant, size }), className)} {...props} />
  )
}
