package cmd

import (
	"log"
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
	ctx, _ := signal.NotifyContext(cmd.Context(), os.Interrupt)
	msgChan := make(chan *zmq.PubMessage)
	done := make(chan bool)
	go zmq.StartPublish(ctx, port, msgChan, done)

	var pubMsg string
	if inlinePubMsg != "" {
		pubMsg = inlinePubMsg
	} else if pubPayloadFile != "" {
		data, err := os.ReadFile(pubPayloadFile)
		if err != nil {
			log.Fatal(err)
		}
		pubMsg = string(data)
	}

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
	for {
		select {
		case <-ctx.Done():
			log.Print("Received interrupt signal, exiting...")
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
