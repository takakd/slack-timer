package enqueue

import (
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/log"
)

// CloudWatchLogsPresenter output logs to CloudWatchLogs.
type CloudWatchLogsPresenter struct {
}

var _ enqueueevent.OutputPort = (*CloudWatchLogsPresenter)(nil)

// NewCloudWatchLogsPresenter creates new struct.
func NewCloudWatchLogsPresenter() *CloudWatchLogsPresenter {
	return &CloudWatchLogsPresenter{}
}

// Output used as outputport by interactor.
func (c CloudWatchLogsPresenter) Output(ac appcontext.AppContext, data enqueueevent.OutputData) {
	if len(data.NotifiedUserIDList) == 0 {
		log.InfoWithContext(ac, "no items to be enqueued")
		return
	}

	for i, v := range data.NotifiedUserIDList {
		log.InfoWithContext(ac, "enqueued", map[string]interface{}{
			"user_id":    v,
			"message_id": data.QueueMessageIDList[i],
		})
	}
}
