import { NextResponse, type NextRequest } from 'next/server'

import { proxyAPI } from '@/lib/server/proxy'

export async function GET(req: NextRequest) {
  return proxyAPI(req, { path: '/admin/mappings/repo-project', searchParams: req.nextUrl.searchParams })
}

export async function POST(req: NextRequest) {
  const body = await req.text()
  return proxyAPI(req, {
    path: '/admin/mappings/repo-project',
    method: 'POST',
    body,
    headers: { 'Content-Type': 'application/json' },
  })
}

export async function PATCH(req: NextRequest) {
  const body = await req.text()
  return proxyAPI(req, {
    path: '/admin/mappings/repo-project',
    method: 'POST',
    body,
    headers: { 'Content-Type': 'application/json' },
  })
}

type DeletePayload = {
  cnb_repo_id?: string
  plane_workspace_id?: string
  plane_project_id?: string
  issue_open_state_id?: string | null
  issue_closed_state_id?: string | null
  label_selector?: string | null
  sync_direction?: string | null
}

export async function DELETE(req: NextRequest) {
  let payload: DeletePayload | null = null
  try {
    payload = (await req.json()) as DeletePayload
  } catch (error) {
    return NextResponse.json(
      { error: { code: 'invalid_json', message: '解析请求失败' } },
      { status: 400 },
    )
  }

  if (!payload?.cnb_repo_id || !payload?.plane_workspace_id || !payload?.plane_project_id) {
    return NextResponse.json(
      { error: { code: 'missing_fields', message: '缺少 cnb_repo_id / plane_workspace_id / plane_project_id' } },
      { status: 400 },
    )
  }

  const body = JSON.stringify({
    cnb_repo_id: payload.cnb_repo_id,
    plane_workspace_id: payload.plane_workspace_id,
    plane_project_id: payload.plane_project_id,
    issue_open_state_id: payload.issue_open_state_id ?? '',
    issue_closed_state_id: payload.issue_closed_state_id ?? '',
    label_selector: payload.label_selector ?? '',
    sync_direction: payload.sync_direction ?? '',
    active: false,
  })

  return proxyAPI(req, {
    path: '/admin/mappings/repo-project',
    method: 'POST',
    body,
    headers: { 'Content-Type': 'application/json' },
  })
}

