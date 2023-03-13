package pubsub

import (
	"context"
	"sync"
)

type Client interface {
	Subscribe(ctx context.Context, channels ...string) error
	Unsubscribe(ctx context.Context, channels ...string) error
}

type Subscriber[T any] interface {
	Dispatch(payload T)
}

type SubscriberMap[T any] map[Subscriber[T]]bool
type SubscriptionMap[T any] map[string]SubscriberMap[T]

type SubscriptionManager[T any] struct {
	mu            sync.Mutex
	pubsub        Client
	subscriptions SubscriptionMap[T]
}

func (m *SubscriptionManager[T]) add(topic string, sub Subscriber[T]) {
	if m.subscriptions[topic] == nil {
		m.subscriptions[topic] = SubscriberMap[T]{sub: true}
	} else {
		m.subscriptions[topic][sub] = true
	}
}

func (m *SubscriptionManager[T]) Add(ctx context.Context, topic string, sub Subscriber[T]) error {
	return m.AddMulti(ctx, []string{topic}, sub)
}

func (m *SubscriptionManager[T]) AddMulti(ctx context.Context, topics []string, sub Subscriber[T]) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var newTopics []string

	for _, topic := range topics {
		m.add(topic, sub)

		if len(m.subscriptions[topic]) == 1 {
			newTopics = append(newTopics, topic)
		}
	}

	return m.pubsub.Subscribe(ctx, newTopics...)
}

func (m *SubscriptionManager[T]) Remove(ctx context.Context, topic string, sub Subscriber[T]) error {
	return m.RemoveMulti(ctx, []string{topic}, sub)
}

func (m *SubscriptionManager[T]) RemoveMulti(ctx context.Context, topics []string, sub Subscriber[T]) error {
	var emptyTopics []string

	for _, topic := range topics {
		if m.subscriptions[topic][sub] {
			delete(m.subscriptions[topic], sub)

			if len(m.subscriptions[topic]) == 0 {
				emptyTopics = append(emptyTopics, topic)

				delete(m.subscriptions, topic)
			}
		}
	}

	return m.pubsub.Unsubscribe(ctx, emptyTopics...)
}

func (m *SubscriptionManager[T]) Dispatch(topic string, payload T) {
	subs := m.subscriptions[topic]

	for sub := range subs {
		sub.Dispatch(payload)
	}
}

func NewSubscriptionManager[T any](client Client) *SubscriptionManager[T] {
	return &SubscriptionManager[T]{
		pubsub:        client,
		subscriptions: SubscriptionMap[T]{},
	}
}
