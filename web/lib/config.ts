export function getAPIBase() {
  // Server-side env (for route handlers)
  const serverBase = process.env.API_BASE
  if (serverBase) return serverBase.replace(/\/$/, '')
  // Client-side env
  if (typeof window !== 'undefined') {
    const clientBase = process.env.NEXT_PUBLIC_API_BASE || ''
    if (clientBase) return clientBase.replace(/\/$/, '')
  }
  // Fallback to local backend
  return 'http://localhost:8080'
}

