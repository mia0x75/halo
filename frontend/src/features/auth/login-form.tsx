import { useState } from "react"
import { IconLogin } from "@tabler/icons-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { graphqlRequest } from "@/lib/api/client"
import { LOGIN_MUTATION, type LoginPayloadNode } from "@/lib/api/queries"
import { setToken } from "@/lib/api/session"

/** 登录表单：拿到令牌后写入本地存储，会话闸门会重新校验并放行。 */
export function LoginForm() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const submit = async (event: React.FormEvent) => {
    event.preventDefault()
    setLoading(true)
    setError(null)
    try {
      const data = await graphqlRequest<{ login: LoginPayloadNode }, { input: { Email: string; Password: string } }>(
        LOGIN_MUTATION,
        { input: { Email: email.trim(), Password: password } },
      )
      setToken(data.login.Token)
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : String(cause))
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={submit} className="flex w-full flex-col gap-4">
      <div className="flex flex-col gap-1.5">
        <Label htmlFor="login-email">邮箱</Label>
        <Input
          id="login-email"
          type="email"
          autoComplete="username"
          autoFocus
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          placeholder="you@example.com"
        />
      </div>

      <div className="flex flex-col gap-1.5">
        <Label htmlFor="login-password">密码</Label>
        <Input
          id="login-password"
          type="password"
          autoComplete="current-password"
          value={password}
          onChange={(event) => setPassword(event.target.value)}
        />
      </div>

      {error ? (
        <p role="alert" className="text-sm text-danger">
          {error}
        </p>
      ) : null}

      <Button type="submit" size="lg" className="w-full" disabled={loading || !email.trim() || !password}>
        <IconLogin size={16} stroke={1.5} />
        {loading ? "登录中…" : "登录"}
      </Button>
    </form>
  )
}
