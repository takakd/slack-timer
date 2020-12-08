package queue

import (
	"errors"
	"fmt"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/config"
	"testing"

	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewSqs(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		concrete := NewSqs(nil)
		assert.IsType(t, &SqsWrapperAdapter{}, concrete.wrp)
	})

	t.Run("mock", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockSqsWrapper(ctrl)
		concrete := NewSqs(mock)
		assert.IsType(t, mock, concrete.wrp)
	})
}

type MessageInputMatcher struct {
	caseValue *sqs.SendMessageInput
	xv        interface{}
}

func (q *MessageInputMatcher) String() string {
	return fmt.Sprintf("%v", q.caseValue)
}
func (q *MessageInputMatcher) Matches(x interface{}) bool {
	another, _ := x.(*sqs.SendMessageInput)
	matched := true
	matched = matched && *q.caseValue.MessageBody == *another.MessageBody
	matched = matched && *q.caseValue.MessageGroupId == *another.MessageGroupId
	matched = matched && *q.caseValue.QueueUrl == *another.QueueUrl
	return matched
}

func TestSqs_Enqueue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		caseMessage := enqueueevent.QueueMessage{
			UserID: "id1",
			Text:   "test text",
		}
		caseSQSUrl := "sqs"
		caseSendMessageInputBody := SqsMessageBody{
			UserID: caseMessage.UserID,
			Text:   caseMessage.Text,
		}
		caseSendMessageInputBodyJSON, _ := json.Marshal(caseSendMessageInputBody)
		caseMessageInput := &sqs.SendMessageInput{
			MessageBody:    aws.String(string(caseSendMessageInputBodyJSON)),
			MessageGroupId: aws.String(_messageGroupID),
			QueueUrl:       aws.String(caseSQSUrl),
		}
		caseMessageOutput := &sqs.SendMessageOutput{
			MessageId: aws.String("msgid1"),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("SQS_URL"), "").Return(caseSQSUrl)
		config.SetConfig(c)

		matcher := &MessageInputMatcher{
			caseValue: caseMessageInput,
		}
		w := NewMockSqsWrapper(ctrl)
		w.EXPECT().SendMessage(matcher).Return(caseMessageOutput, nil)

		q := NewSqs(w)
		r, err := q.Enqueue(caseMessage)
		assert.Equal(t, *caseMessageOutput.MessageId, r)
		assert.NoError(t, err)
	})

	t.Run("ng:failed", func(t *testing.T) {
		caseMessage := enqueueevent.QueueMessage{
			UserID: "id1",
		}
		caseSQSUrl := "sqs"
		caseSendMessageInputBody := SqsMessageBody{
			UserID: caseMessage.UserID,
			Text:   caseMessage.Text,
		}
		caseSendMessageInputBodyJSON, _ := json.Marshal(caseSendMessageInputBody)
		caseMessageInput := &sqs.SendMessageInput{
			MessageBody:    aws.String(string(caseSendMessageInputBodyJSON)),
			MessageGroupId: aws.String(_messageGroupID),
			QueueUrl:       aws.String(caseSQSUrl),
		}
		caseError := errors.New("error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("SQS_URL"), "").Return(caseSQSUrl)
		config.SetConfig(c)

		matcher := &MessageInputMatcher{
			caseValue: caseMessageInput,
		}
		w := NewMockSqsWrapper(ctrl)
		w.EXPECT().SendMessage(matcher).Return(nil, caseError)

		q := NewSqs(w)
		r, err := q.Enqueue(caseMessage)
		assert.Empty(t, r)
		assert.Equal(t, fmt.Errorf("failed to enqueue %w", caseError), err)
	})
}
