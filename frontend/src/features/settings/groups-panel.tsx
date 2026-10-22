import { useMemo, useState } from "react"
import { IconPencil, IconPlus, IconTrash } from "@tabler/icons-react"

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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useGraphQL } from "@/hooks/use-graphql"
import { graphqlRequest } from "@/lib/api/client"
import {
  CLUSTERS_QUERY,
  GROUPS_QUERY,
  REMOVE_GROUP_MUTATION,
  UPDATE_CLUSTER_MUTATION,
  type ClusterNode,
  type GroupNode,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

export interface GroupsPanelProps {
  onCreate: () => void
  onEdit: (uuid: string) => void
}

/** 群集分组表格：新建与编辑走独立页面，删除用确认弹窗。 */
export function GroupsPanel({ onCreate, onEdit }: GroupsPanelProps) {
  const remote = useGraphQL<{ groups: GroupNode[] }>(GROUPS_QUERY, undefined, { groups: [] })
  const clusters = useGraphQL<
    {
      clusters: {
        totalCount: number
        edges: { cursor: string; node: ClusterNode }[]
      } | null
    },
    { first: number }
  >(CLUSTERS_QUERY, { first: 100 }, { clusters: null })
  const [removing, setRemoving] = useState<GroupNode | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [movingUUID, setMovingUUID] = useState<string | null>(null)
  const [moveError, setMoveError] = useState<string | null>(null)

  const groups = useMemo(() => remote.data?.groups ?? [], [remote.data])
  const clusterList = useMemo(
    () => clusters.data?.clusters?.edges.map((edge) => edge.node) ?? [],
    [clusters.data],
  )
  const memberCount = useMemo(() => {
    const counts = new Map<string, number>()
    for (const cluster of clusterList) {
      if (!cluster.GroupUUID) continue
      counts.set(cluster.GroupUUID, (counts.get(cluster.GroupUUID) ?? 0) + 1)
    }
    return counts
  }, [clusterList])

  // 归组只改分组，密码留空表示不修改。
  const moveCluster = async (cluster: ClusterNode, groupUUID: string) => {
    setMovingUUID(cluster.UUID)
    setMoveError(null)
    try {
      await graphqlRequest(
        UPDATE_CLUSTER_MUTATION,
        {
          input: {
            ClusterUUID: cluster.UUID,
            Host: cluster.Host,
            IP: cluster.IP,
            Port: cluster.Port,
            Alias: cluster.Alias,
            User: cluster.User,
            Status: cluster.Status,
            Type: cluster.Type,
            Password: "",
            GroupUUID: groupUUID,
          },
        },
        getToken(),
      )
      clusters.refetch()
    } catch (cause) {
      setMoveError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setMovingUUID(null)
    }
  }

  const remove = async (group: GroupNode) => {
    setBusy(true)
    setError(null)
    try {
      await graphqlRequest(REMOVE_GROUP_MUTATION, { id: group.UUID }, getToken())
      setRemoving(null)
      remote.refetch()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setBusy(false)
    }
  }

  return (
    <SectionCard
      title="群集分组"
      description="群集按分组归类，角色获得分组授权后其成员才能在该分组的群集上提交工单"
      action={
        <Button onClick={onCreate}>
          <IconPlus size={16} stroke={1.5} />
          新建分组
        </Button>
      }
      bodyClassName="p-0"
    >
      {error ? (
        <p className="mb-3 px-5 text-sm text-danger">{error}</p>
      ) : null}
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>分组名称</TableHead>
            <TableHead>说明</TableHead>
            <TableHead>群集数</TableHead>
            <TableHead>状态</TableHead>
            <TableHead>创建时间</TableHead>
            <TableHead className="w-24"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {groups.length === 0 ? (
            <TableRow>
              <TableCell colSpan={6} className="py-8 text-center text-sm text-muted-foreground">
                还没有分组，先新建一个再把群集归进去。
              </TableCell>
            </TableRow>
          ) : null}
          {groups.map((group) => (
            <TableRow key={group.UUID}>
              <TableCell className="font-medium text-strong">{group.Name}</TableCell>
              <TableCell className="text-muted-foreground">{group.Description || "-"}</TableCell>
              <TableCell>{memberCount.get(group.UUID) ?? 0}</TableCell>
              <TableCell>
                <Badge variant={group.Status === 1 ? "success" : "muted"}>
                  {group.Status === 1 ? "正常" : "停用"}
                </Badge>
              </TableCell>
              <TableCell className="text-sm text-muted-foreground">
                {new Date(group.CreateAt * 1000).toLocaleString("zh-CN", { hour12: false })}
              </TableCell>
              <TableCell>
                <div className="flex items-center gap-1">
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    aria-label={`编辑 ${group.Name}`}
                    onClick={() => onEdit(group.UUID)}
                  >
                    <IconPencil size={16} stroke={1.5} />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    aria-label={`删除 ${group.Name}`}
                    onClick={() => setRemoving(group)}
                  >
                    <IconTrash size={16} stroke={1.5} />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      <SectionCard
        title="群集归组"
        description="群集归属的分组决定谁能访问它：角色被授予分组后，持有该角色的成员即可访问"
        className="mt-4"
        bodyClassName="p-0"
      >
        {moveError ? (
          <p className="px-5 pt-4 text-sm text-danger">{moveError}</p>
        ) : null}
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>别名</TableHead>
              <TableHead>地址</TableHead>
              <TableHead>类型</TableHead>
              <TableHead>所属分组</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {clusterList.length === 0 ? (
              <TableRow>
                <TableCell colSpan={4} className="py-8 text-center text-sm text-muted-foreground">
                  还没有接入群集。
                </TableCell>
              </TableRow>
            ) : null}
            {clusterList.map((cluster) => (
              <TableRow key={cluster.UUID}>
                <TableCell className="font-medium">{cluster.Alias}</TableCell>
                <TableCell className="text-muted-foreground">
                  {cluster.IP}:{cluster.Port}
                </TableCell>
                <TableCell className="text-muted-foreground">{cluster.Type}</TableCell>
                <TableCell>
                  <Select
                    value={cluster.GroupUUID ?? "none"}
                    disabled={movingUUID === cluster.UUID}
                    onValueChange={(value) =>
                      void moveCluster(cluster, value === "none" ? "" : value)
                    }
                  >
                    <SelectTrigger
                      aria-label={`${cluster.Alias} 的分组`}
                      className="h-8 max-w-40"
                    >
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="none">未分组</SelectItem>
                      {groups.map((group) => (
                        <SelectItem key={group.UUID} value={group.UUID}>
                          {group.Name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </SectionCard>

      <AlertDialog open={removing !== null} onOpenChange={(open) => !open && setRemoving(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="text-base font-medium text-strong">
              删除分组
            </AlertDialogTitle>
            <AlertDialogDescription className="mt-1 text-sm text-muted-foreground">
              分组下如果还有群集，删除会被拒绝；角色对该分组的授权也会一并失效。确认删除「
              {removing?.Name}」？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <div className="mt-5 flex justify-end gap-2">
            <AlertDialogCancel asChild>
              <Button variant="outline">取消</Button>
            </AlertDialogCancel>
            <AlertDialogAction asChild>
              <Button
                variant="danger"
                disabled={busy}
                onClick={() => {
                  if (removing) void remove(removing)
                }}
              >
                删除
              </Button>
            </AlertDialogAction>
          </div>
        </AlertDialogContent>
      </AlertDialog>
    </SectionCard>
  )
}
