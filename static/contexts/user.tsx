import { createContext } from "preact";
import { type FC, type ReactNode } from "preact/compat";
import { useContext } from "preact/hooks";
import { $user, type UserStore } from "~/stores/user";
import { useStore } from "@nanostores/preact";
import { checkUser } from "~/helpers";
import { useOnce } from "~/hooks";
import { signal } from "@preact/signals";

export type UserContext = { currentUser?: UserStore }

export const signedIn = signal<boolean>(false)
export const Context = createContext<Partial<UserContext>>({})

export const UserProvider: FC<{ children?: ReactNode }> = ({ children }) => {
  useOnce(() => {
    const controller = new AbortController()
    checkUser(controller.signal).then(b => signedIn.value = b)

    return () => controller.abort()
  })

  const currentUser = useStore($user)

  return (
    <Context.Provider value={{
      currentUser
    }}>
      {children}
    </Context.Provider >
  )
}
export const useContextStore = () => useContext(Context)
