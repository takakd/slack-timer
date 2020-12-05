package notify

import (
	"context"
	"fmt"
	"slacktimer/internal/app/adapter/notifycontroller"
	"slacktimer/internal/app/util/appinitializer"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

type LambdaHandler interface {
	LambdaHandler(ctx context.Context, input LambdaInput) error
}

type NotifyLambdaHandler struct {
	ctrl notifycontroller.Handler
}

func NewNotifyLambdaHandler() LambdaHandler {
	h := &NotifyLambdaHandler{}
	h.ctrl = di.Get("notify.Handler").(notifycontroller.Handler)
	return h
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
func (s *SqsMessage) HandlerInput() notifycontroller.HandlerInput {
	return notifycontroller.HandlerInput{
		UserId: s.Body,
		// TODO: Get userid and message from body.
		Message: "test",
	}
}

// SQS calls this function.
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func (n *NotifyLambdaHandler) LambdaHandler(ctx context.Context, input LambdaInput) error {
	appinitializer.AppInit()

	log.Info(fmt.Sprintf("lambda handler input count=%d, recourds=%v", len(input.Records), input.Records))

	count := 0
	for _, m := range input.Records {
		resp := n.ctrl.Handler(ctx, m.HandlerInput())
		if resp.Error != nil {
			count++
		}
	}

	var err error

	if count > 0 {
		err = fmt.Errorf("count=%d", count)
	}

	log.Info("lambda handler output", err)

	return err
}
