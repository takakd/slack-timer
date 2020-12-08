package enqueue

import (
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewCloudWatchLogsPresenter(t *testing.T) {
	assert.NotPanics(t, func() {
		NewCloudWatchLogsPresenter()
	})
}

func TestCloudWatchLogsPresenter_Output(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		caseData := enqueueevent.OutputData{}
		caseData.NotifiedUserIDList = make([]string, 0)

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
		caseData.NotifiedUserIDList = []string{
			"id1", "id2",
		}
		caseData.QueueMessageIDList = []string{
			"mid1", "mid2",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l := log.NewMockLogger(ctrl)
		l.EXPECT().Info("enqueued", map[string]interface{}{
			"user_id":    caseData.NotifiedUserIDList[0],
			"message_id": caseData.QueueMessageIDList[0],
		})
		l.EXPECT().Info("enqueued", map[string]interface{}{
			"user_id":    caseData.NotifiedUserIDList[1],
			"message_id": caseData.QueueMessageIDList[1],
		})
		log.SetDefaultLogger(l)

		o := &CloudWatchLogsPresenter{}
		o.Output(caseData)
	})
}
