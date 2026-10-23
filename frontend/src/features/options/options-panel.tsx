import { useMemo, useState } from "react"
import { IconKey, IconLock, IconPencil } from "@tabler/icons-react"

import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import { useGraphQL } from "@/hooks/use-graphql"
import { graphqlRequest } from "@/lib/api/client"
import { OPTIONS_QUERY, PATCH_OPTION_MUTATION, type OptionNode } from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

/** 选项按分组聚合成卡片，组内按展示顺序排列。 */
function groupOptions(options: OptionNode[]): { name: string; items: OptionNode[] }[] {
  const groups = new Map<string, OptionNode[]>()
  for (const option of options) {
    const name = option.GroupName || "其它"
    const bucket = groups.get(name)
    if (bucket) {
      bucket.push(option)
    } else {
      groups.set(name, [option])
    }
  }
  // 接口已按 sort_order 排好序，分组保持首次出现的顺序，组内也按 sort_order 排。
  return [...groups.entries()].map(([name, items]) => ({
    name,
    items: items.sort((a, b) => a.SortOrder - b.SortOrder || a.Name.localeCompare(b.Name)),
  }))
}

/** 平台参数：每个分组一张卡片，只读项灰显，可写项就地编辑。 */
export function OptionsPanel() {
  const remote = useGraphQL<{ options: OptionNode[] }>(OPTIONS_QUERY, undefined, { options: [] })
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState<string | null>(null)
  const [secret, setSecret] = useState<OptionNode | null>(null)
  const [secretValue, setSecretValue] = useState("")

  const options = useMemo(() => remote.data?.options ?? [], [remote.data])
  const groups = useMemo(() => groupOptions(options), [options])

  // 依赖项没打开时，依赖它的选项不可编辑。
  const enabled = useMemo(() => {
    const map = new Map<string, boolean>()
    for (const option of options) {
      map.set(option.Name, option.Parameter === "true")
    }
    return map
  }, [options])

  const patch = async (option: OptionNode, value: string) => {
    setSaving(option.UUID)
    setError(null)
    try {
      await graphqlRequest(
        PATCH_OPTION_MUTATION,
        { input: { OptionUUID: option.UUID, Parameter: value } },
        getToken(),
      )
      remote.refetch()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setSaving(null)
    }
  }

  const saveSecret = async () => {
    if (!secret) return
    const target = secret
    setSecret(null)
    setSecretValue("")
    if (secretValue === "") return
    await patch(target, secretValue)
  }

  return (
    <div className="flex flex-col gap-4">
      {error ? (
        <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
      ) : null}
      <div className="grid gap-4 md:grid-cols-2">
        {groups.map((group) => (
          <OptionGroup
            key={group.name}
            name={group.name}
            items={group.items}
            enabled={enabled}
            saving={saving}
            onPatch={patch}
            onEditSecret={(option) => {
              setSecret(option)
              setSecretValue("")
            }}
          />
        ))}
      </div>
      <Dialog
        open={secret !== null}
        onOpenChange={(open) => {
          if (!open) {
            setSecret(null)
            setSecretValue("")
          }
        }}
      >
        <DialogContent>
          <DialogTitle className="text-base font-medium text-strong">
            修改 {secret?.Name}
          </DialogTitle>
          <DialogDescription className="mt-1 text-sm text-muted-foreground">
            密码加密后存储，保存后不再回显。
          </DialogDescription>
          <div className="mt-4 flex flex-col gap-1.5">
            <Label htmlFor="option-secret">新值</Label>
            <Input
              id="option-secret"
              type="password"
              value={secretValue}
              onChange={(event) => setSecretValue(event.target.value)}
              autoComplete="new-password"
            />
          </div>
          <div className="mt-5 flex justify-end gap-2">
            <Button variant="outline" onClick={() => setSecret(null)}>
              取消
            </Button>
            <Button disabled={secretValue === ""} onClick={() => void saveSecret()}>
              保存
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}

interface OptionGroupProps {
  name: string
  items: OptionNode[]
  enabled: Map<string, boolean>
  saving: string | null
  onPatch: (option: OptionNode, value: string) => Promise<void>
  onEditSecret: (option: OptionNode) => void
}

/** 一个分组的卡片。 */
function OptionGroup({
  name,
  items,
  enabled,
  saving,
  onPatch,
  onEditSecret,
}: OptionGroupProps) {
  return (
    <SectionCard title={name} bodyClassName="p-0">
      <ul className="divide-y divide-border">
        {items.map((option) => {
          const writable = option.Writable === 1
          const blocked = option.DependsOn !== "" && enabled.get(option.DependsOn) !== true
          const disabled = !writable || blocked
          return (
            <li
              key={option.UUID}
              className="flex flex-wrap items-center gap-3 px-5 py-3"
            >
              <div className="flex min-w-0 flex-1 flex-col gap-0.5">
                <div className="flex items-center gap-2">
                  <span className="font-mono text-xs font-medium">{option.Name}</span>
                  {!writable ? (
                    <Badge variant="muted">
                      <IconLock size={12} stroke={1.5} />
                      只读
                    </Badge>
                  ) : null}
                </div>
                <span className="truncate text-xs text-muted-foreground">
                  {option.Description || "-"}
                </span>
              </div>
              <OptionControl
                option={option}
                disabled={disabled}
                saving={saving === option.UUID}
                onPatch={onPatch}
                onEditSecret={onEditSecret}
              />
            </li>
          )
        })}
      </ul>
    </SectionCard>
  )
}

interface OptionControlProps {
  option: OptionNode
  disabled: boolean
  saving: boolean
  onPatch: (option: OptionNode, value: string) => Promise<void>
  onEditSecret?: (option: OptionNode) => void
}

/** 按 element 渲染对应控件。 */
function OptionControl({
  option,
  disabled,
  saving,
  onPatch,
  onEditSecret,
}: OptionControlProps) {
  if (option.Writable !== 1) {
    // 只读项直接展示取值。
    return <span className="font-mono text-xs text-muted-foreground">{option.Parameter}</span>
  }

  switch (option.Element) {
    case "switch":
      return (
        <Switch
          checked={option.Parameter === "true"}
          disabled={disabled || saving}
          aria-label={option.Name}
          onCheckedChange={(checked) => void onPatch(option, checked ? "true" : "false")}
        />
      )
    case "password":
      return (
        <Button
          variant="outline"
          size="sm"
          disabled={disabled || saving}
          onClick={() => onEditSecret?.(option)}
        >
          <IconKey size={14} stroke={1.5} />
          修改
        </Button>
      )
    case "number":
      return (
        <NumberControl option={option} disabled={disabled} saving={saving} onPatch={onPatch} />
      )
    default:
      return (
        <TextControl option={option} disabled={disabled} saving={saving} onPatch={onPatch} />
      )
  }
}

/** 文本输入：回车或失焦保存。 */
function TextControl({
  option,
  disabled,
  saving,
  onPatch,
}: OptionControlProps) {
  const [draft, setDraft] = useState(option.Parameter)
  const dirty = draft !== option.Parameter
  return (
    <div className="flex items-center gap-2">
      <Input
        value={draft}
        disabled={disabled || saving}
        onChange={(event) => setDraft(event.target.value)}
        onBlur={() => {
          if (dirty && !disabled) void onPatch(option, draft)
        }}
        onKeyDown={(event) => {
          if (event.key === "Enter" && dirty && !disabled) void onPatch(option, draft)
        }}
        className="h-8 w-56"
        aria-label={option.Name}
      />
      {dirty ? (
        <Button
          size="icon-sm"
          aria-label={`保存 ${option.Name}`}
          disabled={disabled || saving}
          onClick={() => void onPatch(option, draft)}
        >
          <IconPencil size={16} stroke={1.5} />
        </Button>
      ) : null}
    </div>
  )
}

/** 数值输入：按 min/max 限制，回车或失焦保存。 */
function NumberControl({
  option,
  disabled,
  saving,
  onPatch,
}: OptionControlProps) {
  const [draft, setDraft] = useState(option.Parameter)
  const dirty = draft !== option.Parameter
  const invalid =
    option.MinValue !== null && Number(draft) < option.MinValue
      ? true
      : option.MaxValue !== null && Number(draft) > option.MaxValue
        ? true
        : false
  return (
    <div className="flex items-center gap-2">
      <Input
        type="number"
        value={draft}
        min={option.MinValue ?? undefined}
        max={option.MaxValue ?? undefined}
        disabled={disabled || saving}
        onChange={(event) => setDraft(event.target.value)}
        onBlur={() => {
          if (dirty && !disabled && !invalid) void onPatch(option, draft)
        }}
        onKeyDown={(event) => {
          if (event.key === "Enter" && dirty && !disabled && !invalid) void onPatch(option, draft)
        }}
        className="h-8 w-32"
        aria-label={option.Name}
      />
      {invalid ? (
        <span className="text-xs text-danger">
          取值范围 {option.MinValue ?? "-"} ~ {option.MaxValue ?? "-"}
        </span>
      ) : null}
    </div>
  )
}
