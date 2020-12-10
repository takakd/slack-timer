package queue

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SqsWrapper wraps AWS SDK for Unit test.
type SqsWrapper interface {
	SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}
