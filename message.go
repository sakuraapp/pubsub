package pubsub

import (
	"github.com/vmihailenco/msgpack/v5"
)

type MessageFilterKind int

type MessageTarget interface {
	Build() string
}

type FilterMap = map[MessageFilterKind]interface{}

type Message struct {
	Filters FilterMap 		`msgpack:"f,omitempty"`
	Payload interface{} 	`msgpack:"p,omitempty"`
}

func (m *Message) WithFilter(kind MessageFilterKind, value interface{}) *Message {
	m.Filters[kind] = value

	return m
}

type rawMessage Message

func (m Message) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal((rawMessage)(m))
}

func (m *Message) UnmarshalBinary(b []byte) error {
	return msgpack.Unmarshal(b, (*rawMessage)(m))
}

func NewMessage(payload interface{}, opts ...FilterMap) Message {
	if len(opts) == 0 {
		opts = append(opts, FilterMap{})
	}

	return Message{
		Payload: payload,
		Filters: opts[0],
	}
}