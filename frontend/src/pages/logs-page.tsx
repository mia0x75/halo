import { useMemo, useState } from "react"
import { IconHistory, IconSearch } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { CursorPagination } from "@/components/ui/cursor-pagination"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import { useCursorPagination } from "@/hooks/use-cursor-pagination"
import { auditLogs } from "@/lib/mock/data"

const pageSize = 8

export function LogsPage() {
  const [keyword, setKeyword] = useState("")
  const [actor, setActor] = useState("all")

  const actors = useMemo(
    () => Array.from(new Set(auditLogs.map((log) => log.actor))),
    [],
  )

  const filtered = useMemo(() => {
    const needle = keyword.trim().toLowerCase()
    return auditLogs.filter((log) => {
      if (actor !== "all" && log.actor !== actor) return false
      if (!needle) return true
      return (
        log.operation.toLowerCase().includes(needle) ||
        log.target.toLowerCase().includes(needle) ||
        log.ip.includes(needle)
      )
    })
  }, [keyword, actor])

  const { page, reset, next, previous } = useCursorPagination(filtered, pageSize)

  return (
    <>
      <PageHeader title="操作日志" description="平台操作与审核动作的审计记录" />

      <SectionCard title="审计流水" description={`共 ${filtered.length} 条`} bodyClassName="p-0">
        <div className="flex flex-wrap items-end gap-3 border-b border-border px-5 py-4">
          <div className="flex flex-col gap-1.5 sm:w-72">
            <Label htmlFor="log-keyword">关键词</Label>
            <div className="relative">
              <IconSearch
                size={16}
                stroke={1.5}
                className="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-muted-foreground"
              />
              <Input
                id="log-keyword"
                value={keyword}
                onChange={(event) => {
                  setKeyword(event.target.value)
                  reset()
                }}
                placeholder="操作、对象或来源 IP"
                className="pl-9"
              />
            </div>
          </div>
          <div className="flex flex-col gap-1.5 sm:w-40">
            <Label htmlFor="log-actor">操作者</Label>
            <Select
              id="log-actor"
              value={actor}
              onChange={(event) => {
                setActor(event.target.value)
                reset()
              }}
            >
              <option value="all">全部操作者</option>
              {actors.map((name) => (
                <option key={name} value={name}>
                  {name}
                </option>
              ))}
            </Select>
          </div>
        </div>

        <ul className="divide-y divide-border">
          {page.edges.map((log) => (
            <li key={log.id} className="flex items-start gap-3 px-5 py-3">
              <span className="mt-0.5 text-muted-foreground">
                <IconHistory size={16} stroke={1.5} />
              </span>
              <div className="min-w-0 flex-1">
                <p className="text-sm">{log.operation}</p>
                <p className="mt-0.5 text-sm text-muted-foreground">
                  {log.actor} · {log.target} · {log.ip}
                </p>
              </div>
              <span className="shrink-0 text-sm text-muted-foreground">{log.createdAt}</span>
            </li>
          ))}
        </ul>

        <CursorPagination
          total={page.totalCount}
          hasPreviousPage={page.pageInfo.hasPreviousPage}
          hasNextPage={page.pageInfo.hasNextPage}
          startCursor={page.pageInfo.startCursor}
          endCursor={page.pageInfo.endCursor}
          onPrevious={previous}
          onNext={next}
        />
      </SectionCard>
    </>
  )
}
