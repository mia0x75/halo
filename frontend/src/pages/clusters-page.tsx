import { useMemo, useState } from "react"
import {
  IconDotsVertical,
  IconEdit,
  IconLayoutGrid,
  IconList,
  IconPlus,
  IconTrash,
} from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { CursorPagination } from "@/components/ui/cursor-pagination"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import { StatusDot } from "@/components/ui/status-dot"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  ClusterFormDialog,
  type ClusterFormValues,
} from "@/features/clusters/cluster-form-dialog"
import { useClusters } from "@/hooks/use-clusters"
import { graphqlRequest } from "@/lib/api/client"
import { REMOVE_CLUSTER_MUTATION } from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

type ViewMode = "card" | "table"

export function ClustersPage() {
  const [view, setView] = useState<ViewMode>("card")
  const [keyword, setKeyword] = useState("")
  const [type, setType] = useState("all")
  const [group, setGroup] = useState("all")
  const [form, setForm] = useState<
    { mode: "create" } | { mode: "edit"; cluster: ClusterFormValues & { uuid: string } } | null
  >(null)
  const [error, setError] = useState<string | null>(null)

  const page = useClusters(8)

  const groups = useMemo(
    () => Array.from(new Set(page.clusters.map((cluster) => cluster.group))),
    [page.clusters],
  )

  const rows = useMemo(() => {
    const needle = keyword.trim().toLowerCase()
    return page.clusters.filter((cluster) => {
      if (type !== "all" && cluster.type !== type) return false
      if (group !== "all" && cluster.group !== group) return false
      if (!needle) return true
      return (
        cluster.alias.toLowerCase().includes(needle) ||
        cluster.host.toLowerCase().includes(needle) ||
        cluster.ip.includes(needle)
      )
    })
  }, [page.clusters, keyword, type, group])

  const remove = async (uuid: string) => {
    setError(null)
    try {
      await graphqlRequest(REMOVE_CLUSTER_MUTATION, { id: uuid }, getToken())
      page.refetch()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    }
  }

  return (
    <>
      <PageHeader
        title="群集管理"
        description="目标数据库所在的主机，按分组隔离"
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
            <ClusterFormDialog
              open={form?.mode === "create"}
              onOpenChange={(open) => setForm(open ? { mode: "create" } : null)}
              onSaved={() => {
                page.refetch()
              }}
              trigger={
                <Button>
                  <IconPlus size={16} stroke={1.5} />
                  接入群集
                </Button>
              }
            />
          </div>
        }
      />

      {error ? (
        <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
      ) : null}

      <SectionCard
        title="群集列表"
        description={`共 ${page.total} 个群集`}
        bodyClassName="p-0"
      >
        <div className="flex flex-wrap items-end gap-3 border-b border-border px-5 py-4">
          <div className="flex flex-col gap-1.5 sm:w-56">
            <Label htmlFor="cluster-keyword">关键词</Label>
            <Input
              id="cluster-keyword"
              value={keyword}
              onChange={(event) => {
                setKeyword(event.target.value)
                page.reset()
              }}
              placeholder="别名、主机名或 IP"
            />
          </div>
          <div className="flex flex-col gap-1.5 sm:w-40">
            <Label htmlFor="cluster-type">类型</Label>
            <Select
              id="cluster-type"
              value={type}
              onChange={(event) => {
                setType(event.target.value)
                page.reset()
              }}
            >
              <option value="all">全部类型</option>
              <option value="mysql">MySQL</option>
              <option value="postgres">PostgreSQL</option>
            </Select>
          </div>
          <div className="flex flex-col gap-1.5 sm:w-44">
            <Label htmlFor="cluster-group">分组</Label>
            <Select
              id="cluster-group"
              value={group}
              onChange={(event) => {
                setGroup(event.target.value)
                page.reset()
              }}
            >
              <option value="all">全部分组</option>
              {groups.map((item) => (
                <option key={item} value={item}>
                  {item}
                </option>
              ))}
            </Select>
          </div>
        </div>

        {view === "card" ? (
          <div className="grid gap-4 p-5 sm:grid-cols-2 xl:grid-cols-3">
            {rows.map((cluster) => (
              <div
                key={cluster.uuid}
                className="flex flex-col rounded-lg border border-border bg-card p-5"
              >
                <div className="flex items-start justify-between gap-2">
                  <div className="flex min-w-0 flex-col leading-tight">
                    <span className="truncate text-sm font-semibold">{cluster.alias}</span>
                    <span className="truncate text-xs text-muted-foreground">{cluster.host}</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <Badge variant={cluster.type === "mysql" ? "info" : "default"}>
                      {cluster.type}
                    </Badge>
                    <ClusterMenu
                      cluster={cluster}
                      onEdit={() =>
                        setForm({
                          mode: "edit",
                          cluster: {
                          uuid: cluster.uuid,
                          alias: cluster.alias,
                          host: cluster.host,
                          ip: cluster.ip,
                          port: cluster.port,
                          user: "",
                          password: "",
                          status: cluster.status === "normal" ? 1 : 2,
                          type: cluster.type,
                            groupUUID: "",
                          },
                        })
                      }
                      onRemove={() => remove(cluster.uuid)}
                    />
                  </div>
                </div>

                <dl className="mt-4 grid grid-cols-2 gap-y-2 text-sm">
                  <dt className="text-muted-foreground">地址</dt>
                  <dd className="text-right font-mono text-xs">
                    {cluster.ip}:{cluster.port}
                  </dd>
                  <dt className="text-muted-foreground">分组</dt>
                  <dd className="text-right">{cluster.group}</dd>
                </dl>

                <div className="mt-4 flex items-center justify-between border-t border-border pt-3">
                  <StatusDot
                    tone={cluster.status === "normal" ? "success" : "danger"}
                    label={cluster.status === "normal" ? "正常" : "已停用"}
                  />
                  <span className="text-xs text-muted-foreground">
                    {cluster.ip.includes("127.0.0.1") ? "本机" : "远端"}
                  </span>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>别名</TableHead>
                <TableHead>主机</TableHead>
                <TableHead>地址</TableHead>
                <TableHead>类型</TableHead>
                <TableHead>分组</TableHead>
                <TableHead>状态</TableHead>
                <TableHead className="w-12"></TableHead>
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
                  <TableCell>
                    <ClusterMenu
                      cluster={cluster}
                      onEdit={() =>
                        setForm({
                          mode: "edit",
                          cluster: {
                          uuid: cluster.uuid,
                          alias: cluster.alias,
                          host: cluster.host,
                          ip: cluster.ip,
                          port: cluster.port,
                          user: "",
                          password: "",
                          status: cluster.status === "normal" ? 1 : 2,
                          type: cluster.type,
                            groupUUID: "",
                          },
                        })
                      }
                      onRemove={() => remove(cluster.uuid)}
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}

        <CursorPagination
          total={page.total}
          hasPreviousPage={page.hasPrevious}
          hasNextPage={page.hasNext}
          startCursor={page.startCursor}
          endCursor={page.endCursor}
          onPrevious={page.previous}
          onNext={page.next}
        />
      </SectionCard>

      {form?.mode === "edit" ? (
        <ClusterFormDialog
          open
          onOpenChange={(open) => {
            if (!open) setForm(null)
          }}
          cluster={form.cluster}
          onSaved={() => {
            setForm(null)
            page.refetch()
          }}
        />
      ) : null}
    </>
  )
}

interface ClusterMenuProps {
  cluster: { uuid: string; alias: string }
  onEdit: () => void
  onRemove: () => void
}

/** 群集行操作菜单。 */
function ClusterMenu({ cluster, onEdit, onRemove }: ClusterMenuProps) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon-sm" aria-label={`${cluster.alias} 的操作`}>
          <IconDotsVertical size={16} stroke={1.5} />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-40">
        <DropdownMenuItem onSelect={onEdit}>
          <IconEdit size={16} stroke={1.5} />
          编辑
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem onSelect={onRemove} className="text-danger">
          <IconTrash size={16} stroke={1.5} />
          删除
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
