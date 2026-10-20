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

export function TrendChart({ categories, series, height = 280, className }: TrendChartProps) {
  const options = useMemo<ApexOptions>(
    () => ({
      chart: { toolbar: { show: false }, zoom: { enabled: false } },
      colors: [chartPalette.primary, chartPalette.success],
      dataLabels: { enabled: false },
      stroke: { width: 1.5, curve: "smooth" },
      fill: { type: "gradient", gradient: { opacityFrom: 0.2, opacityTo: 0.02 } },
      legend: { position: "top", horizontalAlign: "right" },
      grid: { borderColor: "transparent" },
      xaxis: { categories, labels: { style: { colors: chartPalette.textMuted } } },
      yaxis: { labels: { style: { colors: chartPalette.textMuted } } },
    }),
    [categories],
  )

  return <ApexChart type="area" series={series} options={options} height={height} className={className} />
}
