import { Badge } from "@/components/ui/badge"
import { TicketStatus } from "@/lib/api/types"

const statusMeta: Record<TicketStatus, { label: string; variant: "default" | "success" | "warning" | "danger" | "muted" | "outline" }> = {
  [TicketStatus.WaitingForVld]: { label: "待系统审核", variant: "outline" },
  [TicketStatus.VldFailure]: { label: "系统审核不通过", variant: "danger" },
  [TicketStatus.VldWarning]: { label: "系统审核告警", variant: "warning" },
  [TicketStatus.WaitingForMrv]: { label: "待人工审核", variant: "outline" },
  [TicketStatus.MrvFailure]: { label: "人工审核不通过", variant: "danger" },
  [TicketStatus.Ltgm]: { label: "待上线执行", variant: "default" },
  [TicketStatus.Done]: { label: "执行完成", variant: "success" },
  [TicketStatus.ExecFailure]: { label: "执行失败", variant: "danger" },
  [TicketStatus.Closed]: { label: "已关闭", variant: "muted" },
}

export interface StatusBadgeProps {
  status: TicketStatus
  className?: string
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const meta = statusMeta[status]
  return (
    <Badge variant={meta.variant} className={className}>
      {meta.label}
    </Badge>
  )
}
