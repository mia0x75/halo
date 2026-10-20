import { useMemo } from "react"

import { BarChart } from "@/components/charts/bar-chart"
import { chartPalette } from "@/components/charts/palette"
import { severity } from "@/lib/mock/data"

const severityColors: Record<string, string> = {
  致命: chartPalette.danger,
  高危: chartPalette.warning,
  中危: chartPalette.primary,
  低危: chartPalette.neutral,
}

/** 风险分级：命中数量最多的几条规则。 */
export function TopRules() {
  const categories = useMemo(() => severity.map((item) => item.severity), [])
  const series = useMemo(
    () => [{ name: "命中数", data: severity.map((item) => item.count) }],
    [],
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
