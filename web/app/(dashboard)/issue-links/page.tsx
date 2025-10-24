'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'

import { Alert } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const initialForm = {
  plane_issue_id: '',
  cnb_repo_id: '',
  cnb_issue_id: '',
}

type FormState = typeof initialForm

type IssueLink = {
  plane_issue_id: string
  cnb_repo_id?: string | null
  cnb_issue_id?: string | null
  linked_at?: string
  created_at?: string
  updated_at?: string
}

type Feedback = { kind: 'success' | 'error'; message: string }

type Filters = {
  plane_issue_id: string
  cnb_repo_id: string
  cnb_issue_id: string
}

function makeKey(item: IssueLink) {
  return `${item.plane_issue_id}::${item.cnb_repo_id ?? ''}::${item.cnb_issue_id ?? ''}`
}

export default function IssueLinksPage() {
  const [form, setForm] = useState<FormState>(initialForm)
  const [filters, setFilters] = useState<Filters>({ plane_issue_id: '', cnb_repo_id: '', cnb_issue_id: '' })
  const [items, setItems] = useState<IssueLink[]>([])
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [feedback, setFeedback] = useState<Feedback | null>(null)
  const [deletingKey, setDeletingKey] = useState<string | null>(null)

  const querySuffix = useMemo(() => {
    const params = new URLSearchParams()
    params.set('limit', '100')
    const planeIssue = filters.plane_issue_id.trim()
    const repo = filters.cnb_repo_id.trim()
    const issue = filters.cnb_issue_id.trim()
    if (planeIssue) params.set('plane_issue_id', planeIssue)
    if (repo) params.set('cnb_repo_id', repo)
    if (issue) params.set('cnb_issue_id', issue)
    const query = params.toString()
    return query ? `?${query}` : ''
  }, [filters])

  const load = useCallback(async () => {
    setLoading(true)
    setFetchError(null)
    try {
      const res = await fetch(`/api/admin/links/issue-links${querySuffix}`, { cache: 'no-store' })
      if (!res.ok) {
        const message = await readErrorMessage(res, `查询失败（${res.status}）`)
        throw new Error(message)
      }
      const data = await res.json()
      setItems(Array.isArray(data.items) ? (data.items as IssueLink[]) : [])
    } catch (err) {
      setFetchError(err instanceof Error ? err.message : '加载失败')
    } finally {
      setLoading(false)
    }
  }, [querySuffix])

  useEffect(() => {
    load()
  }, [load])

  const total = items.length

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setFeedback(null)
    setSubmitting(true)
    try {
      const payload = {
        plane_issue_id: form.plane_issue_id.trim(),
        cnb_repo_id: form.cnb_repo_id.trim(),
        cnb_issue_id: form.cnb_issue_id.trim(),
      }
      if (!payload.plane_issue_id || !payload.cnb_repo_id || !payload.cnb_issue_id) {
        throw new Error('请填写 Plane Issue ID、CNB 仓库与 Issue IID。')
      }
      const res = await fetch('/api/admin/links/issue-links', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) {
        const message = await readErrorMessage(res, `保存失败（${res.status}）`)
        throw new Error(message)
      }
      setFeedback({ kind: 'success', message: '映射已保存。' })
      setForm(initialForm)
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '保存失败，请稍后再试。' })
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDelete(item: IssueLink) {
    const key = makeKey(item)
    setDeletingKey(key)
    try {
      const payload = {
        plane_issue_id: item.plane_issue_id,
        cnb_repo_id: item.cnb_repo_id ?? '',
        cnb_issue_id: item.cnb_issue_id ?? '',
      }
      const res = await fetch('/api/admin/links/issue-links', {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) {
        const message = await readErrorMessage(res, `删除失败（${res.status}）`)
        throw new Error(message)
      }
      setFeedback({ kind: 'success', message: '映射已删除。' })
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '删除失败，请稍后再试。' })
    } finally {
      setDeletingKey(null)
    }
  }

  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle>Plane Issue ↔ CNB Issue 映射</CardTitle>
          <CardDescription>用于追踪 Plane 工作项与 CNB Issue 的关联关系，便于溯源与手动修复。</CardDescription>
        </CardHeader>
        <CardContent className="grid gap-6">
          <form className="grid gap-4 md:grid-cols-3" onSubmit={handleSubmit}>
            <div className="grid gap-2">
              <Label htmlFor="plane_issue_id">Plane Issue ID</Label>
              <Input
                id="plane_issue_id"
                placeholder="uuid"
                value={form.plane_issue_id}
                onChange={e => setForm(prev => ({ ...prev, plane_issue_id: e.target.value }))}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="cnb_repo_id">CNB 仓库</Label>
              <Input
                id="cnb_repo_id"
                placeholder="group/repo"
                value={form.cnb_repo_id}
                onChange={e => setForm(prev => ({ ...prev, cnb_repo_id: e.target.value }))}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="cnb_issue_id">CNB Issue IID</Label>
              <Input
                id="cnb_issue_id"
                placeholder="数字或自定义标识"
                value={form.cnb_issue_id}
                onChange={e => setForm(prev => ({ ...prev, cnb_issue_id: e.target.value }))}
              />
            </div>
            <div className="md:col-span-3 flex justify-end">
              <Button type="submit" disabled={submitting}>
                {submitting ? '保存中…' : '保存映射'}
              </Button>
            </div>
          </form>
          <div className="grid gap-4 md:grid-cols-3">
            <div className="grid gap-2">
              <Label htmlFor="filter_plane_issue">按 Plane Issue ID 筛选</Label>
              <Input
                id="filter_plane_issue"
                placeholder="uuid"
                value={filters.plane_issue_id}
                onChange={e => setFilters(prev => ({ ...prev, plane_issue_id: e.target.value }))}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="filter_repo">按仓库筛选</Label>
              <Input
                id="filter_repo"
                placeholder="group/repo"
                value={filters.cnb_repo_id}
                onChange={e => setFilters(prev => ({ ...prev, cnb_repo_id: e.target.value }))}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="filter_issue">按 Issue IID 筛选</Label>
              <Input
                id="filter_issue"
                placeholder="iid"
                value={filters.cnb_issue_id}
                onChange={e => setFilters(prev => ({ ...prev, cnb_issue_id: e.target.value }))}
              />
            </div>
          </div>
          <div className="flex items-center justify-between text-sm text-muted-foreground">
            <span>共 {total} 条记录</span>
            <div className="flex items-center gap-2">
              {fetchError ? <Badge variant="destructive">加载失败</Badge> : null}
              <Button type="button" variant="outline" onClick={load} disabled={loading}>
                {loading ? '刷新中…' : '刷新'}
              </Button>
            </div>
          </div>
          {fetchError ? <Alert variant="destructive">{fetchError}</Alert> : null}
        </CardContent>
        <CardFooter>
          {feedback ? <Alert variant={feedback.kind === 'success' ? 'default' : 'destructive'}>{feedback.message}</Alert> : null}
        </CardFooter>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>映射列表</CardTitle>
          <CardDescription>最新更新时间优先排序，可手动解除错误关联。</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Plane Issue ID</TableHead>
                  <TableHead>CNB 仓库</TableHead>
                  <TableHead>CNB Issue IID</TableHead>
                  <TableHead>Linked At</TableHead>
                  <TableHead>更新时间</TableHead>
                  <TableHead className="text-right">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {items.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-sm text-muted-foreground">
                      暂无数据，请调整筛选条件。
                    </TableCell>
                  </TableRow>
                ) : (
                  items.map(item => (
                    <TableRow key={makeKey(item)}>
                      <TableCell className="font-mono text-xs">{item.plane_issue_id}</TableCell>
                      <TableCell className="font-mono text-xs">{item.cnb_repo_id || '—'}</TableCell>
                      <TableCell className="font-mono text-xs">{item.cnb_issue_id || '—'}</TableCell>
                      <TableCell className="text-xs text-muted-foreground">
                        {item.linked_at ? new Date(item.linked_at).toLocaleString() : '—'}
                      </TableCell>
                      <TableCell className="text-xs text-muted-foreground">
                        {item.updated_at ? new Date(item.updated_at).toLocaleString() : '—'}
                      </TableCell>
                      <TableCell className="text-right">
                        <Button
                          type="button"
                          variant="destructive"
                          size="sm"
                          onClick={() => handleDelete(item)}
                          disabled={deletingKey === makeKey(item)}
                        >
                          {deletingKey === makeKey(item) ? '删除中…' : '删除'}
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))
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
