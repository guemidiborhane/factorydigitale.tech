import type { Identifier, Message } from "./message"
import { QueueFunction } from "./queue"

export type Receiver<T = any> = (message: T) => void

type Listener = {
  id: string
  receiver: Receiver
}

export interface SubscriptionInterface {
  confirmed: boolean
  readonly markedForDeath: boolean
  readonly retry: boolean
  retryDelay: number
  listeners: Listener[]
  identifier: Identifier
  channel: string
  unsubscriptionMessage: Message
  addListener: (listener: Receiver) => string
  sendMessage: (data: any) => void
  handle: <T>(data: T) => void
  confirm: () => void
  close: (listenerID?: string) => void
  findListener: (id: string) => Listener | undefined
  subscribe: () => void
}
export const generateRandomID = (): string => {
  return Math.random().toString(36).substring(2, 15)
}

export class Subscription implements SubscriptionInterface {
  listeners: Listener[] = []
  identifier: Identifier
  confirmed: boolean = false
  markedForDeath: boolean = false
  retry: boolean = true
  private retryInterval: NodeJS.Timeout | undefined
  retryDelay: number = 5
  private closeDelay: NodeJS.Timeout | undefined
  private send: QueueFunction

  constructor(identifier: Identifier, receiver: Receiver, send: QueueFunction) {
    this.send = send
    this.identifier = identifier
    this.addListener(receiver)
    this.subscribe()
  }

  subscribe() {
    if (this.confirmed) return

    if (this.retry) {
      this.retryInterval = setTimeout(() => this.handleTimeout(), this.retryDelay * 1_000)
    }
    this.send(this.subscriptionMessage)
  }

  private handleTimeout() {
    console.log(`Subscription ${this.channel} timed out`)
    if (!this.confirmed && this.retry) {
      this.retryInterval = setInterval(() => {
        console.log(`Retrying subscription to ${this.channel}`);
        this.send(this.subscriptionMessage);
      }, this.retryDelay * 1000);
    }
  }

  addListener(receiver: Receiver): string {
    const id = generateRandomID()
    this.listeners.push({ id, receiver })

    if (this.markedForDeath) {
      clearInterval(this.closeDelay)
    }

    return id
  }

  sendMessage(data: any) {
    this.send({
      type: "message",
      identifier: this.identifier,
      data,
    })
  }

  handle<T>(data: T) {
    this.listeners.forEach(listener => listener.receiver(data))
  }

  close(listenerID?: string): void {
    this.listeners = this.listeners.filter(listener => listener.id != listenerID)
    if (this.listeners.length == 0 || listenerID == null) {
      this.closeDelay = setInterval(() => {
        if (this.listeners.length > 0) return
        this.markedForDeath = true
        this.send(this.unsubscriptionMessage)
        clearInterval(this.closeDelay)
      }, 1_000)
    }
  }

  confirm() {
    console.log(`Subscribed to ${this.channel}`)
    this.confirmed = true
    clearTimeout(this.retryInterval)
  }

  findListener(id: string): Listener | undefined {
    return this.listeners.find(listener => listener.id === id)
  }


  private get subscriptionMessage(): Message {
    return {
      type: "subscribe",
      identifier: this.identifier,
    }
  }

  get unsubscriptionMessage(): Message {
    return {
      type: "unsubscribe",
      identifier: this.identifier,
    }
  }

  get channel(): string {
    return this.identifier.channel
  }

}

