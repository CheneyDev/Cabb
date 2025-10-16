'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import type { Route } from 'next'

import { cn } from '@/lib/utils'

type NavItem = { href: Route; label: string }

export function DashboardNav({ items }: { items: NavItem[] }) {
  const pathname = usePathname()
  return (
    <>
      {items.map(item => {
        const active = pathname === item.href || (item.href !== '/' && pathname.startsWith(`${item.href}/`))
        return (
          <Link key={item.href} href={item.href} className={cn('nav-link', active && 'nav-link--active')}>
            {item.label}
          </Link>
        )
      })}
    </>
  )
}
