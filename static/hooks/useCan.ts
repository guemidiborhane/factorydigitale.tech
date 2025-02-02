import { signedIn, useContextStore } from "~/contexts/user"

export function useCan(action?: string): boolean {
  const context = useContextStore()

  if (!action) return true
  if (!signedIn.value) return false;

  const hasPermission = !!context.currentUser?.permissions?.includes(action)

  return hasPermission
}
