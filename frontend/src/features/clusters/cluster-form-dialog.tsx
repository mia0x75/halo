import { useState } from "react"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import { graphqlRequest } from "@/lib/api/client"
import {
  CREATE_CLUSTER_MUTATION,
  GROUPS_QUERY,
  UPDATE_CLUSTER_MUTATION,
  type GroupNode,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"
import { useGraphQL } from "@/hooks/use-graphql"

export interface ClusterFormValues {
  alias: string
  host: string
  ip: string
  port: number
  user: string
  password: string
  status: number
  type: string
  groupUUID: string
}

export interface ClusterFormDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  /** 传入表示编辑，不传表示新建 */
  cluster?: ClusterFormValues & { uuid: string }
  onSaved?: () => void
  trigger?: React.ReactNode
}

const emptyValues: ClusterFormValues = {
  alias: "",
  host: "",
  ip: "",
  port: 3306,
  user: "",
  password: "",
  status: 1,
  type: "mysql",
  groupUUID: "",
}

/** 群集新增/编辑表单。 */
export function ClusterFormDialog({
  open,
  onOpenChange,
  cluster,
  onSaved,
  trigger,
}: ClusterFormDialogProps) {
  const [values, setValues] = useState<ClusterFormValues>(emptyValues)
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  const groups = useGraphQL<{ groups: GroupNode[] }>(GROUPS_QUERY, undefined, { groups: [] })

  // 打开时用当前值重置表单（事件里赋值，不放 effect）。
  const handleOpenChange = (next: boolean) => {
    if (next) {
      setValues(cluster ?? emptyValues)
      setError(null)
    }
    onOpenChange(next)
  }

  const set = <K extends keyof ClusterFormValues>(key: K, value: ClusterFormValues[K]) => {
    setValues((prev) => ({ ...prev, [key]: value }))
  }

  const submit = async () => {
    if (!getToken()) {
      setError("尚未登录，无法保存。请先在右上角登录后重试。")
      return
    }

    setSaving(true)
    setError(null)
    try {
      const input = {
        Alias: values.alias.trim(),
        Host: values.host.trim(),
        IP: values.ip.trim(),
        Port: values.port,
        User: values.user.trim(),
        Password: values.password,
        Status: values.status,
        Type: values.type,
        ...(values.groupUUID ? { GroupUUID: values.groupUUID } : {}),
        ...(cluster ? { ClusterUUID: cluster.uuid } : {}),
      }
      if (cluster) {
        await graphqlRequest(UPDATE_CLUSTER_MUTATION, { input }, getToken())
      } else {
        await graphqlRequest(CREATE_CLUSTER_MUTATION, { input }, getToken())
      }
      onOpenChange(false)
      onSaved?.()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setSaving(false)
    }
  }

  const valid =
    values.alias.trim() && values.host.trim() && values.ip.trim() && values.user.trim() && values.password

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        {trigger ?? (
          <Button>
            接入群集
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="max-w-2xl">
        <DialogTitle className="text-base font-medium text-strong">
          {cluster ? "编辑群集" : "接入群集"}
        </DialogTitle>
        <DialogDescription className="mt-1 text-sm text-muted-foreground">
          群集是目标数据库所在的主机，密码由服务端加密保存。
        </DialogDescription>

        <div className="mt-4 grid gap-3 sm:grid-cols-2">
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="cluster-alias">别名</Label>
            <Input
              id="cluster-alias"
              value={values.alias}
              onChange={(event) => set("alias", event.target.value)}
              placeholder="prod-trade-01"
            />
          </div>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="cluster-host">主机名</Label>
            <Input
              id="cluster-host"
              value={values.host}
              onChange={(event) => set("host", event.target.value)}
              placeholder="db-trade-01.internal"
            />
          </div>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="cluster-ip">IP 地址</Label>
            <Input
              id="cluster-ip"
              value={values.ip}
              onChange={(event) => set("ip", event.target.value)}
              placeholder="10.20.30.11"
            />
          </div>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="cluster-port">端口</Label>
            <Input
              id="cluster-port"
              type="number"
              value={values.port}
              onChange={(event) => set("port", Number(event.target.value))}
            />
          </div>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="cluster-user">连接账户</Label>
            <Input
              id="cluster-user"
              value={values.user}
              onChange={(event) => set("user", event.target.value)}
              placeholder="root"
            />
          </div>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="cluster-password">连接密码</Label>
            <Input
              id="cluster-password"
              type="password"
              autoComplete="new-password"
              value={values.password}
              onChange={(event) => set("password", event.target.value)}
            />
          </div>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="cluster-type">目标库类型</Label>
            <Select
              id="cluster-type"
              value={values.type}
              onChange={(event) => set("type", event.target.value)}
            >
              <option value="mysql">MySQL</option>
              <option value="postgres">PostgreSQL</option>
            </Select>
          </div>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="cluster-group">所属分组</Label>
            <Select
              id="cluster-group"
              value={values.groupUUID}
              onChange={(event) => set("groupUUID", event.target.value)}
            >
              <option value="">未分组</option>
              {(groups.data?.groups ?? []).map((group) => (
                <option key={group.UUID} value={group.UUID}>
                  {group.Name}
                </option>
              ))}
            </Select>
          </div>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="cluster-status">状态</Label>
            <Select
              id="cluster-status"
              value={String(values.status)}
              onChange={(event) => set("status", Number(event.target.value))}
            >
              <option value="1">正常</option>
              <option value="2">禁用</option>
            </Select>
          </div>
        </div>

        {error ? (
          <p className="mt-3 rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
        ) : null}

        <div className="mt-5 flex justify-end gap-2">
          <Button variant="outline" onClick={() => handleOpenChange(false)}>
            取消
          </Button>
          <Button disabled={!valid || saving} onClick={submit}>
            {saving ? "保存中…" : cluster ? "保存" : "接入"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
