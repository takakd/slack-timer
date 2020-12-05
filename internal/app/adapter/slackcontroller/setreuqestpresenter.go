package slackcontroller

import (
	"fmt"
	"net/http"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/log"
)

type LambdaResponseAdaptPresenter struct {
	Resp       *Response
	Error      error
	SavedEvent *enterpriserule.TimerEvent
	StatusCode int
	Body       string
}

func NewLambdaResponseAdaptPresenter() *LambdaResponseAdaptPresenter {
	return &LambdaResponseAdaptPresenter{}
}

func (p *LambdaResponseAdaptPresenter) Output(data updatetimerevent.OutputData) {
	p.Resp = &Response{
		Error: data.Result,
	}

	if data.Result != nil {
		log.Info(fmt.Sprintf("SetRequestOutputPort.Output error=%v", data.Result))
		p.Resp.StatusCode = http.StatusInternalServerError
		p.Resp.Body = "internal server error"
		return
	}

	p.SavedEvent = data.SavedEvent
	p.Resp.StatusCode = http.StatusOK
	p.Resp.Body = "success"
}
