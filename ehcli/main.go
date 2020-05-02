package main

import "fmt"

import (
	"github.com/azure-open-tools/event-hubs/receiver"
	"github.com/azure-open-tools/event-hubs/sender"
)

func main() {
	_ = sender.NewSenderBuilder()
	_ = receiver.NewReceiverBuilder()

	//_ = receiver.ListenMessages("", "", false, "", "", false)
	fmt.Println("define yaml struct to define tests")
	fmt.Println("implement this command line to be used as an unique binary along side with test cases.")
}
