import * as React from 'react'

import { cn } from '@/lib/utils'

type AlertVariant = 'default' | 'success' | 'info' | 'destructive'

export interface AlertProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: AlertVariant
}

const variantStyles: Record<AlertVariant, string> = {
  default: 'border border-[color-mix(in_srgb,var(--border)_75%,transparent)] bg-[color-mix(in_srgb,var(--card)_82%,transparent)] text-foreground',
  success: 'border border-[color-mix(in_srgb,var(--success)_65%,transparent)] bg-[color-mix(in_srgb,var(--success)_22%,transparent)] text-success-foreground',
  info: 'border border-[color-mix(in_srgb,var(--info)_65%,transparent)] bg-[color-mix(in_srgb,var(--info)_22%,transparent)] text-info-foreground',
  destructive: 'border border-[color-mix(in_srgb,var(--destructive)_65%,transparent)] bg-[color-mix(in_srgb,var(--destructive)_22%,transparent)] text-destructive-foreground',
}

export function Alert({ className, variant = 'default', ...props }: AlertProps) {
  return (
    <div
      role="alert"
      className={cn(
        'flex w-full items-start gap-3 rounded-2xl px-4 py-3 text-sm shadow-[0_1px_0_rgba(255,255,255,0.04)] backdrop-blur-sm',
        variantStyles[variant],
        className,
      )}
      {...props}
    />
  )
}
