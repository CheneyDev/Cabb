import type { ReactNode } from 'react'
import { redirect } from 'next/navigation'

import { fetchAdminSession } from '@/lib/server/admin-session'

export default async function AuthLayout({ children }: { children: ReactNode }) {
  const session = await fetchAdminSession()
  if (session.user && session.status === 200) {
    redirect('/')
  }
  return <>{children}</>
}
