import { useT } from "i18n";
import { Form, useSubmit } from "react-router-dom";
import { type SubmitTarget } from "react-router-dom/dist/dom";
import { useContextStore, signedIn } from "~/contexts/user";
import { routes } from "~/router";
import Modal from "ui/Modal";

import { NavLinkStyles } from 'ui/Navbar'
import Styles from './Auth.module.scss'
import clsx from "clsx";
import { useSignal } from "@preact/signals";
import { useWebSocket } from "ws/client";

const UserDisplay = () => {
  if (!signedIn.value) return null

  useWebSocket<string>({
    channel: "test",
    receiver(data) {
      console.log(data, "test")
    }
  })

  const ctx = useContextStore()

  if (process.env.NODE_ENV !== "development") {
    return (
      <span class={clsx(NavLinkStyles.NavLink, NavLinkStyles.NavLinkActive)}>
        {ctx.currentUser?.username}
      </span>
    )
  }

  const show = useSignal<boolean>(false)

  return (
    <>
      <span class={clsx(NavLinkStyles.NavLink, NavLinkStyles.NavLinkActive, 'cursor-pointer')} onClick={() => show.value = true}>
        {ctx.currentUser?.username}
      </span>
      <Modal show={show}>
        <pre dir="ltr">
          {JSON.stringify(ctx.currentUser, undefined, 2)}
        </pre>
      </Modal>
    </>
  )
}

export default function Auth() {
  if (!signedIn.value) return null

  const { t } = useT()
  const submit = useSubmit();
  const handleSignout = (e: Event) => {
    e.preventDefault();
    e.stopPropagation();
    if (confirm(t("misc.confirmation_msg"))) {
      submit(e.currentTarget as SubmitTarget);
    }
  };

  return (
    <>
      <Form action={routes.users_sign_out} method="DELETE" onSubmit={handleSignout} class="h-full">
        <button class={clsx(Styles.Auth, NavLinkStyles.NavLink)}>{t('misc.sign_out')}</button>
      </Form>
      <UserDisplay />
    </>
  );
}
