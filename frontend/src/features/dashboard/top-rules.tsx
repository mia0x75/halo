import { useMemo } from "react"

import { BarChart } from "@/components/charts/bar-chart"
import { chartPalette } from "@/components/charts/palette"
import { severity, trendFor, type Period } from "@/lib/mock/data"

const severityColors: Record<string, string> = {
  致命: chartPalette.danger,
  高危: chartPalette.warning,
  中危: chartPalette.primary,
  低危: chartPalette.neutral,
}

export interface TopRulesProps {
  period: Period
}

/** 风险分级：按时段缩放命中数量最多的几条规则。 */
export function TopRules({ period }: TopRulesProps) {
  const scale = Math.max(1, Math.round(trendFor(period).length / 7))

  const categories = useMemo(() => severity.map((item) => item.severity), [])
  const series = useMemo(
    () => [{ name: "命中数", data: severity.map((item) => item.count * scale) }],
    [scale],
  )
  const colors = useMemo(
    () => severity.map((item) => severityColors[item.severity] ?? chartPalette.primary),
    [],
  )

  return (
    <div className="px-2 py-4">
      <BarChart categories={categories} series={series} colors={colors} height={220} />
    </div>
  )
}
