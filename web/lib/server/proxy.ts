import type { NextRequest } from 'next/server'

import { getAPIBase } from '@/lib/config'

function buildURL(path: string, search?: URLSearchParams | string) {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`
  if (!search) return normalizedPath
  const suffix = typeof search === 'string' ? search.replace(/^\?/, '') : search.toString()
  if (!suffix) return normalizedPath
  return `${normalizedPath}?${suffix}`
}

function parseErrorBody(text: string) {
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

export async function proxyAPI(
  req: NextRequest,
  init: {
    path: string
    method?: string
    body?: string | null
    searchParams?: URLSearchParams | string
    headers?: HeadersInit
  },
) {
  const base = getAPIBase()
  const url = `${base}${buildURL(init.path, init.searchParams)}`
  const headers = new Headers(init.headers ?? {})
  if (!headers.has('Accept')) {
    headers.set('Accept', 'application/json')
  }
  const cookie = req.headers.get('cookie')
  if (cookie) {
    headers.set('Cookie', cookie)
  }
  const forwardedFor = req.headers.get('x-forwarded-for')
  if (forwardedFor) {
    headers.set('X-Forwarded-For', forwardedFor)
  }
  const forwardedProto = req.headers.get('x-forwarded-proto')
  if (forwardedProto) {
    headers.set('X-Forwarded-Proto', forwardedProto)
  }
  const userAgent = req.headers.get('user-agent')
  if (userAgent) {
    headers.set('User-Agent', userAgent)
  }

  let res: Response
  try {
    res = await fetch(url, {
      method: init.method ?? req.method,
      headers,
      body: init.body ?? undefined,
      cache: 'no-store',
    })
  } catch (error) {
    const reason = error instanceof Error ? error.message : 'unknown error'
    const message = '无法连接后端服务，请稍后再试'
    return Response.json(
      {
        error: {
          code: 'upstream_unreachable',
          message,
          details: { reason },
        },
      },
      {
        status: 502,
        headers: {
          'x-proxy-error': message,
        },
      },
    )
  }

  const text = await res.text()
  const proxyRes = new Response(text, {
    status: res.status,
    headers: {
      'content-type': res.headers.get('content-type') ?? 'application/json',
    },
  })
  const requestID = res.headers.get('x-request-id')
  if (requestID) {
    proxyRes.headers.set('x-request-id', requestID)
  }
  const setCookies = res.headers.getSetCookie?.()
  if (setCookies && setCookies.length > 0) {
    for (const cookieValue of setCookies) {
      proxyRes.headers.append('set-cookie', cookieValue)
    }
  }
  if (!res.ok && text) {
    const message = parseErrorBody(text)
    if (message) {
      proxyRes.headers.set('x-proxy-error', message)
    }
  }
  return proxyRes
}
