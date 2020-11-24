package notifycontroller

import (
	"fmt"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/pkg/log"
)

type CloudWatchLogsPresenter struct {
	Error error
}

func (p *CloudWatchLogsPresenter) Output(data *notifyevent.OutputData) {
	if data.Result == nil {
		log.Info(fmt.Sprintf("notified user_id=%s", data.UserId))
	} else {
		log.Error(fmt.Sprintf("failed to notify user_id=%s", data.UserId))
	}
}
