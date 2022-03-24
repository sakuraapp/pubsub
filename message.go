package pubsub

import (
	"fmt"
	"github.com/sakuraapp/shared/pkg/model"
	"github.com/sakuraapp/shared/pkg/resource"
	"github.com/sakuraapp/shared/pkg/resource/permission"
	"github.com/vmihailenco/msgpack/v5"
)

type MessageType int

type MessageTargetKind int
type MessageFilterKind int


type MessageTarget struct {
	Kind MessageTargetKind
	Value interface{}
}

func (m *MessageTarget) Build() string {
	return fmt.Sprintf(targets[m.Kind], m.Value)
}

type FilterMap map[MessageFilterKind]interface{}

func (m FilterMap) WithIgnoredSession(sessionId string) FilterMap {
	m[MessageFilterIgnoredSession] = sessionId

	return m
}

func (m FilterMap) WithPermissions(perm permission.Permission) FilterMap {
	m[MessageFilterPermissions] = perm

	return m
}

func (m FilterMap) WithRoom(roomId model.RoomId) FilterMap {
	m[MessageFilterRoom] = roomId

	return m
}

func NewFilterMap() *FilterMap {
	return &FilterMap{}
}

type MessageOptions struct {
	Filters FilterMap `msgpack:"f,omitempty"`
}

type Message struct {
	Type    MessageType     `msgpack:"t,omitempty"`
	Options *MessageOptions `msgpack:"o,omitempty,inline"`
	Payload resource.Packet `msgpack:"p,omitempty"`
}

type rawMessage Message

func (m Message) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal((rawMessage)(m))
}

func (m *Message) UnmarshalBinary(b []byte) error {
	return msgpack.Unmarshal(b, (*rawMessage)(m))
}

func NewMessage(payload resource.Packet, opts ...*MessageOptions) Message {
	if len(opts) == 0 {
		opts = append(opts, &MessageOptions{})
	}

	return Message{
		Payload: payload,
		Options: opts[0],
	}
}