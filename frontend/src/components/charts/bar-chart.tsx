import { useMemo } from "react"
import type { ApexOptions } from "apexcharts"

import { ApexChart } from "@/components/charts/apex-chart"
import { chartPalette } from "@/components/charts/palette"

export interface BarChartProps {
  categories: string[]
  series: { name: string; data: number[] }[]
  colors?: string[]
  horizontal?: boolean
  height?: number
  className?: string
}

/** 条形图：默认横向，用于分级分布。 */
export function BarChart({
  categories,
  series,
  colors,
  horizontal = true,
  height = 240,
  className,
}: BarChartProps) {
  const options = useMemo<ApexOptions>(
    () => ({
      chart: { toolbar: { show: false }, fontFamily: "Inter, sans-serif" },
      colors: colors ?? [chartPalette.primary],
      dataLabels: { enabled: false },
      legend: { position: "bottom", fontSize: "12px" },
      plotOptions: { bar: { horizontal, borderRadius: 4, barHeight: "60%" } },
      grid: { borderColor: "transparent" },
      xaxis: { categories, labels: { style: { colors: chartPalette.textMuted, fontSize: "11px" } } },
      yaxis: { labels: { style: { colors: chartPalette.textMuted, fontSize: "11px" } } },
    }),
    [categories, colors, horizontal],
  )

  return <ApexChart type="bar" series={series} options={options} height={height} className={className} />
}
