package queue

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SqsWrapper wraps AWS SDK for Unit test.
type SqsWrapper interface {
	SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

// SqsWrapperAdapter dispatches to AWS SDK SQS methods.
type SqsWrapperAdapter struct {
	sqs *sqs.SQS
}

var _ SqsWrapper = (*SqsWrapperAdapter)(nil)

// SendMessage dispatches SDK's method simply.
func (s SqsWrapperAdapter) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return s.sqs.SendMessage(input)
}
