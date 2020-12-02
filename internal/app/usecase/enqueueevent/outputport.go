package enqueueevent

import (
	"fmt"
	"slacktimer/internal/app/util/log"
)

type CloudWatchLogsOutputPort struct {
	Error error
}

func NewCloudWatchLogsOutputPort() *CloudWatchLogsOutputPort {
	return &CloudWatchLogsOutputPort{}
}

func (c *CloudWatchLogsOutputPort) Output(data *OutputData) {
	if len(data.NotifiedUserIdList) == 0 {
		log.Info("no user, so did not enqueue")
	} else {
		for i, v := range data.NotifiedUserIdList {
			log.Info(fmt.Sprintf("enqueue user_id=%s message_id=%s", v, data.QueueMessageIdList[i]))
		}
	}
}
