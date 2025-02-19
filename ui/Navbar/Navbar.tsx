import { useT } from "i18n"
import Auth from "./Auth/Auth"
import NavLink from "./NavLink/NavLink"

import styles from './Navbar.module.scss'
import { routes } from "~/router"

export default function Navbar() {
  const { t } = useT()

  return (
    <nav class={styles.Navbar}>
      <div>
        <NavLink to="/">
          {t('resources.home')}
        </NavLink>
        <NavLink to={routes.movies} can="movies:index">
          {t('resources.movies')}
        </NavLink>
        <NavLink to={routes.movies_favourites} can="movies:favourite">
          {t('resources.favourites')}
        </NavLink>
      </div>

      <div>
        <NavLink to={routes.permissions} can="permissions:index">
          {t("resources.permissions", { other: true })}
        </NavLink>
        <Auth />
      </div>
    </nav>
  )
}
