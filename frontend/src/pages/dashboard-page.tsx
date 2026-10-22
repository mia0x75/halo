import { useMemo, useState } from "react"
import { IconPlus } from "@tabler/icons-react"

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
import { type Period } from "@/features/dashboard/periods"
import { StatCard } from "@/features/dashboard/stat-card"

import { TopRules } from "@/features/dashboard/top-rules"
import { mapTicket, formatUnix } from "@/lib/api/mappers"
import { LOGS_QUERY, TICKETS_QUERY, type LogConnectionNode, type TicketConnectionNode } from "@/lib/api/queries"
import { parseLogOperation } from "@/lib/labels"
import type { SeverityPoint, TopRulePoint } from "@/lib/model"

/** 趋势主图的 DOM id，刷选图按它驱动缩放。 */
const TREND_CHART_ID = "chart-ticket-trend"

/** 指标列表在 schema 里可为 null，空结果不能当成数组直接用。 */
interface MetricsPayload {
  metrics: MetricNode[] | null
}

/** 后端的命中级别：0 提示 / 1 错误 / 2 警告。 */
const levelLabels: Record<number, string> = { 0: "提示", 1: "错误", 2: "警告" }

export interface DashboardPageProps {
  onOpenTicket: (uuid: string) => void
  onCreateTicket: () => void
  onViewTickets: () => void
}

export function DashboardPage({ onOpenTicket, onCreateTicket, onViewTickets }: DashboardPageProps) {
  const [period, setPeriod] = useState<Period>(7)
  // 刷选图初始化时必须能找到主图的 DOM id，等主图挂载完成再渲染它。
  const [trendReady, setTrendReady] = useState(false)
  const clusterPage = useClusters(50)

  const metrics = useGraphQL<MetricsPayload, { scope: number }>(METRICS_QUERY, { scope: 1 }, { metrics: [] })
  const metricBuckets = useMemo(() => metrics.data?.metrics ?? [], [metrics.data])
  const metricValue = (name: string) => metricBuckets.find((item) => item.Name === name)?.Value ?? 0

  // 刷选图固定展示较长区间（最多一年），主图跟随选中的窗口。
  const fullPoints = useMemo(() => {
    const buckets = metricBuckets
    if (buckets.length === 0) return []

    const byDay = new Map<string, { created: number; executed: number; failed: number }>()
    for (const item of buckets) {
      const day = byDay.get(item.Bucket) ?? { created: 0, executed: 0, failed: 0 }
      if (item.Name === "tickets-created") day.created = item.Value
      if (item.Name === "tickets-executed") day.executed = item.Value
      if (item.Name === "tickets-failed") day.failed = item.Value
      byDay.set(item.Bucket, day)
    }

    const days = Array.from(byDay.keys()).sort()
    if (days.length === 0) return []

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
  }, [metricBuckets])

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

  const levelRows = useMemo<SeverityPoint[]>(
    () =>
      (levels.data?.ruleHitLevels ?? []).map((item) => ({
        severity: levelLabels[item.Level] ?? "提示",
        count: item.Count,
      })),
    [levels.data],
  )

  const ruleRows = useMemo<TopRulePoint[]>(
    () => (rules.data?.ruleHitRules ?? []).map((item) => ({ rule: item.RuleName || "-", count: item.Count })),
    [rules.data],
  )

  // 待审核 = 待系统审核 + 待人工审核，两次计数查询合成一个数。
  const waitingVld = useGraphQL<{ tickets: TicketConnectionNode | null }, { first: number; status: number }>(
    TICKETS_QUERY,
    { first: 1, status: 1 },
  )
  const waitingMrv = useGraphQL<{ tickets: TicketConnectionNode | null }, { first: number; status: number }>(
    TICKETS_QUERY,
    { first: 1, status: 4 },
  )
  const pending =
    (waitingVld.data?.tickets?.totalCount ?? 0) + (waitingMrv.data?.tickets?.totalCount ?? 0)

  const recent = useGraphQL<{ tickets: TicketConnectionNode | null }, { first: number }>(TICKETS_QUERY, {
    first: 5,
  })
  const recentTickets = useMemo(
    () => (recent.data?.tickets?.edges ?? []).map((edge) => mapTicket(edge.node)),
    [recent.data],
  )

  const logs = useGraphQL<{ logs: LogConnectionNode | null }, { first: number }>(LOGS_QUERY, { first: 6 })
  const recentLogs = logs.data?.logs?.edges ?? []

  return (
    <>
      <PageHeader
        title="概览"
        description="SQL 变更工单的提交、审核与执行情况"
        action={
          <Button onClick={onCreateTicket}>
            <IconPlus size={16} stroke={1.5} />
            新建工单
          </Button>
        }
      />

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          label="待审核工单"
          value={String(pending)}
          trend={visiblePoints.map((point) => point.submitted)}
        />
        <StatCard
          label="今日已执行"
          value={String(metricValue("tickets-executed"))}
          trend={visiblePoints.map((point) => point.executed)}
        />
        <StatCard
          label="失败工单"
          value={String(metricValue("tickets-failed"))}
          trend={visiblePoints.map((point) => point.failed)}
        />
        <StatCard label="接入群集" value={String(clusterPage.total)} />
      </div>

      <div className="grid gap-4 xl:grid-cols-3">
        <SectionCard
          title="工单趋势"
          className="xl:col-span-2"
          bodyClassName="p-0"
          action={<PeriodSelect value={period} onChange={setPeriod} />}
        >
          <TrendChart
            id={TREND_CHART_ID}
            selectable
            range={period}
            onReady={() => setTrendReady(true)}
            categories={fullPoints.map((point) => point.date)}
            series={[
              { name: "提交", data: fullPoints.map((point) => point.submitted) },
              { name: "执行", data: fullPoints.map((point) => point.executed) },
              { name: "失败", data: fullPoints.map((point) => point.failed) },
            ]}
          />
          {trendReady ? (
            <BrushChart
              target={TREND_CHART_ID}
              categories={fullPoints.map((point) => point.date)}
              series={[{ name: "提交", data: fullPoints.map((point) => point.submitted) }]}
            />
          ) : (
            <div className="h-24" />
          )}
        </SectionCard>

        <SectionCard
          title="风险分级"
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
          action={
            <Button variant="ghost" size="sm" onClick={onViewTickets}>
              查看全部
            </Button>
          }
        >
          <ul className="divide-y divide-border">
            {recentTickets.length === 0 ? (
              <li className="px-5 py-6 text-center text-sm text-muted-foreground">暂无工单</li>
            ) : null}
            {recentTickets.map((ticket) => (
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

        <SectionCard title="操作日志" bodyClassName="p-0">
          <ul className="divide-y divide-border">
            {recentLogs.length === 0 ? (
              <li className="px-5 py-6 text-center text-sm text-muted-foreground">暂无审计记录</li>
            ) : null}
            {recentLogs.map((edge) => {
              const display = parseLogOperation(edge.node.Operation)
              return (
                <li key={edge.node.UUID} className="px-5 py-3">
                  <p className="text-sm">
                    {display.label}
                    {display.target ? ` · ${display.target}` : ""}
                  </p>
                  <p className="mt-0.5 text-sm text-muted-foreground">
                    {edge.node.User?.Name || edge.node.User?.Email || "-"} ·{" "}
                    {formatUnix(edge.node.CreateAt)}
                  </p>
                </li>
              )
            })}
          </ul>
        </SectionCard>
      </div>
    </>
  )
}
