'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import type { ReactNode } from 'react'
import type { Route } from 'next'

import { cn } from '@/lib/utils'

export type NavItem = { href: Route; label: string; icon: ReactNode }

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
      : 'flex flex-nowrap items-center gap-2 overflow-x-auto overflow-y-hidden whitespace-nowrap no-scrollbar'

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
            aria-label={orientation === 'horizontal' ? item.label : undefined}
          >
            <span className="nav-link__icon" aria-hidden="true">
              {item.icon}
            </span>
            <span
              className={cn(
                orientation === 'vertical' ? 'not-sr-only' : 'sr-only md:not-sr-only',
                'nav-link__label',
              )}
            >
              {item.label}
            </span>
          </Link>
        )
      })}
    </div>
  )
}
