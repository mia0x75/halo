import { useMemo, useState } from "react"
import { IconSearch } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { StatusBadge } from "@/components/shared/status-badge"
import { Button } from "@/components/ui/button"
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
import { mockTickets } from "@/lib/api/mock"
import { formatDateTime } from "@/lib/format"

export function TicketsPage() {
  const [keyword, setKeyword] = useState("")

  const tickets = useMemo(() => {
    const needle = keyword.trim().toLowerCase()
    if (!needle) return mockTickets
    return mockTickets.filter(
      (ticket) =>
        ticket.subject.toLowerCase().includes(needle) ||
        ticket.database.toLowerCase().includes(needle) ||
        ticket.cluster.alias.toLowerCase().includes(needle),
    )
  }, [keyword])

  return (
    <>
      <PageHeader
        title="工单"
        description="变更工单的提交与审核状态"
        action={<Button>新建工单</Button>}
      />

      <SectionCard title="工单列表" description={`共 ${tickets.length} 条`} bodyClassName="p-0">
        <div className="flex flex-col gap-3 border-b border-border px-5 py-3 sm:flex-row sm:items-end">
          <div className="flex flex-col gap-1.5 sm:w-72">
            <Label htmlFor="ticket-keyword">搜索</Label>
            <div className="relative">
              <IconSearch
                size={16}
                stroke={1.5}
                className="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-muted-foreground"
              />
              <Input
                id="ticket-keyword"
                value={keyword}
                onChange={(event) => setKeyword(event.target.value)}
                placeholder="主题、库名或群集别名"
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
              <TableHead>提交时间</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {tickets.map((ticket) => (
              <TableRow key={ticket.uuid}>
                <TableCell className="max-w-[320px] truncate font-medium">
                  {ticket.subject}
                </TableCell>
                <TableCell className="text-muted-foreground">{ticket.database}</TableCell>
                <TableCell className="text-muted-foreground">{ticket.cluster.alias}</TableCell>
                <TableCell>
                  <StatusBadge status={ticket.status} />
                </TableCell>
                <TableCell className="text-muted-foreground">
                  {formatDateTime(ticket.createAt)}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </SectionCard>
    </>
  )
}
