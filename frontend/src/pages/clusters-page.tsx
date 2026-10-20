import { IconDatabase, IconPlus, IconServer } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { StatusDot } from "@/components/ui/status-dot"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { clusters, groups } from "@/lib/mock/data"

export function ClustersPage() {
  return (
    <>
      <PageHeader
        title="群集"
        description="目标数据库服务器，按分组隔离"
        action={
          <Button>
            <IconPlus size={16} stroke={1.5} />
            接入群集
          </Button>
        }
      />

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        {clusters.map((cluster) => (
          <div
            key={cluster.uuid}
            className="flex flex-col rounded-lg border border-border bg-card p-5 shadow-xs"
          >
            <div className="flex items-start justify-between gap-2">
              <div className="flex items-center gap-2">
                <IconServer size={20} stroke={1.5} className="text-muted-foreground" />
                <div className="flex flex-col leading-tight">
                  <span className="text-sm font-semibold">{cluster.alias}</span>
                  <span className="text-xs text-muted-foreground">{cluster.host}</span>
                </div>
              </div>
              <Badge variant={cluster.type === "mysql" ? "info" : "default"}>
                {cluster.type}
              </Badge>
            </div>

            <dl className="mt-4 grid grid-cols-2 gap-y-2 text-sm">
              <dt className="text-muted-foreground">地址</dt>
              <dd className="text-right font-mono text-xs">
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
              <span className="inline-flex items-center gap-1 text-xs text-muted-foreground">
                <IconDatabase size={14} stroke={1.5} />
                {groups.find((group) => group.name === cluster.group)?.description}
              </span>
            </div>
          </div>
        ))}
      </div>
    </>
  )
}
