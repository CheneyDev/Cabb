'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'

import { Alert } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Menu, MenuItem, MenuPopup, MenuPositioner, MenuTrigger, MenuPortal } from '@/components/ui/menu'

type Feedback = { kind: 'success' | 'error'; message: string }

type UserMapping = {
  plane_user_id: string
  cnb_user_id?: string | null
  lark_user_id?: string | null
  display_name?: string | null
  connected_at?: string | null
  created_at: string
  updated_at: string
}

type FormRow = {
  plane_user_id: string
  cnb_user_id: string
  display_name: string
}

const initialRow: FormRow = {
  plane_user_id: '',
  cnb_user_id: '',
  display_name: '',
}

export default function UsersPage() {
  const [rows, setRows] = useState<FormRow[]>([initialRow])
  const [mappings, setMappings] = useState<UserMapping[]>([])
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [feedback, setFeedback] = useState<Feedback | null>(null)
  const [editing, setEditing] = useState<UserMapping | null>(null)
  const [editForm, setEditForm] = useState<FormRow>(initialRow)
  const [editPending, setEditPending] = useState(false)
  const [filterPlaneUser, setFilterPlaneUser] = useState('')
  const [filterCNBUser, setFilterCNBUser] = useState('')
  const [search, setSearch] = useState('')
  const [limit, setLimit] = useState(50)

  const querySuffix = useMemo(() => {
    const params = new URLSearchParams()
    const plane = filterPlaneUser.trim()
    const cnb = filterCNBUser.trim()
    const keyword = search.trim()
    if (plane) params.set('plane_user_id', plane)
    if (cnb) params.set('cnb_user_id', cnb)
    if (keyword) params.set('q', keyword)
    if (limit > 0) params.set('limit', String(limit))
    const query = params.toString()
    return query ? `?${query}` : ''
  }, [filterPlaneUser, filterCNBUser, search, limit])

  const load = useCallback(async () => {
    setLoading(true)
    setFetchError(null)
    try {
      const res = await fetch(`/api/admin/mappings/users${querySuffix}`, { cache: 'no-store' })
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
      setMappings(Array.isArray(json.items) ? json.items : [])
    } catch (err) {
      setFetchError(err instanceof Error ? err.message : '加载用户映射失败')
    } finally {
      setLoading(false)
    }
  }, [querySuffix])

  useEffect(() => {
    load()
  }, [load])

  useEffect(() => {
    if (!editing?.cnb_user_id) return
    const match = mappings.find(item => item.cnb_user_id === editing.cnb_user_id)
    if (!match) {
      setEditing(null)
      setEditForm(initialRow)
      return
    }
    setEditForm({
      plane_user_id: match.plane_user_id,
      cnb_user_id: match.cnb_user_id ?? '',
      display_name: match.display_name ?? '',
    })
  }, [mappings, editing?.cnb_user_id])

  function updateRow(index: number, field: keyof FormRow, value: string) {
    setRows(prev => prev.map((row, idx) => (idx === index ? { ...row, [field]: value } : row)))
  }

  function addRow() {
    setRows(prev => [...prev, initialRow])
  }

  function removeRow(index: number) {
    setRows(prev => (prev.length <= 1 ? prev : prev.filter((_, idx) => idx !== index)))
  }

  async function submit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setSubmitting(true)
    setFeedback(null)
    try {
      const payloadRows = rows
        .map(row => ({
          plane_user_id: row.plane_user_id.trim(),
          cnb_user_id: row.cnb_user_id.trim(),
          display_name: row.display_name.trim(),
        }))
        .filter(row => row.plane_user_id && row.cnb_user_id)

      if (payloadRows.length === 0) {
        throw new Error('请至少填写一条有效的 Plane 用户 ID 与 CNB 用户 ID 对。')
      }

      const res = await fetch('/api/admin/mappings/users', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mappings: payloadRows }),
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
      setFeedback({ kind: 'success', message: `已成功保存 ${payloadRows.length} 条映射。` })
      setRows([initialRow])
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '保存失败，请稍后再试。' })
    } finally {
      setSubmitting(false)
    }
  }

  function startEdit(mapping: UserMapping) {
    if (!mapping.cnb_user_id) {
      setFeedback({ kind: 'error', message: '该映射缺少 CNB 用户 ID，无法编辑。' })
      return
    }
    setEditing(mapping)
    setEditForm({
      plane_user_id: mapping.plane_user_id,
      cnb_user_id: mapping.cnb_user_id ?? '',
      display_name: mapping.display_name ?? '',
    })
    setFeedback(null)
  }

  function cancelEdit() {
    setEditing(null)
    setEditForm(initialRow)
  }

  async function submitEdit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!editing) {
      setFeedback({ kind: 'error', message: '请先在列表中选择需要编辑的映射。' })
      return
    }
    const planeUserID = editForm.plane_user_id.trim()
    const cnbUserID = editForm.cnb_user_id.trim()
    if (!planeUserID || !cnbUserID) {
      setFeedback({ kind: 'error', message: 'Plane 用户 ID 与 CNB 用户 ID 均不能为空。' })
      return
    }
    setEditPending(true)
    setFeedback(null)
    try {
      const res = await fetch('/api/admin/mappings/users', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          mappings: [
            {
              plane_user_id: planeUserID,
              cnb_user_id: cnbUserID,
              display_name: editForm.display_name.trim(),
            },
          ],
        }),
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
      setFeedback({ kind: 'success', message: '映射已更新。' })
      cancelEdit()
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '保存失败，请稍后再试。' })
    } finally {
      setEditPending(false)
    }
  }

