import { useMemo } from "react"

import { ApexChart } from "@/components/charts/apex-chart"
import { chartPalette } from "@/components/charts/palette"

export interface SparklineProps {
  data: number[]
  color?: string
  height?: number
  className?: string
}

export function Sparkline({ data, color = chartPalette.primary, height = 48, className }: SparklineProps) {
  const options = useMemo(
    () => ({
      chart: { sparkline: { enabled: true } },
      colors: [color],
      stroke: { width: 1.5 },
      fill: { type: "gradient", gradient: { opacityFrom: 0.24, opacityTo: 0.02 } },
      tooltip: { enabled: false },
    }),
    [color],
  )

  const series = useMemo(() => [{ name: "trend", data }], [data])

  return <ApexChart type="area" series={series} options={options} height={height} className={className} />
}
