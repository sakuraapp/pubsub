package pubsub

type Dispatcher interface {
	Dispatch(topic string, message Message) error
	DispatchTo(target *MessageTarget, message Message) error
}