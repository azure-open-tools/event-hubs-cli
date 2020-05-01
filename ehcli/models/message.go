package models

type Event struct {
	Data          string
	Base64           bool
	Batch            bool
	NumberOfMessages int
	ReplayMessages   bool
	FileMessagePath  string
	Properties       []string
}