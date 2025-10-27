import { NextRequest } from 'next/server'

import { proxyAPI } from '@/lib/server/proxy'

export async function GET(req: NextRequest) {
  return proxyAPI(req, { path: '/admin/links/issues', searchParams: req.nextUrl.searchParams })
}

export async function POST(req: NextRequest) {
  const body = await req.text()
  return proxyAPI(req, {
    path: '/admin/links/issues',
    method: 'POST',
    body,
    headers: { 'Content-Type': 'application/json' },
  })
}

export async function DELETE(req: NextRequest) {
  const body = await req.text()
  return proxyAPI(req, {
    path: '/admin/links/issues',
    method: 'DELETE',
    body,
    headers: { 'Content-Type': 'application/json' },
  })
}
