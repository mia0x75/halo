import {
  IconBell,
  IconChevronDown,
  IconLogout,
  IconMoon,
  IconSearch,
  IconSettings,
  IconSun,
  IconUserCircle,
} from "@tabler/icons-react"

import { Avatar } from "@/components/ui/avatar"
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
import { useTheme } from "@/hooks/use-theme"

/** 顶栏：搜索框、消息、主题切换、账号菜单。 */
export function TopNav() {
  const { theme, toggleTheme } = useTheme()
  const isDark = theme === "dark"

  return (
    <header className="flex h-14 items-center gap-4 border-b border-border bg-card px-4">
      <div className="relative w-full max-w-md">
        <IconSearch
          size={16}
          stroke={1.5}
          className="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-muted-foreground"
        />
        <Input placeholder="搜索工单、群集、规则…" className="h-9 pl-9" />
      </div>

      <div className="ml-auto flex items-center gap-1">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="relative" aria-label="消息">
              <IconBell size={18} stroke={1.5} />
              <span className="absolute top-1.5 right-1.5 size-2 rounded-full bg-danger ring-2 ring-card" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-72">
            <DropdownMenuLabel>消息</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <IconBell size={16} stroke={1.5} />
              <span className="flex flex-col">
                <span>工单「订单归档字段」审核通过</span>
                <span className="text-xs text-muted-foreground">2 分钟前</span>
              </span>
            </DropdownMenuItem>
            <DropdownMenuItem>
              <IconBell size={16} stroke={1.5} />
              <span className="flex flex-col">
                <span>工单「报表视图重建」待审核</span>
                <span className="text-xs text-muted-foreground">1 小时前</span>
              </span>
            </DropdownMenuItem>
            <DropdownMenuItem>
              <IconBell size={16} stroke={1.5} />
              <span className="flex flex-col">
                <span>工单「清理历史数据」执行失败</span>
                <span className="text-xs text-muted-foreground">昨天</span>
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

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button
              type="button"
              className="flex items-center gap-2 rounded-md px-2 py-1.5 transition-colors hover:bg-surface-secondary"
              aria-label="账号菜单"
            >
              <Avatar label="张" size="sm" />
              <span className="hidden text-sm md:inline">张三</span>
              <IconChevronDown size={14} stroke={1.5} className="text-muted-foreground" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-48">
            <DropdownMenuLabel>张三 · 管理员</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <IconUserCircle size={16} stroke={1.5} />
              个人资料
            </DropdownMenuItem>
            <DropdownMenuItem>
              <IconSettings size={16} stroke={1.5} />
              账号设置
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <IconLogout size={16} stroke={1.5} />
              退出登录
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  )
}
