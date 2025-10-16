import type { NextRequest } from 'next/server'

import { proxyAPI } from '@/lib/server/proxy'

export async function POST(req: NextRequest, context: { params: Promise<Record<string, string>> }) {
  const params = await context.params
  const id = params.id
  const body = await req.text()
  return proxyAPI(req, {
    path: `/admin/access/users/${id}/reset-password`,
    method: 'POST',
    body,
    headers: { 'Content-Type': 'application/json' },
  })
}
