import * as React from 'react'

import { cn } from '@/lib/utils'

export interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  containerClassName?: string
}

export const Select = React.forwardRef<HTMLSelectElement, SelectProps>(function Select(
  { className, containerClassName, children, ...props },
  ref,
) {
  return (
    <div className={cn('relative w-full', containerClassName)}>
      <select
        ref={ref}
        className={cn(
          'flex h-10 w-full appearance-none rounded-xl border border-border bg-transparent px-3 text-sm leading-tight text-foreground shadow-[inset_0_1px_0_rgba(255,255,255,0.05)] transition focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40 focus:ring-offset-2 focus:ring-offset-background disabled:cursor-not-allowed disabled:opacity-60',
          className,
        )}
        {...props}
      >
        {children}
      </select>
      <span className="pointer-events-none absolute inset-y-0 right-3 flex items-center text-muted-foreground">▾</span>
    </div>
  )
})
