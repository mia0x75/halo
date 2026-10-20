import { useMemo } from "react"
import { IconArrowLeft } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Button } from "@/components/ui/button"
import { TicketForm, type TicketFormValues } from "@/features/tickets/ticket-form"
import { useGraphQL } from "@/hooks/use-graphql"
import { TICKET_DETAIL_QUERY } from "@/lib/api/queries"

export interface TicketEditPageProps {
  uuid: string
  onBack: () => void
  onSaved: (uuid: string) => void
}

interface TicketDetail {
  UUID: string
  Subject: string
  Content: string
  Database: string
  Status: number
  Cluster: { UUID: string }
  Reviewer: { UUID: string }
}

/** 工单编辑页：保存后重新走系统审核。 */
export function TicketEditPage({ uuid, onBack, onSaved }: TicketEditPageProps) {
  const detail = useGraphQL<{ ticket: TicketDetail | null }, { id: string }>(
    TICKET_DETAIL_QUERY,
    { id: uuid },
    { ticket: null },
  )

  const ticket = detail.data?.ticket ?? null

  const initial = useMemo<TicketFormValues | undefined>(() => {
    if (!ticket) return undefined
    return {
      clusterUUID: ticket.Cluster?.UUID ?? "",
      database: ticket.Database ?? "",
      subject: ticket.Subject ?? "",
      content: ticket.Content ?? "",
      reviewerUUID: ticket.Reviewer?.UUID ?? "",
    }
  }, [ticket])

  return (
    <>
      <PageHeader
        title="编辑工单"
        description="保存后工单回到系统自动审核，重新跑一遍规则检查"
        action={
          <Button variant="outline" onClick={onBack}>
            <IconArrowLeft size={16} stroke={1.5} />
            返回详情
          </Button>
        }
      />

      <SectionCard title="工单内容" description="改完保存即可，语句会按新的正文重新拆分">
        {initial ? (
          <TicketForm
            mode="edit"
            ticketUUID={uuid}
            initial={initial}
            onSubmitted={onSaved}
            onCancel={onBack}
          />
        ) : (
          <p className="text-sm text-muted-foreground">
            {detail.loading ? "正在读取工单信息…" : "未获取到工单信息"}
          </p>
        )}
      </SectionCard>
    </>
  )
}
