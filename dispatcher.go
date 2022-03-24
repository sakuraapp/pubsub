package pubsub

import (
	"github.com/sakuraapp/shared/pkg/model"
)

type Dispatcher interface {
	Dispatch(topic string, message Message) error
	DispatchTo(target *MessageTarget, message Message) error
	DispatchToRoom(roomId model.RoomId, message Message) error
}