import { useMemo } from "react"
import type { ApexOptions } from "apexcharts"

import { ApexChart } from "@/components/charts/apex-chart"
import { chartPalette } from "@/components/charts/palette"

export interface BrushChartProps {
  categories: string[]
  series: { name: string; data: number[] }[]
  /** 主图所在的 DOM id，刷选后主图跟着缩放到选中的区间 */
  target: string
  height?: number
  className?: string
}

/**
 * 范围刷选图。设计稿 charts-advanced 的 "Brush and zoom"：下方小图固定展示
 * 较长时间范围（默认一年），图上有一个可拖动、可伸缩的半透明滑块；滑块选中的
 * 区间通过 chart.brush.target 驱动上方主图缩放。
 */
export function BrushChart({ categories, series, target, height = 96, className }: BrushChartProps) {
  const options = useMemo<ApexOptions>(
    () => ({
      chart: {
        toolbar: { show: false },
        zoom: { enabled: false },
        // v8 的 brush 配置在 chart 之下；主图要先渲染出来。
        brush: { enabled: true, target, autoScaleYaxis: false },
      },
      colors: [chartPalette.primary],
      dataLabels: { enabled: false },
      stroke: { width: 1.5, curve: "smooth" },
      fill: { type: "gradient", gradient: { opacityFrom: 0.3, opacityTo: 0.04 } },
      markers: { size: 0, hover: { size: 8 } },
      tooltip: { enabled: false },
      grid: { borderColor: "transparent", padding: { left: 4, right: 4 } },
      xaxis: {
        categories,
        labels: { style: { colors: chartPalette.textMuted, fontSize: "10px" } },
        tooltip: { enabled: false },
      },
      yaxis: { labels: { style: { colors: chartPalette.textMuted, fontSize: "10px" } } },
    }),
    [categories, target],
  )

  return <ApexChart type="area" series={series} options={options} height={height} className={className} />
}
