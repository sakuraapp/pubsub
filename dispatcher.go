package pubsub

import (
	"github.com/sakuraapp/shared/pkg/model"
)

type LocalDispatcher interface {
	DispatchLocal(msg *Message) error
	DispatchRoomLocal(roomId model.RoomId, msg *Message) error
	HandleServerMessage(msg *Message)
}

type Dispatcher interface {
	LocalDispatcher
	Dispatch(msg *Message) error
	DispatchRoom(roomId model.RoomId, msg *Message) error
}