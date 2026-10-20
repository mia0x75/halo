import { useState } from "react"

import { Badge } from "@/components/ui/badge"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from "@/components/ui/dialog"
import { StatusBadge } from "@/components/shared/status-badge"
import { useGraphQL } from "@/hooks/use-graphql"
import { TICKET_DETAIL_QUERY, type CommentNode, type StatementNode } from "@/lib/api/queries"
import { formatUnix } from "@/lib/api/mappers"
import type { TicketStatus } from "@/lib/mock/data"

export interface TicketDetailDialogProps {
  uuid: string | null
  onClose: () => void
}

interface TicketDetail {
  UUID: string
  Subject: string
  Content: string
  Database: string
  Status: number
  CreateAt: number
  UpdateAt: number | null
  Cluster: { UUID: string; Alias: string; Host: string; IP: string; Port: number }
  User: { UUID: string; Name: string | null; Email: string }
  Reviewer: { UUID: string; Name: string | null }
  Statements: { edges: { node: StatementNode }[] }
  Comments: { edges: { node: CommentNode }[] } | null
}

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

/** 工单详情：基本信息、语句列表与审核意见。 */
export function TicketDetailDialog({ uuid, onClose }: TicketDetailDialogProps) {
  const [tab, setTab] = useState<"statements" | "comments">("statements")

  const detail = useGraphQL<{ ticket: TicketDetail | null }, { id: string }>(
    TICKET_DETAIL_QUERY,
    { id: uuid ?? "" },
    { ticket: null },
  )

  const ticket = detail.data?.ticket ?? null
  const statements = ticket?.Statements.edges.map((edge) => edge.node) ?? []
  const comments = ticket?.Comments?.edges.map((edge) => edge.node) ?? []

  return (
    <Dialog open={uuid !== null} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-w-3xl">
        {ticket ? (
          <>
            <DialogTitle className="text-base font-medium text-strong">{ticket.Subject}</DialogTitle>
            <DialogDescription className="mt-1 text-sm text-muted-foreground">
              {ticket.Database} · {ticket.Cluster.Alias} · 提交人 {ticket.User.Name ?? ticket.User.Email} ·
              审核人 {ticket.Reviewer.Name ?? "-"} · {formatUnix(ticket.CreateAt)}
            </DialogDescription>

            <div className="mt-3 flex items-center gap-2">
              <StatusBadge status={statusByCode[ticket.Status] ?? "waiting-vld"} />
              <Badge variant="muted">工单 {ticket.UUID.slice(0, 8)}</Badge>
            </div>

            <div className="mt-4 flex gap-4 border-b border-border">
              <button
                type="button"
                onClick={() => setTab("statements")}
                className={
                  tab === "statements"
                    ? "-mb-px border-b-2 border-primary pb-2 text-sm font-medium"
                    : "-mb-px border-b-2 border-transparent pb-2 text-sm text-muted-foreground"
                }
              >
                语句（{statements.length}）
              </button>
              <button
                type="button"
                onClick={() => setTab("comments")}
                className={
                  tab === "comments"
                    ? "-mb-px border-b-2 border-primary pb-2 text-sm font-medium"
                    : "-mb-px border-b-2 border-transparent pb-2 text-sm text-muted-foreground"
                }
              >
                审核意见（{comments.length}）
              </button>
            </div>

            {tab === "statements" ? (
              <div className="mt-3 flex max-h-72 flex-col gap-3 overflow-y-auto">
                {statements.length === 0 ? (
                  <p className="py-6 text-center text-sm text-muted-foreground">暂无语句</p>
                ) : null}
                {statements.map((statement) => (
                  <div key={statement.UUID} className="rounded-md border border-border">
                    <div className="flex items-center gap-2 border-b border-border px-3 py-2">
                      <span className="font-mono text-xs text-muted-foreground">
                        #{statement.Sequence}
                      </span>
                      <Badge variant="muted">{statement.TypeDesc}</Badge>
                    </div>
                    <pre className="overflow-x-auto px-3 py-2 font-mono text-xs whitespace-pre-wrap">
                      {statement.Content}
                    </pre>
                    {statement.Report ? (
                      <p className="border-t border-border px-3 py-2 text-xs text-warning">
                        {statement.Report}
                      </p>
                    ) : null}
                  </div>
                ))}
              </div>
            ) : (
              <ul className="mt-3 flex max-h-72 flex-col gap-3 overflow-y-auto">
                {comments.length === 0 ? (
                  <li className="py-6 text-center text-sm text-muted-foreground">暂无审核意见</li>
                ) : null}
                {comments.map((comment) => (
                  <li key={comment.UUID} className="rounded-md border border-border px-3 py-2">
                    <p className="text-sm">{comment.Content}</p>
                    <p className="mt-1 text-xs text-muted-foreground">
                      {comment.User.Name ?? "-"} · {formatUnix(comment.CreateAt)}
                    </p>
                  </li>
                ))}
              </ul>
            )}
          </>
        ) : (
          <>
            <DialogTitle className="text-base font-medium text-strong">工单详情</DialogTitle>
            <DialogDescription className="mt-1 text-sm text-muted-foreground">
              {detail.loading ? "加载中…" : "未获取到工单信息"}
            </DialogDescription>
          </>
        )}
      </DialogContent>
    </Dialog>
  )
}
