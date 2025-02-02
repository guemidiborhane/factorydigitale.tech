import { type ReactNode } from 'preact/compat'
import styles from './Modal.module.scss'
import clsx from 'clsx'
import { type Signal } from '@preact/signals'
import { useEventListener } from '~/hooks'

type Props = {
  show: Signal<boolean>
  children: ReactNode
  closeOnOverlayClick?: boolean
}

export default function Modal({ show, children, closeOnOverlayClick = true }: Props) {
  if (!show.value) return null

  const close = () => show.value = false
  useEventListener(document.body, 'keyup', (e) => {
    if (e.key === 'Escape') close()
  })

  return (
    <div class={styles.Modal}>
      <div onClick={() => closeOnOverlayClick && close()} class={styles.ModalOverLay}></div>
      <button onClick={close} class={styles.ModalCloseButton}>
        <svg class={styles.ModalCloseButtonIcon} viewBox="0 0 24 24" width="48" height="48">
          <path fill="currentColor" d="M12 10.586l4.95-4.95 1.414 1.414-4.95 4.95 4.95 4.95-1.414 1.414-4.95-4.95-4.95 4.95-1.414-1.414 4.95-4.95-4.95-4.95L7.05 5.636z"></path>
        </svg>
      </button>
      <div class={clsx(styles.ModalBox, 'box z-10 w-full md:w-3/4 lg:w-max overflow-y-auto')}>
        {children}
      </div>
    </div >
  )
}
