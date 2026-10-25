import { useState } from "react"
import { IconArrowLeft, IconX } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { graphqlRequest } from "@/lib/api/client"
import { notifyDataChanged } from "@/lib/api/data-events"
import { REVIEW_TICKET_MUTATION } from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

export interface TicketRejectPageProps {
  uuid: string
  onBack: () => void
  onDone: () => void
}

/** 驳回工单页：驳回必须给出原因，因此是整页表单而不是弹窗。 */
export function TicketRejectPage({ uuid, onBack, onDone }: TicketRejectPageProps) {
  const [reason, setReason] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  const submit = async () => {
    if (!reason.trim()) {
      setError("请填写驳回原因")
      return
    }
    setSaving(true)
    setError(null)
    try {
      await graphqlRequest(
        REVIEW_TICKET_MUTATION,
        { input: { TicketUUID: uuid, Outcome: "REJECTED", Reason: reason.trim() } },
        getToken(),
      )
      notifyDataChanged("tickets")
      onDone()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setSaving(false)
    }
  }

  return (
    <>
      <PageHeader
        title="驳回工单"
        description="驳回后工单回到提交人，需要修改后重新提交"
        action={
          <Button variant="outline" onClick={onBack}>
            <IconArrowLeft size={16} stroke={1.5} />
            返回工单详情
          </Button>
        }
      />

      <SectionCard>
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="reject-reason">原因</Label>
          <Textarea
            id="reject-reason"
            value={reason}
            onChange={(event) => setReason(event.target.value)}
            placeholder="例如：缺少回滚脚本，执行范围过大"
            rows={5}
          />
        </div>

        {error ? (
          <p role="alert" className="mt-4 rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">
            {error}
          </p>
        ) : null}

        <div className="mt-5 flex justify-end gap-2">
          <Button variant="outline" onClick={onBack}>
            取消
          </Button>
          <Button variant="danger" onClick={() => void submit()} disabled={saving}>
            <IconX size={16} stroke={1.5} />
            确认驳回
          </Button>
        </div>
      </SectionCard>
    </>
  )
}
