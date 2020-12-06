package queue

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Wrap AWS SDK for Unit test.
type SqsWrapper interface {
	SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

type SqsWrapperAdapter struct {
	sqs *sqs.SQS
}

func (s SqsWrapperAdapter) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return s.sqs.SendMessage(input)
}
