import { useMemo } from "react"

import { BarChart } from "@/components/charts/bar-chart"
import { chartPalette } from "@/components/charts/palette"
import type { SeverityPoint, TopRulePoint } from "@/lib/model"

const severityColors: Record<string, string> = {
  错误: chartPalette.danger,
  警告: chartPalette.warning,
  提示: chartPalette.info,
}

export interface SeverityBodyProps {
  /** 按级别统计的命中数 */
  levels: SeverityPoint[]
}

/** 风险分级：级别分布的柱状图。 */
export function SeverityBody({ levels }: SeverityBodyProps) {
  const categories = useMemo(() => levels.map((item) => item.severity), [levels])
  const series = useMemo(
    () => [{ name: "命中数", data: levels.map((item) => item.count) }],
    [levels],
  )
  const colors = useMemo(
    () => levels.map((item) => severityColors[item.severity] ?? chartPalette.primary),
    [levels],
  )

  if (categories.length === 0) {
    return <p className="px-5 py-6 text-center text-sm text-muted-foreground">该时段没有规则命中</p>
  }

  return (
    <div className="px-2 py-4">
      <BarChart categories={categories} series={series} colors={colors} height={160} />
    </div>
  )
}

export interface TopRulesBodyProps {
  /** 命中最多的规则 */
  rules: TopRulePoint[]
}

/** 命中最多的规则：按命中数降序的列表。 */
export function TopRulesBody({ rules }: TopRulesBodyProps) {
  if (rules.length === 0) {
    return <p className="px-5 py-6 text-center text-sm text-muted-foreground">暂无命中记录</p>
  }

  return (
    <ul className="divide-y divide-border">
      {rules.map((item) => (
        <li key={item.rule} className="flex items-center gap-2 px-5 py-2 text-sm">
          <span className="font-mono text-sm text-muted-foreground">{item.rule}</span>
          <span className="ml-auto font-medium">{item.count}</span>
        </li>
      ))}
    </ul>
  )
}
