import { signal } from "@preact/signals"
import { type Identifier, type Message, Types } from "./message"
import { type Receiver } from "./subscription"
import { SubscriptionManager, type SubscriptionManagerInterface } from "./manager"

const { host, protocol } = window.location
const schema = protocol == 'https:' ? 'wss' : 'ws';

const URL = `${schema}://${host}/ws`
export const isReady = signal(false)

export type SendFunction = (message: Message) => boolean

interface ConsumerInterface {
  send: SendFunction
}

class Consumer implements ConsumerInterface {
  private websocket: WebSocket | undefined
  private SubscriptionsManager: SubscriptionManagerInterface
  private pingEach: number = 30 // seconds
  private pingTimer: NodeJS.Timeout | undefined
  private shouldReconnect: boolean = true
  private reconnectTimer: NodeJS.Timeout | undefined
  private reconnectInterval: number = 5
  private reconnecting: boolean = false

  constructor() {
    this.websocket = this.connect()
    this.SubscriptionsManager = new SubscriptionManager(this.send.bind(this))

    if (this.shouldReconnect) {
      this.reconnectTimer = setTimeout(() => this.constructor(), 5_000)
    }
  }

  private connect(): WebSocket {
    isReady.value = false
    const websocket = new WebSocket(URL)
    console.debug("[websocket] connecting...")
    websocket.onmessage = this.receiver.bind(this)
    websocket.onopen = this.onOpen.bind(this)
    websocket.onclose = this.onClose.bind(this);
    websocket.onerror = this.onError.bind(this);
    return websocket
  }

  subscribe<T>(identifier: Identifier, receiver: Receiver<T>): string {
    return this.SubscriptionsManager.addSubscription(identifier, receiver)
  }

  unsubscribe(subscriptionID: string) {
    return this.SubscriptionsManager.removeSubscription(subscriptionID)
  }

  send(message: Message): boolean {
    if (
      isReady.value
      && this.websocket
      && this.websocket?.readyState == WebSocket.OPEN
    ) {
      message.timestamp = this.timestamp
      this.websocket.send(JSON.stringify(message))

      return true
    }

    return false
  }

  private ping(): void {
    this.send({ type: 'ping' })
  }

  private get timestamp(): number {
    return Math.floor(Date.now() / 1_000)
  }

  private onOpen() {
    clearInterval(this.reconnectTimer)
    this.pingTimer = setInterval(this.ping.bind(this), this.pingEach * 1_000)
  }

  private reconnect() {
    if (this.shouldReconnect) {
      this.reconnecting = true
      this.reconnectTimer = setTimeout(() => {
        this.websocket = this.connect()
      }, this.reconnectInterval * 1_000);
    }
  }

  private onClose() {
    clearInterval(this.pingTimer);
    isReady.value = false;
    this.reconnect()
    this.SubscriptionsManager.closeShop()
    this.websocket?.close()
    this.websocket = undefined
  }

  private onError() {
    console.error("[websocket] error");
    this.websocket?.close();
  }

  private receiver(event: MessageEvent): void {
    const message: Message = JSON.parse(event.data)

    switch (message.type) {
      case Types.CONNECTED:
        isReady.value = true
        console.debug(`[websocket] connected`)
        if (this.reconnecting) {
          this.SubscriptionsManager.reconnect()
          this.reconnecting = false
        }
        break;
      case Types.PONG:
        break;
      default:
        this.SubscriptionsManager.spread(message)
    }
  }
}

export default new Consumer()
