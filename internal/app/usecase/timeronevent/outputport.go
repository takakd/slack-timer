package timeronevent

import (
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/appcontext"
)

// OutputPort defines outputport method of updating timer events usecase.
type OutputPort interface {
	Output(ac appcontext.AppContext, data OutputData)
}

// OutputData is parameter of OutputPort.
type OutputData struct {
	Result     error
	SavedEvent *enterpriserule.TimerEvent
}
