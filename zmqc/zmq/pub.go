package zmq

import (
	"context"
	"fmt"

	"github.com/Opsi/zmqc/zmqc/logger"
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
		logger.Fatalf("Error listening on %s: %s", address, err)
	}

	for msg := range msgChan {
		topicAsBytes := []byte(msg.Topic)
		payloadAsBytes := []byte(msg.Payload)

		msg := zmq4.NewMsgFrom(topicAsBytes, payloadAsBytes)
		if err := pub.Send(msg); err != nil {
			logger.Fatalf("Error sending message: %s", err)
		}
		logger.Infof("Published a %d bytes message to topic %s", len(payloadAsBytes), topicAsBytes)
	}
	done <- true
}
