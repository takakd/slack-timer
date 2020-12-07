package notify

import (
	"errors"
	"fmt"
	"slacktimer/internal/app/usecase/notifyevent"
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
	t.Run("ok", func(t *testing.T) {
		caseData := notifyevent.OutputData{
			UserID: "test user",
			Result: nil,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l := log.NewMockLogger(ctrl)
		l.EXPECT().Info(fmt.Sprintf("done notified user_id=%s", caseData.UserID))
		log.SetDefaultLogger(l)

		o := &CloudWatchLogsPresenter{}
		o.Output(caseData)
	})

	t.Run("ng", func(t *testing.T) {
		caseData := notifyevent.OutputData{
			UserID: "test user",
			Result: errors.New("error"),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l := log.NewMockLogger(ctrl)
		l.EXPECT().Error(fmt.Sprintf("notify user_id=%s: %v", caseData.UserID, caseData.Result))
		log.SetDefaultLogger(l)

		o := &CloudWatchLogsPresenter{}
		o.Output(caseData)
	})
}
