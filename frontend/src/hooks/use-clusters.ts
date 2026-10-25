import { useCallback, useMemo, useState } from "react"

import { useGraphQL } from "@/hooks/use-graphql"
import { CLUSTERS_QUERY, type ClusterNode } from "@/lib/api/queries"
import { PAGE_SIZE } from "@/lib/cursor"
import type { Cluster } from "@/lib/model"

/** 接口返回的群集节点转成页面模型。 */
export function toCluster(node: ClusterNode): Cluster {
  return {
    id: 0,
    uuid: node.UUID,
    alias: node.Alias,
    host: node.Host,
    ip: node.IP,
    port: node.Port,
    type: node.Type === "postgres" ? "postgres" : "mysql",
    group: node.Group ?? "未分组",
    status: node.Status === 1 ? "normal" : "blocked",
    latencyMs: 0,
    version: "",
  }
}

export interface ClusterPage {
  clusters: Cluster[]
  total: number
  hasNext: boolean
  hasPrevious: boolean
  startCursor: string
  endCursor: string
  /** 1-based page number for the current position */
  pageNumber: number
  next: () => void
  previous: () => void
  reset: () => void
  refetch: () => void
  loading: boolean
}

/**
 * 群集列表。群集管理、数据库页与新建工单的群集下拉共用同一份分页数据。
 */
export function useClusters(pageSize = PAGE_SIZE): ClusterPage {
  const [after, setAfter] = useState<string | null>(null)
  const [stack, setStack] = useState<(string | null)[]>([])

  const remote = useGraphQL<
    {
      clusters: {
        totalCount: number
        pageInfo: { hasNextPage: boolean; hasPreviousPage: boolean; startCursor: string; endCursor: string }
        edges: { cursor: string; node: ClusterNode }[]
      } | null
    },
    { first: number; after?: string }
  >(CLUSTERS_QUERY, { first: pageSize, after: after ?? undefined }, { clusters: null })

  const connection = remote.data?.clusters ?? null

  const clusters = useMemo<Cluster[]>(
    () => connection?.edges.map((edge) => toCluster(edge.node)) ?? [],
    [connection],
  )

  const next = useCallback(() => {
    setStack((items) => [...items, after])
    setAfter(connection?.pageInfo.endCursor || null)
  }, [after, connection])

  const previous = useCallback(() => {
    const previousCursor = stack[stack.length - 1] ?? null
    setStack((items) => items.slice(0, -1))
    setAfter(previousCursor)
  }, [stack])

  const reset = useCallback(() => {
    setAfter(null)
    setStack([])
  }, [])

  return {
    clusters,
    total: connection?.totalCount ?? 0,
    hasNext: connection?.pageInfo.hasNextPage ?? false,
    hasPrevious: connection?.pageInfo.hasPreviousPage ?? false,
    startCursor: connection?.pageInfo.startCursor ?? "",
    endCursor: connection?.pageInfo.endCursor ?? "",
    pageNumber: stack.length + 1,
    next,
    previous,
    reset,
    refetch: remote.refetch,
    loading: remote.loading,
  }
}
