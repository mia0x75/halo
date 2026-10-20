import { IconArrowLeft } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { Button } from "@/components/ui/button"
import { TicketDetailView } from "@/features/tickets/ticket-detail-view"

export interface TicketDetailPageProps {
  uuid: string
  onBack: () => void
}

/** 工单详情页：内容会越来越多，独立成页而不是弹窗。 */
export function TicketDetailPage({ uuid, onBack }: TicketDetailPageProps) {
  return (
    <>
      <PageHeader
        title="工单详情"
        description="工单信息、执行语句与审核意见"
        action={
          <Button variant="outline" onClick={onBack}>
            <IconArrowLeft size={16} stroke={1.5} />
            返回列表
          </Button>
        }
      />
      <TicketDetailView uuid={uuid} />
    </>
  )
}
