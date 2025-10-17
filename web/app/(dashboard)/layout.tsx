import Link from 'next/link'
import { redirect } from 'next/navigation'
import type { ReactNode } from 'react'

import { buttonVariants } from '@/components/ui/button-variants'
import type { AdminUser } from '@/lib/server/admin-session'
import { fetchAdminSession } from '@/lib/server/admin-session'

import { DashboardNav, type NavItem } from './_components/nav-links'
import { LogoutButton } from './_components/logout-button'
import { MobileNav } from './_components/MobileNav'
import { ThemeToggle } from '@/components/theme-toggle'

const iconProps = {
  viewBox: '0 0 24 24',
  fill: 'none',
  stroke: 'currentColor',
  strokeWidth: 1.8,
  strokeLinecap: 'round',
  strokeLinejoin: 'round',
} as const

const navItems: NavItem[] = [
  {
    href: '/',
    label: '概览',
    icon: (
      <svg aria-hidden="true" {...iconProps}>
        <path d="M4 5.5h7v6H4z" />
        <path d="M13 5.5h7v4H13z" />
        <path d="M13 11.5h7v7.5H13z" />
        <path d="M4 13.5h7v5.5H4z" />
      </svg>
    ),
  },
  {
    href: '/mappings',
    label: 'Repo↔Project 同步',
    icon: (
      <svg aria-hidden="true" {...iconProps}>
        <path d="M6 8h8" />
        <path d="M11.5 5.5 14 8l-2.5 2.5" />
        <path d="M18 16h-8" />
        <path d="M12.5 18.5 10 16l2.5-2.5" />
        <rect x="4" y="4" width="4" height="4" rx="1" />
        <rect x="16" y="16" width="4" height="4" rx="1" />
      </svg>
    ),
  },
  {
    href: '/users',
    label: '用户映射管理',
    icon: (
      <svg aria-hidden="true" {...iconProps}>
        <circle cx="12" cy="9" r="3.2" />
        <path d="M6.5 19.5c.7-2.6 2.9-4.5 5.5-4.5s4.8 1.9 5.5 4.5" />
        <path d="M4.5 19.5c.4-1.8 1.8-3.2 3.6-3.8" />
        <path d="M19.5 19.5c-.4-1.8-1.8-3.2-3.6-3.8" />
      </svg>
    ),
  },
  {
    href: '/admin-users',
    label: '系统用户管理',
    icon: (
      <svg aria-hidden="true" {...iconProps}>
        <path d="M12 4.5 5 8v4c0 4.1 2.9 7.8 7 8.8 4.1-1 7-4.7 7-8.8V8z" />
        <circle cx="12" cy="11" r="2.7" />
      </svg>
    ),
  },
]

export default async function DashboardLayout({ children }: { children: ReactNode }) {
  const session = await fetchAdminSession()
  if (session.status === 401 || session.status === 403) {
    redirect('/login')
  }
  const user = (session.user ?? null) as AdminUser | null
  if (!user) {
    const message = session.message ?? '暂时无法获取管理员信息，请稍后重试或联系集成负责人。'
    const detail = session.status ? `（状态码 ${session.status}）` : ''
    return (
      <div className="flex min-h-dvh flex-col items-center justify-center bg-[color-mix(in_srgb,var(--background)_92%,transparent)] px-6">
        <div className="mx-auto w-full max-w-md space-y-5 rounded-3xl border border-[color-mix(in_srgb,var(--border)_75%,transparent)] bg-card p-8 text-center shadow-[0_30px_120px_-45px_rgba(79,70,229,0.4)]">
          <div className="space-y-2">
            <h1 className="text-2xl font-semibold tracking-tight text-foreground">管理员鉴权不可用</h1>
            <p className="text-sm leading-relaxed text-muted-foreground">{message}{detail}</p>
          </div>
          <div className="flex flex-col items-center gap-3">
            <Link href="/login" className={buttonVariants({ variant: 'primary' })}>
              返回登录页
            </Link>
            <p className="text-xs text-muted-foreground">若问题持续，请确认后端 `/admin/auth/me` 接口可用并检查网络。</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-dvh">
      <header className="header-surface sticky top-0 z-40">
        <div className="mx-auto flex max-w-6xl flex-col gap-4 px-4 py-4 sm:px-6 md:flex-row md:items-center md:justify-between">
          <div className="flex items-center justify-between gap-3 md:justify-start">
            <Link href="/" className="flex items-center gap-3">
              <span className="inline-flex h-9 w-9 items-center justify-center rounded-2xl bg-[color-mix(in_srgb,var(--primary)_25%,transparent)] text-sm font-semibold uppercase tracking-wide text-primary-foreground shadow-[0_12px_32px_-22px_rgba(79,70,229,0.8)]">
                PI
              </span>
              <div className="flex flex-col text-left">
                <span className="text-sm font-semibold tracking-wide text-foreground">Plane Integration</span>
                <span className="text-xs text-muted-foreground">后台配置中心 · COSS UI</span>
              </div>
            </Link>
            <div className="flex items-center gap-2 md:hidden">
              <MobileNav items={navItems} user={user} />
              <ThemeToggle />
              <LogoutButton />
            </div>
          </div>
          <nav className="hidden w-full md:flex md:flex-1 md:justify-center">
            <DashboardNav items={navItems} />
          </nav>
          <div className="hidden w-full min-w-0 flex-1 items-center justify-end gap-3 md:flex">
            <div className="flex min-w-0 max-w-full flex-col text-right">
              <span className="truncate text-sm font-semibold text-foreground">{user.display_name || user.email}</span>
              <span className="truncate text-xs text-muted-foreground">{user.email}</span>
            </div>
            <ThemeToggle />
            <LogoutButton />
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-6xl px-4 pb-10 pt-8 sm:px-6">{children}</main>
    </div>
  )
}
