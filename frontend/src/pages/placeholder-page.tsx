import type { Icon } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { EmptyState } from "@/components/shared/empty-state"

export interface PlaceholderPageProps {
  title: string
  description: string
  icon: Icon
}

/** 尚未实现的区块占位页，保持导航结构完整。 */
export function PlaceholderPage({ title, description, icon: IconComponent }: PlaceholderPageProps) {
  return (
    <>
      <PageHeader title={title} description={description} />
      <div className="rounded-lg border border-border bg-card shadow-xs">
        <EmptyState
          icon={<IconComponent size={48} stroke={1} />}
          title="该模块尚未开放"
          description="演示版本先提供概览、工单、群集、规则、命中记录与操作日志，其余区块待接口就绪后接入。"
        />
      </div>
    </>
  )
}
