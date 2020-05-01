package models

type Sender struct {
	ConnectionString   string
	Repeat             int
	Interval           int
	MaxNumberOfThreads int
	MessageSettings    Message
}

