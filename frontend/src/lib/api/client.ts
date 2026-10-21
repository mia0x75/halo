/**
 * GraphQL 传输层。只做一件事：把查询发到后端并把 data 取出来。
 * 令牌失效（后端 1401）时顺手清掉本地令牌，外层随即回到登录页。
 */

import { clearToken } from "@/lib/api/session"

const ENDPOINT = "/api/query"

/** 后端错误码：令牌缺失/被收回/过期。 */
const CODE_UNAUTHORIZED = 1401

export interface GraphQLError {
  message: string
}

interface GraphQLResponse<TData> {
  data?: TData
  errors?: GraphQLError[]
}

export class ApiError extends Error {
  readonly code: number | null

  constructor(message: string, code: number | null = null) {
    super(message)
    this.name = "ApiError"
    this.code = code
  }
}

/** 后端把错误码写进 message，形如「错误代码: 1401, 错误信息: ...」。 */
function errorCode(message: string): number | null {
  const matched = /错误代码:\s*(\d+)/.exec(message)
  return matched ? Number(matched[1]) : null
}

export async function graphqlRequest<TData, TVariables extends Record<string, unknown> = Record<string, never>>(
  query: string,
  variables?: TVariables,
  token?: string | null,
): Promise<TData> {
  const headers: Record<string, string> = { "Content-Type": "application/json" }
  // 后端从配置里的 token 名读取（config.json 的 secret.jwt.token_name，默认 Authentication）
  if (token) {
    headers.Authentication = token
  }

  const response = await fetch(ENDPOINT, {
    method: "POST",
    headers,
    body: JSON.stringify({ query, variables }),
  })

  if (!response.ok) {
    throw new ApiError(`GraphQL 请求失败：HTTP ${response.status}`, response.status)
  }

  const payload = (await response.json()) as GraphQLResponse<TData>
  if (payload.errors?.length) {
    const message = payload.errors.map((error) => error.message).join("; ")
    const code = errorCode(payload.errors[0]!.message)
    if (code === CODE_UNAUTHORIZED) {
      clearToken()
    }
    throw new ApiError(message, code)
  }
  if (!payload.data) {
    throw new ApiError("GraphQL 响应缺少 data")
  }
  return payload.data
}
