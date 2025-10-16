import type { NextRequest } from 'next/server'

import { proxyAPI } from '@/lib/server/proxy'

export async function POST(req: NextRequest) {
  const body = await req.text()
  return proxyAPI(req, {
    path: '/admin/auth/login',
    method: 'POST',
    body,
    headers: { 'Content-Type': 'application/json' },
  })
}
