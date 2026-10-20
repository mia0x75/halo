import { useState } from "react"
import { IconLogin, IconLogout, IconPlugConnected } from "@tabler/icons-react"

import { LoginForm } from "@/features/auth/login-form"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { clearToken, getToken } from "@/lib/api/session"

export interface ConnectionMenuProps {
  onChanged: () => void
}

/**
 * 连接状态：未登录时展示「演示数据」，点击打开登录框；
 * 登录后展示「已连接」，菜单里可以退出登录。
 */
export function ConnectionMenu({ onChanged }: ConnectionMenuProps) {
  const [open, setOpen] = useState(false)
  const [signedIn, setSignedIn] = useState(() => getToken() !== null)

  if (!signedIn) {
    return (
      <Dialog
        open={open}
        onOpenChange={(next) => {
          setOpen(next)
          if (!next) {
            setSignedIn(getToken() !== null)
            onChanged()
          }
        }}
      >
        <DialogTrigger asChild>
          <Button variant="ghost" size="sm" className="gap-1.5 text-muted-foreground">
            <IconLogin size={16} stroke={1.5} />
            演示数据
          </Button>
        </DialogTrigger>
        <DialogContent>
          <DialogTitle className="text-base font-medium text-strong">连接后端</DialogTitle>
          <DialogDescription className="mt-1 text-sm text-muted-foreground">
            登录后概览与消息会读取真实接口；未登录或接口不可用时继续使用演示数据。
          </DialogDescription>
          <div className="mt-4 flex justify-center">
            <LoginForm
              onLoggedIn={() => {
                setSignedIn(true)
                setOpen(false)
                onChanged()
              }}
            />
          </div>
        </DialogContent>
      </Dialog>
    )
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="sm" className="gap-1.5 text-success">
          <IconPlugConnected size={16} stroke={1.5} />
          已连接
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-44">
        <DropdownMenuLabel>已连接后端</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          onSelect={() => {
            clearToken()
            setSignedIn(false)
            onChanged()
          }}
        >
          <IconLogout size={16} stroke={1.5} />
          退出登录
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
