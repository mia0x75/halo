import { useMemo, useState } from "react"
import { IconCheck, IconInfoCircle, IconUsers } from "@tabler/icons-react"

import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { useGraphQL } from "@/hooks/use-graphql"
import { graphqlRequest } from "@/lib/api/client"
import {
  CLUSTERS_QUERY,
  GRANT_CLUSTERS_MUTATION,
  GRANT_GROUPS_MUTATION,
  GROUPS_QUERY,
  ROLES_QUERY,
  USERS_QUERY,
  type ClusterGrantNode,
  type ClusterNode,
  type GroupNode,
  type RoleNode,
  type UserNode,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

type Scope = "role" | "user"

export interface AccessPanelProps {
  scope: Scope
}

/**
 * 访问授权：角色→分组、用户→群集。
 * 提交时把勾选结果整份写入，后端按「替换」语义处理，无需先收回。
 */
export function AccessPanel({ scope }: AccessPanelProps) {
  const groups = useGraphQL<{ groups: GroupNode[] }>(GROUPS_QUERY, undefined, { groups: [] })
  const roles = useGraphQL<{ roles: RoleNode[] }>(ROLES_QUERY, undefined, { roles: [] })
  const users = useGraphQL<
    {
      users: {
        totalCount: number
        edges: { cursor: string; node: UserNode }[]
      } | null
    },
    { first: number }
  >(USERS_QUERY, { first: 100 }, { users: null })
  const clusters = useGraphQL<
    {
      clusters: {
        totalCount: number
        edges: { cursor: string; node: ClusterNode }[]
      } | null
    },
    { first: number }
  >(CLUSTERS_QUERY, { first: 100 }, { clusters: null })

  const [targetUUID, setTargetUUID] = useState("")
  const [checked, setChecked] = useState<string[]>([])
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [saved, setSaved] = useState(false)

  const groupList = useMemo(() => groups.data?.groups ?? [], [groups.data])
  const roleList = useMemo(() => roles.data?.roles ?? [], [roles.data])
  const userList = useMemo(() => users.data?.users?.edges.map((edge) => edge.node) ?? [], [users.data])
  const clusterList = useMemo(
    () => clusters.data?.clusters?.edges.map((edge) => edge.node) ?? [],
    [clusters.data],
  )

  const selectedRole = roleList.find((role) => role.UUID === targetUUID)
  const selectedUser = userList.find((user) => user.UUID === targetUUID)

  // 生效权限 = 直接授权 ∪ 角色继承。管理面必须分开呈现两层，
  // 否则「取消直接授权」会被误读成「断访问」。
  const grantsByCluster = useMemo(() => {
    const map = new Map<string, ClusterGrantNode>()
    for (const grant of selectedUser?.ClusterGrants ?? []) {
      map.set(grant.Cluster.UUID, grant)
    }
    return map
  }, [selectedUser])

  const roleGroups = useMemo(
    () => new Set((selectedRole?.Groups ?? []).map((group) => group.UUID)),
    [selectedRole],
  )

  const optionItems = useMemo(
    () =>
      scope === "role"
        ? groupList.map((group) => ({
            uuid: group.UUID,
            name: group.Name,
            extra: group.Description || "-",
            // Role→group has no inheritance layer: the grant itself is what
            // makes the group effective, so read it back from the saved role.
            effective: roleGroups.has(group.UUID),
            sources: [] as string[],
          }))
        : clusterList.map((cluster) => {
            const grant = grantsByCluster.get(cluster.UUID)
            return {
              uuid: cluster.UUID,
              name: cluster.Alias,
              extra: cluster.Group ?? "未分组",
              effective: Boolean(grant?.Direct || grant?.Inherited),
              sources: grant?.Sources ?? [],
            }
          }),
    [clusterList, grantsByCluster, groupList, roleGroups, scope],
  )
  const current =
    scope === "role"
      ? (selectedRole?.Groups ?? []).map((group) => group.UUID)
      : (selectedUser?.GrantedClusters ?? []).map((cluster) => cluster.UUID)

  const switchTarget = (uuid: string) => {
    setTargetUUID(uuid)
    setSaved(false)
    setError(null)
    if (scope === "role") {
      const role = roleList.find((item) => item.UUID === uuid)
      setChecked((role?.Groups ?? []).map((group) => group.UUID))
    } else {
      const user = userList.find((item) => item.UUID === uuid)
      setChecked((user?.GrantedClusters ?? []).map((cluster) => cluster.UUID))
    }
  }

  const toggle = (uuid: string, on: boolean) => {
    setSaved(false)
    setChecked((previous) =>
      on ? [...new Set([...previous, uuid])] : previous.filter((item) => item !== uuid),
    )
  }

  const save = async () => {
    if (!targetUUID) return
    setSaving(true)
    setError(null)
    try {
      if (scope === "role") {
        await graphqlRequest(
          GRANT_GROUPS_MUTATION,
          { input: { RoleUUID: targetUUID, GroupUUIDs: checked } },
          getToken(),
        )
        roles.refetch()
      } else {
        await graphqlRequest(
          GRANT_CLUSTERS_MUTATION,
          { input: { UserUUID: targetUUID, ClusterUUIDs: checked } },
          getToken(),
        )
        users.refetch()
      }
      setSaved(true)
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setSaving(false)
    }
  }

  const dirty = useMemo(() => {
    const before = [...current].sort()
    const after = [...checked].sort()
    return before.length !== after.length || before.some((item, index) => item !== after[index])
  }, [checked, current])

  return (
    <>
      {error ? (
        <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
      ) : null}
      {saved && !error ? (
        <p className="rounded-md bg-success/10 px-3 py-2 text-sm text-success">授权已保存</p>
      ) : null}

      <SectionCard
        title={scope === "role" ? "角色授权" : "用户群集授权"}
        description={
          scope === "role"
            ? "角色获得分组授权后，持有该角色的成员才能在这些分组的群集上提交工单"
            : "生效权限 = 直接授权 + 角色继承的分组权限，两者取并集"
        }
      >
        <div className="flex flex-col gap-4">
          <p className="text-xs text-muted-foreground">
            {scope === "role"
              ? "先选角色，再勾选它能管理的分组"
              : "勾选的是「直接授权」；继承自角色分组的权限由「角色授权」决定"}
          </p>
          <div className="flex flex-col gap-1.5 sm:w-96">
            <Label htmlFor="access-target">
              {scope === "role" ? "角色" : "用户"}
            </Label>
            <Select
              value={targetUUID || "none"}
              onValueChange={(value) => switchTarget(value === "none" ? "" : value)}
            >
              <SelectTrigger id="access-target">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">请选择</SelectItem>
                {scope === "role"
                  ? roleList.map((role) => (
                      <SelectItem key={role.UUID} value={role.UUID}>
                        {role.Name}（{role.Description || "无说明"}）
                      </SelectItem>
                    ))
                  : userList.map((user) => (
                      <SelectItem key={user.UUID} value={user.UUID}>
                        {user.Name || user.Email}
                      </SelectItem>
                    ))}
              </SelectContent>
            </Select>
          </div>

          {targetUUID ? (
            <div className="flex flex-col gap-2">
              <p className="text-sm text-muted-foreground">
                已授权 {checked.length} 项 / 共 {optionItems.length} 项
              </p>
              <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
                {optionItems.map((option) => (
                  <div
                    key={option.uuid}
                    className="flex items-start gap-2 rounded-md border border-border px-3 py-2 text-sm"
                  >
                    <Checkbox
                      id={`grant-${option.uuid}`}
                      checked={checked.includes(option.uuid)}
                      onCheckedChange={(value) => toggle(option.uuid, value === true)}
                    />
                    <label
                      htmlFor={`grant-${option.uuid}`}
                      className="flex min-w-0 cursor-pointer flex-col gap-1"
                    >
                      <span className="flex items-center gap-2">
                        <span className="truncate font-medium text-strong">{option.name}</span>
                        {option.effective ? (
                          <Badge variant="success">生效中</Badge>
                        ) : (
                          <Badge variant="muted">无权限</Badge>
                        )}
                      </span>
                      <span className="truncate text-xs text-muted-foreground">{option.extra}</span>
                      {option.sources.length > 0 ? (
                        <span className="truncate text-xs text-info">
                          经由 {option.sources.join("、")} 继承
                        </span>
                      ) : null}
                    </label>
                  </div>
                ))}
              </div>
              <div className="flex items-center gap-2">
                <Button onClick={() => void save()} disabled={saving || !dirty}>
                  <IconCheck size={16} stroke={1.5} />
                  保存授权
                </Button>
                {dirty ? (
                  <span className="text-sm text-muted-foreground">有未保存的改动</span>
                ) : null}
              </div>
              {scope === "user" ? (
                <p className="flex items-start gap-2 rounded-md bg-surface-secondary px-3 py-2 text-xs text-muted-foreground">
                  <IconInfoCircle size={14} stroke={1.5} className="mt-0.5 shrink-0" />
                  <span>
                    取消「直接授权」只会撤掉这一层。若该账号的角色仍被授予了群集所属分组，
                    权限依然生效——要真正收回，需要到「角色授权」里调整。
                  </span>
                </p>
              ) : null}
            </div>
          ) : (
            <p className="flex items-center gap-2 text-sm text-muted-foreground">
              <IconUsers size={16} stroke={1.5} />
              选一个对象开始配置。
            </p>
          )}
        </div>
      </SectionCard>
    </>
  )
}
