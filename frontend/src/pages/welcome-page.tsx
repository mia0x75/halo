import { useMemo, useState } from "react"
import {
  IconArrowLeft,
  IconArrowRight,
  IconCheck,
  IconChevronDown,
  IconChevronRight,
  IconEye,
  IconEyeOff,
  IconKey,
  IconPlugConnected,
} from "@tabler/icons-react"

import { LogoMark } from "@/components/shared/logo-mark"
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
import { cn } from "@/lib/utils"

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

/** 密码强度：按长度与字符种类粗略打分。 */
function passwordStrength(value: string): { level: number; label: string } {
  let score = 0
  if (value.length >= 8) score += 1
  if (value.length >= 12) score += 1
  if (/[a-z]/.test(value) && /[A-Z]/.test(value)) score += 1
  if (/\d/.test(value) && /[^A-Za-z0-9]/.test(value)) score += 1
  const level = Math.min(score, 4)
  const labels = ["太弱", "偏弱", "一般", "较强", "很强"]
  return { level, label: value ? labels[level] : "" }
}

/** 强度条：四格进度 + 一行文字，只有填了的格子才有颜色。 */
function StrengthMeter({
  level,
  label,
  tone,
}: {
  level: number
  label: string
  /** 指定色调时不用按等级取色，用于「匹配 / 不匹配」这种提示。 */
  tone?: "danger" | "success"
}) {
  const colors = ["bg-danger", "bg-danger", "bg-warning", "bg-info", "bg-success"]
  return (
    <div className="flex w-full flex-col gap-1">
      <div className="flex gap-1">
        {[0, 1, 2, 3].map((index) => (
          <span
            key={index}
            className={cn(
              "h-1.5 flex-1 rounded-full bg-border",
              index < level && (tone === "danger" ? "bg-danger" : tone === "success" ? "bg-success" : colors[level]),
            )}
          />
        ))}
      </div>
      <span className="text-xs text-muted-foreground">{label}</span>
    </div>
  )
}

/** 匹配提示：只有一格，和强度条不是一回事。 */
function MatchHint({ ok }: { ok: boolean }) {
  return (
    <div className="flex items-center gap-2">
      <span className={cn("h-1.5 w-8 rounded-full", ok ? "bg-success" : "bg-danger")} />
      <span className="text-xs text-muted-foreground">{ok ? "匹配" : "不匹配"}</span>
    </div>
  )
}

