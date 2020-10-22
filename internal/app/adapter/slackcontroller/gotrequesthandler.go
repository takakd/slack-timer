package slackcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"proteinreminder/internal/app/usecase/updateproteinevent"
	"proteinreminder/internal/pkg/httputil"
	"proteinreminder/internal/pkg/log"
)

// GotRequestHandler represents the API command "Got".
type GotRequestHandler struct {
	messageEvent *MessageEvent
	usecase      updateproteinevent.Usecase
}

func (gr *GotRequestHandler) Handler(ctx context.Context, w http.ResponseWriter) {
	outputPort := &GotRequestOutputPort{w: w}
	// Update time to drink.
	gr.usecase.UpdateTimeToDrink(ctx, gr.messageEvent.User, outputPort)
	return
}

type GotRequestOutputPort struct {
	w http.ResponseWriter
}

func (g *GotRequestOutputPort) Output(data *updateproteinevent.OutputData) {
	err := data.Result
	errRaised := false
	if errors.Is(err, updateproteinevent.ErrFind) {
		errRaised = true
	} else if errors.Is(err, updateproteinevent.ErrCreate) {
		errRaised = true
	} else if errors.Is(err, updateproteinevent.ErrSave) {
		errRaised = true
	}

	resp := &SlackCallbackResponse{
		Message: "success",
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		errRaised = true
	}

	if errRaised {
		log.Error(err)
		body, err := makeErrorCallbackResponseBody("failed to save event", ErrSaveEvent)
		if err != nil {
			log.Error(err)
			body = []byte("internal error")
		}
		httputil.WriteJsonResponse(g.w, http.StatusBadRequest, body)
		return
	}

	httputil.WriteJsonResponse(g.w, http.StatusOK, respBody)
	return
}
