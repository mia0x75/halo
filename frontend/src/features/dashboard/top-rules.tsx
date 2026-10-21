import { useMemo } from "react"

import { BarChart } from "@/components/charts/bar-chart"
import { chartPalette } from "@/components/charts/palette"
import type { SeverityPoint, TopRulePoint } from "@/lib/model"

const severityColors: Record<string, string> = {
  错误: chartPalette.danger,
  警告: chartPalette.warning,
  提示: chartPalette.info,
}

export interface TopRulesProps {
  /** 按级别统计的命中数 */
  levels: SeverityPoint[]
  /** 命中最多的规则 */
  rules: TopRulePoint[]
}

/** 风险分级：级别分布 + 命中最多的规则。 */
export function TopRules({ levels, rules }: TopRulesProps) {
  const categories = useMemo(() => levels.map((item) => item.severity), [levels])
  const series = useMemo(
    () => [{ name: "命中数", data: levels.map((item) => item.count) }],
    [levels],
  )
  const colors = useMemo(
    () => levels.map((item) => severityColors[item.severity] ?? chartPalette.primary),
    [levels],
  )

  return (
    <div className="flex flex-col gap-4 px-2 py-4">
      {categories.length > 0 ? (
        <BarChart categories={categories} series={series} colors={colors} height={160} />
      ) : (
        <p className="px-3 py-6 text-center text-sm text-muted-foreground">该时段没有规则命中</p>
      )}

      <div>
        <p className="px-3 pb-1.5 text-sm font-medium text-muted-foreground">命中最多的规则</p>
        <ul className="flex flex-col">
          {rules.length === 0 ? (
            <li className="border-t border-border px-3 py-2 text-sm text-muted-foreground">
              暂无命中记录
            </li>
          ) : null}
          {rules.map((item) => (
            <li
              key={item.rule}
              className="flex items-center gap-2 border-t border-border px-3 py-2 text-sm"
            >
              <span className="font-mono text-sm text-muted-foreground">{item.rule}</span>
              <span className="ml-auto font-medium">{item.count}</span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  )
}
