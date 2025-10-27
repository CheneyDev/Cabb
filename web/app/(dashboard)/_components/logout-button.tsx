'use client'

import { useState, useTransition } from 'react'
import { useRouter } from 'next/navigation'

import { Button } from '@/components/ui/button'

function LogoutIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.8"
      strokeLinecap="round"
      aria-hidden="true"
      {...props}
    >
      <path d="M16 17l5-5-5-5" />
      <path d="M21 12H9" />
      <path d="M12 19v2a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V3a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
    </svg>
  )
}

export function LogoutButton() {
  const router = useRouter()
  const [pending, startTransition] = useTransition()
  const [error, setError] = useState<string | null>(null)

  function handleSignOut() {
    setError(null)
    startTransition(async () => {
      try {
        const res = await fetch('/api/admin/auth/logout', { method: 'POST' })
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
          throw new Error(msg || `退出失败（${res.status}）`)
        }
        router.push('/login')
        router.refresh()
      } catch (err) {
        setError(err instanceof Error ? err.message : '退出失败，请稍后重试')
      }
    })
  }

  const label = pending ? '正在退出…' : '退出登录'

  return (
    <div className="flex flex-col items-end gap-1">
      {error && <span className="text-xs text-destructive">{error}</span>}
      <Button
        variant="ghost"
        size="sm"
        onClick={handleSignOut}
        disabled={pending}
        aria-label={label}
        title={label}
      >
        <LogoutIcon className="h-4 w-4" />
        <span className="sr-only">{label}</span>
      </Button>
    </div>
  )
}
