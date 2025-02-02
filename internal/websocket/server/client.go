package websocket

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"github.com/guemidiborhane/factorydigitale.tech/internal/logger"
	"github.com/gofiber/contrib/websocket"
	"github.com/redis/go-redis/v9"
)

type Connection = websocket.Conn

type Client struct {
	Connection   *Connection   `json:"-"`
	Subscription *redis.PubSub `json:"-"`
	LastSeenAt   time.Time     `json:"last_seen_at"`
	Channel      string        `json:"channel"`
	RequestId    string        `json:"request_id"`
	UserID       int           `json:"user_id"`
}

func (c *Client) Send(message *Message) error {
	return c.WriteMessage(websocket.TextMessage, message)
}

func (c *Client) WriteMessage(t int, message *Message) error {
	message.Timestamp = time.Now().Unix()

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return c.Connection.WriteMessage(t, jsonData)
}

func (c *Client) LogAttrs() logger.Attrs {
	return logger.Attrs{
		"user-id":      c.UserID,
		"channel":      c.Channel,
		"last-seen-at": c.LastSeenAt,
		"request-id":   c.RequestId,
	}
}

func (c *Client) Disconnect(h *Hub) error {
	h.Lock()
	defer h.Unlock()

	logger.Warn("Closing connection", c.LogAttrs())
	c.Connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseTryAgainLater, "You've been missing"))
	time.Sleep(1 * time.Second)
	defer c.Connection.Close()

	for _, channel := range h.Channels {
		for i, client := range channel.Subscriptions {
			if client.Connection == c.Connection {
				channel.Subscriptions = append(channel.Subscriptions[:i], channel.Subscriptions[i+1:]...)
			}
		}
	}

	return nil
}

func (c *Client) Confirm(channel *Channel, t string) error {
	identifier, _ := json.Marshal(channel.Attributes)
	message := &Message{
		Type:       fmt.Sprintf("%s_confirmation", t),
		Identifier: identifier,
	}

	return c.Send(message)
}
