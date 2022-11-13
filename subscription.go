package pubsub

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
)

type Subscriber interface {
	Receive(message Message)
}

type SubscriberMap = map[Subscriber]bool
type SubscriptionMap = map[string]SubscriberMap

type SubscriptionManager struct {
	mu sync.Mutex
	pubsub *redis.PubSub
	subscriptions SubscriptionMap
}

func (m *SubscriptionManager) has(topic string, sub Subscriber) bool {
	if m.subscriptions[topic] != nil {
		return m.subscriptions[topic][sub]
	} else {
		return false
	}
}

func (m *SubscriptionManager) add(topic string, sub Subscriber) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.subscriptions[topic] == nil {
		m.subscriptions[topic] = SubscriberMap{sub: true}
	} else {
		m.subscriptions[topic][sub] = true
	}
}

func (m *SubscriptionManager) Subscribe(ctx context.Context, topic string, sub Subscriber) error {
	if !m.has(topic, sub) {
		m.add(topic, sub)

		return m.pubsub.Subscribe(ctx, topic)
	} else {
		return nil
	}
}

func (m *SubscriptionManager) SubscribeMulti(ctx context.Context, topics []string, sub Subscriber) error {
	var newTopics []string

	for _, topic := range topics {
		if !m.has(topic, sub) {
			m.add(topic, sub)

			newTopics = append(newTopics, topic)
		}
	}

	return m.pubsub.Subscribe(ctx, newTopics...)
}