import { useState } from "react"
import { IconBell, IconDatabase, IconMoon, IconSearch, IconSun } from "@tabler/icons-react"

import { ConnectionMenu } from "@/components/layout/connection-menu"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Input } from "@/components/ui/input"
import { graphqlRequest } from "@/lib/api/client"
import { useGraphQL } from "@/hooks/use-graphql"
import { useTheme } from "@/hooks/use-theme"
import { formatUnix } from "@/lib/api/mappers"
import type { MeUser } from "@/lib/api/queries"
import {
  MARK_MESSAGE_READ_MUTATION,
  MESSAGES_QUERY,
  UNREAD_MESSAGE_COUNT_QUERY,
  type MessageNode,
} from "@/lib/api/queries"
import { getToken } from "@/lib/api/session"

export interface TopNavProps {
  user: MeUser
  onSignOut: () => void
  /** 提交搜索：跳到工单列表并按关键词过滤 */
  onSearch: (keyword: string) => void
}

/** 顶栏：品牌、搜索、消息与主题切换。 */
export function TopNav({ user, onSignOut, onSearch }: TopNavProps) {
  const { theme, toggleTheme } = useTheme()
  const isDark = theme === "dark"
  const [keyword, setKeyword] = useState("")

  const unread = useGraphQL<{ unreadMessageCount: number }>(UNREAD_MESSAGE_COUNT_QUERY)
  const unreadCount = unread.data?.unreadMessageCount ?? 0

  const messages = useGraphQL<{ messages: MessageNode[] }, { unreadOnly: boolean; first: number }>(
    MESSAGES_QUERY,
    { unreadOnly: true, first: 5 },
  )
  const items = messages.data?.messages ?? []

  const openMessage = async (message: MessageNode) => {
    try {
      await graphqlRequest<{ markMessageRead: boolean }, { uuid: string }>(
        MARK_MESSAGE_READ_MUTATION,
        { uuid: message.UUID },
        getToken(),
      )
    } finally {
      messages.refetch()
      unread.refetch()
    }
  }

  return (
    <header className="flex h-14 items-center gap-4 border-b border-border bg-card px-4">
      <div className="flex items-center gap-2">
        <span className="flex size-8 items-center justify-center rounded-md bg-primary text-primary-foreground">
          <IconDatabase size={18} stroke={1.5} />
        </span>
        <span className="text-xl font-semibold tracking-tight text-strong">Halo</span>
      </div>

      <form
        className="relative w-full max-w-md"
        onSubmit={(event) => {
          event.preventDefault()
          onSearch(keyword.trim())
        }}
      >
        <IconSearch
          size={16}
          stroke={1.5}
          className="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-muted-foreground"
        />
        <Input
          value={keyword}
          onChange={(event) => setKeyword(event.target.value)}
          placeholder="搜索工单主题或目标库，回车确认"
          aria-label="搜索工单"
          className="pl-9"
        />
      </form>

      <div className="ml-auto flex items-center gap-1">
        <ConnectionMenu user={user} onSignOut={onSignOut} />

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="relative" aria-label="消息">
              <IconBell size={18} stroke={1.5} />
              {unreadCount > 0 ? (
                <span className="absolute top-1 right-1 flex size-4 items-center justify-center rounded-full bg-danger text-[0.625rem] font-medium text-danger-foreground ring-2 ring-card">
                  {unreadCount}
                </span>
              ) : null}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-80">
            <DropdownMenuLabel>未读消息</DropdownMenuLabel>
            <DropdownMenuSeparator />
            {items.length === 0 ? (
              <div className="px-2 py-3 text-sm text-muted-foreground">暂无未读消息</div>
            ) : (
              items.map((message) => (
                <DropdownMenuItem
                  key={message.UUID}
                  onSelect={() => void openMessage(message)}
                  className="items-start"
                >
                  <IconBell size={16} stroke={1.5} className="mt-0.5" />
                  <span className="flex flex-col">
                    <span>{message.Title}</span>
                    <span className="text-sm text-muted-foreground">{formatUnix(message.CreateAt)}</span>
                  </span>
                </DropdownMenuItem>
              ))
            )}
          </DropdownMenuContent>
        </DropdownMenu>

        <Button
          variant="ghost"
          size="icon"
          onClick={toggleTheme}
          aria-label={isDark ? "切换到浅色主题" : "切换到深色主题"}
          title={isDark ? "切换到浅色主题" : "切换到深色主题"}
        >
          {isDark ? <IconSun size={18} stroke={1.5} /> : <IconMoon size={18} stroke={1.5} />}
        </Button>
      </div>
    </header>
  )
}
