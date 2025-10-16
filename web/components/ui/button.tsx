'use client'

import * as React from 'react'

import { cn } from '@/lib/utils'

import { buttonVariants, type ButtonSize, type ButtonVariant } from './button-variants'

export interface ButtonProps extends React.ComponentPropsWithoutRef<'button'> {
  variant?: ButtonVariant
  size?: ButtonSize
  render?: React.ReactElement
}

function assignRef<T>(ref: React.Ref<T> | undefined, value: T | null) {
  if (!ref) return
  if (typeof ref === 'function') {
    ref(value)
  } else {
    ;(ref as React.MutableRefObject<T | null>).current = value
  }
}

function composeRefs<T>(...refs: (React.Ref<T> | undefined)[]) {
  return (value: T | null) => {
    for (const ref of refs) {
      assignRef(ref, value)
    }
  }
}

export const Button = React.forwardRef<HTMLElement, ButtonProps>(function Button(
  { className, variant = 'primary', size = 'md', render, children, type, ...restProps },
  forwardedRef,
) {
  const finalClassName = cn(buttonVariants({ variant, size }), className)

  if (render) {
    const element = render as React.ReactElement & { ref?: React.Ref<HTMLElement> }
    const { className: renderClassName, children: renderChildren, ref: renderRef, ...renderRest } = element.props

    return React.cloneElement(
      element,
      {
        ...renderRest,
        ...restProps,
        className: cn(finalClassName, renderClassName),
        ref: composeRefs(renderRef as React.Ref<HTMLElement>, forwardedRef),
      },
      children ?? renderChildren,
    )
  }

  return (
    <button
      type={type ?? 'button'}
      className={finalClassName}
      ref={composeRefs(forwardedRef)}
      {...restProps}
    >
      {children}
    </button>
  )
})

Button.displayName = 'Button'

export { buttonVariants }
export type { ButtonSize, ButtonVariant }
