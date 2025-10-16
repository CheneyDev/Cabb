'use client'

import * as React from 'react'

import { cn } from '@/lib/utils'

export interface SwitchProps extends React.InputHTMLAttributes<HTMLInputElement> {}

export const Switch = React.forwardRef<HTMLInputElement, SwitchProps>(function Switch({ className, ...props }, ref) {
  return (
    <label className={cn('switch-root', className)}>
      <input ref={ref} type="checkbox" className="switch-input" {...props} />
      <span className="switch-track" />
      <span className="switch-thumb" />
    </label>
  )
})
