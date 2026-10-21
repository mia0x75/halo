import { useMemo, useState } from "react"
import {
  IconCheck,
  IconExternalLink,
  IconInbox,
  IconMail,
  IconMailOpened,
  IconSearch,
} from "@tabler/icons-react"

import { PageHeader } from "@/components/layout/page-header"
import { SectionCard } from "@/components/shared/section-card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Select } from "@/components/ui/select"
import { Tabs } from "@/components/ui/tabs"
import { useGraphQL } from "@/hooks/use-graphql"
import { graphqlRequest } from "@/lib/api/client"
import { formatUnix } from "@/lib/api/mappers"
import { notifyMessagesChanged } from "@/lib/api/message-events"
import {
  MARK_ALL_MESSAGES_READ_MUTATION,
  MARK_MESSAGE_READ_MUTATION,
  MESSAGES_QUERY,
  type MessageNode,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

const inboxSize = 100

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

/** 我的消息：左边列表按未读/分类筛选，右边阅读正文并跳转关联工单。 */
export function MessagesPage({ onOpenTicket }: MessagesPageProps) {
  const remote = useGraphQL<{ messages: MessageNode[] }, { unreadOnly: boolean; first: number }>(
    MESSAGES_QUERY,
    { unreadOnly: false, first: inboxSize },
    { messages: [] },
  )
  const [readOverrides, setReadOverrides] = useState<Record<string, number>>({})
  const [tab, setTab] = useState("all")
  const [category, setCategory] = useState("all")
  const [keyword, setKeyword] = useState("")
  const [selectedUUID, setSelectedUUID] = useState<string | null>(null)
  const [actionError, setActionError] = useState<string | null>(null)

  const messages = useMemo(
    () =>
      (remote.data?.messages ?? []).map((message) =>
        readOverrides[message.UUID]
          ? { ...message, ReadAt: readOverrides[message.UUID] }
          : message,
      ),
    [readOverrides, remote.data],
  )
  const unreadCount = messages.filter((message) => message.ReadAt === null).length

  const availableCategories = useMemo(
    () => Array.from(new Set(messages.map((message) => message.Category))).filter(Boolean),
    [messages],
  )

  const filtered = useMemo(() => {
    const needle = keyword.trim().toLowerCase()
    return messages.filter((message) => {
      if (tab === "unread" && message.ReadAt !== null) return false
      if (category !== "all" && message.Category !== category) return false
      if (needle && !`${message.Title} ${message.Content}`.toLowerCase().includes(needle)) {
        return false
      }
      return true
    })
  }, [category, keyword, messages, tab])

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
      notifyMessagesChanged()
      remote.refetch()
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
      const now = Math.floor(Date.now() / 1000)
      setReadOverrides((previous) => {
        const next = { ...previous }
        for (const message of messages) {
          if (message.ReadAt === null) next[message.UUID] = now
        }
        return next
      })
      setActionError(null)
      notifyMessagesChanged()
      remote.refetch()
    } catch (cause) {
      setActionError(cause instanceof Error ? cause.message : String(cause))
    }
  }

  const openMessage = (message: MessageNode) => {
    setSelectedUUID(message.UUID)
    void markRead(message)
  }

  const linkedTicket = selected ? ticketUUID(selected.Link) : null

  return (
    <>
      <PageHeader
        title="我的消息"
        description="审核、工单状态变更等站内通知，最多展示最近 100 条。"
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

      <SectionCard title="收件箱" description={`未读 ${unreadCount} 条`} bodyClassName="p-0">
        <div className="flex flex-wrap items-center gap-3 border-b border-border px-4 py-3">
          <Tabs
            className="border-b-0"
            value={tab}
            onChange={setTab}
            items={[
              { key: "all", label: "全部", badge: messages.length },
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
            onChange={(event) => setCategory(event.target.value)}
            aria-label="按分类筛选"
            className="w-32"
          >
            <option value="all">全部分类</option>
            {availableCategories.map((item) => (
              <option key={item} value={item}>
                {categoryMeta(item).label}
              </option>
            ))}
          </Select>
        </div>

        <div className="grid min-h-[28rem] md:grid-cols-[22rem_minmax(0,1fr)]">
          <div className="border-b border-border md:border-r md:border-b-0">
            {remote.loading && messages.length === 0 ? (
              <p className="px-4 py-6 text-sm text-muted-foreground">加载中…</p>
            ) : filtered.length === 0 ? (
              <div className="flex flex-col items-center gap-2 px-4 py-12 text-muted-foreground">
                <IconInbox size={28} stroke={1.5} />
                <p className="text-sm">没有符合条件的消息</p>
              </div>
            ) : (
              <ul className="max-h-[36rem] overflow-y-auto">
                {filtered.map((message) => {
                  const meta = categoryMeta(message.Category)
                  const unread = message.ReadAt === null
                  return (
                    <li key={message.UUID}>
                      <button
                        type="button"
                        onClick={() => openMessage(message)}
                        aria-current={message.UUID === selectedUUID ? "true" : undefined}
                        className={`flex w-full flex-col gap-1 border-b border-border px-4 py-3 text-left transition-colors hover:bg-surface-secondary ${
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
                    </li>
                  )
                })}
              </ul>
            )}
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
                  {selected.ReadAt === null ? (
                    <Badge variant="default">未读</Badge>
                  ) : null}
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
    </>
  )
}
