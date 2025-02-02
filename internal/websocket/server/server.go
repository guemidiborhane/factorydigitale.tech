package websocket

import (
	"errors"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/guemidiborhane/factorydigitale.tech/internal/logger"
	"github.com/guemidiborhane/factorydigitale.tech/internal/router"
	"github.com/guemidiborhane/factorydigitale.tech/pkg/users"
)

func (h *Hub) handleStream(c *websocket.Conn, client *Client) {
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			if errors.Is(err, websocket.ErrCloseSent) {
				h.cleanUpClient(client)
			} else {
				log.Println("read:", err)
			}
			break
		}

		h.handleMessage(c, msg, client)
	}
}

func (h *Hub) cleanUpClient(client *Client) {
	for _, channel := range h.Channels {
		channel.Unsubscribe(client)
	}
}

func (h *Hub) RegisterRoutes(a *fiber.App) {
	r := a.Group("/ws",
		h.UpgradeHandler,
		router.EncryptCookies(),
		router.HelmetMiddleware,
		router.CorsMiddleware,
		logger.Middleware,
		router.CompressMiddleware,
		router.EtagMiddleware,
		router.CsrfMiddleware(),
		router.RequestIDMiddleware,
		users.CheckAuthenticated,
	)

	r.Get("/", h.WebSocketHandler())
}
