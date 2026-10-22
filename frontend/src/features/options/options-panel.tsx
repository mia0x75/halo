import { useMemo, useState } from "react"
import { IconCheck, IconPencil, IconX } from "@tabler/icons-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
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
import { OPTIONS_QUERY, PATCH_OPTION_MUTATION, type OptionNode } from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

export interface OptionsPanelProps {
  /** 表格标题为空时只渲染表格，供嵌套页签使用。 */
  showHeader?: boolean
}

/** 平台参数表格：只读项由服务端固定，可写项能在这里改值。 */
export function OptionsPanel({ showHeader = true }: OptionsPanelProps) {
  const remote = useGraphQL<{ options: OptionNode[] }>(OPTIONS_QUERY, undefined, { options: [] })
  const [editing, setEditing] = useState<string | null>(null)
  const [draft, setDraft] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  const options = useMemo(() => remote.data?.options ?? [], [remote.data])

  const startEdit = (option: OptionNode) => {
    setEditing(option.UUID)
    setDraft(option.Parameter)
    setError(null)
  }

  const cancel = () => {
    setEditing(null)
    setDraft("")
  }

  const save = async (option: OptionNode) => {
    setSaving(true)
    setError(null)
    try {
      await graphqlRequest(
        PATCH_OPTION_MUTATION,
        { input: { OptionUUID: option.UUID, Parameter: draft } },
        getToken(),
      )
      cancel()
      remote.refetch()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="flex flex-col">
      {error ? (
        <p className="mb-3 rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
      ) : null}
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>配置项</TableHead>
            <TableHead>值</TableHead>
            <TableHead>说明</TableHead>
            <TableHead>可写</TableHead>
            <TableHead className="w-24"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {options.length === 0 ? (
            <TableRow>
              <TableCell colSpan={5} className="py-8 text-center text-sm text-muted-foreground">
                没有可读的系统选项。
              </TableCell>
            </TableRow>
          ) : null}
          {options.map((option) => (
            <TableRow key={option.UUID}>
              <TableCell className="font-mono text-xs font-medium">{option.Name}</TableCell>
              <TableCell>
                {editing === option.UUID ? (
                  <div className="flex items-center gap-2">
                    {option.Element === "switch" ? (
                      <Select value={draft} onValueChange={(value) => setDraft(value)}>
                        <SelectTrigger
                          aria-label={`${option.Name} 的值`}
                          className="h-8 max-w-32"
                        >
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="true">开启</SelectItem>
                          <SelectItem value="false">关闭</SelectItem>
                        </SelectContent>
                      </Select>
                    ) : (
                      <Input
                        value={draft}
                        onChange={(event) => setDraft(event.target.value)}
                        type={
                          option.Element === "password"
                            ? "password"
                            : option.Element === "number"
                              ? "number"
                              : "text"
                        }
                        className="h-8 max-w-64"
                        aria-label={`${option.Name} 的值`}
                      />
                    )}
                    <Button
                      size="icon-sm"
                      aria-label="保存"
                      disabled={saving}
                      onClick={() => void save(option)}
                    >
                      <IconCheck size={16} stroke={1.5} />
                    </Button>
                    <Button size="icon-sm" aria-label="取消" onClick={cancel}>
                      <IconX size={16} stroke={1.5} />
                    </Button>
                  </div>
                ) : (
                  <span className="font-mono text-xs">{option.Parameter}</span>
                )}
              </TableCell>
              <TableCell className="text-muted-foreground">{option.Description || "-"}</TableCell>
              <TableCell>
                <Badge variant={option.Writable === 1 ? "success" : "muted"}>
                  {option.Writable === 1 ? "可写" : "只读"}
                </Badge>
              </TableCell>
              <TableCell>
                {option.Writable === 1 && editing !== option.UUID ? (
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    aria-label={`编辑 ${option.Name}`}
                    onClick={() => startEdit(option)}
                  >
                    <IconPencil size={16} stroke={1.5} />
                  </Button>
                ) : null}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      {showHeader ? null : null}
    </div>
  )
}
