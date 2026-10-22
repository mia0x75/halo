import { useMemo, useState } from "react"
import { IconChartBar } from "@tabler/icons-react"

import { BrushChart } from "@/components/charts/brush-chart"
import { TrendChart } from "@/components/charts/trend-chart"
import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useGraphQL } from "@/hooks/use-graphql"
import { METRICS_QUERY, type MetricNode } from "@/lib/api/queries"

/** 统计报表里的指标名，与后端每日汇总写的一致。 */
const metricNames = [
  "tickets-created",
  "tickets-executed",
  "tickets-failed",
  "violations-error",
  "violations-warning",
  "violations-notice",
]

const metricLabels: Record<string, string> = {
  "tickets-created": "提交工单",
  "tickets-executed": "执行成功",
  "tickets-failed": "执行失败",
  "violations-error": "错误命中",
  "violations-warning": "警告命中",
  "violations-notice": "提示命中",
}

const TREND_CHART_ID = "chart-stats-trend"

/** 统计报表页：按天汇总的工单与命中指标。 */
export function StatsPage() {
  const [range, setRange] = useState("30")

  const remote = useGraphQL<{ metrics: MetricNode[] | null }, { Scope: number; Names: string[] }>(
    METRICS_QUERY,
    { Scope: 1, Names: metricNames },
    { metrics: [] },
  )

  const metrics = useMemo(() => remote.data?.metrics ?? [], [remote.data])

  // 后端按 bucket（日期）分桶，这里把同一日期的指标排到一起。
  const buckets = useMemo(() => {
    const map = new Map<string, Record<string, number>>()
    for (const metric of metrics) {
      const row = map.get(metric.Bucket) ?? {}
      row[metric.Name] = metric.Value
      map.set(metric.Bucket, row)
    }
    // 图表从左到右是时间顺序，升序排列后再截取最近的一段。
    return [...map.entries()].sort((a, b) => (a[0] < b[0] ? -1 : 1))
  }, [metrics])

  const limited = useMemo(() => {
    const size = Number(range) || buckets.length
    return buckets.slice(Math.max(0, buckets.length - size))
  }, [buckets, range])
  const categories = limited.map(([bucket]) => bucket)

  const series = useMemo(
    () => [
      {
        name: "提交工单",
        data: limited.map(([, row]) => row["tickets-created"] ?? 0),
      },
      {
        name: "执行成功",
        data: limited.map(([, row]) => row["tickets-executed"] ?? 0),
      },
      {
        name: "执行失败",
        data: limited.map(([, row]) => row["tickets-failed"] ?? 0),
      },
      {
        name: "错误命中",
        data: limited.map(([, row]) => row["violations-error"] ?? 0),
      },
    ],
    [limited],
  )

  return (
    <>
      <PageHeader
        title="统计报表"
        description="按天汇总的工单量与规则命中"
        action={
          <Badge variant="default" className="px-2 py-1">
            <IconChartBar size={14} stroke={1.5} className="mr-1" />
            {buckets.length} 天
          </Badge>
        }
      />

      <SectionCard
        title="工单趋势"
        action={
          <div className="flex items-center gap-2">
            <Label htmlFor="stats-range" className="text-muted-foreground">
              区间
            </Label>
            <Select value={range} onValueChange={(value) => setRange(value)}>
              <SelectTrigger id="stats-range" className="h-8 w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="7">近 7 天</SelectItem>
                <SelectItem value="30">近 30 天</SelectItem>
                <SelectItem value="90">近 90 天</SelectItem>
              </SelectContent>
            </Select>
          </div>
        }
      >
        {categories.length === 0 ? (
          <p className="py-10 text-center text-sm text-muted-foreground">
            {remote.loading ? "正在读取…" : "还没有汇总数据，指标由每日汇总作业写入。"}
          </p>
        ) : (
          <>
            <TrendChart
              id={TREND_CHART_ID}
              categories={categories}
              series={series}
              selectable
              range={7}
            />
            <BrushChart
              // 刷选图只有 96px 高，完整的年月日会把 x 轴挤成负高度。
              categories={categories.map((item) => item.slice(5))}
              series={series.slice(0, 1)}
              target={TREND_CHART_ID}
            />
          </>
        )}
      </SectionCard>

      <SectionCard title="每日明细" bodyClassName="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>日期</TableHead>
              {metricNames.map((name) => (
                <TableHead key={name} className="text-right">
                  {metricLabels[name]}
                </TableHead>
              ))}
            </TableRow>
          </TableHeader>
          <TableBody>
            {[...limited].reverse().length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={metricNames.length + 1}
                  className="py-8 text-center text-sm text-muted-foreground"
                >
                  暂无数据。
                </TableCell>
              </TableRow>
            ) : null}
            {[...limited].reverse().map(([bucket, row]) => (
              <TableRow key={bucket}>
                <TableCell className="font-mono text-xs">{bucket}</TableCell>
                {metricNames.map((name) => (
                  <TableCell key={name} className="text-right tabular-nums">
                    {row[name] ?? 0}
                  </TableCell>
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </SectionCard>
    </>
  )
}
