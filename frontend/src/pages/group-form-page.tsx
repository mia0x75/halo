import { useMemo, useState } from "react"
import { IconArrowLeft } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
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
  CREATE_GROUP_MUTATION,
  GROUPS_QUERY,
  UPDATE_GROUP_MUTATION,
  type GroupNode,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

export interface GroupFormPageProps {
  mode: "create" | "edit"
  /** 编辑时带的分组 UUID */
  uuid?: string
  onBack: () => void
  onSaved: () => void
}

interface GroupFormProps {
  initial?: GroupNode
  mode: "create" | "edit"
  uuid?: string
  onSaved: () => void
}

/** 分组表单：初始值来自接口数据，保存后回到设置页。 */
function GroupForm({ initial, mode, uuid, onSaved }: GroupFormProps) {
  const [name, setName] = useState(initial?.Name ?? "")
  const [description, setDescription] = useState(initial?.Description ?? "")
  const [status, setStatus] = useState(String(initial?.Status ?? 1))
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  const submit = async () => {
    if (!name.trim()) {
      setError("分组名称不能为空")
      return
    }
    setSaving(true)
    setError(null)
    try {
      if (mode === "edit" && uuid) {
        await graphqlRequest(
          UPDATE_GROUP_MUTATION,
          {
            input: {
              GroupUUID: uuid,
              Name: name.trim(),
              Description: description.trim(),
              Status: Number(status),
            },
          },
          getToken(),
        )
      } else {
        await graphqlRequest(
          CREATE_GROUP_MUTATION,
          {
            input: {
              Name: name.trim(),
              Description: description.trim(),
              Status: Number(status),
            },
          },
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
          <Label htmlFor="group-name">分组名称</Label>
          <Input
            id="group-name"
            value={name}
            onChange={(event) => setName(event.target.value)}
            placeholder="核心交易"
          />
        </div>
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="group-status">状态</Label>
          <Select value={status} onValueChange={(value) => setStatus(value)}>
            <SelectTrigger id="group-status">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1">正常</SelectItem>
              <SelectItem value="2">停用</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="flex flex-col gap-1.5 sm:col-span-2">
          <Label htmlFor="group-description">分组说明</Label>
          <Input
            id="group-description"
            value={description}
            onChange={(event) => setDescription(event.target.value)}
            placeholder="订单、支付相关库"
          />
        </div>
      </div>

      {error ? (
        <p className="mt-4 rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
      ) : null}

      <div className="mt-5 flex justify-end gap-2">
        <Button onClick={() => void submit()} disabled={saving}>
          {mode === "create" ? "创建" : "保存"}
        </Button>
      </div>
    </>
  )
}

/** 群集分组新建/编辑页。表单是整页，不是弹窗。 */
export function GroupFormPage({ mode, uuid, onBack, onSaved }: GroupFormPageProps) {
  const remote = useGraphQL<{ groups: GroupNode[] }>(GROUPS_QUERY, undefined, { groups: [] })

  const initial = useMemo(() => {
    if (mode !== "edit" || !uuid) return undefined
    return remote.data?.groups.find((group) => group.UUID === uuid)
  }, [mode, remote.data, uuid])

  return (
    <>
      <PageHeader
        title={mode === "create" ? "新建分组" : "编辑分组"}
        description="分组用于把群集归类，并作为角色授权的边界"
        action={
          <Button variant="outline" onClick={onBack}>
            <IconArrowLeft size={16} stroke={1.5} />
            返回系统设置
          </Button>
        }
      />

      <SectionCard>
        {mode === "edit" && !initial ? (
          <p className="text-sm text-muted-foreground">正在读取分组信息…</p>
        ) : (
          <GroupForm
            key={initial?.UUID ?? mode}
            initial={initial}
            mode={mode}
            uuid={uuid}
            onSaved={onSaved}
          />
        )}
      </SectionCard>
    </>
  )
}
