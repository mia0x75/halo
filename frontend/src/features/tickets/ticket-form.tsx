import { useMemo, useState } from "react"
import { IconSend } from "@tabler/icons-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import { Textarea } from "@/components/ui/textarea"
import { useGraphQL } from "@/hooks/use-graphql"
import { graphqlRequest } from "@/lib/api/client"
import {
  CREATE_TICKET_MUTATION,
  ME_CLUSTERS_QUERY,
  ME_REVIEWERS_QUERY,
  UPDATE_TICKET_MUTATION,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

export interface TicketFormValues {
  clusterUUID: string
  database: string
  subject: string
  content: string
  reviewerUUID: string
}

export interface TicketFormProps {
  mode: "create" | "edit"
  /** 编辑时带原工单 UUID */
  ticketUUID?: string
  /** 编辑时的初始值 */
  initial?: TicketFormValues
  onSubmitted: (uuid: string) => void
  onCancel: () => void
}

/** 工单表单：选群集、填库名与主题、正文写 SQL、指定审核人。 */
export function TicketForm({
  mode,
  ticketUUID,
  initial,
  onSubmitted,
  onCancel,
}: TicketFormProps) {
  const [clusterUUID, setClusterUUID] = useState(initial?.clusterUUID ?? "")
  const [database, setDatabase] = useState(initial?.database ?? "")
  const [subject, setSubject] = useState(initial?.subject ?? "")
  const [content, setContent] = useState(initial?.content ?? "")
  const [reviewerUUID, setReviewerUUID] = useState(initial?.reviewerUUID ?? "")
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const remoteClusters = useGraphQL<
    {
      me: {
        Clusters?: {
          edges: {
            node: {
              UUID: string
              Alias: string
              Host: string
              IP: string
              Port: number
              Type: string
              Group: string | null
            }
          }[]
        } | null
      } | null
    },
    { first: number }
  >(ME_CLUSTERS_QUERY, { first: 50 }, { me: null })

  const remoteReviewers = useGraphQL<
    { me: { Reviewers?: { UUID: string; Name: string | null; Email: string }[] | null } | null }
  >(ME_REVIEWERS_QUERY, undefined, { me: null })

  const clusterOptions = useMemo(() => {
    const edges = remoteClusters.data?.me?.Clusters?.edges ?? []
    const options = edges.map((edge) => ({
      uuid: edge.node.UUID,
      label: `${edge.node.Alias} · ${edge.node.Type} (${edge.node.IP}:${edge.node.Port})`,
    }))
    // 原工单的群集可能不在当前可选列表里，补一条免得回填的值丢失。
    if (initial?.clusterUUID && !options.some((option) => option.uuid === initial.clusterUUID)) {
      options.unshift({ uuid: initial.clusterUUID, label: "原工单群集" })
    }
    return options
  }, [remoteClusters.data, initial])

  const reviewerOptions = useMemo(() => {
    const reviewers = remoteReviewers.data?.me?.Reviewers ?? []
    const options = reviewers.map((user) => ({
      uuid: user.UUID,
      label: user.Name || user.Email,
    }))
    if (initial?.reviewerUUID && !options.some((option) => option.uuid === initial.reviewerUUID)) {
      options.unshift({ uuid: initial.reviewerUUID, label: "原工单审核人" })
    }
    return options
  }, [remoteReviewers.data, initial])

  // 默认选第一项：渲染时推导，不在 effect 里改状态。
  const selectedCluster = clusterUUID || clusterOptions[0]?.uuid || ""
  const selectedReviewer = reviewerUUID || reviewerOptions[0]?.uuid || ""

  const submit = async () => {
    // 没登录时接口会拒绝，这里先给出明确提示，避免只看到 1401 错误码。
    if (!getToken()) {
      setError("登录状态已失效，请重新登录后再提交。")
      return
    }

    setSubmitting(true)
    setError(null)
    try {
      const input = {
        ClusterUUID: selectedCluster,
        Database: database.trim(),
        Subject: subject.trim(),
        Content: content,
        ReviewerUUID: selectedReviewer,
      }
      if (mode === "edit") {
        await graphqlRequest(
          UPDATE_TICKET_MUTATION,
          { input: { ...input, TicketUUID: ticketUUID } },
          getToken(),
        )
        onSubmitted(ticketUUID ?? "")
        return
      }
      const data = await graphqlRequest<
        { createTicket: { UUID: string } },
        {
          input: {
            ClusterUUID: string
            Database: string
            Subject: string
            Content: string
            ReviewerUUID: string
          }
        }
      >(CREATE_TICKET_MUTATION, { input }, getToken())
      onSubmitted(data.createTicket.UUID)
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setSubmitting(false)
    }
  }

  const valid =
    selectedCluster && database.trim() && subject.trim() && content.trim() && selectedReviewer

  return (
    <div className="flex flex-col gap-4">
      <div className="grid gap-3 sm:grid-cols-2">
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="ticket-cluster">目标群集</Label>
          <Select
            id="ticket-cluster"
            value={selectedCluster}
            onChange={(event) => setClusterUUID(event.target.value)}
            disabled={clusterOptions.length === 0}
          >
            {clusterOptions.length === 0 ? <option value="">没有可用群集</option> : null}
            {clusterOptions.map((option) => (
              <option key={option.uuid} value={option.uuid}>
                {option.label}
              </option>
            ))}
          </Select>
        </div>
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="ticket-database">目标库</Label>
          <Input
            id="ticket-database"
            value={database}
            onChange={(event) => setDatabase(event.target.value)}
            placeholder="app_order"
          />
        </div>
      </div>

      <div className="flex flex-col gap-1.5">
        <Label htmlFor="ticket-subject">主题</Label>
        <Input
          id="ticket-subject"
          value={subject}
          onChange={(event) => setSubject(event.target.value)}
          placeholder="订单表新增归档字段"
        />
      </div>

      <div className="flex flex-col gap-1.5">
        <Label htmlFor="ticket-content">工单内容</Label>
        <Textarea
          id="ticket-content"
          value={content}
          onChange={(event) => setContent(event.target.value)}
          placeholder="ALTER TABLE orders ADD COLUMN archived_at DATETIME NULL;"
          className="min-h-56 font-mono"
        />
      </div>

      <div className="flex flex-col gap-1.5">
        <Label htmlFor="ticket-reviewer">审核人</Label>
        <Select
          id="ticket-reviewer"
          value={selectedReviewer}
          onChange={(event) => setReviewerUUID(event.target.value)}
        >
          {reviewerOptions.map((option) => (
            <option key={option.uuid} value={option.uuid}>
              {option.label}
            </option>
          ))}
        </Select>
      </div>

      {error ? (
        <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
      ) : null}

      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={onCancel}>
          取消
        </Button>
        <Button disabled={!valid || submitting} onClick={submit}>
          <IconSend size={16} stroke={1.5} />
          {submitting ? "提交中…" : mode === "edit" ? "保存工单" : "提交工单"}
        </Button>
      </div>
    </div>
  )
}
