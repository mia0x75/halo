import { IconArrowRight } from "@tabler/icons-react"

import { TrendChart } from "@/components/charts/trend-chart"
import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { StatusBadge } from "@/components/shared/status-badge"
import { Subheader } from "@/components/shared/subheader"
import { Button } from "@/components/ui/button"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { KpiCard } from "@/features/dashboard/kpi-card"
import { SeverityBreakdownChart } from "@/features/dashboard/severity-breakdown-chart"
import { formatNumber, formatRelativeTime } from "@/lib/format"
import { mockDashboardSummary } from "@/lib/api/mock"

export function DashboardPage() {
  const summary = mockDashboardSummary
  const categories = summary.trend.map((point) => point.date)

  return (
    <>
      <PageHeader
        title="概览"
        description="SQL 变更工单的提交、审核与执行情况"
        action={
          <Button>
            新建工单
            <IconArrowRight size={16} stroke={1.5} />
          </Button>
        }
      />

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <KpiCard
          label="待审核工单"
          value={formatNumber(summary.pendingTickets)}
          hint="系统审核 + 人工审核"
          trend={summary.trend.map((point) => point.submitted)}
        />
        <KpiCard
          label="今日已执行"
          value={formatNumber(summary.executedToday)}
          hint="上线执行完成"
          trend={summary.trend.map((point) => point.executed)}
        />
        <KpiCard
          label="失败工单"
          value={formatNumber(summary.failedTickets)}
          hint="需要重新提交"
        />
        <KpiCard label="接入群集" value={formatNumber(summary.clusters)} hint="MySQL + PostgreSQL" />
      </div>

      <div className="grid gap-4 xl:grid-cols-3">
        <SectionCard
          title="工单趋势"
          description="近 7 天提交与执行数量"
          className="xl:col-span-2"
          bodyClassName="p-0"
        >
          <TrendChart
            categories={categories}
            series={[
              { name: "提交", data: summary.trend.map((point) => point.submitted) },
              { name: "执行", data: summary.trend.map((point) => point.executed) },
            ]}
          />
        </SectionCard>

        <SectionCard title="风险分级" description="按规则命中级别统计">
          <SeverityBreakdownChart data={summary.severity} />
        </SectionCard>
      </div>

      <div className="grid gap-4 xl:grid-cols-3">
        <SectionCard
          title="最近工单"
          className="xl:col-span-2"
          bodyClassName="p-0"
          action={<Button variant="ghost" size="sm">查看全部</Button>}
        >
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>主题</TableHead>
                <TableHead>群集</TableHead>
                <TableHead>状态</TableHead>
                <TableHead>提交时间</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {summary.recentTickets.map((ticket) => (
                <TableRow key={ticket.uuid}>
                  <TableCell className="max-w-[280px] truncate font-medium">
                    {ticket.subject}
                  </TableCell>
                  <TableCell className="text-muted-foreground">{ticket.cluster.alias}</TableCell>
                  <TableCell>
                    <StatusBadge status={ticket.status} />
                  </TableCell>
                  <TableCell className="text-muted-foreground">
                    {formatRelativeTime(ticket.createAt)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </SectionCard>

        <SectionCard title="操作日志" description="最近审计记录">
          <ul className="flex flex-col gap-3">
            {summary.recentLogs.map((log) => (
              <li key={log.uuid} className="flex flex-col gap-1 border-b border-border pb-3 last:border-0 last:pb-0">
                <span className="text-sm">{log.operation}</span>
                <span className="text-xs text-muted-foreground">
                  {log.user.name ?? log.user.email} · {formatRelativeTime(log.createAt)}
                </span>
              </li>
            ))}
          </ul>
        </SectionCard>
      </div>

      <Subheader>系统说明</Subheader>
      <p className="text-xs text-muted-foreground">
        当前页面数据来自类型化 mock（src/lib/api/mock.ts）。接入后端 GraphQL 后，仅需替换数据来源，
        页面组件无需改动。
      </p>
    </>
  )
}
