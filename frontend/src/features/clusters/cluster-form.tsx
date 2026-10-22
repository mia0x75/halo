import { useState } from "react"
import { IconDeviceFloppy, IconPlugConnected } from "@tabler/icons-react"

import { Button } from "@/components/ui/button"
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
  CREATE_CLUSTER_MUTATION,
  GROUPS_QUERY,
  TEST_CLUSTER_QUERY,
  UPDATE_CLUSTER_MUTATION,
  type ClusterTestResult,
  type GroupNode,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

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

export interface ClusterFormProps {
  mode: "create" | "edit"
  /** 编辑时带上原记录 */
  initial?: ClusterFormValues & { uuid: string }
  onSaved: () => void
  onCancel: () => void
}

/** 群集接入/编辑表单。 */
export function ClusterForm({ mode, initial, onSaved, onCancel }: ClusterFormProps) {
  const [values, setValues] = useState<ClusterFormValues>(initial ?? emptyValues)
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)
  const [testing, setTesting] = useState(false)
  const [testResult, setTestResult] = useState<ClusterTestResult | null>(null)

  const groups = useGraphQL<{ groups: GroupNode[] }>(GROUPS_QUERY, undefined, { groups: [] })

  const set = <K extends keyof ClusterFormValues>(key: K, value: ClusterFormValues[K]) => {
    setValues((prev) => ({ ...prev, [key]: value }))
  }

  const testConnection = async () => {
    setTesting(true)
    setTestResult(null)
    setError(null)
    try {
      const result = await graphqlRequest<{ testCluster: ClusterTestResult }, { input: unknown }>(
        TEST_CLUSTER_QUERY,
        {
          input: {
            IP: values.ip.trim(),
            Port: values.port,
            User: values.user.trim(),
            Password: values.password,
            Type: values.type,
          },
        },
        getToken(),
      )
      setTestResult(result.testCluster)
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setTesting(false)
    }
  }

  const submit = async () => {
    setSaving(true)
    setError(null)
    try {
      const input = {
        ...(mode === "edit" && initial ? { ClusterUUID: initial.uuid } : {}),
        Alias: values.alias.trim(),
        Host: values.host.trim(),
        IP: values.ip.trim(),
        Port: values.port,
        User: values.user.trim(),
        // 编辑时留空表示保持原密码。
        ...(values.password ? { Password: values.password } : {}),
        Status: values.status,
        Type: values.type,
        ...(values.groupUUID ? { GroupUUID: values.groupUUID } : {}),
      }
      if (mode === "edit") {
        await graphqlRequest(UPDATE_CLUSTER_MUTATION, { input }, getToken())
      } else {
        await graphqlRequest(CREATE_CLUSTER_MUTATION, { input }, getToken())
      }
      onSaved()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setSaving(false)
    }
  }

  const connectable = values.ip.trim() && values.user.trim() && values.password
  const valid =
    values.alias.trim() && values.host.trim() && values.ip.trim() && values.user.trim() &&
    (mode === "create" ? values.password : true)

  return (
    <div className="flex flex-col gap-5">
      <div className="grid gap-4 sm:grid-cols-2">
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
            placeholder={mode === "edit" ? "留空表示不修改" : ""}
          />
        </div>
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="cluster-type">目标库类型</Label>
          <Select value={values.type} onValueChange={(value) => set("type", value)}>
            <SelectTrigger id="cluster-type">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="mysql">MySQL</SelectItem>
              <SelectItem value="postgres">PostgreSQL</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="cluster-group">所属分组</Label>
          <Select
            value={values.groupUUID || "none"}
            onValueChange={(value) => set("groupUUID", value === "none" ? "" : value)}
          >
            <SelectTrigger id="cluster-group">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="none">未分组</SelectItem>
              {(groups.data?.groups ?? []).map((group) => (
                <SelectItem key={group.UUID} value={group.UUID}>
                  {group.Name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="cluster-status">状态</Label>
          <Select
            value={String(values.status)}
            onValueChange={(value) => set("status", Number(value))}
          >
            <SelectTrigger id="cluster-status">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1">正常</SelectItem>
              <SelectItem value="2">禁用</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="flex flex-wrap items-center gap-3 rounded-md border border-border bg-surface-secondary px-4 py-3">
        <Button
          variant="outline"
          disabled={!connectable || testing}
          onClick={testConnection}
        >
          <IconPlugConnected size={16} stroke={1.5} />
          {testing ? "测试中…" : "测试连接"}
        </Button>
        {testResult ? (
          <div className="flex min-w-0 flex-col gap-0.5 text-sm">
            <span className={testResult.Success ? "text-success" : "text-danger"}>
              {testResult.Success ? "连接成功" : "连接失败"}
              {testResult.Success ? ` · ${testResult.LatencyMs} ms` : ""}
            </span>
            <span className="truncate text-xs text-muted-foreground">
              {testResult.Success
                ? `${testResult.Version} · 可见用户库 ${testResult.Databases.length} 个`
                : testResult.Message}
            </span>
          </div>
        ) : (
          <span className="text-sm text-muted-foreground">
            测试连接不会保存群集，只校验地址、账户与密码。
          </span>
        )}
      </div>

      {error ? (
        <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
      ) : null}

      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={onCancel}>
          取消
        </Button>
        <Button disabled={!valid || saving} onClick={submit}>
          <IconDeviceFloppy size={16} stroke={1.5} />
          {saving ? "保存中…" : "保存"}
        </Button>
      </div>
    </div>
  )
}
