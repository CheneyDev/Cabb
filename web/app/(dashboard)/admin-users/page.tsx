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

type AdminUser = {
  id: string
  email: string
  display_name: string
  role: string
  active: boolean
  last_login_at?: string | null
  created_at: string
  updated_at: string
}

type Feedback = { kind: 'success' | 'error'; message: string }

type CreateForm = { email: string; display_name: string; password: string }

type UpdateForm = { id: string; display_name: string; role: string; active: boolean }

type ResetForm = { id: string; password: string }

const initialCreate: CreateForm = { email: '', display_name: '', password: '' }
const initialUpdate: UpdateForm = { id: '', display_name: '', role: 'admin', active: true }
const initialReset: ResetForm = { id: '', password: '' }

export default function AdminUsersPage() {
  const [users, setUsers] = useState<AdminUser[]>([])
  const [loading, setLoading] = useState(true)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [feedback, setFeedback] = useState<Feedback | null>(null)

  const [createForm, setCreateForm] = useState<CreateForm>(initialCreate)
  const [createPending, setCreatePending] = useState(false)

  const [updateForm, setUpdateForm] = useState<UpdateForm>(initialUpdate)
  const [updatePending, setUpdatePending] = useState(false)

  const [resetForm, setResetForm] = useState<ResetForm>(initialReset)
  const [resetPending, setResetPending] = useState(false)

  const [togglingId, setTogglingId] = useState<string | null>(null)

  const load = useCallback(async () => {
    setLoading(true)
    setFetchError(null)
    try {
      const res = await fetch('/api/admin/access/admin-users', { cache: 'no-store' })
      if (!res.ok) {
        const message = await resolveError(res, `加载系统用户失败（${res.status}）`)
        throw new Error(message)
      }
      const json = await res.json()
      const items: AdminUser[] = Array.isArray(json?.items) ? json.items : []
      setUsers(items)
    } catch (err) {
      setFetchError(err instanceof Error ? err.message : '加载系统用户失败')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  useEffect(() => {
    if (!updateForm.id) return
    const target = users.find(item => item.id === updateForm.id)
    if (!target) return
    setUpdateForm(prev => ({ ...prev, display_name: target.display_name, role: target.role, active: target.active }))
  }, [users, updateForm.id])

  const userOptions = useMemo(() => users.map(user => ({ value: user.id, label: `${user.display_name || user.email}` })), [users])

  function handleCreateChange(field: keyof CreateForm, value: string) {
    setCreateForm(prev => ({ ...prev, [field]: value }))
  }

  function handleUpdateSelect(id: string) {
    setUpdateForm(prev => ({ ...prev, id }))
  }

  function handleResetSelect(id: string) {
    setResetForm(prev => ({ ...prev, id }))
  }

  async function submitCreate(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setCreatePending(true)
    setFeedback(null)
    try {
      const payload = {
        email: createForm.email.trim(),
        display_name: createForm.display_name.trim(),
        password: createForm.password,
        role: 'admin',
      }
      const res = await fetch('/api/admin/access/admin-users', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) {
        const message = await resolveError(res, `创建失败（${res.status}）`)
        throw new Error(message)
      }
      setCreateForm(initialCreate)
      setFeedback({ kind: 'success', message: '系统用户已创建，默认状态为启用。' })
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '创建系统用户失败' })
    } finally {
      setCreatePending(false)
    }
  }

  async function submitUpdate(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!updateForm.id) {
      setFeedback({ kind: 'error', message: '请选择需要更新的用户。' })
      return
    }
    setUpdatePending(true)
    setFeedback(null)
    try {
      const res = await fetch(`/api/admin/access/admin-users/${updateForm.id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          display_name: updateForm.display_name.trim(),
          role: updateForm.role,
          active: updateForm.active,
        }),
      })
      if (!res.ok) {
        const message = await resolveError(res, `更新失败（${res.status}）`)
        throw new Error(message)
      }
      setFeedback({ kind: 'success', message: '系统用户信息已更新。' })
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '更新系统用户失败' })
    } finally {
      setUpdatePending(false)
    }
  }

  async function submitReset(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!resetForm.id) {
      setFeedback({ kind: 'error', message: '请选择需要重置密码的用户。' })
      return
    }
    setResetPending(true)
    setFeedback(null)
    try {
      const res = await fetch(`/api/admin/access/admin-users/${resetForm.id}/reset-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password: resetForm.password }),
      })
      if (!res.ok) {
        const message = await resolveError(res, `重置密码失败（${res.status}）`)
        throw new Error(message)
      }
      setResetForm(initialReset)
      setFeedback({ kind: 'success', message: '密码已重置，新密码立即生效。' })
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '重置密码失败' })
    } finally {
      setResetPending(false)
    }
  }

  async function handleToggleActive(id: string, nextValue: boolean) {
    setTogglingId(id)
    setFeedback(null)
    try {
      const res = await fetch(`/api/admin/access/admin-users/${id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ active: nextValue }),
      })
      if (!res.ok) {
        const message = await resolveError(res, `更新状态失败（${res.status}）`)
        throw new Error(message)
      }
      setFeedback({ kind: 'success', message: `已${nextValue ? '启用' : '停用'}该系统用户。` })
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '更新用户状态失败' })
    } finally {
      setTogglingId(null)
    }
  }

  return (
    <div className="grid gap-6">
      <div className="grid items-stretch gap-6 lg:grid-cols-2">
        <Card className="flex h-full flex-col">
          <CardHeader>
            <CardTitle>创建系统用户</CardTitle>
            <CardDescription>用于访问 Plane 集成后台的管理员账号，默认角色为 admin。</CardDescription>
          </CardHeader>
          <CardContent className="flex-1">
            <form id="form-create-admin" className="grid gap-4" onSubmit={submitCreate}>
              <div className="grid gap-2">
                <Label htmlFor="admin_email">邮箱</Label>
                <Input
                  id="admin_email"
                  placeholder="admin@example.com"
                  value={createForm.email}
                  onChange={event => handleCreateChange('email', event.target.value)}
                  type="email"
                  required
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="admin_display_name">显示名</Label>
                <Input
                  id="admin_display_name"
                  placeholder="如：集成管理员"
                  value={createForm.display_name}
                  onChange={event => handleCreateChange('display_name', event.target.value)}
                  required
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="admin_password">初始密码</Label>
                <Input
                  id="admin_password"
                  type="password"
                  placeholder="至少 8 位，包含字母数字"
                  value={createForm.password}
                  onChange={event => handleCreateChange('password', event.target.value)}
                  required
                  minLength={8}
                />
              </div>
              
            </form>
          </CardContent>
          <CardFooter className="flex flex-col gap-2 text-xs text-muted-foreground sm:flex-row sm:items-center sm:justify-between">
            <span>账号创建后可在下方列表中查看并调整状态。</span>
            <Button type="submit" form="form-create-admin" disabled={createPending}>
              {createPending ? '创建中…' : '创建用户'}
            </Button>
          </CardFooter>
        </Card>

        <Card className="flex h-full flex-col">
          <CardHeader>
            <CardTitle>更新显示名 / 角色 / 状态</CardTitle>
            <CardDescription>选择已存在的系统用户并调整资料，角色当前仅支持 admin。</CardDescription>
          </CardHeader>
          <CardContent className="flex-1">
            <form id="form-update-admin" className="grid gap-4" onSubmit={submitUpdate}>
              <div className="grid gap-2">
                <Label htmlFor="update_user">选择用户</Label>
                <Select
                  id="update_user"
                  value={updateForm.id}
                  onChange={event => handleUpdateSelect(event.target.value)}
                >
                  <option value="">请选择系统用户</option>
                  {userOptions.map(option => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </Select>
              </div>
              <div className="grid gap-2">
                <Label htmlFor="update_display_name">显示名</Label>
                <Input
                  id="update_display_name"
                  placeholder="如：集成管理员"
                  value={updateForm.display_name}
                  onChange={event => setUpdateForm(prev => ({ ...prev, display_name: event.target.value }))}
                  disabled={!updateForm.id}
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="update_role">角色</Label>
                <Select
                  id="update_role"
                  value={updateForm.role}
                  onChange={event => setUpdateForm(prev => ({ ...prev, role: event.target.value }))}
                  disabled={!updateForm.id}
                >
                  <option value="admin">admin</option>
                </Select>
              </div>
              <div className="flex flex-col gap-3 rounded-xl border border-dashed border-border px-4 py-3 sm:flex-row sm:items-center sm:justify-between">
                <div className="flex flex-col text-sm">
                  <span className="font-medium text-foreground">账号启用状态</span>
                  <span className="text-xs text-muted-foreground">停用后不可登录，已登录会立即失效。</span>
                </div>
                <Switch
                  checked={updateForm.active}
                  onChange={event => setUpdateForm(prev => ({ ...prev, active: event.target.checked }))}
                  disabled={!updateForm.id}
                />
              </div>
              
            </form>
          </CardContent>
          <CardFooter className="flex items-center justify-end">
            <Button type="submit" form="form-update-admin" disabled={updatePending || !updateForm.id}>
              {updatePending ? '更新中…' : '保存更新'}
            </Button>
          </CardFooter>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>重置系统用户密码</CardTitle>
          <CardDescription>选择目标账号并设置新的登录密码，执行后立即生效。</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="grid gap-4 md:grid-cols-[2fr,2fr,1fr]" onSubmit={submitReset}>
            <div className="grid gap-2">
              <Label htmlFor="reset_user">系统用户</Label>
              <Select
                id="reset_user"
                value={resetForm.id}
                onChange={event => handleResetSelect(event.target.value)}
              >
                <option value="">请选择系统用户</option>
                {userOptions.map(option => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </Select>
            </div>
            <div className="grid gap-2">
              <Label htmlFor="reset_password">新密码</Label>
              <Input
                id="reset_password"
                type="password"
                placeholder="至少 8 位"
                value={resetForm.password}
                onChange={event => setResetForm(prev => ({ ...prev, password: event.target.value }))}
                minLength={8}
                required
              />
            </div>
            <div className="flex items-end">
              <Button type="submit" disabled={resetPending || !resetForm.id} className="w-full">
                {resetPending ? '重置中…' : '重置密码'}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {feedback && <Alert variant={feedback.kind === 'success' ? 'success' : 'destructive'}>{feedback.message}</Alert>}
      {fetchError && <Alert variant="destructive">{fetchError}</Alert>}

      <Card>
        <CardHeader>
          <CardTitle>系统用户列表</CardTitle>
          <CardDescription>查看管理员账号的启用状态、登录记录与创建时间。</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex flex-col gap-2 text-sm text-muted-foreground sm:flex-row sm:items-center sm:justify-between">
            <span>共 {users.length} 个系统用户，{users.filter(user => user.active).length} 个处于启用状态。</span>
            <Badge variant="info">实时数据</Badge>
          </div>
          <div className="overflow-x-auto">
            <Table className="min-w-[720px]">
              <TableHeader>
                <TableRow>
                  <TableHead>显示名</TableHead>
                  <TableHead>邮箱</TableHead>
                  <TableHead>角色</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>最近登录</TableHead>
                  <TableHead>创建时间</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {users.map(user => (
                  <TableRow key={user.id}>
                    <TableCell className="font-medium text-foreground">{user.display_name || '—'}</TableCell>
                    <TableCell>{user.email}</TableCell>
                    <TableCell>
                      <Badge variant="muted">{user.role}</Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <Switch
                          checked={user.active}
                          onChange={event => handleToggleActive(user.id, event.target.checked)}
                          disabled={togglingId === user.id}
                        />
                        <span className="text-xs text-muted-foreground">{user.active ? '启用中' : '已停用'}</span>
                      </div>
                    </TableCell>
                    <TableCell>{formatDate(user.last_login_at)}</TableCell>
                    <TableCell>{formatDate(user.created_at)}</TableCell>
                  </TableRow>
                ))}
                {!loading && users.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-sm text-muted-foreground">
                      尚未创建系统用户。
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

async function resolveError(res: Response, fallback: string) {
  const proxyMessage = res.headers.get('x-proxy-error')
  if (proxyMessage) return proxyMessage
  const ct = res.headers.get('content-type') || ''
  const text = await res.text()
  if (!text) return fallback
  if (ct.includes('application/json')) {
    try {
      const json = JSON.parse(text)
      if (json?.error?.message) return String(json.error.message)
      if (json?.message) return String(json.message)
    } catch {}
  }
  const snippet = text.trim().replace(/\s+/g, ' ')
  if (snippet.startsWith('<')) return fallback
  return snippet.length > 200 ? snippet.slice(0, 200) + '…' : snippet
}

function formatDate(value?: string | null) {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN', { hour12: false })
}
