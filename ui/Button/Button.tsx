import type { HTMLAttributes } from "preact/compat";
import styles, { type ClassNames } from './Button.module.scss'
import clsx from "clsx";

interface Props extends HTMLAttributes<HTMLButtonElement> {
  type?: Exclude<ClassNames, 'Button'>;
  full?: boolean
}

export default function Button({ children, type, full = true, class: className, ...props }: Props) {
  return (
    <button class={clsx(
      styles.Button,
      type && styles[type],
      full && styles.Full,
      className
    )} {...props}>
      {children}
    </button>
  )
}
