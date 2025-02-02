import { FileEdit, List, Plus, ShieldQuestion, Trash } from "lucide-preact"
import styles from "./Card.module.scss"
import { useOnce } from "~/hooks"
import { ReactNode } from "preact/compat"
import clsx from "clsx"
import { useSignal } from "@preact/signals"


const PermissionsIcons = {
  index: List,
  create: Plus,
  update: FileEdit,
  destroy: Trash
}

type ValidIcons = keyof typeof PermissionsIcons

const iconize = (p: string, hasPermission: boolean): ReactNode => {
  const Element = PermissionsIcons[p as ValidIcons] || ShieldQuestion

  return (
    <abbr title={p}>
      <Element class={clsx(hasPermission ? styles.CardIconValid : styles.CardIconInvalid)} />
    </abbr>
  )
}

export default function Card({ role, permissions }: { role: string, permissions: { [key: string]: string[] } }) {
  const defaultPermissions = Object.keys(PermissionsIcons);
  const columnsCount = useSignal<number>(defaultPermissions.length);

  const getPermissions = (p: string[]): string[] => {
    const filtered = p.filter((p) => !defaultPermissions.includes(p));
    return [...defaultPermissions, ...filtered];
  };
  useOnce(() => {
    Object.values(permissions).forEach((p) => {
      const pp = getPermissions(p)

      columnsCount.value = Math.max(columnsCount.value, pp.length)
    })
  })

  return (
    <div class={styles.Card} >
      <h5 class={styles.CardTitle}>
        {role}
      </h5>

      <table dir="ltr">
        <tbody>
          {Object.entries(permissions).map(([resource, permissions]) => (
            <tr>
              <td>{resource}</td>
              {getPermissions(permissions).map(p => {
                return <td>{iconize(p, permissions.includes(p))}</td>
              })}
              {columnsCount.value - getPermissions(permissions).length > 0 && new Array(columnsCount.value - getPermissions(permissions).length).fill(<td />)}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

