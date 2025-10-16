import type { NextRequest } from 'next/server'

import { proxyAPI } from '@/lib/server/proxy'

export async function GET(req: NextRequest) {
  return proxyAPI(req, { path: '/admin/access/users', searchParams: req.nextUrl.searchParams })
}

export async function POST(req: NextRequest) {
  const body = await req.text()
  return proxyAPI(req, {
    path: '/admin/access/users',
    method: 'POST',
    body,
    headers: { 'Content-Type': 'application/json' },
  })
}
