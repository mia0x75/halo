import { useMemo, useState } from "react"
import { IconSearch } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { StatusBadge } from "@/components/shared/status-badge"
import { Button } from "@/components/ui/button"
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
import { mapTicket } from "@/lib/api/mappers"
import {
  MY_TICKETS_QUERY,
  PENDING_REVIEWS_QUERY,
  type TicketConnectionNode,
} from "@/lib/api/queries"
import { tickets as mockTickets, type Ticket } from "@/lib/mock/data"

const pageSize = 8

export interface TicketQueuePageProps {
  /** 我的工单 / 待我审核 */
  mode: "mine" | "review"
}

/** 工单队列页：我的工单与待我审核共用一套表格与游标分页。 */
export function TicketQueuePage({ mode }: TicketQueuePageProps) {
  const [keyword, setKeyword] = useState("")
  const [after, setAfter] = useState<string | null>(null)
  const [stack, setStack] = useState<(string | null)[]>([])

  const query = mode === "mine" ? MY_TICKETS_QUERY : PENDING_REVIEWS_QUERY
  const remote = useGraphQL<
    { me: { MyTickets?: TicketConnectionNode | null; PendingReviews?: TicketConnectionNode | null } | null },
    { first: number; after?: string }
  >(query, { first: pageSize, after: after ?? undefined }, { me: null })

  const connection = mode === "mine" ? remote.data?.me?.MyTickets : remote.data?.me?.PendingReviews

  const rows = useMemo<Ticket[]>(() => {
    if (!connection) {
      // 接口不可用时回落到假数据：我的工单取全部，待我审核取待审核的部分
      const fallback = mode === "mine" ? mockTickets : mockTickets.filter((t) => t.status === "waiting-mrv")
      const needle = keyword.trim().toLowerCase()
      if (!needle) return fallback.slice(0, pageSize)
      return fallback
        .filter(
          (ticket) =>
            ticket.subject.toLowerCase().includes(needle) ||
            ticket.database.toLowerCase().includes(needle),
        )
        .slice(0, pageSize)
    }
    return connection.edges.map((edge) => mapTicket(edge.node))
  }, [connection, keyword, mode])

  const total = connection?.totalCount ?? rows.length
  const hasNext = connection?.pageInfo.hasNextPage ?? false
  const hasPrevious = connection?.pageInfo.hasPreviousPage ?? false
  const endCursor = connection?.pageInfo.endCursor ?? ""

  const title = mode === "mine" ? "我的工单" : "待我审核"
  const description =
    mode === "mine" ? "我提交与参与处理的变更工单" : "等待我人工审核的工单"

  return (
    <>
      <PageHeader
        title={title}
        description={description}
        action={
          mode === "mine" ? (
            <Button>新建工单</Button>
          ) : (
            <Button variant="outline">批量审核</Button>
          )
        }
      />

      <SectionCard title={title} description={`共 ${total} 条`} bodyClassName="p-0">
        <div className="flex flex-col gap-3 border-b border-border px-5 py-4">
          <div className="flex flex-col gap-1.5 sm:w-72">
            <Label htmlFor={`${mode}-keyword`}>关键词</Label>
            <div className="relative">
              <IconSearch
                size={16}
                stroke={1.5}
                className="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-muted-foreground"
              />
              <Input
                id={`${mode}-keyword`}
                value={keyword}
                onChange={(event) => {
                  setKeyword(event.target.value)
                  setAfter(null)
                  setStack([])
                }}
                placeholder="主题或目标库"
                className="pl-9"
              />
            </div>
          </div>
        </div>

        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>主题</TableHead>
              <TableHead>目标库</TableHead>
              <TableHead>群集</TableHead>
              <TableHead>状态</TableHead>
              <TableHead>提交人</TableHead>
              <TableHead>创建时间</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rows.map((ticket) => (
              <TableRow key={ticket.uuid}>
                <TableCell className="font-medium">{ticket.subject}</TableCell>
                <TableCell className="text-muted-foreground">{ticket.database}</TableCell>
                <TableCell className="text-muted-foreground">{ticket.cluster}</TableCell>
                <TableCell>
                  <StatusBadge status={ticket.status} />
                </TableCell>
                <TableCell className="text-muted-foreground">{ticket.submitter}</TableCell>
                <TableCell className="text-muted-foreground">{ticket.createdAt}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>

        <CursorPagination
          total={total}
          hasPreviousPage={hasPrevious}
          hasNextPage={hasNext}
          startCursor={connection?.pageInfo.startCursor}
          endCursor={endCursor}
          onPrevious={() => {
            const previousCursor = stack[stack.length - 1] ?? null
            setStack((items) => items.slice(0, -1))
            setAfter(previousCursor)
          }}
          onNext={() => {
            setStack((items) => [...items, after])
            setAfter(endCursor || null)
          }}
        />
      </SectionCard>
    </>
  )
}
