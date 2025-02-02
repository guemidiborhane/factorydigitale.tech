import { useActionData } from "react-router-dom"
import styles from './Input.module.scss'
import { type Key, useT } from "i18n"
import type { APIError } from "~/helpers/http"
import { HTMLAttributes } from "preact/compat"
import clsx from "clsx"

type Props = HTMLAttributes<HTMLInputElement> & {
  name: string
  type?: string
  label?: Key
}

const capitalize = (str: string) => str.charAt(0).toUpperCase() + str.slice(1)

export default function Input({ name, type = "text", label, ...props }: Props) {
  const { t } = useT()
  const data = useActionData() as APIError
  const errors = data && data.message instanceof Object ? data.message : {}
  const hasErrors = Object.keys(errors).length > 0

  return (
    <div class={clsx(hasErrors && styles.HasError)}>
      <label for={name} class={clsx(styles.Label, errors[name] == 'required' && styles.LabelRequired)}>{capitalize(label && t(label) || name)}</label>
      <input type={type} name={name} id={name} class={styles.Input} {...props} dir="auto" />
      {
        hasErrors && (
          <span class={styles.HelpBlock}>
            <ul>
              <li>{t(`validations.${errors[name]}` as Key, { field: label && t(label) })}</li>
            </ul>
          </span>
        )
      }
    </div >
  )
}
