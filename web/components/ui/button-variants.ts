import { cn } from '@/lib/utils'

export type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost' | 'destructive'
export type ButtonSize = 'xs' | 'sm' | 'md' | 'lg'

const variantStyles: Record<ButtonVariant, string> = {
  primary:
    'border border-transparent bg-primary text-primary-foreground shadow-[0_18px_48px_-30px_rgba(79,70,229,0.75)] hover:-translate-y-px hover:bg-[color-mix(in_srgb,var(--primary)_88%,transparent)] hover:shadow-[0_20px_52px_-28px_rgba(79,70,229,0.65)] active:translate-y-0 active:shadow-[0_16px_36px_-28px_rgba(79,70,229,0.7)]',
  secondary:
    'border border-[color-mix(in_srgb,var(--border)_80%,transparent)] bg-[color-mix(in_srgb,var(--secondary)_90%,transparent)] text-secondary-foreground shadow-[0_1px_0_rgba(255,255,255,0.08)] hover:bg-[color-mix(in_srgb,var(--secondary)_78%,transparent)]',
  outline:
    'border border-border bg-[color-mix(in_srgb,var(--background)_94%,transparent)] text-foreground shadow-[0_1px_0_rgba(255,255,255,0.05)] hover:bg-[color-mix(in_srgb,var(--foreground)_6%,transparent)]',
  ghost:
    'border border-transparent bg-transparent text-muted-foreground hover:bg-[color-mix(in_srgb,var(--foreground)_6%,transparent)] hover:text-foreground',
  destructive:
    'border border-transparent bg-destructive text-destructive-foreground shadow-[0_16px_40px_-28px_rgba(220,38,38,0.55)] hover:bg-[color-mix(in_srgb,var(--destructive)_82%,transparent)]',
}

const sizeStyles: Record<ButtonSize, string> = {
  xs: 'h-8 px-3 text-xs',
  sm: 'h-9 px-3.5 text-sm',
  md: 'h-10 px-4 text-sm',
  lg: 'h-12 px-5 text-base',
}

export function buttonVariants({
  variant = 'primary',
  size = 'md',
}: { variant?: ButtonVariant; size?: ButtonSize } = {}) {
  return cn(
    'inline-flex items-center justify-center gap-2 rounded-full font-medium transition-all duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:opacity-60',
    variantStyles[variant],
    sizeStyles[size],
  )
}
