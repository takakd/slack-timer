package queue

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SqsWrapperAdapter dispatches to AWS SDK SQS methods.
type SqsWrapperAdapter struct {
	sqs *sqs.SQS
}

var _ SqsWrapper = (*SqsWrapperAdapter)(nil)

// NewSqsWrapperAdapter creates new struct.
func NewSqsWrapperAdapter() *SqsWrapperAdapter {
	return &SqsWrapperAdapter{
		sqs: sqs.New(session.New()),
	}
}

// SendMessage dispatches SDK's method simply.
func (s SqsWrapperAdapter) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return s.sqs.SendMessage(input)
}
