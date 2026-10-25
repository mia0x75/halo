import { useCallback, useMemo, useState } from "react"

import { PAGE_SIZE, paginate, type CursorPage } from "@/lib/cursor"

/**
 * 游标分页状态：只保留当前位置的游标与回退栈，翻页不重算偏移。
 * 过滤条件变化时用 reset 回到第一页。
 */
export function useCursorPagination<T extends { id: number }>(
  items: T[],
  pageSize = PAGE_SIZE,
): {
  page: CursorPage<T>
  /** 1-based page number for the current position */
  pageNumber: number
  reset: () => void
  next: () => void
  previous: () => void
} {
  const [cursor, setCursor] = useState("")
  const [history, setHistory] = useState<string[]>([])

  const page = useMemo(() => paginate(items, cursor, pageSize), [items, cursor, pageSize])

  const reset = useCallback(() => {
    setCursor("")
    setHistory([])
  }, [])

  const next = useCallback(() => {
    setHistory((stack) => [...stack, cursor])
    setCursor(page.pageInfo.endCursor)
  }, [cursor, page.pageInfo.endCursor])

  const previous = useCallback(() => {
    const previousCursor = history[history.length - 1] ?? ""
    setHistory((stack) => stack.slice(0, -1))
    setCursor(previousCursor)
  }, [history])

  return { page, pageNumber: history.length + 1, reset, next, previous }
}
