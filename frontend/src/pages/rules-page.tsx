import { useMemo, useState } from "react"
import { IconSearch, IconShieldCheck } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { CursorPagination } from "@/components/ui/cursor-pagination"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useCursorPagination } from "@/hooks/use-cursor-pagination"
import { rules } from "@/lib/mock/data"

const pageSize = 8

const levelTone = {
  错误: "danger",
  警告: "warning",
  提示: "info",
} as const

export function RulesPage() {
  const [keyword, setKeyword] = useState("")
  const [level, setLevel] = useState("all")
  const [group, setGroup] = useState("all")
  const [onlyEnabled, setOnlyEnabled] = useState(false)

  const groups = useMemo(() => Array.from(new Set(rules.map((rule) => rule.group))), [])

  const filtered = useMemo(() => {
    const needle = keyword.trim().toLowerCase()
    return rules.filter((rule) => {
      if (level !== "all" && rule.level !== level) return false
      if (group !== "all" && rule.group !== group) return false
      if (onlyEnabled && !rule.enabled) return false
      if (!needle) return true
      return (
        rule.name.toLowerCase().includes(needle) ||
        rule.description.toLowerCase().includes(needle)
      )
    })
  }, [keyword, level, group, onlyEnabled])

  const { page, reset, next, previous } = useCursorPagination(filtered, pageSize)

  return (
    <>
      <PageHeader
        title="规则库"
        description="校验规则的清单，命中明细按规则与级别聚合"
        action={
          <Badge variant="default" className="px-2 py-1">
            <IconShieldCheck size={14} stroke={1.5} />
            248 条已启用
          </Badge>
        }
      />

      <SectionCard title="规则清单" description={`共 ${filtered.length} 条`} bodyClassName="p-0">
        <div className="flex flex-wrap items-end gap-3 border-b border-border px-5 py-4">
          <div className="flex flex-col gap-1.5 sm:w-64">
            <Label htmlFor="rule-keyword">关键词</Label>
            <div className="relative">
              <IconSearch
                size={16}
                stroke={1.5}
                className="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-muted-foreground"
              />
              <Input
                id="rule-keyword"
                value={keyword}
                onChange={(event) => {
                  setKeyword(event.target.value)
                  reset()
                }}
                placeholder="规则名或说明"
                className="pl-9"
              />
            </div>
          </div>
          <div className="flex flex-col gap-1.5 sm:w-36">
            <Label htmlFor="rule-level">级别</Label>
            <Select
              id="rule-level"
              value={level}
              onChange={(event) => {
                setLevel(event.target.value)
                reset()
              }}
            >
              <option value="all">全部级别</option>
              <option value="错误">错误</option>
              <option value="警告">警告</option>
              <option value="提示">提示</option>
            </Select>
          </div>
          <div className="flex flex-col gap-1.5 sm:w-40">
            <Label htmlFor="rule-group">分组</Label>
            <Select
              id="rule-group"
              value={group}
              onChange={(event) => {
                setGroup(event.target.value)
                reset()
              }}
            >
              <option value="all">全部分组</option>
              {groups.map((item) => (
                <option key={item} value={item}>
                  {item}
                </option>
              ))}
            </Select>
          </div>
          <label className="flex items-center gap-2 pb-2 text-sm">
            <input
              type="checkbox"
              className="size-4 rounded border-border accent-primary"
              checked={onlyEnabled}
              onChange={(event) => {
                setOnlyEnabled(event.target.checked)
                reset()
              }}
            />
            只看已启用
          </label>
        </div>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>规则名</TableHead>
              <TableHead>分组</TableHead>
              <TableHead>级别</TableHead>
              <TableHead>状态</TableHead>
              <TableHead>说明</TableHead>
              <TableHead className="text-right">近 7 天命中</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {page.edges.map((rule) => (
              <TableRow key={rule.name}>
                <TableCell className="font-mono text-xs font-medium">{rule.name}</TableCell>
                <TableCell className="text-muted-foreground">{rule.group}</TableCell>
                <TableCell>
                  <Badge variant={levelTone[rule.level]}>{rule.level}</Badge>
                </TableCell>
                <TableCell className="text-muted-foreground">
                  {rule.enabled ? "已启用" : "已停用"}
                </TableCell>
                <TableCell className="text-muted-foreground">{rule.description}</TableCell>
                <TableCell className="text-right font-medium">{rule.hitCount}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <CursorPagination
          total={page.totalCount}
          hasPreviousPage={page.pageInfo.hasPreviousPage}
          hasNextPage={page.pageInfo.hasNextPage}
          startCursor={page.pageInfo.startCursor}
          endCursor={page.pageInfo.endCursor}
          onPrevious={previous}
          onNext={next}
        />
      </SectionCard>
    </>
  )
}
