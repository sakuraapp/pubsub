package pubsub

import (
	"context"
	"sync"
)

type Client interface {
	Subscribe(ctx context.Context, channels ...string) error
	Unsubscribe(ctx context.Context, channels ...string) error
}

type SubscriptionMap = map[string]int

type SubscriptionManager struct {
	mu sync.Mutex
	pubsub Client
	subscriptions SubscriptionMap
}

func (m *SubscriptionManager) Add(ctx context.Context, topic string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.subscriptions[topic] += 1

	if m.subscriptions[topic] == 1 {
		return m.pubsub.Subscribe(ctx, topic)
	} else {
		return nil
	}
}

func (m *SubscriptionManager) AddMulti(ctx context.Context, topics []string) error {
	var newTopics []string

	for _, topic := range topics {
		m.subscriptions[topic] += 1

		if m.subscriptions[topic] == 1 {
			newTopics = append(newTopics, topic)
		}
	}

	return m.pubsub.Subscribe(ctx, newTopics...)
}

func (m *SubscriptionManager) Remove(ctx context.Context, topic string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.subscriptions[topic] > 0 {
		m.subscriptions[topic] -= 1

		if m.subscriptions[topic] == 0 {
			return m.pubsub.Unsubscribe(ctx, topic)
		}
	}

	return nil
}

func (m *SubscriptionManager) RemoveMulti(ctx context.Context, topics []string) error {
	var emptyTopics []string

	for _, topic := range topics {
		if m.subscriptions[topic] > 0 {
			m.subscriptions[topic] -= 1

			if m.subscriptions[topic] == 0 {
				emptyTopics = append(emptyTopics, topic)
			}
		}
	}

	return m.pubsub.Unsubscribe(ctx, emptyTopics...)
}

func NewSubscriptionManager(client Client) *SubscriptionManager {
	return &SubscriptionManager{
		pubsub: client,
		subscriptions: SubscriptionMap{},
	}
}