package pubsub

const (
	NodeTopic = "gateway.%v"
	SessionTopic = "session.%v"
	UserTopic = "user.%v"
	RoomTopic = "room.%v"
)

const (
	NormalMessage MessageType = iota
	ServerMessage
)

const (
	MessageTargetNode MessageTargetKind = iota
	MessageTargetSession
	MessageTargetUser
	MessageTargetRoom
)

const (
	MessageFilterIgnoredSession MessageFilterKind = iota
	MessageFilterPermissions
	MessageFilterRoom
)

var targets = map[MessageTargetKind]string {
	MessageTargetNode: NodeTopic,
	MessageTargetSession: SessionTopic,
	MessageTargetUser: UserTopic,
	MessageTargetRoom: RoomTopic,
}