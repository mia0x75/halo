import { useMemo, useState } from "react"
import { IconChecks } from "@tabler/icons-react"

import { TrendChart } from "@/components/charts/trend-chart"
import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { StatusBadge } from "@/components/shared/status-badge"
import { Button } from "@/components/ui/button"
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
import { TopRules } from "@/features/dashboard/top-rules"
import {
  auditLogs,
  clusterTrend,
  clusters,
  periodLabels,
  tickets,
  trendFor,
  type Period,
} from "@/lib/mock/data"

/** 后端的命中级别：0 提示 / 1 错误 / 2 警告。 */
const levelLabels: Record<number, string> = { 0: "提示", 1: "错误", 2: "警告" }

export function DashboardPage() {
  const [period, setPeriod] = useState<Period>(7)
  const points = trendFor(period)

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
  const executed = points[points.length - 1]!.executed
  const failed = points[points.length - 1]!.failed

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

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          label="待审核工单"
          value={String(pending)}
          delta={12}
          hint="较昨日"
          trend={points.map((point) => point.submitted)}
        />
        <StatCard
          label="今日已执行"
          value={String(metricValue("tickets-executed", executed))}
          delta={8}
          hint="较昨日"
          trend={points.map((point) => point.executed)}
        />
        <StatCard
          label="失败工单"
          value={String(metricValue("tickets-failed", failed))}
          delta={-4}
          hint="较昨日"
          trend={points.map((point) => point.failed)}
        />
        <StatCard
          label="接入群集"
          value={String(clusters.length)}
          hint={`${clusters.filter((c) => c.status === "normal").length} 个正常`}
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
            categories={points.map((point) => point.date)}
            series={[
              { name: "提交", data: points.map((point) => point.submitted) },
              { name: "执行", data: points.map((point) => point.executed) },
              { name: "失败", data: points.map((point) => point.failed) },
            ]}
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
              <li key={ticket.uuid} className="flex items-center gap-3 px-5 py-3">
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
