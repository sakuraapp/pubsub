package pubsub

type MessageTarget interface {
	Build() string
}

type Dispatcher interface {
	Dispatch(topic string, message interface{}) error
	DispatchTo(target MessageTarget, message interface{}) error
}
