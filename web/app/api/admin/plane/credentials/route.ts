import { NextResponse, type NextRequest } from 'next/server'

import { proxyAPI } from '@/lib/server/proxy'

export async function GET(req: NextRequest) {
  return proxyAPI(req, { path: '/admin/plane/credentials' })
}

export async function POST(req: NextRequest) {
  return proxyAPI(req, { path: '/admin/plane/credentials', method: 'POST' })
}
