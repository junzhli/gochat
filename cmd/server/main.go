package main

import (
	"flag"
	"fmt"
	"gochat-system/internal/protocol"
	"gochat-system/internal/stream"
	"io"
	"log"
	"net"
	"sync"
)

type ChatServer interface {
	Listen(address string) error
	Broadcast(command interface{})
	Start()
	Close()
}

type client struct {
	conn   net.Conn
	name   string
	writer *stream.CommandWriter
}

type TcpChatServer struct {
	listener net.Listener
	clients  []*client
	mutex    *sync.Mutex
}

func NewServer() ChatServer {
	return &TcpChatServer{
		listener: nil,
		clients:  make([]*client, 0),
		mutex:    &sync.Mutex{},
	}
}

func (s *TcpChatServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.listener = l
	log.Printf("Listening on %v", address)
	return nil
}

func (s *TcpChatServer) Close() {
	s.listener.Close()
}

func (s *TcpChatServer) Start() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("connection error: %v\n", conn)
		}

		client := s.accept(conn)
		go s.serve(client)
	}
}

func (s *TcpChatServer) accept(conn net.Conn) *client {
	log.Printf("Accepting connection from %v | total clients %v", conn.RemoteAddr().String(), len(s.clients)+1)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &client{
		conn:   conn,
		writer: stream.NewCommandWriter(conn),
	}

	go s.Broadcast(protocol.OnlineMessage{
		RemoteAddr: conn.RemoteAddr().String(),
	})

	s.clients = append(s.clients, client)

	return client
}

func (s *TcpChatServer) remove(client *client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// remove the connections from clients array
	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}
	log.Printf("Closing connection from %v", client.conn.RemoteAddr().String())
	client.conn.Close()
}

// serve method says hello to new client
func (s *TcpChatServer) serve(client *client) {
	cmdReader := stream.NewCommandReader(client.conn)
	defer s.remove(client)

	for {
		cmd, err := cmdReader.Read()

		if err != nil && err != io.EOF {
			log.Printf("Read error occurred: %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.SendCommand:
				go s.Broadcast(protocol.MessageCommand{
					Message: v.Message,
					Name:    client.name,
				})
			case protocol.NameCommand:
				client.name = v.Name
				go s.Broadcast(protocol.SetNickNameMessage{
					RemoteAddr: client.conn.RemoteAddr().String(),
					Name:       client.name,
				})
			}
		}

		if err == io.EOF {
			go s.Broadcast(protocol.LeaveMessage{
				Name: client.name,
			})
			break
		}
	}
}

func (s *TcpChatServer) Broadcast(command interface{}) {
	for _, client := range s.clients {
		err := client.writer.Write(command)
		if err != nil {
			log.Printf("Unable to broadcast to client: %v: %v", client.conn.RemoteAddr(), client.name)
		}
	}
}

func main() {
	ip := flag.String("ip", "", "Specify a network interface address to start listening")
	port := flag.String("port", "3333", "Specify a listening port")
	flag.Parse()

	fmt.Print("Chat console (Server)\n")
	fmt.Print("===============\n")

	s := NewServer()
	err := s.Listen(fmt.Sprintf("%v:%v", *ip, *port))
	if err != nil {
		log.Fatalf("Unable to listen on %v:%v: %v", *ip, *port, err)
	}

	s.Start()
}
