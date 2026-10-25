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
  /** 初始显示区间：末尾 N 个点，不传则显示全部 */
  range?: number
  /** 图表挂载完成回调，供刷选图延后初始化 */
  onReady?: () => void
}

/** 趋势图：面积折线，多系列。 */
export function TrendChart({
  categories,
  series,
  height = 280,
  className,
  id,
  selectable = false,
  range,
  onReady,
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
      // 常态不画点，鼠标悬停时在曲线上标出该点；
      // ApexCharts 的 size 是半径，4 即直径 8px，与图例、提示框一致。
      markers: { size: 0, hover: { size: 4 } },
      colors: [
        chartPalette.primary,
        chartPalette.success,
        chartPalette.danger,
        chartPalette.warning,
        chartPalette.info,
        chartPalette.secondary,
      ],
      dataLabels: { enabled: false },
      stroke: { width: 1.5, curve: "smooth" },
      fill: { type: "gradient", gradient: { opacityFrom: 0.2, opacityTo: 0.02 } },
      legend: {
        position: "top",
        horizontalAlign: "right",
        fontSize: "12px",
        labels: { colors: "var(--halo-text)" },
        // 图例用小圆点，圆点与文字之间留出间距
        markers: { size: 4, shape: "circle" },
        itemMargin: { horizontal: 8, vertical: 0 },
      },
      // 紧凑提示框：单行展示，缩小行距
      tooltip: { compact: true, intersect: false, shared: true, style: { fontSize: "11px" } },
      grid: { borderColor: "transparent", padding: { left: 4, right: 4 } },
      xaxis: {
        categories,
        // range 表示从末尾算起的点数，初始只显示最近的一段；
        // 点数不足时不能照搬，否则整段数据会被挤到一侧。
        ...(range ? { range: Math.min(range, categories.length) } : {}),
        labels: { style: { colors: chartPalette.textMuted, fontSize: "11px" } },
      },
      yaxis: { labels: { style: { colors: chartPalette.textMuted, fontSize: "11px" } } },
    }),
    [categories, selectable, range],
  )

  return (
    <ApexChart
      id={id}
      type="area"
      series={series}
      options={options}
      height={height}
      className={className}
      onReady={onReady}
    />
  )
}
