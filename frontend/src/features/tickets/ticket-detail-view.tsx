import { useMemo, useState } from "react"

import { Badge } from "@/components/ui/badge"
import { StatusBadge } from "@/components/shared/status-badge"
import { useGraphQL } from "@/hooks/use-graphql"
import { formatUnix } from "@/lib/api/mappers"
import { TICKET_DETAIL_QUERY, type CommentNode, type StatementNode } from "@/lib/api/queries"
import { parseReport, reportLevelLabel } from "@/features/tickets/ticket-report"
import {
  ticketComments,
  ticketStatements,
  tickets,
  type MockComment,
  type MockStatement,
  type TicketStatus,
} from "@/lib/mock/data"

export interface TicketDetailViewProps {
  uuid: string | null
}

interface StatementView {
  uuid: string
  sequence: number
  content: string
  typeDesc: string
  status: number
  report: string
}

interface CommentView {
  uuid: string
  content: string
  createAt: string
  user: string
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

/** 工单详情内容：基本信息、语句列表与审核意见。 */
export function TicketDetailView({ uuid }: TicketDetailViewProps) {
  const [tab, setTab] = useState<"statements" | "comments">("statements")

  const detail = useGraphQL<{ ticket: TicketDetail | null }, { id: string }>(
    TICKET_DETAIL_QUERY,
    { id: uuid ?? "" },
    { ticket: null },
  )

  const ticket = detail.data?.ticket ?? null
  // 接口没有数据时（未登录或接口不可用）回落到假数据，详情仍然可看。
  const mock = ticket ? null : mockTicket(uuid)
  const statements: StatementView[] = useMemo(
    () =>
      ticket
        ? (ticket.Statements?.edges?.map((edge) => toStatementView(edge.node)) ?? [])
        : (mock?.statements ?? []),
    [ticket, mock],
  )
  const comments: CommentView[] = useMemo(
    () =>
      ticket
        ? (ticket.Comments?.edges?.map((edge) => toCommentView(edge.node)) ?? [])
        : (mock?.comments ?? []),
    [ticket, mock],
  )

  if (!ticket && !mock) {
    return (
      <div className="rounded-lg border border-border bg-card p-5 text-sm text-muted-foreground shadow-xs">
        {detail.loading ? "加载中…" : "未获取到工单信息"}
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="rounded-lg border border-border bg-card p-5 shadow-xs">
        <h2 className="text-lg font-medium text-strong">{ticket?.Subject ?? mock?.subject}</h2>
        <p className="mt-1 text-sm text-muted-foreground">
          {ticket
            ? `${ticket.Database} · ${ticket.Cluster.Alias} · 提交人 ${ticket.User.Name ?? ticket.User.Email} · 审核人 ${ticket.Reviewer.Name ?? "-"} · ${formatUnix(ticket.CreateAt)}`
            : `${mock?.database} · ${mock?.cluster} · 提交人 ${mock?.submitter} · 审核人 ${mock?.reviewer} · ${mock?.createdAt}`}
        </p>
        <div className="mt-3 flex flex-wrap items-center gap-2">
          <StatusBadge status={statusByCode[ticket?.Status ?? 0] ?? "waiting-vld"} />
          <Badge variant="muted">工单 {(ticket?.UUID ?? mock?.uuid ?? "").slice(0, 8)}</Badge>
          {ticket?.Cluster ? (
            <Badge variant="outline">
              {ticket.Cluster.IP}:{ticket.Cluster.Port}
            </Badge>
          ) : null}
        </div>
      </div>

      <div className="rounded-lg border border-border bg-card shadow-xs">
        <div className="flex gap-4 border-b border-border px-5 pt-3">
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
          <div className="flex flex-col gap-3 p-5">
            {statements.length === 0 ? (
              <p className="py-6 text-center text-sm text-muted-foreground">暂无语句</p>
            ) : null}
            {statements.map((statement) => (
              <div key={statement.uuid} className="rounded-md border border-border">
                <div className="flex items-center gap-2 border-b border-border px-3 py-2">
                  <span className="font-mono text-xs text-muted-foreground">
                    #{statement.sequence}
                  </span>
                  <Badge variant="muted">{statement.typeDesc}</Badge>
                </div>
                <pre className="overflow-x-auto px-3 py-2 font-mono text-xs whitespace-pre-wrap">
                  {statement.content}
                </pre>
                {parseReport(statement.report).map((hit, index) => (
                  <p
                    key={`${statement.uuid}-hit-${index}`}
                    className="border-t border-border px-3 py-2 text-xs text-warning"
                  >
                    <span className="mr-1 font-medium">[{reportLevelLabel(hit.level)}]</span>
                    {hit.description}
                    {hit.rule ? (
                      <span className="ml-1 font-mono text-muted-foreground">（{hit.rule}）</span>
                    ) : null}
                  </p>
                ))}
              </div>
            ))}
          </div>
        ) : (
          <ul className="flex flex-col gap-3 p-5">
            {comments.length === 0 ? (
              <li className="py-6 text-center text-sm text-muted-foreground">暂无审核意见</li>
            ) : null}
            {comments.map((comment) => (
              <li key={comment.uuid} className="rounded-md border border-border px-3 py-2">
                <p className="text-sm">{comment.content}</p>
                <p className="mt-1 text-xs text-muted-foreground">
                  {comment.user} · {comment.createAt}
                </p>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}

/** 接口的语句节点转成页面视图。 */
function toStatementView(node: StatementNode): StatementView {
  return {
    uuid: node.UUID,
    sequence: node.Sequence,
    content: node.Content,
    typeDesc: node.TypeDesc,
    status: node.Status,
    report: node.Report ?? "",
  }
}

/** 接口的意见节点转成页面视图。 */
function toCommentView(node: CommentNode): CommentView {
  return {
    uuid: node.UUID,
    content: node.Content,
    createAt: formatUnix(node.CreateAt),
    user: node.User?.Name ?? "-",
  }
}

interface MockTicketView {
  uuid: string
  subject: string
  database: string
  cluster: string
  submitter: string
  reviewer: string
  createdAt: string
  statements: MockStatement[]
  comments: MockComment[]
}

/** 假数据里的工单详情。 */
function mockTicket(uuid: string | null): MockTicketView | null {
  if (!uuid) return null
  const ticket = tickets.find((item) => item.uuid === uuid)
  if (!ticket) return null
  return {
    uuid: ticket.uuid,
    subject: ticket.subject,
    database: ticket.database,
    cluster: ticket.cluster,
    submitter: ticket.submitter,
    reviewer: ticket.reviewer,
    createdAt: ticket.createdAt,
    statements: ticketStatements[ticket.uuid] ?? [],
    comments: ticketComments[ticket.uuid] ?? [],
  }
}
