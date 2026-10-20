import { useMemo, useState } from "react"

import { BrushChart } from "@/components/charts/brush-chart"
import { TrendChart } from "@/components/charts/trend-chart"
import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { StatusBadge } from "@/components/shared/status-badge"
import { Button } from "@/components/ui/button"
import { useClusters } from "@/hooks/use-clusters"
import { useGraphQL } from "@/hooks/use-graphql"
import {
  METRICS_QUERY,
  RULE_HIT_LEVELS_QUERY,
  RULE_HIT_RULES_QUERY,
  type MetricNode,
  type RuleHitCountNode,
} from "@/lib/api/queries"
import { PeriodSelect } from "@/features/dashboard/period-select"
import { StatCard } from "@/features/dashboard/stat-card"
import { TicketCreateDialog } from "@/features/tickets/ticket-create-dialog"

import { TopRules } from "@/features/dashboard/top-rules"
import {
  auditLogs,
  clusterTrend,
  periodLabels,
  tickets,
  trendFor,
  type Period,
} from "@/lib/mock/data"

/** 趋势主图的 DOM id，刷选图按它驱动缩放。 */
const TREND_CHART_ID = "chart-ticket-trend"

/** 后端的命中级别：0 提示 / 1 错误 / 2 警告。 */
const levelLabels: Record<number, string> = { 0: "提示", 1: "错误", 2: "警告" }

export interface DashboardPageProps {
  onOpenTicket: (uuid: string) => void
}

