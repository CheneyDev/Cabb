'use client'

import type { Route } from 'next'
import { Suspense, useMemo, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'

import { Alert } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

function LoginContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const nextParam = searchParams.get('next') || '/'
  const target = useMemo<Route>(() => {
    return nextParam.startsWith('/') ? (nextParam as Route) : ('/' as Route)
  }, [nextParam])

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [pending, setPending] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setPending(true)
    setError(null)
    try {
      const payload = { email: email.trim(), password }
      const res = await fetch('/api/admin/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) {
        const message = await resolveError(res, `登录失败（${res.status}）`)
        throw new Error(message)
      }
      router.push(target)
      router.refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败，请稍后再试')
    } finally {
      setPending(false)
    }
  }

  return (
    <div className="auth-surface">
      <Card className="auth-card">
        <CardHeader className="space-y-4 text-center">
          <div className="mx-auto inline-flex h-12 w-12 items-center justify-center rounded-3xl bg-[color-mix(in_srgb,var(--primary)_22%,transparent)] text-lg font-semibold uppercase tracking-wider text-primary-foreground shadow-[0_12px_35px_-18px_rgba(79,70,229,0.65)]">
            PI
          </div>
          <div className="space-y-2">
            <CardTitle className="text-2xl font-semibold tracking-tight text-foreground">Plane 集成后台登录</CardTitle>
            <CardDescription className="text-sm text-muted-foreground">
              使用分配的管理员邮箱与密码进入 Plane ↔ CNB 集成配置中心。
            </CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          <form className="grid gap-4" onSubmit={submit}>
            <div className="grid gap-2 text-left">
              <Label htmlFor="login_email">管理员邮箱</Label>
              <Input
                id="login_email"
                type="email"
                placeholder="admin@example.com"
                value={email}
                onChange={event => setEmail(event.target.value)}
                required
              />
            </div>
            <div className="grid gap-2 text-left">
              <Label htmlFor="login_password">密码</Label>
              <Input
                id="login_password"
                type="password"
                placeholder="请输入密码"
                value={password}
                onChange={event => setPassword(event.target.value)}
                required
                minLength={8}
              />
            </div>
            <div className="flex items-center justify-between text-xs text-muted-foreground">
              <span>首次使用请联系集成负责人获取账号。</span>
              <span>支持 SSO 整合预留位</span>
            </div>
            {error && <Alert variant="destructive">{error}</Alert>}
            <Button type="submit" disabled={pending} className="w-full">
              {pending ? '登录中…' : '登录'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}

export default function LoginPage() {
  return (
    <Suspense
      fallback={
        <div className="auth-surface">
          <Card className="auth-card">
            <CardContent className="text-center text-sm text-muted-foreground">正在加载登录表单…</CardContent>
          </Card>
        </div>
      }
    >
      <LoginContent />
    </Suspense>
  )
}

async function resolveError(res: Response, fallback: string) {
  const proxyMessage = res.headers.get('x-proxy-error')
  if (proxyMessage) {
    return proxyMessage
  }
  const text = await res.text()
  if (!text) return fallback
  try {
    const json = JSON.parse(text)
    if (json?.error?.message) return String(json.error.message)
    if (json?.message) return String(json.message)
    return fallback
  } catch (err) {
    return text
  }
}
