package zmq

import (
	"context"
	"fmt"
	"os"

	"github.com/go-zeromq/zmq4"
)

type PubMessage struct {
	Topic   string
	Payload string
}

func StartPublish(ctx context.Context, port uint, msgChan <-chan *PubMessage, done chan<- bool) {
	pub := zmq4.NewPub(ctx)
	defer pub.Close()

	address := fmt.Sprintf("tcp://*:%d", port)
	if err := pub.Listen(address); err != nil {
		fmt.Printf("Error listening on %s: %s\n", address, err)
		os.Exit(1)
	}

	for msg := range msgChan {
		topicAsBytes := []byte(msg.Topic)
		payloadAsBytes := []byte(msg.Payload)

		msg := zmq4.NewMsgFrom(topicAsBytes, payloadAsBytes)
		if err := pub.Send(msg); err != nil {
			fmt.Printf("Error sending message: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Published a %d bytes message to topic %s\n", len(payloadAsBytes), topicAsBytes)
	}
	done <- true
}
