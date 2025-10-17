import './globals.css'

export const metadata = {
  title: 'Plane Integration Admin',
  description: 'Manage Plane ↔ CNB mappings',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN" className="dark">
      <body className="min-h-dvh bg-background text-foreground antialiased">
        <script
          dangerouslySetInnerHTML={{
            __html: `(()=>{try{var s=localStorage.getItem('theme');var d=window.matchMedia&&window.matchMedia('(prefers-color-scheme: dark)').matches;var useDark=s? s==='dark': d;var c=document.documentElement.classList;useDark?c.add('dark'):c.remove('dark')}catch(e){}})();`,
          }}
        />
        {children}
      </body>
    </html>
  )
}
