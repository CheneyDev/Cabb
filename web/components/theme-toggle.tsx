'use client'

import { useEffect, useState } from 'react'
import { Switch } from '@/components/ui/switch'

export function ThemeToggle() {
  const [mounted, setMounted] = useState(false)
  const [checked, setChecked] = useState(true) // default dark to match current UI

  useEffect(() => {
    setMounted(true)
    try {
      const saved = localStorage.getItem('theme')
      if (saved === 'light') {
        setChecked(false)
        document.documentElement.classList.remove('dark')
      } else if (saved === 'dark') {
        setChecked(true)
        document.documentElement.classList.add('dark')
      } else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        setChecked(true)
        document.documentElement.classList.add('dark')
      } else {
        setChecked(false)
        document.documentElement.classList.remove('dark')
      }
    } catch {
      // ignore
    }
  }, [])

  const toggle = (next: boolean) => {
    setChecked(next)
    try {
      if (next) {
        document.documentElement.classList.add('dark')
        localStorage.setItem('theme', 'dark')
      } else {
        document.documentElement.classList.remove('dark')
        localStorage.setItem('theme', 'light')
      }
    } catch {
      // ignore storage errors
    }
  }

  // Avoid hydration mismatch: render switch only after mount
  if (!mounted) {
    return (
      <div className="flex items-center gap-2 text-xs text-muted-foreground" aria-hidden>
        主题
      </div>
    )
  }

  return (
    <div className="flex items-center gap-2">
      <span className="text-xs text-muted-foreground">浅色</span>
      <Switch
        aria-label="切换深色模式"
        checked={checked}
        onChange={(e) => toggle((e.target as HTMLInputElement).checked)}
      />
      <span className="text-xs text-muted-foreground">深色</span>
    </div>
  )
}

