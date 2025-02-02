package websocket

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Broker struct {
	*redis.Client
}

func NewBroker(redisClient *redis.Client) *Broker {
	return &Broker{
		Client: redisClient,
	}
}

func (b *Broker) Publish(channel string, message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return b.Client.Publish(ctx, channel, message).Err()
}

func (b *Broker) Subscribe(channel string) *redis.PubSub {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return b.Client.Subscribe(ctx, channel)
}

func (b *Broker) StartListening(channel *Channel, pubsub *redis.PubSub, stopCh <-chan struct{}) {
	go func() {
		defer pubsub.Close()

		ch := pubsub.Channel()
		for {
			select {
			case msg := <-ch:
				for _, client := range channel.Subscriptions {
					var message Message
					if err := ParseMessage([]byte(msg.Payload), &message); err != nil {
						log.Printf("Error unmarshalling message: %v", err)
						continue
					}
					client.Send(&message)
				}
			case <-stopCh:
				return
			}
		}
	}()
}
