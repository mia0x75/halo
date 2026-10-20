import { useMemo, useState } from "react"
import { IconLayoutGrid, IconList, IconPlus } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { CursorPagination } from "@/components/ui/cursor-pagination"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { StatusDot } from "@/components/ui/status-dot"
import { useGraphQL } from "@/hooks/use-graphql"
import { CLUSTERS_QUERY, type ClusterNode } from "@/lib/api/queries"
import { clusters as mockClusters, type Cluster } from "@/lib/mock/data"

const pageSize = 6

type ViewMode = "card" | "table"

export function ClustersPage() {
  const [view, setView] = useState<ViewMode>("card")
  const [after, setAfter] = useState<string | null>(null)
  const [stack, setStack] = useState<(string | null)[]>([])

  const remote = useGraphQL<
    { clusters: { totalCount: number; pageInfo: { hasNextPage: boolean; hasPreviousPage: boolean; startCursor: string; endCursor: string }; edges: { cursor: string; node: ClusterNode }[] } | null },
    { first: number; after?: string }
  >(CLUSTERS_QUERY, { first: pageSize, after: after ?? undefined }, { clusters: null })

  const connection = remote.data?.clusters ?? null

  const rows = useMemo<Cluster[]>(() => {
    if (connection) {
      return connection.edges.map((edge) => ({
        id: 0,
        uuid: edge.node.UUID,
        alias: edge.node.Alias,
        host: edge.node.Host,
        ip: edge.node.IP,
        port: edge.node.Port,
        type: "mysql",
        group: edge.node.Group ?? "未分组",
        status: edge.node.Status === 1 ? "normal" : "blocked",
        latencyMs: 0,
        version: "",
      }))
    }
    return mockClusters.slice(0, pageSize)
  }, [connection])

  const total = connection?.totalCount ?? mockClusters.length
  const hasNext = connection?.pageInfo.hasNextPage ?? false
  const hasPrevious = connection?.pageInfo.hasPreviousPage ?? false
  const endCursor = connection?.pageInfo.endCursor ?? ""

  return (
    <>
      <PageHeader
        title="群集"
        description="目标数据库服务器，按分组隔离"
        action={
          <div className="flex items-center gap-2">
            <div className="flex rounded-md border border-border">
              <Button
                variant={view === "card" ? "subtle" : "ghost"}
                size="icon-sm"
                aria-label="卡片视图"
                aria-pressed={view === "card"}
                onClick={() => setView("card")}
              >
                <IconLayoutGrid size={16} stroke={1.5} />
              </Button>
              <Button
                variant={view === "table" ? "subtle" : "ghost"}
                size="icon-sm"
                aria-label="列表视图"
                aria-pressed={view === "table"}
                onClick={() => setView("table")}
              >
                <IconList size={16} stroke={1.5} />
              </Button>
            </div>
            <Button>
              <IconPlus size={16} stroke={1.5} />
              接入群集
            </Button>
          </div>
        }
      />

      {view === "card" ? (
        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
          {rows.map((cluster) => (
            <div
              key={cluster.uuid}
              className="flex flex-col rounded-lg border border-border bg-card p-5 shadow-xs"
            >
              <div className="flex items-start justify-between gap-2">
                <div className="flex min-w-0 flex-col leading-tight">
                  <span className="truncate text-sm font-semibold">{cluster.alias}</span>
                  <span className="truncate text-xs text-muted-foreground">{cluster.host}</span>
                </div>
                <Badge variant={cluster.type === "mysql" ? "info" : "default"}>{cluster.type}</Badge>
              </div>

              <dl className="mt-4 grid grid-cols-2 gap-y-2 text-sm">
                <dt className="text-muted-foreground">地址</dt>
                <dd className="text-right font-mono text-xs">
                  {cluster.ip}:{cluster.port}
                </dd>
                <dt className="text-muted-foreground">分组</dt>
                <dd className="text-right">{cluster.group}</dd>
                <dt className="text-muted-foreground">版本</dt>
                <dd className="text-right">{cluster.version || "-"}</dd>
              </dl>

              <div className="mt-4 flex items-center justify-between border-t border-border pt-3">
                <StatusDot
                  tone={cluster.status === "normal" ? "success" : "danger"}
                  label={cluster.status === "normal" ? "正常" : "已停用"}
                />
                <span className="text-xs text-muted-foreground">分组：{cluster.group}</span>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <SectionCard title="群集列表" bodyClassName="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>别名</TableHead>
                <TableHead>主机</TableHead>
                <TableHead>地址</TableHead>
                <TableHead>类型</TableHead>
                <TableHead>分组</TableHead>
                <TableHead>状态</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {rows.map((cluster) => (
                <TableRow key={cluster.uuid}>
                  <TableCell className="font-medium">{cluster.alias}</TableCell>
                  <TableCell className="text-muted-foreground">{cluster.host}</TableCell>
                  <TableCell className="font-mono text-sm text-muted-foreground">
                    {cluster.ip}:{cluster.port}
                  </TableCell>
                  <TableCell>
                    <Badge variant={cluster.type === "mysql" ? "info" : "default"}>
                      {cluster.type}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-muted-foreground">{cluster.group}</TableCell>
                  <TableCell>
                    <StatusDot
                      tone={cluster.status === "normal" ? "success" : "danger"}
                      label={cluster.status === "normal" ? "正常" : "已停用"}
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </SectionCard>
      )}

      <div className="rounded-lg border border-border bg-card shadow-xs">
        <CursorPagination
          divided={false}
          total={total}
          hasPreviousPage={hasPrevious}
          hasNextPage={hasNext}
          startCursor={connection?.pageInfo.startCursor}
          endCursor={endCursor}
          onPrevious={() => {
            const previousCursor = stack[stack.length - 1] ?? null
            setStack((items) => items.slice(0, -1))
            setAfter(previousCursor)
          }}
          onNext={() => {
            setStack((items) => [...items, after])
            setAfter(endCursor || null)
          }}
        />
      </div>
    </>
  )
}
