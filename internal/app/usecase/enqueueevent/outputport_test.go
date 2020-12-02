package enqueueevent

import (
	"github.com/golang/mock/gomock"
	"slacktimer/internal/app/util/log"
	"testing"
)

func TestCloudWatchLogsOutputPort_Output(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		caseData := &OutputData{}
		caseData.NotifiedUserIdList = make([]string, 0)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l := log.NewMockLogger(ctrl)
		l.EXPECT().Info("no user, so did not enqueue")
		log.SetDefaultLogger(l)

		o := &CloudWatchLogsOutputPort{}
		o.Output(caseData)
	})

	t.Run("exist", func(t *testing.T) {
		caseData := &OutputData{}
		caseData.NotifiedUserIdList = []string{
			"id1", "id2",
		}
		caseData.QueueMessageIdList = []string{
			"mid1", "mid2",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l := log.NewMockLogger(ctrl)
		l.EXPECT().Info("enqueue user_id=id1 message_id=mid1")
		l.EXPECT().Info("enqueue user_id=id2 message_id=mid2")
		log.SetDefaultLogger(l)

		o := &CloudWatchLogsOutputPort{}
		o.Output(caseData)
	})
}
