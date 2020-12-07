// Package notify the entry point of the notify lambda function.
package notify

import (
	"context"
	"slacktimer/internal/app/adapter/notify"
)

// LambdaHandler defines the interface called by AWS Lambda.
type LambdaHandler interface {
	Handle(ctx context.Context, input LambdaInput) error
}

// LambdaInput is passed from SQS.
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html
type LambdaInput struct {
	// Parameters include multiple SQS messages.
	Records []SqsMessage `json:"records"`
}

// SqsMessage is one of SQS message in handler parameters.
type SqsMessage struct {
	MessageID     string            `json:"messageId"`
	ReceiptHandle string            `json:"receiptHandle"`
	Body          string            `json:"body"`
	Attributes    map[string]string `json:"attributes"`
	// TODO: check schema
	MessageAttributes map[string]interface{} `json:"messageAttributes"`
	MD5OfBody         string                 `json:"md5OfBody"`
	EventSource       string                 `json:"eventSource"`
	EventSourceArn    string                 `json:"eventSourceARN"`
	AwsRegion         string                 `json:"awsRegion"`
}

// HandleInput convert to the data for controller.
func (s SqsMessage) HandleInput() notify.HandleInput {
	return notify.HandleInput{
		UserID: s.Body,
		// TODO: Get userid and message from body.
		Message: "test",
	}
}
