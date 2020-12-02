package cmd

import (
	"crypto/tls"
	"os"

	"github.com/Azure/go-amqp"
	"github.com/juju/loggo"
	"github.com/spf13/cobra"
)

var cfgFile string
var verbose bool
var insecure bool
var amqpConnect string
var anonymous bool

var out = loggo.GetLogger("cmd")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "amqp10-tool",
	Short: "Tool to interact with AMQP 1.0 broker",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			out.SetLogLevel(loggo.DEBUG)
		} else {
			out.SetLogLevel(loggo.INFO)
		}
		out.Debugf("Starting\n")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		out.Criticalf("%s", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&amqpConnect, "broker", "b", "amqp://localhost:5672", "amqp connection string")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false, "skip certificate validation")
	rootCmd.PersistentFlags().BoolVarP(&anonymous, "anonymous", "a", false, "anonymous connection. mandatory for unsecured ActiveMQ")

}

// function to be called on fatal errors, this kills the app
func failOnError(err error, msg string) {
	if err != nil {
		out.Errorf("%s: %s\n", msg, err)
		os.Exit(1)
	}
}

func connect() *amqp.Client {
	options := make([]amqp.ConnOption, 0)
	cfg := new(tls.Config)
	if insecure {
		out.Debugf("insecure flag on, will skip certificate validation")
		cfg.InsecureSkipVerify = true
	}
	//options = append(options, amqp.ConnTLSConfig(cfg))
	if anonymous {
		out.Debugf("anonymous flag on, will send a SASL Anonymous")
		options = append(options, amqp.ConnSASLAnonymous())
	}

	out.Debugf("Connecting to %s\n", amqpConnect)
	conn, err := amqp.Dial(amqpConnect, options...)
	failOnError(err, "Failed to connect to amqp")
	return conn
}
