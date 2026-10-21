import { useCallback, useMemo, useState } from "react"
import { IconSearch } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { CursorPagination } from "@/components/ui/cursor-pagination"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useGraphQL } from "@/hooks/use-graphql"
import { formatUnix } from "@/lib/api/mappers"
import {
  STATEMENTS_QUERY,
  TICKETS_QUERY,
  type StatementNode,
  type TicketNode,
} from "@/lib/api/queries"
import { parseReport, reportLevelLabel } from "@/features/tickets/ticket-report"

/** 语句状态跟工单状态一个口径。 */
const statusLabels: Record<number, string> = {
  1: "待系统审核",
  2: "系统审核不通过",
  3: "系统审核有警告",
  4: "待人工审核",
  5: "人工审核不通过",
  6: "审核通过",
  7: "执行完成",
  8: "执行失败",
  9: "已关闭",
}

/** 语句管理页：所有工单拆分出来的 SQL 语句。 */
export function StatementsPage() {
  const [keyword, setKeyword] = useState("")
  const [after, setAfter] = useState<string | null>(null)
  const [stack, setStack] = useState<(string | null)[]>([])

  const remote = useGraphQL<
    {
      statements: {
        totalCount: number
        pageInfo: { hasNextPage: boolean; hasPreviousPage: boolean; startCursor: string; endCursor: string }
        edges: { cursor: string; node: StatementNode }[]
      } | null
    },
    { first: number; after?: string }
  >(STATEMENTS_QUERY, { first: 10, after: after ?? undefined }, { statements: null })

  const ticketList = useGraphQL<
    { tickets: { edges: { node: TicketNode }[] } | null },
    { first: number }
  >(TICKETS_QUERY, { first: 100 }, { tickets: null })

  // 语句只带工单的自增 ID，主题要从工单列表里对出来。
  const ticketById = useMemo(() => {
    const map = new Map<number, string>()
    for (const edge of ticketList.data?.tickets?.edges ?? []) {
      map.set(edge.node.TicketID, edge.node.Subject)
    }
    return map
  }, [ticketList.data])

  const connection = remote.data?.statements ?? null
  const rows = useMemo(() => {
    const needle = keyword.trim().toLowerCase()
    return (connection?.edges ?? [])
      .map((edge) => edge.node)
      .filter((node) => {
        if (!needle) return true
        const subject = ticketById.get(node.TicketID) ?? ""
        return (
          node.Content.toLowerCase().includes(needle) ||
          subject.toLowerCase().includes(needle)
        )
      })
  }, [connection, keyword, ticketById])

  const nextPage = useCallback(() => {
    setStack((items) => [...items, after])
    setAfter(connection?.pageInfo.endCursor || null)
  }, [after, connection])

  const previousPage = useCallback(() => {
    const previousCursor = stack[stack.length - 1] ?? null
    setStack((items) => items.slice(0, -1))
    setAfter(previousCursor)
  }, [stack])

  return (
    <>
      <PageHeader
        title="语句管理"
        description="工单拆分出来的 SQL 语句与规则命中"
        action={
          <Badge variant="default" className="px-2 py-1">
            共 {connection?.totalCount ?? 0} 条
          </Badge>
        }
      />

      <SectionCard title="语句列表" description="最新的在最前面" bodyClassName="p-0">
        <div className="flex flex-wrap items-end gap-3 border-b border-border px-5 py-4">
          <div className="flex flex-col gap-1.5 sm:w-72">
            <Label htmlFor="statement-keyword">关键词</Label>
            <div className="relative">
              <IconSearch
                size={16}
                stroke={1.5}
                className="pointer-events-none absolute top-1/2 left-2.5 -translate-y-1/2 text-muted-foreground"
              />
              <Input
                id="statement-keyword"
                value={keyword}
                onChange={(event) => setKeyword(event.target.value)}
                placeholder="SQL 内容或工单主题"
                className="pl-8"
              />
            </div>
          </div>
        </div>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>工单</TableHead>
              <TableHead>语句</TableHead>
              <TableHead>类型</TableHead>
              <TableHead>状态</TableHead>
              <TableHead className="w-40">时间</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rows.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="py-8 text-center text-sm text-muted-foreground">
                  {remote.loading ? "正在读取…" : "没有匹配��语句。"}
                </TableCell>
              </TableRow>
            ) : null}
            {rows.map((node) => (
              <TableRow key={node.UUID}>
                <TableCell>
                  <span className="block max-w-40 truncate">
                    {ticketById.get(node.TicketID) ?? `工单 #${node.TicketID}`}
                  </span>
                  <span className="font-mono text-xs text-muted-foreground">
                    #{node.TicketID} · 语句 {node.Sequence}
                  </span>
                </TableCell>
                <TableCell>
                  <pre className="max-w-xl overflow-x-auto font-mono text-xs whitespace-pre-wrap">
                    {node.Content}
                  </pre>
                  {parseReport(node.Report ?? "").map((hit, index) => (
                    <p
                      key={`${node.UUID}-hit-${index}`}
                      className="mt-1 text-xs text-warning"
                    >
                      <span className="mr-1 font-medium">[{reportLevelLabel(hit.level)}]</span>
                      {hit.description}
                    </p>
                  ))}
                </TableCell>
                <TableCell>
                  <Badge variant="muted">{node.TypeDesc}</Badge>
                </TableCell>
                <TableCell className="text-muted-foreground">
                  {statusLabels[node.Status] ?? node.Status}
                </TableCell>
                <TableCell className="text-muted-foreground">{formatUnix(node.CreateAt)}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <CursorPagination
          total={connection?.totalCount ?? 0}
          hasPreviousPage={connection?.pageInfo.hasPreviousPage ?? false}
          hasNextPage={connection?.pageInfo.hasNextPage ?? false}
          onPrevious={previousPage}
          onNext={nextPage}
        />
      </SectionCard>
    </>
  )
}
