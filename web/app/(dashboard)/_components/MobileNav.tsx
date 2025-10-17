'use client'

import Link from 'next/link'
import { useEffect, useState } from 'react'
import { usePathname } from 'next/navigation'

import { cn } from '@/lib/utils'
import {
  Menu,
  MenuGroup,
  MenuItem,
  MenuPortal,
  MenuPositioner,
  MenuPopup,
  MenuSeparator,
  MenuTrigger,
} from '@/components/ui/menu'

import type { NavItem } from './nav-links'

type MobileNavProps = {
  items: NavItem[]
  user: { display_name?: string | null; email: string }
}

export function MobileNav({ items, user }: MobileNavProps) {
  const [open, setOpen] = useState(false)
  const pathname = usePathname()

  useEffect(() => {
    setOpen(false)
  }, [pathname])

  return (
    <div className="relative md:hidden">
      <Menu open={open} onOpenChange={nextOpen => setOpen(nextOpen)}>
        <MenuTrigger aria-label={open ? '收起菜单' : '展开菜单'} data-open={open} className="min-w-0 px-2">
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
          <span className="sr-only">{open ? '收起菜单' : '展开菜单'}</span>
        </MenuTrigger>
        <MenuPortal>
          <MenuPositioner>
            <MenuPopup className="w-[min(18rem,92vw)] space-y-3 p-4">
              <div className="flex flex-col gap-1">
                <span className="text-sm font-semibold text-foreground">{user.display_name || user.email}</span>
                <span className="break-all text-xs text-muted-foreground">{user.email}</span>
              </div>
              <MenuSeparator className="my-2" />
              <MenuGroup>
                {items.map(item => {
                  const active = pathname === item.href || (item.href !== '/' && pathname.startsWith(`${item.href}/`))
                  return (
                    <MenuItem
                      key={item.href}
                      data-active={active}
                      render={<Link href={item.href} />}
                      onClick={() => setOpen(false)}
                      className={cn(
                        'justify-start gap-3 text-sm font-medium text-foreground/80',
                        active && 'text-primary',
                      )}
                    >
                      <span className="nav-link__icon text-muted-foreground" aria-hidden="true">
                        {item.icon}
                      </span>
                      <span className="truncate">{item.label}</span>
                    </MenuItem>
                  )
                })}
              </MenuGroup>
            </MenuPopup>
          </MenuPositioner>
        </MenuPortal>
      </Menu>
    </div>
  )
}
