import { useState } from "react"
import { IconArrowLeft, IconArrowRight, IconCheck, IconDatabase, IconKey } from "@tabler/icons-react"

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

const STEPS = [
  { key: "database", title: "设置数据库连接", hint: "元数据库存放工单、规则与系统选项" },
  { key: "secret", title: "设置签名密钥", hint: "令牌用 JWT 密钥签发，密码用加密密钥加密" },
  { key: "admin", title: "创建 root 用户", hint: "用这个账号登录平台，拥有全部权限" },
] as const

const MIN_JWT_LENGTH = 32
const MIN_CRYPTO_LENGTH = 16
const MIN_PASSWORD_LENGTH = 8

/** 首次安装向导：分三步填数据库、密钥与 root 用户。 */
export function WelcomePage({ onDone }: WelcomePageProps) {
  const [step, setStep] = useState(0)
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

  // 每一步自己校验，没填完不许往下走。
  const stepError = (index: number): string | null => {
    if (index === 0 && addr.trim() === "") {
      return "请填写元数据库连接串"
    }
    if (index === 1) {
      if (secretKey.length < MIN_JWT_LENGTH) {
        return `JWT 密钥至少需要 ${MIN_JWT_LENGTH} 个字符`
      }
      if (cryptoKey.length < MIN_CRYPTO_LENGTH) {
        return `加密密钥至少需要 ${MIN_CRYPTO_LENGTH} 个字符`
      }
    }
    if (index === 2) {
      if (!email.includes("@")) {
        return "请填写 root 用户的邮箱"
      }
      if (password.length < MIN_PASSWORD_LENGTH) {
        return `root 用户密码至少需要 ${MIN_PASSWORD_LENGTH} 位`
      }
    }
    return null
  }

  const goNext = () => {
    const problem = stepError(step)
    if (problem) {
      setError(problem)
      return
    }
    setError(null)
    setStep((current) => Math.min(current + 1, STEPS.length - 1))
  }

  const goBack = () => {
    setError(null)
    setStep((current) => Math.max(current - 1, 0))
  }

  const submit = async () => {
    const problem = stepError(2)
    if (problem) {
      setError(problem)
      return
    }
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

  const percent = Math.round(((step + 1) / STEPS.length) * 100)

  return (
    <div className="min-h-svh bg-background px-4 py-8 text-foreground">
      <div className="mx-auto flex w-full max-w-2xl flex-col gap-6">
        <div className="flex flex-col items-center gap-3 text-center">
          <span className="flex size-11 items-center justify-center rounded-md bg-primary text-primary-foreground">
            <IconDatabase size={24} stroke={1.5} />
          </span>
          <h1 className="text-xl font-medium text-strong">欢迎使用 Halo</h1>
          <p className="max-w-md text-sm text-muted-foreground">
            初始化只需三步：连上元数据库、生成签名密钥、创建 root 用户。完成后配置会写入服务器配置文件。
          </p>
        </div>

        <section className="overflow-hidden rounded-lg border border-border bg-card shadow-xs">
          <div className="flex flex-col gap-1 border-b border-border px-6 py-4">
            <h2 className="text-base font-medium text-strong">
              第 {step + 1} 步 · {STEPS[step].title}
            </h2>
            <p className="text-sm text-muted-foreground">{STEPS[step].hint}</p>
          </div>

          <div className="flex flex-col gap-4 px-6 py-5">
            {step === 0 ? (
              <>
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
                      autoFocus
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
              </>
            ) : null}

            {step === 1 ? (
              <>
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="flex flex-col gap-1.5">
                    <Label htmlFor="setup-token-name">令牌请求头</Label>
                    <Input
                      id="setup-token-name"
                      value={tokenName}
                      onChange={(event) => setTokenName(event.target.value)}
                      autoFocus
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
                  <Label htmlFor="setup-secret-key">JWT 密钥（至少 {MIN_JWT_LENGTH} 字符）</Label>
                  <div className="flex items-center gap-2">
                    <Input
                      id="setup-secret-key"
                      value={secretKey}
                      onChange={(event) => setSecretKey(event.target.value)}
                    />
                    <Button
                      type="button"
                      variant="outline"
                      size="icon"
                      aria-label="生成 JWT 密钥"
                      onClick={() => setSecretKey(randomSecret(32))}
                    >
                      <IconKey size={16} stroke={1.5} />
                    </Button>
                  </div>
                </div>
                <div className="flex flex-col gap-1.5">
                  <Label htmlFor="setup-crypto">加密密钥（至少 {MIN_CRYPTO_LENGTH} 字符）</Label>
                  <div className="flex items-center gap-2">
                    <Input
                      id="setup-crypto"
                      value={cryptoKey}
                      onChange={(event) => setCryptoKey(event.target.value)}
                    />
                    <Button
                      type="button"
                      variant="outline"
                      size="icon"
                      aria-label="生成加密密钥"
                      onClick={() => setCryptoKey(randomSecret(16))}
                    >
                      <IconKey size={16} stroke={1.5} />
                    </Button>
                  </div>
                </div>
              </>
            ) : null}

            {step === 2 ? (
              <>
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="flex flex-col gap-1.5">
                    <Label htmlFor="setup-admin-email">邮箱</Label>
                    <Input
                      id="setup-admin-email"
                      type="email"
                      value={email}
                      onChange={(event) => setEmail(event.target.value)}
                      placeholder="admin@example.com"
                      autoFocus
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
                  <Label htmlFor="setup-admin-password">密码（至少 {MIN_PASSWORD_LENGTH} 位）</Label>
                  <Input
                    id="setup-admin-password"
                    type="password"
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                  />
                </div>
              </>
            ) : null}

            {error ? (
              <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
            ) : null}
          </div>

          <div className="flex items-center justify-between gap-2 border-t border-border bg-surface-secondary px-6 py-4">
            <Button
              variant="ghost"
              onClick={goBack}
              disabled={step === 0 || submitting}
            >
              <IconArrowLeft size={16} stroke={1.5} />
              上一步
            </Button>
            {step < STEPS.length - 1 ? (
              <Button onClick={goNext}>
                下一步
                <IconArrowRight size={16} stroke={1.5} />
              </Button>
            ) : (
              <Button onClick={() => void submit()} disabled={submitting}>
                <IconCheck size={16} stroke={1.5} />
                {submitting ? "正在初始化…" : "完成初始化"}
              </Button>
            )}
          </div>
        </section>
      </div>

      {/* 左下角进度条 */}
      <div className="fixed bottom-5 left-5 flex w-56 flex-col gap-2">
        <div
          className="h-1 w-full overflow-hidden rounded-full bg-border"
          role="progressbar"
          aria-valuenow={percent}
          aria-valuemin={0}
          aria-valuemax={100}
          aria-label={`初始化进度 ${percent}%`}
        >
          <div
            className="h-full rounded-full bg-primary transition-[width]"
            style={{ width: `${percent}%` }}
          />
        </div>
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <span>
            {step + 1} / {STEPS.length} · {STEPS[step].title}
          </span>
          <span>{percent}%</span>
        </div>
      </div>
    </div>
  )
}
