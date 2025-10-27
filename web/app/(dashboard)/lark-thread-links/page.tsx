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
import { Menu, MenuItem, MenuPopup, MenuPositioner, MenuTrigger, MenuSeparator, MenuPortal } from '@/components/ui/menu'

const initialForm = {
  lark_thread_id: '',
  plane_issue_id: '',
  plane_project_id: '',
  workspace_slug: '',
  sync_enabled: false,
}

type FormState = typeof initialForm

type ThreadLink = {
  lark_thread_id: string
  lark_chat_id?: string | null
  lark_chat_name?: string | null
  plane_issue_id: string
  plane_issue_name?: string | null
  plane_project_id?: string | null
  plane_project_name?: string | null
  plane_project_identifier?: string | null
  plane_project_slug?: string | null
  workspace_slug?: string | null
  workspace_name?: string | null
  sync_enabled: boolean
  linked_at?: string
  created_at?: string
  updated_at?: string
}

type Feedback = { kind: 'success' | 'error'; message: string }

type Filters = {
  plane_issue_id: string
  lark_thread_id: string
  sync_enabled: 'all' | 'true' | 'false'
}

function safeTrim(value?: string | null) {
  return value?.trim() ?? ''
}

function firstNonEmpty(...values: (string | null | undefined)[]) {
  for (const value of values) {
    const trimmed = value?.trim()
    if (trimmed) {
      return trimmed
    }
  }
  return ''
}

function makeKey(item: ThreadLink) {
  return `${item.lark_thread_id}::${item.plane_issue_id}`
}

