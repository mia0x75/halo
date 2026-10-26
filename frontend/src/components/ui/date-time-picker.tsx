import { useEffect, useRef, useState } from "react"
import { IconCalendar, IconChevronLeft, IconChevronRight } from "@tabler/icons-react"

import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

/**
 * 预约用日期时间选择器，交互对齐 tw-elements 的 datetimepicker：
 * 触发输入框（右侧日历图标）→ 弹出日历 → 选完日期日历关闭、弹出时/分
 * 面板 → 点「确认」才把值写回。草稿值在面板内部维护。
 * 值与 min 都是墙上时间 `YYYY-MM-DDTHH:mm`，与预约接口入参格式一致。
 * 不传 min 时兜底为「当前时间 + 1 分钟」：预约永远只允许未来。
 */
export interface DateTimePickerProps {
  value: string
  onChange: (value: string) => void
  /** 最早可选时间；早于它的日期与时分被禁用。预约场景传 scheduleMinValue()。 */
  min: string
  className?: string
  id?: string
}

type Panel = "date" | "time"

const WEEKDAYS = ["一", "二", "三", "四", "五", "六", "日"]
const HOURS = Array.from({ length: 24 }, (_, i) => i)
const MINUTES = Array.from({ length: 60 }, (_, i) => i)
const ITEM_HEIGHT = 36

function pad(value: number): string {
  return value.toString().padStart(2, "0")
}

/** `YYYY-MM-DDTHH:mm` → 本地 Date；非法输入返回 null。 */
function parseLocal(value: string): Date | null {
  const [datePart, timePart] = value.split("T")
  if (!datePart || !timePart) return null
  const [y, m, d] = datePart.split("-").map(Number)
  const [hh, mm] = timePart.split(":").map(Number)
  if ([y, m, d, hh, mm].some((n) => !Number.isFinite(n))) return null
  const date = new Date(y, m - 1, d, hh, mm)
  return Number.isFinite(date.getTime()) ? date : null
}

function dayKey(date: Date): string {
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
}

function startOfMonth(date: Date): Date {
  return new Date(date.getFullYear(), date.getMonth(), 1)
}

/** 周一为一周第一天时的前置空位数。 */
function leadingBlanks(first: Date): number {
  return (first.getDay() + 6) % 7
}

