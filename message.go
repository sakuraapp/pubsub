package pubsub

import (
	"github.com/sakuraapp/shared/pkg/model"
	"github.com/sakuraapp/shared/pkg/resource"
	"github.com/sakuraapp/shared/pkg/resource/permission"
	"github.com/vmihailenco/msgpack/v5"
)

type MessageType int

const (
	NormalMessage MessageType = iota
	BroadcastMessage
	ServerMessage
)

type MessageTarget struct {
	UserIds           []model.UserId       `msgpack:"u,omitempty"`
	RoomId            model.RoomId         `msgpack:"r,omitempty"`
	Permissions       permission.Permission `msgpack:"p,omitempty"`
	IgnoredSessionIds map[string]bool       `msgpack:"i,omitempty"`
}

type Message struct {
	Type   MessageType      `msgpack:"t,omitempty"`
	Target *MessageTarget   `msgpack:"tr,omitempty"`
	Data   *resource.Packet `msgpack:"d,omitempty"`
	Origin string           `msgpack:"o,omitempty"` // source/origin node of the message
}

type rawMessage Message

func (m Message) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal((rawMessage)(m))
}

func (m *Message) UnmarshalBinary(b []byte) error {
	return msgpack.Unmarshal(b, (*rawMessage)(m))
}