import { useMemo, useState } from "react"
import { IconChevronLeft, IconChevronRight } from "@tabler/icons-react"

import type { ExecutionTicketNode } from "@/lib/api/queries"
import { leadingBlanks, localDayKey } from "@/lib/datetime"
import { cn } from "@/lib/utils"

const WEEKDAYS = ["一", "二", "三", "四", "五", "六", "日"]
// 格子里默认最多展示 3 条工单片，多余折叠成「+N」。
const MAX_CHIPS = 3

export interface TicketCalendarProps {
  tickets: ExecutionTicketNode[]
  onOpenTicket: (uuid: string) => void
}

interface ScheduledEntry {
  ticket: ExecutionTicketNode
  day: string
  time: string
  ms: number
}

/** 已预约工单的月历视图：按计划执行时间把工单归入日历格。 */
export function TicketCalendar({ tickets, onOpenTicket }: TicketCalendarProps) {
  const [today] = useState(() => new Date())
  const [month, setMonth] = useState(() => new Date(today.getFullYear(), today.getMonth(), 1))
  const [expandedDay, setExpandedDay] = useState<string | null>(null)

  const todayKey = localDayKey(today)

  const byDay = useMemo(() => {
    const map = new Map<string, ScheduledEntry[]>()
    for (const ticket of tickets) {
      const nextRun = ticket.Cron?.NextRun
      if (!nextRun) continue
      const ms = Date.parse(nextRun)
      if (!Number.isFinite(ms)) continue
      const date = new Date(ms)
      const entry: ScheduledEntry = {
        ticket,
        day: localDayKey(date),
        time: date.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit", hour12: false }),
        ms,
      }
      const bucket = map.get(entry.day)
      if (bucket) bucket.push(entry)
      else map.set(entry.day, [entry])
    }
    for (const bucket of map.values()) bucket.sort((a, b) => a.ms - b.ms)
    return map
  }, [tickets])

  const monthDays = Array.from(
    { length: new Date(month.getFullYear(), month.getMonth() + 1, 0).getDate() },
    (_, i) => i + 1,
  )
  const blanks = leadingBlanks(month)
  // 尾部补空格让总格数凑满整行，末行去边框的选择器才能准确命中。
  const trailingBlanks = (7 - ((blanks + monthDays.length) % 7)) % 7
  const monthLabel = `${month.getFullYear()} 年 ${month.getMonth() + 1} 月`
  const inCurrentMonth =
    month.getFullYear() === today.getFullYear() && month.getMonth() === today.getMonth()

  return (
    <div className="p-3">
      <div className="mb-3 flex items-center justify-between">
        <div className="flex items-center gap-1">
          <button
            type="button"
            aria-label="上个月"
            onClick={() => setMonth((m) => new Date(m.getFullYear(), m.getMonth() - 1, 1))}
            className="rounded-md p-1 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
          >
            <IconChevronLeft size={16} stroke={1.5} />
          </button>
          <span className="min-w-24 text-center text-sm font-medium text-strong">{monthLabel}</span>
          <button
            type="button"
            aria-label="下个月"
            onClick={() => setMonth((m) => new Date(m.getFullYear(), m.getMonth() + 1, 1))}
            className="rounded-md p-1 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
          >
            <IconChevronRight size={16} stroke={1.5} />
          </button>
        </div>
        {inCurrentMonth ? null : (
          <button
            type="button"
            onClick={() => setMonth(new Date(today.getFullYear(), today.getMonth(), 1))}
            className="rounded-md px-2 py-1 text-xs text-primary transition-colors hover:bg-accent"
          >
            回到本月
          </button>
        )}
      </div>

      <div className="overflow-hidden rounded-md border border-border">
        <div className="grid grid-cols-7 border-b border-border bg-surface-secondary text-center text-xs text-muted-foreground">
          {WEEKDAYS.map((day) => (
            <span key={day} className="py-1.5">
              {day}
            </span>
          ))}
        </div>
        <div className="grid grid-cols-7 [&>*:nth-child(7n)]:border-r-0 [&>*:nth-last-child(-n+7)]:border-b-0">
          {Array.from({ length: blanks }, (_, i) => (
            <span key={`blank-${i}`} className="min-h-24 border-b border-r border-border bg-surface-secondary/40" />
          ))}
          {monthDays.map((day) => {
            const stamp = `${month.getFullYear()}-${String(month.getMonth() + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`
            const entries = byDay.get(stamp) ?? []
            const isToday = stamp === todayKey
            const expanded = expandedDay === stamp
            const visible = expanded ? entries : entries.slice(0, MAX_CHIPS)
            return (
              <div
                key={stamp}
                className={cn(
                  "flex min-h-24 flex-col gap-1 border-b border-r border-border p-1",
                  isToday && "bg-primary/5",
                )}
              >
                <span
                  className={cn(
                    "flex h-5 w-5 items-center justify-center rounded-full text-xs",
                    isToday ? "bg-primary font-medium text-primary-foreground" : "text-muted-foreground",
                  )}
                >
                  {day}
                </span>
                {visible.map((entry) => (
                  <button
                    key={entry.ticket.UUID}
                    type="button"
                    title={`${entry.time} ${entry.ticket.Subject}`}
                    onClick={() => onOpenTicket(entry.ticket.UUID)}
                    className={cn(
                      "flex w-full items-baseline gap-1 rounded-md bg-accent px-1 py-0.5 text-left text-xs transition-colors hover:bg-primary/10",
                      isToday && "bg-primary/10",
                    )}
                  >
                    <span className="shrink-0 font-medium tabular-nums text-primary">{entry.time}</span>
                    <span className="truncate text-foreground">{entry.ticket.Subject}</span>
                  </button>
                ))}
                {entries.length > MAX_CHIPS ? (
                  <button
                    type="button"
                    onClick={() => setExpandedDay(expanded ? null : stamp)}
                    className="self-start rounded-md px-1 py-0.5 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                  >
                    {expanded ? "收起" : `+${entries.length - MAX_CHIPS} 更多`}
                  </button>
                ) : null}
              </div>
            )
          })}
          {Array.from({ length: trailingBlanks }, (_, i) => (
            <span key={`tail-${i}`} className="min-h-24 border-b border-r border-border bg-surface-secondary/40" />
          ))}
        </div>
      </div>
    </div>
  )
}
