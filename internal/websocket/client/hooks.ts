import { useOnce } from '~/hooks'
import consumer from './consumer'
import type { Identifier } from './consumer/message'

type WebSocketArgs<T> = Identifier & {
  receiver: (data: T) => void
}

export function useWebSocket<T>({ receiver, ...identity }: WebSocketArgs<T>) {
  useOnce(() => {
    const subId = consumer.subscribe(identity, receiver)

    return () => consumer.unsubscribe(subId)
  })
}
