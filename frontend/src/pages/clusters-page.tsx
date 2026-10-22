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
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { StatusDot } from "@/components/ui/status-dot"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useClusters } from "@/hooks/use-clusters"
import { graphqlRequest } from "@/lib/api/client"
import { REMOVE_CLUSTER_MUTATION } from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

type ViewMode = "card" | "table"

export interface ClustersPageProps {
  /** 跳转到接入群集页 */
  onCreate: () => void
  /** 跳转到编辑群集页 */
  onEdit: (uuid: string) => void
}

export function ClustersPage({ onCreate, onEdit }: ClustersPageProps) {
  const [view, setView] = useState<ViewMode>("card")
  const [keyword, setKeyword] = useState("")
  const [type, setType] = useState("all")
  const [group, setGroup] = useState("all")
  const [removing, setRemoving] = useState<string | null>(null)
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

  const remove = async () => {
    if (!removing) return
    setError(null)
    try {
      await graphqlRequest(REMOVE_CLUSTER_MUTATION, { id: removing }, getToken())
      setRemoving(null)
      page.refetch()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
      setRemoving(null)
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
            <Button onClick={onCreate}>
              <IconPlus size={16} stroke={1.5} />
              接入群集
            </Button>
          </div>
        }
      />

      {error ? (
        <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
      ) : null}

      <SectionCard bodyClassName="p-0">
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
              value={type}
              onValueChange={(value) => {
                setType(value)
                page.reset()
              }}
            >
              <SelectTrigger id="cluster-type">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">全部类型</SelectItem>
                <SelectItem value="mysql">MySQL</SelectItem>
                <SelectItem value="postgres">PostgreSQL</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="flex flex-col gap-1.5 sm:w-44">
            <Label htmlFor="cluster-group">分组</Label>
            <Select
              value={group}
              onValueChange={(value) => {
                setGroup(value)
                page.reset()
              }}
            >
              <SelectTrigger id="cluster-group">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">全部分组</SelectItem>
                {groups.map((item) => (
                  <SelectItem key={item} value={item}>
                    {item}
                  </SelectItem>
                ))}
              </SelectContent>
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
                      alias={cluster.alias}
                      onEdit={() => onEdit(cluster.uuid)}
                      onRemove={() => setRemoving(cluster.uuid)}
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
                      alias={cluster.alias}
                      onEdit={() => onEdit(cluster.uuid)}
                      onRemove={() => setRemoving(cluster.uuid)}
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
          onPrevious={page.previous}
          onNext={page.next}
        />
      </SectionCard>

      <AlertDialog open={removing !== null} onOpenChange={(open) => !open && setRemoving(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="text-base font-medium text-strong">
              删除群集
            </AlertDialogTitle>
            <AlertDialogDescription className="mt-1 text-sm text-muted-foreground">
              删除后该群集不再可用，已有的工单记录保留。确认删除？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <div className="mt-5 flex justify-end gap-2">
            <AlertDialogCancel asChild>
              <Button variant="outline">取消</Button>
            </AlertDialogCancel>
            <AlertDialogAction asChild>
              <Button variant="danger" onClick={remove}>
                删除
              </Button>
            </AlertDialogAction>
          </div>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}

interface ClusterMenuProps {
  alias: string
  onEdit: () => void
  onRemove: () => void
}

/** 群集行操作菜单。 */
function ClusterMenu({ alias, onEdit, onRemove }: ClusterMenuProps) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon-sm" aria-label={`${alias} 的操作`}>
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