export default function LarkThreadLinksPage() {
  const [form, setForm] = useState<FormState>(initialForm)
  const [filters, setFilters] = useState<Filters>({ plane_issue_id: '', lark_thread_id: '', sync_enabled: 'all' })
  const [items, setItems] = useState<ThreadLink[]>([])
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [feedback, setFeedback] = useState<Feedback | null>(null)
  const [actionKey, setActionKey] = useState<string | null>(null)

  const querySuffix = useMemo(() => {
    const params = new URLSearchParams()
    params.set('limit', '100')
    const planeIssue = filters.plane_issue_id.trim()
    const thread = filters.lark_thread_id.trim()
    if (planeIssue) params.set('plane_issue_id', planeIssue)
    if (thread) params.set('lark_thread_id', thread)
    if (filters.sync_enabled !== 'all') {
      params.set('sync_enabled', filters.sync_enabled)
    }
    const query = params.toString()
    return query ? `?${query}` : ''
  }, [filters])

  const load = useCallback(async () => {
    setLoading(true)
    setFetchError(null)
    try {
      const res = await fetch(`/api/admin/links/lark-thread-links${querySuffix}`, { cache: 'no-store' })
      if (!res.ok) {
        const message = await readErrorMessage(res, `查询失败（${res.status}）`)
        throw new Error(message)
      }
      const data = await res.json()
      setItems(Array.isArray(data.items) ? (data.items as ThreadLink[]) : [])
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
        lark_thread_id: form.lark_thread_id.trim(),
        plane_issue_id: form.plane_issue_id.trim(),
        plane_project_id: form.plane_project_id.trim(),
        workspace_slug: form.workspace_slug.trim(),
        sync_enabled: form.sync_enabled,
      }
      if (!payload.lark_thread_id || !payload.plane_issue_id) {
        throw new Error('请填写 Lark Thread ID 与 Plane Issue ID。')
      }
      const res = await fetch('/api/admin/links/lark-thread-links', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) {
        const message = await readErrorMessage(res, `保存失败（${res.status}）`)
        throw new Error(message)
      }
      setFeedback({ kind: 'success', message: '线程映射已保存。' })
      setForm(initialForm)
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '保存失败，请稍后再试。' })
    } finally {
      setSubmitting(false)
    }
  }

  async function toggleSync(item: ThreadLink) {
    setActionKey(makeKey(item))
    try {
      const payload = {
        lark_thread_id: item.lark_thread_id,
        plane_issue_id: item.plane_issue_id,
        plane_project_id: item.plane_project_id ?? '',
        workspace_slug: item.workspace_slug ?? '',
        sync_enabled: !item.sync_enabled,
      }
      const res = await fetch('/api/admin/links/lark-thread-links', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) {
        const message = await readErrorMessage(res, `更新失败（${res.status}）`)
        throw new Error(message)
      }
      setFeedback({ kind: 'success', message: `已${item.sync_enabled ? '关闭' : '开启'}线程同步。` })
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '更新失败，请稍后再试。' })
    } finally {
      setActionKey(null)
    }
  }

  async function handleDelete(item: ThreadLink) {
    setActionKey(makeKey(item) + ':delete')
    try {
      const res = await fetch('/api/admin/links/lark-thread-links', {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ lark_thread_id: item.lark_thread_id }),
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
      setActionKey(null)
    }
  }

  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle>飞书线程 ↔ Plane Issue 映射</CardTitle>
          <CardDescription>绑定后可控制线程自动同步 Plane 评论，或手动解除错误绑定。</CardDescription>
        </CardHeader>
        <CardContent className="grid gap-6">
          <form className="grid gap-4 md:grid-cols-2" onSubmit={handleSubmit}>
            <div className="grid gap-2">
              <Label htmlFor="lark_thread_id">Lark Thread ID</Label>
              <Input
                id="lark_thread_id"
                placeholder="om_xxx"
                value={form.lark_thread_id}
                onChange={e => setForm(prev => ({ ...prev, lark_thread_id: e.target.value }))}
              />
            </div>
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
              <Label htmlFor="plane_project_id">Plane Project ID（可选）</Label>
              <Input
                id="plane_project_id"
                placeholder="uuid"
                value={form.plane_project_id}
                onChange={e => setForm(prev => ({ ...prev, plane_project_id: e.target.value }))}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="workspace_slug">Plane Workspace Slug（可选）</Label>
              <Input
                id="workspace_slug"
                placeholder="workspace"
                value={form.workspace_slug}
                onChange={e => setForm(prev => ({ ...prev, workspace_slug: e.target.value }))}
              />
            </div>
            <div className="flex items-center gap-3">
              <Switch
                id="sync_enabled"
                checked={form.sync_enabled}
                onChange={e => setForm(prev => ({ ...prev, sync_enabled: e.target.checked }))}
              />
              <Label htmlFor="sync_enabled" className="text-sm text-muted-foreground">
                开启后，线程中新消息将自动同步至 Plane
              </Label>
            </div>
            <div className="md:col-span-2 flex justify-end">
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
              <Label htmlFor="filter_thread">按线程 ID 筛选</Label>
              <Input
                id="filter_thread"
                placeholder="om_xxx"
                value={filters.lark_thread_id}
                onChange={e => setFilters(prev => ({ ...prev, lark_thread_id: e.target.value }))}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="filter_sync">同步状态</Label>
              <Select
                id="filter_sync"
                value={filters.sync_enabled}
                onChange={e => setFilters(prev => ({ ...prev, sync_enabled: e.target.value as Filters['sync_enabled'] }))}
              >
                <option value="all">全部</option>
                <option value="true">仅自动同步</option>
                <option value="false">仅关闭同步</option>
              </Select>
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
          <CardTitle>线程映射列表</CardTitle>
          <CardDescription>按更新时间倒序排列，可快速开启/关闭自动同步或解除绑定。</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Lark Chat</TableHead>
                  <TableHead>Lark Thread ID</TableHead>
                  <TableHead>Plane Issue</TableHead>
                  <TableHead>Plane Project</TableHead>
                  <TableHead>Workspace</TableHead>
                  <TableHead>自动同步</TableHead>
                  <TableHead>Linked At</TableHead>
                  <TableHead>更新时间</TableHead>
                  <TableHead className="text-right w-[200px]">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {items.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={9} className="text-center text-sm text-muted-foreground">
                      暂无数据，请调整筛选条件。
                    </TableCell>
                  </TableRow>
                ) : (
                  items.map(item => {
                    const key = makeKey(item)
                    const toggleBusy = actionKey === key
                    const deleteBusy = actionKey === `${key}:delete`
                    const chatName = safeTrim(item.lark_chat_name)
                    const chatId = safeTrim(item.lark_chat_id)
                    const issueName = safeTrim(item.plane_issue_name)
                    const issueId = safeTrim(item.plane_issue_id)
                    const projectName = firstNonEmpty(item.plane_project_name, item.plane_project_identifier, item.plane_project_slug)
                    const projectId = safeTrim(item.plane_project_id)
                    const projectIdentifier = safeTrim(item.plane_project_identifier)
                    const projectSlug = safeTrim(item.plane_project_slug)
                    const workspaceDisplay = firstNonEmpty(item.workspace_name, item.workspace_slug)
                    const rawWorkspaceName = safeTrim(item.workspace_name)
                    const workspaceSlug = safeTrim(item.workspace_slug)
                    const projectDetails: { key: string; value: string; monospace?: boolean }[] = []
                    if (projectName && projectId) {
                      projectDetails.push({ key: 'id', value: projectId, monospace: true })
                    }
                    if (projectIdentifier && projectIdentifier !== projectName) {
                      projectDetails.push({ key: 'identifier', value: `标识符: ${projectIdentifier}` })
                    }
                    if (projectSlug && projectSlug !== projectIdentifier && projectSlug !== projectName) {
                      projectDetails.push({ key: 'slug', value: `Slug: ${projectSlug}` })
                    }
                    const workspaceShowSlug = rawWorkspaceName && workspaceSlug && rawWorkspaceName !== workspaceSlug
                    return (
                      <TableRow key={key}>
                        <TableCell>
                          {chatName || chatId ? (
                            <div className="flex flex-col">
                              <span className="font-medium">{chatName || chatId}</span>
                              {chatName && chatId ? (
                                <span className="font-mono text-xs text-muted-foreground">{chatId}</span>
                              ) : null}
                            </div>
                          ) : (
                            <span className="text-muted-foreground">—</span>
                          )}
                        </TableCell>
                        <TableCell className="font-mono text-xs">{item.lark_thread_id}</TableCell>
                        <TableCell>
                          {issueName || issueId ? (
                            <div className="flex flex-col">
                              <span className="font-medium">{issueName || issueId}</span>
                              {issueName ? (
                                <span className="font-mono text-xs text-muted-foreground">{issueId}</span>
                              ) : null}
                            </div>
                          ) : (
                            <span className="text-muted-foreground">—</span>
                          )}
                        </TableCell>
                        <TableCell>
                          {projectName || projectId ? (
                            <div className="flex flex-col gap-1">
                              <span className="font-medium">{projectName || projectId}</span>
                              {projectDetails.length > 0 ? (
                                <div className="flex flex-col text-xs text-muted-foreground">
                                  {projectDetails.map(detail => (
                                    <span
                                      key={detail.key}
                                      className={detail.monospace ? 'font-mono' : undefined}
                                    >
                                      {detail.value}
                                    </span>
                                  ))}
                                </div>
                              ) : null}
                            </div>
                          ) : (
                            <span className="text-muted-foreground">—</span>
                          )}
                        </TableCell>
                        <TableCell>
                          {workspaceDisplay ? (
                            <div className="flex flex-col gap-1">
                              <span className="font-medium">{workspaceDisplay}</span>
                              {workspaceShowSlug ? (
                                <span className="font-mono text-xs text-muted-foreground">{workspaceSlug}</span>
                              ) : null}
                            </div>
                          ) : (
                            <span className="text-muted-foreground">—</span>
                          )}
                        </TableCell>
                        <TableCell>
                          <Badge variant={item.sync_enabled ? 'success' : 'muted'}>
                            {item.sync_enabled ? '已开启' : '已关闭'}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-xs text-muted-foreground">
                          {item.linked_at ? new Date(item.linked_at).toLocaleString() : '—'}
                        </TableCell>
                        <TableCell className="text-xs text-muted-foreground">
                          {item.updated_at ? new Date(item.updated_at).toLocaleString() : '—'}
                        </TableCell>
                        <TableCell className="w-[200px]">
                          <div className="flex items-center justify-end gap-1">
                            <Button
                              type="button"
                              variant="outline"
                              size="sm"
                              onClick={() => toggleSync(item)}
                              disabled={toggleBusy || deleteBusy}
                              title={item.sync_enabled ? '关闭自动同步' : '开启自动同步'}
                            >
                              {toggleBusy ? '更新中…' : item.sync_enabled ? '关闭同步' : '开启同步'}
                            </Button>
                            <Menu>
                              <MenuTrigger
                                aria-label="更多操作"
                                title="更多操作"
                                className="!min-w-0 !px-0 !py-0 w-9 h-9 rounded-full border-transparent bg-transparent hover:bg-[color-mix(in_srgb,var(--foreground)_6%,transparent)]"
                              >
                                <MoreIcon className="h-4 w-4" />
                              </MenuTrigger>
                              <MenuPortal>
                                <MenuPositioner>
                                  <MenuPopup className="p-1 min-w-[10rem]">
                                  <MenuItem
                                    onSelect={() => handleDelete(item)}
                                    className="justify-start text-destructive-foreground"
                                    disabled={deleteBusy || toggleBusy}
                                  >
                                    <span className="inline-flex items-center gap-2"><TrashIcon className="h-4 w-4" /> 删除</span>
                                  </MenuItem>
                                  </MenuPopup>
                                </MenuPositioner>
                              </MenuPortal>
                            </Menu>
                          </div>
                        </TableCell>
                      </TableRow>
                    )
                  })
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

// region: inline icons (kept locally to avoid deps)
function MoreIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" aria-hidden="true" {...props}>
      <circle cx="6" cy="12" r="1.5" />
      <circle cx="12" cy="12" r="1.5" />
      <circle cx="18" cy="12" r="1.5" />
    </svg>
  )
}

function TrashIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true" {...props}>
      <path d="M3 6h18" />
      <path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
      <path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6" />
      <path d="M10 11v6M14 11v6" />
    </svg>
  )
}
// endregion
