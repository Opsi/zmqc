package cmd

import (
	"log"
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
	ctx, _ := signal.NotifyContext(cmd.Context(), os.Interrupt)
	msgChan := make(chan *zmq.SubMessage)

	go zmq.StartSubscribe(ctx, host, port, subTopic, msgChan)

	for {
		select {
		case <-ctx.Done():
			log.Print("Received interrupt signal, exiting...")
			return
		case msg, ok := <-msgChan:
			if !ok {
				return
			}
			if printTopic {
				log.Printf("Topic '%s':\n%s", msg.Topic, msg.Payload)
			} else {
				log.Printf("%s", msg.Payload)
			}
		}
	}
}

func init() {
	subCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Host to connect to")
	subCmd.Flags().StringVarP(&subTopic, "topic", "t", "", "Topic to subscribe to")
	subCmd.Flags().BoolVarP(&printTopic, "print-topic", "T", false, "Prints the topic of the messages")
}
