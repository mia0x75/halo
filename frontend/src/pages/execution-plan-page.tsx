import { useEffect, useMemo, useState, type ReactNode } from "react"
import { IconCalendar, IconListDetails } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { BlockNote } from "@/components/shared/status-note"
import { StatusBadge } from "@/components/shared/status-badge"
import { TicketCalendar } from "@/components/shared/ticket-calendar"
import { Badge } from "@/components/ui/badge"
import { CursorPagination } from "@/components/ui/cursor-pagination"
import { Tabs } from "@/components/ui/tabs"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useGraphQL } from "@/hooks/use-graphql"
import { subscribeDataChanged } from "@/lib/api/data-events"
import { formatUnix } from "@/lib/api/mappers"
import {
  EXECUTION_TICKETS_QUERY,
  type ExecutionTicketConnectionNode,
  type ExecutionTicketNode,
} from "@/lib/api/queries"
import { PAGE_SIZE } from "@/lib/cursor"

type Tab = "pending" | "scheduled"
type View = "table" | "calendar"

// 日历模式一次拉全当月之外的全部预约（预约量有限），不走游标翻页；
// 后端 resolver 对 first 无上限，500 足够覆盖并留有余量。
const CALENDAR_FETCH_SIZE = 500

export interface ExecutionPlanPageProps {
  onOpenTicket: (uuid: string) => void
}

