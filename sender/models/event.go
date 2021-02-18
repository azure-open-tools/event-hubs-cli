package models

import "time"

type EventModel struct {
	Data              string      			 `json:"Data"`
	PartitionKey      interface{} 			 `json:"PartitionKey"`
	Properties        map[string]interface{} `json:"properties"`
	ID                string 				 `json:"ID"`
	SystemProperties  *SystemProperties      `json:"SystemProperties"`
}

type SystemProperties struct {
	SequenceNumber *int64              `json:"x-opt-sequence-number"` // unique sequence number of the message
	EnqueuedTime   *time.Time          `json:"x-opt-enqueued-time"`   // time the message landed in the message queue
	Offset         *int64              `json:"x-opt-offset"`
	PartitionID    *int16              `json:"x-opt-partition-id"` // This value will always be nil. For information related to the event's partition refer to the PartitionKey field in this type
	PartitionKey   *string             `json:"x-opt-partition-key"`
	// Nil for messages other than from Azure IoT Hub. deviceId of the device that sent the message.
	IoTHubDeviceConnectionID *string   `json:"iothub-connection-device-id"`
	// Nil for messages other than from Azure IoT Hub. Used to distinguish devices with the same deviceId, when they have been deleted and re-created.
	IoTHubAuthGenerationID *string     `json:"iothub-connection-auth-generation-id"`
	// Nil for messages other than from Azure IoT Hub. Contains information about the authentication method used to authenticate the device sending the message.
	IoTHubConnectionAuthMethod *string `json:"iothub-connection-auth-method"`
	// Nil for messages other than from Azure IoT Hub. moduleId of the device that sent the message.
	IoTHubConnectionModuleID *string   `json:"iothub-connection-module-id"`
	// Nil for messages other than from Azure IoT Hub. The time the Device-to-Cloud message was received by IoT Hub.
	IoTHubEnqueuedTime *time.Time      `json:"iothub-enqueuedtime"`
}