/** 输入框展示值：2026-10-28 周三 22:15。 */
function displayValue(value: string): string {
  const date = parseLocal(value)
  if (!date) return ""
  return `${dayKey(date)} 周${WEEKDAYS[(date.getDay() + 6) % 7]} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}

const ITEM =
  "flex h-8 w-8 items-center justify-center rounded-full text-sm transition-colors"

export function DateTimePicker({ value, onChange, min, className, id }: DateTimePickerProps) {
  // min 由调用方保证有效；解析失败时退化为 epoch，等同「全部日期可选」。
  const minDate = parseLocal(min) ?? new Date(0)
  const minDayKey = dayKey(minDate)
  // 今天的高亮在打开期间不变，取一次即可，避免每次渲染读时钟。
  const [todayKey] = useState(() => dayKey(new Date()))

  const [open, setOpen] = useState(false)
  const [panel, setPanel] = useState<Panel>("date")
  const [draft, setDraft] = useState(value)
  const [month, setMonth] = useState(() => startOfMonth(parseLocal(value) ?? minDate))
  const wrapRef = useRef<HTMLDivElement>(null)

  const draftDate = parseLocal(draft) ?? minDate

  useEffect(() => {
    if (!open) return
    const onPointerDown = (event: PointerEvent) => {
      if (!wrapRef.current?.contains(event.target as Node)) setOpen(false)
    }
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") setOpen(false)
    }
    document.addEventListener("pointerdown", onPointerDown)
    document.addEventListener("keydown", onKeyDown)
    return () => {
      document.removeEventListener("pointerdown", onPointerDown)
      document.removeEventListener("keydown", onKeyDown)
    }
  }, [open])

  const openPanel = () => {
    // 每次打开都从日历开始；无有效值时以最早可选时间起步。
    const base = parseLocal(value) ?? minDate
    setDraft(value || dayKey(base) + "T" + `${pad(base.getHours())}:${pad(base.getMinutes())}`)
    setMonth(startOfMonth(base))
    setPanel("date")
    setOpen(true)
  }

  const pickDay = (day: number) => {
    const stamp = `${month.getFullYear()}-${pad(month.getMonth() + 1)}-${pad(day)}`
    let hour = draftDate.getHours()
    let minute = draftDate.getMinutes()
    // 时间与禁用最多的日期同天时，把沿用下来的时间顶到最早可选值，
    // 否则确认按钮默认就是禁用状态。
    if (stamp === minDayKey) {
      if (hour < minDate.getHours() || (hour === minDate.getHours() && minute < minDate.getMinutes())) {
        hour = minDate.getHours()
        minute = minDate.getMinutes()
      }
    }
    setDraft(`${stamp}T${pad(hour)}:${pad(minute)}`)
    setPanel("time")
  }

  const pickTime = (hour: number, minute: number) => {
    setDraft(`${dayKey(draftDate)}T${pad(hour)}:${pad(minute)}`)
  }

  const confirm = () => {
    if (!parseLocal(draft)) return
    onChange(draft)
    setOpen(false)
  }

  const sameDayAsMin = dayKey(draftDate) === minDayKey
  const draftValid = parseLocal(draft) !== null

  return (
    <div ref={wrapRef} className={cn("relative", className)}>
      <div className="relative">
        <input
          id={id}
          readOnly
          value={displayValue(value)}
          placeholder="选择执行时间"
          onClick={() => (open ? setOpen(false) : openPanel())}
          className="h-[2.375rem] w-full cursor-pointer rounded-md border border-input bg-card px-4 py-[0.5625rem] pr-10 text-sm leading-5 text-foreground transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-3 focus-visible:ring-ring/50"
        />
        <button
          type="button"
          aria-label="选择执行时间"
          onClick={() => (open ? setOpen(false) : openPanel())}
          className="absolute right-1 top-1/2 -translate-y-1/2 rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"        >
          <IconCalendar size={16} stroke={1.5} />
        </button>
      </div>

      {open ? (
        <div
          role="dialog"
          aria-label="选择执行时间"
          className={cn(
            "absolute left-0 top-full z-20 mt-1 w-[min(19rem,100vw-2rem)]",
            "rounded-lg border border-border bg-card p-3 shadow-md",
            "animate-in fade-in-0 zoom-in-95",
          )}
        >
          {panel === "date" ? (
            <>
              <div className="mb-1 flex items-center justify-between">
                <button
                  type="button"
                  aria-label="上个月"
                  disabled={startOfMonth(minDate) >= month}
                  onClick={() =>
                    setMonth((m) => new Date(m.getFullYear(), m.getMonth() - 1, 1))
                  }
                  className="rounded-md p-1 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground disabled:pointer-events-none disabled:opacity-40"
                >
                  <IconChevronLeft size={16} stroke={1.5} />
                </button>
                <span className="text-sm font-medium text-strong">
                  {month.getFullYear()} 年 {month.getMonth() + 1} 月
                </span>
                <button
                  type="button"
                  aria-label="下个月"
                  onClick={() =>
                    setMonth((m) => new Date(m.getFullYear(), m.getMonth() + 1, 1))
                  }
                  className="rounded-md p-1 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                >
                  <IconChevronRight size={16} stroke={1.5} />
                </button>
              </div>
              <div className="grid grid-cols-7 text-center text-xs text-muted-foreground">
                {WEEKDAYS.map((day) => (
                  <span key={day} className="py-1">
                    {day}
                  </span>
                ))}
              </div>
              <div className="grid grid-cols-7 justify-items-center gap-y-1">
                {Array.from({ length: leadingBlanks(month) }, (_, i) => (
                  <span key={`blank-${i}`} />
                ))}
                {Array.from(
                  { length: new Date(month.getFullYear(), month.getMonth() + 1, 0).getDate() },
                  (_, i) => i + 1,
                ).map((day) => {
                  const stamp = `${month.getFullYear()}-${pad(month.getMonth() + 1)}-${pad(day)}`
                  const isSelected = dayKey(draftDate) === stamp
                  const isToday = todayKey === stamp
                  const isDisabled = stamp < minDayKey
                  return (
                    <button
                      key={stamp}
                      type="button"
                      disabled={isDisabled}
                      onClick={() => pickDay(day)}
                      className={cn(
                        ITEM,
                        // 选中日期用正圆（Tabler 风格是圆角矩形，按要求改成正圆）。
                        isSelected
                          ? "bg-primary font-medium text-primary-foreground"
                          : "text-foreground hover:bg-accent",
                        isToday && !isSelected && "font-semibold text-primary",
                        isDisabled && "pointer-events-none opacity-35",
                      )}
                    >
                      {day}
                    </button>
                  )
                })}
              </div>
            </>
          ) : (
            <>
              <div className="mb-2 text-center text-sm text-strong">
                {displayValue(draft) || "选择时间"}
              </div>
              <div className="flex items-stretch gap-2">
                <TimeColumn
                  label="时"
                  items={HOURS}
                  selected={draftDate.getHours()}
                  disabled={(hour) => hour < minDate.getHours() && sameDayAsMin}
                  onSelect={(hour) => pickTime(hour, draftDate.getMinutes())}
                />
                <div className="flex items-center text-lg font-medium text-strong">:</div>
                <TimeColumn
                  label="分"
                  items={MINUTES}
                  selected={draftDate.getMinutes()}
                  disabled={(minute) =>
                    sameDayAsMin &&
                    draftDate.getHours() === minDate.getHours() &&
                    minute < minDate.getMinutes()
                  }
                  onSelect={(minute) => pickTime(draftDate.getHours(), minute)}
                />
              </div>
              <div className="mt-3 flex justify-end gap-2">
                <Button variant="outline" size="sm" onClick={() => setPanel("date")}>
                  重选日期
                </Button>
                <Button size="sm" disabled={!draftValid} onClick={confirm}>
                  确认
                </Button>
              </div>
            </>
          )}
        </div>
      ) : null}
    </div>
  )
}

interface TimeColumnProps {
  label: string
  items: number[]
  selected: number
  disabled: (value: number) => boolean
  onSelect: (value: number) => void
}

function TimeColumn({ label, items, selected, disabled, onSelect }: TimeColumnProps) {
  const scrollRef = useRef<HTMLDivElement>(null)

  // 滚轮停在选中值上；用容器 scrollTop 计算，避免 scrollIntoView
  // 连带滚动外层对话框。
  useEffect(() => {
    const container = scrollRef.current
    if (!container) return
    const index = items.indexOf(selected)
    if (index < 0) return
    const target = index * ITEM_HEIGHT - container.clientHeight / 2 + ITEM_HEIGHT / 2 + 4
    container.scrollTop = Math.max(0, target)
  }, [selected, items])

  return (
    <div className="flex min-w-0 flex-1 flex-col items-center gap-1">
      <span className="text-xs text-muted-foreground">{label}</span>
      <div
        ref={scrollRef}
        className="h-36 w-full overflow-y-auto rounded-md border border-border py-1 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
      >
        {items.map((item) => (
          <button
            key={item}
            type="button"
            disabled={disabled(item)}
            onClick={() => onSelect(item)}
            className={cn(
              ITEM,
              "mx-auto my-0.5 h-8 w-8 text-xs",
              item === selected
                ? "bg-primary font-medium text-primary-foreground"
                : "text-foreground hover:bg-accent",
              disabled(item) && "pointer-events-none opacity-35",
            )}
          >
            {pad(item)}
          </button>
        ))}
      </div>
    </div>
  )
}
