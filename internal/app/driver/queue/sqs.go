// package queue provides features of AWS SQS that are used in the app.
package queue

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/config"
)

const (
	messageGroupId = "fifo"
)

type Sqs struct {
	wrp SqsWrapper
}

// TODO: not null parameter, get from DI in the function
// Set wrp to null. In case unit test, set mock interface.
func NewSqs(wrp SqsWrapper) enqueueevent.Queue {
	if wrp == nil {
		wrp = &SqsWrapperAdapter{
			sqs: sqs.New(session.New()),
		}
	}
	return &Sqs{
		wrp: wrp,
	}
}

func (s Sqs) Enqueue(message enqueueevent.QueueMessage) (string, error) {
	r, err := s.wrp.SendMessage(&sqs.SendMessageInput{
		MessageBody:    aws.String(message.UserId),
		MessageGroupId: aws.String(messageGroupId),
		QueueUrl:       aws.String(config.Get("SQS_URL", "")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to enqueue %w", err)
	}

	return *r.MessageId, nil
}