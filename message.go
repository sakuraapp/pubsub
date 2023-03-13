package pubsub

import "github.com/vmihailenco/msgpack/v5"

type MessageFilterKind int

type FilterMap = map[MessageFilterKind]interface{}

type Message[T any] struct {
	Filters FilterMap `msgpack:"f,omitempty"`
	Payload T         `msgpack:"p,omitempty"`
}

func (m *Message[T]) WithFilter(kind MessageFilterKind, value T) *Message[T] {
	m.Filters[kind] = value

	return m
}

type rawMessage[T any] Message[T]

func (m *Message[T]) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal((*rawMessage[T])(m))
}

func (m *Message[T]) UnmarshalBinary(b []byte) error {
	return msgpack.Unmarshal(b, (*rawMessage[T])(m))
}

func NewMessage[T any](payload T, opts ...FilterMap) *Message[T] {
	if len(opts) == 0 {
		opts = append(opts, FilterMap{})
	}

	return &Message[T]{
		Payload: payload,
		Filters: opts[0],
	}
}
