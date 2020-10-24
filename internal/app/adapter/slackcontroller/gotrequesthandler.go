package slackcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/pkg/httputil"
	"slacktimer/internal/pkg/log"
)

// GotRequestHandler represents the API command "Got".
type GotRequestHandler struct {
	messageEvent *MessageEvent
	usecase      updatetimerevent.Usecase
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

func (g *GotRequestOutputPort) Output(data *updatetimerevent.OutputData) {
	err := data.Result
	errRaised := false
	if errors.Is(err, updatetimerevent.ErrFind) {
		errRaised = true
	} else if errors.Is(err, updatetimerevent.ErrCreate) {
		errRaised = true
	} else if errors.Is(err, updatetimerevent.ErrSave) {
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
		writeErrorCallbackResponse(g.w, body)
		return
	}

	httputil.WriteJsonResponse(g.w, nil, http.StatusOK, respBody)
	return
}
