package queue

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"slacktimer/internal/app/usecase/enqueueevent"
)

const (
	messageGroupId = "fifo"
)

type SQSWrapper interface {
	SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

type SQSWrapperAdapter struct {
	sqs *sqs.SQS
}

func (s *SQSWrapperAdapter) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return s.sqs.SendMessage(input)
}

type SQSMessageQueue struct {
	wrp SQSWrapper
}

// Set wrp to null. In case unit test, set mock interface.
func NewSQSMessageQueue(wrp SQSWrapper) enqueueevent.Queue {
	if wrp == nil {
		wrp = &SQSWrapperAdapter{
			sqs: sqs.New(session.New()),
		}
	}
	return &SQSMessageQueue{
		wrp: wrp,
	}
}

func (s *SQSMessageQueue) Enqueue(message *enqueueevent.QueueMessage) (string, error) {
	r, err := s.wrp.SendMessage(&sqs.SendMessageInput{
		// TODO: message
		MessageBody:    aws.String(""),
		MessageGroupId: aws.String(messageGroupId),
	})
	if err != nil {
		return "", fmt.Errorf("failed to enqueue %w", err)
	}

	return *r.MessageId, nil
}
