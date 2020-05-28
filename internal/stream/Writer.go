package stream

import (
	"fmt"
	"gochat-system/internal/protocol"
	"io"
	"net/url"
)

// CommandWriter is a writer used for client/server
type CommandWriter struct {
	writer io.Writer
}

// NewCommandWriter constructs a new CommandWriter
func NewCommandWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{
		writer: writer,
	}
}

func (w *CommandWriter) writeString(msg string) error {
	_, err := w.writer.Write([]byte(msg))
	return err
}

func (w *CommandWriter) Write(command interface{}) error {
	var err error

	switch v := command.(type) {
	case protocol.SendCommand:
		err = w.writeString(fmt.Sprintf("SEND %v\n", v.Message))
	case protocol.MessageCommand:
		err = w.writeString(fmt.Sprintf("MESSAGE %v %v\n", v.Name, url.QueryEscape(v.Message)))
	case protocol.NameCommand:
		err = w.writeString(fmt.Sprintf("NAME %v\n", v.Name))
	case protocol.LeaveMessage:
		err = w.writeString(fmt.Sprintf("LEAVE %v\n", v.Name))
	case protocol.OnlineMessage:
		err = w.writeString(fmt.Sprintf("ONLINE %v\n", v.RemoteAddr))
	case protocol.SetNickNameMessage:
		err = w.writeString(fmt.Sprintf("NICKNAME %v %v\n", v.RemoteAddr, v.Name))
	default:
		err = &protocol.UnknownCommand{}
	}

	return err
}
