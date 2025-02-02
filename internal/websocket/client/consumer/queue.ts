import { Signal, effect, signal } from "@preact/signals";
import { type Message as M } from "./message";
import { SendFunction, isReady } from ".";
import { generateRandomID } from "./subscription";

type Message = {
  payload: M
  id: string
  sent: boolean
}
type MessageQueue = Signal<Message[]>
export type QueueFunction = (message: M) => void

export interface QueueManager {
  add: QueueFunction
}

const q: MessageQueue = signal([])

export class Queue implements QueueManager {
  private send: SendFunction
  private receipts: string[] = []
  private cleaningInterval = 50

  constructor(send: SendFunction) {
    this.send = send

    effect(() => {
      this.handleQueue(q, isReady)
    })

    setInterval(this.clean.bind(this), this.cleaningInterval)
  }

  add(payload: M) {
    const currentQueue = q.value
    const message = {
      payload,
      id: generateRandomID(),
      sent: false
    }
    q.value = [...currentQueue, message]
  }

  private clean() {
    const cleaningNeeded = q.value.some(m => this.receipts.includes(m.id))
    if (!cleaningNeeded) return

    q.value = q.value.filter(m => !this.receipts.includes(m.id))
  }

  private handleQueue(q: MessageQueue, ready: Signal<boolean>) {
    if (q.value.length == 0) return
    if (!ready) {
      setTimeout(() => this.handleQueue(q, ready), 500)
      return
    }

    q.value.forEach(m => {
      if (this.receipts.includes(m.id)) return

      if (this.send(m.payload)) {
        this.receipts.push(m.id)
      }
    })
  }
}
