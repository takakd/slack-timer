package notify

import (
	"fmt"
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
		log.Error(fmt.Sprintf("notify user_id=%s: %v", data.UserID, data.Result))
		return
	}

	log.Info(fmt.Sprintf("done notified user_id=%s", data.UserID))
}
