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

type ChatClient interface {
	Dial(address string) error
	Send(command interface{}) error
	SendMessage(message string) error
	GetName() string
	SetName(name string) error
	Start()
	Close()
	Incoming() chan interface{}
}

type TcpChatClient struct {
	conn      net.Conn
	cmdReader *stream.CommandReader
	cmdWriter *stream.CommandWriter
	name      string
	incoming  chan interface{}
}

func NewClient() ChatClient {
	return &TcpChatClient{
		incoming: make(chan interface{}),
	}
}

func (c *TcpChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.conn = conn
	c.cmdReader = stream.NewCommandReader(conn)
	c.cmdWriter = stream.NewCommandWriter(conn)
	return nil
}

func (c *TcpChatClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

func (c *TcpChatClient) SendMessage(message string) error {
	return c.Send(protocol.SendCommand{
		Message: message,
	})
}

func (c *TcpChatClient) GetName() string {
	return c.name
}

func (c *TcpChatClient) SetName(name string) error {
	c.name = name
	return c.Send(protocol.NameCommand{
		Name: name,
	})
}

func (c *TcpChatClient) Start() {
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

func (c *TcpChatClient) Close() {
	c.conn.Close()
}

func (c *TcpChatClient) Incoming() chan interface{} {
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
	c.SetName(text[:len(text)-1])

	fmt.Printf("Wellcome %v\n", text)
	go listenOnIncomingMessages(c)
	for {
		//fmt.Printf("--> ")
		textBytes, _, err := reader.ReadLine()
		if err != nil {
			log.Fatalf("Unable to read message from input: %v", err)
		}
		c.SendMessage(string(textBytes))
	}
}
