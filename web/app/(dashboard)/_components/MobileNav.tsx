'use client'

import { useEffect, useRef, useState } from 'react'
import { usePathname } from 'next/navigation'

import { DashboardNav, type NavItem } from './nav-links'

type MobileNavProps = {
  items: NavItem[]
  user: { display_name?: string | null; email: string }
}

export function MobileNav({ items, user }: MobileNavProps) {
  const [open, setOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)
  const pathname = usePathname()

  useEffect(() => {
    setOpen(false)
  }, [pathname])

  useEffect(() => {
    if (!open) return

    const handleOutsideClick = (event: MouseEvent | TouchEvent) => {
      if (!containerRef.current) return
      const target = event.target
      if (!(target instanceof Node)) return
      if (!containerRef.current.contains(target)) {
        setOpen(false)
      }
    }

    document.addEventListener('mousedown', handleOutsideClick)
    document.addEventListener('touchstart', handleOutsideClick)

    return () => {
      document.removeEventListener('mousedown', handleOutsideClick)
      document.removeEventListener('touchstart', handleOutsideClick)
    }
  }, [open])

  return (
    <div ref={containerRef} className="relative md:hidden">
      <button
        type="button"
        onClick={() => setOpen(prev => !prev)}
        aria-expanded={open}
        aria-label="切换菜单"
        aria-controls="mobile-dashboard-nav"
        aria-haspopup="menu"
        className="inline-flex items-center gap-2 rounded-full border border-[color-mix(in_srgb,var(--border)_65%,transparent)] bg-[color-mix(in_srgb,var(--background)_85%,transparent)] px-3 py-1.5 text-sm font-medium text-foreground shadow-[0_18px_36px_-30px_rgba(79,70,229,0.8)] transition"
      >
        <svg
          aria-hidden="true"
          focusable="false"
          className="h-4 w-4 text-muted-foreground"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.8"
          strokeLinecap="round"
        >
          <path d="M4 7h16M4 12h16M4 17h16" />
        </svg>
        <span>{open ? '收起' : '菜单'}</span>
      </button>
      {open && (
        <div
          id="mobile-dashboard-nav"
          role="menu"
          className="absolute right-0 top-full z-50 mt-3 min-w-[16rem] max-w-[92vw]"
        >
          <div className="space-y-3 rounded-3xl border border-[color-mix(in_srgb,var(--border)_75%,transparent)] bg-[color-mix(in_srgb,var(--background)_94%,transparent)] p-4 shadow-[0_32px_90px_-45px_rgba(79,70,229,0.55)]">
            <div className="flex flex-col gap-1">
              <span className="text-sm font-semibold text-foreground">{user.display_name || user.email}</span>
              <span className="break-all text-xs text-muted-foreground">{user.email}</span>
            </div>
            <nav>
              <DashboardNav items={items} orientation="vertical" onNavigate={() => setOpen(false)} />
            </nav>
          </div>
        </div>
      )}
    </div>
  )
}
