// Package queue provides features of AWS SQS that are used in the app.
package queue

import (
	"fmt"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	_messageGroupID = "fifo"
)

// Sqs implements Queue interface with SQS.
type Sqs struct {
	wrp SqsWrapper
}

var _ enqueueevent.Queue = (*Sqs)(nil)

// NewSqs create new struct.
// TODO: not null parameter, get from DI in the function
// Set wrp to null. In case unit test, set mock interface.
func NewSqs(wrp SqsWrapper) *Sqs {
	if wrp == nil {
		wrp = &SqsWrapperAdapter{
			sqs: sqs.New(session.New()),
		}
	}
	return &Sqs{
		wrp: wrp,
	}
}

// Enqueue enqueues a message to SQS.
func (s Sqs) Enqueue(message enqueueevent.QueueMessage) (string, error) {
	r, err := s.wrp.SendMessage(&sqs.SendMessageInput{
		MessageBody:    aws.String(message.UserId),
		MessageGroupId: aws.String(_messageGroupID),
		QueueUrl:       aws.String(config.Get("SQS_URL", "")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to enqueue %w", err)
	}

	return *r.MessageId, nil
}
