package notifycontroller

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/util/log"
	"testing"
)

func TestCloudWatchLogsPresenter_Output(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		caseData := &notifyevent.OutputData{
			UserId: "test user",
			Result: nil,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l := log.NewMockLogger(ctrl)
		l.EXPECT().Info(fmt.Sprintf("notified user_id=%s", caseData.UserId))
		log.SetDefaultLogger(l)

		o := &CloudWatchLogsPresenter{}
		o.Output(caseData)
	})

	t.Run("ng", func(t *testing.T) {
		caseData := &notifyevent.OutputData{
			UserId: "test user",
			Result: errors.New("error"),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l := log.NewMockLogger(ctrl)
		l.EXPECT().Error(fmt.Sprintf("failed to notify user_id=%s: %v", caseData.UserId, caseData.Result))
		log.SetDefaultLogger(l)

		o := &CloudWatchLogsPresenter{}
		o.Output(caseData)
	})
}
