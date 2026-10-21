import { useMemo, useState } from "react"
import { IconPlus, IconSearch } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { StatusBadge } from "@/components/shared/status-badge"
import { Button } from "@/components/ui/button"
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
import { useGraphQL } from "@/hooks/use-graphql"

import { mapTicket } from "@/lib/api/mappers"
import { TICKETS_QUERY, type TicketConnectionNode } from "@/lib/api/queries"
import { ticketStatusMeta } from "@/lib/labels"
import type { Ticket, TicketStatus } from "@/lib/model"

const pageSize = 8

export interface TicketsPageProps {
  onOpenTicket: (uuid: string) => void
  onCreateTicket: () => void
  /** 顶栏搜索带进来的关键词 */
  initialKeyword?: string
}

export function TicketsPage({ onOpenTicket, onCreateTicket, initialKeyword }: TicketsPageProps) {
  const [keyword, setKeyword] = useState(initialKeyword ?? "")
  const [status, setStatus] = useState<TicketStatus | "all">("all")
  const [after, setAfter] = useState<string | null>(null)
  const [stack, setStack] = useState<(string | null)[]>([])

  const remote = useGraphQL<
    { tickets: TicketConnectionNode | null },
    { first: number; after?: string; status?: number; keyword?: string }
  >(
    TICKETS_QUERY,
    {
      first: pageSize,
      after: after ?? undefined,
      status: status === "all" ? 0 : statusCode(status),
      keyword: keyword.trim() || undefined,
    },
    { tickets: null },
  )

  const connection = remote.data?.tickets ?? null

  const rows = useMemo<Ticket[]>(
    () => connection?.edges.map((edge) => mapTicket(edge.node)) ?? [],
    [connection],
  )

  const total = connection?.totalCount ?? rows.length
  const hasNext = connection?.pageInfo.hasNextPage ?? false
  const hasPrevious = connection?.pageInfo.hasPreviousPage ?? false
  const endCursor = connection?.pageInfo.endCursor ?? ""

  const reset = () => {
    setAfter(null)
    setStack([])
  }

  return (
    <>
      <PageHeader
        title="工单"
        description="变更工单的提交与审核状态"
        action={
          <Button onClick={onCreateTicket}>
            <IconPlus size={16} stroke={1.5} />
            新建工单
          </Button>
        }
      />

      <SectionCard title="工单列表" description={`共 ${total} 条`} bodyClassName="p-0">
        <div className="flex flex-wrap items-end gap-3 border-b border-border px-5 py-4">
          <div className="flex flex-col gap-1.5 sm:w-72">
            <Label htmlFor="ticket-keyword">关键词</Label>
            <div className="relative">
              <IconSearch
                size={16}
                stroke={1.5}
                className="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-muted-foreground"
              />
              <Input
                id="ticket-keyword"
                value={keyword}
                onChange={(event) => {
                  setKeyword(event.target.value)
                  reset()
                }}
                placeholder="主题、目标库或群集"
                className="pl-9"
              />
            </div>
          </div>
          <div className="flex flex-col gap-1.5 sm:w-52">
            <Label htmlFor="ticket-status">状态</Label>
            <Select
              id="ticket-status"
              value={status}
              onChange={(event) => {
                setStatus(event.target.value as TicketStatus | "all")
                reset()
              }}
            >
              <option value="all">全部状态</option>
              {Object.entries(ticketStatusMeta).map(([key, meta]) => (
                <option key={key} value={key}>
                  {meta.label}
                </option>
              ))}
            </Select>
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

/** 状态枚举值：与后端 mm_tickets.status 一致。 */
function statusCode(status: TicketStatus): number {
  const codes: Record<TicketStatus, number> = {
    "waiting-vld": 1,
    "vld-failure": 2,
    "vld-warning": 3,
    "waiting-mrv": 4,
    "mrv-failure": 5,
    ltgm: 6,
    done: 7,
    "exec-failure": 8,
    closed: 9,
  }
  return codes[status]
}
