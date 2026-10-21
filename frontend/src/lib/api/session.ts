/**
 * 会话令牌。令牌放在 localStorage，刷新后仍然有效；
 * 失效由后端返回 1401，传输层会清掉令牌并通知外层回到登录页。
 */

const TOKEN_KEY = "halo-token"

/**
 * 演示模式的伪令牌。它不会发给后端：传输层看到它就转到本地演示执行器，
 * 因此演示环境可以脱离服务端独立运行。
 */
export const DEMO_TOKEN = "demo"

type Listener = (token: string | null) => void

const listeners = new Set<Listener>()

export function getToken(): string | null {
  try {
    return localStorage.getItem(TOKEN_KEY)
  } catch {
    return null
  }
}

function notify(token: string | null): void {
  for (const listener of listeners) {
    listener(token)
  }
}

/** 订阅令牌变化，返回退订函数。 */
export function subscribeToken(listener: Listener): () => void {
  listeners.add(listener)
  return () => {
    listeners.delete(listener)
  }
}

export function setToken(token: string): void {
  try {
    localStorage.setItem(TOKEN_KEY, token)
  } catch {
    // 隐私模式下写不进去；内存里也没有副本，等于未登录。
    return
  }
  notify(token)
}

export function clearToken(): void {
  try {
    localStorage.removeItem(TOKEN_KEY)
  } catch {
    // 忽略：读不到令牌时外层同样按未登录处理。
  }
  notify(null)
}

export function isDemoSession(): boolean {
  return getToken() === DEMO_TOKEN
}

/** 进入演示：不需要后端，数据由前端内置的演示执行器提供。 */
export function startDemoSession(): void {
  setToken(DEMO_TOKEN)
}
