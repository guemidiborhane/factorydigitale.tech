import { Message, type Identifier } from "./message"
import { type Receiver, Subscription, type SubscriptionInterface } from "./subscription"
import { Queue, QueueManager } from "./queue"
import { SendFunction } from "."


export interface SubscriptionManagerInterface {
  addSubscription: (identifier: Identifier, receiver: Receiver) => string
  removeSubscription: (id: string) => void
  confirm: (msg: Message) => void
  spread: (message: Message) => void
  closeShop: () => void
  reconnect: () => void
}

export class SubscriptionManager implements SubscriptionManagerInterface {
  private subscriptions: SubscriptionInterface[] = []
  private queue: QueueManager

  constructor(send: SendFunction) {
    this.queue = new Queue(send)
  }

  addSubscription(identifier: Identifier, receiver: Receiver): string {
    let id: string
    const existingSubscription = this.subscriptions.find(sub => sub.channel == identifier.channel)

    if (existingSubscription) {
      id = existingSubscription.addListener(receiver)
    } else {
      const subscription = new Subscription(identifier, receiver, this.queue.add)
      id = subscription.listeners[0]?.id
      this.subscriptions.push(subscription)
    }

    return id
  }

  removeSubscription(id: string) {
    this.subscriptions.find(sub => !!sub.findListener(id))?.close(id)
  }

  closeShop() {
    this.subscriptions.forEach(sub => {
      this.queue.add(sub.unsubscriptionMessage)
    })
  }

  reconnect() {
    this.subscriptions.forEach(sub => {
      sub.confirmed = false
      sub.subscribe()
    })
  }

  confirm(msg: Message) {
    if (!('identifier' in msg)) return
    const { type, identifier: { channel: ch } } = msg

    switch (type.replace("_confirmation", "")) {
      case "subscription":
        this.subscriptions = this.subscriptions.map(sub => {
          if (sub.channel == ch) sub.confirm()
          return sub
        })
        break;
      case "unsubscription":
        console.log(`Unsubscribed from ${ch}`)
        this.subscriptions = this.subscriptions.filter(sub => {
          return !sub.markedForDeath && sub.channel != ch
        })
        break;
    }
  }

  spread(message: Message): void {
    if (message.type.endsWith("_confirmation")) {
      this.confirm(message)
      return
    }

    if (message.type != 'message') return
    this.subscriptions.find(sub => sub.channel == message.identifier.channel)?.handle(message.data)
  }
}
