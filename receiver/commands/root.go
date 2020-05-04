package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var largs = listenArgs{
	connString           : "",
	consumerGroup        : "",
	dataFilter           : []string{},
	propertyFilter       : []string{},
	partitionIds         : []string{},
	verbose              : false,
	splitOutputIntoFiles : false,
}

var rootCmd = &cobra.Command{
	Use:   "ehrt",
	Long: `
you can use this tool to start a listener (consumer) against Azure Event Hubs.

Straight to the point:

	receiver -c "<conn-string>" 

the command above uses the $Default consumer group, for development environment maybe will not be an issue.
when the listener start, it will print the Data field from the event to the stdout.

in environment which you have many listeners better if you use a specific consumer group to your need, to avoid disconnect 
others, you can use the command with -g flag to use your consumer group for debugging:

	receiver -g "debug-consumer" -c "<conn-string>" 

with a -v flag you can print out the whole Event struct (json)

	receiver -v -g "debug-consumer" -c "<conn-string>" 

if you desire to avoid print all incoming events, you can filter by the data field or properties field like this:

data filter:

	receiver -v -d "my message content filter" -g "debug-consumer" -c "<conn-string>"

property filter:

	receiver -v -p "property filter" -g "debug-consumer" -c "<conn-string>"

you can also provide more than one filter, just providing more than one -d or -p flag, like this:

	receiver -v -p "property filter1" -p "property filter2" -p "property filter3" -g "debug-consumer" -c "<conn-string>"

you can save all the output sending them to a file like this:

	receiver -v -g "<consumer-group>" -c "<conn-string>" > output.json

it will not generate a well formatted json, but each line of that file will be a valid json since you are using -v flag.
with that file you can play alongside the sender (producer) to replay all the file content to an event hub you wish.

write events to multiple json files using -s flag:

	receiver -v -s -f ${PWD}/ -p "<filter1>" -g "<consumer-group>" -c "<conn-string>"

the command above will generate a json file for each event arrived that match the property filter. in case of a path was not
provided the current folder will be used.

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunListen(largs)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&largs.connString, "connstring", "c", "", "event hub listen profile connection string (not the namespace)")
	rootCmd.Flags().StringVarP(&largs.consumerGroup, "consumer-group", "g", "$Default", "consumer group. (recommended flag. otherwise $Default consumer group will be used)")
	rootCmd.Flags().StringSliceVarP(&largs.dataFilter, "data-filter", "d", []string{}, "text to be used to filter data that contain such string. (can be used multiple times on the command. -d <filter1> -d <filter2> -d <filter3>)")
	rootCmd.Flags().StringSliceVarP(&largs.propertyFilter, "property-filter", "p", []string{}, "text to be used to filter event properties that contain such string. (can be used multiple times on the command. -p <filter1> -p <filter2> -p <filter3>)")
	rootCmd.Flags().StringSliceVarP(&largs.partitionIds, "partition-ids", "i", []string{}, "partition id you want to connect to. leave it out and you are going to listen to all partition. ex: -i \"0\" -i \"5\" -i \"18\"")
	rootCmd.Flags().BoolVarP(&largs.verbose, "verbose", "v", false, "print out the whole event(data, properties, system properties) with data field as base64 string")
	rootCmd.Flags().BoolVarP(&largs.splitOutputIntoFiles, "split-output-files", "s", false, "for each event received it will be saved in an isolated json file. Works only with verbose (-v) flag")
	rootCmd.Flags().StringVarP(&largs.pathToWriteFiles, "path-files", "f", "", "path where the files will be written with Split-Output-Files flag. if not provided, current directory(PWD) will be used.")

	_ = rootCmd.MarkFlagRequired("connstring")
	rootCmd.SetVersionTemplate(EventHubAscii)
	rootCmd.Version = Version
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}