function formatDate(value?: string | null) {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN', { hour12: false })
}

// region: inline icons for menu
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
// endregion

  const connectedCount = mappings.filter(item => Boolean(item.connected_at)).length

  const editingLabel = useMemo(() => {
    if (!editing) return '选择列表中的映射进行编辑。'
    const display = editForm.display_name || editing.display_name || ''
    const cnb = editing.cnb_user_id || editForm.cnb_user_id || ''
    const plane = editForm.plane_user_id
    const segments = [display || plane, cnb ? `CNB: ${cnb}` : '', `Plane: ${plane}`].filter(Boolean)
    return segments.join(' ｜ ')
  }, [editing, editForm])

  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle>批量维护 Plane ↔ CNB 用户映射</CardTitle>
          <CardDescription>手动绑定后，评论与任务更新会准确归属到对应成员。</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="grid gap-4" onSubmit={submit}>
            <div className="space-y-4">
              {rows.map((row, index) => (
                <div key={index} className="grid gap-4 rounded-3xl border border-border px-4 py-4 md:grid-cols-3">
                  <div className="grid gap-2">
                    <Label htmlFor={`plane_user_${index}`}>Plane 用户 ID</Label>
                    <Input
                      id={`plane_user_${index}`}
                      placeholder="Plane 用户 UUID"
                      value={row.plane_user_id}
                      onChange={event => updateRow(index, 'plane_user_id', event.target.value)}
                    />
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor={`cnb_user_${index}`}>CNB 用户 ID</Label>
                    <Input
                      id={`cnb_user_${index}`}
                      placeholder="CNB 用户名或 ID"
                      value={row.cnb_user_id}
                      onChange={event => updateRow(index, 'cnb_user_id', event.target.value)}
                    />
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor={`display_name_${index}`}>显示名（可选）</Label>
                    <Input
                      id={`display_name_${index}`}
                      placeholder="用于后台备注的名称"
                      value={row.display_name}
                      onChange={event => updateRow(index, 'display_name', event.target.value)}
                    />
                  </div>
                  <div className="flex flex-col gap-2 text-sm text-muted-foreground md:col-span-3 md:flex-row md:items-center md:justify-end">
                    {rows.length > 1 && (
                      <Button type="button" variant="ghost" onClick={() => removeRow(index)}>
                        移除此行
                      </Button>
                    )}
                    {index === rows.length - 1 && (
                      <Button type="button" variant="outline" onClick={addRow}>
                        新增一行
                      </Button>
                    )}
                  </div>
                </div>
              ))}
            </div>
            <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
              <span className="text-sm text-muted-foreground">提交后列表会立即刷新，确保最新绑定生效。</span>
              <Button type="submit" disabled={submitting}>
                {submitting ? '保存中…' : '保存映射'}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div className="flex flex-col gap-2 md:flex-row md:items-start md:justify-between">
            <div className="space-y-1">
              <CardTitle>编辑已有用户映射</CardTitle>
              <CardDescription>{editingLabel}</CardDescription>
            </div>
            {editing && <Badge variant="outline">编辑模式</Badge>}
          </div>
        </CardHeader>
        <CardContent>
          <form className="grid gap-4 md:grid-cols-3" onSubmit={submitEdit}>
            <div className="grid gap-2">
              <Label htmlFor="edit_plane_user">Plane 用户 ID</Label>
              <Input
                id="edit_plane_user"
                placeholder="Plane 用户 UUID"
                value={editForm.plane_user_id}
                onChange={event => setEditForm(prev => ({ ...prev, plane_user_id: event.target.value }))}
                disabled={!editing}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit_cnb_user">CNB 用户 ID</Label>
              <Input
                id="edit_cnb_user"
                placeholder="CNB 用户名或 ID"
                value={editForm.cnb_user_id}
                onChange={event => setEditForm(prev => ({ ...prev, cnb_user_id: event.target.value }))}
                disabled={!editing}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="edit_display_name">显示名（可选）</Label>
              <Input
                id="edit_display_name"
                placeholder="用于后台备注的名称"
                value={editForm.display_name}
                onChange={event => setEditForm(prev => ({ ...prev, display_name: event.target.value }))}
                disabled={!editing}
              />
            </div>
            <div className="md:col-span-3 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
              <span className="text-sm text-muted-foreground">
                {editing
                  ? '保存后会覆盖该条映射的 Plane 用户或显示名配置。'
                  : '从下方列表选择映射后即可在此处更新绑定信息。'}
              </span>
              <div className="flex items-center gap-2">
                <Button type="button" variant="ghost" onClick={cancelEdit} disabled={!editing}>
                  取消
                </Button>
                <Button type="submit" disabled={!editing || editPending}>
                  {editPending ? '保存中…' : '保存更新'}
                </Button>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>

      {feedback && <Alert variant={feedback.kind === 'success' ? 'success' : 'destructive'}>{feedback.message}</Alert>}

      <Card>
        <CardHeader>
          <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <div className="space-y-1">
              <CardTitle>用户映射列表</CardTitle>
              <CardDescription>按照 Plane 用户、CNB 用户或关键字过滤，快速定位绑定状态。</CardDescription>
            </div>
            <div className="flex items-center gap-2">
              <Badge variant="success">已连接 {connectedCount}</Badge>
              <Badge variant="outline">全部 {mappings.length}</Badge>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-3 md:grid-cols-[2fr_2fr_2fr_auto] md:items-end">
            <div className="grid gap-2">
              <Label htmlFor="filterPlane">Plane 用户 ID</Label>
              <Input id="filterPlane" value={filterPlaneUser} onChange={event => setFilterPlaneUser(event.target.value)} />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="filterCNB">CNB 用户 ID</Label>
              <Input id="filterCNB" value={filterCNBUser} onChange={event => setFilterCNBUser(event.target.value)} />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="search">搜索（显示名 / ID）</Label>
              <Input id="search" value={search} onChange={event => setSearch(event.target.value)} />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="limit">返回条数</Label>
              <Input
                id="limit"
                type="number"
                min={1}
                max={200}
                value={limit}
                onChange={event => setLimit(Number(event.target.value) || 50)}
              />
            </div>
            <div className="flex flex-wrap gap-2 md:col-span-4 md:justify-end">
              <Button type="button" variant="outline" onClick={() => load()} disabled={loading}>
                {loading ? '加载中…' : '刷新'}
              </Button>
              <Button
                type="button"
                variant="ghost"
                onClick={() => {
                  setFilterPlaneUser('')
                  setFilterCNBUser('')
                  setSearch('')
                  setLimit(50)
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
                  <TableHead>Plane 用户</TableHead>
                  <TableHead>CNB 用户</TableHead>
                  <TableHead>显示名</TableHead>
                  <TableHead>Lark 用户</TableHead>
                  <TableHead>连接时间</TableHead>
                  <TableHead>最近更新</TableHead>
                  <TableHead className="text-right w-[160px]">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {mappings.map(item => (
                  <TableRow key={`${item.plane_user_id}-${item.cnb_user_id ?? 'unknown'}`}>
                    <TableCell>
                      <div className="flex flex-col gap-1">
                        <span className="font-medium text-foreground">{item.plane_user_id}</span>
                        <span className="text-xs text-muted-foreground">创建: {formatDate(item.created_at)}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <span className="text-sm text-foreground">{item.cnb_user_id || '—'}</span>
                    </TableCell>
                    <TableCell>
                      <span className="text-sm text-muted-foreground">{item.display_name || '—'}</span>
                    </TableCell>
                    <TableCell>
                      <span className="text-sm text-muted-foreground">{item.lark_user_id || '—'}</span>
                    </TableCell>
                    <TableCell>
                      {item.connected_at ? <Badge variant="success">{formatDate(item.connected_at)}</Badge> : <Badge variant="muted">未连接</Badge>}
                    </TableCell>
                    <TableCell>
                      <span className="text-sm text-foreground">{formatDate(item.updated_at)}</span>
                    </TableCell>
                    <TableCell className="w-[160px]">
                      <div className="flex items-center justify-end">
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
                                onClick={() => startEdit(item)}
                                className="justify-start"
                                disabled={!item.cnb_user_id}
                              >
                                <span className="inline-flex items-center gap-2"><EditIcon className="h-4 w-4" /> 编辑</span>
                              </MenuItem>
                              </MenuPopup>
                            </MenuPositioner>
                          </MenuPortal>
                        </Menu>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
                {mappings.length === 0 && !loading && (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center text-sm text-muted-foreground">
                      暂无符合条件的用户映射。
                    </TableCell>
                  </TableRow>
                )}
                {loading && (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center text-sm text-muted-foreground">
                      正在加载用户映射…
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
        <CardFooter className="text-sm text-muted-foreground">
          映射数据来源于数据库 <code>user_mappings</code> 表，最新更新时间以 UTC 显示。
        </CardFooter>
      </Card>
    </div>
  )
}
