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
  /** DOM id，刷选图通过它指向主图 */
  id?: string
  /** 图表挂载完成回调，供依赖它的刷选图延后初始化 */
  onReady?: () => void
}

/**
 * Owns the ApexCharts lifecycle so features never touch the library directly.
 * Grid, axis and tooltip colors follow the theme through CSS overrides in
 * index.css, therefore re-rendering on theme change is unnecessary.
 */
export function ApexChart({ type, series, options, height = 240, className, id, onReady }: ApexChartProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  // 回调每次渲染可能是新引用，不能放进 effect 依赖，否则图表会重建。
  const onReadyRef = useRef(onReady)
  useEffect(() => {
    onReadyRef.current = onReady
  }, [onReady])

  useEffect(() => {
    const container = containerRef.current
    if (!container) return

    const chart = new ApexCharts(container, {
      ...options,
      series,
      // 图表文字由 SVG 渲染，不继承 CSS 的 font-family，需要单独指定；
      // 放在这里统一设置，各图表无需重复声明。
      chart: {
        fontFamily:
          '"寒蝉团圆体 Round", system-ui, -apple-system, "Segoe UI", Roboto, sans-serif',
        ...options?.chart,
        // brush 通过 chart.id 查找目标实例，DOM id 并不参与查找。
        ...(id ? { id } : {}),
        type,
        height,
      },
    })

    let cancelled = false
    void chart.render().then(() => {
      if (!cancelled) onReadyRef.current?.()
    })

    return () => {
      cancelled = true
      chart.destroy()
    }
  }, [options, series, type, height, id])

  return <div id={id} ref={containerRef} className={cn("w-full", className)} />
}
