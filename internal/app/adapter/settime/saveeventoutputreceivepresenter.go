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
	Resp       *Response
	Error      error
	SavedEvent *enterpriserule.TimerEvent
	StatusCode int
	Body       string
}

func NewSaveEventOutputReceivePresenter() *SaveEventOutputReceivePresenter {
	return &SaveEventOutputReceivePresenter{}
}

func (p *SaveEventOutputReceivePresenter) Output(data updatetimerevent.OutputData) {
	p.Resp = &Response{
		Error: data.Result,
	}

	if data.Result != nil {
		log.Info(fmt.Sprintf("SaveEventOutputReceivePresenter.Output error=%v", data.Result))
		p.Resp.StatusCode = http.StatusInternalServerError
		p.Resp.Body = "internal server error"
		return
	}

	p.SavedEvent = data.SavedEvent
	p.Resp.StatusCode = http.StatusOK
	p.Resp.Body = "success"
}
