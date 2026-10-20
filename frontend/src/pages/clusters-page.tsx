import { useMemo, useState } from "react"
import { IconLayoutGrid, IconList, IconPlus, IconServer } from "@tabler/icons-react"

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
import { useCursorPagination } from "@/hooks/use-cursor-pagination"
import { clusters } from "@/lib/mock/data"

const pageSize = 6

type ViewMode = "card" | "table"

export function ClustersPage() {
  const [view, setView] = useState<ViewMode>("card")
  const { page, next, previous } = useCursorPagination(clusters, pageSize)

  const groups = useMemo(() => new Set(clusters.map((cluster) => cluster.group)), [])

  return (
    <>
      <PageHeader
        title="群集"
        description={`目标数据库服务器，按分组隔离，共 ${groups.size} 个分组`}
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
          {page.edges.map((cluster) => (
            <div
              key={cluster.uuid}
              className="flex flex-col rounded-lg border border-border bg-card p-5 shadow-xs"
            >
              <div className="flex items-start justify-between gap-2">
                <div className="flex items-center gap-2">
                  <IconServer size={20} stroke={1.5} className="text-muted-foreground" />
                  <div className="flex flex-col leading-tight">
                    <span className="text-sm font-semibold">{cluster.alias}</span>
                    <span className="text-sm text-muted-foreground">{cluster.host}</span>
                  </div>
                </div>
                <Badge variant={cluster.type === "mysql" ? "info" : "default"}>
                  {cluster.type}
                </Badge>
              </div>

              <dl className="mt-4 grid grid-cols-2 gap-y-2 text-sm">
                <dt className="text-muted-foreground">地址</dt>
                <dd className="text-right font-mono text-sm">
                  {cluster.ip}:{cluster.port}
                </dd>
                <dt className="text-muted-foreground">分组</dt>
                <dd className="text-right">{cluster.group}</dd>
                <dt className="text-muted-foreground">版本</dt>
                <dd className="text-right">{cluster.version}</dd>
                <dt className="text-muted-foreground">延迟</dt>
                <dd className="text-right">{cluster.latencyMs} ms</dd>
              </dl>

              <div className="mt-4 flex items-center justify-between border-t border-border pt-3">
                <StatusDot
                  tone={cluster.status === "normal" ? "success" : "danger"}
                  label={cluster.status === "normal" ? "正常" : "已停用"}
                />
                <span className="text-sm text-muted-foreground">分组：{cluster.group}</span>
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
                <TableHead>版本</TableHead>
                <TableHead>状态</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {page.edges.map((cluster) => (
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
                  <TableCell className="text-muted-foreground">{cluster.version}</TableCell>
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
          total={page.totalCount}
          hasPreviousPage={page.pageInfo.hasPreviousPage}
          hasNextPage={page.pageInfo.hasNextPage}
          startCursor={page.pageInfo.startCursor}
          endCursor={page.pageInfo.endCursor}
          onPrevious={previous}
          onNext={next}
        />
      </div>
    </>
  )
}
