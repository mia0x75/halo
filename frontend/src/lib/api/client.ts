/**
 * GraphQL 传输层。只做一件事：把查询发到后端并把 data 取出来。
 * 不引入查询缓存或状态库，失败时由调用方回落到假数据。
 */

const ENDPOINT = "/api/query"

export interface GraphQLError {
  message: string
}

interface GraphQLResponse<TData> {
  data?: TData
  errors?: GraphQLError[]
}

export class ApiError extends Error {
  constructor(message: string) {
    super(message)
    this.name = "ApiError"
  }
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
    throw new ApiError(`GraphQL 请求失败：HTTP ${response.status}`)
  }

  const payload = (await response.json()) as GraphQLResponse<TData>
  if (payload.errors?.length) {
    throw new ApiError(payload.errors.map((error) => error.message).join("; "))
  }
  if (!payload.data) {
    throw new ApiError("GraphQL 响应缺少 data")
  }
  return payload.data
}
