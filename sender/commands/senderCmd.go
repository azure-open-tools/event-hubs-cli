package commands

import (
	"bufio"
	"context"
	"encoding/json"
	"event-hubs-cli/sender/common"
	"event-hubs-cli/sender/models"
	"fmt"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/azure-open-tools/event-hubs/sender"
	"github.com/vbauerster/mpb/v5"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type (
	 senderArgs struct {
		message            string
		base64             bool
		batch              bool
		connStr            string
		properties         []string
		numberOfMessages   int
		repeat             int
		interval           int
		replayMessages     bool
		fileMessagePath    string
	}

	 senderCli struct {
		numberOfMessages   int64
		sender             *sender.Sender
		sendWaitGroup      *sync.WaitGroup
		sendProgress       *mpb.Progress
		sendBar            *mpb.Bar
		sendBatchBar	   map[int]*mpb.Bar
		start              time.Time
	}
)

var (
	mCli *senderCli
)

func RunSender(args senderArgs) error {
	var err error

	if args.replayMessages {
		err = newSenderCli(args.connStr, args.properties, args.base64, int64(args.numberOfMessages))
		if err == nil {
			return mCli.replayMessage(args.fileMessagePath, args.repeat, args.interval)
		}
	} else {
		err = newSenderCli(args.connStr, args.properties, args.base64, int64(args.numberOfMessages))
		if err == nil {
			return mCli.sendMessage(args.message, args.batch, args.repeat, args.interval)
		}
	}

	return err
}

func newSenderCli(connStr string, properties []string, base64 bool, numberOfMessages int64) error {
	builder := sender.NewSenderBuilder()
	builder.SetConnectionString(getConnString(connStr))
	builder.AddProperties(properties)
	builder.SetBase64(base64)
	builder.SetNumberOfMessages(numberOfMessages)
	builder.SetOnAfterSendMessage(OnAfterSendMessage)
	builder.SetOnBeforeSendMessage(OnBeforeSendMessage)
	builder.SetOnAfterSendBatchMessage(OnAfterSendBatchMessage)
	builder.SetOnBeforeSendBatchMessage(OnBeforeSendBatchMessage)

	snd, err := builder.GetSender()
	if err == nil {
		cli := &senderCli{
			sender:             snd,
			numberOfMessages:   numberOfMessages,
			sendWaitGroup:      &sync.WaitGroup{},
			sendBatchBar: 		make(map[int]*mpb.Bar),
		}
		cli.sendProgress = mpb.New(mpb.WithWidth(64), mpb.WithWaitGroup(cli.sendWaitGroup))
		mCli = cli

		return nil
	}

	return err
}

func getConnString(connString string) string {
	if len(strings.TrimSpace(connString)) > 0 {
		return connString
	} else {
		return os.Getenv("EVENTHUB_SEND_CONNSTR")
	}
}

func (cli *senderCli) sendMessage(message string, batch bool, repeat int, interval int) error {
	if !batch {
		return cli.send(message, repeat, interval)
	} else {
		return cli.sendBatch(message, repeat, interval)
	}
}

func (cli *senderCli) send(message string, repeat int, interval int) error {
	var err error

	if repeat == 0 {
		repeat = 1
	}

	cli.start = time.Now()
	cli.sendWaitGroup.Add(repeat)
	for i := 1; i <= repeat; i++ {
		cli.sendBar = common.GetBar(cli.numberOfMessages, i, cli.sendProgress)
		err = cli.sender.SendMessage(message, context.Background())

		if err != nil {
			return err
		}
		cli.sendWaitGroup.Done()

		if interval > 0 {
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}

	cli.sendProgress.Wait()

	return err
}

func (cli *senderCli) sendBatch(message string, repeat int, interval int) error {
	var err error

	if repeat == 0 {
		repeat = 1
	}

	cli.start = time.Now()
	cli.sendWaitGroup.Add(repeat)
	for i := 0; i < repeat; i++ {
		err = cli.sender.SendBatchMessage(message, context.Background())

		if err != nil {
			return err
		}

		cli.sendWaitGroup.Done()
		if interval > 0 {
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}

	cli.sendProgress.Wait()

	return err
}

func (cli *senderCli) replayMessage(filePath string, repeat int, interval int) error {
	var err error

	for i := 1; i <= repeat; i++ {
		err = cli.replayMessageFile(filePath)
		if err != nil {
			return err
		}

		if interval > 0 {
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}

	return nil
}

func OnBeforeSendMessage(*eventhub.Event) {

}

func OnAfterSendMessage(event *eventhub.Event) {
	if event != nil && mCli.sendBar != nil{
		mCli.sendBar.Increment()
		mCli.sendBar.DecoratorEwmaUpdate(time.Since(mCli.start))
		//mCli.sendBar.DecoratorAverageAdjust(mCli.start)
	}
}

func OnBeforeSendBatchMessage(batchSize int, workerIndex int) {
	if mCli != nil && mCli.sendBatchBar != nil {
		if _, exist := mCli.sendBatchBar[workerIndex]; !exist {
			mCli.sendBatchBar[workerIndex] = common.GetBar(int64(batchSize), workerIndex, mCli.sendProgress)
		}
	}
}

func OnAfterSendBatchMessage(batchSizeSent int, workerIndex int) {
	if mCli != nil && mCli.sendBatchBar != nil {
		if _, exist := mCli.sendBatchBar[workerIndex]; exist {
			mCli.sendBatchBar[workerIndex].IncrBy(batchSizeSent)
			mCli.sendBatchBar[workerIndex].DecoratorEwmaUpdate(time.Since(mCli.start))
		}
	}
}

func (cli *senderCli) replayMessageFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	defer common.CloseFile(file)

	if err != nil {
		return err
	}

	rd := bufio.NewReader(file)

	for err != io.EOF {
		line, err := rd.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		model := models.EventModel{}
		_ = json.Unmarshal([]byte(line), &model)

		cli.sender.AddProperties(model.Properties)
		err = cli.sender.SendMessage(model.Data, context.Background())

		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}