/**
 * 本地游标分页工具。接口分页一律走服务端游标；这里只给「一次拉全量再本地翻页」
 * 的页面（规则库、命中记录）用，编码格式与后端 EncodeCursor 保持一致。
 */

// Keep paginated lists consistent without changing lookup or preview limits.
export const PAGE_SIZE = 10

export interface CursorPage<T> {
  edges: T[]
  pageInfo: {
    hasPreviousPage: boolean
    hasNextPage: boolean
    startCursor: string
    endCursor: string
  }
  totalCount: number
}

export function encodeCursor(id: number | string): string {
  return btoa(`cursor${id}`)
}

export function decodeCursor(cursor: string): string | null {
  try {
    const raw = atob(cursor)
    return raw.startsWith("cursor") ? raw.slice("cursor".length) : null
  } catch {
    return null
  }
}

/**
 * 按游标切分一页。cursor 为空表示第一页；after 表示「取该条之后的内容」。
 * items 需要按 id 升序排列（游标语义依赖稳定的顺序）。
 */
export function paginate<T extends { id: number }>(
  items: T[],
  cursor: string,
  first: number,
): CursorPage<T> {
  const totalCount = items.length
  const all = [...items].sort((a, b) => a.id - b.id)

  let start = 0
  if (cursor) {
    const raw = decodeCursor(cursor)
    const id = raw === null ? Number.NaN : Number(raw)
    const index = all.findIndex((item) => item.id === id)
    start = index >= 0 ? index + 1 : 0
  }

  const edges = all.slice(start, start + first)
  const hasNextPage = start + first < all.length
  const hasPreviousPage = start > 0

  return {
    edges,
    pageInfo: {
      hasPreviousPage,
      hasNextPage,
      startCursor: edges.length > 0 ? encodeCursor(edges[0]!.id) : "",
      endCursor: edges.length > 0 ? encodeCursor(edges[edges.length - 1]!.id) : "",
    },
    totalCount,
  }
}
