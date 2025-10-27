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
import {
  Menu,
  MenuItem,
  MenuPopup,
  MenuPositioner,
  MenuTrigger,
  MenuSeparator,
  MenuPortal,
  MenuGroup,
  MenuGroupLabel,
} from '@/components/ui/menu'

const directionOptions = [
  { value: 'cnb_to_plane', label: '仅 CNB → Plane' },
  { value: 'bidirectional', label: '双向同步' },
]

type Mapping = {
  plane_project_id: string
  plane_workspace_id: string
  plane_workspace_slug?: string | null
  plane_workspace_name?: string | null
  plane_project_name?: string | null
  plane_project_identifier?: string | null
  plane_project_slug?: string | null
  cnb_repo_id: string
  issue_open_state_id?: string | null
  issue_closed_state_id?: string | null
  active: boolean
  sync_direction?: string | null
  label_selector?: string | null
}

type Feedback = { kind: 'success' | 'error'; message: string }

type ActiveFilter = 'all' | 'active' | 'inactive'

function makeKey(item: { cnb_repo_id: string; plane_project_id: string }) {
  return `${item.cnb_repo_id}::${item.plane_project_id}`
}

function trimValue(value?: string | null) {
  const trimmed = value?.trim()
  return trimmed ? trimmed : ''
}

function formatWorkspaceLabel(item: Mapping) {
  return trimValue(item.plane_workspace_name) || trimValue(item.plane_workspace_slug) || item.plane_workspace_id
}

function formatWorkspaceDetails(item: Mapping) {
  const details: string[] = []
  const slug = trimValue(item.plane_workspace_slug)
  if (slug) {
    details.push(`Slug: ${slug}`)
  }
  details.push(`ID: ${item.plane_workspace_id}`)
  return details.join(' · ')
}

function formatProjectLabel(item: Mapping) {
  return (
    trimValue(item.plane_project_name) ||
    trimValue(item.plane_project_identifier) ||
    trimValue(item.plane_project_slug) ||
    item.plane_project_id
  )
}

