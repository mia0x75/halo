import { useMemo, useState } from "react"
import { IconArrowLeft } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Input } from "@/components/ui/input"
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
  CREATE_USER_MUTATION,
  GRANT_CLUSTERS_MUTATION,
  GRANT_REVIEWERS_MUTATION,
  GRANT_ROLES_MUTATION,
  ROLES_QUERY,
  UPDATE_USER_MUTATION,
  USERS_QUERY,
  type ClusterNode,
  type RoleNode,
  type UserNode,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

export interface UserFormPageProps {
  mode: "create" | "edit"
  /** 编辑时带的用户 UUID */
  uuid?: string
  onBack: () => void
  onSaved: () => void
}

interface UserFormProps {
  mode: "create" | "edit"
  uuid?: string
  initial?: UserNode
  onSaved: () => void
}

/** 用户表单：基本资料 + 角色、审核人、群集授权一次提交。 */
function UserForm({ mode, uuid, initial, onSaved }: UserFormProps) {
  const roles = useGraphQL<{ roles: RoleNode[] }>(ROLES_QUERY, undefined, { roles: [] })
  const users = useGraphQL<
    { users: { totalCount: number; edges: { cursor: string; node: UserNode }[] } | null },
    { first: number }
  >(USERS_QUERY, { first: 100 }, { users: null })
  const clusters = useGraphQL<
    { clusters: { totalCount: number; edges: { cursor: string; node: ClusterNode }[] } | null },
    { first: number }
  >(CLUSTERS_QUERY, { first: 100 }, { clusters: null })

  const [email, setEmail] = useState(initial?.Email ?? "")
  const [name, setName] = useState(initial?.Name ?? "")
  const [phone, setPhone] = useState("")
  const [password, setPassword] = useState("")
  const [status, setStatus] = useState(String(initial?.Status ?? 1))
  const [roleUUIDs, setRoleUUIDs] = useState<string[]>(
    (initial?.Roles ?? []).map((role) => role.UUID),
  )
  const [reviewerUUIDs, setReviewerUUIDs] = useState<string[]>(
    (initial?.Reviewers ?? []).map((reviewer) => reviewer.UUID),
  )
  const [clusterUUIDs, setClusterUUIDs] = useState<string[]>(
    (initial?.GrantedClusters ?? []).map((cluster) => cluster.UUID),
  )
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  const roleList = useMemo(() => roles.data?.roles ?? [], [roles.data])
  const userList = useMemo(
    () => users.data?.users?.edges.map((edge) => edge.node) ?? [],
    [users.data],
  )
  const clusterList = useMemo(
    () => clusters.data?.clusters?.edges.map((edge) => edge.node) ?? [],
    [clusters.data],
  )

  const toggle = (
    list: string[],
    setList: (next: string[]) => void,
    uuid: string,
    on: boolean,
  ) => {
    setList(on ? [...new Set([...list, uuid])] : list.filter((item) => item !== uuid))
  }

  const submit = async () => {
    if (!email.trim()) {
      setError("邮箱不能为空")
      return
    }
    if (mode === "create" && password.length < 3) {
      setError("新建账号的密码至少 3 位")
      return
    }
    setSaving(true)
    setError(null)
    try {
      if (mode === "create") {
        await graphqlRequest(
          CREATE_USER_MUTATION,
          {
            input: {
              Email: email.trim(),
              Password: password,
              Name: name.trim(),
              Phone: phone.trim() ? Number(phone.trim()) : 0,
              RoleUUIDs: roleUUIDs,
              ClusterUUIDs: clusterUUIDs,
              ReviewerUUIDs: reviewerUUIDs,
              AvatarUUID: "",
              Status: Number(status),
            },
          },
          getToken(),
        )
      } else if (uuid) {
        await graphqlRequest(
          UPDATE_USER_MUTATION,
          {
            input: {
              UserUUID: uuid,
              Email: email.trim(),
              Password: password,
              Status: Number(status),
              Name: name.trim(),
              Phone: phone.trim() ? Number(phone.trim()) : 0,
              AvatarUUID: "",
            },
          },
          getToken(),
        )
        await graphqlRequest(
          GRANT_ROLES_MUTATION,
          { input: { UserUUID: uuid, RoleUUIDs: roleUUIDs } },
          getToken(),
        )
        await graphqlRequest(
          GRANT_REVIEWERS_MUTATION,
          { input: { UserUUID: uuid, ReviewerUUIDs: reviewerUUIDs } },
          getToken(),
        )
        await graphqlRequest(
          GRANT_CLUSTERS_MUTATION,
          { input: { UserUUID: uuid, ClusterUUIDs: clusterUUIDs } },
          getToken(),
        )
      }
      onSaved()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setSaving(false)
    }
  }

  return (
    <>
      <div className="grid gap-4 sm:grid-cols-2">
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="user-email">邮箱</Label>
          <Input
            id="user-email"
            type="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            placeholder="user@example.com"
          />
        </div>
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="user-name">姓名</Label>
          <Input
            id="user-name"
            value={name}
            onChange={(event) => setName(event.target.value)}
            placeholder="张三"
          />
        </div>
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="user-phone">电话</Label>
          <Input
            id="user-phone"
            value={phone}
            onChange={(event) => setPhone(event.target.value)}
            placeholder="选填"
          />
        </div>
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="user-status">状态</Label>
          <Select value={status} onValueChange={(value) => setStatus(value)}>
            <SelectTrigger id="user-status">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1">正常</SelectItem>
              <SelectItem value="2">禁用</SelectItem>
              <SelectItem value="3">待验证</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="flex flex-col gap-1.5 sm:col-span-2">
          <Label htmlFor="user-password">
            {mode === "create" ? "初始密码" : "重置密码"}
          </Label>
          <Input
            id="user-password"
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            placeholder={mode === "create" ? "至少 3 位" : "留空表示不修改"}
          />
        </div>
      </div>

      <div className="mt-6 grid gap-6 lg:grid-cols-3">
        <div className="flex flex-col gap-2">
          <p className="text-sm font-medium text-strong">角色</p>
          <p className="text-xs text-muted-foreground">决定能做什么</p>
          {roleList.map((role) => (
            <label
              key={role.UUID}
              className="flex items-start gap-2 rounded-md border border-border px-3 py-2 text-sm"
            >
              <Checkbox
                id={`role-${role.UUID}`}
                checked={roleUUIDs.includes(role.UUID)}
                onCheckedChange={(value) => toggle(roleUUIDs, setRoleUUIDs, role.UUID, value === true)}
              />
              <span className="flex min-w-0 cursor-pointer flex-col">
                <span className="truncate font-medium text-strong">{role.Name}</span>
                <span className="truncate text-xs text-muted-foreground">
                  {role.Description || "-"}
                </span>
              </span>
            </label>
          ))}
        </div>

        <div className="flex flex-col gap-2">
          <p className="text-sm font-medium text-strong">审核人</p>
          <p className="text-xs text-muted-foreground">决定谁来审这个账号的工单</p>
          {userList
            .filter((user) => user.UUID !== uuid)
            .map((user) => (
              <label
                key={user.UUID}
                className="flex items-center gap-2 rounded-md border border-border px-3 py-2 text-sm"
              >
                <Checkbox
                  id={`reviewer-${user.UUID}`}
                  checked={reviewerUUIDs.includes(user.UUID)}
                  onCheckedChange={(value) =>
                    toggle(reviewerUUIDs, setReviewerUUIDs, user.UUID, value === true)
                  }
                />
                <span className="flex min-w-0 cursor-pointer flex-col">
                  <span className="truncate font-medium text-strong">
                    {user.Name || user.Email}
                  </span>
                  <span className="truncate text-xs text-muted-foreground">{user.Email}</span>
                </span>
              </label>
            ))}
        </div>

        <div className="flex flex-col gap-2">
          <p className="text-sm font-medium text-strong">群集授权</p>
          <p className="text-xs text-muted-foreground">直接授权，可在「系统设置」里按分组授权</p>
          {clusterList.map((cluster) => (
            <label
              key={cluster.UUID}
              className="flex items-center gap-2 rounded-md border border-border px-3 py-2 text-sm"
            >
              <Checkbox
                id={`cluster-${cluster.UUID}`}
                checked={clusterUUIDs.includes(cluster.UUID)}
                onCheckedChange={(value) =>
                  toggle(clusterUUIDs, setClusterUUIDs, cluster.UUID, value === true)
                }
              />
              <span className="flex min-w-0 cursor-pointer flex-col">
                <span className="truncate font-medium text-strong">{cluster.Alias}</span>
                <span className="truncate text-xs text-muted-foreground">
                  {cluster.Group ?? "未分组"}
                </span>
              </span>
            </label>
          ))}
        </div>
      </div>

      {error ? (
        <p className="mt-4 rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
      ) : null}

      <div className="mt-6 flex justify-end gap-2">
        <Button onClick={() => void submit()} disabled={saving}>
          {mode === "create" ? "创建账号" : "保存"}
        </Button>
      </div>
    </>
  )
}

