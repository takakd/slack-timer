package notify

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

func (s *SqsMessage) HandlerInput() notifycontroller.HandlerInput {
	return notifycontroller.HandlerInput{
		UserId: s.Body,
		// TODO: Get userid and message from body.
		Message: "test",
	}
}

// Lambda callback
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func NotifyLambdaHandler(ctx context.Context, input LambdaInput) error {
	appinit.AppInit()

	log.Info(fmt.Sprintf("lambda handler input count=%d, recourds=%v", len(input.Records), input.Records))

	count := 0
	for _, m := range input.Records {
		h := notifycontroller.NewHandler()
		resp := h.Handler(ctx, m.HandlerInput())
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
