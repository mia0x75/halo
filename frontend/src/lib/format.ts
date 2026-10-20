const MINUTE = 60
const HOUR = 60 * MINUTE

function pad(value: number): string {
  return value.toString().padStart(2, "0")
}

/**
 * Accepts unix seconds (the server stores timestamps that way) or a Date.
 * Recent items read as relative time, older ones as an absolute stamp.
 */
export function formatRelativeTime(input: number | Date, now: number = Date.now()): string {
  const timestamp = typeof input === "number" ? input * 1000 : input.getTime()
  const seconds = Math.floor((now - timestamp) / 1000)

  if (seconds < 0) return formatDateTime(new Date(timestamp))
  if (seconds < MINUTE) return "just now"
  if (seconds < HOUR) {
    const minutes = Math.floor(seconds / MINUTE)
    return `${minutes} mins ago`
  }
  if (seconds < 24 * HOUR) {
    const hours = Math.floor(seconds / HOUR)
    return `${hours} hours ago`
  }
  return formatDateTime(new Date(timestamp))
}

export function formatDateTime(input: number | Date): string {
  const date = typeof input === "number" ? new Date(input * 1000) : input
  const y = date.getFullYear()
  const m = pad(date.getMonth() + 1)
  const d = pad(date.getDate())
  const hh = pad(date.getHours())
  const mm = pad(date.getMinutes())
  return `${y}-${m}-${d} ${hh}:${mm}`
}

export function formatNumber(value: number): string {
  return new Intl.NumberFormat("zh-CN").format(value)
}

export function formatPercent(value: number, fractionDigits = 1): string {
  return `${value.toFixed(fractionDigits)}%`
}
