package sqs

import "time"

// SQSEventRecord represents a record in an SQS event.
// It is used to unmarshal an SNS event from AWS.
type SQSEventRecord struct {
	Type             string    `json:"Type"`             // The type of the event (e.g., "Notification").
	MessageId        string    `json:"MessageId"`        // The unique identifier for the message.
	TopicArn         string    `json:"TopicArn"`         // The ARN of the SNS topic that published the message.
	Message          string    `json:"Message"`          // The actual message content.
	Timestamp        time.Time `json:"Timestamp"`        // The time when the message was published.
	SignatureVersion string    `json:"SignatureVersion"` // The version of the signature used.
	Signature        string    `json:"Signature"`        // The signature of the message.
	SigningCertURL   string    `json:"SigningCertURL"`   // The URL to the certificate used to sign the message.
	UnsubscribeURL   string    `json:"UnsubscribeURL"`   // The URL to unsubscribe from the topic.
}
