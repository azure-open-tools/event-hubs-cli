package main

import (
	"event-hubs-cli/sender/commands"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	commands.Execute()
}