export interface ReportHit {
  /** 0 提示 1 错误 2 警告 */
  level: number
  description: string
  rule: string
}

const levelLabels: Record<number, string> = { 0: "提示", 1: "错误", 2: "警告" }

/** 解析语句的命中报告：后端存的是 JSON 数组，这里还原成可读列表。 */
export function parseReport(report: string | null | undefined): ReportHit[] {
  if (!report) return []
  try {
    const parsed: unknown = JSON.parse(report)
    if (!Array.isArray(parsed)) return []
    return parsed.flatMap((item) => {
      if (typeof item !== "object" || item === null) return []
      // 后端存的是 Go 字段名（Level/Description/Rule），大小写都兼容一下
      const hit = item as Record<string, unknown>
      const level = hit.Level ?? hit.level
      const description = hit.Description ?? hit.description
      const rule = hit.Rule ?? hit.rule
      return [
        {
          level: typeof level === "number" ? level : 0,
          description: typeof description === "string" ? description : "",
          rule: typeof rule === "string" ? rule : "",
        },
      ]
    })
  } catch {
    return []
  }
}

export function reportLevelLabel(level: number): string {
  return levelLabels[level] ?? "提示"
}
