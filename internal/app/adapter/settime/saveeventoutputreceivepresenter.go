package settime

import (
	"fmt"
	"net/http"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/log"
)

// Pass the output data to the controller.
type SaveEventOutputReceivePresenter struct {
	Resp       Response
	Error      error
	SavedEvent *enterpriserule.TimerEvent
	StatusCode int
	Body       string
}

func NewSaveEventOutputReceivePresenter() *SaveEventOutputReceivePresenter {
	return &SaveEventOutputReceivePresenter{}
}

func (s *SaveEventOutputReceivePresenter) Output(data updatetimerevent.OutputData) {
	s.Resp = Response{
		Error: data.Result,
	}

	if data.Result != nil {
		log.Info(fmt.Sprintf("SaveEventOutputReceivePresenter.Output error=%v", data.Result))
		s.Resp.StatusCode = http.StatusInternalServerError
		s.Resp.Body = "internal server error"
		return
	}

	s.SavedEvent = data.SavedEvent
	s.Resp.StatusCode = http.StatusOK
	s.Resp.Body = "success"
}
