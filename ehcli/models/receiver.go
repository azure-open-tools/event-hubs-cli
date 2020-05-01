package models

type Receiver struct {
	ConnectionString     string
	ConsumerGroup        string
	DataFilter           []string
	PropertyFilter       []string
	Verbose              bool
	SplitOutputIntoFiles bool
	PathToWriteFile      string
}