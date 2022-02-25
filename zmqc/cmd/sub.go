package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Opsi/zmqc/zmqc/zmq"
	"github.com/spf13/cobra"
)

var (
	subTopic   string
	printTopic bool

	subCmd = &cobra.Command{
		Use:   "sub",
		Short: "Subscribe to a topic",
		Run:   subscribe,
	}
)

func subscribe(cmd *cobra.Command, args []string) {
	socket := zmq.NewSocket(zmq.SUB)
	socket.Connect(host, port)

	msgChan := make(chan *zmq.SubMessage)

	socket.StartSubscribe(subTopic, msgChan)

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
	subCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Host to connect to")
	subCmd.Flags().StringVarP(&subTopic, "topic", "t", "", "Topic to subscribe to")
	subCmd.Flags().BoolVarP(&printTopic, "print-topic", "T", false, "Prints the topic of the messages")
}
