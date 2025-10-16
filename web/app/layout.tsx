export const metadata = {
  title: 'Plane Integration Admin',
  description: 'Manage Plane ↔ CNB mappings',
}

import './globals.css'

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <body>
        <div className="min-h-dvh">
          <header className="sticky top-0 z-40 border-b border-border bg-foreground/80 backdrop-blur">
            <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-3">
              <div className="flex items-center gap-2">
                <div className="h-6 w-6 rounded-md bg-gradient-to-br from-primary to-accent" />
                <span className="text-sm font-semibold tracking-wide">Plane Integration</span>
              </div>
              <nav className="flex items-center gap-4 text-sm text-neutral-300">
                <a className="hover:text-white" href="/">概览</a>
                <a className="hover:text-white" href="/mappings">Repo↔Project 同步</a>
              </nav>
            </div>
          </header>
          <main className="mx-auto max-w-6xl px-6 py-6">{children}</main>
        </div>
      </body>
    </html>
  )
}

