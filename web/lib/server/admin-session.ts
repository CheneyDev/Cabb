import { cookies, headers } from 'next/headers'

import { getAPIBase } from '@/lib/config'

export type AdminUser = {
  id: string
  email: string
  display_name: string
  role: string
  active: boolean
}

export type AdminSessionResult = {
  user: AdminUser | null
  status: number | null
  message?: string
  expiresAt?: string | null
}

async function buildHeadersFromRequest() {
  const headerInit = new Headers({ Accept: 'application/json' })
  const cookieStore = await cookies()
  const headerStore = await headers()
  const cookieHeader = cookieStore
    .getAll()
    .map(cookie => `${cookie.name}=${cookie.value}`)
    .join('; ')
  if (cookieHeader) {
    headerInit.set('Cookie', cookieHeader)
  }
  const forwardedFor = headerStore.get('x-forwarded-for')
  if (forwardedFor) {
    headerInit.set('X-Forwarded-For', forwardedFor)
  }
  const forwardedProto = headerStore.get('x-forwarded-proto')
  if (forwardedProto) {
    headerInit.set('X-Forwarded-Proto', forwardedProto)
  }
  const userAgent = headerStore.get('user-agent')
  if (userAgent) {
    headerInit.set('User-Agent', userAgent)
  }
  return headerInit
}

function parseErrorMessage(text: string | null) {
  if (!text) return undefined
  try {
    const json = JSON.parse(text)
    if (json?.error?.message) return String(json.error.message)
    if (json?.message) return String(json.message)
  } catch (err) {
    return text
  }
  return undefined
}

export async function fetchAdminSession(): Promise<AdminSessionResult> {
  const base = getAPIBase()
  const headersInit = await buildHeadersFromRequest()
  try {
    const res = await fetch(`${base}/admin/auth/me`, {
      headers: headersInit,
      cache: 'no-store',
    })
    const text = await res.text()
    if (!res.ok) {
      const message = parseErrorMessage(text) ?? `加载管理员信息失败（${res.status}）`
      return { user: null, status: res.status, message }
    }
    let data: unknown
    if (text) {
      try {
        data = JSON.parse(text)
      } catch (err) {
        return { user: null, status: res.status, message: '解析管理员信息失败' }
      }
    }
    const json = data as { user?: AdminUser; session?: { expires_at?: string | null } } | undefined
    const user = json?.user ?? null
    if (!user) {
      return { user: null, status: res.status, message: '未找到管理员信息，请重新登录' }
    }
    const expiresAt = json?.session?.expires_at ?? null
    return {
      user,
      status: res.status,
      expiresAt: typeof expiresAt === 'string' ? expiresAt : null,
    }
  } catch (err) {
    return {
      user: null,
      status: null,
      message: err instanceof Error ? err.message : '管理员鉴权失败，请稍后重试',
    }
  }
}
