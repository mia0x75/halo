import type { TicketNode } from "@/lib/api/queries"
import type { Ticket, TicketStatus } from "@/lib/model"

/** 后端的工单状态枚举（与 mm_tickets.status 一致）。 */
const statusByCode: Record<number, TicketStatus> = {
  1: "waiting-vld",
  2: "vld-failure",
  3: "vld-warning",
  4: "waiting-mrv",
  5: "mrv-failure",
  6: "ltgm",
  7: "done",
  8: "exec-failure",
  9: "closed",
}

/** 把接口返回的工单转成页面用的模型。 */
export function mapTicket(node: TicketNode): Ticket {
  return {
    id: 0,
    uuid: node.UUID,
    subject: node.Subject,
    database: node.Database,
    cluster: node.Cluster?.Alias ?? "-",
    status: statusByCode[node.Status] ?? "waiting-vld",
    submitter: node.User?.Name ?? node.User?.UUID ?? "-",
    reviewer: node.Reviewer?.Name ?? "-",
    createdAt: formatUnix(node.CreateAt),
    updatedAt: formatUnix(node.CreateAt),
    statementCount: 0,
    violationCount: 0,
  }
}

/** 把 UNIX 秒格式化成「YYYY-MM-DD HH:mm」。 */
export function formatUnix(value: number): string {
  if (!value) return "-"
  const date = new Date(value * 1000)
  const pad = (input: number) => input.toString().padStart(2, "0")
  return (
    `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ` +
    `${pad(date.getHours())}:${pad(date.getMinutes())}`
  )
}
