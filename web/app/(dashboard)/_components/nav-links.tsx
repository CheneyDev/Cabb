'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import type { Route } from 'next'

import { cn } from '@/lib/utils'

export type NavItem = { href: Route; label: string }

type DashboardNavProps = {
  items: NavItem[]
  orientation?: 'horizontal' | 'vertical'
  onNavigate?: () => void
}

export function DashboardNav({ items, orientation = 'horizontal', onNavigate }: DashboardNavProps) {
  const pathname = usePathname()
  const containerClass =
    orientation === 'vertical'
      ? 'flex flex-col gap-1'
      : 'flex flex-wrap items-center gap-2'

  return (
    <div className={containerClass}>
      {items.map(item => {
        const active = pathname === item.href || (item.href !== '/' && pathname.startsWith(`${item.href}/`))
        return (
          <Link
            key={item.href}
            href={item.href}
            onClick={() => onNavigate?.()}
            className={cn('nav-link', orientation === 'vertical' && 'nav-link--vertical', active && 'nav-link--active')}
          >
            {item.label}
          </Link>
        )
      })}
    </div>
  )
}
