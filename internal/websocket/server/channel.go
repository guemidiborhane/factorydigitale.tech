package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/goccy/go-json"

	"github.com/guemidiborhane/factorydigitale.tech/internal/logger"
)

type ChannelAttributes struct {
	Name     string                        `json:"channel"`
	Receiver func(m *Message, ch *Channel) `json:"-" binding:"required"`
}

type Channel struct {
	Attributes    *ChannelAttributes
	Hub           *Hub      `json:"-"`
	Subscriptions []*Client `json:"-"`
	sync.Mutex    `json:"-"`
	stopCh        chan struct{} `json:"-"`
}

func (h *Hub) NewChannel(attributes *ChannelAttributes) *Channel {
	h.Lock()
	defer h.Unlock()

	name := attributes.Name

	if attributes.Receiver == nil {
		attributes.Receiver = func(m *Message, ch *Channel) {}
	}

	channel, exists := h.Channels[name]
	if !exists {
		channel = &Channel{
			Attributes:    attributes,
			Hub:           h,
			Subscriptions: []*Client{},
		}
		h.Channels[name] = channel
	}

	return channel
}

func (c *Channel) Name() string {
	return c.Attributes.Name
}

func (c *Channel) Identifier() []byte {
	jsonData, err := json.Marshal(c.Attributes)
	if err != nil {
		log.Printf("Error occurred during marshal of identifier: %s", err.Error())
		return nil
	}
	return jsonData
}

func (c *Channel) Broadcast(msg any) error {
	messageJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	message := &Message{
		Type:       MessageType,
		Identifier: c.Identifier(),
		Data:       messageJson,
		Timestamp:  time.Now().Unix(),
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	logger.Info("Streaming message", logger.Attrs{
		"channel": c.Name(),
		"payload": string(jsonData),
	})

	return c.Hub.Broker.Publish(c.Name(), jsonData)
}

func (c *Channel) Subscribe(client *Client) {
	c.Lock()
	defer c.Unlock()

	for _, existingClient := range c.Subscriptions {
		if existingClient.Connection == client.Connection {
			return
		}
	}

	client.Channel = c.Name()
	c.Subscriptions = append(c.Subscriptions, client)
	client.Confirm(c, "subscription")
	logger.Info("Client subscribed", client.LogAttrs())

	if len(c.Subscriptions) == 1 {
		c.stopCh = make(chan struct{})
		pubsub := c.Hub.Broker.Subscribe(c.Name())
		c.Hub.Broker.StartListening(c, pubsub, c.stopCh)
	}
}

func (c *Channel) Unsubscribe(client *Client) {
	c.Lock()
	defer c.Unlock()

	for i, sub := range c.Subscriptions {
		if sub.Connection == client.Connection {
			c.Subscriptions = append(c.Subscriptions[:i], c.Subscriptions[i+1:]...)
			if len(c.Subscriptions) == 0 && c.stopCh != nil {
				close(c.stopCh)
				c.stopCh = nil
			}
			sub.Confirm(c, "unsubscription")
			logger.Info("Client unsubscribed", sub.LogAttrs())
			break
		}
	}
}
