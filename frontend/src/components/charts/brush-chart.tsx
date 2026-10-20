import { useMemo } from "react"
import type { ApexOptions } from "apexcharts"

import { ApexChart } from "@/components/charts/apex-chart"
import { chartPalette } from "@/components/charts/palette"

export interface BrushChartProps {
  categories: string[]
  series: { name: string; data: number[] }[]
  /** 主图所在的 DOM id，刷选后主图跟着缩放 */
  target: string
  height?: number
  className?: string
}

/**
 * 范围刷选图。设计稿 charts-advanced 的 "Brush and zoom"：下方的小图负责
 * 拖拽选择区间，区间通过 chart.brush.target 驱动上方主图的缩放。
 */
export function BrushChart({ categories, series, target, height = 90, className }: BrushChartProps) {
  const options = useMemo<ApexOptions>(
    () => ({
      chart: {
        toolbar: { show: false },
        zoom: { enabled: false },
        // 指向主图，刷选后主图缩放到选中的区间
        brush: { target },
      },
      brush: { enabled: true, autoScaleYaxis: false },
      colors: [chartPalette.primary],
      dataLabels: { enabled: false },
      markers: { size: 0, hover: { size: 8 } },
      stroke: { width: 1.5, curve: "smooth" },
      fill: { type: "gradient", gradient: { opacityFrom: 0.25, opacityTo: 0.02 } },
      tooltip: { enabled: false },
      grid: { borderColor: "transparent", padding: { left: 4, right: 4 } },
      xaxis: { categories, labels: { style: { colors: chartPalette.textMuted, fontSize: "10px" } } },
      yaxis: { labels: { style: { colors: chartPalette.textMuted, fontSize: "10px" } } },
    }),
    [categories, target],
  )

  return <ApexChart type="area" series={series} options={options} height={height} className={className} />
}
