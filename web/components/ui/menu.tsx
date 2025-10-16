'use client'

import * as React from 'react'
import { Menu as MenuPrimitive } from '@base-ui-components/react/menu'

import { cn } from '@/lib/utils'

type ExtractState<T> = T extends (state: infer S) => any ? S : never

function mergeClassNames<State>(
  base: string,
  className?: string | ((state: State) => string | undefined),
) {
  if (typeof className === 'function') {
    return (state: State) => cn(base, className(state))
  }
  return cn(base, className)
}

const baseItemStyles =
  'group flex w-full select-none items-center justify-between gap-2 rounded-2xl px-3 py-2 text-sm text-muted-foreground transition-colors duration-150 data-[highlighted]:bg-[color-mix(in_srgb,var(--primary)_16%,transparent)] data-[highlighted]:text-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-45 data-[active=true]:bg-[color-mix(in_srgb,var(--primary)_18%,transparent)] data-[active=true]:text-primary'

const baseTriggerStyles =
  'inline-flex min-w-[8rem] items-center justify-center gap-2 rounded-full border border-[color-mix(in_srgb,var(--border)_70%,transparent)] bg-[color-mix(in_srgb,var(--background)_88%,transparent)] px-3 py-1.5 text-sm font-medium text-foreground shadow-[0_18px_36px_-30px_rgba(79,70,229,0.8)] transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background data-[open=true]:bg-[color-mix(in_srgb,var(--background)_82%,transparent)]'

const basePopupStyles =
  'z-50 min-w-[15rem] max-w-[92vw] rounded-3xl border border-[color-mix(in_srgb,var(--border)_78%,transparent)] bg-[color-mix(in_srgb,var(--background)_96%,transparent)] p-2 shadow-[0_32px_90px_-45px_rgba(79,70,229,0.55)] backdrop-blur-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=open]:fade-in data-[state=closed]:fade-out data-[side=bottom]:slide-in-from-top-2'

type MenuTriggerElement = React.ElementRef<typeof MenuPrimitive.Trigger>
type MenuTriggerProps = React.ComponentPropsWithoutRef<typeof MenuPrimitive.Trigger>

type MenuTriggerState = ExtractState<NonNullable<MenuTriggerProps['className']>>

export const MenuTrigger = React.forwardRef<MenuTriggerElement, MenuTriggerProps>(function MenuTrigger(
  { className, ...props },
  ref,
) {
  return (
    <MenuPrimitive.Trigger
      ref={ref}
      className={mergeClassNames<MenuTriggerState>(baseTriggerStyles, className)}
      {...props}
    />
  )
})

MenuTrigger.displayName = 'MenuTrigger'

type MenuPopupElement = React.ElementRef<typeof MenuPrimitive.Popup>
type MenuPopupProps = React.ComponentPropsWithoutRef<typeof MenuPrimitive.Popup>

type MenuPopupState = ExtractState<NonNullable<MenuPopupProps['className']>>

export const MenuPopup = React.forwardRef<MenuPopupElement, MenuPopupProps>(function MenuPopup(
  { className, ...props },
  ref,
) {
  return (
    <MenuPrimitive.Popup
      ref={ref}
      className={mergeClassNames<MenuPopupState>(basePopupStyles, className)}
      {...props}
    />
  )
})

MenuPopup.displayName = 'MenuPopup'

type MenuItemElement = React.ElementRef<typeof MenuPrimitive.Item>
type MenuItemProps = React.ComponentPropsWithoutRef<typeof MenuPrimitive.Item>

type MenuItemState = ExtractState<NonNullable<MenuItemProps['className']>>

export const MenuItem = React.forwardRef<MenuItemElement, MenuItemProps>(function MenuItem(
  { className, ...props },
  ref,
) {
  return (
    <MenuPrimitive.Item
      ref={ref}
      className={mergeClassNames<MenuItemState>(baseItemStyles, className)}
      {...props}
    />
  )
})

MenuItem.displayName = 'MenuItem'

type MenuCheckboxItemElement = React.ElementRef<typeof MenuPrimitive.CheckboxItem>
type MenuCheckboxItemProps = React.ComponentPropsWithoutRef<typeof MenuPrimitive.CheckboxItem>

type MenuCheckboxItemState = ExtractState<NonNullable<MenuCheckboxItemProps['className']>>

export const MenuCheckboxItem = React.forwardRef<MenuCheckboxItemElement, MenuCheckboxItemProps>(
  function MenuCheckboxItem({ className, ...props }, ref) {
    return (
      <MenuPrimitive.CheckboxItem
        ref={ref}
        className={mergeClassNames<MenuCheckboxItemState>(baseItemStyles, className)}
        {...props}
      />
    )
  },
)

MenuCheckboxItem.displayName = 'MenuCheckboxItem'

