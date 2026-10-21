import { Badge } from "@/components/ui/badge"
import { ticketStatusMeta } from "@/lib/labels"
import type { TicketStatus } from "@/lib/model"

export interface StatusBadgeProps {
  status: TicketStatus
  className?: string
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const meta = ticketStatusMeta[status]
  return (
    <Badge variant={meta.tone} className={className}>
      {meta.label}
    </Badge>
  )
}
