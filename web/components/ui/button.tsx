import * as React from 'react'

import { cn } from '@/lib/utils'

type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost'
type ButtonSize = 'sm' | 'md' | 'lg'

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant
  size?: ButtonSize
}

const variantStyles: Record<ButtonVariant, string> = {
  primary:
    'border border-transparent bg-primary text-primary-foreground shadow-[0_1px_0_rgba(255,255,255,0.08),0_18px_45px_-24px_rgba(79,70,229,0.75)] hover:-translate-y-[1px] hover:bg-[color-mix(in_srgb,var(--primary)_88%,transparent)] hover:shadow-[0_18px_50px_-22px_rgba(79,70,229,0.6)] active:translate-y-0 active:shadow-[0_12px_30px_-24px_rgba(79,70,229,0.7)]',
  secondary:
    'border border-[color-mix(in_srgb,var(--border)_80%,transparent)] bg-[color-mix(in_srgb,var(--secondary)_85%,transparent)] text-secondary-foreground shadow-[0_1px_0_rgba(255,255,255,0.08)] hover:bg-[color-mix(in_srgb,var(--secondary)_65%,transparent)]',
  outline:
    'border border-border bg-transparent text-foreground shadow-[0_1px_0_rgba(255,255,255,0.05)] hover:bg-[color-mix(in_srgb,var(--foreground)_6%,transparent)]',
  ghost: 'border border-transparent bg-transparent text-muted-foreground hover:bg-[color-mix(in_srgb,var(--foreground)_6%,transparent)] hover:text-foreground',
}

const sizeStyles: Record<ButtonSize, string> = {
  sm: 'h-9 px-3 text-xs',
  md: 'h-10 px-4 text-sm',
  lg: 'h-12 px-5 text-base',
}

export function buttonVariants({ variant = 'primary', size = 'md' }: { variant?: ButtonVariant; size?: ButtonSize } = {}) {
  return cn(
    'inline-flex items-center justify-center gap-2 rounded-full font-medium transition-all duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:opacity-60',
    variantStyles[variant],
    sizeStyles[size],
  )
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  { className, variant = 'primary', size = 'md', ...props },
  ref,
) {
  return <button ref={ref} className={cn(buttonVariants({ variant, size }), className)} {...props} />
})
