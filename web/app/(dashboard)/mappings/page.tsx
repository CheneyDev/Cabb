'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'

import { Alert } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const directionOptions = [
  { value: 'cnb_to_plane', label: '仅 CNB → Plane' },
  { value: 'bidirectional', label: '双向同步' },
]

type Mapping = {
  plane_project_id: string
  plane_workspace_id: string
  cnb_repo_id: string
  issue_open_state_id?: string | null
  issue_closed_state_id?: string | null
  active: boolean
  sync_direction?: string | null
  label_selector?: string | null
}

type Feedback = { kind: 'success' | 'error'; message: string }

type ActiveFilter = 'all' | 'active' | 'inactive'

const initialForm = {
  cnb_repo_id: '',
  plane_workspace_id: '',
  plane_project_id: '',
  issue_open_state_id: '',
  issue_closed_state_id: '',
  sync_direction: 'cnb_to_plane',
  label_selector: '',
  active: true,
}

type FormState = typeof initialForm

export default function MappingsPage() {
  const [items, setItems] = useState<Mapping[]>([])
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [feedback, setFeedback] = useState<Feedback | null>(null)
  const [filterProject, setFilterProject] = useState('')
  const [filterRepo, setFilterRepo] = useState('')
  const [activeFilter, setActiveFilter] = useState<ActiveFilter>('all')
  const [form, setForm] = useState<FormState>(initialForm)

  const querySuffix = useMemo(() => {
    const params = new URLSearchParams()
    const project = filterProject.trim()
    const repo = filterRepo.trim()
    if (project) params.set('plane_project_id', project)
    if (repo) params.set('cnb_repo_id', repo)
    if (activeFilter === 'active') params.set('active', 'true')
    if (activeFilter === 'inactive') params.set('active', 'false')
    const query = params.toString()
    return query ? `?${query}` : ''
  }, [filterProject, filterRepo, activeFilter])

  const load = useCallback(async () => {
    setLoading(true)
    setFetchError(null)
    try {
      const res = await fetch(`/api/admin/mappings/repo-project${querySuffix}`, { cache: 'no-store' })
      if (!res.ok) {
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
          msg = snippet.startsWith('<') ? '' : snippet.slice(0, 200)
        }
        throw new Error(msg || `查询失败（${res.status}）`)
      }
      const json = await res.json()
      setItems(Array.isArray(json.items) ? json.items : [])
    } catch (err) {
      setFetchError(err instanceof Error ? err.message : '加载映射失败')
    } finally {
      setLoading(false)
    }
  }, [querySuffix])

  useEffect(() => {
    load()
  }, [load])

  const activeCount = items.filter(item => item.active).length
  const inactiveCount = items.length - activeCount

  async function submit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setSubmitting(true)
    setFeedback(null)
    try {
      const payload: Record<string, unknown> = {
        cnb_repo_id: form.cnb_repo_id.trim(),
        plane_workspace_id: form.plane_workspace_id.trim(),
        plane_project_id: form.plane_project_id.trim(),
        active: form.active,
      }
      if (form.sync_direction) payload.sync_direction = form.sync_direction
      if (form.label_selector) payload.label_selector = form.label_selector.trim()
      if (form.issue_open_state_id) payload.issue_open_state_id = form.issue_open_state_id.trim()
      if (form.issue_closed_state_id) payload.issue_closed_state_id = form.issue_closed_state_id.trim()

      const res = await fetch('/api/admin/mappings/repo-project', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) {
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
          msg = snippet.startsWith('<') ? '' : snippet.slice(0, 200)
        }
        throw new Error(msg || `保存失败（${res.status}）`)
      }
      setFeedback({ kind: 'success', message: '映射已保存并刷新。' })
      setForm(prev => ({ ...prev, cnb_repo_id: '' }))
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '保存失败，请稍后再试。' })
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle>新增或更新 Repo ↔ Project 映射</CardTitle>
          <CardDescription>填写必填字段后提交即可覆盖同名仓库的最新策略。</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="grid gap-4 md:grid-cols-2" onSubmit={submit}>
            <div className="grid gap-2">
              <Label htmlFor="cnb_repo_id">CNB 仓库（org/repo）</Label>
              <Input
                id="cnb_repo_id"
                required
                placeholder="1024hub/plane-integration"
                value={form.cnb_repo_id}
                onChange={event => setForm(prev => ({ ...prev, cnb_repo_id: event.target.value }))}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="plane_workspace_id">Plane Workspace ID</Label>
              <Input
                id="plane_workspace_id"
                required
                placeholder="workspace uuid"
                value={form.plane_workspace_id}
                onChange={event => setForm(prev => ({ ...prev, plane_workspace_id: event.target.value }))}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="plane_project_id">Plane Project ID</Label>
              <Input
                id="plane_project_id"
                required
                placeholder="project uuid"
                value={form.plane_project_id}
                onChange={event => setForm(prev => ({ ...prev, plane_project_id: event.target.value }))}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="sync_direction">同步方向</Label>
              <Select
                id="sync_direction"
                value={form.sync_direction}
                onChange={event => setForm(prev => ({ ...prev, sync_direction: event.target.value }))}
              >
                {directionOptions.map(option => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </Select>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="issue_open_state_id">Issue 打开状态 ID（可选）</Label>
              <Input
                id="issue_open_state_id"
                placeholder="用于新建 Issue 对应的 Plane 状态"
                value={form.issue_open_state_id}
                onChange={event => setForm(prev => ({ ...prev, issue_open_state_id: event.target.value }))}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="issue_closed_state_id">Issue 关闭状态 ID（可选）</Label>
              <Input
                id="issue_closed_state_id"
                placeholder="用于关闭 Issue 对应的 Plane 状态"
                value={form.issue_closed_state_id}
                onChange={event => setForm(prev => ({ ...prev, issue_closed_state_id: event.target.value }))}
              />
            </div>
            <div className="grid gap-2 md:col-span-2">
              <Label htmlFor="label_selector">标签选择器（逗号分隔，可选）</Label>
              <Input
                id="label_selector"
                placeholder="* 或 逗号分隔的标签名"
                value={form.label_selector}
                onChange={event => setForm(prev => ({ ...prev, label_selector: event.target.value }))}
              />
            </div>
            <div className="flex items-center gap-3 md:col-span-2">
              <Switch
                checked={form.active}
                onChange={event => setForm(prev => ({ ...prev, active: event.target.checked }))}
                aria-label="是否启用映射"
              />
              <span className="text-sm text-muted-foreground">启用映射（停用后同步将跳过该仓库）</span>
            </div>
            {feedback && (
              <div className="md:col-span-2">
                <Alert variant={feedback.kind === 'success' ? 'success' : 'destructive'}>{feedback.message}</Alert>
              </div>
            )}
            <div className="md:col-span-2 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
              <span className="text-sm text-muted-foreground">保存后列表将自动刷新。</span>
              <div className="flex items-center gap-2">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => {
                    setForm(initialForm)
                    setFeedback(null)
                  }}
                >
                  重置表单
                </Button>
                <Button type="submit" disabled={submitting}>
                  {submitting ? '保存中…' : '保存映射'}
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
              <CardTitle>现有映射概览</CardTitle>
              <CardDescription>支持按 Plane 项目或 CNB 仓库筛选，并快速查看同步状态。</CardDescription>
            </div>
            <div className="flex flex-wrap items-center gap-2">
              <Badge variant="success">启用 {activeCount}</Badge>
              <Badge variant="muted">停用 {inactiveCount}</Badge>
              <Badge variant="outline">总计 {items.length}</Badge>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-3 md:grid-cols-[2fr_2fr_auto_auto] md:items-end">
            <div className="grid gap-2">
              <Label htmlFor="filterProject">按 Plane Project ID 过滤</Label>
              <Input id="filterProject" value={filterProject} onChange={event => setFilterProject(event.target.value)} />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="filterRepo">按 CNB 仓库过滤</Label>
              <Input id="filterRepo" value={filterRepo} onChange={event => setFilterRepo(event.target.value)} />
            </div>
            <div className="flex flex-col gap-2">
              <span className="text-sm font-medium text-muted-foreground">显示状态</span>
              <div className="flex flex-wrap items-center gap-2">
                {(['all', 'active', 'inactive'] as ActiveFilter[]).map(value => (
                  <Button
                    key={value}
                    variant={activeFilter === value ? 'secondary' : 'ghost'}
                    onClick={() => setActiveFilter(value)}
                  >
                    {value === 'all' ? '全部' : value === 'active' ? '仅启用' : '仅停用'}
                  </Button>
                ))}
              </div>
            </div>
            <div className="flex flex-wrap gap-2 md:justify-end">
              <Button type="button" variant="outline" onClick={() => load()} disabled={loading}>
                {loading ? '加载中…' : '刷新'}
              </Button>
              <Button
                type="button"
                variant="ghost"
                onClick={() => {
                  setFilterProject('')
                  setFilterRepo('')
                  setActiveFilter('all')
                }}
              >
                清除筛选
              </Button>
            </div>
          </div>
          {fetchError && <Alert variant="destructive">{fetchError}</Alert>}
          <div className="overflow-x-auto">
            <Table className="min-w-[720px]">
              <TableHeader>
                <TableRow className="bg-transparent">
                  <TableHead>CNB 仓库</TableHead>
                  <TableHead>Plane Workspace / Project</TableHead>
                  <TableHead>同步方向</TableHead>
                  <TableHead>标签选择器</TableHead>
                  <TableHead>Issue 状态映射</TableHead>
                  <TableHead>状态</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {items.map(item => (
                  <TableRow key={`${item.cnb_repo_id}-${item.plane_project_id}`}>
                    <TableCell>
                      <span className="font-medium text-foreground">{item.cnb_repo_id}</span>
                    </TableCell>
                    <TableCell>
                      <div className="flex flex-col gap-1">
                        <span className="text-sm text-foreground">Workspace: {item.plane_workspace_id}</span>
                        <span className="text-xs text-muted-foreground">Project ID: {item.plane_project_id}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline">{item.sync_direction || 'cnb_to_plane'}</Badge>
                    </TableCell>
                    <TableCell>
                      <span className="text-sm text-muted-foreground">{item.label_selector || '—'}</span>
                    </TableCell>
                    <TableCell>
                      <div className="flex flex-col text-xs text-muted-foreground">
                        <span>Open: {item.issue_open_state_id || '—'}</span>
                        <span>Closed: {item.issue_closed_state_id || '—'}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      {item.active ? <Badge variant="success">启用</Badge> : <Badge variant="muted">停用</Badge>}
                    </TableCell>
                  </TableRow>
                ))}
                {items.length === 0 && !loading && (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-sm text-muted-foreground">
                      暂无符合条件的映射。
                    </TableCell>
                  </TableRow>
                )}
                {loading && (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-sm text-muted-foreground">
                      正在加载映射数据…
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
