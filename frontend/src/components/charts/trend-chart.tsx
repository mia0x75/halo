import { useMemo } from "react"
import type { ApexOptions } from "apexcharts"

import { ApexChart } from "@/components/charts/apex-chart"
import { chartPalette } from "@/components/charts/palette"

export interface TrendChartProps {
  categories: string[]
  series: { name: string; data: number[] }[]
  height?: number
  className?: string
  /** 主图 DOM id，供刷选图指向 */
  id?: string
  /** 打开选择工具，配合下方的刷选图缩放 */
  selectable?: boolean
}

/** 趋势图：面积折线，多系列。 */
export function TrendChart({
  categories,
  series,
  height = 280,
  className,
  id,
  selectable = false,
}: TrendChartProps) {
  const options = useMemo<ApexOptions>(
    () => ({
      chart: {
        toolbar: {
          show: false,
          // 刷选图通过 chart.brush.target 驱动主图缩放，主图需要打开选择工具
          tools: { selection: selectable, zoom: false },
        },
        zoom: { enabled: selectable, type: "xy" },
      },
      // 常态不画点，鼠标悬停时在曲线上标出该点
      markers: { size: 0, hover: { size: 8 } },
      colors: [chartPalette.primary, chartPalette.success, chartPalette.danger],
      dataLabels: { enabled: false },
      stroke: { width: 1.5, curve: "smooth" },
      fill: { type: "gradient", gradient: { opacityFrom: 0.2, opacityTo: 0.02 } },
      legend: {
        position: "top",
        horizontalAlign: "right",
        fontSize: "12px",
        // 图例用小圆点，圆点与文字之间留出间距
        markers: { size: 8, shape: "circle" },
        itemMargin: { horizontal: 8, vertical: 0 },
      },
      // 紧凑提示框：单行展示，缩小行距
      tooltip: { compact: true, intersect: false, shared: true },
      grid: { borderColor: "transparent", padding: { left: 4, right: 4 } },
      xaxis: { categories, labels: { style: { colors: chartPalette.textMuted, fontSize: "11px" } } },
      yaxis: { labels: { style: { colors: chartPalette.textMuted, fontSize: "11px" } } },
    }),
    [categories, selectable],
  )

  return (
    <ApexChart
      id={id}
      type="area"
      series={series}
      options={options}
      height={height}
      className={className}
    />
  )
}
