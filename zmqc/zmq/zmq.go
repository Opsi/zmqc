package zmq

import (
	"fmt"
	"os"

	zmq4 "github.com/pebbe/zmq4"
)

type Socket struct {
	_socket zmq4.Socket
}

type SocketType uint8

const (
	SUB SocketType = iota
	PUB
)

type SubMessage struct {
	Topic   string
	Payload string
}

type PubMessage struct {
	Topic   string
	Payload string
}

func (st SocketType) convert() zmq4.Type {
	switch st {
	case SUB:
		return zmq4.SUB
	case PUB:
		return zmq4.PUB
	default:
		panic("Unknown socket type: " + string(st))
	}
}

func NewSocket(t SocketType) *Socket {
	socket, err := zmq4.NewSocket(t.convert())
	if err != nil {
		fmt.Printf("Error creating socket: %s\n", err)
		os.Exit(1)
	}
	return &Socket{
		_socket: *socket,
	}
}

func (s *Socket) Connect(host string, port uint) {
	address := fmt.Sprintf("tcp://%s:%d", host, port)
	if err := s._socket.Connect(address); err != nil {
		fmt.Printf("Error connecting to %s: %s\n", address, err)
		os.Exit(1)
	}
}

func (s *Socket) StartSubscribe(topic string, msgChan chan<- *SubMessage) {
	if err := s._socket.SetSubscribe(topic); err != nil {
		fmt.Printf("Error setting topic: %s\n", err)
		os.Exit(1)
	}
	go func() {
		for {
			msg, err := s._socket.RecvMessage(0)
			if err != nil {
				fmt.Printf("Error receiving message: %s\n", err)
				os.Exit(1)
			}
			if len(msg) != 2 {
				fmt.Printf("Invalid message: %v\n", msg)
				os.Exit(1)
			}
			msgChan <- &SubMessage{
				Topic:   msg[0],
				Payload: msg[1],
			}
		}
	}()
}

func (s *Socket) Bind(port uint) {
	address := fmt.Sprintf("tcp://*:%d", port)
	if err := s._socket.Bind(address); err != nil {
		fmt.Printf("Error binding to %s: %s\n", address, err)
		os.Exit(1)
	}
}

func (s *Socket) StartPublish(msgChan <-chan *PubMessage, done chan<- bool) {
	for msg := range msgChan {
		topicAsBytes := []byte(msg.Topic)
		payloadAsBytes := []byte(msg.Payload)
		_, err := s._socket.SendBytes(topicAsBytes, zmq4.SNDMORE)
		if err != nil {
			fmt.Printf("Error sending message: %s\n", err)
			os.Exit(1)
		}
		_, err = s._socket.SendBytes(payloadAsBytes, 0)
		if err != nil {
			fmt.Printf("Error sending message: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Published a %d bytes message to topic '%s'\n", len(payloadAsBytes), msg.Topic)
	}
	done <- true
}
