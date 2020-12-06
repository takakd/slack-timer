package updatetimerevent

import "slacktimer/internal/app/enterpriserule"

type OutputPort interface {
	Output(data OutputData)
}

type OutputData struct {
	Result     error
	SavedEvent *enterpriserule.TimerEvent
}
