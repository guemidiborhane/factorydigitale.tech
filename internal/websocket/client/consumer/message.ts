export type Identifier = {
  channel: string
}

export enum Types {
  CONNECTED = 'connected',
  PONG = 'pong'
}

export type Message<T = any> = { timestamp?: number } & (
  {
    type: "connected" | "ping" | "pong"
  }
  | {
    type: "subscribe" | "subscription_confirmation" | "unsubscribe" | "unsubscription_confirmation"
    identifier: Identifier
  }
  | {
    type: "message"
    identifier: Identifier
    data: T
  }
)