export function DashboardPage({ onOpenTicket }: DashboardPageProps) {
  const [period, setPeriod] = useState<Period>(7)
  const mockPoints = trendFor(period)
  const clusterPage = useClusters(50)

  // 接口未登录或不可用时回落到假数据，演示与联调都能看。
  const metrics = useGraphQL<{ metrics: MetricNode[] }, { scope: number }>(
    METRICS_QUERY,
    { scope: 1 },
    { metrics: [] },
  )
  const metricValue = (name: string, fallback: number) => {
    const found = metrics.data?.metrics.find((item) => item.Name === name)
    return found ? found.Value : fallback
  }

  // 刷选图固定展示较长区间（最多一年），主图跟随选中的窗口。
  // 数据优先取接口的日桶；接口不可用时用假数据补齐。
  const fullPoints = useMemo(() => {
    const buckets = metrics.data?.metrics ?? []
    if (buckets.length === 0) return trendFor(365)

    const byDay = new Map<string, { created: number; executed: number; failed: number }>()
    for (const item of buckets) {
      const day = byDay.get(item.Bucket) ?? { created: 0, executed: 0, failed: 0 }
      if (item.Name === "tickets-created") day.created = item.Value
      if (item.Name === "tickets-executed") day.executed = item.Value
      if (item.Name === "tickets-failed") day.failed = item.Value
      byDay.set(item.Bucket, day)
    }

    const days = Array.from(byDay.keys()).sort()
    if (days.length === 0) return trendFor(365)

    const points = days.map((day) => ({
      date: day.slice(5),
      submitted: byDay.get(day)?.created ?? 0,
      executed: byDay.get(day)?.executed ?? 0,
      failed: byDay.get(day)?.failed ?? 0,
    }))

    // 点太多时按步长采样，曲线才有意义
    const max = 365
    if (points.length <= max) return points
    const step = Math.ceil(points.length / max)
    return points.filter((_, index) => index % step === 0)
  }, [metrics.data])

  // 主图默认显示最近 period 天，拖动刷选滑块可切换到别的窗口
  const visiblePoints = useMemo(() => fullPoints.slice(-period), [fullPoints, period])

  const levels = useGraphQL<{ ruleHitLevels: RuleHitCountNode[] }, { from: number; to: number }>(
    RULE_HIT_LEVELS_QUERY,
    { from: 0, to: 0 },
    { ruleHitLevels: [] },
  )
  const rules = useGraphQL<
    { ruleHitRules: RuleHitCountNode[] },
    { from: number; to: number; first: number }
  >(
    RULE_HIT_RULES_QUERY,
    { from: 0, to: 0, first: 5 },
    { ruleHitRules: [] },
  )

  const levelRows = useMemo(() => {
    const source = levels.data?.ruleHitLevels ?? []
    if (source.length === 0) return undefined
    return source.map((item) => ({
      severity: levelLabels[item.Level] ?? "提示",
      count: item.Count,
    }))
  }, [levels.data])

  const ruleRows = useMemo(() => {
    const source = rules.data?.ruleHitRules ?? []
    if (source.length === 0) return undefined
    return source.map((item) => ({ rule: item.RuleName ?? "-", count: item.Count }))
  }, [rules.data])

  const pending = tickets.filter((ticket) =>
    ["waiting-vld", "waiting-mrv"].includes(ticket.status),
  ).length
  const executed = mockPoints[mockPoints.length - 1]!.executed
  const failed = mockPoints[mockPoints.length - 1]!.failed

  return (
    <>
      <PageHeader
        title="概览"
        description="SQL 变更工单的提交、审核与执行情况"
        action={
          <TicketCreateDialog
            onCreated={() => {
              metrics.refetch()
            }}
          />
        }
      />

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          label="待审核工单"
          value={String(pending)}
          delta={12}
          hint="较昨日"
          trend={visiblePoints.map((point) => point.submitted)}
        />
        <StatCard
          label="今日已执行"
          value={String(metricValue("tickets-executed", executed))}
          delta={8}
          hint="较昨日"
          trend={visiblePoints.map((point) => point.executed)}
        />
        <StatCard
          label="失败工单"
          value={String(metricValue("tickets-failed", failed))}
          delta={-4}
          hint="较昨日"
          trend={visiblePoints.map((point) => point.failed)}
        />
        <StatCard
          label="接入群集"
          value={String(clusterPage.total)}
          hint={`${clusterPage.clusters.filter((c) => c.status === "normal").length} 个正常`}
          trend={clusterTrend}
        />
      </div>

      <div className="grid gap-4 xl:grid-cols-3">
        <SectionCard
          title="工单趋势"
          description={`${periodLabels[period]}提交、执行与失败数量`}
          className="xl:col-span-2"
          bodyClassName="p-0"
          action={<PeriodSelect value={period} onChange={setPeriod} />}
        >
          <TrendChart
            id={TREND_CHART_ID}
            selectable
            range={period}
            categories={fullPoints.map((point) => point.date)}
            series={[
              { name: "提交", data: fullPoints.map((point) => point.submitted) },
              { name: "执行", data: fullPoints.map((point) => point.executed) },
              { name: "失败", data: fullPoints.map((point) => point.failed) },
            ]}
          />
          <BrushChart
            target={TREND_CHART_ID}
            categories={fullPoints.map((point) => point.date)}
            series={[{ name: "提交", data: fullPoints.map((point) => point.submitted) }]}
          />
        </SectionCard>

        <SectionCard
          title="风险分级"
          description={`${periodLabels[period]}规则命中分布`}
          bodyClassName="p-0"
          action={<PeriodSelect value={period} onChange={setPeriod} />}
        >
          <TopRules levels={levelRows} rules={ruleRows} />
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
              <li
                key={ticket.uuid}
                className="flex cursor-pointer items-center gap-3 px-5 py-3"
                onClick={() => onOpenTicket(ticket.uuid)}
              >
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-medium">{ticket.subject}</p>
                  <p className="truncate text-sm text-muted-foreground">
                    {ticket.database} · {ticket.cluster} · {ticket.submitter}
                  </p>
                </div>
                <StatusBadge status={ticket.status} />
                <span className="w-28 shrink-0 text-right text-sm text-muted-foreground">
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
                <p className="mt-0.5 text-sm text-muted-foreground">
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
