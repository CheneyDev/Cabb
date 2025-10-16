import { NextRequest } from 'next/server'
import { getAPIBase } from '@/lib/config'

export async function GET(req: NextRequest) {
  const base = getAPIBase()
  const qs = req.nextUrl.searchParams
  const url = `${base}/admin/mappings/repo-project?${qs.toString()}`
  const r = await fetch(url, { headers: { 'Accept': 'application/json' }, cache: 'no-store' })
  const body = await r.text()
  return new Response(body, { status: r.status, headers: { 'content-type': r.headers.get('content-type') || 'application/json' } })
}

export async function POST(req: NextRequest) {
  const base = getAPIBase()
  const url = `${base}/admin/mappings/repo-project`
  const body = await req.text()
  const r = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
    body,
  })
  const text = await r.text()
  return new Response(text, { status: r.status, headers: { 'content-type': r.headers.get('content-type') || 'application/json' } })
}