type MenuRadioItemElement = React.ElementRef<typeof MenuPrimitive.RadioItem>
type MenuRadioItemProps = React.ComponentPropsWithoutRef<typeof MenuPrimitive.RadioItem>

type MenuRadioItemState = ExtractState<NonNullable<MenuRadioItemProps['className']>>

export const MenuRadioItem = React.forwardRef<MenuRadioItemElement, MenuRadioItemProps>(function MenuRadioItem(
  { className, ...props },
  ref,
) {
  return (
    <MenuPrimitive.RadioItem
      ref={ref}
      className={mergeClassNames<MenuRadioItemState>(baseItemStyles, className)}
      {...props}
    />
  )
})

MenuRadioItem.displayName = 'MenuRadioItem'

type MenuGroupElement = React.ElementRef<typeof MenuPrimitive.Group>
type MenuGroupProps = React.ComponentPropsWithoutRef<typeof MenuPrimitive.Group>

type MenuGroupState = ExtractState<NonNullable<MenuGroupProps['className']>>

export const MenuGroup = React.forwardRef<MenuGroupElement, MenuGroupProps>(function MenuGroup(
  { className, ...props },
  ref,
) {
  return (
    <MenuPrimitive.Group
      ref={ref}
      className={mergeClassNames<MenuGroupState>('flex flex-col gap-1', className)}
      {...props}
    />
  )
})

MenuGroup.displayName = 'MenuGroup'

type MenuGroupLabelElement = React.ElementRef<typeof MenuPrimitive.GroupLabel>
type MenuGroupLabelProps = React.ComponentPropsWithoutRef<typeof MenuPrimitive.GroupLabel>

type MenuGroupLabelState = ExtractState<NonNullable<MenuGroupLabelProps['className']>>

export const MenuGroupLabel = React.forwardRef<MenuGroupLabelElement, MenuGroupLabelProps>(function MenuGroupLabel(
  { className, ...props },
  ref,
) {
  return (
    <MenuPrimitive.GroupLabel
      ref={ref}
      className={mergeClassNames<MenuGroupLabelState>(
        'px-3 pt-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground',
        className,
      )}
      {...props}
    />
  )
})

MenuGroupLabel.displayName = 'MenuGroupLabel'

type MenuSeparatorElement = React.ElementRef<typeof MenuPrimitive.Separator>
type MenuSeparatorProps = React.ComponentPropsWithoutRef<typeof MenuPrimitive.Separator>

type MenuSeparatorState = ExtractState<NonNullable<MenuSeparatorProps['className']>>

export const MenuSeparator = React.forwardRef<MenuSeparatorElement, MenuSeparatorProps>(function MenuSeparator(
  { className, ...props },
  ref,
) {
  return (
    <MenuPrimitive.Separator
      ref={ref}
      className={mergeClassNames<MenuSeparatorState>(
        'my-2 h-px bg-[color-mix(in_srgb,var(--border)_82%,transparent)]',
        className,
      )}
      {...props}
    />
  )
})

MenuSeparator.displayName = 'MenuSeparator'

type MenuSubTriggerElement = React.ElementRef<typeof MenuPrimitive.SubmenuTrigger>
type MenuSubTriggerProps = React.ComponentPropsWithoutRef<typeof MenuPrimitive.SubmenuTrigger>

type MenuSubTriggerState = ExtractState<NonNullable<MenuSubTriggerProps['className']>>

export const MenuSubTrigger = React.forwardRef<MenuSubTriggerElement, MenuSubTriggerProps>(function MenuSubTrigger(
  { className, ...props },
  ref,
) {
  return (
    <MenuPrimitive.SubmenuTrigger
      ref={ref}
      className={mergeClassNames<MenuSubTriggerState>(baseItemStyles, className)}
      {...props}
    />
  )
})

MenuSubTrigger.displayName = 'MenuSubTrigger'

type MenuSubPopupElement = React.ElementRef<typeof MenuPrimitive.Popup>
type MenuSubPopupProps = React.ComponentPropsWithoutRef<typeof MenuPrimitive.Popup>

type MenuSubPopupState = ExtractState<NonNullable<MenuSubPopupProps['className']>>

export const MenuSubPopup = React.forwardRef<MenuSubPopupElement, MenuSubPopupProps>(function MenuSubPopup(
  { className, ...props },
  ref,
) {
  return (
    <MenuPrimitive.Popup
      ref={ref}
      className={mergeClassNames<MenuSubPopupState>(basePopupStyles, className)}
      {...props}
    />
  )
})

MenuSubPopup.displayName = 'MenuSubPopup'

export const Menu = MenuPrimitive.Root
export const MenuRadioGroup = MenuPrimitive.RadioGroup
export const MenuSub = MenuPrimitive.SubmenuRoot
