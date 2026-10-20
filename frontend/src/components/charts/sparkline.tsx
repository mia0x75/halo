import { useMemo } from "react"
import type { ApexOptions } from "apexcharts"

import { ApexChart } from "@/components/charts/apex-chart"
import { chartPalette } from "@/components/charts/palette"

export interface SparklineProps {
  data: number[]
  color?: string
  height?: number
  className?: string
}

/** 迷你趋势图：用于指标卡。 */
export function Sparkline({ data, color = chartPalette.primary, height = 40, className }: SparklineProps) {
  const options = useMemo(
    () =>
      ({
        chart: { sparkline: { enabled: true } },
        colors: [color],
        stroke: { width: 1.5, curve: "smooth" },
        // 指标卡里的迷你趋势图面积用更实的渐变，颜色更重
        fill: { type: "gradient", gradient: { opacityFrom: 0.5, opacityTo: 0.08 } },
        tooltip: { enabled: false },
      }) as ApexOptions,
    [color],
  )

  const series = useMemo(() => [{ name: "trend", data }], [data])

  return <ApexChart type="area" series={series} options={options} height={height} className={className} />
}
