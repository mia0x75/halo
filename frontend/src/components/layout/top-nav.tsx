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
import { useGraphQL } from "@/hooks/use-graphql"
import { useTheme } from "@/hooks/use-theme"
import { UNREAD_MESSAGE_COUNT_QUERY } from "@/lib/api/queries"

export interface TopNavProps {
  onConnectionChange: () => void
}

/** 顶栏：品牌、搜索、消息与主题切换。 */
export function TopNav({ onConnectionChange }: TopNavProps) {
  const { theme, toggleTheme } = useTheme()
  const isDark = theme === "dark"

  const unread = useGraphQL<{ unreadMessageCount: number }>(
    UNREAD_MESSAGE_COUNT_QUERY,
    undefined,
    { unreadMessageCount: 3 },
  )
  const unreadCount = unread.data?.unreadMessageCount ?? 3

  return (
    <header className="flex h-14 items-center gap-4 border-b border-border bg-card px-4">
      <div className="flex items-center gap-2">
        <span className="flex size-8 items-center justify-center rounded-md bg-primary text-primary-foreground">
          <IconDatabase size={18} stroke={1.5} />
        </span>
        <span className="text-xl font-semibold tracking-tight text-strong">Halo</span>
      </div>

      <div className="relative w-full max-w-md">
        <IconSearch
          size={16}
          stroke={1.5}
          className="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-muted-foreground"
        />
        <Input placeholder="搜索工单、群集、规则…" className="pl-9" />
      </div>

      <div className="ml-auto flex items-center gap-1">
        <ConnectionMenu onChanged={onConnectionChange} />

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
          <DropdownMenuContent align="end" className="w-72">
            <DropdownMenuLabel>消息</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <IconBell size={16} stroke={1.5} />
              <span className="flex flex-col">
                <span>工单「订单归档字段」审核通过</span>
                <span className="text-sm text-muted-foreground">2 分钟前</span>
              </span>
            </DropdownMenuItem>
            <DropdownMenuItem>
              <IconBell size={16} stroke={1.5} />
              <span className="flex flex-col">
                <span>工单「报表视图重建」待审核</span>
                <span className="text-sm text-muted-foreground">1 小时前</span>
              </span>
            </DropdownMenuItem>
            <DropdownMenuItem>
              <IconBell size={16} stroke={1.5} />
              <span className="flex flex-col">
                <span>工单「清理历史数据」执行失败</span>
                <span className="text-sm text-muted-foreground">昨天</span>
              </span>
            </DropdownMenuItem>
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
