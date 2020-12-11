package settime

import (
	"net/http"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/usecase/timeroffevent"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/log"
)

// OffEventOutputReceivePresenter output logs to CloudWatchLogs.
type OffEventOutputReceivePresenter struct {
	Resp       Response
	Error      error
	SavedEvent *enterpriserule.TimerEvent
	StatusCode int
	Body       string
}

var _ timeroffevent.OutputPort = (*OffEventOutputReceivePresenter)(nil)

// NewOffEventOutputReceivePresenter creates new struct.
func NewOffEventOutputReceivePresenter() *OffEventOutputReceivePresenter {
	return &OffEventOutputReceivePresenter{}
}

// Output receives interactor outputs and keep them inside.
func (s *OffEventOutputReceivePresenter) Output(ac appcontext.AppContext, data timeroffevent.OutputData) {
	s.Resp = Response{
		Error: data.Result,
	}

	if data.Result != nil {
		log.ErrorWithContext(ac, "settime offevent outputport", data.Result)
		s.Resp.StatusCode = http.StatusInternalServerError
		s.Resp.Body = "internal server error"
		return
	}

	s.SavedEvent = data.SavedEvent
	s.Resp.StatusCode = http.StatusOK
	s.Resp.Body = "success"

}
