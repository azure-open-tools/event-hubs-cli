package commands

import (
	"bufio"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"event-hubs-cli/sender/common"
	"event-hubs-cli/sender/models"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/azure-open-tools/event-hubs/sender"
	"github.com/google/uuid"
	"github.com/vbauerster/mpb/v5"
)

type (
	senderArgs struct {
		message          string
		base64           bool
		batch            bool
		connStr          string
		properties       []string
		numberOfMessages int
		repeat           int
		interval         int
		replayMessages   bool
		fileMessagePath  string
		templateFile     bool
	}

	senderCli struct {
		numberOfMessages int64
		sender           *sender.Sender
		sendWaitGroup    *sync.WaitGroup
		sendProgress     *mpb.Progress
		sendBar          *mpb.Bar
		sendBatchBar     map[int]*mpb.Bar
		start            time.Time
	}
)

var (
	mCli *senderCli
)

func RunSender(args senderArgs) error {
	err := newSenderCli(args.connStr, args.properties, args.base64, int64(args.numberOfMessages))
	if err == nil {
		if args.replayMessages {
			return mCli.replayMessage(args.fileMessagePath, args.repeat, args.interval)
		} else if args.templateFile {
			return mCli.templateMessage(args.fileMessagePath, args.repeat, args.interval)
		} else {
			if len(strings.TrimSpace(args.message)) == 0 {
				return errors.New("you must to provide a content to the -m/--message parameter")
			}

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
			sender:           snd,
			numberOfMessages: numberOfMessages,
			sendWaitGroup:    &sync.WaitGroup{},
			sendBatchBar:     make(map[int]*mpb.Bar),
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
		start := time.Now()

		cli.sendBar = common.GetBar(cli.numberOfMessages, i, cli.sendProgress)
		err = cli.sender.SendMessage(message, context.Background())

		if err != nil {
			fmt.Println(err)
		}
		cli.sendWaitGroup.Done()

		if interval > 0 {
			sleep(start, interval)
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
		start := time.Now()

		err = cli.sender.SendBatchMessage(message, context.Background())

		if err != nil {
			fmt.Println(err)
		}

		cli.sendWaitGroup.Done()
		if interval > 0 {
			sleep(start, interval)
		}
	}

	cli.sendProgress.Wait()

	return err
}

func (cli *senderCli) replayMessage(filePath string, repeat int, interval int) error {
	var err error

	for i := 1; i <= repeat; i++ {
		start := time.Now()

		err = cli.replayMessageFile(filePath)
		if err != nil {
			return err
		}

		if interval > 0 {
			sleep(start, interval)
		}
	}

	return nil
}

func (cli *senderCli) templateMessage(filePath string, repeat int, interval int) error {
	var err error

	for i := 1; i <= repeat; i++ {
		start := time.Now()

		err = cli.templateMessageFile(filePath)
		if err != nil {
			return err
		}

		if interval > 0 {
			sleep(start, interval)
		}
	}

	return nil
}

func sleep(start time.Time, interval int) {
	elapsed := time.Since(start)
	timeToSleep := (time.Duration(interval) * time.Millisecond) - elapsed
	time.Sleep(timeToSleep)
}

func OnBeforeSendMessage(*eventhub.Event) {

}

func OnAfterSendMessage(event *eventhub.Event) {
	if event != nil && mCli.sendBar != nil {
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

func findPayloadField(fields []string) (string, error) {
	for _, field := range fields {
		if strings.Contains(field, ":") {
			keyVal := strings.Split(field, ":")

			if keyVal[0] == "payload" {
				return keyVal[1], nil
			}
		}
	}
	return "", errors.New("payload not found")
}

func findPropertyFields(fields []string) map[string]string {
	result := make(map[string]string)
	for _, field := range fields {
		if strings.Contains(field, ":") {
			keyVal := strings.Split(field, ":")

			if keyVal[0] != "payload" {
				result[keyVal[0]] = keyVal[1]
			}
		}
	}
	return result
}

func (cli *senderCli) templateMessageFile(filePath string) error {
	var eofErr error
	var line string
	var event *eventhub.Event
	var events []*eventhub.Event
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	defer common.CloseFile(file)

	if err != nil {
		return err
	}

	rd := bufio.NewReader(file)

	for eofErr != io.EOF {
		line, eofErr = rd.ReadString('\n')

		// if the line contains something
		if len(line) > 0 {

			guid := uuid.New().String()

			fields := strings.Split(line, ";")
			payload, pErr := findPayloadField(fields)

			if pErr != nil {
				fmt.Println("Payload not found. Skipping message")
				continue
			}

			event = createAnEvent(true, payload)
			event.ID = guid

			if len(fields) > 1 {
				properties := findPropertyFields(fields)

				for k, v := range properties {
					if strings.Contains(v, "[epoch]") {
						v = strings.ReplaceAll(v, "[epoch]", strconv.FormatInt(time.Now().Unix(), 10))
					}
					if strings.Contains(v, "[guid]") {
						v = strings.ReplaceAll(v, "[guid]", guid)
					}

					// add deviceId additionally to SystemProperties
					if k == "deviceId" {
						deviceId := v
						if event.SystemProperties == nil {
							event.SystemProperties = &eventhub.SystemProperties{}
						}
						event.SystemProperties.IoTHubDeviceConnectionID = &deviceId
					}
					event.Set(k, v)
				}
			}

			events = append(events, event)
		}
	}

	return cli.sender.SendEventsAsBatch(context.Background(), &events)
}

func createAnEvent(base64 bool, message string) *eventhub.Event {
	var event *eventhub.Event

	if base64 {
		decoded, err := b64.StdEncoding.DecodeString(message)
		if err == nil {
			event = eventhub.NewEvent(decoded)
		} else {
			event = eventhub.NewEvent([]byte(message))
		}
	} else {
		event = eventhub.NewEvent([]byte(message))
	}

	return event
}
