'use client'

import { useEffect } from 'react'

import { Button } from '@/components/ui/button'

export default function GlobalError({ error, reset }: { error: Error; reset: () => void }) {
  useEffect(() => {
    // 把原始错误打印到控制台，便于线上排查
    // 注意：不要在 UI 中暴露敏感信息
    // eslint-disable-next-line no-console
    console.error('[GlobalError]', error)
  }, [error])

  return (
    <html>
      <body>
        <div className="flex min-h-dvh flex-col items-center justify-center gap-6 bg-background px-6 text-center">
          <div className="space-y-2">
            <h1 className="text-2xl font-semibold tracking-tight text-foreground">应用发生异常</h1>
            <p className="text-sm leading-relaxed text-muted-foreground">
              页面加载时遇到客户端错误。请刷新重试或返回首页。
            </p>
          </div>
          <div className="flex items-center gap-3">
            <Button variant="secondary" onClick={() => (window.location.href = '/')}>返回首页</Button>
            <Button onClick={() => reset()}>刷新</Button>
          </div>
          <details className="max-w-3xl rounded-2xl border border-border bg-card p-4 text-left text-xs text-muted-foreground">
            <summary className="cursor-pointer text-foreground">错误详情（仅供排查）</summary>
            <pre className="mt-3 whitespace-pre-wrap break-all font-mono">{String(error?.stack || error?.message || error)}</pre>
          </details>
        </div>
      </body>
    </html>
  )
}

