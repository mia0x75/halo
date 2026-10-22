import { useState } from "react"
import { IconDatabase, IconKey, IconRefresh, IconUser } from "@tabler/icons-react"

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
import { randomSecret, runSetup, type SetupPayload } from "@/lib/api/setup"

export interface WelcomePageProps {
  /** 初始化完成后的回调，通常是刷新页面。 */
  onDone: () => void
}

/** 首次安装向导：写元数据库连接、密钥与初始管理员。 */
export function WelcomePage({ onDone }: WelcomePageProps) {
  const [driver, setDriver] = useState("mysql")
  const [addr, setAddr] = useState("")
  const [maxIdle, setMaxIdle] = useState("10")
  const [maxConnections, setMaxConnections] = useState("100")
  const [waitTimeout, setWaitTimeout] = useState("300")
  const [tokenName, setTokenName] = useState("Authentication")
  const [issuer, setIssuer] = useState("halo")
  const [secretKey, setSecretKey] = useState("")
  const [cryptoKey, setCryptoKey] = useState("")
  const [email, setEmail] = useState("")
  const [name, setName] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const submit = async (event: React.FormEvent) => {
    event.preventDefault()
    setError(null)
    setSubmitting(true)
    const payload: SetupPayload = {
      database: {
        driver,
        addr,
        max_idle: Number(maxIdle) || 10,
        max_connections: Number(maxConnections) || 100,
        wait_timeout: Number(waitTimeout) || 300,
      },
      secret: {
        token_name: tokenName,
        issuer,
        secret_key: secretKey,
        crypto: cryptoKey,
      },
      admin: { email, name, password },
    }
    try {
      await runSetup(payload)
      onDone()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="min-h-svh bg-background px-4 py-8 text-foreground">
      <div className="mx-auto w-full max-w-2xl">
        <form className="overflow-hidden rounded-lg border border-border bg-card shadow-xs" onSubmit={(event) => void submit(event)}>
          <div className="flex flex-col items-center gap-3 px-6 py-8 text-center">
            <span className="flex size-11 items-center justify-center rounded-md bg-primary text-primary-foreground">
              <IconDatabase size={24} stroke={1.5} />
            </span>
            <h1 className="text-xl font-medium text-strong">欢迎使用 Halo</h1>
            <p className="max-w-md text-sm text-muted-foreground">
              填写元数据库连接、签名密钥与初始管理员账号。配置会写入服务器配置文件，随后自动执行数据库迁移。
            </p>
          </div>

          <HrText label="元数据库" />
          <div className="flex flex-col gap-4 px-6 pb-6">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="setup-driver">数据库类型</Label>
                <Select value={driver} onValueChange={(value) => setDriver(value)}>
                  <SelectTrigger id="setup-driver">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="mysql">MySQL</SelectItem>
                    <SelectItem value="postgres">PostgreSQL</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="setup-addr">连接串</Label>
                <Input
                  id="setup-addr"
                  value={addr}
                  onChange={(event) => setAddr(event.target.value)}
                  placeholder="user:pass@tcp(host:3306)/halo?charset=utf8mb4&parseTime=true"
                  required
                />
              </div>
            </div>
            <div className="grid gap-4 sm:grid-cols-3">
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="setup-max-idle">最大空闲连接</Label>
                <Input
                  id="setup-max-idle"
                  type="number"
                  value={maxIdle}
                  onChange={(event) => setMaxIdle(event.target.value)}
                />
              </div>
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="setup-max-connections">最大连接数</Label>
                <Input
                  id="setup-max-connections"
                  type="number"
                  value={maxConnections}
                  onChange={(event) => setMaxConnections(event.target.value)}
                />
              </div>
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="setup-wait-timeout">连接存活（秒）</Label>
                <Input
                  id="setup-wait-timeout"
                  type="number"
                  value={waitTimeout}
                  onChange={(event) => setWaitTimeout(event.target.value)}
                />
              </div>
            </div>
          </div>

          <HrText label="签名与加密密钥" />
          <div className="flex flex-col gap-4 px-6 pb-6">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="setup-token-name">令牌请求头</Label>
                <Input
                  id="setup-token-name"
                  value={tokenName}
                  onChange={(event) => setTokenName(event.target.value)}
                />
              </div>
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="setup-issuer">令牌签发者</Label>
                <Input
                  id="setup-issuer"
                  value={issuer}
                  onChange={(event) => setIssuer(event.target.value)}
                />
              </div>
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="setup-secret-key">JWT 密钥（至少 32 字符）</Label>
              <div className="flex items-center gap-2">
                <Input
                  id="setup-secret-key"
                  value={secretKey}
                  onChange={(event) => setSecretKey(event.target.value)}
                  required
                />
                <Button
                  type="button"
                  variant="outline"
                  size="icon"
                  aria-label="生成 JWT 密钥"
                  onClick={() => setSecretKey(randomSecret(32))}
                >
                  <IconRefresh size={16} stroke={1.5} />
                </Button>
              </div>
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="setup-crypto">加密密钥（至少 16 字符）</Label>
              <div className="flex items-center gap-2">
                <Input
                  id="setup-crypto"
                  value={cryptoKey}
                  onChange={(event) => setCryptoKey(event.target.value)}
                  required
                />
                <Button
                  type="button"
                  variant="outline"
                  size="icon"
                  aria-label="生成加密密钥"
                  onClick={() => setCryptoKey(randomSecret(16))}
                >
                  <IconRefresh size={16} stroke={1.5} />
                </Button>
              </div>
            </div>
          </div>

          <HrText label="初始管理员" />
          <div className="flex flex-col gap-4 px-6 pb-6">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="setup-admin-email">邮箱</Label>
                <Input
                  id="setup-admin-email"
                  type="email"
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                  placeholder="admin@example.com"
                  required
                />
              </div>
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="setup-admin-name">姓名</Label>
                <Input
                  id="setup-admin-name"
                  value={name}
                  onChange={(event) => setName(event.target.value)}
                  placeholder="管理员"
                />
              </div>
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="setup-admin-password">密码（至少 8 位）</Label>
              <Input
                id="setup-admin-password"
                type="password"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                required
              />
            </div>
          </div>

          <div className="flex flex-col gap-3 border-t border-border bg-surface-secondary px-6 py-4">
            {error ? (
              <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
            ) : null}
            <div className="flex items-center justify-end gap-2">
              <Button type="submit" disabled={submitting}>
                <IconUser size={16} stroke={1.5} />
                {submitting ? "正在初始化…" : "完成初始化"}
              </Button>
            </div>
          </div>
        </form>

        <p className="mt-4 flex items-center justify-center gap-1.5 text-xs text-muted-foreground">
          <IconKey size={13} stroke={1.5} />
          密钥只写入服务器配置文件，不会回显。
        </p>
      </div>
    </div>
  )
}

/** 分节标题：一条横线中间夹小字。 */
function HrText({ label }: { label: string }) {
  return (
    <div className="flex items-center gap-3 px-6 pb-2 text-xs font-medium tracking-wide text-muted-foreground uppercase">
      <span className="h-px flex-1 bg-border" />
      {label}
      <span className="h-px flex-1 bg-border" />
    </div>
  )
}
