package settime

import (
	"net/http"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/log"
)

// SaveEventOutputReceivePresenter passes the output data to the controller.
type SaveEventOutputReceivePresenter struct {
	Resp       Response
	Error      error
	SavedEvent *enterpriserule.TimerEvent
	StatusCode int
	Body       string
}

var _ updatetimerevent.OutputPort = (*SaveEventOutputReceivePresenter)(nil)

// NewSaveEventOutputReceivePresenter create new struct.
func NewSaveEventOutputReceivePresenter() *SaveEventOutputReceivePresenter {
	return &SaveEventOutputReceivePresenter{}
}

// Output receives interactor outputs and keep them inside.
func (s *SaveEventOutputReceivePresenter) Output(ac appcontext.AppContext, data updatetimerevent.OutputData) {
	s.Resp = Response{
		Error: data.Result,
	}

	if data.Result != nil {
		log.ErrorWithContext(ac, "settime outputport", data.Result)
		s.Resp.StatusCode = http.StatusInternalServerError
		s.Resp.Body = "internal server error"
		return
	}

	s.SavedEvent = data.SavedEvent
	s.Resp.StatusCode = http.StatusOK
	s.Resp.Body = "success"
}
