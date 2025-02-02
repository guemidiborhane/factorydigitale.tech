package websocket

import (
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	Broker   *Broker
	Channels map[string]*Channel
	sync.Mutex
}

func NewHub(redisClient *redis.Client) *Hub {
	h := &Hub{
		Broker:   NewBroker(redisClient),
		Channels: make(map[string]*Channel),
	}

	go h.StartInactiveClientCheck()

	return h
}

const TIMEOUT time.Duration = 30 * time.Second

func (h *Hub) StartInactiveClientCheck() {
	ticker := time.NewTicker(TIMEOUT)
	for range ticker.C {
		h.checkInactiveClients()
	}
}

func (h *Hub) checkInactiveClients() {
	h.Lock()
	defer h.Unlock()
	now := time.Now()
	for _, channel := range h.Channels {
		for _, client := range channel.Subscriptions {
			if now.Sub(client.LastSeenAt) > TIMEOUT {
				client.Disconnect(h)
			}
		}
	}
}
