import type { NextRequest } from 'next/server'

import { proxyAPI } from '@/lib/server/proxy'

export async function POST(req: NextRequest) {
  return proxyAPI(req, { path: '/admin/auth/logout', method: 'POST' })
}
