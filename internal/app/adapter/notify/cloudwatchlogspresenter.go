package notify

import (
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/util/appcontext"
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
func (c CloudWatchLogsPresenter) Output(ac appcontext.AppContext, data notifyevent.OutputData) {
	if data.Result != nil {
		log.ErrorWithContext(ac, "notify", map[string]interface{}{
			"data": data,
		})
		return
	}
	log.InfoWithContext(ac, "done notified", map[string]interface{}{
		"user_id": data.UserID,
	})
}
