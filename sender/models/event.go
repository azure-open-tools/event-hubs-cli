package models

import "time"

type EventModel struct {
	Data         string      			`json:"Data"`
	PartitionKey interface{} 			`json:"PartitionKey"`
	Properties   map[string]interface{} `json:"properties"`
	ID           string 				`json:"ID"`
	SystemProperties struct {
		SequenceNumber int         		`json:"SequenceNumber"`
		EnqueuedTime   time.Time   		`json:"EnqueuedTime"`
		Offset         int64       		`json:"Offset"`
		PartitionID    interface{} 		`json:"PartitionID"`
		PartitionKey   interface{} 		`json:"PartitionKey"`
	} `json:"SystemProperties"`
}