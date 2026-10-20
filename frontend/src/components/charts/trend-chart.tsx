import { useMemo } from "react"
import type { ApexOptions } from "apexcharts"

import { ApexChart } from "@/components/charts/apex-chart"
import { chartPalette } from "@/components/charts/palette"

export interface TrendChartProps {
  categories: string[]
  series: { name: string; data: number[] }[]
  height?: number
  className?: string
}

/** 趋势图：面积折线，多系列。 */
export function TrendChart({ categories, series, height = 280, className }: TrendChartProps) {
  const options = useMemo<ApexOptions>(
    () => ({
      chart: { toolbar: { show: false }, zoom: { enabled: false }, fontFamily: "Inter, sans-serif" },
      colors: [chartPalette.primary, chartPalette.success, chartPalette.danger],
      dataLabels: { enabled: false },
      stroke: { width: 1.5, curve: "smooth" },
      fill: { type: "gradient", gradient: { opacityFrom: 0.2, opacityTo: 0.02 } },
      legend: { position: "top", horizontalAlign: "right", fontSize: "12px" },
      grid: { borderColor: "transparent", padding: { left: 4, right: 4 } },
      xaxis: { categories, labels: { style: { colors: chartPalette.textMuted, fontSize: "11px" } } },
      yaxis: { labels: { style: { colors: chartPalette.textMuted, fontSize: "11px" } } },
    }),
    [categories],
  )

  return <ApexChart type="area" series={series} options={options} height={height} className={className} />
}
