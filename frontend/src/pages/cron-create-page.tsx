import { useState } from "react"
import { IconArrowLeft, IconPlus } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { InlineNote } from "@/components/shared/status-note"
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
import { graphqlRequest } from "@/lib/api/client"
import { CREATE_CRON_MUTATION } from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

export interface CronCreatePageProps {
  onBack: () => void
  onSaved: () => void
}

/** 系统作业类型。命令由服务端固定，页面只选类型与间隔。 */
const jobOptions = [
  { value: "METRICS_ROLLUP", label: "每日汇总工单与规则命中" },
  { value: "METRICS_CLEANUP", label: "清理过期指标与命中明细" },
  { value: "RETENTION_CLEANUP", label: "按保留期清理日志/站内信/任务记录" },
]

/** 新建系统作业页。 */
export function CronCreatePage({ onBack, onSaved }: CronCreatePageProps) {
  const [name, setName] = useState("")
  const [job, setJob] = useState("METRICS_ROLLUP")
  const [interval, setInterval] = useState("24h")
  const [days, setDays] = useState("90")
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  const submit = async () => {
    setSaving(true)
    setError(null)
    try {
      await graphqlRequest(
        CREATE_CRON_MUTATION,
        {
          input: {
            Name: name.trim(),
            Job: job,
            Interval: interval.trim(),
            ...(job === "METRICS_CLEANUP" && days ? { Days: Number(days) } : {}),
          },
        },
        getToken(),
      )
      onSaved()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setSaving(false)
    }
  }

  const valid = name.trim() && interval.trim()

  return (
    <>
      <PageHeader
        title="新建作业"
        description="系统作业的命令由服务端固定，这里只选类型与执行间隔"
        action={
          <Button variant="outline" onClick={onBack}>
            <IconArrowLeft size={16} stroke={1.5} />
            返回
          </Button>
        }
      />

      <SectionCard>
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="cron-name">作业名称</Label>
            <Input
              id="cron-name"
              value={name}
              onChange={(event) => setName(event.target.value)}
              placeholder="每日指标汇总"
            />
          </div>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="cron-job">作业类型</Label>
            <Select value={job} onValueChange={(value) => setJob(value)}>
              <SelectTrigger id="cron-job">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {jobOptions.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="cron-interval">执行间隔</Label>
            <Input
              id="cron-interval"
              value={interval}
              onChange={(event) => setInterval(event.target.value)}
              placeholder="24h"
            />
          </div>
          {job === "METRICS_CLEANUP" ? (
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="cron-days">保留天数</Label>
              <Input
                id="cron-days"
                type="number"
                value={days}
                onChange={(event) => setDays(event.target.value)}
              />
            </div>
          ) : null}
        </div>

        <div className="mt-5 flex flex-wrap items-center justify-end gap-2">
          {error ? <InlineNote className="mr-auto">{error}</InlineNote> : null}
          <Button variant="outline" onClick={onBack}>
            取消
          </Button>
          <Button disabled={!valid || saving} onClick={submit}>
            <IconPlus size={16} stroke={1.5} />
            {saving ? "创建中…" : "创建"}
          </Button>
        </div>
      </SectionCard>
    </>
  )
}
