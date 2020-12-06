package notify

import (
	"fmt"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/util/log"
)

// Output to CloudWatchLogs.
type CloudWatchLogsPresenter struct {
	Error error
}

func NewCloudWatchLogsPresenter() notifyevent.OutputPort {
	return &CloudWatchLogsPresenter{}
}

func (c CloudWatchLogsPresenter) Output(data notifyevent.OutputData) {
	if data.Result != nil {
		log.Error(fmt.Sprintf("notify user_id=%s: %v", data.UserId, data.Result))
		return
	}

	log.Info(fmt.Sprintf("done notified user_id=%s", data.UserId))
}
