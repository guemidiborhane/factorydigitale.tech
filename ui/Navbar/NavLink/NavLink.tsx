import { FC, type ReactNode } from "preact/compat";
import { Link, useLocation } from "react-router-dom";

import styles from './NavLink.module.scss'
import Can from "ui/Can";
import clsx from "clsx";

const normalize = (str?: string) => str?.replace(/\//g, '')

const isActive = (to: string) => {
  const { pathname } = useLocation()

  return normalize(to) == normalize(pathname)
}

const NavLink: FC<{ to: string, children: ReactNode, can?: string }> = ({ to, children, can }) => {
  return (
    <Can action={can}>
      <Link to={to} class={clsx(styles.NavLink, isActive(to) && styles.NavLinkActive)}>
        {children}
      </Link>
    </Can>
  )
}

export default NavLink
