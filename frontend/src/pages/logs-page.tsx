import { IconHistory } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { auditLogs } from "@/lib/mock/data"

export function LogsPage() {
  return (
    <>
      <PageHeader title="操作日志" description="平台操作与审核动作的审计记录" />

      <SectionCard title="审计流水" bodyClassName="p-0">
        <ul className="divide-y divide-border">
          {auditLogs.map((log) => (
            <li key={log.id} className="flex items-start gap-3 px-5 py-3">
              <span className="mt-0.5 text-muted-foreground">
                <IconHistory size={16} stroke={1.5} />
              </span>
              <div className="min-w-0 flex-1">
                <p className="text-sm">{log.operation}</p>
                <p className="mt-0.5 text-xs text-muted-foreground">
                  {log.actor} · {log.target} · {log.ip}
                </p>
              </div>
              <span className="shrink-0 text-xs text-muted-foreground">{log.createdAt}</span>
            </li>
          ))}
        </ul>
      </SectionCard>
    </>
  )
}
