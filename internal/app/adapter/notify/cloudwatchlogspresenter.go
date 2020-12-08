package notify

import (
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/util/log"
)

// CloudWatchLogsPresenter output logs to CloudWatchLogs.
type CloudWatchLogsPresenter struct {
}

var _ notifyevent.OutputPort = (*CloudWatchLogsPresenter)(nil)

// NewCloudWatchLogsPresenter create new struct.
func NewCloudWatchLogsPresenter() *CloudWatchLogsPresenter {
	return &CloudWatchLogsPresenter{}
}

// Output used as outputport by interactor.
func (c CloudWatchLogsPresenter) Output(data notifyevent.OutputData) {
	if data.Result != nil {
		log.Error("notify", map[string]interface{}{
			"data": data,
		})
		return
	}
	log.Info("done notified", map[string]interface{}{
		"user_id": data.UserID,
	})
}
