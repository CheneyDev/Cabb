"use client"

import { useEffect, useMemo, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

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

export default function MappingsPage() {
  const [items, setItems] = useState<Mapping[]>([])
  const [loading, setLoading] = useState(false)
  const [filterProject, setFilterProject] = useState('')
  const [form, setForm] = useState({
    cnb_repo_id: '',
    plane_workspace_id: '',
    plane_project_id: '',
    sync_direction: 'cnb_to_plane',
    label_selector: '*',
    active: true,
  })

  const qs = useMemo(() => {
    const p = new URLSearchParams()
    if (filterProject) p.set('plane_project_id', filterProject)
    return p.toString()
  }, [filterProject])

  async function load() {
    setLoading(true)
    try {
      const res = await fetch(`/api/admin/mappings/repo-project?${qs}`, { cache: 'no-store' })
      const json = await res.json()
      setItems(json.items || [])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { load() // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [qs])

  async function submit(e: React.FormEvent) {
    e.preventDefault()
    const res = await fetch('/api/admin/mappings/repo-project', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form),
    })
    if (res.ok) {
      setForm({ ...form, cnb_repo_id: '' })
      await load()
    } else {
      const t = await res.text()
      alert(`保存失败: ${t}`)
    }
  }

  return (
    <div className="grid gap-6">
      <section className="card p-4">
        <h2 className="mb-3 text-lg font-semibold">新增/更新 Repo↔Project 映射</h2>
        <form onSubmit={submit} className="grid grid-cols-1 gap-3 md:grid-cols-2">
          <label className="grid gap-1">
            <span className="label">CNB 仓库（org/repo）</span>
            <Input required placeholder="1024hub/plane-test" value={form.cnb_repo_id} onChange={e => setForm({ ...form, cnb_repo_id: e.target.value })} />
          </label>
          <label className="grid gap-1">
            <span className="label">Plane Workspace ID</span>
            <Input required placeholder="uuid" value={form.plane_workspace_id} onChange={e => setForm({ ...form, plane_workspace_id: e.target.value })} />
          </label>
          <label className="grid gap-1">
            <span className="label">Plane Project ID</span>
            <Input required placeholder="uuid" value={form.plane_project_id} onChange={e => setForm({ ...form, plane_project_id: e.target.value })} />
          </label>
          <label className="grid gap-1">
            <span className="label">同步方向</span>
            <select className="input" value={form.sync_direction} onChange={e => setForm({ ...form, sync_direction: e.target.value })}>
              <option value="cnb_to_plane">cnb_to_plane</option>
              <option value="bidirectional">bidirectional</option>
            </select>
          </label>
          <label className="grid gap-1">
            <span className="label">标签选择器</span>
            <Input placeholder="* 或 逗号分隔标签名" value={form.label_selector} onChange={e => setForm({ ...form, label_selector: e.target.value })} />
          </label>
          <label className="mt-6 flex items-center gap-2">
            <input type="checkbox" checked={form.active} onChange={e => setForm({ ...form, active: e.target.checked })} />
            <span className="text-sm">Active</span>
          </label>
          <div className="md:col-span-2">
            <Button type="submit" className="w-full md:w-auto">保存映射</Button>
          </div>
        </form>
      </section>

      <section className="card p-4">
        <div className="mb-3 flex items-center justify-between gap-2">
          <h2 className="text-lg font-semibold">映射列表</h2>
          <div className="flex items-center gap-2">
            <Input placeholder="按 Plane Project ID 过滤" value={filterProject} onChange={e => setFilterProject(e.target.value)} />
            <Button variant="outline" onClick={load} disabled={loading}>{loading ? '加载中...' : '刷新'}</Button>
          </div>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="text-neutral-300">
              <tr className="border-b border-border">
                <th className="p-2">CNB Repo</th>
                <th className="p-2">Plane Workspace</th>
                <th className="p-2">Plane Project</th>
                <th className="p-2">方向</th>
                <th className="p-2">选择器</th>
                <th className="p-2">Active</th>
              </tr>
            </thead>
            <tbody>
              {items.map((m, i) => (
                <tr key={`${m.cnb_repo_id}-${m.plane_project_id}-${i}`} className="border-b border-border/60">
                  <td className="p-2">{m.cnb_repo_id}</td>
                  <td className="p-2">{m.plane_workspace_id}</td>
                  <td className="p-2">{m.plane_project_id}</td>
                  <td className="p-2">{m.sync_direction || 'cnb_to_plane'}</td>
                  <td className="p-2">{m.label_selector || ''}</td>
                  <td className="p-2">{m.active ? '✓' : ''}</td>
                </tr>
              ))}
              {items.length === 0 && (
                <tr>
                  <td className="p-3 text-neutral-400" colSpan={6}>暂无数据</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  )
}

