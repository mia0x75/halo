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
import { useCursorPagination } from "@/hooks/use-cursor-pagination"
import { ticketStatusMeta, tickets, type TicketStatus } from "@/lib/mock/data"

const pageSize = 6

export function TicketsPage() {
  const [keyword, setKeyword] = useState("")
  const [status, setStatus] = useState<TicketStatus | "all">("all")
  const filtered = useMemo(() => {
    const needle = keyword.trim().toLowerCase()
    return tickets.filter((ticket) => {
      if (status !== "all" && ticket.status !== status) return false
      if (!needle) return true
      return (
        ticket.subject.toLowerCase().includes(needle) ||
        ticket.database.toLowerCase().includes(needle) ||
        ticket.cluster.toLowerCase().includes(needle)
      )
    })
  }, [keyword, status])

  const { page, reset, next, previous } = useCursorPagination(filtered, pageSize)

  return (
    <>
      <PageHeader
        title="工单"
        description="变更工单的提交与审核状态"
        action={
          <Button>
            <IconPlus size={16} stroke={1.5} />
            新建工单
          </Button>
        }
      />

      <SectionCard
        title="工单列表"
        description={`共 ${filtered.length} 条`}
        bodyClassName="p-0"
      >
        <div className="flex flex-col gap-3 border-b border-border px-5 py-4 sm:flex-row sm:items-end">
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
            {page.edges.map((ticket) => (
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
