package enqueue

import (
	"fmt"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/log"
)

type CloudWatchLogsPresenter struct {
	Error error
}

func NewCloudWatchLogsPresenter() enqueueevent.OutputPort {
	return &CloudWatchLogsPresenter{}
}

func (c *CloudWatchLogsPresenter) Output(data enqueueevent.OutputData) {
	if len(data.NotifiedUserIdList) == 0 {
		log.Info("no items to be enqueued")
	} else {
		for i, v := range data.NotifiedUserIdList {
			log.Info(fmt.Sprintf("enqueued user_id=%s message_id=%s", v, data.QueueMessageIdList[i]))
		}
	}
}
