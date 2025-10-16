import './globals.css'

export const metadata = {
  title: 'Plane Integration Admin',
  description: 'Manage Plane ↔ CNB mappings',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <body className="min-h-dvh bg-background text-foreground antialiased">{children}</body>
    </html>
  )
}
