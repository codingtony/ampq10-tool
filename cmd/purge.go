package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-amqp"
	"github.com/spf13/cobra"
)

var queueName string
var timeout int

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Read all the messages that are in a queue",
	PreRun: func(cmd *cobra.Command, args []string) {
		if queueName == "" {
			failOnError(fmt.Errorf("queue name is mandatory"), "Missing parameter")
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		conn := connect()
		//defer conn.Close()
		session, err := conn.NewSession()
		ctx := context.Background()
		failOnError(err, "Failed to open a channel")

		receiver, err := session.NewReceiver(
			amqp.LinkSourceAddress(queueName),
			amqp.LinkCredit(100),
			amqp.LinkBatching(true),
		)
		count := 0
		defer func() {
			out.Infof("Number of messages read from %s : %d\n", queueName, count)
		}()
		out.Debugf("Reading messages from %s\n", queueName)
		for {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
			defer cancel()
			msg, err := receiver.Receive(ctxWithTimeout)
			if err != nil {
				out.Errorf("Reading message from AMQP:%s\n", err)
			}
			// Accept message
			msg.Accept(context.Background())
			count++
		}

	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)
	purgeCmd.Flags().StringVarP(&queueName, "queue", "q", "", "Queue name")
	purgeCmd.Flags().IntVar(&timeout, "timeout", 5, "time to wait in seconds for new messages")
}
