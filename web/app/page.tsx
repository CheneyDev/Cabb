import Link from 'next/link'

export default async function Page() {
  return (
    <div className="grid gap-6">
      <section className="card p-6">
        <h1 className="mb-2 text-xl font-semibold">概览</h1>
        <p className="text-sm text-neutral-300">管理 Plane ↔ CNB 的集成配置、同步策略与映射。</p>
        <div className="mt-4 flex gap-3">
          <Link href="/mappings" className="btn btn-primary">管理 Repo↔Project 映射</Link>
        </div>
      </section>

      <section className="card p-6">
        <h2 className="mb-2 text-lg font-semibold">运行与环境</h2>
        <ul className="list-disc pl-5 text-sm text-neutral-300">
          <li>后端健康检查：<code className="mx-1">/healthz</code></li>
          <li>环境变量（前端）：<code className="mx-1">NEXT_PUBLIC_API_BASE</code> 指向后端，或使用内置 API 代理。</li>
          <li>Render 部署：<code className="mx-1">render.yaml</code> 已包含前后端服务。</li>
        </ul>
      </section>
    </div>
  )
}

