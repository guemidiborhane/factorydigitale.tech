import { Key, useT } from "i18n";
import toast, { Toaster as T, ToastIcon } from "react-hot-toast";
import { resolveValue, type Toast, type ToasterProps } from "react-hot-toast/headless";

import styles from './Toaster.module.scss'
import clsx from "clsx";

export default function Toaster(props: ToasterProps) {
  const { t: tt } = useT()
  const options = props ?? {}
  options.position = 'bottom-left'
  options.toastOptions = {
    duration: 500
  }

  return (
    <T {...options}>
      {(t: Toast) => (
        <div
          class={clsx(styles.Toast, t.visible ? styles.ToastAnimationIn : styles.ToastAnimationOut)}
        >
          <div class={styles.ToastBody}>
            <ToastIcon toast={t} />
            <p class={styles.ToastBodyMessage}>
              {
                tt(`paths.${resolveValue(t.message, t)}.${t.type}` as Key)
                || tt(`paths.${resolveValue(t.message, t)}` as Key)
                || resolveValue(t.message, t)}
            </p>
          </div>
          <div class={styles.ToastClose}>
            <button
              onClick={() => toast.dismiss(t.id)}
              class={styles.ToastCloseButton}
            >
              {tt('misc.close')}
            </button>
          </div>
        </div>
      )}
    </T>
  );
};
