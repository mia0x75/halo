import { useCallback, useMemo, useState } from "react"
import { IconClock, IconPlus, IconTrash } from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { CursorPagination } from "@/components/ui/cursor-pagination"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { useGraphQL } from "@/hooks/use-graphql"
import { formatUnix } from "@/lib/api/mappers"
import { graphqlRequest } from "@/lib/api/client"
import {
  CRONS_QUERY,
  PENDING_TASKS_QUERY,
  REMOVE_CRON_MUTATION,
  TASK_EXECUTIONS_QUERY,
  type CronNode,
  type TaskExecutionNode,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

type Tab = "crons" | "tasks" | "executions"

export interface CronsPageProps {
  /** 跳到新建作业页 */
  onCreate: () => void
}

export function CronsPage({ onCreate }: CronsPageProps) {
  const [tab, setTab] = useState<Tab>("crons")
  const [removing, setRemoving] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  const [after, setAfter] = useState<string | null>(null)
  const [stack, setStack] = useState<(string | null)[]>([])

  const crons = useGraphQL<
    {
      crons: {
        totalCount: number
        pageInfo: { hasNextPage: boolean; hasPreviousPage: boolean; startCursor: string; endCursor: string }
        edges: { cursor: string; node: CronNode }[]
      } | null
    },
    { first: number; after?: string }
  >(CRONS_QUERY, { first: 10, after: after ?? undefined }, { crons: null })

  const tasks = useGraphQL<{ tasks: CronNode[] }>(PENDING_TASKS_QUERY, undefined, { tasks: [] })

  const executions = useGraphQL<
    {
      taskExecutions: TaskExecutionNode[]
      taskExecutionSummary: { Status: number; Count: number }[]
    },
    { first: number }
  >(TASK_EXECUTIONS_QUERY, { first: 50 }, { taskExecutions: [], taskExecutionSummary: [] })

  const cronList = crons.data?.crons?.edges.map((edge) => edge.node) ?? []
  const connection = crons.data?.crons ?? null

  const remove = async () => {
    if (!removing) return
    setError(null)
    try {
      await graphqlRequest(REMOVE_CRON_MUTATION, { id: removing }, getToken())
      setRemoving(null)
      crons.refetch()
      tasks.refetch()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
      setRemoving(null)
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

  const summary = useMemo(
    () =>
      (executions.data?.taskExecutionSummary ?? []).map((item) => ({
        status: item.Status,
        count: item.Count,
      })),
    [executions.data],
  )

  return (
    <>
      <PageHeader
        title="定时任务"
        description="系统作业、工单预约与执行记录"
        action={
          <Button onClick={onCreate}>
            <IconPlus size={16} stroke={1.5} />
            新建作业
          </Button>
        }
      />

      {error ? (
        <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
      ) : null}

      <SectionCard title="作业与执行" description="系统作业按固定间隔运行，命令由服务端固定" bodyClassName="p-0">
        <div className="flex gap-4 border-b border-border px-5 pt-3">
          {(
            [
              ["crons", "计划任务", cronList.length],
              ["tasks", "等待执行", tasks.data?.tasks?.length ?? 0],
              ["executions", "执行记录", executions.data?.taskExecutions?.length ?? 0],
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

        {tab === "crons" ? (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>名称</TableHead>
                <TableHead>命令</TableHead>
                <TableHead>间隔</TableHead>
                <TableHead>下次执行</TableHead>
                <TableHead>状态</TableHead>
                <TableHead className="w-12"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {cronList.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="py-8 text-center text-sm text-muted-foreground">
                    还没有计划任务，点右上角新建一个系统作业。
                  </TableCell>
                </TableRow>
              ) : null}
              {cronList.map((cron) => (
                <TableRow key={cron.UUID}>
                  <TableCell className="font-medium">{cron.Name}</TableCell>
                  <TableCell className="font-mono text-xs text-muted-foreground">
                    {[cron.Cmd, cron.Params].filter(Boolean).join(" ")}
                  </TableCell>
                  <TableCell className="text-muted-foreground">{cron.Interval}</TableCell>
                  <TableCell className="text-muted-foreground">
                    {cron.NextRun ? formatUnix(Date.parse(cron.NextRun) / 1000) : "-"}
                  </TableCell>
                  <TableCell>
                    <Badge variant={cron.Status === "PENDING" ? "info" : "muted"}>
                      {cron.Status}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="ghost"
                      size="icon-sm"
                      aria-label={`删除 ${cron.Name}`}
                      onClick={() => setRemoving(cron.UUID)}
                    >
                      <IconTrash size={16} stroke={1.5} />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        ) : null}

        {tab === "tasks" ? (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>名称</TableHead>
                <TableHead>命令</TableHead>
                <TableHead>间隔</TableHead>
                <TableHead>下次执行</TableHead>
                <TableHead>状态</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {(tasks.data?.tasks ?? []).length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="py-8 text-center text-sm text-muted-foreground">
                    当前没有等待执行的任务。
                  </TableCell>
                </TableRow>
              ) : null}
              {(tasks.data?.tasks ?? []).map((task) => (
                <TableRow key={task.UUID}>
                  <TableCell className="font-medium">{task.Name}</TableCell>
                  <TableCell className="font-mono text-xs text-muted-foreground">
                    {[task.Cmd, task.Params].filter(Boolean).join(" ")}
                  </TableCell>
                  <TableCell className="text-muted-foreground">{task.Interval}</TableCell>
                  <TableCell className="text-muted-foreground">
                    {task.NextRun ? formatUnix(Date.parse(task.NextRun) / 1000) : "-"}
                  </TableCell>
                  <TableCell>
                    <Badge variant="info">
                      <IconClock size={12} stroke={1.5} className="mr-1" />
                      {task.Status}
                    </Badge>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        ) : null}

        {tab === "executions" ? (
          <>
            {summary.length > 0 ? (
              <div className="flex flex-wrap items-center gap-2 border-b border-border px-5 py-3">
                {summary.map((item) => (
                  <Badge key={item.status} variant="muted">
                    {executionStatusLabel(item.status)} {item.count}
                  </Badge>
                ))}
              </div>
            ) : null}
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>名称</TableHead>
                  <TableHead>类型</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>尝试</TableHead>
                  <TableHead>计划时间</TableHead>
                  <TableHead>耗时</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(executions.data?.taskExecutions ?? []).length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="py-8 text-center text-sm text-muted-foreground">
                      还没有执行记录。
                    </TableCell>
                  </TableRow>
                ) : null}
                {(executions.data?.taskExecutions ?? []).map((item) => (
                  <TableRow key={item.UUID}>
                    <TableCell className="font-medium">{item.Name}</TableCell>
                    <TableCell className="text-muted-foreground">{item.Type}</TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          item.Status === 1 ? "success" : item.Status === 2 ? "info" : "danger"
                        }
                      >
                        {executionStatusLabel(item.Status)}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-muted-foreground">{item.Attempt}</TableCell>
                    <TableCell className="text-muted-foreground">
                      {formatUnix(item.ScheduledAt)}
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {item.FinishedAt && item.StartedAt
                        ? `${item.FinishedAt - item.StartedAt}s`
                        : "-"}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </>
        ) : null}

        {tab === "crons" ? (
          <CursorPagination
            total={connection?.totalCount ?? 0}
            hasPreviousPage={connection?.pageInfo.hasPreviousPage ?? false}
            hasNextPage={connection?.pageInfo.hasNextPage ?? false}
            onPrevious={previousPage}
            onNext={nextPage}
          />
        ) : null}
      </SectionCard>

      <AlertDialog open={removing !== null} onOpenChange={(open) => !open && setRemoving(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="text-base font-medium text-strong">删除作业</AlertDialogTitle>
            <AlertDialogDescription className="mt-1 text-sm text-muted-foreground">
              删除后该作业不再调度，确认删除？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <div className="mt-5 flex justify-end gap-2">
            <AlertDialogCancel asChild>
              <Button variant="outline">取消</Button>
            </AlertDialogCancel>
            <AlertDialogAction asChild>
              <Button variant="danger" onClick={remove}>
                删除
              </Button>
            </AlertDialogAction>
          </div>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}

/** 任务执行状态：0 失败 1 成功 2 运行中。 */
function executionStatusLabel(status: number): string {
  if (status === 1) return "成功"
  if (status === 2) return "运行中"
  return "失败"
}
