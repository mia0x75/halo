/**
 * 会话令牌。令牌放在 localStorage，刷新后仍然有效；
 * 没有令牌时接口调用会失败，页面自动回落到假数据。
 */

const TOKEN_KEY = "halo-token"

export function getToken(): string | null {
  try {
    return localStorage.getItem(TOKEN_KEY)
  } catch {
    return null
  }
}

export function setToken(token: string): void {
  try {
    localStorage.setItem(TOKEN_KEY, token)
  } catch {
    // 隐私模式下写入会失败，忽略即可（后续请求会自动回落到假数据）。
  }
}

export function clearToken(): void {
  try {
    localStorage.removeItem(TOKEN_KEY)
  } catch {
    // 同上。
  }
}
