import { IconArrowLeft } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Button } from "@/components/ui/button"
import { TicketForm } from "@/features/tickets/ticket-form"

export interface TicketCreatePageProps {
  onBack: () => void
  /** 提交成功后拿到工单 UUID */
  onCreated: (uuid: string) => void
}

/** 新建工单页。表单是整页，不是弹窗。 */
export function TicketCreatePage({ onBack, onCreated }: TicketCreatePageProps) {
  return (
    <>
      <PageHeader
        title="新建工单"
        description="提交后由系统自动审核，通过后进入人工审核环节"
        action={
          <Button variant="outline" onClick={onBack}>
            <IconArrowLeft size={16} stroke={1.5} />
            返回
          </Button>
        }
      />

      <SectionCard description="选择目标群集与审核人，正文写要执行的 SQL">
        <TicketForm mode="create" onSubmitted={onCreated} onCancel={onBack} />
      </SectionCard>
    </>
  )
}
