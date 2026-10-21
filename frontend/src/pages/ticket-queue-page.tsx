import { useMemo, useState } from "react"
import { IconPlus } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { StatusBadge } from "@/components/shared/status-badge"
import { Button } from "@/components/ui/button"
import { CursorPagination } from "@/components/ui/cursor-pagination"
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
import type { Ticket } from "@/lib/model"

const pageSize = 8

export interface TicketQueuePageProps {
  /** 我的工单 / 待我审核 */
  mode: "mine" | "review"
  onOpenTicket: (uuid: string) => void
  onCreateTicket?: () => void
}

/** 工单队列页：我的工单与待我审核共用一套表格与游标分页。 */
export function TicketQueuePage({ mode, onOpenTicket, onCreateTicket }: TicketQueuePageProps) {
  const [after, setAfter] = useState<string | null>(null)
  const [stack, setStack] = useState<(string | null)[]>([])

  const query = mode === "mine" ? MY_TICKETS_QUERY : PENDING_REVIEWS_QUERY
  const remote = useGraphQL<
    { me: { MyTickets?: TicketConnectionNode | null; PendingReviews?: TicketConnectionNode | null } | null },
    { first: number; after?: string }
  >(query, { first: pageSize, after: after ?? undefined }, { me: null })

  const connection = mode === "mine" ? remote.data?.me?.MyTickets : remote.data?.me?.PendingReviews

  const rows = useMemo<Ticket[]>(
    () => connection?.edges.map((edge) => mapTicket(edge.node)) ?? [],
    [connection],
  )

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
          mode === "mine" && onCreateTicket ? (
            <Button onClick={onCreateTicket}>
              <IconPlus size={16} stroke={1.5} />
              新建工单
            </Button>
          ) : null
        }
      />

      <SectionCard title={title} description={`共 ${total} 条`} bodyClassName="p-0">
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
              <TableRow
                key={ticket.uuid}
                className="cursor-pointer"
                onClick={() => onOpenTicket(ticket.uuid)}
              >
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
