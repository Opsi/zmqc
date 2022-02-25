package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Opsi/zmqc/zmqc/zmq"
	"github.com/spf13/cobra"
)

var (
	pubTopic       string
	inlinePubMsg   string
	pubPayloadFile string
	repeatFreq     time.Duration

	pubCmd = &cobra.Command{
		Use:   "pub",
		Short: "Publish to a topic",
		Run:   publish,
	}
)

func publish(cmd *cobra.Command, args []string) {
	socket := zmq.NewSocket(zmq.PUB)
	socket.Bind(port)

	var pubMsg string
	if inlinePubMsg != "" {
		pubMsg = inlinePubMsg
	} else if pubPayloadFile != "" {
		data, err := os.ReadFile(pubPayloadFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pubMsg = string(data)
	}

	msgChan := make(chan *zmq.PubMessage)
	done := make(chan bool)
	go socket.StartPublish(msgChan, done)

	msgChan <- &zmq.PubMessage{
		Topic:   pubTopic,
		Payload: pubMsg,
	}

	if repeatFreq <= 0 {
		close(msgChan)
		<-done
		return
	}

	ticker := time.NewTicker(repeatFreq)
	ctx, _ := signal.NotifyContext(cmd.Context(), os.Interrupt)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Received interrupt signal, exiting...")
			return
		case <-ticker.C:
			msgChan <- &zmq.PubMessage{
				Topic:   pubTopic,
				Payload: pubMsg,
			}
		}
	}
}

func init() {
	pubCmd.Flags().StringVarP(&pubTopic, "topic", "t", "", "Topic to publish to")
	pubCmd.Flags().StringVarP(&inlinePubMsg, "message", "m", "", "Message to publish")
	pubCmd.Flags().StringVarP(&pubPayloadFile, "file", "f", "", "File to read payload from")
	pubCmd.Flags().DurationVarP(&repeatFreq, "freq", "q", 0, "Send Payload over and over in the given frequency")
}