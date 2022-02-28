package zmq

import (
	"context"
	"fmt"
	"log"

	"github.com/go-zeromq/zmq4"
)

type SubMessage struct {
	Topic   []byte
	Payload []byte
}

func StartSubscribe(ctx context.Context, host string, port uint, topic string, msgChan chan<- *SubMessage) {
	sub := zmq4.NewSub(ctx)
	defer sub.Close()

	address := fmt.Sprintf("tcp://%s:%d", host, port)
	if err := sub.Dial(address); err != nil {
		log.Fatalf("Error connecting to %s: %s", address, err)
	}

	if err := sub.SetOption(zmq4.OptionSubscribe, topic); err != nil {
		log.Fatalf("Error subscribing to topic: %s", err)
	}
	for {
		msg, err := sub.Recv()
		if err != nil {
			log.Fatalf("Error receiving message: %s", err)
		}
		if len(msg.Frames) != 2 {
			log.Fatalf("Invalid message: %v", msg)
		}
		msgChan <- &SubMessage{
			Topic:   msg.Frames[0],
			Payload: msg.Frames[1],
		}
	}
}
