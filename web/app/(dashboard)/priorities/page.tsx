'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'

import { Alert } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

type Feedback = { kind: 'success' | 'error'; message: string }

type MappingItem = {
  id?: number
  scope_kind: string
  scope_id?: string | null
  mapping_type: string
  left: { system: string; type: string; key: string }
  right: { system: string; type: string; key: string }
  bidirectional: boolean
  active: boolean
  created_at?: string
  updated_at?: string
}

const planePriorityKeys = ['urgent', 'high', 'medium', 'low', 'none'] as const
const cnbPriorityOptions = [
  { value: '', label: '空（无）' },
  { value: 'P3', label: 'P3（低）' },
  { value: 'P2', label: 'P2（中）' },
  { value: 'P1', label: 'P1（高）' },
  { value: 'P0', label: 'P0（紧急）' },
  { value: '-1P', label: '-1P（负优先级1）' },
  { value: '-2P', label: '-2P（负优先级2）' },
]

type FormState = {
  plane_project_id: string
  values: Record<(typeof planePriorityKeys)[number], string>
}

const initialForm: FormState = {
  plane_project_id: '',
  values: {
    urgent: 'P0',
    high: 'P1',
    medium: 'P2',
    low: 'P3',
    none: '',
  },
}

export default function PriorityMappingsPage() {
  const [form, setForm] = useState<FormState>(initialForm)
  const [loading, setLoading] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [feedback, setFeedback] = useState<Feedback | null>(null)
  const [rows, setRows] = useState<MappingItem[]>([])
  const [globalDefaults, setGlobalDefaults] = useState<Record<string, string>>({})

  const projectQuery = useMemo(() => {
    const params = new URLSearchParams()
    params.set('scope_kind', 'plane_project')
    if (form.plane_project_id.trim()) params.set('scope_id', form.plane_project_id.trim())
    params.set('mapping_type', 'priority')
    return params.toString()
  }, [form.plane_project_id])

  const load = useCallback(async () => {
    if (!form.plane_project_id.trim()) {
      setRows([])
      setFetchError('请先填写 Plane Project ID。')
      return
    }
    setLoading(true)
    setFetchError(null)
    try {
      // 1) 项目级映射
      const res = await fetch(`/api/admin/mappings?${projectQuery}`, { cache: 'no-store' })
      if (!res.ok) throw new Error(`查询失败（${res.status}）`)
      const data = await res.json()
      const items = Array.isArray(data.items) ? (data.items as MappingItem[]) : []
      setRows(items)
      // 初始化表单值（如项目级无配置，则用全局默认）
      if (items.length > 0) {
        const next: Record<string, string> = { ...initialForm.values }
        for (const it of items) {
          const k = it.left?.key?.toLowerCase()
          const v = it.right?.key ?? ''
          if (k && planePriorityKeys.includes(k as any)) next[k] = v
        }
        setForm(prev => ({ ...prev, values: next }))
      } else {
        // 2) 拉取全局默认
        const g = await fetch(`/api/admin/mappings?${new URLSearchParams({ scope_kind: 'global', mapping_type: 'priority' }).toString()}`, { cache: 'no-store' })
        if (g.ok) {
          const gj = await g.json()
          const next: Record<string, string> = { ...initialForm.values }
          const defs: Record<string, string> = {}
          for (const it of gj.items as MappingItem[]) {
            const k = it.left?.key?.toLowerCase()
            const v = it.right?.key ?? ''
            if (k && planePriorityKeys.includes(k as any)) {
              next[k] = v
              defs[k] = v
            }
          }
          setForm(prev => ({ ...prev, values: next }))
          setGlobalDefaults(defs)
        }
      }
    } catch (err) {
      setFetchError(err instanceof Error ? err.message : '加载失败')
    } finally {
      setLoading(false)
    }
  }, [projectQuery, form.plane_project_id])

  useEffect(() => {
    // 仅当填写了 project 时加载
    if (form.plane_project_id.trim()) {
      load()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  async function submit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setFeedback(null)
    if (!form.plane_project_id.trim()) {
      setFeedback({ kind: 'error', message: '请先填写 Plane Project ID。' })
      return
    }
    try {
      const items = planePriorityKeys.map(k => ({
        left: { system: 'plane', type: 'priority', key: k },
        right: { system: 'cnb', type: 'priority', key: form.values[k] ?? '' },
        bidirectional: true,
        active: true,
      }))
      const payload = {
        scope_kind: 'plane_project',
        scope_id: form.plane_project_id.trim(),
        mapping_type: 'priority',
        items,
      }
      const res = await fetch('/api/admin/mappings', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) {
        const text = await res.text()
        throw new Error(text || `保存失败（${res.status}）`)
      }
      setFeedback({ kind: 'success', message: '优先级映射已保存。' })
      await load()
    } catch (err) {
      setFeedback({ kind: 'error', message: err instanceof Error ? err.message : '保存失败，请稍后再试。' })
    }
  }

  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle>Plane ↔ CNB 优先级映射</CardTitle>
          <CardDescription>默认全局映射已预置，可按项目进行覆盖配置。</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="grid gap-6" onSubmit={submit}>
            <div className="grid gap-2 md:max-w-sm">
              <Label htmlFor="plane_project_id">Plane Project ID</Label>
              <div className="flex items-center gap-2">
                <Input
                  id="plane_project_id"
                  placeholder="project uuid"
                  value={form.plane_project_id}
                  onChange={e => setForm(prev => ({ ...prev, plane_project_id: e.target.value }))}
                />
                <Button type="button" variant="outline" onClick={() => form.plane_project_id.trim() && load()} disabled={loading}>
                  {loading ? '加载中…' : '加载'}
                </Button>
              </div>
            </div>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {planePriorityKeys.map(key => (
                <div key={key} className="grid gap-2">
                  <Label htmlFor={`cnb_${key}`}>Plane: {key}</Label>
                  <Select
                    id={`cnb_${key}`}
                    value={form.values[key]}
                    onChange={e => setForm(prev => ({ ...prev, values: { ...prev.values, [key]: e.target.value } }))}
                  >
                    {cnbPriorityOptions.map(opt => (
                      <option key={opt.value} value={opt.value}>
                        {opt.label}
                      </option>
                    ))}
                  </Select>
                  {globalDefaults[key] && (
                    <span className="text-xs text-muted-foreground">默认：{globalDefaults[key]}</span>
                  )}
                </div>
              ))}
            </div>
            {feedback && <Alert variant={feedback.kind === 'success' ? 'success' : 'destructive'}>{feedback.message}</Alert>}
            <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
              <span className="text-sm text-muted-foreground">保存后将覆盖该 Plane 项目下的优先级映射（若未填写则按全局默认）。</span>
              <Button type="submit">保存映射</Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <div className="space-y-1">
              <CardTitle>当前项目的映射规则</CardTitle>
              <CardDescription>显示项目级自定义（若为空，则实际生效为全局默认）。</CardDescription>
            </div>
            <div className="flex items-center gap-2">
              <Badge variant="outline">{rows.length} 条</Badge>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {fetchError && <Alert variant="destructive">{fetchError}</Alert>}
          <div className="overflow-x-auto">
            <Table className="min-w-[640px]">
              <TableHeader>
                <TableRow className="bg-transparent">
                  <TableHead>左侧（Plane）</TableHead>
                  <TableHead>右侧（CNB）</TableHead>
                  <TableHead>双向</TableHead>
                  <TableHead>状态</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {rows.map(item => (
                  <TableRow key={item.id ?? `${item.left.key}-${item.right.key}`}>
                    <TableCell>
                      <div className="text-sm text-foreground">{item.left.type}:{' '}{item.left.key}</div>
                      <div className="text-xs text-muted-foreground">system: {item.left.system}</div>
                    </TableCell>
                    <TableCell>
                      <div className="text-sm text-foreground">{item.right.type}:{' '}{item.right.key}</div>
                      <div className="text-xs text-muted-foreground">system: {item.right.system}</div>
                    </TableCell>
                    <TableCell>
                      {item.bidirectional ? <Badge variant="success">是</Badge> : <Badge variant="muted">否</Badge>}
                    </TableCell>
                    <TableCell>
                      {item.active ? <Badge variant="success">启用</Badge> : <Badge variant="muted">停用</Badge>}
                    </TableCell>
                  </TableRow>
                ))}
                {rows.length === 0 && !loading && (
                  <TableRow>
                    <TableCell colSpan={4} className="text-center text-sm text-muted-foreground">
                      该项目暂无自定义优先级映射，当前使用全局默认。
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
        <CardFooter className="text-sm text-muted-foreground">
          映射数据来自数据库 <code>integration_mappings</code>（mapping_type='priority'）。
        </CardFooter>
      </Card>
    </div>
  )
}

