'use client'

import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'

function SunIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true" {...props}>
      <circle cx="12" cy="12" r="4" />
      <path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41" />
    </svg>
  )
}

function MoonIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden="true" {...props}>
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
    </svg>
  )
}

export function ThemeToggle() {
  const [mounted, setMounted] = useState(false)
  const [isDark, setIsDark] = useState(true) // default dark to match current UI

  useEffect(() => {
    setMounted(true)
    try {
      const saved = localStorage.getItem('theme')
      if (saved === 'light') {
        setIsDark(false)
        document.documentElement.classList.remove('dark')
      } else if (saved === 'dark') {
        setIsDark(true)
        document.documentElement.classList.add('dark')
      } else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        setIsDark(true)
        document.documentElement.classList.add('dark')
      } else {
        setIsDark(false)
        document.documentElement.classList.remove('dark')
      }
    } catch {
      // ignore
    }
  }, [])

  const toggle = () => {
    const next = !isDark
    setIsDark(next)
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

  const label = isDark ? '切换为浅色模式' : '切换为深色模式'

  if (!mounted) {
    return (
      <Button variant="ghost" size="sm" aria-hidden disabled>
        <MoonIcon className="h-4 w-4" />
      </Button>
    )
  }

  return (
    <Button variant="ghost" size="sm" aria-label={label} title={label} onClick={toggle}>
      {isDark ? <MoonIcon className="h-4 w-4" /> : <SunIcon className="h-4 w-4" />}
    </Button>
  )
}
