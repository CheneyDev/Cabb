'use client'

import { useCallback, useEffect, useState } from 'react'

import { Alert } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type Credential = {
  id: string
  plane_workspace_id: string
  workspace_slug: string
  workspace_name?: string | null
  kind: string
  token_masked: string
  created_at: string
  updated_at: string
}

type Feedback = { kind: 'success' | 'error'; message: string }

export default function PlaneCredentialsPage() {
  const [items, setItems] = useState<Credential[]>([])
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [deleting, setDeleting] = useState<string | null>(null)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [feedback, setFeedback] = useState<Feedback | null>(null)

  const [form, setForm] = useState({
    plane_workspace_id: '',
    workspace_slug: '',
    token: '',
  })

  const load = useCallback(async () => {
    setLoading(true)
    setFetchError(null)
    try {
      const res = await fetch('/api/admin/plane/credentials', { cache: 'no-store' })
      if (!res.ok) {
        const message = await readErrorMessage(res, `查询失败（${res.status}）`)
        throw new Error(message)
      }
      const json = await res.json()
      setItems(Array.isArray(json.items) ? json.items : [])
    } catch (err) {
      setFetchError(err instanceof Error ? err.message : '加载凭据失败')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setSubmitting(true)
    setFeedback(null)

    try {
      const payload = {
        plane_workspace_id: form.plane_workspace_id.trim(),
        workspace_slug: form.workspace_slug.trim(),
        token: form.token.trim(),
      }

      if (!payload.plane_workspace_id || !payload.workspace_slug || !payload.token) {
        throw new Error('请填写所有必填字段')
      }

      const res = await fetch('/api/admin/plane/credentials', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })

      if (!res.ok) {
        const message = await readErrorMessage(res, `保存失败（${res.status}）`)
        throw new Error(message)
      }

      const result = await res.json()
      setFeedback({
        kind: 'success',
        message: result.message || '凭据已保存，Plane 出站功能已启用',
      })
      setForm({ plane_workspace_id: '', workspace_slug: '', token: '' })
      await load()
    } catch (err) {
      setFeedback({
        kind: 'error',
        message: err instanceof Error ? err.message : '保存失败，请稍后再试',
      })
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDelete(id: string, workspace_slug: string) {
    if (!window.confirm(`确定要删除 Workspace "${workspace_slug}" 的凭据吗？删除后将无法调用 Plane API。`)) {
      return
    }

    setDeleting(id)
    setFeedback(null)

    try {
      const res = await fetch(`/api/admin/plane/credentials/${id}`, {
        method: 'DELETE',
      })

      if (!res.ok) {
        const message = await readErrorMessage(res, `删除失败（${res.status}）`)
        throw new Error(message)
      }

      const result = await res.json()
      setFeedback({
        kind: 'success',
        message: result.message || '凭据已删除',
      })
      await load()
    } catch (err) {
      setFeedback({
        kind: 'error',
        message: err instanceof Error ? err.message : '删除失败',
      })
    } finally {
      setDeleting(null)
    }
  }

  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <div className="flex flex-col gap-2 md:flex-row md:items-start md:justify-between">
            <div className="space-y-1">
              <CardTitle>Plane Service Token 凭据管理</CardTitle>
              <CardDescription>
                配置 Workspace 的 Service Token 后，可启用向 Plane 写回功能（CNB/飞书 → Plane）。
                <br />
                删除凭据将禁用该 Workspace 的出站调用。
              </CardDescription>
            </div>
            <Badge variant="info">Credentials</Badge>
          </div>
        </CardHeader>
        <CardContent>
          <form className="grid gap-4 md:grid-cols-3" onSubmit={handleSubmit}>
            <div className="grid gap-2">
              <Label htmlFor="plane_workspace_id">Plane Workspace ID（必填）</Label>
              <Input
                id="plane_workspace_id"
                required
                placeholder="uuid"
                value={form.plane_workspace_id}
                onChange={e => setForm(prev => ({ ...prev, plane_workspace_id: e.target.value }))}
              />
              <span className="text-xs text-muted-foreground">Workspace 的 UUID</span>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="workspace_slug">Workspace Slug（必填）</Label>
              <Input
                id="workspace_slug"
                required
                placeholder="my-workspace"
                value={form.workspace_slug}
                onChange={e => setForm(prev => ({ ...prev, workspace_slug: e.target.value }))}
              />
              <span className="text-xs text-muted-foreground">Workspace 的 URL 标识</span>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="token">Service Token（必填）</Label>
              <Input
                id="token"
                type="password"
                required
                placeholder="plane_token_xxx"
                value={form.token}
                onChange={e => setForm(prev => ({ ...prev, token: e.target.value }))}
              />
              <span className="text-xs text-muted-foreground">从 Plane 获取的 API Token</span>
            </div>
            {feedback && (
              <div className="md:col-span-3">
                <Alert variant={feedback.kind === 'success' ? 'success' : 'destructive'}>{feedback.message}</Alert>
              </div>
            )}
            <div className="md:col-span-3 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
              <span className="text-sm text-muted-foreground">
                提交后会验证 Token 有效性，成功后即可启用 Plane 出站功能。
              </span>
              <div className="flex items-center gap-2">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => setForm({ plane_workspace_id: '', workspace_slug: '', token: '' })}
                >
                  清空
                </Button>
                <Button type="submit" disabled={submitting}>
                  {submitting ? '验证并保存中…' : '保存凭据'}
                </Button>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <div className="space-y-1">
              <CardTitle>已配置的凭据</CardTitle>
              <CardDescription>当前所有 Workspace 的 Service Token 配置。</CardDescription>
            </div>
            <div className="flex items-center gap-2">
              <Badge variant="success">已配置 {items.length}</Badge>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {fetchError && <Alert variant="destructive">{fetchError}</Alert>}
          <div className="overflow-x-auto">
            <Table className="min-w-[600px]">
              <TableHeader>
                <TableRow className="bg-transparent">
                  <TableHead>Workspace</TableHead>
                  <TableHead>Workspace ID</TableHead>
                  <TableHead>Token（已脱敏）</TableHead>
                  <TableHead>更新时间</TableHead>
                  <TableHead className="text-right">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {items.map(item => (
                  <TableRow key={item.id}>
                    <TableCell>
                      <div className="flex flex-col gap-1">
                        <span className="font-medium text-foreground">
                          {item.workspace_name || item.workspace_slug}
                        </span>
                        {item.workspace_name && (
                          <span className="text-xs text-muted-foreground">Slug: {item.workspace_slug}</span>
                        )}
                      </div>
                    </TableCell>
                    <TableCell>
                      <span className="text-xs font-mono text-muted-foreground">{item.plane_workspace_id}</span>
                    </TableCell>
                    <TableCell>
                      <span className="font-mono text-sm text-muted-foreground">{item.token_masked}</span>
                    </TableCell>
                    <TableCell>
                      <span className="text-sm text-muted-foreground">
                        {new Date(item.updated_at).toLocaleString('zh-CN')}
                      </span>
                    </TableCell>
                    <TableCell className="text-right">
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleDelete(item.id, item.workspace_slug)}
                        disabled={deleting === item.id}
                      >
                        {deleting === item.id ? '删除中…' : '删除'}
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
                {items.length === 0 && !loading && (
                  <TableRow>
                    <TableCell colSpan={5} className="text-center text-sm text-muted-foreground">
                      暂无凭据配置，请添加 Service Token 以启用 Plane 出站功能。
                    </TableCell>
                  </TableRow>
                )}
                {loading && (
                  <TableRow>
                    <TableCell colSpan={5} className="text-center text-sm text-muted-foreground">
                      正在加载凭据数据…
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

async function readErrorMessage(res: Response, fallback: string) {
  const text = await res.text()
  const ct = res.headers.get('content-type') || ''
  let msg = res.statusText || ''
  if (ct.includes('application/json')) {
    try {
      const json = JSON.parse(text)
      msg = json?.error?.message || json?.message || msg
    } catch {}
  }
  if (!msg) {
    const snippet = text.trim().replace(/\s+/g, ' ')
    if (snippet && !snippet.startsWith('<')) {
      msg = snippet.slice(0, 200)
    }
  }
  return msg || fallback
}
