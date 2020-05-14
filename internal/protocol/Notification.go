package protocol

// LeaveMessage is for notifying clients that someone is going to leave
type LeaveMessage struct {
	Name string
}

// OnlineMessage is for notifying clients that new client is connected without nicknamed
type OnlineMessage struct {
	RemoteAddr string
}

// SetNickNameMessage is used for notifying clients that the connected client gets nicknamed
type SetNickNameMessage struct {
	RemoteAddr string
	Name       string
}
