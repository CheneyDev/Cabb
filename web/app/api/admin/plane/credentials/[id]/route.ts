import { NextResponse, type NextRequest } from 'next/server'

import { proxyAPI } from '@/lib/server/proxy'

export async function DELETE(req: NextRequest, { params }: { params: { id: string } }) {
  return proxyAPI(req, { path: `/admin/plane/credentials/${params.id}`, method: 'DELETE' })
}
