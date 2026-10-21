import { useMemo, useState } from "react"
import { IconCheck, IconUsers } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import { useGraphQL } from "@/hooks/use-graphql"
import { graphqlRequest } from "@/lib/api/client"
import {
  CLUSTERS_QUERY,
  GRANT_CLUSTERS_MUTATION,
  GRANT_GROUPS_MUTATION,
  GROUPS_QUERY,
  ROLES_QUERY,
  USERS_QUERY,
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

  const optionItems = useMemo(
    () =>
      scope === "role"
        ? groupList.map((group) => ({
            uuid: group.UUID,
            name: group.Name,
            extra: group.Description || "-",
          }))
        : clusterList.map((cluster) => ({
            uuid: cluster.UUID,
            name: cluster.Alias,
            extra: cluster.Group ?? "未分组",
          })),
    [clusterList, groupList, scope],
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
      <PageHeader
        title={scope === "role" ? "角色授权" : "用户群集授权"}
        description={
          scope === "role"
            ? "角色获得分组授权后，持有该角色的成员才能在这些分组的群集上提交工单"
            : "直接授权给用户的群集，与角色继承的分组权限叠加"
        }
      />

      {error ? (
        <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
      ) : null}
      {saved && !error ? (
        <p className="rounded-md bg-success/10 px-3 py-2 text-sm text-success">授权已保存</p>
      ) : null}

      <SectionCard
        title="选择对象"
        description={
          scope === "role"
            ? "先选角色，再勾选它能管理的分组"
            : "先选用户，再勾选直接授权给他的群集"
        }
      >
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-1.5 sm:w-96">
            <Label htmlFor="access-target">
              {scope === "role" ? "角色" : "用户"}
            </Label>
            <Select
              id="access-target"
              value={targetUUID}
              onChange={(event) => switchTarget(event.target.value)}
            >
              <option value="">请选择</option>
              {scope === "role"
                ? roleList.map((role) => (
                    <option key={role.UUID} value={role.UUID}>
                      {role.Name}（{role.Description || "无说明"}）
                    </option>
                  ))
                : userList.map((user) => (
                    <option key={user.UUID} value={user.UUID}>
                      {user.Name || user.Email}
                    </option>
                  ))}
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
                      className="flex min-w-0 cursor-pointer flex-col"
                    >
                      <span className="truncate font-medium text-strong">{option.name}</span>
                      <span className="truncate text-xs text-muted-foreground">{option.extra}</span>
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
