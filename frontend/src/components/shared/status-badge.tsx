import { Badge } from "@/components/ui/badge"
import { ticketStatusMeta, type TicketStatus } from "@/lib/mock/data"

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