function formatProjectDetails(item: Mapping) {
  const details: string[] = []
  const identifier = trimValue(item.plane_project_identifier)
  const slug = trimValue(item.plane_project_slug)
  if (identifier) {
    details.push(`标识符: ${identifier}`)
  }
  if (slug && slug !== identifier) {
    details.push(`Slug: ${slug}`)
  }
  details.push(`ID: ${item.plane_project_id}`)
  return details.join(' · ')
}

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
  const [editingKey, setEditingKey] = useState<string | null>(null)
  const [actionKey, setActionKey] = useState<string | null>(null)

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
        const message = await readErrorMessage(res, `查询失败（${res.status}）`)
        throw new Error(message)
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
  const isEditing = Boolean(editingKey)
  const editingLabel = useMemo(() => {
    if (!editingKey) return ''
    const target = items.find(item => makeKey(item) === editingKey)
    if (!target) return ''
    return `${target.cnb_repo_id} → ${formatProjectLabel(target)}`
  }, [items, editingKey])

  function mappingToForm(item: Mapping): FormState {
    return {
      cnb_repo_id: item.cnb_repo_id,
      plane_workspace_id: item.plane_workspace_id,
      plane_project_id: item.plane_project_id,
      issue_open_state_id: item.issue_open_state_id ?? '',
      issue_closed_state_id: item.issue_closed_state_id ?? '',
      sync_direction: item.sync_direction ?? 'cnb_to_plane',
      label_selector: item.label_selector ?? '',
      active: item.active,
    }
  }

  function buildPayloadFromForm(state: FormState) {
    return {
      cnb_repo_id: state.cnb_repo_id.trim(),
      plane_workspace_id: state.plane_workspace_id.trim(),
      plane_project_id: state.plane_project_id.trim(),
      issue_open_state_id: state.issue_open_state_id.trim(),
      issue_closed_state_id: state.issue_closed_state_id.trim(),
      sync_direction: state.sync_direction,
      label_selector: state.label_selector.trim(),
      active: state.active,
    }
  }

  function handleEdit(item: Mapping) {
    setForm(mappingToForm(item))
    setEditingKey(makeKey(item))
    setFeedback(null)
  }

  // region: icons (inline SVG to avoid extra deps)
  function MoreIcon(props: React.SVGProps<SVGSVGElement>) {
    return (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" aria-hidden="true" {...props}>
        <circle cx="6" cy="12" r="1.5" />
        <circle cx="12" cy="12" r="1.5" />
        <circle cx="18" cy="12" r="1.5" />
      </svg>
    )
  }

  function EditIcon(props: React.SVGProps<SVGSVGElement>) {
    return (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true" {...props}>
        <path d="M12 20h9" />
        <path d="M16.5 3.5a2.121 2.121 0 1 1 3 3L7 19l-4 1 1-4 12.5-12.5z" />
      </svg>
    )
  }

  function ToggleOnIcon(props: React.SVGProps<SVGSVGElement>) {
    return (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true" {...props}>
        <rect x="2" y="6" width="20" height="12" rx="6" />
        <circle cx="16" cy="12" r="4" />
      </svg>
    )
  }

  function ToggleOffIcon(props: React.SVGProps<SVGSVGElement>) {
    return (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true" {...props}>
        <rect x="2" y="6" width="20" height="12" rx="6" />
        <circle cx="8" cy="12" r="4" />
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

  function ChevronDownIcon(props: React.SVGProps<SVGSVGElement>) {
    return (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true" {...props}>
        <path d="M6 9l6 6 6-6" />
      </svg>
    )
  }

  function handleResetForm() {
    const workspace = form.plane_workspace_id
    const project = form.plane_project_id
    setForm({ ...initialForm, plane_workspace_id: workspace, plane_project_id: project })
    setEditingKey(null)
    setFeedback(null)
  }

  async function submit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setSubmitting(true)
    setFeedback(null)
    try {
      const payload = buildPayloadFromForm(form)
      const method = isEditing ? 'PATCH' : 'POST'
      const res = await fetch('/api/admin/mappings/repo-project', {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) {
        const message = await readErrorMessage(res, `保存失败（${res.status}）`)
        throw new Error(message)
      }
      const workspace = form.plane_workspace_id
      const project = form.plane_project_id
      setFeedback({ kind: 'success', message: isEditing ? '映射已更新。' : '映射已保存并刷新。' })
      if (isEditing) {
        setEditingKey(null)
        setForm({ ...initialForm, plane_workspace_id: workspace, plane_project_id: project })
      } else {
        setForm(prev => ({ ...prev, cnb_repo_id: '' }))
      }
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '保存失败，请稍后再试。' })
    } finally {
      setSubmitting(false)
    }
  }

  async function handleToggleActive(item: Mapping, nextActive: boolean) {
    const key = makeKey(item)
    setActionKey(key)
    setFeedback(null)
    try {
      const payload = { ...buildPayloadFromForm(mappingToForm(item)), active: nextActive }
      const res = await fetch('/api/admin/mappings/repo-project', {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) {
        const message = await readErrorMessage(res, `更新失败（${res.status}）`)
        throw new Error(message)
      }
      if (editingKey === key) {
        setForm(prev => ({ ...prev, active: nextActive }))
      }
      setFeedback({ kind: 'success', message: `已${nextActive ? '启用' : '停用'}该映射。` })
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '更新映射状态失败。' })
    } finally {
      setActionKey(null)
    }
  }

  async function handleDelete(item: Mapping) {
    if (!window.confirm(`确定要删除映射 ${item.cnb_repo_id} → ${formatProjectLabel(item)} 吗？此操作将标记为停用。`)) {
      return
    }
    const key = makeKey(item)
    setActionKey(key)
    setFeedback(null)
    try {
      const res = await fetch('/api/admin/mappings/repo-project', {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ...buildPayloadFromForm(mappingToForm(item)) }),
      })
      if (!res.ok) {
        const message = await readErrorMessage(res, `删除失败（${res.status}）`)
        throw new Error(message)
      }
      if (editingKey === key) {
        const workspace = form.plane_workspace_id
        const project = form.plane_project_id
        setEditingKey(null)
        setForm({ ...initialForm, plane_workspace_id: workspace, plane_project_id: project })
      }
      setFeedback({ kind: 'success', message: '映射已删除（标记为停用）。' })
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '删除映射失败。' })
    } finally {
      setActionKey(null)
    }
  }

  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <div className="flex flex-col gap-2 md:flex-row md:items-start md:justify-between">
            <div className="space-y-1">
              <CardTitle>新增或更新 Repo ↔ Project 映射</CardTitle>
              <CardDescription>
                {isEditing
                  ? `正在编辑：${editingLabel || '选中的映射' }。如需更换仓库或项目，请先删除后重新创建。`
                  : '填写必填字段后提交即可覆盖同名仓库的最新策略。'}
              </CardDescription>
            </div>
            {isEditing && <Badge variant="outline">编辑模式</Badge>}
          </div>
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
                disabled={isEditing}
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
                disabled={isEditing}
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
              <span className="text-sm text-muted-foreground">
                {isEditing ? '保存后将覆盖该映射的同步设置。' : '保存后列表将自动刷新。'}
              </span>
              <div className="flex items-center gap-2">
                <Button type="button" variant="ghost" onClick={handleResetForm}>
                  重置表单
                </Button>
                {isEditing && (
                  <Button
                    type="button"
                    variant="ghost"
                    onClick={() => {
                      setEditingKey(null)
                      setForm(prev => ({ ...initialForm, plane_workspace_id: prev.plane_workspace_id, plane_project_id: prev.plane_project_id }))
                    }}
                  >
                    取消编辑
                  </Button>
                )}
                <Button type="submit" disabled={submitting}>
                  {submitting ? '保存中…' : isEditing ? '保存更新' : '保存映射'}
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
                  <TableHead className="text-right w-[200px]">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {items.map(item => (
                  <TableRow key={`${item.cnb_repo_id}-${item.plane_project_id}`}>
                    <TableCell>
                      <span className="font-medium text-foreground">{item.cnb_repo_id}</span>
                    </TableCell>
                    <TableCell>
                      <div className="flex flex-col gap-2">
                        <div className="flex flex-col gap-1">
                          <span className="text-sm text-foreground">Workspace: {formatWorkspaceLabel(item)}</span>
                          <span className="text-xs text-muted-foreground">{formatWorkspaceDetails(item)}</span>
                        </div>
                        <div className="flex flex-col gap-1">
                          <span className="text-sm text-foreground">Project: {formatProjectLabel(item)}</span>
                          <span className="text-xs text-muted-foreground">{formatProjectDetails(item)}</span>
                        </div>
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
                    <TableCell className="w-[200px]">
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          className="shrink-0"
                          variant="ghost"
                          size="sm"
                          onClick={() => handleEdit(item)}
                          title="编辑映射"
                        >
                          编辑
                        </Button>
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
                              <MenuPopup className="p-1 min-w-[12rem]">
                                <MenuGroup>
                                  <MenuGroupLabel>快速操作</MenuGroupLabel>
                                  <MenuItem
                                    onClick={() => handleToggleActive(item, !item.active)}
                                    className="justify-start"
                                    data-active={item.active ? true : undefined}
                                    disabled={actionKey === makeKey(item)}
                                  >
                                    {item.active ? (
                                      <span className="inline-flex items-center gap-2"><ToggleOffIcon className="h-4 w-4" /> 停用</span>
                                    ) : (
                                      <span className="inline-flex items-center gap-2"><ToggleOnIcon className="h-4 w-4" /> 启用</span>
                                    )}
                                  </MenuItem>
                                  <MenuItem onClick={() => handleEdit(item)} className="justify-start">
                                    <span className="inline-flex items-center gap-2"><EditIcon className="h-4 w-4" /> 编辑</span>
                                  </MenuItem>
                                </MenuGroup>
                                <MenuSeparator />
                                <MenuGroup>
                                  <MenuGroupLabel>危险操作</MenuGroupLabel>
                                  <MenuItem
                                    onClick={() => handleDelete(item)}
                                    className="justify-start text-destructive-foreground"
                                    disabled={actionKey === makeKey(item)}
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
                ))}
                {items.length === 0 && !loading && (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center text-sm text-muted-foreground">
                      暂无符合条件的映射。
                    </TableCell>
                  </TableRow>
                )}
                {loading && (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center text-sm text-muted-foreground">
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
