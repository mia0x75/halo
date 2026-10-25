import { useEffect, useMemo, useState } from "react"
import {
  IconCheck,
  IconExternalLink,
  IconInbox,
  IconMail,
  IconMailOpened,
  IconSearch,
  IconTrash,
} from "@tabler/icons-react"

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
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Tabs } from "@/components/ui/tabs"
import { useGraphQL } from "@/hooks/use-graphql"
import { graphqlRequest } from "@/lib/api/client"
import { formatUnix } from "@/lib/api/mappers"
import { notifyDataChanged, subscribeDataChanged } from "@/lib/api/data-events"
import {
  MARK_ALL_MESSAGES_READ_MUTATION,
  MARK_MESSAGE_READ_MUTATION,
  MESSAGES_QUERY,
  REMOVE_MESSAGE_MUTATION,
  UNREAD_MESSAGE_COUNT_QUERY,
  type MessageConnectionNode,
  type MessageNode,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"
import { PAGE_SIZE } from "@/lib/cursor"

interface CategoryMeta {
  label: string
  variant: "info" | "warning" | "muted"
}

const categories: Record<string, CategoryMeta> = {
  ticket: { label: "工单", variant: "info" },
  review: { label: "审核", variant: "warning" },
  system: { label: "系统", variant: "muted" },
}

function categoryMeta(category: string): CategoryMeta {
  return categories[category] ?? { label: category || "通知", variant: "muted" }
}

function ticketUUID(link: string): string | null {
  const matched = /^\/tickets\/([^/?#]+)/.exec(link)
  return matched?.[1] ?? null
}

export interface MessagesPageProps {
  onOpenTicket: (uuid: string) => void
}

/**
 * 我的消息：过滤与分页都在服务端做，列表按未读/分类/关键词筛选，
 * 右侧阅读正文并跳转关联工单。
 */
export function MessagesPage({ onOpenTicket }: MessagesPageProps) {
  const [tab, setTab] = useState("all")
  const [category, setCategory] = useState("all")
  const [keyword, setKeyword] = useState("")
  const [search, setSearch] = useState("")
  const [after, setAfter] = useState<string | null>(null)
  const [stack, setStack] = useState<(string | null)[]>([])
  const [readOverrides, setReadOverrides] = useState<Record<string, number>>({})
  const [selectedUUID, setSelectedUUID] = useState<string | null>(null)
  const [actionError, setActionError] = useState<string | null>(null)
  const [removing, setRemoving] = useState<MessageNode | null>(null)

  // 搜索输入防抖，避免每敲一个字都打一次接口。
  useEffect(() => {
    const timer = setTimeout(() => {
      setSearch(keyword.trim())
      setAfter(null)
      setStack([])
    }, 300)
    return () => clearTimeout(timer)
  }, [keyword])

  const remote = useGraphQL<
    { messages: MessageConnectionNode },
    { unreadOnly: boolean; category: string; keyword: string; first: number; after?: string }
  >(
    MESSAGES_QUERY,
    {
      unreadOnly: tab === "unread",
      category: category === "all" ? "" : category,
      keyword: search,
      first: PAGE_SIZE,
      after: after ?? undefined,
    },
    { messages: { totalCount: 0, pageInfo: { hasNextPage: false, hasPreviousPage: false, startCursor: "", endCursor: "" }, edges: [] } },
  )

  const messages = useMemo(
    () =>
      (remote.data?.messages.edges ?? []).map((edge) => {
        const message = edge.node
        return readOverrides[message.UUID]
          ? { ...message, ReadAt: readOverrides[message.UUID] }
          : message
      }),
    [readOverrides, remote.data],
  )

  const pageInfo = remote.data?.messages.pageInfo

  // 未读数取全局计数，不受当前筛选影响。
  const unread = useGraphQL<{ unreadMessageCount: number }>(UNREAD_MESSAGE_COUNT_QUERY)
  const unreadCount = unread.data?.unreadMessageCount ?? 0

  // 顶栏弹窗与消息页都可能改动已读状态，页面跟着主题事件一起重取，
  // 否则「未读」标签上的数字会一直停在挂载时的值。
  const refetchMessages = remote.refetch
  const refetchUnread = unread.refetch
  useEffect(
    () =>
      subscribeDataChanged("messages", () => {
        refetchMessages()
        refetchUnread()
      }),
    [refetchMessages, refetchUnread],
  )

  const availableCategories = useMemo(
    () => Array.from(new Set(messages.map((message) => message.Category))).filter(Boolean),
    [messages],
  )

  const selected = messages.find((message) => message.UUID === selectedUUID) ?? null

  const markRead = async (message: MessageNode) => {
    if (message.ReadAt !== null) return
    try {
      await graphqlRequest<{ markMessageRead: boolean }, { uuid: string }>(
        MARK_MESSAGE_READ_MUTATION,
        { uuid: message.UUID },
        getToken(),
      )
      setReadOverrides((previous) => ({
        ...previous,
        [message.UUID]: Math.floor(Date.now() / 1000),
      }))
      setActionError(null)
      notifyDataChanged("messages")
    } catch (cause) {
      setActionError(cause instanceof Error ? cause.message : String(cause))
    }
  }

  const markAllRead = async () => {
    try {
      await graphqlRequest<{ markAllMessagesRead: number }>(
        MARK_ALL_MESSAGES_READ_MUTATION,
        {},
        getToken(),
      )
      setReadOverrides((previous) => {
        const next = { ...previous }
        const now = Math.floor(Date.now() / 1000)
        for (const message of messages) {
          if (message.ReadAt === null) next[message.UUID] = now
        }
        return next
      })
      setActionError(null)
      notifyDataChanged("messages")
    } catch (cause) {
      setActionError(cause instanceof Error ? cause.message : String(cause))
    }
  }

  const removeMessage = async (message: MessageNode) => {
    try {
      await graphqlRequest<{ removeMessage: boolean }, { uuid: string }>(
        REMOVE_MESSAGE_MUTATION,
        { uuid: message.UUID },
        getToken(),
      )
      setRemoving(null)
      if (selectedUUID === message.UUID) setSelectedUUID(null)
      notifyDataChanged("messages")
    } catch (cause) {
      setActionError(cause instanceof Error ? cause.message : String(cause))
    }
  }

  const openMessage = (message: MessageNode) => {
    setSelectedUUID(message.UUID)
    void markRead(message)
  }

  const linkedTicket = selected ? ticketUUID(selected.Link) : null

  const resetPage = () => {
    setAfter(null)
    setStack([])
  }

  return (
    <>
      <PageHeader
        title="我的消息"
        description="审核、工单状态变更等站内通知"
        action={
          <Button
            variant="outline"
            onClick={() => void markAllRead()}
            disabled={unreadCount === 0}
          >
            <IconCheck size={16} stroke={1.5} />
            全部标记为已读
          </Button>
        }
      />

      {actionError ? (
        <p role="alert" className="text-sm text-danger">
          {actionError}
        </p>
      ) : null}

      <SectionCard bodyClassName="p-0">
        <div className="flex flex-wrap items-center gap-3 border-b border-border px-4 py-3">
          <Tabs
            className="border-b-0"
            value={tab}
            onChange={(key) => {
              setTab(key)
              resetPage()
            }}
            items={[
              { key: "all", label: "全部" },
              { key: "unread", label: "未读", badge: unreadCount },
            ]}
          />
          <div className="relative ml-auto w-full max-w-xs">
            <IconSearch
              size={16}
              stroke={1.5}
              className="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-muted-foreground"
            />
            <Input
              value={keyword}
              onChange={(event) => setKeyword(event.target.value)}
              placeholder="搜索标题或内容"
              aria-label="搜索消息"
              className="pl-9"
            />
          </div>
          <Select
            value={category}
            onValueChange={(value) => {
              setCategory(value)
              resetPage()
            }}
          >
            <SelectTrigger aria-label="按分类筛选" className="w-32">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">全部分类</SelectItem>
              {availableCategories.map((item) => (
                <SelectItem key={item} value={item}>
                  {categoryMeta(item).label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="grid min-h-[28rem] md:grid-cols-[22rem_minmax(0,1fr)]">
          <div className="border-b border-border md:border-r md:border-b-0">
            {remote.loading && messages.length === 0 ? (
              <p className="px-4 py-6 text-sm text-muted-foreground">加载中…</p>
            ) : messages.length === 0 ? (
              <div className="flex flex-col items-center gap-2 px-4 py-12 text-muted-foreground">
                <IconInbox size={28} stroke={1.5} />
                <p className="text-sm">没有符合条件的消息</p>
              </div>
            ) : (
              <ul className="max-h-[36rem] overflow-y-auto">
                {messages.map((message) => {
                  const meta = categoryMeta(message.Category)
                  const unread = message.ReadAt === null
                  return (
                    <li key={message.UUID} className="flex items-start gap-1 pr-2">
                      <button
                        type="button"
                        onClick={() => openMessage(message)}
                        aria-current={message.UUID === selectedUUID ? "true" : undefined}
                        className={`flex min-w-0 flex-1 flex-col gap-1 border-b border-border px-4 py-3 text-left transition-colors hover:bg-surface-secondary ${
                          message.UUID === selectedUUID ? "bg-surface-secondary" : ""
                        }`}
                      >
                        <span className="flex items-center gap-2">
                          {unread ? (
                            <span
                              className="size-2 shrink-0 rounded-full bg-primary"
                              aria-label="未读"
                            />
                          ) : (
                            <IconMailOpened
                              size={14}
                              stroke={1.5}
                              className="shrink-0 text-muted-foreground"
                            />
                          )}
                          <span
                            className={`truncate text-sm ${unread ? "font-medium text-strong" : "text-foreground"}`}
                          >
                            {message.Title}
                          </span>
                          <Badge variant={meta.variant} className="ml-auto shrink-0">
                            {meta.label}
                          </Badge>
                        </span>
                        <span className="line-clamp-2 text-sm text-muted-foreground">
                          {message.Content}
                        </span>
                        <span className="text-xs text-muted-foreground">
                          {formatUnix(message.CreateAt)}
                        </span>
                      </button>
                      <Button
                        variant="ghost"
                        size="icon-sm"
                        aria-label={`删除 ${message.Title}`}
                        onClick={() => setRemoving(message)}
                      >
                        <IconTrash size={16} stroke={1.5} />
                      </Button>
                    </li>
                  )
                })}
              </ul>
            )}
            <CursorPagination
              divided={false}
              page={stack.length + 1}
              empty={messages.length === 0}
              hasPreviousPage={pageInfo?.hasPreviousPage ?? false}
              hasNextPage={pageInfo?.hasNextPage ?? false}
              onPrevious={() => {
                const previousCursor = stack[stack.length - 1] ?? null
                setStack((items) => items.slice(0, -1))
                setAfter(previousCursor)
              }}
              onNext={() => {
                setStack((items) => [...items, after])
                setAfter(pageInfo?.endCursor || null)
              }}
              onReset={() => {
                setAfter(null)
                setStack([])
              }}
            />
          </div>

          <div className="flex flex-col">
            {selected ? (
              <article className="flex flex-col gap-4 p-6">
                <div className="flex flex-wrap items-center gap-2">
                  <Badge variant={categoryMeta(selected.Category).variant}>
                    {categoryMeta(selected.Category).label}
                  </Badge>
                  <span className="text-sm text-muted-foreground">
                    {formatUnix(selected.CreateAt)}
                  </span>
                  {selected.ReadAt === null ? <Badge variant="default">未读</Badge> : null}
                </div>
                <h2 className="text-lg font-semibold text-strong">{selected.Title}</h2>
                <p className="whitespace-pre-wrap text-sm leading-6 text-foreground">
                  {selected.Content}
                </p>
                {linkedTicket ? (
                  <Button
                    variant="outline"
                    className="self-start"
                    onClick={() => onOpenTicket(linkedTicket)}
                  >
                    <IconExternalLink size={16} stroke={1.5} />
                    打开关联工单
                  </Button>
                ) : null}
              </article>
            ) : (
              <div className="flex flex-1 flex-col items-center justify-center gap-2 p-6 text-muted-foreground">
                <IconMail size={32} stroke={1.5} />
                <p className="text-sm">选择左侧消息查看内容</p>
              </div>
            )}
          </div>
        </div>
      </SectionCard>

      <AlertDialog open={removing !== null} onOpenChange={(open) => !open && setRemoving(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="text-base font-medium text-strong">
              删除消息
            </AlertDialogTitle>
            <AlertDialogDescription className="mt-1 text-sm text-muted-foreground">
              删除后这条消息不再出现在收件箱。确认删除「{removing?.Title}」？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <div className="mt-5 flex justify-end gap-2">
            <AlertDialogCancel asChild>
              <Button variant="outline">取消</Button>
            </AlertDialogCancel>
            <AlertDialogAction asChild>
              <Button
                variant="danger"
                onClick={() => {
                  if (removing) void removeMessage(removing)
                }}
              >
                删除
              </Button>
            </AlertDialogAction>
          </div>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}
