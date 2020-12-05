package notify

import (
	"fmt"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/util/log"
)

type CloudWatchLogsPresenter struct {
	Error error
}

func NewCloudWatchLogsPresenter() notifyevent.OutputPort {
	return &CloudWatchLogsPresenter{}
}

func (p *CloudWatchLogsPresenter) Output(data notifyevent.OutputData) {
	if data.Result == nil {
		log.Info(fmt.Sprintf("done notified user_id=%s", data.UserId))
	} else {
		log.Error(fmt.Sprintf("notify user_id=%s: %v", data.UserId, data.Result))
	}
}
