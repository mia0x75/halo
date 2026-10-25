import { useMemo, useState } from "react"
import { IconCheck, IconMessagePlus, IconX } from "@tabler/icons-react"

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { StatusBadge } from "@/components/shared/status-badge"
import { BlockNote, InlineNote } from "@/components/shared/status-note"
import { isInlineNote } from "@/lib/status-note"
import { useGraphQL } from "@/hooks/use-graphql"
import { graphqlRequest } from "@/lib/api/client"
import { notifyDataChanged } from "@/lib/api/data-events"
import { formatUnix } from "@/lib/api/mappers"
import {
  CREATE_COMMENT_MUTATION,
  REVIEW_TICKET_MUTATION,
  TICKET_DETAIL_QUERY,
  type CommentNode,
  type StatementNode,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"
import { parseReport, reportLevelLabel } from "@/features/tickets/ticket-report"
import type { TicketStatus } from "@/lib/model"

export interface TicketDetailViewProps {
  uuid: string | null
  /** 进入编辑页 */
  onEdit?: () => void
  /** 当前登录用户，用于判断能否审核 */
  viewerUUID?: string
  /** 驳回工单：需要填写原因，走独立页面 */
  onReject?: (uuid: string) => void
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
  Cluster: { UUID: string; Alias: string; Host: string; IP: string; Port: number; Type: string; Group: string | null }
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

/** 已关闭、已执行成功、执行失败的工单不可再改。 */
const editableStatuses = new Set([1, 2, 3, 4, 5, 6])

/** 只有存在审核意见的阶段才允许补充意见，与后端一致。 */
const commentableStatuses = new Set([3, 4, 5])

/** 工单详情内容：基本信息、语句列表与审核意见。 */
export function TicketDetailView({
  uuid,
  onEdit,
  viewerUUID,
  onReject,
}: TicketDetailViewProps) {
  const [tab, setTab] = useState<"statements" | "comments">("statements")
  const [approving, setApproving] = useState(false)
  const [reviewBusy, setReviewBusy] = useState(false)
  const [reviewError, setReviewError] = useState<string | null>(null)
  const [comment, setComment] = useState("")
  const [commentBusy, setCommentBusy] = useState(false)
  const [commentError, setCommentError] = useState<string | null>(null)

  const detail = useGraphQL<{ ticket: TicketDetail | null }, { id: string }>(
    TICKET_DETAIL_QUERY,
    { id: uuid ?? "" },
    { ticket: null },
  )

  const ticket = detail.data?.ticket ?? null
  const statements: StatementView[] = useMemo(
    () => ticket?.Statements?.edges?.map((edge) => toStatementView(edge.node)) ?? [],
    [ticket],
  )
  const comments: CommentView[] = useMemo(
    () => ticket?.Comments?.edges?.map((edge) => toCommentView(edge.node)) ?? [],
    [ticket],
  )

  const reviewErrorInline = reviewError ? isInlineNote(reviewError) : false
  const isReviewer = Boolean(viewerUUID && ticket && ticket.Reviewer?.UUID === viewerUUID)
  const canReview = ticket?.Status === 4 && isReviewer

  const approve = async () => {
    if (!ticket) return
    setReviewBusy(true)
    setReviewError(null)
    try {
      await graphqlRequest(
        REVIEW_TICKET_MUTATION,
        { input: { TicketUUID: ticket.UUID, Outcome: "APPROVED", Reason: "" } },
        getToken(),
      )
      setApproving(false)
      detail.refetch()
      notifyDataChanged("tickets")
    } catch (cause) {
      setReviewError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setReviewBusy(false)
    }
  }

  const submitComment = async () => {
    if (!ticket || !comment.trim()) return
    setCommentBusy(true)
    setCommentError(null)
    try {
      await graphqlRequest(
        CREATE_COMMENT_MUTATION,
        { input: { TicketUUID: ticket.UUID, Content: comment.trim() } },
        getToken(),
      )
      setComment("")
      detail.refetch()
    } catch (cause) {
      setCommentError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setCommentBusy(false)
    }
  }

  if (!ticket) {
    return (
      <div className="rounded-lg border border-border bg-card p-5 text-sm text-muted-foreground shadow-xs">
        {detail.loading ? "加载中…" : "未获取到工单信息"}
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="rounded-lg border border-border bg-card p-5 shadow-xs">
        <h2 className="text-lg font-medium text-strong">{ticket.Subject}</h2>
        <p className="mt-1 text-sm text-muted-foreground">
          {`${ticket.Database} · ${ticket.Cluster.Alias} · 提交人 ${ticket.User.Name ?? ticket.User.Email} · 审核人 ${ticket.Reviewer.Name ?? "-"} · ${formatUnix(ticket.CreateAt)}`}
        </p>
        <div className="mt-3 flex flex-wrap items-center gap-2">
          <StatusBadge status={statusByCode[ticket.Status] ?? "waiting-vld"} />
          <Badge variant="muted">工单 {ticket.UUID.slice(0, 8)}</Badge>
          {ticket.Cluster ? (
            <>
              <Badge variant="outline">
                {ticket.Cluster.IP}:{ticket.Cluster.Port}
              </Badge>
              <Badge variant="muted">{ticket.Cluster.Type}</Badge>
              {ticket.Cluster.Group ? (
                <Badge variant="muted">{ticket.Cluster.Group}</Badge>
              ) : null}
            </>
          ) : null}
          {onEdit && editableStatuses.has(ticket.Status) ? (
            <Button variant="outline" size="sm" onClick={onEdit}>
              编辑
            </Button>
          ) : null}
          {canReview ? (
            <>
              <Button size="sm" onClick={() => setApproving(true)}>
                <IconCheck size={16} stroke={1.5} />
                审核通过
              </Button>
              {onReject ? (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => onReject(ticket.UUID)}
                >
                  <IconX size={16} stroke={1.5} />
                  驳回
                </Button>
              ) : null}
            </>
          ) : null}
          {reviewError && reviewErrorInline ? (
            <InlineNote className="ml-1">{reviewError}</InlineNote>
          ) : null}
        </div>
      </div>
      {reviewError && !reviewErrorInline ? <BlockNote>{reviewError}</BlockNote> : null}

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
                {/* Level 1 is a hard failure; keeping it amber made errors look
                    like warnings at a glance. */}
                {parseReport(statement.report).map((hit, index) => (
                  <p
                    key={`${statement.uuid}-hit-${index}`}
                    className={
                      hit.level === 1
                        ? "border-t border-border px-3 py-2 text-xs text-danger"
                        : "border-t border-border px-3 py-2 text-xs text-warning"
                    }
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
          <div className="flex flex-col gap-3 p-5">
            {comments.length === 0 ? (
              <p className="py-2 text-center text-sm text-muted-foreground">暂无审核意见</p>
            ) : null}
            <ul className="flex flex-col gap-3">
              {comments.map((item) => (
                <li key={item.uuid} className="rounded-md border border-border px-3 py-2">
                  <p className="text-sm">{item.content}</p>
                  <p className="mt-1 text-xs text-muted-foreground">
                    {item.user} · {item.createAt}
                  </p>
                </li>
              ))}
            </ul>
            {commentableStatuses.has(ticket.Status) ? (
              <div className="flex flex-col gap-2 border-t border-border pt-4">
                <div className="flex items-center gap-2">
                  <Input
                    value={comment}
                    onChange={(event) => setComment(event.target.value)}
                    placeholder="写下审核意见"
                    aria-label="审核意见内容"
                  />
                  <Button
                    onClick={() => void submitComment()}
                    disabled={commentBusy || !comment.trim()}
                  >
                    <IconMessagePlus size={16} stroke={1.5} />
                    提交意见
                  </Button>
                </div>
                {commentError ? (
                  <p role="alert" className="text-sm text-danger">
                    {commentError}
                  </p>
                ) : null}
              </div>
            ) : null}
          </div>
        )}
      </div>

      <AlertDialog open={approving} onOpenChange={setApproving}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="text-base font-medium text-strong">
              审核通过
            </AlertDialogTitle>
            <AlertDialogDescription className="mt-1 text-sm text-muted-foreground">
              通过后工单进入待执行状态，可以立即执行或预约执行。确认通过？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <div className="mt-5 flex justify-end gap-2">
            <AlertDialogCancel asChild>
              <Button variant="outline">取消</Button>
            </AlertDialogCancel>
            <AlertDialogAction asChild>
              <Button disabled={reviewBusy} onClick={() => void approve()}>
                确认通过
              </Button>
            </AlertDialogAction>
          </div>
        </AlertDialogContent>
      </AlertDialog>
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
