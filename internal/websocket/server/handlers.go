package websocket

import (
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/internal/logger"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/users/auth"
)

func (h *Hub) UpgradeHandler(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (h *Hub) WebSocketHandler() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		client := &Client{
			Connection: c,
			LastSeenAt: time.Now(),
			RequestId:  c.Locals("request-id").(string),
			UserID:     c.Locals(auth.USER_ID).(int),
		}

		logger.Debug("Client connected", client.LogAttrs())

		defer func() {
			c.Close()
			h.cleanUpClient(&Client{Connection: c})
		}()

		client.Send(&Message{Type: "connected"})
		h.handleStream(c, client)
		h.cleanUpClient(client)
	})
}

func (h *Hub) TestHandler(c *fiber.Ctx) error {
	ch := c.Params("channel")
	channel := h.NewChannel(&ChannelAttributes{Name: ch})

	return channel.Broadcast(c.Params("message"))
}
