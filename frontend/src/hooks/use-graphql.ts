import { useCallback, useEffect, useMemo, useState } from "react"

import { graphqlRequest } from "@/lib/api/client"
import { getToken } from "@/lib/api/session"

export interface QueryState<TData> {
  data: TData | null
  error: string | null
  loading: boolean
  refetch: () => void
}

/**
 * 轻量查询 hook：请求一次、失败时把 fallback 交给页面。
 * 没有查询缓存，也不做自动重试——演示阶段够用。
 */
export function useGraphQL<TData, TVariables extends Record<string, unknown> = Record<string, never>>(
  query: string,
  variables?: TVariables,
  fallback?: TData,
): QueryState<TData> {
  const [data, setData] = useState<TData | null>(fallback ?? null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [nonce, setNonce] = useState(0)
  // 变量每次渲染都是新对象，序列化成字符串后再放进依赖，
  // 只有内容变化才会重新发请求。
  const serialized = useMemo(() => JSON.stringify(variables ?? {}), [variables])

  useEffect(() => {
    let cancelled = false
    const token = getToken()
    if (!token) {
      // 没有令牌就不请求，页面继续用假数据。
      setData(fallback ?? null)
      return
    }

    setLoading(true)
    graphqlRequest<TData, TVariables>(query, variables, token)
      .then((result) => {
        if (cancelled) return
        setData(result)
        setError(null)
      })
      .catch((cause: unknown) => {
        if (cancelled) return
        setError(cause instanceof Error ? cause.message : String(cause))
        setData(fallback ?? null)
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })

    return () => {
      cancelled = true
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query, nonce, serialized])

  const refetch = useCallback(() => setNonce((value) => value + 1), [])

  return { data, error, loading, refetch }
}
