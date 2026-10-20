import { useMemo } from "react"
import type { ApexOptions } from "apexcharts"

import { ApexChart } from "@/components/charts/apex-chart"
import { chartPalette } from "@/components/charts/palette"
import type { SeverityBreakdown } from "@/lib/api/types"

const severityMeta: Record<SeverityBreakdown["severity"], { label: string; color: string }> = {
  critical: { label: "致命", color: chartPalette.danger },
  high: { label: "高危", color: chartPalette.warning },
  medium: { label: "中危", color: chartPalette.primary },
  low: { label: "低危", color: chartPalette.neutral },
}

export interface SeverityBreakdownChartProps {
  data: SeverityBreakdown[]
  height?: number
}

export function SeverityBreakdownChart({ data, height = 240 }: SeverityBreakdownChartProps) {
  const options = useMemo<ApexOptions>(
    () => ({
      chart: { toolbar: { show: false } },
      colors: data.map((item) => severityMeta[item.severity].color),
      dataLabels: { enabled: false },
      legend: { position: "bottom" },
      plotOptions: { bar: { horizontal: true, borderRadius: 4, barHeight: "60%" } },
      grid: { borderColor: "transparent" },
      xaxis: {
        categories: data.map((item) => severityMeta[item.severity].label),
        labels: { style: { colors: chartPalette.textMuted } },
      },
      yaxis: { labels: { style: { colors: chartPalette.textMuted } } },
    }),
    [data],
  )

  const series = useMemo(() => [{ name: "工单数", data: data.map((item) => item.count) }], [data])

  return <ApexChart type="bar" series={series} options={options} height={height} />
}
