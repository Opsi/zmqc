package cmd

import (
	"fmt"
	"os"
	"os/signal"

	zmq "github.com/pebbe/zmq4"
	"github.com/spf13/cobra"
)

var (
	topic      string
	printTopic bool

	subCmd = &cobra.Command{
		Use:   "sub",
		Short: "Subscribe to a ZMQ topic",
		Run:   subscribe,
	}
)

type Message struct {
	Topic   string
	Payload string
}

func subscribe(cmd *cobra.Command, args []string) {
	socket, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		fmt.Printf("Error creating socket: %s\n", err)
		return
	}
	if err := socket.SetSubscribe(topic); err != nil {
		fmt.Printf("Error setting topic: %s\n", err)
		return
	}
	address := fmt.Sprintf("tcp://%s:%d", host, port)
	if err := socket.Connect(address); err != nil {
		fmt.Printf("Error connecting to %s: %s\n", address, err)
		return
	}

	msgChan := make(chan Message)
	go func() {
		defer close(msgChan)
		for {
			msg, err := socket.RecvMessage(0)
			if err != nil {
				fmt.Printf("Error receiving message: %s\n", err)
				return
			}
			if len(msg) != 2 {
				fmt.Printf("Invalid message: %v\n", msg)
				return
			}
			msgChan <- Message{
				Topic:   msg[0],
				Payload: msg[1],
			}
		}
	}()

	ctx, _ := signal.NotifyContext(cmd.Context(), os.Interrupt)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Received interrupt signal, exiting...")
			return
		case msg, ok := <-msgChan:
			if !ok {
				return
			}
			if printTopic {
				fmt.Printf("Topic '%s':\n", msg.Topic)
			}
			fmt.Printf("%s\n", msg.Payload)
		}
	}
}

func init() {
	subCmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic to subscribe to")
	subCmd.Flags().BoolVarP(&printTopic, "print-topic", "T", false, "Prints the topic of the messages")
}
