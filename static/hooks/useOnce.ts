import { type EffectCallback, useEffect } from "preact/hooks";

export default function useOnce(fn: EffectCallback) {
  useEffect(fn, [])
}
