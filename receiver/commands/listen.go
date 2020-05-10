package commands

import (
	"context"
	"errors"
	"event-hubs-cli/receiver/common"
	"fmt"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/azure-open-tools/event-hubs/receiver"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type (
	listenArgs struct {
		connString           string
		consumerGroup        string
		dataFilter           []string
		propertyFilter       []string
		partitionIds		 []string
		verbose	             bool
		splitOutputIntoFiles bool
		pathToWriteFiles     string
	}

	receiverCli struct {
		verbose              bool
		splitOutputIntoFiles bool
		pathToWriteFiles     string
	}
)

var (
	mCli *receiverCli
)

func RunListen(args listenArgs) error {
	builder := receiver.NewReceiverBuilder()

	if builder != nil {
		builder.AddDataFilters(args.dataFilter)
		builder.AddPropertyFilters(args.propertyFilter)
		builder.SetConnectionString(getConnString(args.connString))
		builder.SetConsumerGroup(args.consumerGroup)
		builder.AddListenerPartitionIds(args.partitionIds)
		builder.SetReceiverHandler(OnReceiverHandler)

		rcv, err := builder.GetReceiver()
		if rcv != nil && err == nil {
			mCli = &receiverCli{
				verbose: args.verbose,
				splitOutputIntoFiles: args.splitOutputIntoFiles,
				pathToWriteFiles: args.pathToWriteFiles,
			}

			ctx := context.Background()
			startListener(args, rcv, ctx)

			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			<-signalChan

			stopListener(rcv, ctx)
		} else {
			return err
		}
	}

	return errors.New("was not possible to start the listener")
}

func getConnString(connString string) string {
	if len(strings.TrimSpace(connString)) > 0 {
		return connString
	} else {
		return os.Getenv("EVENTHUB_LISTEN_CONNSTR")
	}
}

func startListener(args listenArgs, rcv *receiver.Receiver, ctx context.Context) {
	fmt.Println("starting listeners, this can take a couple of seconds...")

	err := rcv.StartListener(ctx)
	if err != nil {
		fmt.Println(err)
		stopListener(rcv, ctx)
	} else {
		fmt.Printf("listening with consumer group: %v\n\n", args.consumerGroup)
	}
}

func stopListener(rcv *receiver.Receiver, ctx context.Context) {
	fmt.Println("\nstopping listeners")

	err := rcv.StopListener(ctx)
	if err == nil {
		fmt.Println("finished")
	} else {
		fmt.Println(err)
	}
}

func OnReceiverHandler(_ context.Context, event *eventhub.Event) error {
	if mCli != nil {
		return common.PrintEvent(event, mCli.verbose, mCli.splitOutputIntoFiles, mCli.pathToWriteFiles)
	}

	return nil
}