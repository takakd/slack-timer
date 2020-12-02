package lambdahandler

import (
	"context"
	"fmt"
	"slacktimer/internal/app/adapter/notifycontroller"
	"slacktimer/internal/app/util/appinit"
	"slacktimer/internal/app/util/log"
)

// Lambda handler input data
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html
type LambdaInput struct {
	Records []SqsMessage `json:"records"`
}

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

func (s *SqsMessage) HandlerInput() *notifycontroller.HandlerInput {
	return &notifycontroller.HandlerInput{
		UserId: s.Body,
		// TODO: Get userid and message from body.
		Message: "test",
	}
}

// Lambda callback
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func NotifyLambdaHandler(ctx context.Context, input LambdaInput) error {
	appinit.AppInit()

	log.Debug(fmt.Sprintf("handler, input.Records=%v", input.Records))

	count := 0
	for _, m := range input.Records {
		log.Debug(fmt.Sprintf("record %v", m))

		h := notifycontroller.NewHandler()
		i := m.HandlerInput()
		resp := h.Handler(ctx, i)
		if resp.Error != nil {
			count++
		}
	}

	if count > 0 {
		return fmt.Errorf("error happend count=%d", count)
	}

	return nil
}
