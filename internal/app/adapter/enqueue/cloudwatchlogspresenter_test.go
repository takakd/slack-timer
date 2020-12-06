package enqueue

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/log"
	"testing"
)

func TestNewCloudWatchLogsPresenter(t *testing.T) {
	assert.NotPanics(t, func() {
		NewCloudWatchLogsPresenter()
	})
}

func TestCloudWatchLogsPresenter_Output(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		caseData := enqueueevent.OutputData{}
		caseData.NotifiedUserIdList = make([]string, 0)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l := log.NewMockLogger(ctrl)
		l.EXPECT().Info("no items to be enqueued")
		log.SetDefaultLogger(l)

		o := &CloudWatchLogsPresenter{}
		o.Output(caseData)
	})

	t.Run("exist", func(t *testing.T) {
		caseData := enqueueevent.OutputData{}
		caseData.NotifiedUserIdList = []string{
			"id1", "id2",
		}
		caseData.QueueMessageIdList = []string{
			"mid1", "mid2",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l := log.NewMockLogger(ctrl)
		l.EXPECT().Info("enqueued user_id=id1 message_id=mid1")
		l.EXPECT().Info("enqueued user_id=id2 message_id=mid2")
		log.SetDefaultLogger(l)

		o := &CloudWatchLogsPresenter{}
		o.Output(caseData)
	})
}
