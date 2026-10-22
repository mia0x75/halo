import { useMemo, useState } from "react"
import { IconSearch, IconShieldCheck } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { CursorPagination } from "@/components/ui/cursor-pagination"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useCursorPagination } from "@/hooks/use-cursor-pagination"
import { useGraphQL } from "@/hooks/use-graphql"
import { RULES_QUERY, RULE_HIT_RULES_QUERY, type RuleNode } from "@/lib/api/queries"
import { ruleGroupOf, ruleLevelLabels, ruleScopeOf } from "@/lib/labels"

const pageSize = 10

const levelTone = {
  错误: "danger",
  警告: "warning",
  提示: "info",
} as const

interface RuleView {
  id: number
  name: string
  group: string
  level: string
  description: string
  hitCount: number
  scope: string
  message: string
}

export function RulesPage() {
  const [keyword, setKeyword] = useState("")
  const [level, setLevel] = useState("all")
  const [group, setGroup] = useState("all")

  const remote = useGraphQL<{ rules: RuleNode[] }>(RULES_QUERY, undefined, { rules: [] })
  const hitCounts = useGraphQL<
    { ruleHitRules: { RuleName: string; Count: number }[] },
    { first: number }
  >(RULE_HIT_RULES_QUERY, { first: 300 }, { ruleHitRules: [] })

  const countByRule = useMemo(() => {
    const map = new Map<string, number>()
    for (const item of hitCounts.data?.ruleHitRules ?? []) {
      map.set(item.RuleName, item.Count)
    }
    return map
  }, [hitCounts.data])

  const views = useMemo<RuleView[]>(
    () =>
      (remote.data?.rules ?? []).map((rule, index) => ({
        id: index,
        name: rule.Name,
        group: ruleGroupOf(rule.Name),
        level: ruleLevelLabels[rule.Level] ?? "提示",
        description: rule.Description,
        hitCount: countByRule.get(rule.Name) ?? 0,
        scope: ruleScopeOf(rule.Bitwise),
        message: rule.Message,
      })),
    [remote.data, countByRule],
  )

  const groups = useMemo(() => Array.from(new Set(views.map((rule) => rule.group))), [views])

  const filtered = useMemo(() => {
    const needle = keyword.trim().toLowerCase()
    return views.filter((rule) => {
      if (level !== "all" && rule.level !== level) return false
      if (group !== "all" && rule.group !== group) return false
      if (!needle) return true
      return (
        rule.name.toLowerCase().includes(needle) ||
        rule.description.toLowerCase().includes(needle)
      )
    })
  }, [views, keyword, level, group])

  const { page, reset, next, previous } = useCursorPagination(filtered, pageSize)

  return (
    <>
      <PageHeader
        title="规则库"
        description="校验规则的清单，命中明细按规则与级别聚合"
        action={
          <Badge variant="default" className="px-2 py-1">
            <IconShieldCheck size={14} stroke={1.5} className="mr-1" />
            共 {views.length} 条规则
          </Badge>
        }
      />

      <SectionCard bodyClassName="p-0">
        <div className="flex flex-wrap items-end gap-3 border-b border-border px-5 py-4">
          <div className="flex flex-col gap-1.5 sm:w-64">
            <Label htmlFor="rule-keyword">关键词</Label>
            <div className="relative">
              <IconSearch
                size={16}
                stroke={1.5}
                className="pointer-events-none absolute top-1/2 left-2.5 -translate-y-1/2 text-muted-foreground"
              />
              <Input
                id="rule-keyword"
                value={keyword}
                onChange={(event) => {
                  setKeyword(event.target.value)
                  reset()
                }}
                placeholder="规则名或说明"
                className="pl-8"
              />
            </div>
          </div>
          <div className="flex flex-col gap-1.5 sm:w-40">
            <Label htmlFor="rule-level">级别</Label>
            <Select
              value={level}
              onValueChange={(value) => {
                setLevel(value)
                reset()
              }}
            >
              <SelectTrigger id="rule-level">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">全部级别</SelectItem>
                <SelectItem value="错误">错误</SelectItem>
                <SelectItem value="警告">警告</SelectItem>
                <SelectItem value="提示">提示</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="flex flex-col gap-1.5 sm:w-44">
            <Label htmlFor="rule-group">分组</Label>
            <Select
              value={group}
              onValueChange={(value) => {
                setGroup(value)
                reset()
              }}
            >
              <SelectTrigger id="rule-group">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">全部分组</SelectItem>
                {groups.map((item) => (
                  <SelectItem key={item} value={item}>
                    {item}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>规则</TableHead>
              <TableHead>分组</TableHead>
              <TableHead>级别</TableHead>
              <TableHead>说明</TableHead>
              <TableHead>适用范围</TableHead>
              <TableHead className="text-right">命中</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {page.edges.map((rule) => (
              <TableRow key={rule.name}>
                <TableCell className="font-mono text-xs font-medium">{rule.name}</TableCell>
                <TableCell className="text-muted-foreground">{rule.group}</TableCell>
                <TableCell>
                  <Badge variant={levelTone[rule.level as keyof typeof levelTone] ?? "muted"}>
                    {rule.level}
                  </Badge>
                </TableCell>
                <TableCell className="max-w-md">
                  <span className="block truncate" title={rule.description}>
                    {rule.description}
                  </span>
                </TableCell>
                <TableCell className="text-muted-foreground">{rule.scope}</TableCell>
                <TableCell className="text-right tabular-nums">{rule.hitCount}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <CursorPagination
          total={filtered.length}
          hasPreviousPage={page.pageInfo.hasPreviousPage}
          hasNextPage={page.pageInfo.hasNextPage}
          onPrevious={previous}
          onNext={next}
        />
      </SectionCard>
    </>
  )
}
