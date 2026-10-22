/**
 * 首次安装向导的传输层。初始化接口不是 GraphQL：那时还没有元数据库，
 * 也没有令牌可用。
 */

export interface SetupStatus {
  /** 是否已经配置完成。false 时前端要渲染 /welcome 向导。 */
  configured: boolean
  /** 配置文件路径，向导会把它写进去。 */
  configFile: string
  driver: string
  addr: string
}

export interface SetupPayload {
  database: {
    driver: string
    addr: string
    max_idle: number
    max_connections: number
    wait_timeout: number
  }
  secret: {
    token_name: string
    issuer: string
    secret_key: string
    crypto: string
  }
  admin: {
    email: string
    name: string
    password: string
  }
}

/** 读取初始化状态。 */
export async function fetchSetupStatus(): Promise<SetupStatus> {
  const response = await fetch("/api/setup/status", {
    method: "GET",
    headers: { Accept: "application/json" },
  })
  if (!response.ok) {
    throw new Error(`读取初始化状态失败：${response.status}`)
  }
  return (await response.json()) as SetupStatus
}

/** 执行初始化。失败时抛出后端返回的错误信息。 */
export async function runSetup(payload: SetupPayload): Promise<void> {
  const response = await fetch("/api/setup", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as { error?: string } | null
    throw new Error(body?.error ?? `初始化失败：${response.status}`)
  }
}

/** 试连元数据库，不写任何东西。 */
export async function testSetupConnection(payload: SetupPayload["database"]): Promise<void> {
  const response = await fetch("/api/setup/test", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as { error?: string } | null
    throw new Error(body?.error ?? `连接测试失败：${response.status}`)
  }
}

/** 生成一段十六进制随机串，用作密钥。 */
export function randomSecret(bytes: number): string {
  const buffer = new Uint8Array(bytes)
  crypto.getRandomValues(buffer)
  return Array.from(buffer, (item) => item.toString(16).padStart(2, "0")).join("")
}
