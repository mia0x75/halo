import { useMemo, useState } from "react"
import { IconPlus } from "@tabler/icons-react"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
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
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"
import { clusters as mockClusters, users as mockUsers } from "@/lib/mock/data"

export interface TicketCreateDialogProps {
  onCreated?: (uuid: string) => void
  trigger?: React.ReactNode
}

/** 新建工单：选群集、填库名与主题、正文写 SQL、指定审核人。 */
export function TicketCreateDialog({ onCreated, trigger }: TicketCreateDialogProps) {
  const [open, setOpen] = useState(false)
  const [clusterUUID, setClusterUUID] = useState("")
  const [database, setDatabase] = useState("")
  const [subject, setSubject] = useState("")
  const [content, setContent] = useState("")
  const [reviewerUUID, setReviewerUUID] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const remoteClusters = useGraphQL<
    { me: { Clusters?: { edges: { node: { UUID: string; Alias: string; Host: string; IP: string; Port: number; Type: string } }[] } | null } | null },
    { first: number }
  >(ME_CLUSTERS_QUERY, { first: 50 }, { me: null })

  const remoteReviewers = useGraphQL<
    { me: { Reviewers?: { UUID: string; Name: string | null; Email: string }[] | null } | null }
  >(ME_REVIEWERS_QUERY, undefined, { me: null })

  const clusterOptions = useMemo(() => {
    const edges = remoteClusters.data?.me?.Clusters?.edges
    if (!edges || edges.length === 0) {
      return mockClusters.map((cluster) => ({
        uuid: cluster.uuid,
        label: `${cluster.alias} (${cluster.ip}:${cluster.port})`,
      }))
    }
    return edges.map((edge) => ({
      uuid: edge.node.UUID,
      label: `${edge.node.Alias} (${edge.node.IP}:${edge.node.Port})`,
    }))
  }, [remoteClusters.data])

  const reviewerOptions = useMemo(() => {
    const reviewers = remoteReviewers.data?.me?.Reviewers
    if (!reviewers || reviewers.length === 0) {
      return mockUsers.map((user) => ({ uuid: user.email, label: user.name }))
    }
    return reviewers.map((user) => ({
      uuid: user.UUID,
      label: user.Name || user.Email,
    }))
  }, [remoteReviewers.data])

  // 默认选第一项：渲染时推导，不在 effect 里改状态。
  const selectedCluster = clusterUUID || clusterOptions[0]?.uuid || ""
  const selectedReviewer = reviewerUUID || reviewerOptions[0]?.uuid || ""

  const submit = async () => {
    // 没登录时接口会拒绝，这里先给出明确提示，避免只看到 1401 错误码。
    if (!getToken()) {
      setError("尚未登录，无法提交工单。请先在右上角登录后重试。")
      return
    }

    setSubmitting(true)
    setError(null)
    try {
      const data = await graphqlRequest<
        { createTicket: { UUID: string } },
        { input: { ClusterUUID: string; Database: string; Subject: string; Content: string; ReviewerUUID: string } }
      >(
        CREATE_TICKET_MUTATION,
        {
          input: {
            ClusterUUID: selectedCluster,
            Database: database.trim(),
            Subject: subject.trim(),
            Content: content,
            ReviewerUUID: selectedReviewer,
          },
        },
        getToken(),
      )
      setOpen(false)
      setDatabase("")
      setSubject("")
      setContent("")
      onCreated?.(data.createTicket.UUID)
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setSubmitting(false)
    }
  }

  const valid =
    selectedCluster && database.trim() && subject.trim() && content.trim() && selectedReviewer

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {trigger ?? (
          <Button>
            <IconPlus size={16} stroke={1.5} />
            新建工单
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="max-w-xl">
        <DialogTitle className="text-base font-medium text-strong">新建工单</DialogTitle>
        <DialogDescription className="mt-1 text-sm text-muted-foreground">
          提交后由系统自动审核，通过后进入人工审核环节。
        </DialogDescription>

        <div className="mt-4 flex flex-col gap-3">
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="ticket-cluster">目标群集</Label>
              <Select
                id="ticket-cluster"
                value={selectedCluster}
                onChange={(event) => setClusterUUID(event.target.value)}
              >
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
              className="min-h-40 font-mono"
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <Label htmlFor="ticket-reviewer">审核人</Label>
            <Select
              id="ticket-reviewer"
              value={reviewerUUID}
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
        </div>

        <div className="mt-5 flex justify-end gap-2">
          <Button variant="outline" onClick={() => setOpen(false)}>
            取消
          </Button>
          <Button disabled={!valid || submitting} onClick={submit}>
            {submitting ? "提交中…" : "提交工单"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
