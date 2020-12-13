package settime

import (
	"net/http"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/usecase/timeronevent"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/log"
)

// OnEventOutputReceivePresenter output logs to CloudWatchLogs.
type OnEventOutputReceivePresenter struct {
	Resp       Response
	Error      error
	SavedEvent *enterpriserule.TimerEvent
	StatusCode int
	Body       string
}

var _ timeronevent.OutputPort = (*OnEventOutputReceivePresenter)(nil)

// NewOnEventOutputReceivePresenter creates new struct.
func NewOnEventOutputReceivePresenter() *OnEventOutputReceivePresenter {
	return &OnEventOutputReceivePresenter{}
}

// Output receives interactor outputs and keep them inside.
func (s *OnEventOutputReceivePresenter) Output(ac appcontext.AppContext, data timeronevent.OutputData) {
	s.Resp = Response{
		Error: data.Result,
	}

	if data.Result != nil {
		log.ErrorWithContext(ac, "settime offevent outputport", data.Result.Error())
		s.Resp.StatusCode = http.StatusInternalServerError
		s.Resp.Body = "internal server error"
		return
	}

	s.SavedEvent = data.SavedEvent
	s.Resp.StatusCode = http.StatusOK
	s.Resp.Body = "success"

}
