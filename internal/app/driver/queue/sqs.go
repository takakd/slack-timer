// Package queue provides features of AWS SQS that are used in the app.
package queue

import (
	"fmt"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/config"

	"encoding/json"

	"slacktimer/internal/app/util/di"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

const (
	_messageGroupID = "fifo"
)

// Sqs implements Queue interface with SQS.
type Sqs struct {
	wrp SqsWrapper
}

var _ enqueueevent.Queue = (*Sqs)(nil)

// NewSqs creates new struct.
func NewSqs() *Sqs {
	return &Sqs{
		wrp: di.Get("queue.SqsWrapper").(SqsWrapper),
	}
}

// Enqueue enqueues a message to SQS.
func (s Sqs) Enqueue(message enqueueevent.QueueMessage) (string, error) {
	body := NewSqsMessageBody()
	body.UserID = message.UserID
	body.Text = message.Text
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal error %w", err)
	}

	r, err := s.wrp.SendMessage(&sqs.SendMessageInput{
		MessageBody:            aws.String(string(bodyJSON)),
		MessageGroupId:         aws.String(_messageGroupID),
		QueueUrl:               aws.String(config.Get("SQS_URL", "")),
		MessageDeduplicationId: aws.String(newMessageDeduplicationID(message)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to enqueue %w", err)
	}

	return *r.MessageId, nil
}

func newMessageDeduplicationID(message enqueueevent.QueueMessage) string {
	return uuid.New().String()
}
