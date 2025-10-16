import * as React from 'react'

import { cn } from '@/lib/utils'

type BadgeVariant = 'default' | 'outline' | 'success' | 'destructive' | 'info' | 'muted'

export interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: BadgeVariant
}

const variantStyles: Record<BadgeVariant, string> = {
  default: 'border border-[color-mix(in_srgb,var(--secondary)_65%,transparent)] bg-[color-mix(in_srgb,var(--secondary)_75%,transparent)] text-secondary-foreground',
  outline: 'border border-border text-foreground',
  success: 'border border-[color-mix(in_srgb,var(--success)_65%,transparent)] bg-[color-mix(in_srgb,var(--success)_22%,transparent)] text-success-foreground',
  destructive: 'border border-[color-mix(in_srgb,var(--destructive)_65%,transparent)] bg-[color-mix(in_srgb,var(--destructive)_22%,transparent)] text-destructive-foreground',
  info: 'border border-[color-mix(in_srgb,var(--info)_65%,transparent)] bg-[color-mix(in_srgb,var(--info)_22%,transparent)] text-info-foreground',
  muted: 'border border-[color-mix(in_srgb,var(--border)_80%,transparent)] bg-[color-mix(in_srgb,var(--foreground)_6%,transparent)] text-muted-foreground',
}

export function Badge({ className, variant = 'default', ...props }: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-medium uppercase tracking-wide',
        variantStyles[variant],
        className,
      )}
      {...props}
    />
  )
}
