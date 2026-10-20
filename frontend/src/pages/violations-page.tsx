import { useMemo, useState } from "react"
import { IconSearch } from "@tabler/icons-react"

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
import { violations } from "@/lib/mock/data"

const pageSize = 8

const levelTone = {
  错误: "danger",
  警告: "warning",
  提示: "info",
} as const

export function ViolationsPage() {
  const [keyword, setKeyword] = useState("")
  const [level, setLevel] = useState("all")
  const [rule, setRule] = useState("all")

  const ruleNames = useMemo(
    () => Array.from(new Set(violations.map((violation) => violation.rule))),
    [],
  )

  const filtered = useMemo(() => {
    const needle = keyword.trim().toLowerCase()
    return violations.filter((violation) => {
      if (level !== "all" && violation.level !== level) return false
      if (rule !== "all" && violation.rule !== rule) return false
      if (!needle) return true
      return (
        violation.ticket.toLowerCase().includes(needle) ||
        violation.message.toLowerCase().includes(needle)
      )
    })
  }, [keyword, level, rule])

  const { page, reset, next, previous } = useCursorPagination(filtered, pageSize)

  return (
    <>
      <PageHeader title="命中记录" description="规则命中明细，可按工单与规则追溯" />

      <SectionCard title="命中明细" description={`共 ${filtered.length} 条`} bodyClassName="p-0">
        <div className="flex flex-wrap items-end gap-3 border-b border-border px-5 py-4">
          <div className="flex flex-col gap-1.5 sm:w-72">
            <Label htmlFor="violation-keyword">关键词</Label>
            <div className="relative">
              <IconSearch
                size={16}
                stroke={1.5}
                className="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-muted-foreground"
              />
              <Input
                id="violation-keyword"
                value={keyword}
                onChange={(event) => {
                  setKeyword(event.target.value)
                  reset()
                }}
                placeholder="工单或提示内容"
                className="pl-9"
              />
            </div>
          </div>
          <div className="flex flex-col gap-1.5 sm:w-36">
            <Label htmlFor="violation-level">级别</Label>
            <Select
              id="violation-level"
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
          <div className="flex flex-col gap-1.5 sm:w-48">
            <Label htmlFor="violation-rule">规则</Label>
            <Select
              id="violation-rule"
              value={rule}
              onChange={(event) => {
                setRule(event.target.value)
                reset()
              }}
            >
              <option value="all">全部规则</option>
              {ruleNames.map((name) => (
                <option key={name} value={name}>
                  {name}
                </option>
              ))}
            </Select>
          </div>
        </div>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>规则</TableHead>
              <TableHead>级别</TableHead>
              <TableHead>工单</TableHead>
              <TableHead>提示</TableHead>
              <TableHead>时间</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {page.edges.map((violation) => (
              <TableRow key={violation.id}>
                <TableCell className="font-mono text-xs font-medium">{violation.rule}</TableCell>
                <TableCell>
                  <Badge variant={levelTone[violation.level]}>{violation.level}</Badge>
                </TableCell>
                <TableCell className="text-muted-foreground">{violation.ticket}</TableCell>
                <TableCell>{violation.message}</TableCell>
                <TableCell className="text-muted-foreground">{violation.createdAt}</TableCell>
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
