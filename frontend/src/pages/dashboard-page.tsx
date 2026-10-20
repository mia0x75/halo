import { IconAlertTriangle, IconChecks, IconClockHour4, IconServer } from "@tabler/icons-react"

import { TrendChart } from "@/components/charts/trend-chart"
import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { StatusBadge } from "@/components/shared/status-badge"
import { Button } from "@/components/ui/button"
import { StatCard } from "@/features/dashboard/stat-card"
import { TopRules } from "@/features/dashboard/top-rules"
import { auditLogs, clusters, tickets, trend } from "@/lib/mock/data"

export function DashboardPage() {
  const pending = tickets.filter((ticket) =>
    ["waiting-vld", "waiting-mrv"].includes(ticket.status),
  ).length
  const executed = trend[trend.length - 1]!.executed
  const failed = trend[trend.length - 1]!.failed

  return (
    <>
      <PageHeader
        title="概览"
        description="SQL 变更工单的提交、审核与执行情况"
        action={
          <Button>
            新建工单
            <IconChecks size={16} stroke={1.5} />
          </Button>
        }
      />

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="待审核工单"
          value={String(pending)}
          delta={12}
          hint="较昨日"
          trend={trend.map((point) => point.submitted)}
          icon={<IconClockHour4 size={20} stroke={1.5} />}
        />
        <StatCard
          label="今日已执行"
          value={String(executed)}
          delta={8}
          hint="较昨日"
          trend={trend.map((point) => point.executed)}
          icon={<IconChecks size={20} stroke={1.5} />}
        />
        <StatCard
          label="失败工单"
          value={String(failed)}
          delta={-4}
          hint="较昨日"
          trend={trend.map((point) => point.failed)}
          icon={<IconAlertTriangle size={20} stroke={1.5} />}
        />
        <StatCard
          label="接入群集"
          value={String(clusters.length)}
          hint={`${clusters.filter((c) => c.status === "normal").length} 个正常`}
          icon={<IconServer size={20} stroke={1.5} />}
        />
      </div>

      <div className="grid gap-4 xl:grid-cols-3">
        <SectionCard
          title="工单趋势"
          description="近 7 天提交、执行与失败数量"
          className="xl:col-span-2"
          bodyClassName="p-0"
          action={<Button variant="ghost" size="sm">近 7 天</Button>}
        >
          <TrendChart
            categories={trend.map((point) => point.date)}
            series={[
              { name: "提交", data: trend.map((point) => point.submitted) },
              { name: "执行", data: trend.map((point) => point.executed) },
              { name: "失败", data: trend.map((point) => point.failed) },
            ]}
          />
        </SectionCard>

        <SectionCard title="风险分级" description="近 7 天规则命中分布" bodyClassName="p-0">
          <TopRules />
        </SectionCard>
      </div>

      <div className="grid gap-4 xl:grid-cols-3">
        <SectionCard
          title="最近工单"
          className="xl:col-span-2"
          bodyClassName="p-0"
          action={<Button variant="ghost" size="sm">查看全部</Button>}
        >
          <ul className="divide-y divide-border">
            {tickets.slice(0, 5).map((ticket) => (
              <li key={ticket.uuid} className="flex items-center gap-3 px-5 py-3">
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-medium">{ticket.subject}</p>
                  <p className="truncate text-xs text-muted-foreground">
                    {ticket.database} · {ticket.cluster} · {ticket.submitter}
                  </p>
                </div>
                <StatusBadge status={ticket.status} />
                <span className="w-28 shrink-0 text-right text-xs text-muted-foreground">
                  {ticket.createdAt}
                </span>
              </li>
            ))}
          </ul>
        </SectionCard>

        <SectionCard title="操作日志" description="最近审计记录" bodyClassName="p-0">
          <ul className="divide-y divide-border">
            {auditLogs.slice(0, 6).map((log) => (
              <li key={log.id} className="px-5 py-3">
                <p className="text-sm">{log.operation}</p>
                <p className="mt-0.5 text-xs text-muted-foreground">
                  {log.actor} · {log.ip} · {log.createdAt}
                </p>
              </li>
            ))}
          </ul>
        </SectionCard>
      </div>
    </>
  )
}
