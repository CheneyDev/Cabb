import * as React from 'react'

import { cn } from '@/lib/utils'

export function Table({ className, ...props }: React.HTMLAttributes<HTMLTableElement>) {
  return <table className={cn('w-full border-separate border-spacing-y-3 text-sm', className)} {...props} />
}

export function TableHeader({ className, ...props }: React.HTMLAttributes<HTMLTableSectionElement>) {
  return <thead className={cn('text-left text-xs font-semibold uppercase tracking-wide text-muted-foreground', className)} {...props} />
}

export function TableBody({ className, ...props }: React.HTMLAttributes<HTMLTableSectionElement>) {
  return <tbody className={cn('align-top text-sm text-foreground', className)} {...props} />
}

export function TableRow({ className, ...props }: React.HTMLAttributes<HTMLTableRowElement>) {
  return (
    <tr
      className={cn(
        'overflow-hidden rounded-2xl border border-[color-mix(in_srgb,var(--border)_75%,transparent)] bg-[color-mix(in_srgb,var(--card)_82%,transparent)] backdrop-blur-sm transition hover:border-[color-mix(in_srgb,var(--border)_90%,transparent)]',
        className,
      )}
      {...props}
    />
  )
}

export function TableHead({ className, ...props }: React.ThHTMLAttributes<HTMLTableCellElement>) {
  return <th className={cn('px-4 pb-2 text-left text-[0.72rem] font-semibold uppercase tracking-[0.08em] text-muted-foreground', className)} {...props} />
}

export function TableCell({ className, ...props }: React.TdHTMLAttributes<HTMLTableCellElement>) {
  return <td className={cn('px-4 py-3 align-middle text-sm text-foreground', className)} {...props} />
}
