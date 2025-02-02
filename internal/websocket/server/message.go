package websocket

import (
	"log"
	"time"

	"github.com/goccy/go-json"

	"github.com/gofiber/contrib/websocket"
)

type Identifier map[string]string

type Message struct {
	Type       string          `json:"type"`
	Identifier json.RawMessage `json:"identifier,omitempty"`
	Data       json.RawMessage `json:"data,omitempty"`
	Timestamp  int64           `json:"timestamp"`
}

var (
	SubscriptionType   = "subscribe"
	UnsubscriptionType = "unsubscribe"
	MessageType        = "message"
	PingType           = "ping"
	PongType           = "pong"
	CloseType          = "bye-bye"
)

func ParseMessage(msg []byte, message *Message) error {
	if err := json.Unmarshal(msg, &message); err != nil {
		return err
	}

	return nil
}

func (m *Message) Json() ([]byte, error) {
	json, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return json, nil
}

func (m *Message) getIdentifier() ChannelAttributes {
	var id ChannelAttributes
	json.Unmarshal(m.Identifier, &id)
	return id
}

func (h *Hub) handleMessage(c *Connection, msg []byte, client *Client) {
	client.LastSeenAt = time.Now()

	var message Message
	if err := ParseMessage(msg, &message); err != nil {
		log.Printf("Error during handling of message: %s", err.Error())
	}

	ch := message.getIdentifier().Name

	switch message.Type {
	case SubscriptionType:
		if ch != "" {
			h.NewChannel(&ChannelAttributes{Name: ch}).Subscribe(client)
		}
	case UnsubscriptionType:
		if ch != "" {
			h.NewChannel(&ChannelAttributes{Name: ch}).Unsubscribe(client)
		}
	case MessageType:
		if ch != "" {
			channel, exists := h.Channels[ch]
			if exists && channel.Attributes.Receiver != nil {
				channel.Attributes.Receiver(&message, channel)
			}
		}
	case PingType:
		if err := client.WriteMessage(websocket.PongMessage, &Message{Type: PongType}); err != nil {
			log.Printf("Error during sending pong: %s", err.Error())
		}
	default:
		log.Printf("Unknown message type: %s", message.Type)
	}
}
