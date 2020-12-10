package enqueue

import (
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/log"
	"testing"

	"slacktimer/internal/app/util/appcontext"

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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseData := enqueueevent.OutputData{}
		caseData.NotifiedUserIDList = make([]string, 0)

		ml := log.NewMockLogger(ctrl)
		ml.EXPECT().InfoWithContext(ac, "no items to be enqueued")
		log.SetDefaultLogger(ml)

		o := &CloudWatchLogsPresenter{}
		o.Output(ac, caseData)
	})

	t.Run("exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseData := enqueueevent.OutputData{}
		caseData.NotifiedUserIDList = []string{
			"id1", "id2",
		}
		caseData.QueueMessageIDList = []string{
			"mid1", "mid2",
		}

		ml := log.NewMockLogger(ctrl)
		ml.EXPECT().InfoWithContext(ac, "enqueued", map[string]interface{}{
			"user_id":    caseData.NotifiedUserIDList[0],
			"message_id": caseData.QueueMessageIDList[0],
		})
		ml.EXPECT().InfoWithContext(ac, "enqueued", map[string]interface{}{
			"user_id":    caseData.NotifiedUserIDList[1],
			"message_id": caseData.QueueMessageIDList[1],
		})
		log.SetDefaultLogger(ml)

		o := &CloudWatchLogsPresenter{}
		o.Output(appcontext.TODO(), caseData)
	})
}
