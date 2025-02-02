import { useOnce } from ".";

export default function useEventListener<
  E extends HTMLElement,
  K extends keyof HTMLElementEventMap
>(element: E, event: K, listener: (event: HTMLElementEventMap[K]) => void): void {
  useOnce(() => {
    element.addEventListener(event, listener)

    return () => element.removeEventListener(event, listener)
  })
}
