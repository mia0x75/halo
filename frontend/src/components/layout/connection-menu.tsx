import { IconLogout, IconRefresh, IconUserCircle } from "@tabler/icons-react"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { clearToken, isDemoSession } from "@/lib/api/session"
import { resetDemoData } from "@/lib/demo/execute"
import type { MeUser } from "@/lib/api/queries"

export interface ConnectionMenuProps {
  user: MeUser
  onSignOut: () => void
}

/** 当前账号：展示邮箱与角色；演示模式额外提供数据重置。 */
export function ConnectionMenu({ user, onSignOut }: ConnectionMenuProps) {
  const name = user.Name?.trim() || user.Email
  const roles = user.Roles ?? []
  const demo = isDemoSession()

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="sm" className="gap-1.5">
          <IconUserCircle size={16} stroke={1.5} />
          <span className="max-w-40 truncate">{name}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-52">
        <DropdownMenuLabel className="flex flex-col gap-0.5">
          <span className="truncate text-sm font-medium text-strong">{name}</span>
          <span className="truncate text-xs font-normal text-muted-foreground">{user.Email}</span>
        </DropdownMenuLabel>
        {roles.length > 0 ? (
          <div className="px-2 pb-1 text-xs text-muted-foreground">
            {roles.map((role) => role.Name).join(" · ")}
          </div>
        ) : null}
        <DropdownMenuSeparator />
        {demo ? (
          <DropdownMenuItem
            onSelect={() => {
              resetDemoData()
              window.location.reload()
            }}
          >
            <IconRefresh size={16} stroke={1.5} />
            重置演示数据
          </DropdownMenuItem>
        ) : null}
        <DropdownMenuItem
          onSelect={() => {
            clearToken()
            onSignOut()
          }}
        >
          <IconLogout size={16} stroke={1.5} />
          退出登录
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
