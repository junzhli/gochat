package stream

import (
	"bufio"
	"gochat-system/internal/protocol"
	"io"
	"log"
	"net/url"
)

// CommandReader is a reader used for client/server
type CommandReader struct {
	reader *bufio.Reader
}

// NewCommandReader constructs a new CommandReader
func NewCommandReader(reader io.Reader) *CommandReader {
	return &CommandReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *CommandReader) Read() (interface{}, error) {
	commandName, err := r.reader.ReadString(' ')
	if err != nil {
		return nil, err
	}

	switch commandName {
	case "MESSAGE ":
		user, err := r.reader.ReadString(' ')
		if err != nil {
			return nil, err
		}

		escapedMessage, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		message, err := url.QueryUnescape(escapedMessage[:len(escapedMessage)-1])
		if err != nil {
			return nil, err
		}
		return protocol.MessageCommand{
			Name:    user[:len(user)-1],
			Message: message,
		}, nil
	case "SEND ":
		message, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return protocol.SendCommand{
			Message: message[:len(message)-1],
		}, nil
	case "NAME ":
		user, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return protocol.NameCommand{
			Name: user[:len(user)-1],
		}, nil
	case "LEAVE ":
		user, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return protocol.LeaveMessage{
			Name: user[:len(user)-1],
		}, nil
	case "ONLINE ":
		remoteAddr, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return protocol.OnlineMessage{
			RemoteAddr: remoteAddr[:len(remoteAddr)-1],
		}, nil
	case "NICKNAME ":
		remoteAddr, err := r.reader.ReadString(' ')
		if err != nil {
			return nil, err
		}

		name, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		return protocol.SetNickNameMessage{
			RemoteAddr: remoteAddr[:len(remoteAddr)-1],
			Name:       name[:len(name)-1],
		}, nil
	default:
		log.Printf("Unknown command: %v", commandName)
		return nil, &protocol.UnknownCommand{}
	}
}
