import { useEffect, useState } from "react"

import { graphqlRequest } from "@/lib/api/client"
import { clearToken, getToken, subscribeToken } from "@/lib/api/session"
import { ME_QUERY, type MeUser } from "@/lib/api/queries"

export type SessionState =
  | { status: "anonymous" }
  | { status: "checking" }
  | { status: "authenticated"; user: MeUser }

/**
 * 会话闸门：本地有令牌还不够，得问一次 me 确认它没过期。
 * 传输层遇到 1401 会清掉令牌，这里跟着回到未登录。
 */
export function useSession(): SessionState {
  const [token, setToken] = useState<string | null>(() => getToken())
  const [user, setUser] = useState<MeUser | null>(null)

  useEffect(
    () =>
      subscribeToken((next) => {
        setToken(next)
        setUser(null)
      }),
    [],
  )

  useEffect(() => {
    if (token == null) return
    let cancelled = false
    graphqlRequest<{ me: MeUser | null }>(ME_QUERY, undefined, token)
      .then((data) => {
        if (cancelled) return
        if (data.me == null) {
          clearToken()
          return
        }
        setUser(data.me)
      })
      .catch(() => {
        if (!cancelled) clearToken()
      })
    return () => {
      cancelled = true
    }
  }, [token])

  if (token == null) return { status: "anonymous" }
  if (user == null) return { status: "checking" }
  return { status: "authenticated", user }
}
