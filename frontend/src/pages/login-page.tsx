import { IconDatabase, IconMoon, IconPlayerPlay, IconSun } from "@tabler/icons-react"

import { LoginForm } from "@/features/auth/login-form"
import { Button } from "@/components/ui/button"
import { useTheme } from "@/hooks/use-theme"
import { startDemoSession } from "@/lib/api/session"

/** 登录页：未登录时唯一的入口，登录成功由会话闸门切到主界面。 */
export function LoginPage() {
  const { theme, toggleTheme } = useTheme()
  const isDark = theme === "dark"

  return (
    <div className="flex min-h-svh flex-col bg-background text-foreground">
      <div className="flex justify-end px-4 py-3">
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
      <main className="flex flex-1 items-center justify-center px-4 pb-12">
        <div className="flex w-full max-w-sm flex-col gap-6">
          <div className="flex flex-col gap-4">
            <div className="flex items-center gap-2">
              <span className="flex size-9 items-center justify-center rounded-md bg-primary text-primary-foreground">
                <IconDatabase size={20} stroke={1.5} />
              </span>
              <span className="text-2xl font-semibold tracking-tight text-strong">Halo</span>
            </div>
            <div className="flex flex-col gap-1">
              <h1 className="text-xl font-medium text-strong">登录</h1>
              <p className="text-sm text-muted-foreground">SQL 变更工单与审计平台</p>
            </div>
          </div>

          <div className="rounded-lg border border-border bg-card p-6 shadow-xs">
            <LoginForm />
          </div>

          <div className="flex flex-col gap-2">
            <div className="flex items-center gap-3 text-xs text-muted-foreground">
              <span className="h-px flex-1 bg-border" />
              或者
              <span className="h-px flex-1 bg-border" />
            </div>
            <Button variant="outline" size="lg" className="w-full" onClick={startDemoSession}>
              <IconPlayerPlay size={16} stroke={1.5} />
              进入演示系统
            </Button>
            <p className="text-center text-xs text-muted-foreground">
              演示数据保存在浏览器内，每天重置一次，不连接真实数据库。
            </p>
          </div>

          <p className="text-xs text-muted-foreground">请使用平台账号登录，未登录无法访问任何页面。</p>
        </div>
      </main>
    </div>
  )
}
