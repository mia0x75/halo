/** Shared helpers for local wall-time strings (`YYYY-MM-DDTHH:mm`). */

function pad(value: number): string {
  return value.toString().padStart(2, "0")
}

/** 预约最早可选时间：当前时间 + 1 分钟，留一点提交余量。 */
export function scheduleMinValue(): string {
  const date = new Date(Date.now() + 60_000)
  return (
    `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}` +
    `T${pad(date.getHours())}:${pad(date.getMinutes())}`
  )
}

/** 本地时区的日键 `YYYY-MM-DD`，用于把时间点归入日历格。 */
export function localDayKey(date: Date): string {
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
}

/** 周一为一周第一天时的前置空位数。 */
export function leadingBlanks(first: Date): number {
  return (first.getDay() + 6) % 7
}