/** 密码输入框：右侧眼睛图标切换明文。 */
function PasswordInput({
  id,
  value,
  onChange,
  autoComplete,
}: {
  id: string
  value: string
  onChange: (value: string) => void
  autoComplete?: string
}) {
  const [visible, setVisible] = useState(false)
  return (
    <div className="relative">
      <Input
        id={id}
        type={visible ? "text" : "password"}
        value={value}
        autoComplete={autoComplete}
        onChange={(event) => onChange(event.target.value)}
        className="pr-10"
      />
      <button
        type="button"
        aria-label={visible ? "隐藏密码" : "显示密码"}
        onClick={() => setVisible((current) => !current)}
        className="absolute top-1/2 right-3 -translate-y-1/2 text-muted-foreground transition-colors hover:text-foreground"
      >
        {visible ? <IconEyeOff size={16} stroke={1.5} /> : <IconEye size={16} stroke={1.5} />}
      </button>
    </div>
  )
}

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
  // 密钥默认生成一套，管理员可以随时换。
  const [secretKey, setSecretKey] = useState(() => randomSecret(32))
  const [cryptoKey, setCryptoKey] = useState(() => randomSecret(16))
  const [email, setEmail] = useState("")
  const [name, setName] = useState("")
  const [password, setPassword] = useState("")
  const [passwordConfirm, setPasswordConfirm] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [testing, setTesting] = useState(false)
  const [testResult, setTestResult] = useState<{ ok: boolean; message: string } | null>(null)
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
      if (password !== passwordConfirm) {
        return "两次输入的密码不一致"
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
    setTestResult(null)
    setTesting(true)
    try {
      await testSetupConnection(databasePayload())
      setTestResult({ ok: true, message: "连接成功" })
    } catch (cause) {
      // 只在按钮旁提示成败，详细原因看服务端日志。
      const message = cause instanceof Error ? cause.message : "连接失败，请检查日志"
      setTestResult({ ok: false, message })
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
      {/* Root font size is 14px, so rem here maps to 14px: 30rem is ~420px.
          Fields are stacked one per row, so the card stays narrower than the
          old 576px one without squeezing any input. */}
      <div className="mx-auto flex w-full max-w-[30rem] flex-col gap-6">
        <div className="flex flex-col items-center gap-3 text-center">
          <LogoMark size={48} className="shrink-0" />
          <h1 className="text-xl font-medium text-strong">欢迎使用 Halo</h1>
        </div>

        <section className="overflow-hidden rounded-lg border border-border bg-card shadow-xs">
          <div className="flex items-center justify-between gap-3 border-b border-border px-5 py-3.5">
            <h2 className="text-base font-medium text-strong">
              第 {step + 1} 步 · {STEPS[step].title}
            </h2>
            {step === 0 ? (
              <Select value={driver} onValueChange={(value) => setDriver(value)}>
                <SelectTrigger id="setup-driver" className="h-8 w-32">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="mysql">MySQL</SelectItem>
                  <SelectItem value="postgres">PostgreSQL</SelectItem>
                </SelectContent>
              </Select>
            ) : null}
          </div>

          <div className="flex flex-col gap-4 px-5 py-4">
            {step === 0 ? (
              <>
                <div className="grid grid-cols-[minmax(0,1fr)_6rem] gap-3">
                  <div className="flex flex-col gap-1.5">
                    <Label htmlFor="setup-db-host">地址</Label>
                    <Input
                      id="setup-db-host"
                      value={dbHost}
                      onChange={(event) => setDbHost(event.target.value)}
                      autoFocus
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
                </div>

                <div className="flex flex-col gap-1.5">
                  <Label htmlFor="setup-db-name">数据库</Label>
                  <Input
                    id="setup-db-name"
                    value={dbName}
                    onChange={(event) => setDbName(event.target.value)}
                    placeholder="halodb"
                  />
                </div>

                <div className="flex flex-col gap-1.5">
                  <Label htmlFor="setup-db-user">用户</Label>
                  <Input
                    id="setup-db-user"
                    value={dbUser}
                    onChange={(event) => setDbUser(event.target.value)}
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

                <div className="flex flex-col gap-1">
                  <button
                    type="button"
                    onClick={() => setShowPool((value) => !value)}
                    className="flex items-center gap-1.5 self-start text-sm text-muted-foreground transition-colors hover:text-foreground"
                  >
                    {showPool ? (
                      <IconChevronDown size={14} stroke={1.5} />
                    ) : (
                      <IconChevronRight size={14} stroke={1.5} />
                    )}
                    连接池参数
                  </button>
                  {showPool ? (
                    <div className="grid grid-cols-3 gap-3 pt-1">
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

                <div className="flex flex-wrap items-center gap-3">
                  <Button variant="outline" onClick={() => void testConnection()} disabled={testing}>
                    <IconPlugConnected size={16} stroke={1.5} />
                    {testing ? "测试中…" : "测试连接"}
                  </Button>
                  {testResult ? (
                    <span
                      className={testResult.ok ? "text-sm text-success" : "text-sm text-danger"}
                    >
                      {testResult.message}
                    </span>
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
                      className="size-[2.375rem] shrink-0"
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
                      className="size-[2.375rem] shrink-0"
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
                  <PasswordInput
                    id="setup-admin-password"
                    value={password}
                    onChange={setPassword}
                    autoComplete="new-password"
                  />
                  {password ? (
                    <StrengthMeter
                      level={passwordStrength(password).level}
                      label={passwordStrength(password).label}
                    />
                  ) : null}
                </div>
                <div className="flex flex-col gap-1.5">
                  <Label htmlFor="setup-admin-password-confirm">确认密码</Label>
                  <PasswordInput
                    id="setup-admin-password-confirm"
                    value={passwordConfirm}
                    onChange={setPasswordConfirm}
                    autoComplete="new-password"
                  />
                  {passwordConfirm ? (
                    <MatchHint ok={passwordConfirm === password} />
                  ) : null}
                </div>
              </>
            ) : null}

            {error ? (
              <p className="rounded-md bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
            ) : null}
          </div>

          <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border bg-surface-secondary px-5 py-4">
            {/* 进度信息在左下角 */}
            <div className="flex shrink-0 flex-col gap-1.5">
              <div
                className="h-2 w-24 overflow-hidden rounded-full bg-border"
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
              <span className="whitespace-nowrap text-xs text-muted-foreground">
                {step + 1} / {STEPS.length} · {STEPS[step].title} · {percent}%
              </span>
            </div>
            {/* 步骤按钮在右下角 */}
            <div className="flex items-center gap-2">
              <Button onClick={goBack} disabled={step === 0 || submitting}>
                <IconArrowLeft size={16} stroke={1.5} />
                上一步
              </Button>
              {step < STEPS.length - 1 ? (
                <Button onClick={goNext} disabled={stepError(step) !== null}>
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
