package updatetimerevent

import "slacktimer/internal/app/enterpriserule"

// OutputPort defines outputport method of updating timer events usecase.
type OutputPort interface {
	Output(data OutputData)
}

// OutputData is parameter of OutputPort.
type OutputData struct {
	Result     error
	SavedEvent *enterpriserule.TimerEvent
}
