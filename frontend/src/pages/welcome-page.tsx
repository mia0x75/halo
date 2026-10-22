import { useMemo, useState } from "react"
import {
  IconArrowLeft,
  IconArrowRight,
  IconCheck,
  IconChevronDown,
  IconChevronRight,
  IconDatabase,
  IconKey,
  IconPlugConnected,
} from "@tabler/icons-react"

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
import {
  randomSecret,
  runSetup,
  testSetupConnection,
  type SetupPayload,
} from "@/lib/api/setup"

export interface WelcomePageProps {
  /** 初始化完成后的回调，通常是刷新页面。 */
  onDone: () => void
}

const STEPS = [
  { key: "database", title: "设置数据库连接" },
  { key: "secret", title: "设置签名密钥" },
  { key: "admin", title: "创建 root 用户" },
] as const

const MIN_JWT_LENGTH = 32
const MIN_CRYPTO_LENGTH = 16
const MIN_PASSWORD_LENGTH = 8

/** 首次安装向导：分三步填数据库、密钥与 root 用户。 */
export function WelcomePage({ onDone }: WelcomePageProps) {
  const [step, setStep] = useState(0)
  const [driver, setDriver] = useState("mysql")
  const [dbUser, setDbUser] = useState("")
  const [dbPassword, setDbPassword] = useState("")
  const [dbHost, setDbHost] = useState("127.0.0.1")
  const [dbPort, setDbPort] = useState("3306")
  const [dbName, setDbName] = useState("")
  const [showPool, setShowPool] = useState(false)
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
  const [testing, setTesting] = useState(false)
  const [tested, setTested] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  // 连接串由各字段拼出来，用户只填看得懂的部分。
  const dsn = useMemo(() => {
    const user = encodeURIComponent(dbUser)
    const password = dbPassword === "" ? "" : `:${encodeURIComponent(dbPassword)}`
    if (driver === "postgres") {
      return `postgres://${user}${password}@${dbHost}:${dbPort}/${dbName}?sslmode=disable`
    }
    return `${dbUser}:${dbPassword}@tcp(${dbHost}:${dbPort})/${dbName}?charset=utf8mb4&parseTime=true&loc=UTC`
  }, [driver, dbUser, dbPassword, dbHost, dbPort, dbName])

  const databasePayload = (): SetupPayload["database"] => ({
    driver,
    addr: dsn,
    max_idle: Number(maxIdle) || 10,
    max_connections: Number(maxConnections) || 100,
    wait_timeout: Number(waitTimeout) || 300,
  })

  // 每一步自己校验，没填完不许往下走。
  const stepError = (index: number): string | null => {
    if (index === 0) {
      if (dbUser === "" || dbHost === "" || dbPort === "" || dbName === "") {
        return "请填写用户、地址、端口与数据库名"
      }
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

  const testConnection = async () => {
    const problem = stepError(0)
    if (problem) {
      setError(problem)
      return
    }
    setError(null)
    setTested(null)
    setTesting(true)
    try {
      await testSetupConnection(databasePayload())
      setTested("连接成功")
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setTesting(false)
    }
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
      database: databasePayload(),
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
        </div>

        <section className="overflow-hidden rounded-lg border border-border bg-card shadow-xs">
          <div className="flex flex-col gap-1 border-b border-border px-6 py-4">
            <h2 className="text-base font-medium text-strong">
              第 {step + 1} 步 · {STEPS[step].title}
            </h2>
          </div>

          <div className="flex flex-col gap-4 px-6 py-5">
            {step === 0 ? (
              <>
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="flex flex-col gap-1.5">
                    <Label htmlFor="setup-db-user">用户</Label>
                    <Input
                      id="setup-db-user"
                      value={dbUser}
                      onChange={(event) => setDbUser(event.target.value)}
                      autoFocus
                    />
                  </div>
                  <div className="flex flex-col gap-1.5">
                    <Label htmlFor="setup-db-password">密码</Label>
                    <Input
                      id="setup-db-password"
                      type="password"
                      value={dbPassword}
                      onChange={(event) => setDbPassword(event.target.value)}
                    />
                  </div>
                  <div className="flex flex-col gap-1.5">
                    <Label htmlFor="setup-db-host">地址</Label>
                    <Input
                      id="setup-db-host"
                      value={dbHost}
                      onChange={(event) => setDbHost(event.target.value)}
                    />
                  </div>
                  <div className="flex flex-col gap-1.5">
                    <Label htmlFor="setup-db-port">端口</Label>
                    <Input
                      id="setup-db-port"
                      value={dbPort}
                      onChange={(event) => setDbPort(event.target.value)}
                    />
                  </div>
                  <div className="flex flex-col gap-1.5">
                    <Label htmlFor="setup-db-name">数据库</Label>
                    <Input
                      id="setup-db-name"
                      value={dbName}
                      onChange={(event) => setDbName(event.target.value)}
                    />
                  </div>
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
                </div>

                <div className="flex flex-col gap-1.5">
                  <Label htmlFor="setup-dsn">连接串</Label>
                  <Input
                    id="setup-dsn"
                    readOnly
                    value={dsn}
                    className="font-mono text-xs"
                    onFocus={(event) => event.currentTarget.select()}
                  />
                </div>

                <div className="flex flex-wrap items-center gap-3">
                  <Button variant="outline" onClick={() => void testConnection()} disabled={testing}>
                    <IconPlugConnected size={16} stroke={1.5} />
                    {testing ? "测试中…" : "测试连接"}
                  </Button>
                  {tested ? <span className="text-sm text-success">{tested}</span> : null}
                </div>

                <div className="flex flex-col gap-1">
                  <button
                    type="button"
                    onClick={() => setShowPool((value) => !value)}
                    className="flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
                  >
                    {showPool ? (
                      <IconChevronDown size={14} stroke={1.5} />
                    ) : (
                      <IconChevronRight size={14} stroke={1.5} />
                    )}
                    连接池参数
                  </button>
                  {showPool ? (
                    <div className="grid gap-4 pt-1 sm:grid-cols-3">
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
                  ) : null}
                </div>
              </>
            ) : null}

            {step === 1 ? (
              <>
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

          <div className="flex items-center justify-between gap-4 border-t border-border bg-surface-secondary px-6 py-4">
            {/* 进度信息在左下角 */}
            <div className="flex w-40 flex-col gap-1.5">
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
              <span className="text-xs text-muted-foreground">
                {step + 1} / {STEPS.length} · {STEPS[step].title} · {percent}%
              </span>
            </div>
            {/* 步骤按钮在右下角 */}
            <div className="flex items-center gap-2">
              <Button variant="ghost" onClick={goBack} disabled={step === 0 || submitting}>
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
          </div>
        </section>
      </div>
    </div>
  )
}
