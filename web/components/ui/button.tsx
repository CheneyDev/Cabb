import * as React from 'react'

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'outline'
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  { className = '', variant = 'primary', ...props },
  ref,
) {
  const base = 'btn'
  const style = variant === 'outline' ? 'btn-outline' : 'btn-primary'
  return <button ref={ref} className={`${base} ${style} ${className}`} {...props} />
})

