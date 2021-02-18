package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var largs = senderArgs {
	message            : "",
	base64             : false,
	batch              : false,
	connStr            : "",
	properties         : []string{},
	numberOfMessages   : 0,
	repeat             : 0,
	interval           : 0,
	replayMessages     : false,
	fileMessagePath    : "",
	templateFile       : false,
}

var rootCmd = &cobra.Command{
	Use:   "ehst",
	Long: `
you can use this tool to start a sender (producer) against Azure Event Hubs.

Straight to the point:

ehst .....

`,
	Example: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunSender(largs)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&largs.connStr, "connstring", "c", "", "event hub listen profile connection string (not the namespace)")
	rootCmd.Flags().StringVarP(&largs.message, "message", "m", "", "the message string to send")
	rootCmd.Flags().StringSliceVarP(&largs.properties, "properties", "p", []string{}, "text to be used to filter event properties that contain such string. (can be used multiple times on the command. -p <filter1> -p <filter2> -p <filter3>)")
	rootCmd.Flags().StringVarP(&largs.fileMessagePath, "file-path", "f", "", "path to the file with the messages extracted from EH to be replayed")
	rootCmd.Flags().IntVarP(&largs.numberOfMessages, "number-of-messages", "n", 1, "amount of messages to be send")
	rootCmd.Flags().IntVarP(&largs.repeat, "repeat", "r", 1, "amount of repeats the sender must to send")
	rootCmd.Flags().IntVarP(&largs.interval, "interval", "i", -1, "amount of time among the repetitions in milliseconds")
	rootCmd.Flags().BoolVarP(&largs.base64, "base64", "b", false, "if present the sender will try to decode the base64 to byte[] before send. Useful when you are using protobuffer messages")
	rootCmd.Flags().BoolVarP(&largs.replayMessages, "replay-messages", "e", false, "replay messages")
	rootCmd.Flags().BoolVarP(&largs.batch, "batch", "g", false, "send messages in batches")
	rootCmd.Flags().BoolVarP(&largs.templateFile, "template", "t", false, "uses a template file name in order to send certain amount of events")

	rootCmd.SetVersionTemplate(EventHubAscii)
	rootCmd.Version = Version
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}