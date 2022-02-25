package zmq

import (
	"context"
	"fmt"
	"os"

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
		fmt.Printf("Error connecting to %s: %s\n", address, err)
		os.Exit(1)
	}

	if err := sub.SetOption(zmq4.OptionSubscribe, topic); err != nil {
		fmt.Printf("Error subscribing to topic: %s\n", err)
		os.Exit(1)
	}
	for {
		msg, err := sub.Recv()
		if err != nil {
			fmt.Printf("Error receiving message: %s\n", err)
			os.Exit(1)
		}
		if len(msg.Frames) != 2 {
			fmt.Printf("Invalid message: %v\n", msg)
			os.Exit(1)
		}
		msgChan <- &SubMessage{
			Topic:   msg.Frames[0],
			Payload: msg.Frames[1],
		}
	}
}
