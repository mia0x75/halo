import ApexCharts from "apexcharts"
import type { ApexOptions } from "apexcharts"
import { useEffect, useRef } from "react"

import { cn } from "@/lib/utils"

export type ApexChartType = NonNullable<ApexOptions["chart"]>["type"]
export type ApexSeries = NonNullable<ApexOptions["series"]>

export interface ApexChartProps {
  type: ApexChartType
  series: ApexSeries
  options?: ApexOptions
  height?: number
  className?: string
}

/**
 * Owns the ApexCharts lifecycle so features never touch the library directly.
 * Grid, axis and tooltip colors follow the theme through CSS overrides in
 * index.css, therefore re-rendering on theme change is unnecessary.
 */
export function ApexChart({ type, series, options, height = 240, className }: ApexChartProps) {
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const container = containerRef.current
    if (!container) return

    const chart = new ApexCharts(container, {
      ...options,
      series,
      chart: { ...options?.chart, type, height },
    })

    void chart.render()

    return () => {
      chart.destroy()
    }
  }, [options, series, type, height])

  return <div ref={containerRef} className={cn("w-full", className)} />
}