/** 新建/编辑账号页。表单是整页，不是弹窗。 */
export function UserFormPage({ mode, uuid, onBack, onSaved }: UserFormPageProps) {
  const remote = useGraphQL<
    { users: { totalCount: number; edges: { cursor: string; node: UserNode }[] } | null },
    { first: number }
  >(USERS_QUERY, { first: 100 }, { users: null })

  const initial = useMemo(() => {
    if (mode !== "edit" || !uuid) return undefined
    return remote.data?.users?.edges.find((edge) => edge.node.UUID === uuid)?.node
  }, [mode, remote.data, uuid])

  return (
    <>
      <PageHeader
        title={mode === "create" ? "新建账号" : "编辑账号"}
        description="账号资料、角色、审核人与群集授权一次配好"
        action={
          <Button variant="outline" onClick={onBack}>
            <IconArrowLeft size={16} stroke={1.5} />
            返回用户列表
          </Button>
        }
      />

      <SectionCard>
        {mode === "edit" && !initial ? (
          <p className="text-sm text-muted-foreground">正在读取账号信息…</p>
        ) : (
          <UserForm
            key={initial?.UUID ?? mode}
            mode={mode}
            uuid={uuid}
            initial={initial}
            onSaved={onSaved}
          />
        )}
      </SectionCard>
    </>
  )
}
