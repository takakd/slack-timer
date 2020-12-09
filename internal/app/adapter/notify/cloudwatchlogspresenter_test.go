package notify

import (
	"errors"
	"slacktimer/internal/app/usecase/notifyevent"
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
	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseData := notifyevent.OutputData{
			UserID: "test user",
			Result: nil,
		}

		ml := log.NewMockLogger(ctrl)
		ml.EXPECT().
			InfoWithContext(ac, "done notified", map[string]interface{}{
				"user_id": caseData.UserID,
			})
		log.SetDefaultLogger(ml)

		o := &CloudWatchLogsPresenter{}
		o.Output(appcontext.TODO(), caseData)
	})

	t.Run("ng", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseData := notifyevent.OutputData{
			UserID: "test user",
			Result: errors.New("error"),
		}

		ml := log.NewMockLogger(ctrl)
		ml.EXPECT().
			ErrorWithContext(ac, "notify", map[string]interface{}{
				"data": caseData,
			})
		log.SetDefaultLogger(ml)

		o := &CloudWatchLogsPresenter{}
		o.Output(appcontext.TODO(), caseData)
	})
}
