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
import { useGraphQL } from "@/hooks/use-graphql"
import { formatUnix } from "@/lib/api/mappers"
import {
  RULE_HITS_QUERY,
  TICKETS_QUERY,
  type RuleHitNode,
  type TicketNode,
} from "@/lib/api/queries"
import { ruleLevelLabels } from "@/lib/labels"

const pageSize = 10

const levelTone = {
  错误: "danger",
  警告: "warning",
  提示: "info",
} as const

interface HitView {
  id: number
  rule: string
  level: string
  ticketId: number
  ticketSubject: string
  ticketUUID: string
  sequence: number
  message: string
  createdAt: number
}

export function ViolationsPage() {
  const [keyword, setKeyword] = useState("")
  const [level, setLevel] = useState("all")
  const [rule, setRule] = useState("all")

  const remote = useGraphQL<{ ruleHits: RuleHitNode[] }, { first: number }>(
    RULE_HITS_QUERY,
    { first: 200 },
    { ruleHits: [] },
  )
  const ticketList = useGraphQL<
    { tickets: { edges: { node: TicketNode }[] } | null },
    { first: number }
  >(TICKETS_QUERY, { first: 100 }, { tickets: null })

  // 命中记录只带工单的自增 ID，主题要从工单列表里对出来。
  const ticketById = useMemo(() => {
    const map = new Map<number, { subject: string; uuid: string }>()
    for (const edge of ticketList.data?.tickets?.edges ?? []) {
      map.set(edge.node.TicketID, { subject: edge.node.Subject, uuid: edge.node.UUID })
    }
    return map
  }, [ticketList.data])

  const views = useMemo<HitView[]>(
    () =>
      (remote.data?.ruleHits ?? []).map((hit, index) => {
        const ticket = ticketById.get(hit.TicketID)
        return {
          id: index,
          rule: hit.RuleName,
          level: ruleLevelLabels[hit.Level] ?? "提示",
          ticketId: hit.TicketID,
          ticketSubject: ticket?.subject ?? `工单 #${hit.TicketID}`,
          ticketUUID: ticket?.uuid ?? "",
          sequence: hit.Sequence,
          message: hit.Message,
          createdAt: hit.CreateAt,
        }
      }),
    [remote.data, ticketById],
  )

  const ruleNames = useMemo(() => Array.from(new Set(views.map((hit) => hit.rule))), [views])

  const filtered = useMemo(() => {
    const needle = keyword.trim().toLowerCase()
    return views.filter((hit) => {
      if (level !== "all" && hit.level !== level) return false
      if (rule !== "all" && hit.rule !== rule) return false
      if (!needle) return true
      return (
        hit.rule.toLowerCase().includes(needle) ||
        hit.message.toLowerCase().includes(needle) ||
        hit.ticketSubject.toLowerCase().includes(needle)
      )
    })
  }, [views, keyword, level, rule])

  const { page, reset, next, previous } = useCursorPagination(filtered, pageSize)

  return (
    <>
      <PageHeader
        title="命中记录"
        description="规则命中的明细，按时间倒序"
        action={
          <Badge variant="default" className="px-2 py-1">
            共 {views.length} 条命中
          </Badge>
        }
      />

      <SectionCard title="命中明细" description="每条命中对应工单里的一个语句" bodyClassName="p-0">
        <div className="flex flex-wrap items-end gap-3 border-b border-border px-5 py-4">
          <div className="flex flex-col gap-1.5 sm:w-64">
            <Label htmlFor="hit-keyword">关键词</Label>
            <div className="relative">
              <IconSearch
                size={16}
                stroke={1.5}
                className="pointer-events-none absolute top-1/2 left-2.5 -translate-y-1/2 text-muted-foreground"
              />
              <Input
                id="hit-keyword"
                value={keyword}
                onChange={(event) => {
                  setKeyword(event.target.value)
                  reset()
                }}
                placeholder="规则名、提示或工单主题"
                className="pl-8"
              />
            </div>
          </div>
          <div className="flex flex-col gap-1.5 sm:w-40">
            <Label htmlFor="hit-level">级别</Label>
            <Select
              id="hit-level"
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
            <Label htmlFor="hit-rule">规则</Label>
            <Select
              id="hit-rule"
              value={rule}
              onChange={(event) => {
                setRule(event.target.value)
                reset()
              }}
            >
              <option value="all">全部规则</option>
              {ruleNames.map((item) => (
                <option key={item} value={item}>
                  {item}
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
              <TableHead className="w-40">时间</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {page.edges.map((hit) => (
              <TableRow key={`${hit.rule}-${hit.ticketId}-${hit.sequence}-${hit.id}`}>
                <TableCell className="font-mono text-xs font-medium">{hit.rule}</TableCell>
                <TableCell>
                  <Badge variant={levelTone[hit.level as keyof typeof levelTone] ?? "muted"}>
                    {hit.level}
                  </Badge>
                </TableCell>
                <TableCell>
                  <span className="block max-w-48 truncate" title={hit.ticketSubject}>
                    {hit.ticketSubject}
                  </span>
                  <span className="font-mono text-xs text-muted-foreground">
                    #{hit.ticketId} · 语句 {hit.sequence}
                  </span>
                </TableCell>
                <TableCell className="max-w-md">
                  <span className="block truncate" title={hit.message}>
                    {hit.message}
                  </span>
                </TableCell>
                <TableCell className="text-muted-foreground">
                  {formatUnix(hit.createdAt)}
                </TableCell>
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
