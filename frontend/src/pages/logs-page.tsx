import { useCallback, useMemo, useState } from "react"
import { IconSearch } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
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
import { useGraphQL } from "@/hooks/use-graphql"
import { formatUnix } from "@/lib/api/mappers"
import { LOGS_QUERY, type LogNode } from "@/lib/api/queries"
import { parseLogOperation } from "@/lib/labels"

const pageSize = 10

interface LogView {
  id: number
  actor: string
  label: string
  target: string
  raw: string
  createdAt: number
}

export function LogsPage() {
  const [keyword, setKeyword] = useState("")
  const [actor, setActor] = useState("all")
  const [after, setAfter] = useState<string | null>(null)
  const [stack, setStack] = useState<(string | null)[]>([])

  const remote = useGraphQL<
    {
      logs: {
        pageInfo: { hasNextPage: boolean; hasPreviousPage: boolean; startCursor: string; endCursor: string }
        edges: { cursor: string; node: LogNode }[]
      } | null
    },
    { first: number; after?: string }
  >(LOGS_QUERY, { first: pageSize, after: after ?? undefined }, { logs: null })

  const connection = remote.data?.logs ?? null

  const views = useMemo<LogView[]>(
    () =>
      (connection?.edges ?? [])
        // A single malformed edge must not blank the whole page.
        .filter((edge): edge is { cursor: string; node: LogNode } => edge?.node != null)
        .map((edge, index) => {
          const parsed = parseLogOperation(edge.node.Operation)
          return {
            id: edge.node.LogID || index,
            actor: edge.node.User?.Name || edge.node.User?.Email || "已删除用户",
            label: parsed.label,
            target: parsed.target,
            raw: edge.node.Operation,
            createdAt: edge.node.CreateAt,
          }
        }),
    [connection],
  )

  const actors = useMemo(() => Array.from(new Set(views.map((log) => log.actor))), [views])

  const filtered = useMemo(() => {
    const needle = keyword.trim().toLowerCase()
    return views.filter((log) => {
      if (actor !== "all" && log.actor !== actor) return false
      if (!needle) return true
      return (
        log.label.toLowerCase().includes(needle) ||
        log.target.toLowerCase().includes(needle) ||
        log.raw.toLowerCase().includes(needle)
      )
    })
  }, [views, keyword, actor])

  const next = useCallback(() => {
    setStack((items) => [...items, after])
    setAfter(connection?.pageInfo.endCursor || null)
  }, [after, connection])

  const previous = useCallback(() => {
    const previousCursor = stack[stack.length - 1] ?? null
    setStack((items) => items.slice(0, -1))
    setAfter(previousCursor)
  }, [stack])

  // Page number tracks how many times the user paged forward in this session.
  // No COUNT(*) is issued: it would be slow on a large filtered table.
  const page = stack.length + 1

  const reset = useCallback(() => {
    setAfter(null)
    setStack([])
  }, [])

  const filteredEmpty = filtered.length === 0
  const hasFilter = keyword.trim() !== "" || actor !== "all"

  return (
    <>
      <PageHeader
        title="操作日志"
        description="系统里的关键操作都留了痕，登录、提交、审核、执行都在内"
      />

      <SectionCard bodyClassName="p-0">
        <div className="flex flex-wrap items-end gap-3 border-b border-border px-5 py-4">
          <div className="flex flex-col gap-1.5 sm:w-64">
            <Label htmlFor="log-keyword">关键词</Label>
            <div className="relative">
              <IconSearch
                size={16}
                stroke={1.5}
                className="pointer-events-none absolute top-1/2 left-2.5 -translate-y-1/2 text-muted-foreground"
              />
              <Input
                id="log-keyword"
                value={keyword}
                onChange={(event) => setKeyword(event.target.value)}
                placeholder="操作、对象或原始内容"
                className="pl-8"
              />
            </div>
          </div>
          <div className="flex flex-col gap-1.5 sm:w-48">
            <Label htmlFor="log-actor">操作人</Label>
            <Select value={actor} onValueChange={(value) => setActor(value)}>
              <SelectTrigger id="log-actor">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">全部操作人</SelectItem>
                {actors.map((item) => (
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
              <TableHead>操作人</TableHead>
              <TableHead>操作</TableHead>
              <TableHead>对象</TableHead>
              <TableHead className="w-40">时间</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filtered.map((log) => (
              <TableRow key={`${log.id}-${log.createdAt}`}>
                <TableCell className="font-medium">{log.actor}</TableCell>
                <TableCell>{log.label}</TableCell>
                <TableCell className="text-muted-foreground">
                  {log.target ? <span className="font-mono text-xs">{log.target}</span> : "-"}
                </TableCell>
                <TableCell className="text-muted-foreground">{formatUnix(log.createdAt)}</TableCell>
              </TableRow>
            ))}
            {filteredEmpty ? (
              <TableRow>
                <TableCell colSpan={4} className="py-10 text-center text-sm text-muted-foreground">
                  {hasFilter ? "没有符合条件的记录，试试调整筛选条件" : "暂无操作记录"}
                </TableCell>
              </TableRow>
            ) : null}
          </TableBody>
        </Table>

        <CursorPagination
          page={page}
          empty={views.length === 0}
          hasPreviousPage={connection?.pageInfo.hasPreviousPage ?? false}
          hasNextPage={connection?.pageInfo.hasNextPage ?? false}
          onPrevious={previous}
          onNext={next}
          onReset={reset}
        />
      </SectionCard>
    </>
  )
}
