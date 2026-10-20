import { useMemo, useState } from "react"
import { IconCheck, IconPencil, IconSettings, IconX } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
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

/** 系统选项页：平台级配置项，值可改。 */
export function OptionsPage() {
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
    <>
      <PageHeader
        title="系统选项"
        description="平台级配置项，值保存在元数据库里"
        action={
          <Badge variant="default" className="px-2 py-1">
            <IconSettings size={14} stroke={1.5} className="mr-1" />
            {options.length} 项
          </Badge>
        }
      />

      {error ? (
        <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
      ) : null}

      <SectionCard title="配置项" description="只读项由服务端固定，可写项能在这里改值" bodyClassName="p-0">
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
                      <Input
                        value={draft}
                        onChange={(event) => setDraft(event.target.value)}
                        className="h-8 max-w-64"
                        aria-label={`${option.Name} 的值`}
                      />
                      <Button
                        size="icon-sm"
                        aria-label="保存"
                        disabled={saving}
                        onClick={() => save(option)}
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
      </SectionCard>
    </>
  )
}
