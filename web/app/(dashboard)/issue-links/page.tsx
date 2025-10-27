'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'

import { Alert } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Menu, MenuItem, MenuPopup, MenuPositioner, MenuTrigger, MenuPortal, MenuGroup, MenuGroupLabel } from '@/components/ui/menu'

const initialForm = {
  plane_issue_id: '',
  cnb_repo_id: '',
  cnb_issue_id: '',
}

type FormState = typeof initialForm

type IssueLink = {
  plane_issue_id: string
  plane_issue_name?: string | null
  cnb_repo_id?: string | null
  cnb_issue_id?: string | null
  plane_project_id?: string | null
  plane_project_name?: string | null
  plane_project_identifier?: string | null
  plane_project_slug?: string | null
  plane_workspace_id?: string | null
  plane_workspace_name?: string | null
  plane_workspace_slug?: string | null
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
                  <TableHead>Plane Issue</TableHead>
                  <TableHead>Plane Project</TableHead>
                  <TableHead>Workspace</TableHead>
                  <TableHead>CNB Issue</TableHead>
                  <TableHead>Linked At</TableHead>
                  <TableHead>更新时间</TableHead>
                  <TableHead className="text-right w-[160px]">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {items.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center text-sm text-muted-foreground">
                      暂无数据，请调整筛选条件。
                    </TableCell>
                  </TableRow>
                ) : (
                  items.map(item => {
                    const issueName = safeTrim(item.plane_issue_name)
                    const issueId = safeTrim(item.plane_issue_id)
                    const projectName = firstNonEmpty(
                      item.plane_project_name,
                      item.plane_project_identifier,
                      item.plane_project_slug,
                    )
                    const projectId = safeTrim(item.plane_project_id)
                    const projectIdentifier = safeTrim(item.plane_project_identifier)
                    const projectSlug = safeTrim(item.plane_project_slug)
                    const projectDisplay = projectName || projectId
                    const workspaceName = safeTrim(item.plane_workspace_name)
                    const workspaceSlug = safeTrim(item.plane_workspace_slug)
                    const workspaceId = safeTrim(item.plane_workspace_id)
                    const workspaceDisplay = firstNonEmpty(workspaceName, workspaceSlug) || workspaceId
                    const projectDetails: { key: string; value: string; monospace?: boolean }[] = []
                    if (projectId && projectDisplay !== projectId) {
                      projectDetails.push({ key: 'id', value: `ID: ${projectId}`, monospace: true })
                    }
                    if (projectIdentifier && projectIdentifier !== projectDisplay) {
                      projectDetails.push({ key: 'identifier', value: `标识符: ${projectIdentifier}` })
                    }
                    if (
                      projectSlug &&
                      projectSlug !== projectDisplay &&
                      projectSlug !== projectIdentifier
                    ) {
                      projectDetails.push({ key: 'slug', value: `Slug: ${projectSlug}` })
                    }
                    const workspaceDetails: { key: string; value: string; monospace?: boolean }[] = []
                    if (workspaceSlug && workspaceSlug !== workspaceDisplay) {
                      workspaceDetails.push({ key: 'slug', value: `Slug: ${workspaceSlug}`, monospace: true })
                    }
                    if (workspaceId && workspaceId !== workspaceDisplay) {
                      workspaceDetails.push({ key: 'id', value: `ID: ${workspaceId}`, monospace: true })
                    }
                    const cnbRepo = safeTrim(item.cnb_repo_id)
                    const cnbIssue = safeTrim(item.cnb_issue_id)
                    const rowKey = makeKey(item)
                    return (
                      <TableRow key={rowKey}>
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
                          {projectDisplay ? (
                            <div className="flex flex-col gap-1">
                              <span className="font-medium">{projectDisplay}</span>
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
                              {workspaceDetails.length > 0 ? (
                                <div className="flex flex-col text-xs text-muted-foreground">
                                  {workspaceDetails.map(detail => (
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
                          {cnbRepo || cnbIssue ? (
                            <div className="flex flex-col gap-1">
                              <span className="font-medium">{cnbRepo || `Issue #${cnbIssue}`}</span>
                              {cnbIssue ? (
                                <span className="font-mono text-xs text-muted-foreground">
                                  {cnbRepo ? `IID: ${cnbIssue}` : cnbIssue}
                                </span>
                              ) : null}
                            </div>
                          ) : (
                            <span className="text-muted-foreground">—</span>
                          )}
                        </TableCell>
                        <TableCell className="text-xs text-muted-foreground">
                          {item.linked_at ? new Date(item.linked_at).toLocaleString() : '—'}
                        </TableCell>
                        <TableCell className="text-xs text-muted-foreground">
                          {item.updated_at ? new Date(item.updated_at).toLocaleString() : '—'}
                        </TableCell>
                        <TableCell className="w-[160px]">
                          <div className="flex items-center justify-end gap-1">
                            <Menu>
                              <MenuTrigger
                                aria-label="更多操作"
                                title="更多操作"
                                className="!min-w-0 h-9 px-3"
                              >
                                <span className="text-sm">更多</span>
                                <ChevronDownIcon className="h-4 w-4" />
                              </MenuTrigger>
                              <MenuPortal>
                                <MenuPositioner>
                                  <MenuPopup className="p-1 min-w-[12rem]" align="end" sideOffset={6}>
                                    <MenuGroup>
                                      <MenuGroupLabel>危险操作</MenuGroupLabel>
                                      <MenuItem
                                        onClick={() => handleDelete(item)}
                                        className="justify-start text-destructive-foreground"
                                        disabled={deletingKey === rowKey}
                                      >
                                        <span className="inline-flex items-center gap-2"><TrashIcon className="h-4 w-4" /> 删除</span>
                                      </MenuItem>
                                    </MenuGroup>
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

// region: inline icons
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
function ChevronDownIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true" {...props}>
      <path d="M6 9l6 6 6-6" />
    </svg>
  )
}
// endregion