/** 执行计划页：待上线执行与已预约的工单集中在一处查看。 */
export function ExecutionPlanPage({ onOpenTicket }: ExecutionPlanPageProps) {
  const [tab, setTab] = useState<Tab>("pending")
  const [view, setView] = useState<View>("table")
  const [after, setAfter] = useState<string | null>(null)
  const [stack, setStack] = useState<(string | null)[]>([])

  const calendarMode = tab === "scheduled" && view === "calendar"

  // 两个 Tab 走同一个查询，scheduled 参数决定过滤哪一半；
  // 工单在详情页被审核/执行后徽章数量要跟着刷新。
  const pending = useGraphQL<
    { executionTickets: ExecutionTicketConnectionNode | null },
    { first: number; after?: string; scheduled?: boolean }
  >(EXECUTION_TICKETS_QUERY, { first: PAGE_SIZE, after: after ?? undefined }, { executionTickets: null })

  const scheduled = useGraphQL<
    { executionTickets: ExecutionTicketConnectionNode | null },
    { first: number; after?: string; scheduled?: boolean }
  >(
    EXECUTION_TICKETS_QUERY,
    // 日历模式一次取全（不分页），表格模式维持每页 10 条。
    calendarMode
      ? { first: CALENDAR_FETCH_SIZE, scheduled: true }
      : { first: PAGE_SIZE, after: after ?? undefined, scheduled: true },
    { executionTickets: null },
  )

  useEffect(() => {
    const unsubscribe = subscribeDataChanged("tickets", () => {
      pending.refetch()
      scheduled.refetch()
    })
    return unsubscribe
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pending.refetch, scheduled.refetch])

  const active = tab === "pending" ? pending : scheduled
  const connection = active.data?.executionTickets ?? null
  const rows = useMemo<ExecutionTicketNode[]>(
    () => connection?.edges.map((edge) => edge.node) ?? [],
    [connection],
  )

  return (
    <>
      <PageHeader
        title="执行计划"
        description="审核通过等待上线执行的工单与已经预约的工单"
      />

      <Tabs
        className="px-1"
        value={tab}
        onChange={(key) => {
          setTab(key as Tab)
          setAfter(null)
          setStack([])
        }}
        items={[
          { key: "pending", label: "待上线", badge: pending.data?.executionTickets?.totalCount ?? 0 },
          { key: "scheduled", label: "已预约", badge: scheduled.data?.executionTickets?.totalCount ?? 0 },
        ]}
      />

      {active.error ? <BlockNote className="mx-4">{active.error}</BlockNote> : null}

      <SectionCard
        bodyClassName="p-0"
        action={
          tab === "scheduled" ? (
            <ViewToggle
              value={view}
              onChange={(next) => {
                setView(next)
                setAfter(null)
                setStack([])
              }}
            />
          ) : undefined
        }
      >
        {calendarMode ? (
          <TicketCalendar tickets={rows} onOpenTicket={onOpenTicket} />
        ) : (
          <>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>主题</TableHead>
                  <TableHead>目标库</TableHead>
                  <TableHead>群集</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>审核人</TableHead>
                  {tab === "scheduled" ? <TableHead>计划执行时间</TableHead> : null}
                  <TableHead>创建时间</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {rows.map((ticket) => (
                  <TableRow
                    key={ticket.UUID}
                    className="cursor-pointer"
                    onClick={() => onOpenTicket(ticket.UUID)}
                  >
                    <TableCell className="font-medium">{ticket.Subject}</TableCell>
                    <TableCell className="text-muted-foreground">{ticket.Database}</TableCell>
                    <TableCell className="text-muted-foreground">{ticket.Cluster?.Alias ?? "-"}</TableCell>
                    <TableCell>
                      <StatusBadge status="ltgm" />
                    </TableCell>
                    <TableCell className="text-muted-foreground">{ticket.Reviewer?.Name ?? "-"}</TableCell>
                    {tab === "scheduled" ? (
                      <TableCell className="text-muted-foreground">
                        {ticket.Cron ? (
                          <Badge variant="info">{formatCronTime(ticket.Cron.NextRun)}</Badge>
                        ) : (
                          "-"
                        )}
                      </TableCell>
                    ) : null}
                    <TableCell className="text-muted-foreground">{formatUnix(ticket.CreateAt)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            <CursorPagination
              page={stack.length + 1}
              empty={rows.length === 0}
              hasPreviousPage={connection?.pageInfo.hasPreviousPage ?? false}
              hasNextPage={connection?.pageInfo.hasNextPage ?? false}
              onPrevious={() => {
                const previousCursor = stack[stack.length - 1] ?? null
                setStack((items) => items.slice(0, -1))
                setAfter(previousCursor)
              }}
              onNext={() => {
                setStack((items) => [...items, after])
                setAfter(connection?.pageInfo.endCursor || null)
              }}
              onReset={() => {
                setAfter(null)
                setStack([])
              }}
            />
          </>
        )}
      </SectionCard>
    </>
  )
}

interface ViewToggleProps {
  value: View
  onChange: (view: View) => void
}

/** 列表 / 日历视图切换：分段按钮，样式对齐现有 outline Button。 */
function ViewToggle({ value, onChange }: ViewToggleProps) {
  const items: { key: View; label: string; icon: ReactNode }[] = [
    { key: "table", label: "列表", icon: <IconListDetails size={14} stroke={1.5} /> },
    { key: "calendar", label: "日历", icon: <IconCalendar size={14} stroke={1.5} /> },
  ]
  return (
    <div className="flex items-center rounded-md border border-border p-0.5" role="group" aria-label="视图切换">
      {items.map((item) => (
        <button
          key={item.key}
          type="button"
          aria-pressed={value === item.key}
          onClick={() => onChange(item.key)}
          className={
            value === item.key
              ? "flex items-center gap-1.5 rounded-[5px] bg-accent px-2.5 py-1 text-xs font-medium text-strong"
              : "flex items-center gap-1.5 rounded-[5px] px-2.5 py-1 text-xs text-muted-foreground transition-colors hover:text-foreground"
          }
        >
          {item.icon}
          {item.label}
        </button>
      ))}
    </div>
  )
}

/** 预约时间展示：接口给的是 RFC3339 串，按本地时区格式化。 */
function formatCronTime(nextRun: string | null): string {
  if (!nextRun) return "-"
  const ms = Date.parse(nextRun)
  if (!Number.isFinite(ms)) return nextRun
  return formatUnix(ms / 1000)
}
