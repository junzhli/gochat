package protocol

// SendCommand is for sending new message
type SendCommand struct {
	Message string
}

// NameCommand is for naming itself
type NameCommand struct {
	Name string
}

// MessageCommand is for modifying message
type MessageCommand struct {
	Name    string
	Message string
}

// UnknownCommand is a type of error for unknown message
type UnknownCommand struct {
}

func (w *UnknownCommand) Error() string {
	return "Unknown message"
}
