import type { ReactNode } from "react"

import type { NoteTone } from "@/lib/status-note"
import { cn } from "@/lib/utils"

const inlineTone: Record<NoteTone, string> = {
  success: "text-success",
  danger: "text-danger",
}

const blockTone: Record<NoteTone, string> = {
  success: "bg-success/10 text-success",
  danger: "bg-danger/10 text-danger",
}

export interface NoteProps {
  tone?: NoteTone
  children: ReactNode
  className?: string
}

export function InlineNote({ tone = "danger", children, className }: NoteProps) {
  return (
    <span
      role={tone === "danger" ? "alert" : "status"}
      className={cn("text-sm", inlineTone[tone], className)}
    >
      {children}
    </span>
  )
}

export function BlockNote({ tone = "danger", children, className }: NoteProps) {
  return (
    <p
      role={tone === "danger" ? "alert" : "status"}
      className={cn("rounded-md px-3 py-2 text-sm", blockTone[tone], className)}
    >
      {children}
    </p>
  )
}
