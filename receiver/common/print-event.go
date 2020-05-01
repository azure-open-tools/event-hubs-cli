package common

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"strings"
)

func PrintEvent(event *eventhub.Event, verbose bool, splitFiles bool, path string) error {
	if verbose {
		return printVerbose(event, splitFiles, path)
	} else {
		return printData(event)
	}
}

func printVerbose(event *eventhub.Event, splitFiles bool, path string) error {
	jsn, err := json.Marshal(event)

	if err == nil {
		if splitFiles {
			err = writeToFile(event, path, jsn)
			if err != nil {
				return err
			}
		} else {
			fmt.Println(string(jsn))
		}
	}

	return err
}

func writeToFile(event *eventhub.Event, path string, jsn []byte) error {
	var err error
	if len(path) == 0 || len(strings.TrimSpace(path)) == 0 {
		path, err = os.Getwd()
		if err != nil {
			return err
		}
		path = path + "/"
	}

	file, _ := json.MarshalIndent(jsn, "", " ")
	err = ioutil.WriteFile(path+"event-id_"+event.ID+".json", file, 0644)
	if err == nil {
		if len(strings.TrimSpace(event.ID)) == 0 {
			event.ID = uuid.New().String()
		}
		fmt.Println(path + "event-id_" + event.ID + ".json")
	}

	return err
}

func printData(event *eventhub.Event) error {
	data := string(event.Data)
	dec, errDecode := b64.StdEncoding.DecodeString(data)

	if errDecode == nil {
		fmt.Println(string(dec))
	} else {
		fmt.Println(data)
	}

	return errDecode
}