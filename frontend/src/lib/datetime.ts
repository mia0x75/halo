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
