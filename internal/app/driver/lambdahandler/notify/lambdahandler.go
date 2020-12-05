package notify

import (
	"context"
	"slacktimer/internal/app/adapter/notify"
)

type LambdaHandler interface {
	Handle(ctx context.Context, input LambdaInput) error
}

// Lambda handler input data
// SQS passes this.
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html
type LambdaInput struct {
	// Lambda handler parameters include multiple SQS messages.
	Records []SqsMessage `json:"records"`
}

// One SQS message in handler parameters.
type SqsMessage struct {
	MessageId     string            `json:"messageId"`
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

// To a controller input data.
func (s *SqsMessage) HandleInput() notify.HandleInput {
	return notify.HandleInput{
		UserId: s.Body,
		// TODO: Get userid and message from body.
		Message: "test",
	}
}
