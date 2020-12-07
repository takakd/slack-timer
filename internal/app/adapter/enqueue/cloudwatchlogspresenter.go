package enqueue

import (
	"fmt"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/log"
)

// CloudWatchLogsPresenter output logs to CloudWatchLogs.
type CloudWatchLogsPresenter struct {
}

var _ enqueueevent.OutputPort = (*CloudWatchLogsPresenter)(nil)

// NewCloudWatchLogsPresenter create new struct.
func NewCloudWatchLogsPresenter() *CloudWatchLogsPresenter {
	return &CloudWatchLogsPresenter{}
}

// Output used as outputport by interactor.
func (c CloudWatchLogsPresenter) Output(data enqueueevent.OutputData) {
	if len(data.NotifiedUserIDList) == 0 {
		log.Info("no items to be enqueued")
		return
	}

	for i, v := range data.NotifiedUserIDList {
		log.Info(fmt.Sprintf("enqueued user_id=%s message_id=%s", v, data.QueueMessageIDList[i]))
	}
}
