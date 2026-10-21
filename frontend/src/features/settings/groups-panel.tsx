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
  GROUPS_QUERY,
  REMOVE_GROUP_MUTATION,
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
  const [removing, setRemoving] = useState<GroupNode | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  const groups = useMemo(() => remote.data?.groups ?? [], [remote.data])

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
            <TableHead>状态</TableHead>
            <TableHead>创建时间</TableHead>
            <TableHead className="w-24"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {groups.length === 0 ? (
            <TableRow>
              <TableCell colSpan={5} className="py-8 text-center text-sm text-muted-foreground">
                还没有分组，先新建一个再把群集归进去。
              </TableCell>
            </TableRow>
          ) : null}
          {groups.map((group) => (
            <TableRow key={group.UUID}>
              <TableCell className="font-medium text-strong">{group.Name}</TableCell>
              <TableCell className="text-muted-foreground">{group.Description || "-"}</TableCell>
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
