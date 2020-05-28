package main

import (
	"bufio"
	"flag"
	"fmt"
	"gochat-system/internal/protocol"
	"gochat-system/internal/stream"
	"io"
	"log"
	"net"
	"os"
)

// ChatClient provides client program with clear explanation of interfaces
type ChatClient interface {
	// Dial is called to reach server
	Dial(address string) error
	// Send command to online users
	Send(command interface{}) error
	// SendMessage is a wrapper function that send text messages to chat room with Send method
	SendMessage(message string) error
	// GetName returns client's nickname
	GetName() string
	// SetName is a wrapper function that register its nickname to chat room with Send method
	SetName(name string) error
	// Start method listens on tcp connection
	Start()
	// Close method terminates its tcp connection
	Close()
	// Incoming returns channel that receives commands/notifications from server
	Incoming() chan interface{}
}

type tcpChatClient struct {
	conn      net.Conn
	cmdReader *stream.CommandReader
	cmdWriter *stream.CommandWriter
	name      string
	incoming  chan interface{}
}

// NewClient creates an instance of ChatClient
func NewClient() ChatClient {
	return &tcpChatClient{
		incoming: make(chan interface{}),
	}
}

func (c *tcpChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.conn = conn
	c.cmdReader = stream.NewCommandReader(conn)
	c.cmdWriter = stream.NewCommandWriter(conn)
	return nil
}

func (c *tcpChatClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

func (c *tcpChatClient) SendMessage(message string) error {
	return c.Send(protocol.SendCommand{
		Message: message,
	})
}

func (c *tcpChatClient) GetName() string {
	return c.name
}

func (c *tcpChatClient) SetName(name string) error {
	c.name = name
	return c.Send(protocol.NameCommand{
		Name: name,
	})
}

func (c *tcpChatClient) Start() {
	for {
		cmd, err := c.cmdReader.Read()
		if err == io.EOF {
			log.Fatal("Server disconnected")
		}
		if err != nil {
			log.Printf("Read error occurred: %v", err)
		}

		if cmd != nil {
			c.incoming <- cmd
		}
	}
}

func (c *tcpChatClient) Close() {
	err := c.conn.Close()
	if err != nil {
		fmt.Printf("Connection terminated with error occurred: %v", err)
	}
}

func (c *tcpChatClient) Incoming() chan interface{} {
	return c.incoming
}

func listenOnIncomingMessages(c ChatClient) {
	for any := range c.Incoming() {
		switch m := any.(type) {
		case protocol.MessageCommand:
			fmt.Printf("\n[User] %v: %v\n", m.Name, m.Message)
		case protocol.OnlineMessage:
			fmt.Printf("[System] New client is online: %v\n", m.RemoteAddr)
		case protocol.LeaveMessage:
			fmt.Printf("[System] %v left\n", m.Name)
		case protocol.SetNickNameMessage:
			fmt.Printf("[System] Client %v set nickname: %v\n", m.RemoteAddr, m.Name)
		default:
			fmt.Printf("[System] Unknown message: %v\n", m)
		}

	}
}

func main() {
	ip := flag.String("ip", "localhost", "Specify a remote address to connect")
	port := flag.String("port", "3333", "Specify a remote port")
	flag.Parse()

	fmt.Print("Chat console (Client)\n")
	fmt.Print("===============\n")

	c := NewClient()
	fmt.Printf("Connecting to server %v:%v...", *ip, *port)
	if err := c.Dial(fmt.Sprintf("%v:%v", *ip, *port)); err != nil {
		fmt.Printf("fail\n")
		os.Exit(1)
	}
	fmt.Printf("success\n")

	go c.Start()
	defer c.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Set up your name: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Read name error occurred: %v", err)
	}
	err = c.SetName(text[:len(text)-1])
	if err != nil {
		fmt.Printf("Unable to set name. Reason: %v", err)
	}

	fmt.Printf("Wellcome %v\n", text)
	go listenOnIncomingMessages(c)
	for {
		//fmt.Printf("--> ")
		textBytes, _, err := reader.ReadLine()
		if err != nil {
			log.Fatalf("Unable to read message from input: %v", err)
		}
		err = c.SendMessage(string(textBytes))
		if err != nil {
			fmt.Printf("Unable to send message. Reason: %v", err)
		}
	}
}
