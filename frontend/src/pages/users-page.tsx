import { useCallback, useMemo, useState } from "react"
import { IconPencil, IconPlus, IconSearch, IconUsers } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { CursorPagination } from "@/components/ui/cursor-pagination"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
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
  PATCH_USER_STATUS_MUTATION,
  ROLES_QUERY,
  USERS_QUERY,
  type RoleNode,
  type UserNode,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

type Tab = "users" | "roles"

/** 用户状态：1 正常 2 禁用 3 未激活。 */
const userStatusLabels: Record<number, string> = {
  1: "正常",
  2: "禁用",
  3: "未激活",
}

export interface UsersPageProps {
  onCreateUser: () => void
  onEditUser: (uuid: string) => void
}

/** 用户与角色页：账号清单与角色定义。 */
export function UsersPage({ onCreateUser, onEditUser }: UsersPageProps) {
  const [tab, setTab] = useState<Tab>("users")
  const [keyword, setKeyword] = useState("")
  const [after, setAfter] = useState<string | null>(null)
  const [stack, setStack] = useState<(string | null)[]>([])
  const [actionError, setActionError] = useState<string | null>(null)
  const [busyUser, setBusyUser] = useState<string | null>(null)

  const users = useGraphQL<
    {
      users: {
        totalCount: number
        pageInfo: { hasNextPage: boolean; hasPreviousPage: boolean; startCursor: string; endCursor: string }
        edges: { cursor: string; node: UserNode }[]
      } | null
    },
    { first: number; after?: string }
  >(USERS_QUERY, { first: 10, after: after ?? undefined }, { users: null })

  const roles = useGraphQL<{ roles: RoleNode[] }>(ROLES_QUERY, undefined, { roles: [] })

  const connection = users.data?.users ?? null
  const userList = useMemo(() => {
    const needle = keyword.trim().toLowerCase()
    const list = connection?.edges.map((edge) => edge.node) ?? []
    if (!needle) return list
    return list.filter(
      (user) =>
        user.Email.toLowerCase().includes(needle) ||
        (user.Name ?? "").toLowerCase().includes(needle),
    )
  }, [connection, keyword])
  const roleList = roles.data?.roles ?? []

  const toggleStatus = async (user: UserNode) => {
    const next = user.Status === 1 ? 2 : 1
    setBusyUser(user.UUID)
    setActionError(null)
    try {
      await graphqlRequest(
        PATCH_USER_STATUS_MUTATION,
        { input: { UserUUID: user.UUID, Status: next } },
        getToken(),
      )
      users.refetch()
    } catch (cause) {
      setActionError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setBusyUser(null)
    }
  }

  const nextPage = useCallback(() => {
    setStack((items) => [...items, after])
    setAfter(connection?.pageInfo.endCursor || null)
  }, [after, connection])

  const previousPage = useCallback(() => {
    const previousCursor = stack[stack.length - 1] ?? null
    setStack((items) => items.slice(0, -1))
    setAfter(previousCursor)
  }, [stack])

  return (
    <>
      <PageHeader
        title="用户与角色"
        description="账号、角色与审核关系"
        action={
          <div className="flex items-center gap-2">
            <Badge variant="default" className="px-2 py-1">
              <IconUsers size={14} stroke={1.5} className="mr-1" />
              {userList.length} 个账号 / {roleList.length} 个角色
            </Badge>
            <Button onClick={onCreateUser}>
              <IconPlus size={16} stroke={1.5} />
              新建账号
            </Button>
          </div>
        }
      />

      {actionError ? (
        <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{actionError}</p>
      ) : null}

      <SectionCard title="账号与角色" description="角色决定能做什么，审核关系决定能审谁的工单" bodyClassName="p-0">
        <div className="flex gap-4 border-b border-border px-5 pt-3">
          {(
            [
              ["users", "用户", userList.length],
              ["roles", "角色", roleList.length],
            ] as const
          ).map(([key, label, count]) => (
            <button
              key={key}
              type="button"
              onClick={() => setTab(key)}
              className={
                tab === key
                  ? "-mb-px border-b-2 border-primary pb-2 text-sm font-medium"
                  : "-mb-px border-b-2 border-transparent pb-2 text-sm text-muted-foreground"
              }
            >
              {label}（{count}）
            </button>
          ))}
        </div>

        {tab === "users" ? (
          <>
            <div className="flex flex-wrap items-end gap-3 border-b border-border px-5 py-4">
              <div className="flex flex-col gap-1.5 sm:w-64">
                <Label htmlFor="user-keyword">关键词</Label>
                <div className="relative">
                  <IconSearch
                    size={16}
                    stroke={1.5}
                    className="pointer-events-none absolute top-1/2 left-2.5 -translate-y-1/2 text-muted-foreground"
                  />
                  <Input
                    id="user-keyword"
                    value={keyword}
                    onChange={(event) => {
                      setKeyword(event.target.value)
                      setAfter(null)
                      setStack([])
                    }}
                    placeholder="邮箱或姓名"
                    className="pl-8"
                  />
                </div>
              </div>
            </div>

            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>邮箱</TableHead>
                  <TableHead>姓名</TableHead>
                  <TableHead>角色</TableHead>
                  <TableHead>可审核人</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead className="w-32">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {userList.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="py-8 text-center text-sm text-muted-foreground">
                      没有匹配的账号。
                    </TableCell>
                  </TableRow>
                ) : null}
                {userList.map((user) => (
                  <TableRow key={user.UUID}>
                    <TableCell className="font-medium">{user.Email}</TableCell>
                    <TableCell>{user.Name ?? "-"}</TableCell>
                    <TableCell>
                      <div className="flex flex-wrap gap-1">
                        {(user.Roles ?? []).map((role) => (
                          <Badge key={role.UUID} variant="muted">
                            {role.Name}
                          </Badge>
                        ))}
                      </div>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {(user.Reviewers ?? []).map((reviewer) => reviewer.Name ?? reviewer.Email).join("、") ||
                        "-"}
                    </TableCell>
                    <TableCell>
                      <Badge variant={user.Status === 1 ? "success" : "muted"}>
                        {userStatusLabels[user.Status] ?? user.Status}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        <Button
                          variant="ghost"
                          size="icon-sm"
                          aria-label={`编辑 ${user.Email}`}
                          onClick={() => onEditUser(user.UUID)}
                        >
                          <IconPencil size={16} stroke={1.5} />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          disabled={busyUser === user.UUID}
                          onClick={() => void toggleStatus(user)}
                        >
                          {user.Status === 1 ? "禁用" : "启用"}
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            <CursorPagination
              total={connection?.totalCount ?? 0}
              hasPreviousPage={connection?.pageInfo.hasPreviousPage ?? false}
              hasNextPage={connection?.pageInfo.hasNextPage ?? false}
              onPrevious={previousPage}
              onNext={nextPage}
            />
          </>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>角色</TableHead>
                <TableHead>说明</TableHead>
                <TableHead>可管理的分组</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {roleList.map((role) => (
                <TableRow key={role.UUID}>
                  <TableCell className="font-medium">{role.Name}</TableCell>
                  <TableCell className="text-muted-foreground">{role.Description || "-"}</TableCell>
                  <TableCell>
                    <div className="flex flex-wrap gap-1">
                      {(role.Groups ?? []).map((group) => (
                        <Badge key={group.UUID} variant="muted">
                          {group.Name}
                        </Badge>
                      ))}
                      {(role.Groups ?? []).length === 0 ? "-" : null}
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </SectionCard>
    </>
  )
}
