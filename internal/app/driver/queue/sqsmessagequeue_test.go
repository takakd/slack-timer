package queue

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/usecase/enqueueevent"
	"testing"
)

func TestNewSQSMessageQueue(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		q := NewSQSMessageQueue(nil)
		concrete, ok := q.(*SQSMessageQueue)
		assert.True(t, ok)
		assert.IsType(t, &SQSWrapperAdapter{}, concrete.wrp)
	})

	t.Run("mock", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockSQSWrapper(ctrl)
		repo := NewSQSMessageQueue(mock)
		concrete, ok := repo.(*SQSMessageQueue)
		assert.True(t, ok)
		assert.IsType(t, mock, concrete.wrp)
	})
}

func TestSQSMessageQueue_Enqueue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		caseMessage := &enqueueevent.QueueMessage{
			"id1",
		}
		caseMessageInput := &sqs.SendMessageInput{
			MessageBody:    aws.String(""),
			MessageGroupId: aws.String(messageGroupId),
		}
		caseMessageOutput := &sqs.SendMessageOutput{
			MessageId: aws.String("msgid1"),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := NewMockSQSWrapper(ctrl)
		w.EXPECT().SendMessage(caseMessageInput).Return(caseMessageOutput, nil)

		q := NewSQSMessageQueue(w)
		r, err := q.Enqueue(caseMessage)
		assert.Equal(t, *caseMessageOutput.MessageId, r)
		assert.NoError(t, err)
	})

	t.Run("ng:failed", func(t *testing.T) {
		caseMessage := &enqueueevent.QueueMessage{
			"id1",
		}
		caseMessageInput := &sqs.SendMessageInput{
			MessageBody:    aws.String(""),
			MessageGroupId: aws.String(messageGroupId),
		}
		caseError := errors.New("error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := NewMockSQSWrapper(ctrl)
		w.EXPECT().SendMessage(caseMessageInput).Return(nil, caseError)

		q := NewSQSMessageQueue(w)
		r, err := q.Enqueue(caseMessage)
		assert.Empty(t, r)
		assert.Equal(t, fmt.Errorf("failed to enqueue %w", caseError), err)
	})
}
