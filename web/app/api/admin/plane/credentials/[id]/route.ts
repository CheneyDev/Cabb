import { NextResponse, type NextRequest } from 'next/server'

import { proxyAPI } from '@/lib/server/proxy'

export async function DELETE(req: NextRequest, { params }: { params: Promise<{ id: string }> }) {
  const { id } = await params
  return proxyAPI(req, { path: `/admin/plane/credentials/${id}`, method: 'DELETE' })
}
