import type { ReactNode } from "preact/compat";
import { useCan } from "~/hooks/useCan";

export default function Can({ action, children }: { action?: string, children: ReactNode }) {
  if (!useCan(action)) return null

  return <>{children}</>
}
